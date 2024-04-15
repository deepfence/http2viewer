package main

import (
	"deepfence/http2viewer/pkg/config"
	"deepfence/http2viewer/pkg/http2"
	"flag"
	"fmt"
	"log"
	"sync"
)

func main() {
	configPath := flag.String("config", "", "Specify configuration to use, mandatory")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Panic(err.Error())
	}

	//ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	//if err != nil {
	//	log.Panic(err.Error())
	//}

	wg := sync.WaitGroup{}
	for _, s := range cfg.InputSources {
		wg.Add(1)
		go func(entry config.ConfigSource) {
			defer wg.Done()
			var capturer http2.HTTP2Capturer
			if entry.JSONCondig != nil {
				capturer = http2.NewJSONCapturer(
					entry.FilePath,
					entry.JSONCondig.StreamIDKey,
					entry.JSONCondig.DataKey,
				)
			} else if entry.PCAPConfig != nil {
				capturer = http2.NewPCAPCapturer(entry.FilePath)
			}

			for message := range http2.StreamHTTP2Message(capturer) {
				fmt.Printf("%v\n", message.String())
			}
		}(s)
	}
	wg.Wait()
	//<-ctx.Done()
}
