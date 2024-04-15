package http2

import (
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type PCAPCapturer struct {
	StreamIDKey     string
	RawHTTP2DataKey string
	Dec             PCAPWrapper

	eth     layers.Ethernet
	ip4     layers.IPv4
	tcp     layers.TCP
	payload gopacket.Payload
	parser  *gopacket.DecodingLayerParser
}

func NewPCAPCapturer(path string) *PCAPCapturer {
	handle, err := pcap.OpenOffline(path)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil
	}
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	res := &PCAPCapturer{
		Dec: PCAPWrapper{
			packetSource: packetSource,
		},
	}
	parser := gopacket.NewDecodingLayerParser(
		layers.LayerTypeEthernet,
		&res.eth, &res.ip4, &res.tcp, &res.payload)
	res.parser = parser
	return res
}

type PCAPWrapper struct {
	packetSource *gopacket.PacketSource
}

func (w PCAPWrapper) Decode(input any) error {
	p, err := w.packetSource.NextPacket()
	if err != nil {
		return err
	}
	input = p
	return nil
}

type PCAPCapture struct {
	streamID     string
	rawHTTP2Data string
}

func (jc PCAPCapture) StreamID() string {
	return jc.streamID
}

func (jc PCAPCapture) RawHTTP2Data() string {
	return jc.rawHTTP2Data
}

func (capturer PCAPCapturer) NextCapture() (HTTP2Capture, error) {
	packet, err := capturer.Dec.packetSource.NextPacket()
	if err != nil {
		return PCAPCapture{}, err
	}
	streamID := "unknown"
	//if v, has := line[capturer.StreamIDKey]; has {
	//	streamID, _ = v.(string)
	//}

	decodedUnderlying := [4]gopacket.LayerType{}
	decoded := decodedUnderlying[:0]
	err = capturer.parser.DecodeLayers(packet.Data(), &decoded)
	return PCAPCapture{
		streamID:     streamID,
		rawHTTP2Data: string(capturer.payload),
	}, nil
}
