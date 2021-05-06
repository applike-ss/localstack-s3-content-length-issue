package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"strings"
)

func main() {
	config := &aws.Config{
		CredentialsChainVerboseErrors: aws.Bool(true),
		MaxRetries:                    aws.Int(30),
		Region:                        aws.String(endpoints.EuCentral1RegionID),
		Endpoint:                      aws.String("http://127.0.0.1:4566"),
		S3ForcePathStyle:              aws.Bool(true),
	}

	newSession, err := session.NewSession(config)

	if err != nil {
		panic(err)
	}
	client := s3.New(newSession)

	bucketName := "my-bucket"

	cbi := s3.CreateBucketInput{
		ACL:    aws.String("public"),
		Bucket: aws.String(bucketName),
	}
	_, err = client.CreateBucket(&cbi)
	if err != nil {
		if awsErr, ok := err.(awserr.RequestFailure); ok {
			code := awsErr.Code()
			if code != "BucketAlreadyExists" {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
	}

	reader := strings.NewReader("foobarbaz")
	filename := "my-file.txt"

	poi := s3.PutObjectInput{
		ACL:    aws.String("public-read"),
		Body:   reader,
		Key:    aws.String(filename),
		Bucket: aws.String(bucketName),
	}

	_, err = client.PutObject(&poi)
	if err != nil {
		fmt.Println(err.Error())
	}

	goi := s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
	}

	goo, err := client.GetObject(&goi)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(fmt.Sprintf("content length is: %v", goo.ContentLength))
}
