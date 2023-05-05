package csv

import (
	"encoding/csv"
	"io"
	"sync"

	"github.com/nagstler/serverless-csv-processor/internal/config"
	"github.com/nagstler/serverless-csv-processor/internal/sqs"
)

func ProcessRows(config *config.Config, rows chan map[string]string, wg *sync.WaitGroup) {
	defer wg.Done()

	for row := range rows {
		sqs.SendMessage(config.SQSQueueURL, row)
	}
}

func ReadCSV(reader *csv.Reader, config *config.Config) {
	rows := make(chan map[string]string, 100)
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go ProcessRows(config, rows, &wg)
	}

	headers, err := reader.Read()
	if err != nil {
		// handle error
		return
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			// handle error
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
