package s3

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/pkg/errors"

	"github.com/aeroideaservices/focus/media/plugin/actions"
)

type fileStorage struct {
	bucket string
	client *s3.Client
}

// NewFileStorage конструктор
func NewFileStorage(client *s3.Client, bucket string) (*fileStorage, error) {
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
func (s fileStorage) Upload(ctx context.Context, media *actions.UploadFile) (err error) {
	_, err = s.client.PutObject(
		ctx, &s3.PutObjectInput{
			Bucket:      aws.String(s.bucket),
			Key:         aws.String(media.Key),
			Body:        media.File,
			ContentType: aws.String(media.ContentType),
		},
	)

	return errors.WithStack(err)
}

// UploadList загрузка нескольких файлов
func (s fileStorage) UploadList(ctx context.Context, medias ...actions.UploadFile) (err error) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	tasks := make(chan *actions.UploadFile)
	go func(tasks chan *actions.UploadFile, wg *sync.WaitGroup) {
		defer wg.Done()
		for media := range tasks {
			err = s.Upload(ctx, media)
			if err != nil {
				//break
			}
		}
	}(tasks, wg)

	for _, media := range medias {
		tasks <- &media
	}
	close(tasks)
	wg.Wait()

	return err
}

// Move перемещение файла
func (s fileStorage) Move(ctx context.Context, oldKey string, newKey string) error {
	_, err := s.client.CopyObject(
		ctx, &s3.CopyObjectInput{
			Bucket:     aws.String(s.bucket),
			CopySource: aws.String(fmt.Sprintf("%s/%s", s.bucket, oldKey)),
			Key:        aws.String(newKey),
		},
	)
	if err != nil {
		return err
	}

	err = s.Delete(ctx, oldKey)

	return errors.WithStack(err)
}

// Delete удаление файлов
func (s fileStorage) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	objectIdentifiers := make([]types.ObjectIdentifier, len(keys))
	for i, key := range keys {
		objectIdentifiers[i] = types.ObjectIdentifier{Key: aws.String(key)}
	}
	_, err := s.client.DeleteObjects(
		ctx, &s3.DeleteObjectsInput{
			Bucket: aws.String(s.bucket),
			Delete: &types.Delete{
				Objects: objectIdentifiers,
			},
		},
	)

	return errors.WithStack(err)
}

// GetSize получение размера файла
func (s fileStorage) GetSize(ctx context.Context, key string) (int64, error) {
	output, err := s.client.GetObject(
		ctx, &s3.GetObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(key),
		},
	)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return output.ContentLength, nil
}

// GetSizeDeprecated deprecated
func (s fileStorage) GetSizeDeprecated(ctx context.Context, key string) (int64, error) {
	output, err := s.client.GetObjectAttributes(
		ctx, &s3.GetObjectAttributesInput{
			Bucket:           aws.String(s.bucket),
			Key:              aws.String(key),
			ObjectAttributes: []types.ObjectAttributes{types.ObjectAttributesObjectSize},
		},
	)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return output.ObjectSize, nil
}

func (s fileStorage) DownloadFile(ctx context.Context, key string, fileName string) error {
	result, err := s.client.GetObject(
		ctx, &s3.GetObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(key),
		},
	)
	if err != nil {
		log.Printf("Couldn't get object %v:%v. Here's why: %v\n", s.bucket, key, err)
		return err
	}
	defer result.Body.Close()
	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("Couldn't create file %v. Here's why: %v\n", fileName, err)
		return err
	}
	defer file.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Printf("Couldn't read object body from %v. Here's why: %v\n", key, err)
	}
	_, err = file.Write(body)
	return err
}
