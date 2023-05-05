package handler

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/nagstler/serverless-csv-processor/internal/config"
	"github.com/nagstler/serverless-csv-processor/internal/csv"
)

func HandleS3Event(s3Event events.S3Event) error {
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

		cfg, err := config.LoadConfig(configObj.Body)
		configObj.Body.Close()
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
		csv.ReadCSV(reader, cfg)
	}

	return nil
}
