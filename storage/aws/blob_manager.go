package aws

import (
	"context"
	"strings"

	baseaws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/bradhe/blobd/blobs"
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
	sess := getAWSSession()
	svc := s3.New(sess)

	resp, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: baseaws.String(bm.bucket),
		Key:    baseaws.String(bm.path(blob.Id)),
	})

	if err != nil {
		return err
	}

	blob.Body = resp.Body
	return nil
}

func (bm *BlobManager) upload(blob *blobs.Blob) error {
	sess := getAWSSession()
	svc := s3manager.NewUploader(sess)

	// we use a slighlty
	_, err := svc.Upload(&s3manager.UploadInput{
		Body:   blob.Body,
		Bucket: baseaws.String(bm.bucket),
		Key:    baseaws.String(bm.path(blob.Id)),
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
