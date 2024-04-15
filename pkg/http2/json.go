package http2

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
)

type JSONCapturer struct {
	StreamIDKey     string
	RawHTTP2DataKey string
	Dec             *json.Decoder
	BufDec          *base64.Encoding // Optional
}

func NewJSONCapturer(path, streamKey, dataKey string) JSONCapturer {
	b, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return JSONCapturer{}
	}
	return JSONCapturer{
		StreamIDKey:     streamKey,
		RawHTTP2DataKey: dataKey,
		Dec:             json.NewDecoder(bytes.NewReader(b)),
		BufDec:          base64.StdEncoding,
	}
}

type JSONCapture struct {
	streamID     string
	rawHTTP2Data string
}

func (jc JSONCapture) StreamID() string {
	return jc.streamID
}

func (jc JSONCapture) RawHTTP2Data() string {
	return jc.rawHTTP2Data
}

func (capturer JSONCapturer) NextCapture() (HTTP2Capture, error) {
	line := map[string]interface{}{}
	err := capturer.Dec.Decode(&line)
	if err != nil {
		return JSONCapture{}, err
	}
	streamID := "unknown"
	if v, has := line[capturer.StreamIDKey]; has {
		streamID, _ = v.(string)
	}
	rawHTTP2Data := "unknown"
	if v, has := line[capturer.RawHTTP2DataKey]; has {
		rawHTTP2Data, _ = v.(string)
	}
	if capturer.BufDec != nil {
		b, err := capturer.BufDec.DecodeString(rawHTTP2Data)
		if err != nil {
			return JSONCapture{}, err
		}
		rawHTTP2Data = string(b)
	}
	return JSONCapture{
		streamID:     streamID,
		rawHTTP2Data: rawHTTP2Data,
	}, nil
}
