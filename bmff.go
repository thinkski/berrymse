package main

import (
	"bytes"
	"io"
)

func writeInt(w io.Writer, v int, n int) {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[n-i-1] = byte(v & 0xff)
		v >>= 8
	}
	w.Write(b)
}

func writeString(w io.Writer, s string) {
	w.Write([]byte(s))
}

func writeTag(w io.Writer, tag string, cb func(w io.Writer)) {
	var b bytes.Buffer
	cb(&b)                    // callback
	writeInt(w, b.Len()+8, 4) // box size
	writeString(w, tag)       // box type
	w.Write(b.Bytes())        // box content
}

func writeFTYP(w io.Writer) {
	writeTag(w, "ftyp", func(w io.Writer) {
		writeString(w, "isom")     // major brand
		writeInt(w, 1, 4)          // minor version
		writeString(w, "isomavc1") // compatible brands
	})
}

func writeMOOV(w io.Writer, width, height uint16) {
	writeTag(w, "moov", func(w io.Writer) {
		writeMVHD(w)
		writeTRAK(w, width, height)
		writeMVEX(w)
	})
}

func writeMVHD(w io.Writer) {
	writeTag(w, "mvhd", func(w io.Writer) {
		writeInt(w, 0, 4)          // version and flags
		writeInt(w, 0, 4)          // creation time
		writeInt(w, 0, 4)          // modification time
		writeInt(w, 1000, 4)       // timescale
		writeInt(w, -1, 4)         // duration (all 1s == live)
		writeInt(w, 0x00010000, 4) // rate (1.0 == normal)
		writeInt(w, 0x0100, 2)     // volume (1.0 == normal)
		writeInt(w, 0, 2)          // reserved
		writeInt(w, 0, 4)          // reserved
		writeInt(w, 0, 4)          // reserved
		writeInt(w, 0x00010000, 4) // matrix
		writeInt(w, 0x0, 4)        // matrix
		writeInt(w, 0x0, 4)        // matrix
		writeInt(w, 0x0, 4)        // matrix
		writeInt(w, 0x00010000, 4) // matrix
		writeInt(w, 0x0, 4)        // matrix
		writeInt(w, 0x0, 4)        // matrix
		writeInt(w, 0x0, 4)        // matrix
		writeInt(w, 0x40000000, 4) // matrix
		writeInt(w, 0, 4)          // pre-defined
		writeInt(w, 0, 4)          // pre-defined
		writeInt(w, 0, 4)          // pre-defined
		writeInt(w, 0, 4)          // pre-defined
		writeInt(w, 0, 4)          // pre-defined
		writeInt(w, 0, 4)          // pre-defined
		writeInt(w, 1, 4)          // next track id
	})
}

func writeTRAK(w io.Writer, width, height uint16) {
	writeTag(w, "trak", func(w io.Writer) {
		writeTKHD(w, width, height)
		writeMDIA(w, width, height)
	})
}

func writeTKHD(w io.Writer, width, height uint16) {
	writeTag(w, "tkhd", func(w io.Writer) {
		writeInt(w, 3, 4)           // version and flags (track enabled)
		writeInt(w, 0, 4)           // creation time
		writeInt(w, 0, 4)           // modification time
		writeInt(w, 1, 4)           // track id
		writeInt(w, 0, 4)           // reserved
		writeInt(w, -1, 4)          // duration (all 1s == live)
		writeInt(w, 0, 4)           // reserved
		writeInt(w, 0, 4)           // reserved
		writeInt(w, 0, 2)           // layer
		writeInt(w, 0, 2)           // alternate group
		writeInt(w, 0, 2)           // volume (ignored for video tracks)
		writeInt(w, 0, 2)           // reserved
		writeInt(w, 0x00010000, 4)  // matrix
		writeInt(w, 0x0, 4)         // matrix
		writeInt(w, 0x0, 4)         // matrix
		writeInt(w, 0x0, 4)         // matrix
		writeInt(w, 0x00010000, 4)  // matrix
		writeInt(w, 0x0, 4)         // matrix
		writeInt(w, 0x0, 4)         // matrix
		writeInt(w, 0x0, 4)         // matrix
		writeInt(w, 0x40000000, 4)  // matrix
		writeInt(w, int(width), 4)  // width
		writeInt(w, int(height), 4) // height
	})
}

func writeMDIA(w io.Writer, width, height uint16) {
	writeTag(w, "mdia", func(w io.Writer) {
		writeMDHD(w)
		writeHDLR(w)
		writeMINF(w, width, height)
	})
}

func writeMDHD(w io.Writer) {
	writeTag(w, "mdhd", func(w io.Writer) {
		writeInt(w, 0, 4)     // version and flags
		writeInt(w, 0, 4)     // creation time
		writeInt(w, 0, 4)     // modification time
		writeInt(w, 90000, 4) // timescale
		writeInt(w, -1, 4)    // duration (all 1s == live)
		writeInt(w, 0, 2)     // language
		writeInt(w, 0, 2)     // pre-defined
	})
}

func writeHDLR(w io.Writer) {
	writeTag(w, "hdlr", func(w io.Writer) {
		writeInt(w, 0, 4)              // version and flags
		writeInt(w, 0, 4)              // pre-defined
		writeString(w, "vide")         // handler type
		writeInt(w, 0, 4)              // reserved
		writeInt(w, 0, 4)              // reserved
		writeInt(w, 0, 4)              // reserved
		writeString(w, "VideoHandler") // name
	})
}

