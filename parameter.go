package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/ssm"
)

// Parameter got from SSM
type Parameter struct {
	Bucket      *string `json:"bucket"`
	Prefix      *string `json:"prefix"`
	ImageName   *string `json:"imageName"`
	FileName    *string `json:"fileName"`
	ZipFileName *string `json:"zipFileName"`
}

// NewParameter get Parameter from SSM
func NewParameter(ssmapi *ssm.SSM, repositoryName, branchName string) *Parameter {
	param := new(Parameter)
	path := fmt.Sprintf("/git/backlog/%s/%s", repositoryName, branchName)
	b := true
	out, err := ssmapi.GetParameter(&ssm.GetParameterInput{
		Name:           &path,
		WithDecryption: &b,
	})
	if err == nil {
		outParam := out.Parameter
		value := outParam.Value
		if err := json.Unmarshal([]byte(*value), param); err != nil {
			log.Fatal(err)
		}
	}
	return param
}
