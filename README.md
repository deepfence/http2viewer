# http2viewer

HTTP2 message viewer that reconstructs request/response out of HTTP2 stream.

## How it works

http2viewer can handle multiple source formats. It currently supports JSON & PCAP (clear text) formats.
http2viewer reconstructs HTTP2 message by aggregating raw buffers and recreates the HTTP2 frames.
Each frame is then analyzed and reconstructed into the original request/response HTTP2 messages.
This tool only works for offline payloads but can be extended for online reading.
Note that the viewer only shows HEADERS & BODY for each messages, control frames are ignored.

## How to use

Simply set the config.yaml:
```
input_sources:
  - file_path: /path/to/payload.json
    json_config:
      stream_id_key: "stream_id"
      data_key: "buffer"
```

Then tun the tool:

```
./http2viewer --config=config.yaml
```
## Example

payload.json:
```
{"buffer":"UFJJICogSFRUUC8yLjAKClNNCgo=","stream_id":0}
```

config.yaml:
```
input_sources:
  - file_path: payload.json
    json_config:
      stream_id_key: "stream_id"
      data_key: "buffer"
```

Output:
```
====== MESSAGE START ======
====== HEADERS ======
user-agent: curl/7.88.1
====== HEADERS END ======
====== BODY ======
{}
====== BODY END ======
====== MESSAGE END ======
```

Please have a look at the payload provided in the example folder, it is viewing the HTTP2 stream initiated by:
```
curl -v --http2 https://go.dev/serverpush
```
## What's next

- [ ] Add master key decrypt support (especially needed for PCAP)
- [ ] Add Online support
