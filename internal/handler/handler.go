package handler

import (
	encCSV "encoding/csv" // Rename the import
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/nagstler/serverless-csv-processor/internal/config"
	csvpkg "github.com/nagstler/serverless-csv-processor/internal/csv" // Use import alias
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

		if key == "config.json" {
			continue
		}

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
		reader := encCSV.NewReader(csvFile) // Use the renamed import
		csvpkg.ReadCSV(reader, cfg)         // Use the import alias

		// Move the processed file to the archive folder
		err = moveToArchive(s3Client, bucket, key, cfg.ArchiveFolder)
		if err != nil {
			log.Printf("Error archiving file: %v", err)
		}
	}

	return nil
}

func moveToArchive(s3Client *s3.S3, bucket, key, archiveFolder string) error {
	archiveKey := fmt.Sprintf("%s/%s", archiveFolder, key)
	_, err := s3Client.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(bucket),
		CopySource: aws.String(fmt.Sprintf("%s/%s", bucket, key)),
		Key:        aws.String(archiveKey),
	})
	if err != nil {
		return fmt.Errorf("unable to move processed file to archive folder: %v", err)
	}

	// Delete the original file after copying
	_, err = s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("unable to delete original file after archiving: %v", err)
	}

	log.Printf("File archived: %s", archiveKey)
	return nil
}
