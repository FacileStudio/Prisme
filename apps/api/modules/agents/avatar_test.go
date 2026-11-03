package agents

import (
	"bytes"
	"image/png"
	"testing"
)

func TestBlobPNGIsASquareWithTheBlobDrawn(t *testing.T) {
	data, err := blobPNG("Pixel")
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("not a PNG: %v", err)
	}
	bounds := decoded.Bounds()
	if bounds.Dx() != blobSize || bounds.Dy() != blobSize {
		t.Fatalf("size = %dx%d, want %dx%d", bounds.Dx(), bounds.Dy(), blobSize, blobSize)
	}
	_, _, _, alpha := decoded.At(blobSize/2, blobSize/2).RGBA()
	if alpha == 0 {
		t.Error("the centre pixel is transparent: no blob was drawn")
	}
}