func writeMINF(w io.Writer, width, height uint16) {
	writeTag(w, "minf", func(w io.Writer) {
		writeVMHD(w)
		writeDINF(w)
		writeSTBL(w, width, height)
	})
}

func writeVMHD(w io.Writer) {
	writeTag(w, "vmhd", func(w io.Writer) {
		writeInt(w, 1, 4) // version and flags
		writeInt(w, 0, 2) // graphics mode
		writeInt(w, 0, 2) // opcolor
		writeInt(w, 0, 2) // opcolor
		writeInt(w, 0, 2) // opcolor
	})
}

func writeDINF(w io.Writer) {
	writeTag(w, "dinf", func(w io.Writer) {
		writeDREF(w)
	})
}

func writeDREF(w io.Writer) {
	writeTag(w, "dref", func(w io.Writer) {
		writeInt(w, 0, 4) // version and flags
		writeInt(w, 1, 4) // entry count
		writeTag(w, "url ", func(w io.Writer) {
			writeInt(w, 1, 4) // version and flags
		})
	})
}

func writeSTBL(w io.Writer, width, height uint16) {
	writeTag(w, "stbl", func(w io.Writer) {
		writeSTSD(w, width, height)
		writeSTTS(w)
		writeSTSC(w)
		writeSTCO(w)
	})
}

func writeSTSD(w io.Writer, width, height uint16) {
	writeTag(w, "stsd", func(w io.Writer) {
		writeInt(w, 0, 4) // version and flags
		writeInt(w, 1, 4) // entry count
		writeTag(w, "avc1", func(w io.Writer) {
			writeInt(w, 0, 4)           // reserved
			writeInt(w, 0, 2)           // reserved
			writeInt(w, 1, 2)           // data-ref index
			writeInt(w, 0, 2)           // pre-defined
			writeInt(w, 0, 2)           // reserved
			writeInt(w, 0, 4)           // pre-defined
			writeInt(w, 0, 4)           // pre-defined
			writeInt(w, 0, 4)           // pre-defined
			writeInt(w, int(width), 2)  // width
			writeInt(w, int(height), 2) // height
			writeInt(w, 0x00480000, 4)  // horizontal resolution 72dpi
			writeInt(w, 0x00480000, 4)  // vertical resolution 72dpi
			writeInt(w, 0, 4)           // data size (= 0)
			writeInt(w, 1, 2)           // frame count (= 1)
			w.Write(make([]byte, 32))   // compressor name
			writeInt(w, 0x18, 2)        // depth
			writeInt(w, 0xffff, 2)      // pre-defined

			writeTag(w, "avcC", func(w io.Writer) {
				w.Write([]byte{
					//	0x64, 0x00, 0x28, 0xac, 0x2b, 0x40, 0x3c, 0x01, 0x13, 0xf2, 0xc0,
					//	0x3c, 0x48, 0x9a, 0x80,
					//	0x28, 0xee, 0x02, 0x5c, 0xb0,
					0x01, 0x4d, 0x40, 0x28, 0xff,
					0xe1, 0x00, 0x16, 0x27, 0x4d, 0x40, 0x28, 0xa9, 0x18, 0x0f, 0x00, 0x44, 0xfc, 0xb8, 0x03, 0x50,
					0x10, 0x10, 0x1b, 0x6c, 0x2b, 0x5e, 0xf7, 0xc0, 0x40, 0x01, 0x00, 0x04, 0x28, 0xde, 0x09, 0xc8,
				})
				// write PPS and SPS data here
				//                            01 4d40 28ff  ......-avcC.M@(.
				// e100 1627 4d40 28a9 180f 0044 fcb8 0350  ...'M@(....D...P
				// 1010 1b6c 2b5e f7c0 4001 0004 28de 09c8
			})
		})
	})
}

func writeSTTS(w io.Writer) {
	writeTag(w, "stts", func(w io.Writer) {
		writeInt(w, 0, 4) // version and flags
		writeInt(w, 0, 4) // entry count
	})
}

func writeSTSC(w io.Writer) {
	writeTag(w, "stsc", func(w io.Writer) {
		writeInt(w, 0, 4) // version and flags
		writeInt(w, 0, 4) // entry count
	})
}

func writeSTCO(w io.Writer) {
	writeTag(w, "stsc", func(w io.Writer) {
		writeInt(w, 0, 4) // version and flags
		writeInt(w, 0, 4) // entry count
	})
}

func writeMVEX(w io.Writer) {
	writeTag(w, "mvex", func(w io.Writer) {
		writeTREX(w)
	})
}

func writeTREX(w io.Writer) {
	writeTag(w, "trex", func(w io.Writer) {
		writeInt(w, 0, 4) // version and flags
		writeInt(w, 1, 4) // track id
		writeInt(w, 1, 4) // default sample description index
		writeInt(w, 0, 4) // default sample duration
		writeInt(w, 0, 4) // default sample size
		writeInt(w, 0, 4) // default sample flags
	})
}

func writeMOOF(w io.Writer) {
}

func writeMDAT(w io.Writer) {
}
