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
)

const (
	messageNumber    = 1000
	skuNumber        = 1000
	sourceNumber     = 100
	awsApiGatewayUrl = "https://28097ldii7.execute-api.us-west-2.amazonaws.com/beta/jsonrpc/"
)

func main() {
	start := time.Now()
	log.Printf("Started at")

	templates := []string{"templates/stock-adjustment.json", "templates/stock-update.json"}
	clients := []string{"HH", "DYSON", "DEVLYN", "LUMA", "TOUS"}

	wg := sync.WaitGroup{}

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sendStockMessageToSqs(templates, clients, i)
		}()
	}

	wg.Wait()

	duration := time.Since(start)
	fmt.Printf("Took %v seconds \n", duration.Seconds())
}

func sendStockMessageToSqs(templates []string, clients []string, routine int) {
	for i := 0; i < messageNumber; i++ {
		file := templates[rand.Intn(len(templates))]
		client := clients[rand.Intn(len(clients))]

		fmt.Printf("%v - Sending %d %s to %s \n", routine, i, file, client)

		body, err := ioutil.ReadFile(file)

		if err != nil {
			log.Printf("Cannot read %s /n", file)
		}

		stockUpdateValues := map[string]interface{}{
			"Source": fmt.Sprintf("SOURCE-%v", rand.Intn(sourceNumber)),
			"Sku":    fmt.Sprintf("SKU-%v", rand.Intn(skuNumber)),
			"Now":    time.Now(),
			"Qty":    rand.Intn(10000),
			"Diff":   rand.Int()%10 - 5,
		}

		tmpl, err := template.New("stock-update").Parse(string(body))

		if err != nil {
			panic(err)
		}

		buffer := bytes.NewBuffer(nil)
		err = tmpl.Execute(buffer, stockUpdateValues)

		req, _ := http.NewRequest(http.MethodPost, awsApiGatewayUrl + client, bytes.NewReader([]byte(buffer.String())))
		_, err = http.DefaultClient.Do(req)

		if err != nil {
			panic(err)
		}
	}
}
