package main

import (
	"os"
	"testing"
)

func TestWriteFTYP(t *testing.T) {
	f, err := os.Create("test.mp4")
	if nil != err {
		t.Error(err)
	}
	defer f.Close()

	writeFTYP(f)
	writeMOOV(f, 1280, 720)
}
