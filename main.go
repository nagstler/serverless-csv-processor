package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Config struct {
	HeaderMapping map[string]string `json:"header_mapping"`
	SQSQueueURL   string            `json:"sqs_queue_url"`
}

func processRows(config *Config, rows chan map[string]string, wg *sync.WaitGroup) {
	defer wg.Done()

	sess := session.Must(session.NewSession())
	sqsClient := sqs.New(sess)

	for row := range rows {
		jsonData, err := json.Marshal(row)
		if err != nil {
			log.Printf("Error marshaling row: %v", err)
			continue
		}

		sendParams := &sqs.SendMessageInput{
			MessageBody: aws.String(string(jsonData)),
			QueueUrl:    aws.String(config.SQSQueueURL),
		}

		_, err = sqsClient.SendMessage(sendParams)
		if err != nil {
			log.Printf("Error sending message to SQS: %v", err)
		}
	}
}

func handleS3Event(s3Event events.S3Event) error {
	log.Println("Lambda function triggered by S3 event")
	sess := session.Must(session.NewSession())
	downloader := s3manager.NewDownloader(sess)
	s3Client := s3.New(sess)

	for _, record := range s3Event.Records {
		log.Printf("Processing record: %+v", record)
		bucket := record.S3.Bucket.Name
		key := record.S3.Object.Key

		// Download and parse the configuration file
		configObj, err := s3Client.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("config.json"),
		})
		if err != nil {
			log.Fatalf("Unable to download config file: %v", err)
		}

		configData := &bytes.Buffer{}
		io.Copy(configData, configObj.Body)
		configObj.Body.Close()

		var config Config
		err = json.Unmarshal(configData.Bytes(), &config)
		if err != nil {
			log.Fatalf("Unable to parse config file: %v", err)
		}

		// Download and process the CSV file
		csvFile, err := os.Create("/tmp/tempfile.csv")
		if err != nil {
			log.Fatalf("Unable to create file: %v", err)
		}
		defer csvFile.Close()

		_, err = downloader.Download(csvFile, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			log.Fatalf("Unable to download CSV file: %v", err)
		}

		log.Printf("Processing file: %s", key)

		csvFile.Seek(0, 0)
		reader := csv.NewReader(csvFile)
		headers, err := reader.Read()
		if err != nil {
			log.Fatalf("Unable to read CSV headers: %v", err)
		}

		// Process rows concurrently
		rows := make(chan map[string]string, 100)
		var wg sync.WaitGroup
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go processRows(&config, rows, &wg)
		}

		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Error reading CSV row: %v", err)
				continue
			}

			row := make(map[string]string)
			for i, value := range record {
				header := headers[i]
				if mappedHeader, ok := config.HeaderMapping[header]; ok {
					row[mappedHeader] = value
				} else {
					row[header] = value
				}
			}
			rows <- row
		}

		close(rows)
		wg.Wait()
	}

	return nil
}

func main() {
	lambda.Start(handleS3Event)
}
