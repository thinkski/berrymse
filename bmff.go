package main

import (
	"bytes"
	"io"
)

// References:
// ISO/IEC 14496 Part 12
// ISO/IEC 14496 Part 15

// ISO/IEC 14496 Part 14 is not used.

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
		writeString(w, "isom")                 // major brand
		writeInt(w, 0x200, 4)                  // minor version
		writeString(w, "isomiso2iso5avc1mp41") // compatible brands
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
		writeInt(w, 0, 4)          // duration (all 1s == unknown)
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
		writeInt(w, -1, 4)         // next track id
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
		writeInt(w, 7, 4)               // version and flags (track enabled)
		writeInt(w, 0, 4)               // creation time
		writeInt(w, 0, 4)               // modification time
		writeInt(w, 1, 4)               // track id
		writeInt(w, 0, 4)               // reserved
		writeInt(w, 0, 4)               // duration
		writeInt(w, 0, 4)               // reserved
		writeInt(w, 0, 4)               // reserved
		writeInt(w, 0, 2)               // layer
		writeInt(w, 0, 2)               // alternate group
		writeInt(w, 0, 2)               // volume (ignored for video tracks)
		writeInt(w, 0, 2)               // reserved
		writeInt(w, 0x00010000, 4)      // matrix
		writeInt(w, 0x0, 4)             // matrix
		writeInt(w, 0x0, 4)             // matrix
		writeInt(w, 0x0, 4)             // matrix
		writeInt(w, 0x00010000, 4)      // matrix
		writeInt(w, 0x0, 4)             // matrix
		writeInt(w, 0x0, 4)             // matrix
		writeInt(w, 0x0, 4)             // matrix
		writeInt(w, 0x40000000, 4)      // matrix
		writeInt(w, int(width)<<16, 4)  // width (fixed-point 16.16 format)
		writeInt(w, int(height)<<16, 4) // height (fixed-point 16.16 format)
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
		writeInt(w, 0, 4)      // version and flags
		writeInt(w, 0, 4)      // creation time
		writeInt(w, 0, 4)      // modification time
		writeInt(w, 10000, 4)  // timescale
		writeInt(w, 0, 4)      // duration
		writeInt(w, 0x55c4, 2) // language ('und' == undefined)
		writeInt(w, 0, 2)      // pre-defined
	})
}

func writeHDLR(w io.Writer) {
	writeTag(w, "hdlr", func(w io.Writer) {
		writeInt(w, 0, 4)                        // version and flags
		writeInt(w, 0, 4)                        // pre-defined
		writeString(w, "vide")                   // handler type
		writeInt(w, 0, 4)                        // reserved
		writeInt(w, 0, 4)                        // reserved
		writeInt(w, 0, 4)                        // reserved
		writeString(w, "MicroMSE Video Handler") // name
		writeInt(w, 0, 1)                        // null-terminator
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
		writeSTSZ(w)
		writeSTSC(w)
		writeSTTS(w)
		writeSTCO(w)
	})
}

