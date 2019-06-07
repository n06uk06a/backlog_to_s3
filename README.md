# Backlog to S3

When this is called by Backlog webhook, this archive source files and locate S3 Bucket.

## make API

### build
```
$ make build
```

### package & deploy
```
$ BUCKET_NAME=YourBucketName make package
$ STACK_NAME=YourStackName make deploy
```

## Backlog Setting
* Go ProjectSetting->Git Setting->Webhook URL
* Set URL to endpoint of API Gateway.
