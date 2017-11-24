package consumer

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"log"
	"encoding/json"
	"github.com/magento-mcom/send-messages/app"
)

var sqsSession *sqs.SQS

type Message struct{}

type Consumer struct {
	sqs    *sqs.SQS
	config app.Config
}

func (consumer *Consumer) getSqsSession() () {
	if sqsSession == nil {
		sess, err := session.NewSession(&aws.Config{
			Region:      aws.String(consumer.config.Consumer.Region),
			Credentials: credentials.NewSharedCredentials("", consumer.config.Consumer.Profile),
		})

		if err != nil {
			log.Fatal(err)
		}

		sqsSession = consumer.sqs.New(sess)
	}

	return sqsSession
}

func (consumer *sqsConsumer) SendMessages() []common.ReindexRequest {
	svc := consumer.getSqsSession()

	result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:              &consumer.config.Consumer.Queue,
		MessageAttributeNames: aws.StringSlice([]string{"client"}),
		MaxNumberOfMessages:   aws.Int64(10),
		VisibilityTimeout:     aws.Int64(36000), // 10 hours
		WaitTimeSeconds:       aws.Int64(10),
	})

	if err != nil {
		log.Fatal(err)
	}

	if len(result.Messages) == 0 {
		log.Println("No message in the queue")
		return nil
	}

	var buffer []common.ReindexRequest

	for _, message := range result.Messages {
		var reindexRequest common.ReindexRequest
		if err := json.Unmarshal([]byte(*message.Body), &reindexRequest); err != nil {
			log.Printf("error unmarshalling: %v", err)
			continue
		}

		buffer = append(buffer, reindexRequest)

		deleteMessage(svc, consumer.config.Consumer.Queue, message)
	}

	return buffer
}
