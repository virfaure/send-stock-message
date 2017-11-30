package app

import (
	"github.com/magento-mcom/send-stock-message/configuration"
	"github.com/magento-mcom/send-stock-message/consumer"
	"fmt"
	"io/ioutil"
	"log"
	"bytes"
	"math/rand"
	"text/template"
)

func NewLoader(c configuration.Config) Loader {
	return Loader{c, nil}
}

type Loader struct {
	c        configuration.Config
	consumer consumer.Consumer
}

func (l *Loader) Consumer() consumer.Consumer {

	if l.consumer == nil {
		l.consumer = consumer.NewSQSConsumer(l.c)
	}

	return l.consumer
}

func (l *Loader) SendReindexRequestToSqs(config configuration.Config, reindexTemplate string, clients []string, routine int) {
	for i := 0; i < config.Messages; i++ {
		client := clients[rand.Intn(len(clients))]

		fmt.Printf("%v - Sending %d %s to %s \n", routine, i, reindexTemplate, client)

		body, err := ioutil.ReadFile(reindexTemplate)

		if err != nil {
			log.Printf("Cannot read %s /n", reindexTemplate)
		}

		reindexRequestValues := map[string]interface{}{
			"Source": fmt.Sprintf("SOURCE-%v", rand.Intn(config.Sources)),
			"Sku":       fmt.Sprintf("SKU-%v", rand.Intn(config.Skus)),
			"Client":    client,
		}

		tmpl, err := template.New("reindex-request").Parse(string(body))

		if err != nil {
			log.Fatal(err)
			return
		}

		buffer := bytes.NewBuffer(nil)
		err = tmpl.Execute(buffer, reindexRequestValues)

		if err != nil {
			log.Fatal(err)
		}

		err = l.Consumer().SendReindexMessage(buffer.String())

		if err != nil {
			log.Fatal(err)
		}
	}
}

func (l *Loader) SendExportRequestToSqs(config configuration.Config, exportTemplate string, clients []string, routine int) {
	for i := 0; i < config.Messages; i++ {
		client := clients[rand.Intn(len(clients))]

		fmt.Printf("%v - Sending %d %s to %s \n", routine, i, exportTemplate, client)

		body, err := ioutil.ReadFile(exportTemplate)

		if err != nil {
			log.Printf("Cannot read %s /n", exportTemplate)
		}

		exportRequestValues := map[string]interface{}{
			"Aggregate": fmt.Sprintf("STORE-%v", rand.Intn(config.Sources)),
			"Sku":       fmt.Sprintf("SKU-%v", rand.Intn(config.Skus)),
			"Client":    client,
		}

		tmpl, err := template.New("export-request").Parse(string(body))

		if err != nil {
			log.Fatal(err)
			return
		}

		buffer := bytes.NewBuffer(nil)
		err = tmpl.Execute(buffer, exportRequestValues)

		if err != nil {
			log.Fatal(err)
		}

		err = l.Consumer().SendExportMessage(buffer.String())

		if err != nil {
			log.Fatal(err)
		}
	}
}
