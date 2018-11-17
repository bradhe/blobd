package aws

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	baseaws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/bradhe/blobd/blobs"
	"github.com/bradhe/blobd/iox"
	"github.com/bradhe/blobd/storage/managers"
)

type BlobManager struct {
	ctx    context.Context
	table  string
	region string
	bucket string
	prefix string
}

func (bm *BlobManager) path(id blobs.Id) string {
	prefix := bm.prefix

	if strings.HasSuffix(prefix, "/") {
		prefix = strings.TrimSuffix(prefix, "/")
	}

	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}

	return prefix + "/" + id.String()
}

func (bm *BlobManager) download(blob *blobs.Blob) error {
	sess := newAWSSession()
	svc := s3.New(sess)

	resp, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: baseaws.String(bm.bucket),
		Key:    baseaws.String(bm.path(blob.Id)),
	})

	if err != nil {
		return err
	}

	// This indicates what media type was originally presented to S3 by our
	// stack. Typically, it's encrypted data.
	var mediaType string

	if resp.ContentType != nil {
		mediaType = *resp.ContentType
	} else {
		mediaType = "application/octet-stream"
	}

	blob.Body = iox.MakeContentReader(mediaType, resp.Body)

	return nil
}

func (bm *BlobManager) upload(blob *blobs.Blob) error {
	sess := newAWSSession()
	svc := s3manager.NewUploader(sess)

	// this should technically be handled by other stuff in the stack but...
	exp := blob.Expiration()

	if exp.Before(time.Now()) {
		return managers.ErrInvalidExpiration
	}

	// we use a slighlty
	_, err := svc.Upload(&s3manager.UploadInput{
		Body:        blob.Body,
		Bucket:      baseaws.String(bm.bucket),
		Key:         baseaws.String(bm.path(blob.Id)),
		Expires:     baseaws.Time(exp),
		ContentType: baseaws.String(blob.Body.ContentType()),
	})

	return err
}

func parseIntAttribute(name string, out *dynamodb.UpdateItemOutput) int {
	if out == nil {
		return 0
	}

	if attr, ok := out.Attributes[name]; ok {
		if attr.N == nil {
			return 0
		}

		i, _ := strconv.Atoi(*attr.N)
		return i
	} else {
		return 0
	}
}

func getTTL() string {
	return fmt.Sprintf("%d", blobs.MaxExpirationFromNow().Unix())
}

func newUpdateItemInput(tableName string, id blobs.Id) *dynamodb.UpdateItemInput {
	return &dynamodb.UpdateItemInput{
		TableName: baseaws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"BlobId": {
				S: baseaws.String(id.String()),
			},
		},
		UpdateExpression: baseaws.String("SET TotalReads = if_not_exists(TotalReads, :zero) + :inc, BlobTTL = if_not_exists(BlobTTL, :ttl)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":inc":  {N: baseaws.String("1")},
			":zero": {N: baseaws.String("0")},
			":ttl":  {N: baseaws.String(getTTL())},
		},
		ReturnValues: baseaws.String("UPDATED_NEW"),
	}
}

func (bm *BlobManager) isDownloadable(id blobs.Id) (bool, error) {
	sess := newAWSSession()
	svc := dynamodb.New(sess)

	if out, err := svc.UpdateItem(newUpdateItemInput(bm.table, id)); err != nil {
		log.WithError(err).Error("UpdateItem failed")
		return false, err
	} else {
		if parseIntAttribute("TotalReads", out) > 1 {
			return false, nil
		}

		return true, nil
	}
}

func (bm *BlobManager) Get(id blobs.Id) (*blobs.Blob, error) {
	if ok, err := bm.isDownloadable(id); err != nil {
		return nil, err
	} else if !ok {
		return nil, managers.ErrNotFound
	}

	blob := blobs.Blob{
		Id: id,
	}

	if err := bm.download(&blob); err != nil {
		return nil, err
	}

	return &blob, nil
}

func (bm *BlobManager) Create(blob *blobs.Blob) error {
	blob.Id = blobs.NewId()

	if err := bm.upload(blob); err != nil {
		// destroy the id since it failed.
		blob.Id = blobs.Id{}
		return err
	}

	return nil
}

func (bm *BlobManager) Update(blob *blobs.Blob) error {
	return bm.upload(blob)
}
