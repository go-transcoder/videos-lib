package ffmpeg

import (
	"bytes"
	"fmt"
	"os/exec"
)

type Command func(path, inputFile, outputDir string) error

func (ffmpegApi Command) Exec(inputFile, outputDir string) error {
	cmd := exec.Command(
		"/bin/sh",
		"./convert_video_cpu.sh",
		inputFile,
		outputDir,
	)
	// Create buffers to capture output
	var stdout, stderr bytes.Buffer

	// Set the output buffers for the command
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command
	err := cmd.Run()

	// Check for errors
	if err != nil {
		fmt.Println("ffmpeg stderr:", stderr.String())
		return err
	}

	// Print standard output
	fmt.Println("ffmpeg stdout:", stdout.String())

	return nil
}
