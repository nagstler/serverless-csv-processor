package sqs

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func SendMessage(queueURL string, message map[string]string) error {
	sess := session.Must(session.NewSession())
	sqsClient := sqs.New(sess)

	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return err
	}

	sendParams := &sqs.SendMessageInput{
		MessageBody: aws.String(string(jsonData)),
		QueueUrl:    aws.String(queueURL),
	}

	_, err = sqsClient.SendMessage(sendParams)
	if err != nil {
		log.Printf("Error sending message to SQS: %v", err)
		return err
	}

	return nil
}
