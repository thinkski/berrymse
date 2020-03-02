package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

var (
	flagListen string
)

func init() {
	flag.StringVar(&flagListen, "l", "localhost:8000", "listen on host:port")
}

type ByteStream interface {
	InitializationSegment() []byte
	MediaSegment() []byte
}

type MP4 struct {
	file *os.File
}

func OpenMP4(filename string) (*MP4, error) {
	f, err := os.Open(filename)
	if nil != err {
		return nil, err
	}

	return &MP4{f}, nil
}

type Box interface {
}

func (mp4 *MP4) ReadBox() ([]byte, error) {
	hdr := make([]byte, 8)

	n, err := mp4.file.Read(hdr)
	if nil != err || 8 != n {
		return nil, err
	}

	size := binary.BigEndian.Uint32(hdr[0:4])

	b := make([]byte, size-8)
	mp4.file.Read(b)

	return append(hdr, b...), nil
}

func (mp4 *MP4) ReadSegment() ([]byte, error) {
	segment := make([]byte, 0)

	for {
		box, err := mp4.ReadBox()
		if nil != err {
			return nil, err
		}

		typ := string(box[4:8])

		switch typ {
		case "ftyp":
			segment = append(segment, box...)
		case "moov":
			segment = append(segment, box...)
			return segment, nil
		case "moof":
			segment = append(segment, box...)
		case "mdat":
			segment = append(segment, box...)
			return segment, nil
		}
	}
}

func PrintBoxTypes(box []byte, level int) {
	i := 0

	for i+8 < len(box) {
		size := binary.BigEndian.Uint32(box[i : i+4])
		typ := string(box[i+4 : i+8])

		fmt.Println(strings.Repeat("  ", level), typ, size)
		switch typ {
		case "moov", "trak", "edts", "mdia", "minf", "mvex", "udta", "stbl", "dinf":
			PrintBoxTypes(box[i+8:i+int(size-8)], level+1)
		}

		i += int(size)
	}
}

func (mp4 *MP4) Close() {
	mp4.file.Close()
}

// initialization segment: single ftyp followed by single moov

// media segment: optional styp, single moof, one or more mdat

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if nil != err {
		log.Println(err)
		return
	}
	defer conn.Close()

	mp4, err := OpenMP4("testdata/bunny.mp4")
	if nil != err {
		panic(err)
	}
	defer mp4.Close()

	for {
		_, _, err := conn.ReadMessage()
		if nil != err {
			log.Println(err)
			break
		}

		segment, err := mp4.ReadSegment()
		PrintBoxTypes(segment, 0)
		if nil != err {
			log.Println(err)
			break
		}

		err = conn.WriteMessage(websocket.BinaryMessage, segment)
		if nil != err {
			log.Println(err)
			break
		}
	}
}

func main() {
	flag.Parse()

	host, port, err := net.SplitHostPort(flagListen)
	if nil != err {
		log.Fatal(err)
	}

	http.HandleFunc("/websocket", websocketHandler)
	http.Handle("/", http.FileServer(http.Dir("web/static/")))
	fmt.Printf("Listening on http://%v:%v\n", host, port)
	log.Fatal(http.ListenAndServe(flagListen, nil))
}
