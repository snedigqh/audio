package wav_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mattetti/audio/wav"
)

func TestClip_Read(t *testing.T) {
	expectations := []struct {
		input       string
		totalFrames int
	}{
		{"fixtures/kick.wav", 4484},
	}

	for _, exp := range expectations {
		path, _ := filepath.Abs(exp.input)
		f, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		d := wav.NewDecoder(f)
		clip := d.Clip()
		totalFrames := int(clip.Size())
		if totalFrames != exp.totalFrames {
			t.Fatalf("Expected %d frames, got %d\n", exp.totalFrames, totalFrames)
		}
		readFrames := 0

		bufSize := 4096
		buf := make([]byte, bufSize)
		var n int
		for readFrames < totalFrames {
			n, err = clip.Read(buf)
			if err != nil || n == 0 {
				break
			}
			readFrames += n
		}
		if readFrames != totalFrames {
			t.Fatalf("file expected to have %d frames, only read %d, off by %d frames\n", totalFrames, readFrames, (totalFrames - readFrames))
		}

	}
}