package service

import (
	"context"
	"encoding/json"
	"goflow/entity"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSConsumer struct {
	client   *sqs.Client
	queueURL string
}

func NewSQSConsumer(client *sqs.Client, queueURL string) *SQSConsumer {
	return &SQSConsumer{
		client:   client,
		queueURL: queueURL,
	}
}

func (c *SQSConsumer) Consume(ctx context.Context) (*entity.Event, error) {
	result, err := c.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            &c.queueURL,
		MaxNumberOfMessages: 1,
		WaitTimeSeconds:     20,
	})

	if err != nil {
		return nil, err
	}

	if len(result.Messages) == 0 {
		return nil, nil
	}

	msg := result.Messages[0]

	var event entity.Event
	if err := json.Unmarshal([]byte(*msg.Body), &event); err != nil {
		return nil, err
	}

	event.MessageID = *msg.ReceiptHandle

	return &event, nil
}

func (c *SQSConsumer) Acknowledge(ctx context.Context, messageID string) error {
	_, err := c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &c.queueURL,
		ReceiptHandle: &messageID,
	})
	return err
}

func (c *SQSConsumer) Close() error {
	return nil
}
