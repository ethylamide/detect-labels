package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
)

type Request struct {
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
}

type Response struct {
	Labels []*rekognition.Label `json:"Labels"`
}

func Handler(request Request) (Response, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	svc := rekognition.New(sess)

	fmt.Printf("Bucket: %s\n", request.Bucket)
	fmt.Printf("Name: %s\n", request.Name)

	input := &rekognition.DetectLabelsInput{
		Image: &rekognition.Image{
			S3Object: &rekognition.S3Object{
				Bucket: aws.String(request.Bucket),
				Name:   aws.String(request.Name),
			},
		},
		MaxLabels:     aws.Int64(123),
		MinConfidence: aws.Float64(70.000000),
	}

	result, err := svc.DetectLabels(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rekognition.ErrCodeInvalidS3ObjectException:
				fmt.Println(rekognition.ErrCodeInvalidS3ObjectException, aerr.Error())
			case rekognition.ErrCodeInvalidParameterException:
				fmt.Println(rekognition.ErrCodeInvalidParameterException, aerr.Error())
			case rekognition.ErrCodeImageTooLargeException:
				fmt.Println(rekognition.ErrCodeImageTooLargeException, aerr.Error())
			case rekognition.ErrCodeAccessDeniedException:
				fmt.Println(rekognition.ErrCodeAccessDeniedException, aerr.Error())
			case rekognition.ErrCodeInternalServerError:
				fmt.Println(rekognition.ErrCodeInternalServerError, aerr.Error())
			case rekognition.ErrCodeThrottlingException:
				fmt.Println(rekognition.ErrCodeThrottlingException, aerr.Error())
			case rekognition.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(rekognition.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case rekognition.ErrCodeInvalidImageFormatException:
				fmt.Println(rekognition.ErrCodeInvalidImageFormatException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return Response{}, err
	}

	return Response{Labels: result.Labels}, nil
}

func main() {
	// bucket := os.Args[1]
	// name := os.Args[2]

	// resp, err := Handler(Request{Bucket: bucket, Name: name})

	// if err == nil {
	// 	fmt.Println(resp)
	// }

	lambda.Start(Handler)
}
