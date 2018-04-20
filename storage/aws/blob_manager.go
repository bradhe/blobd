package aws

import (
	"context"
	"strings"
	"time"

	baseaws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/bradhe/blobd/blobs"
	"github.com/bradhe/blobd/iox"
	"github.com/bradhe/blobd/storage/managers"
)

type BlobManager struct {
	ctx    context.Context
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

	// AWS uses the Expires field to track the expiration of the object. The
	// format conforms to RFC1123. This is based on anecdotal (e.g. testing)
	// evidence and I couldn't find anything in the docs to support it.
	if resp.Expires != nil {
		if t, err := time.Parse(time.RFC1123, *resp.Expires); err == nil {
			blob.ExpiresAt = t
		} else {
			log.WithError(err).Error("failed to parse expiration from Amazon")
		}
	}

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

func (bm *BlobManager) Get(id blobs.Id) (*blobs.Blob, error) {
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
