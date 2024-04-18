package smil

import (
	"testing"
)

func TestSmil_createSmilFile(t *testing.T) {
	var f Command

	err := f.Exec("../../resources/test/output")

	if err != nil {
		t.Fatalf("Error while creating smil file err: %v", err)
	}
}
