# Serverless CSV Processor
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

ServerlessCSVProcessor is a serverless application that processes large CSV files stored in an Amazon S3 bucket, maps the headers according to a configuration file, and pushes the records to an Amazon SQS queue. This project is built using AWS Lambda and is written in Golang.
## Features
- Efficient processing of large CSV files using concurrent workers
- Customizable header mapping through a JSON configuration file
- Triggered automatically when a new CSV file is uploaded to the S3 bucket
## Prerequisites
- Go 1.11 or later
- AWS CLI
- AWS account with access to Lambda, S3, and SQS services 
- [aws-lambda-go](https://github.com/aws/aws-lambda-go)  package
## Getting Started 
1. Clone the repository to your local machine.

```bash

git clone https://github.com/username/ServerlessCSVProcessor.git
cd ServerlessCSVProcessor
``` 
2. Initialize a new Go module by choosing a module name, which is usually the import path for your project. The module name should be unique to avoid conflicts with other projects. A common convention is to use your repository URL, like `github.com/username/project-name`. Replace `<module-name>` with a suitable name for your module, such as `github.com/johndoe/ServerlessCSVProcessor`.

```go

go mod init <module-name>
``` 
3. Download and install the required dependencies, including the `aws-lambda-go` package.

```go

go get -u github.com/aws/aws-lambda-go

go get -u github.com/aws/aws-sdk-go/aws 
go get -u github.com/aws/aws-sdk-go/aws/session 
go get -u github.com/aws/aws-sdk-go/service/s3 
go get -u github.com/aws/aws-sdk-go/service/s3/s3manager 
go get -u github.com/aws/aws-sdk-go/service/sqs


``` 
4. Compile the `main.go` file to create the Lambda binary.

```bash

GOOS=linux GOARCH=amd64 go build -o main main.go
``` 
5. Follow the deployment instructions in the [Deployment](#deployment)  section to set up your serverless CSV processor on AWS.

## Usage 
1. Upload your CSV file to the configured S3 bucket. 
2. Add a `config.json` file to the same S3 bucket, containing the header mappings and SQS queue URL. For example:

```json

{
    "header_mapping": {
        "OriginalHeader1": "MappedHeader1",
        "OriginalHeader2": "MappedHeader2"
    },
    "sqs_queue_url": "https://sqs.region.amazonaws.com/your-account-id/your-queue-name"
}
``` 
3. Once the CSV file is uploaded, the Lambda function will be triggered automatically, processing the CSV and pushing the records to the specified SQS queue.
## Deployment

**AWS Management Console:** 
You can manually create and configure the necessary AWS resources through the AWS Management Console:
- Create an S3 bucket to store your CSV files and the configuration file.
- Create an SQS queue to receive the parsed CSV records. 
- Compile your `main.go` file to create the Lambda binary.
- Create a Lambda function, set the runtime to "Go", and upload the compiled binary.
- In the Lambda function configuration, add an S3 trigger with the appropriate event type (e.g., "ObjectCreated"), and specify the S3 bucket you created earlier.
- Configure the necessary IAM roles and permissions for Lambda, S3, and SQS. 

There are multiple ways to deploy the serverless CSV processor to AWS. Refer to the [AWS Lambda Deployment](https://docs.aws.amazon.com/lambda/latest/dg/gettingstarted-awscli.html)  documentation for a detailed guide on deployment methods.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
## License

The gem is available as open source under the terms of the [MIT License](https://opensource.org/licenses/MIT).
