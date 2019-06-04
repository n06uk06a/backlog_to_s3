GO111MODULE=on

.PHONY: deps clean build package deploy

deps:
	go get -u ./...

clean: 
	rm -rf ./build/*
	
build:
	GOOS=linux GOARCH=amd64 go build -o build/ci ./

package:
	aws cloudformation package \
		--template-file template.yaml \
		--s3-bucket lambda-testing.eatas.co.jp \
		--output-template-file package.yaml

deploy:
	aws cloudformation deploy \
		--template-file package.yaml \
		--stack-name CIStartFunction \
		--capabilities CAPABILITY_IAM
