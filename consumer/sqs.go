package consumer

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"github.com/magento-mcom/send-messages/configuration"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

var sqsSession *sqs.SQS

func NewSQSConsumer(config configuration.Config) Consumer {
	return &sqsConsumer{config}
}

type sqsConsumer struct {
	config configuration.Config
}

type Message struct{}

func (consumer *sqsConsumer) getSqsSession(queue string) (*sqs.SQS) {
	if sqsSession == nil {
		sess, err := session.NewSession(&aws.Config{
			Region:      aws.String(consumer.config.Queues.Region),
			Endpoint:    aws.String(queue),
			Credentials: credentials.NewSharedCredentials("", consumer.config.Queues.Profile),
		})

		if err != nil {
			log.Fatal(err)
		}

		sqsSession = sqs.New(sess)
	}

	return sqsSession
}

func (consumer *sqsConsumer) SendReindexMessage(body string) error {
	svc := consumer.getSqsSession(consumer.config.Queues.Reindex)

	params := &sqs.SendMessageInput{
		MessageBody:  aws.String(string(body[:])),
		QueueUrl:     aws.String(consumer.config.Queues.Reindex),
		DelaySeconds: aws.Int64(3),
	}

	_, err := svc.SendMessage(params)

	if err != nil {
		return err
	}

	return nil
}

func (consumer *sqsConsumer) SendExportMessage(body string) error {
	svc := consumer.getSqsSession(consumer.config.Queues.Export)

	params := &sqs.SendMessageInput{
		MessageBody:  aws.String(string(body[:])),
		QueueUrl:     aws.String(consumer.config.Queues.Export),
		DelaySeconds: aws.Int64(3),
	}

	_, err := svc.SendMessage(params)

	if err != nil {
		return err
	}

	return nil
}
