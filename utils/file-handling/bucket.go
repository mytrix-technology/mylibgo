package filestore

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type FileHandler struct {
	FileName    string
	ContentType string
	Size        int64
	ContentMD5  string
}

type FileStore interface {
	DownloadToFile(src string, dstPath string) error
	Download(src string, w io.WriterAt) (n int64, err error)
	UploadWithContext(ctx context.Context, file io.Reader, fileHandler FileHandler) (string, error)
}

type AWSBucket struct {
	accessKeyID string
	secretKey   string
	BucketName  string
	AWSEndpoint string
	AWSRegion   string
}

func NewAWSBucket(name, awsEndpoint, accessKeyID, secretKey, awsRegion string) *AWSBucket {
	return &AWSBucket{
		accessKeyID: accessKeyID,
		secretKey:   secretKey,
		BucketName:  name,
		AWSEndpoint: awsEndpoint,
		AWSRegion:   awsRegion,
	}
}

func (b *AWSBucket) DownloadToFile(src string, dstPath string) error {
	outfile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer outfile.Close()

	creds := credentials.NewStaticCredentials(b.accessKeyID, b.secretKey, "")
	sess, _ := session.NewSession(&aws.Config{
		Credentials: creds,
		Endpoint:    aws.String(b.AWSEndpoint),
		Region:      aws.String(b.AWSRegion),
	})

	downloader := s3manager.NewDownloader(sess)

	if _, err := downloader.Download(outfile, &s3.GetObjectInput{
		Bucket: aws.String(b.BucketName),
		Key:    aws.String(src),
	}); err != nil {
		return err
	}

	return nil
}

func (b *AWSBucket) Download(src string, w io.WriterAt) (n int64, err error) {
	creds := credentials.NewStaticCredentials(b.accessKeyID, b.secretKey, "")
	sess, _ := session.NewSession(&aws.Config{
		Credentials: creds,
		Endpoint:    aws.String(b.AWSEndpoint),
		Region:      aws.String(b.AWSRegion),
	})

	downloader := s3manager.NewDownloader(sess)

	return downloader.Download(w, &s3.GetObjectInput{
		Bucket: aws.String(b.BucketName),
		Key:    aws.String(src),
	})
}

func (b *AWSBucket) UploadWithContext(ctx context.Context, file io.Reader, handler FileHandler) (string, error) {
	var err error

	if handler.FileName == "" {
		return "", fmt.Errorf("filename cannot be zero string")
	}
	bucket := b.BucketName

	creds := credentials.NewStaticCredentials(b.accessKeyID, b.secretKey, "")
	sess, _ := session.NewSession(&aws.Config{
		Credentials: creds,
		Endpoint:    aws.String(b.AWSEndpoint),
		Region:      aws.String(endpoints.ApSoutheast2RegionID),
	})

	uploader := s3manager.NewUploader(sess)

	input := &s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(handler.FileName),
		Body:   file,
		ACL:    aws.String("public-read"),
	}

	if handler.ContentType != "" {
		input.ContentType = &handler.ContentType
	}

	if handler.ContentMD5 != "" {
		input.ContentMD5 = &handler.ContentMD5
	}

	out, err := uploader.UploadWithContext(ctx, input)
	if err != nil {
		return "", err
	}

	return out.Location, nil
}
