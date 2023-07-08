package minio_uploader

import (
	"context"
	"fmt"
	"log"

	"go-minio-sample/pkg/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	endpoint        string
	accessKeyID     string
	secretAccessKey string
	useSSL          bool
	bucketName      string
	location        string
	minioClient     *minio.Client
)

func Init() (err error) {
	endpoint = config.ConfigData.MinioEndpoint
	accessKeyID = config.ConfigData.AWSAccessKey
	secretAccessKey = config.ConfigData.AWSSecretAccessKey
	useSSL = true
	bucketName = config.ConfigData.AWSBucketName
	location = config.ConfigData.AWSLocation

	minioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	return err
}

func MakeBucket(ctx context.Context) (err error) {
	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			return err
		}
	}
	log.Printf("Successfully created %s\n", bucketName)
	return nil
}

func PutObject(objectName, filePath, contentType string) error {
	ctx := context.Background()
	err := Init()
	if err != nil {
		return err
	}

	err = MakeBucket(ctx)
	if err != nil {
		return err
	}

	info, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return err
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
	return nil
}

func ListObjects(prefix string) (objectData []minio.ObjectInfo, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = Init()
	if err != nil {
		return nil, err
	}

	objectCh := minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}
		fmt.Println(object)
		objectData = append(objectData, object)
	}

	return objectData, nil
}

func GetObject(objectName, filePath string) (err error) {
	err = minioClient.FGetObject(context.Background(), bucketName, objectName, filePath, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func CopyObject(ctx context.Context, src minio.CopySrcOptions, dst minio.CopyDestOptions) (uploadInfo minio.UploadInfo, err error) {
	// Copy object call
	uploadInfo, err = minioClient.CopyObject(context.Background(), dst, src)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully copied object:", uploadInfo)
	return uploadInfo, nil
}
