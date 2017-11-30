package main

import (
	"fmt"
	"math/rand"
	"io/ioutil"
	"log"
	"time"
	"text/template"
	"bytes"
	"net/http"
	"sync"
	"flag"
	"github.com/magento-mcom/send-messages/configuration"
)

func main() {
	filename := flag.String("config", "config.yml", "Configuration file")
	flag.Parse()

	config, err := configuration.Load(*filename)

	if err != nil {
		panic(fmt.Errorf("Failed to load configuration: %v", err))
	}

	start := time.Now()
	log.Printf("Started at")

	templates := config.Templates
	clients := config.Clients

	wg := sync.WaitGroup{}

	for i := 0; i < config.Routines; i++ {
		wg.Add(1)
		go func(counter int) {
			defer wg.Done()
			sendStockMessageToEndpoint(config, templates, clients, counter)
		}(i)
	}

	wg.Wait()

	defer func() {
		duration := time.Since(start)
		fmt.Printf("%d messages sent and it took %v seconds \n", config.Routines * config.Messages, duration.Seconds())
	}()

}

func sendStockMessageToEndpoint(config configuration.Config, templates []string, clients []string, routine int) {
	for i := 0; i < config.Messages; i++ {
		file := templates[rand.Intn(len(templates))]
		client := clients[rand.Intn(len(clients))]

		fmt.Printf("%v - Sending %d %s to %s \n", routine, i, file, client)

		body, err := ioutil.ReadFile(file)

		if err != nil {
			log.Printf("Cannot read %s /n", file)
			log.Println(err)
		}

		stockUpdateValues := map[string]interface{}{
			"Source": fmt.Sprintf("SOURCE-%v", rand.Intn(config.Sources)),
			"Sku":    fmt.Sprintf("SKU-%v", rand.Intn(config.Skus)),
			"Now":    time.Now().Format(time.RFC850),
			"Qty":    rand.Intn(10000),
			"Diff":   rand.Int()%10 - 5,
		}

		tmpl, err := template.New("stock-update").Parse(string(body))

		if err != nil {
			log.Fatal(err)
		}

		buffer := bytes.NewBuffer(nil)
		err = tmpl.Execute(buffer, stockUpdateValues)

		req, _ := http.NewRequest(http.MethodPost, config.Endpoint+client, bytes.NewReader([]byte(buffer.String())))
		_, err = http.DefaultClient.Do(req)

		if err != nil {
			log.Fatal(err)
		}
	}
}
