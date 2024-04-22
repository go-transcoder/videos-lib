package videos_lib

import (
	_ "embed"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-transcoder/videos-lib/internal/ffmpeg"
	s3_internal "github.com/go-transcoder/videos-lib/internal/s3"
	"github.com/go-transcoder/videos-lib/internal/smil"
	"github.com/go-transcoder/videos-lib/internal/thumbnails"
	"mime/multipart"
)

func UploadMultipart(cfg aws.Config, bucketName string, file multipart.FileHeader, path string, name string) (filename string, err error) {
	// TODO: the creation of the S3BucketApi should be a singleton
	// 	create a method that gets the S3BucketApi if already created
	s3BucketApi := s3_internal.S3BucketApi{
		S3Client:   s3.NewFromConfig(cfg),
		BucketName: bucketName,
	}

	return s3BucketApi.UploadMultipart(file, path, name)
}

func DownloadFile(cfg aws.Config, bucketName string, objectKey string, fileName string) error {
	s3BucketApi := s3_internal.S3BucketApi{
		S3Client:   s3.NewFromConfig(cfg),
		BucketName: bucketName,
	}

	return s3BucketApi.DownloadFile(objectKey, fileName)
}

func UploadVideoDir(cfg aws.Config, bucketName string, videoDirPath string, prefix string) error {
	s3BucketApi := s3_internal.S3BucketApi{
		S3Client:   s3.NewFromConfig(cfg),
		BucketName: bucketName,
	}

	return s3BucketApi.UploadVideoDir(videoDirPath, prefix)
}

//go:embed internal/ffmpeg/convert_video_cpu.sh
var convertScript string

func FfmpegTranscode(inputFile string, outputDir string) error {
	var ffmpegApi ffmpeg.Command

	return ffmpegApi.Exec(convertScript, inputFile, outputDir)
}

func ExtractThumbnails(inFileName, outputDir string) error {
	var extractThumbs thumbnails.ExtractThumbs

	return extractThumbs.Exec(inFileName, outputDir)
}

func CreateSmil(outputPath string) error {
	var smilApi smil.Command

	return smilApi.Exec(outputPath)
}
