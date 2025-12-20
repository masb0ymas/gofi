package services

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Service struct {
	Client *s3.Client
}

func (s S3Service) CreateBucket(name string) (string, error) {
	ctx := context.Background()

	bucket, err := s.Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: &name,
	})
	if err != nil {
		return "", fmt.Errorf("error creating bucket: %s", err.Error())
	}

	fmt.Println("Bucket created:", *bucket.Location)

	return *bucket.Location, nil
}
