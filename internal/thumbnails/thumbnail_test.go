package thumbnails

import (
	"testing"
)

func TestExtractThumbs_CreateThumbs(t *testing.T) {
	var extractThumbs ExtractThumbs

	err := extractThumbs.Exec("../../resources/test/sample_test.mp4", "../../resources/test/output")

	if err != nil {
		t.Fatalf("Error while extracting images err: %v", err)
	}
}
