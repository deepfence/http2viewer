all: http2viewer

http2viewer:
	go build -o http2viewer cmd/http2viewer/main.go

clean:
	rm http2viewer

.PHONY: http2viewer all
