package main

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/ssm"
)

var s3api *s3.S3
var ssmapi *ssm.SSM

var tempdir = os.TempDir()

func init() {
	sess := session.Must(session.NewSession())
	s3api = s3.New(sess)
	ssmapi = ssm.New(sess)
}

// HandleRequest handle request from API Gateway
func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// get payload from event
	eventJSON, err := json.Marshal(req)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(eventJSON))
	event := NewBacklogGitWebhookEvent(req.Body)

	repositoryName := event.Repository.Name
	branchName := strings.TrimPrefix(event.Ref, "refs/heads/")
	commitHash := event.After

	params := NewParameter(ssmapi, repositoryName, branchName)
	if params.Bucket == nil {
		log.Println("Nothing to do")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotModified,
		}, nil
	}

	mapping := map[string]string{
		"commitHash": commitHash,
		"imageName":  *params.ImageName,
		"branchName": branchName,
	}
	fp := createZip(*params.FileName, mapping)

	// Upload to S3
	resp, err := s3api.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(*params.Bucket),
		Key:    aws.String(fmt.Sprintf("%s%s", *params.Prefix, *params.ZipFileName)),
		Body:   fp,
	})
	if err != nil {
		log.Println(err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}
	respJSON, _ := json.Marshal(resp)
	log.Println(string(respJSON))

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusAccepted,
	}, nil
}

func createZip(fileName string, mapping map[string]string) io.ReadSeeker {
	t, err := template.New("repositoryTemplate").Parse(`export COMMIT_HASH={{.commitHash}}
export IMAGE_NAME={{.imageName}}
export BRANCH_NAME={{.branchName}}
`)
	if err != nil {
		log.Fatal(err)
	}
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)
	zf, err := zw.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(zf, mapping)
	zw.Close()

	return bytes.NewReader(buf.Bytes())
}

func main() {
	lambda.Start(HandleRequest)
}
