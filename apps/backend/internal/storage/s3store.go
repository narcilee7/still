package storage

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
	"github.com/still-mvp/still/apps/backend/pkg/config"
)

// S3Store stores uploads in S3 or any S3-compatible service (MinIO, R2).
type S3Store struct {
	cfg      *config.Config
	client   *s3.Client
	presign  *s3.PresignClient
	bucket   string
	endpoint string
}

// NewS3Store creates a new S3-backed store.
func NewS3Store(cfg *config.Config) (*S3Store, error) {
	awsCfg, err := loadAWSConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("load aws config failed: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.S3Endpoint != "" {
			o.UsePathStyle = true
		}
	})

	store := &S3Store{
		cfg:      cfg,
		client:   client,
		presign:  s3.NewPresignClient(client),
		bucket:   cfg.S3Bucket,
		endpoint: cfg.S3Endpoint,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := store.ensureBucket(ctx); err != nil {
		return nil, fmt.Errorf("ensure bucket failed: %w", err)
	}

	if cfg.S3Endpoint != "" {
		if err := store.setPublicReadPolicy(ctx); err != nil {
			return nil, fmt.Errorf("set bucket policy failed: %w", err)
		}
	}

	return store, nil
}

// GenerateUploadURL returns a presigned PUT URL and the future public URL.
func (s *S3Store) GenerateUploadURL(ctx context.Context, filename, contentType string) (string, string, error) {
	key := "uploads/" + uuid.New().String() + path.Ext(filename)

	req, err := s.presign.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
		ACL:         types.ObjectCannedACLPublicRead,
	}, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		return "", "", fmt.Errorf("presign put object failed: %w", err)
	}

	return req.URL, s.PublicURL(ctx, key), nil
}

// PublicURL returns the public URL for a stored key.
func (s *S3Store) PublicURL(ctx context.Context, key string) string {
	if s.endpoint != "" {
		return strings.TrimSuffix(s.endpoint, "/") + "/" + s.bucket + "/" + key
	}
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.cfg.S3Region, key)
}

func loadAWSConfig(cfg *config.Config) (aws.Config, error) {
	opts := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(cfg.S3Region),
	}

	if cfg.S3Endpoint != "" {
		opts = append(opts,
			awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				cfg.S3AccessKeyID,
				cfg.S3SecretAccessKey,
				"",
			)),
			awsconfig.WithEndpointResolverWithOptions(
				aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:               cfg.S3Endpoint,
						HostnameImmutable: true,
					}, nil
				}),
			),
		)
	}

	return awsconfig.LoadDefaultConfig(context.Background(), opts...)
}

func (s *S3Store) ensureBucket(ctx context.Context) error {
	_, err := s.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(s.bucket),
	})
	if err != nil {
		var owned *types.BucketAlreadyOwnedByYou
		var exists *types.BucketAlreadyExists
		if errors.As(err, &owned) || errors.As(err, &exists) {
			return nil
		}
		return err
	}
	return nil
}

func (s *S3Store) setPublicReadPolicy(ctx context.Context) error {
	policy := fmt.Sprintf(`{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": "*",
      "Action": ["s3:GetObject"],
      "Resource": ["arn:aws:s3:::%s/*"]
    }
  ]
}`, s.bucket)

	_, err := s.client.PutBucketPolicy(ctx, &s3.PutBucketPolicyInput{
		Bucket: aws.String(s.bucket),
		Policy: aws.String(policy),
	})
	return err
}

