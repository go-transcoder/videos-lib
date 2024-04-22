package s3

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
)

type S3BucketApi struct {
	S3Client   *s3.Client
	BucketName string
}

func (s3BucketApi S3BucketApi) DeleteFile(objectKey string) error {
	_, err := s3BucketApi.S3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s3BucketApi.BucketName),
		Key:    aws.String(objectKey),
	})

	if err != nil {
		return err
	}

	return nil
}

func (s3BucketApi S3BucketApi) DownloadFile(objectKey string, fileName string) error {
	result, err := s3BucketApi.S3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s3BucketApi.BucketName),
		Key:    aws.String(objectKey),
	})

	if err != nil {
		log.Printf("Couldn't get object %v:%v. Here's why: %v\n", s3BucketApi.BucketName, objectKey, err)
		return err
	}

	file, err := os.Create(fileName)

	if err != nil {
		log.Printf("Couldn't create file %v. Here's why: %v\n", fileName, err)
		return err
	}

	defer file.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Printf("Couldn't read object body from %v. Here's why: %v\n", objectKey, err)
	}
	_, err = file.Write(body)
	return err
}

func (s3BucketApi S3BucketApi) UploadVideoDir(videoDirPath string, prefix string) error {
	walker := make(fileWalk)
	go func() {
		// Gather the files to videos by walking the path recursively
		if err := filepath.Walk(videoDirPath, walker.Walk); err != nil {
			log.Fatalln("Walk failed:", err)
		}
		close(walker)
	}()

	for path := range walker {
		rel, err := filepath.Rel(videoDirPath, path)

		if err != nil {
			log.Fatalln("Unable to get relative path:", path, err)
		}
		file, err := os.Open(path)

		if err != nil {
			log.Println("Failed opening file", path, err)
			continue
		}

		defer file.Close()

		result, err := s3BucketApi.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(s3BucketApi.BucketName),
			Key:    aws.String(filepath.Join(prefix, rel)),
			Body:   file,
		})
		if err != nil {
			log.Fatalln("Failed to videos", path, err)
		}
		log.Println("Uploaded", path, result.ResultMetadata)
	}
	return nil
}

func (s3BucketApi S3BucketApi) UploadMultipart(file multipart.FileHeader, path string, fileName string) (string, error) {
	myFile, err := file.Open()

	if err != nil {
		return "", err
	}
	defer myFile.Close()

	// first videos to tmp
	tmpFile, err := os.CreateTemp("", fileName)

	if err != nil {
		return fileName, err
	}
	defer os.Remove(tmpFile.Name()) // Clean up

	// copy the multipart to the tmp file
	_, err = io.Copy(tmpFile, myFile)

	if err != nil {
		return "", err
	}

	destination := fmt.Sprintf("%s/%s", path, fileName)

	f, err := os.Open(tmpFile.Name())
	defer f.Close()

	_, err = s3BucketApi.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s3BucketApi.BucketName),
		Key:    aws.String(destination),
		Body:   f,
	})

	if err != nil {
		return fileName, err
	}

	return fileName, nil
}
