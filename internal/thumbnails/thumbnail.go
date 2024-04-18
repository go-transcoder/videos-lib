package thumbnails

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/tidwall/gjson"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"io"
	"os"
)

type ExtractThumbs func(inFileName string, outputDir string) error

func (extractThumbs ExtractThumbs) Exec(inFileName, outputDir string) error {
	a, err := ffmpeg_go.Probe(inFileName)

	if err != nil {
		return err
	}

	totalDuration := gjson.Get(a, "format.duration").Float()

	screenshotTime := make([]float64, 0)
	getScreenshotSeconds(totalDuration, &screenshotTime)

	// make sure that the thumbnails dir is created
	thumbnailDir := fmt.Sprintf("%s/thumbnails", outputDir)

	_, err = os.Stat(thumbnailDir)

	if os.IsNotExist(err) {
		err := os.Mkdir(thumbnailDir, 0755)
		if err != nil {
			return err
		}
	}

	for i, v := range screenshotTime {
		reader := readFrameAsJpeg(inFileName, int(v))

		img, err := imaging.Decode(reader)
		if err != nil {
			return err
		}
		err = imaging.Save(img, fmt.Sprintf("%s/thumbnails/thumb%v.jpeg", outputDir, i+1))
		if err != nil {
			return err
		}
	}

	return nil
}

func getScreenshotSeconds(totalDuration float64, screenshotTimes *[]float64) {
	// Getting the screenshot times
	*screenshotTimes = []float64{0, 3}

	if totalDuration > 60 {
		for time := float64(60); totalDuration/time > 1; time += 60 {
			*screenshotTimes = append(*screenshotTimes, time)
		}
	} else {
		*screenshotTimes = append(*screenshotTimes, totalDuration/2)
	}
	*screenshotTimes = append(*screenshotTimes, totalDuration-3)
}

func readFrameAsJpeg(inFileName string, frameNum int) io.Reader {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg_go.Input(inFileName).
		Filter("select", ffmpeg_go.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg_go.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		panic(err)
	}
	return buf
}
