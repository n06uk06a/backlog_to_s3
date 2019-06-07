GO111MODULE=on

.PHONY: build package deploy

build:
	GOOS=linux GOARCH=amd64 go build -o build/backlog_to_s3 ./

package:
	aws cloudformation package \
		--template-file template.yaml \
		--s3-bucket ${BUCKET_NAME} \
		--output-template-file package.yaml

deploy:
	aws cloudformation deploy \
		--template-file package.yaml \
		--stack-name ${STACK_NAME} \
		--capabilities CAPABILITY_IAM
