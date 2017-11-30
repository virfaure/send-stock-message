package main

import (
	"fmt"
	"log"
	"time"
	"sync"
	"flag"
	"github.com/magento-mcom/send-messages/app"
	"github.com/magento-mcom/send-messages/configuration"
)

func main() {
	filename := flag.String("config", "config.yml", "Configuration file")
	flag.Parse()

	config, err := configuration.Load(*filename)

	if err != nil {
		panic(fmt.Errorf("Failed to load configuration: %v", err))
	}

	l := app.NewLoader(config)

	start := time.Now()
	log.Printf("Started at")

	reindexTemplate := "templates/reindex-request.json"
	clients := []string{"HH", "DYSON", "DEVLYN", "LUMA", "TOUS"}

	wg := sync.WaitGroup{}

	for i := 0; i < config.Routines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.SendReindexRequestToSqs(config, reindexTemplate, clients, i)
		}()
	}

	wg.Wait()

	duration := time.Since(start)
	fmt.Printf("Took %v seconds \n", duration.Seconds())
}
