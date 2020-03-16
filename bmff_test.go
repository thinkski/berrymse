package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestWriteFTYP(t *testing.T) {
	dst, err := os.Create("testdata/acai.mp4")
	if nil != err {
		t.Error(err)
	}
	defer dst.Close()

	writeFTYP(dst)
	writeMOOV(dst, 1280, 720)

	src, err := ioutil.ReadFile("testdata/acai.264")
	if nil != err {
		t.Error(err)
	}

	nals := bytes.Split(src, []byte{0, 0, 0, 1})

	n := 1
	for _, nal := range nals {
		if len(nal) > 0 {
			if nal[0]&0x1f < 6 {
				writeMOOF(dst, n, nal)
				writeMDAT(dst, nal)
				n += 1
			}
		}
	}
}
