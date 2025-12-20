package main

import (
	"context"
	"log"

	"gofi/internal/config"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func newS3Client(cfg config.ConfigS3) *s3.Client {
	ctx := context.Background()

	// Load AWS configuration with static credentials
	awsCfg, err := awsConfig.LoadDefaultConfig(ctx,
		awsConfig.WithRegion(cfg.Region),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.ClientID,
			cfg.ClientSecret,
			"", // session token (empty for static credentials)
		)),
	)

	if err != nil {
		log.Fatal("failed to load AWS config", "error", err.Error())
	}

	// Create S3 client
	client := s3.NewFromConfig(awsCfg)

	return client
}
