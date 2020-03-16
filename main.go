// MicroMSE: Low-latency live video via Media Source Extensions (MSE)
// Copyright (C) 2020 Chris Hiszpanski
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Command line flag parameters
var (
	flagListen string
)

func init() {
	flag.StringVar(
		&flagListen, "l", "localhost:8000", "listen on host:port",
	)
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
)

// Websocket parameters
var upgrader = websocket.Upgrader{
	// Tune read buffers for short acknowledgement messages
	ReadBufferSize: 256,

	// Tune write buffers to comfortably fit most all B and P frames.
	WriteBufferSize: 8192,

	// Allow any origin for demo purposes
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// client structure
type client struct {
	hub   *hub
	conn  *websocket.Conn // Websocket connection
	frags chan []byte     // Buffered channel of outbound MP4 fragments
}

// hub maintains a set of active clients and broadcasts video to clients
type hub struct {
	clients    map[*client]bool // registered clients
	nals       chan []byte      // NAL units from camera source
	register   chan *client     // register requests from clients
	unregister chan *client     // unregister requests from clients
}

// newHub instantiates a new hub
func newHub() *hub {
	return &hub{
		clients:    make(map[*client]bool),
		nals:       make(chan []byte),
		register:   make(chan *client),
		unregister: make(chan *client),
	}
}

// run processes register and unregister requests, and nal units
func (h *hub) run() {
	for {
		select {
		// Register request
		case c := <-h.register:
			h.clients[c] = true

		// Unregister request
		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.frags)
			}

		// New NAL from source
		case nal := <-h.nals:
			for c := range h.clients {
				// TODO Convert NAL unit into MP4 fragment

				select {
				// Encode NAL unit into MP4 fragment
				case c.frags <- frag:

				// Buffered channel full. Drop client.
				default:
					close(c.frags)
					delete(h.clients, c)
				}
			}
		}
	}
}

type source struct {
}

func newSource() *source {
}

// Handle websocket client connections
func serveWs(h *hub, w http.ResponseWriter, r *http.Request) {
	// Upgrade websocket connection from HTTP to TCP
	conn, err := upgrader.Upgrade(w, r, nil)
	if nil != err {
		log.Println(err)
		return
	}

	// Instantiate client
	c := &client{hub: h, conn: conn, frags: make(chan []byte, 30)}
	c.hub.register <- c

	// Go routine writes new MP4 fragment to client websocket
	go func(c *client) {
		defer func() {
			c.conn.Close()
		}()

		for {
			select {
			case frag, ok := <-c.frags:
				c.conn.SetWriteDeadline(time.Now().Add(writeWait))
				if !ok {
					// Hub closed the channel
					c.conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}

				// Write next segment
				nw, err := c.conn.NextWriter(websocket.BinaryMessage)
				if nil != err {
					return
				}
				nw.Write(frag)

				// Close writer
				if err := nw.Close(); nil != err {
					return
				}
			case <-ticker.C:
			}
		}
	}(c)
}

func main() {
	flag.Parse()

	// Parse host:port into host and port
	host, port, err := net.SplitHostPort(flagListen)
	if nil != err {
		log.Fatal(err)
	}

	// One-to-many hub broadcasts NAL units as MP4 fragments to clients
	hub := newHub()
	go hub.run()

	// TODO Open source

	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	http.Handle("/", http.FileServer(http.Dir("web/static/")))

	// Start server
	fmt.Printf("Listening on http://%v:%v\n", host, port)
	log.Fatal(http.ListenAndServe(flagListen, nil))
}
