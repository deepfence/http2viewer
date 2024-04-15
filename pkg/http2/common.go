package http2

import (
	"bytes"
	"io"
	"log"
	"strings"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"

	lru "github.com/hashicorp/golang-lru/v2/expirable"
)

type HTTP2Message struct {
	Headers   map[string][]string `json:"headers"`
	Body      string              `json:"body"`
	Direction string              `json:"direction"`
}

func (hm *HTTP2Message) String() string {
	res := "\n======" + " MESSAGE START " + "======" + "\n"
	res += "======" + " HEADERS " + "======" + "\n"
	for k, v := range hm.Headers {
		res += k + ": " + strings.Join(v, ", ")
		res += "\n"
	}
	res += "======" + " HEADERS END " + "======" + "\n"
	res += "======" + " BODY " + "======" + "\n"
	res += hm.Body
	res += "\n"
	res += "======" + " BODY END " + "======" + "\n"
	res += "======" + " MESSAGE END " + "======" + "\n"
	return res
}

type HTTP2Capture interface {
	StreamID() string
	RawHTTP2Data() string
}

type HTTP2Capturer interface {
	NextCapture() (HTTP2Capture, error)
}

func StreamHTTP2Message(capturer HTTP2Capturer) <-chan HTTP2Message {
	res := make(chan HTTP2Message)
	captures := make(chan HTTP2Capture)

	go func() {
		defer close(captures)
		defer log.Println("Close captured")
		for mess, err := capturer.NextCapture(); err == nil; mess, err = capturer.NextCapture() {
			if err == io.EOF {
				return
			}
			captures <- mess
		}
	}()

	go func() {
		defer close(res)
		defer log.Println("Close framer")

		streams := lru.NewLRU[string, []byte](8*1024, func(string, []byte) {}, 2*time.Minute)

		for capture := range captures {
			b, _ := streams.Get(capture.StreamID())
			buffer := append(b, []byte(capture.RawHTTP2Data())...)

			wbuf := bytes.NewBuffer([]byte{})
			offset := 0
			if bytes.HasPrefix(buffer, []byte(http2.ClientPreface)) {
				offset = len(http2.ClientPreface)
			}
			framer := http2.NewFramer(wbuf, bytes.NewReader(buffer[offset:]))

			completed := false
			currentHeaders := map[string][]string{}
			for frame, err := framer.ReadFrame(); err == nil; frame, err = framer.ReadFrame() {
				offset += int(frame.Header().Length) + 9
				switch f := frame.(type) {
				case *http2.HeadersFrame:
					hbuffer := hpack.NewDecoder(2048, nil)
					decoded, _ := hbuffer.DecodeFull(f.HeaderBlockFragment())
					for _, v := range decoded {
						currentHeaders[v.Name] = append(currentHeaders[v.Name], v.Value)
					}
				case *http2.DataFrame:
					res <- HTTP2Message{
						Headers:   currentHeaders,
						Body:      string(f.Data()),
						Direction: "",
					}
					completed = f.Flags.Has(http2.FlagDataEndStream)
					if len(buffer) == offset { // Avoid keeping large buffer in memory
						buffer = buffer[:0]
					} else {
						buffer = buffer[offset:]
					}
					offset = 0
				}
			}

			if completed {
				streams.Remove(capture.StreamID())
			} else {
				streams.Add(capture.StreamID(), buffer)
			}

		}
	}()

	return res
}
