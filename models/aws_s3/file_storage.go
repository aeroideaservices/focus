package s3

import (
	"context"
	"github.com/aeroideaservices/focus/models/plugin/actions"
	"github.com/aeroideaservices/focus/services/errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type fileStorage struct {
	bucket string
	client *s3.Client
}

// NewMediaStorage конструктор
func NewMediaStorage(client *s3.Client, bucket string) (*fileStorage, error) {
	_, err := client.HeadBucket(context.Background(), &s3.HeadBucketInput{Bucket: aws.String(bucket)})
	if err != nil {
		return nil, err
	}

	return &fileStorage{
		bucket: bucket,
		client: client,
	}, nil
}

// Upload загрузка файла
func (s fileStorage) Upload(ctx context.Context, media *actions.CreateFile) (err error) {
	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(media.Key),
		Body:        media.File,
		ContentType: aws.String(media.ContentType),
	})
	if err != nil {
		return errors.NoType.Wrap(err, "error uploading file")
	}

	return nil
}

// Delete удаление файла
func (s fileStorage) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	objectIdentifiers := make([]types.ObjectIdentifier, len(keys))
	for i, key := range keys {
		objectIdentifiers[i] = types.ObjectIdentifier{Key: aws.String(key)}
	}
	_, err := s.client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(s.bucket),
		Delete: &types.Delete{
			Objects: objectIdentifiers,
		},
	})

	if err != nil {
		return errors.NoType.Wrap(err, "error deleting file")
	}

	return nil
}

// GetSize получение размера файла
func (s fileStorage) GetSize(ctx context.Context, key string) (int64, error) {
	output, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return 0, errors.NoType.Wrap(err, "error getting file size")
	}

	return output.ContentLength, nil
}