// Sample Table Box
func writeSTSD(w io.Writer, width, height uint16) {
	writeTag(w, "stsd", func(w io.Writer) {
		writeInt(w, 0, 6) // reserved
		writeInt(w, 1, 2) // deta reference index
		writeTag(w, "avc1", func(w io.Writer) {
			writeInt(w, 0, 6)           // reserved
			writeInt(w, 1, 2)           // data reference index
			writeInt(w, 0, 2)           // pre-defined
			writeInt(w, 0, 2)           // reserved
			writeInt(w, 0, 4)           // pre-defined
			writeInt(w, 0, 4)           // pre-defined
			writeInt(w, 0, 4)           // pre-defined
			writeInt(w, int(width), 2)  // width
			writeInt(w, int(height), 2) // height
			writeInt(w, 0x00480000, 4)  // horizontal resolution: 72 dpi
			writeInt(w, 0x00480000, 4)  // vertical resolution: 72 dpi
			writeInt(w, 0, 4)           // data size: 0
			writeInt(w, 1, 2)           // frame count: 1
			w.Write(make([]byte, 32))   // compressor name
			writeInt(w, 0x18, 2)        // depth
			writeInt(w, 0xffff, 2)      // pre-defined

			// Raspberry Pi 3B+ SPS/PPS for H.264 Main 4.0
			sps := []byte{
				0x27, 0x64, 0x00, 0x28, 0xac, 0x2b, 0x40, 0x28,
				0x02, 0xdd, 0x00, 0xf1, 0x22, 0x6a,
			}
			pps := []byte{
				0x28, 0xee, 0x02, 0x5c, 0xb0, 0x00,
			}

			// MPEG-4 Part 15 extension
			// See ISO/IEC 14496-15:2004 5.3.4.1.2
			writeTag(w, "avcC", func(w io.Writer) {
				writeInt(w, 1, 1)    // configuration version
				writeInt(w, 0x64, 1) // H.264 profile (0x64 == high)
				writeInt(w, 0x00, 1) // H.264 profile compatibility
				writeInt(w, 0x28, 1) // H.264 level (0x28 == 4.0)
				writeInt(w, 0xff, 1) // nal unit length - 1 (upper 6 bits == 1)
				writeInt(w, 0xe1, 1) // number of sps (upper 3 bits == 1)
				writeInt(w, len(sps), 2)
				w.Write(sps)
				writeInt(w, 1, 1) // number of pps
				writeInt(w, len(pps), 2)
				w.Write(pps)
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

func writeSTSZ(w io.Writer) {
	writeTag(w, "stsz", func(w io.Writer) {
		writeInt(w, 0, 4) // version and flags
		writeInt(w, 0, 4) // sample size
		writeInt(w, 0, 4) // sample count
	})
}

func writeSTCO(w io.Writer) {
	writeTag(w, "stco", func(w io.Writer) {
		writeInt(w, 0, 4) // version and flags
		writeInt(w, 0, 4) // entry count
	})
}

// Movie Extends Box
func writeMVEX(w io.Writer) {
	writeTag(w, "mvex", func(w io.Writer) {
		writeMEHD(w)
		writeTREX(w)
	})
}

// Movie Extends Header Box
func writeMEHD(w io.Writer) {
	writeTag(w, "mehd", func(w io.Writer) {
		writeInt(w, 0, 4) // version and flags
		writeInt(w, 0, 4) // fragment duration
	})
}

// Track Extends Box
func writeTREX(w io.Writer) {
	writeTag(w, "trex", func(w io.Writer) {
		writeInt(w, 0, 4)          // version and flags
		writeInt(w, 1, 4)          // track id
		writeInt(w, 1, 4)          // default sample description index
		writeInt(w, 0, 4)          // default sample duration
		writeInt(w, 0, 4)          // default sample size
		writeInt(w, 0x00010000, 4) // default sample flags
	})
}

// Movie Fragment Box
func writeMOOF(w io.Writer, seq int, data []byte) {
	writeTag(w, "moof", func(w io.Writer) {
		writeMFHD(w, seq)
		writeTRAF(w, seq, data)
	})
}

// Movie Fragment Header Box
func writeMFHD(w io.Writer, seq int) {
	writeTag(w, "mfhd", func(w io.Writer) {
		writeInt(w, 0, 4)   // version and flags
		writeInt(w, seq, 4) // sequence number
	})
}

// Track Fragment Box
func writeTRAF(w io.Writer, seq int, data []byte) {
	writeTag(w, "traf", func(w io.Writer) {
		writeTFHD(w)
		writeTFDT(w, seq)
		writeTRUN(w, data)
	})
}

// Track Fragment Header Box
func writeTFHD(w io.Writer) {
	writeTag(w, "tfhd", func(w io.Writer) {
		writeInt(w, 0x020020, 4)   // version and flags
		writeInt(w, 1, 4)          // track ID
		writeInt(w, 0x01010000, 4) // default sample flags
	})
}

// Track Fragment Base Media Decode Time Box
func writeTFDT(w io.Writer, seq int) {
	writeTag(w, "tfdt", func(w io.Writer) {
		writeInt(w, 0x01000000, 4) // version and flags
		writeInt(w, 330*seq, 8)    // base media decode time
	})
}

// Track Run Box
func writeTRUN(w io.Writer, data []byte) {
	writeTag(w, "trun", func(w io.Writer) {
		writeInt(w, 0x00000305, 4) // version and flags
		writeInt(w, 1, 4)          // sample count
		writeInt(w, 0x70, 4)       // data offset
		if (len(data) > 0) && (data[0]&0x1f == 0x5) {
			writeInt(w, 0x02000000, 4) // first sample flags (i-frame)
		} else {
			writeInt(w, 0x01010000, 4) // first sample flags (not i-frame)
		}
		writeInt(w, 330, 4)         // sample duration
		writeInt(w, 4+len(data), 4) // sample size
	})
}

// Media Data Box
func writeMDAT(w io.Writer, data []byte) {
	writeTag(w, "mdat", func(w io.Writer) {
		writeInt(w, len(data), 4)
		w.Write(data)
	})
}
