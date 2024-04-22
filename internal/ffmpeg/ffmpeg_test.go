package ffmpeg

import (
	_ "embed"
	"os"
	"testing"
)

//go:embed convert_video_cpu.sh
var convertScript string

func TestFfmpegApi_FfmpegCommandExec(t *testing.T) {
	var ffmpegApi Command

	outputPath := "../../resources/test/output"

	err := ffmpegApi.Exec(convertScript, "../../resources/test/sample_test.mp4", outputPath)

	if err != nil {
		t.Fatalf("ffmpegCommandExec error: , %v", err)
	}

	// the output file should exist
	_, err = os.Stat(outputPath)

	if os.IsNotExist(err) {
		t.Fatalf("File %s is expected to be in path. error: %v", outputPath, err)
	}
}
