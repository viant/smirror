# smirror - Serverless Cloud Storage Mirror

[![GoReportCard](https://goreportcard.com/badge/github.com/viant/smirror)](https://goreportcard.com/report/github.com/viant/smirror)
[![GoDoc](https://godoc.org/github.com/viant/smirror?status.svg)](https://godoc.org/github.com/viant/smirror)

This library is compatible with Go 1.11+

Please refer to [`CHANGELOG.md`](CHANGELOG.md) if you encounter breaking changes.


- [Motivation](#motivation)
- [Introduction](#introduction)
- [Usage](#usage)
   * [Google Storage to S3](#google-storage-to-s3)
   * [Google Storage to Pubsub](#google-storage-to-pubsub)
   * [S3 to Google Storage](#s3-to-google-storage)
- [Deployment](#deployment)
- [End to end testing](#end-to-end-testing)
- [Monitoring and limitation](#monitoring-and-limitation)
- [Code Coverage](#code-coverage)
- [License](#license)
- [Credits and Acknowledgements](#credits-and-acknowledgements)


## Motivation

When dealing with various cloud providers, it is a frequent use case to move seamlessly data from one cloud storage to another. 
In some scenarios, you may also need to split transferred content into a few smaller chunks. 
In any cases facilitating compression and post processing for both successful and failed transfer would be just additional requirement.


## Introduction

This project provide serverless implementation for cloud storage mirror. All external secrets/credentials are secured with KMS. 

**Google Storage to S3 Mirror**

[![Google storage to S3 mirror](images/g3Tos3Mirror.png)](images/g3Tos3Mirror.png)


## Usage

### Google Storage to S3

To mirror data from google storage that match /data/ prefix and '.csv.gz' suffix to s3://destBucket/data
preserving parent folder (folderDepth:1) the following configuration can be used with Mirror cloud function

[@gs://sourceBucket/config/config.json](usage/gs_to_s3/config.json)
```json
{
  "Mirrors": {
    "BaseURL": "gs://${gsConfigBucket}/StorageMirror/dataflow/",
    "Rules": [
      {
        "Source": {
          "Prefix": "/data/",
          "Suffix": ".csv.gz"
        },
        "Dest": {
          "URL": "s3://destBucket/data",
          "Credentials": {
            "URL": "gs://sourceBucket/secret/s3-cred.json.enc",
            "Key": "projects/my_project/locations/us-central1/keyRings/my_ring/cryptoKeys/my_key"
          }
        },
        "OnSuccess": [
          {
            "Action": "delete"
          }
        ],
        "OnFailure": [
          {
            "Action": "move",
            "URL": "gs://sourceBucket/data/errors/"
          }
        ],
        "Codec": "gzip",
        "PreserveDepth": 1
      },
      {
        "Source": {
          "Filter": "^\/[a-z]+/data/\\d+/",
          "Suffix": ".csv.gz"
        },
        "Dest": {
          "URL": "s3://destBucket/data/chunks/",
          "Credentials": {
            "URL": "gs://sourceBucket/secret/s3-cred.json.enc",
            "Key": "projects/my_project/locations/us-central1/keyRings/my_ring/cryptoKeys/my_key"
          }
        },
        "Split": {
          "MaxLines": 10000,
          "Template": "%s_%05d"
        },
        "OnSuccess": [
          {
            "Action": "delete"
          }
        ],
        "OnFailure": [
          {
            "Action": "move",
            "URL": "gs://sourceBucket/data/errors/"
          }
        ],
        "Codec": "gzip",
        "PreserveDepth": 1
      }
    ]
  }
}
```


### S3 to Google Storage

[![Google storage to S3 mirror](images/s3to_gs_mirror.png)](images/s3to_gs_mirror.png)


To mirror data from S3 that match /data/ prefix and '.csv.gz' suffix to gs://destBucket/data
preserving parent folder (folderDepth:1) the following configuration can be used with Mirror cloud function

[@gs://sourceBucket/config/config.json](usage/s3_to_gs/config.json)
```json
{
  "Mirrors": {
    "BaseURL": "gs://${gsConfigBucket}/StorageMirror/dataflow/",
    "Rules": [
      {
        "Source": {
          "Prefix": "/data/",
          "Suffix": ".csv.gz"
        },
        "Dest": {
          "URL": "gs://destBucket/data",
          "Credentials": {
            "Parameter": "smirror.gs",
            "Key": "smirror"
          }
        },
        "OnSuccess": [
          {
            "Action": "delete"
          }
        ],
        "OnFailure": [
          {
            "Action": "move",
            "URL": "s3://sourceBucket/data/errors/"
          }
        ],
        "Codec": "gzip",
        "PreserveDepth": 1
      },
      {
        "Source": {
          "Prefix": "/large/data/",
          "Suffix": ".csv.gz"
        },
        "Dest": {
          "URL": "gs://destBucket/data/chunks/",
          "Credentials": {
            "Parameter": "smirror.gs",
            "Key": "smirror"
          }
        },
        "Split": {
          "MaxLines": 10000,
          "Template": "%s_%05d"
        },
        "OnSuccess": [
          {
            "Action": "delete"
          }
        ],
        "OnFailure": [
          {
            "Action": "move",
            "URL": "s3://sourceBucket/data/errors/"
          }
        ],
        "Codec": "gzip",
        "PreserveDepth": 1
      }
    ]
  }
}
```


 ### Google Storage To Pubsub

[![Google storage to Pubsub](images/g3ToPubsub.png)](images/g3ToPubsub.png)


To mirror data from google storage that match /data/ prefix and '.csv' suffix to pubsub 'myTopic' topic
the following configuration can be used with Mirror cloud function

[@gs://sourceBucket/config/config.json](usage/gs_to_pubsub/config.json)
```json
{
  "Mirrors": {
    "Rules": [
      {
        "Source": {
          "Prefix": "/data/p6",
          "Suffix": ".csv"
        },
        "Dest": {
          "Topic": "myTopic"
        },
        "Split": {
          "MaxLines": 1000
        },
        "OnSuccess": [
          {
            "Action": "delete"
          }
        ],
        "OnFailure": [
          {
            "Action": "move",
            "URL": "gs:///${gsTriggerBucket}/e2e-mirror/errors/"
          }
        ],
        "PreserveDepth": 1
      }
    ]
  }
}
```


##Deployment

The following are used by storage mirror services:

**Prerequisites**

- _$configBucket_: bucket storing storage mirror configuration and mirror rules
- _$triggerBucket_: bucket storing data that needs to be mirror, event triggered by GCP
- _$opsBucket_: bucker string error, processed mirrors 
-  config:Mirrors.BaseURL: location storing routes rules as JSON Array



The following [Deployment](deployment/README.md) details generic deployment.


#### Google Cloud Function 

- With **endly cli**

```bash
endly deploy
```

[@deploy.yaml](usage/deploy/gcp/deploy.yaml)
```yaml
init:
  appPath: $Pwd(../..)
  gsTriggerBucket: ${s3TriggerBucket}

pipeline:

  package:
    action: exec:run
    target: $target
    commands:
      - unset GOPATH
      - cd ${appPath}
      - export GO111MODULE=on
      - go mod vendor

  deploy:
    action: gcp/cloudfunctions:deploy
    credentials: gcp-e2e
    '@name': StorageMirror
    entryPoint: StorageMirror
    runtime: go111
    availableMemoryMb: 2048
    timeout: 540s
    eventTrigger:
      eventType: google.storage.object.finalize
      resource: projects/_/buckets/${gsTriggerBucket}
    environmentVariables:
      LOGGING: 'true'
      CONFIG: gs://${gsConfigBucket}/StorageMirror/config.json
    source:
      URL: ${appPath}/
    sleepTimeMs: 5000
```



- With **gcloud cli**

```bash
unset GOPATH
export GO111MODULE=on
go mod vendor

gcloud functions deploy MyGsBucketToS3Mirror --entry-point Fn \ 
    --trigger-resource ${gsTriggerBucket} 
    --trigger-event google.storage.object.finalize \
    --set-env-vars=LOGGING=true,CONFIG=gs://gsTriggerBucket/mirror/config/gs.json \
    --memory=512M \
    --timeout=540s \
    --runtime=go111 
```




###### Deploying lambda

- With **endly cli**

```bash
endly deploy
```

[@deploy.yaml](usage/deploy/aws/deploy.yaml)
```yaml
init:
  appPath: $Pwd(/tmp/smirror)
  gsTriggerBucket: mys3TriggerBucket
  s3ConfigBucket: mys3ConfigBucket
  codeZip: ${appPath}/aws/smirror.zip
  functionName: StorageMirror
  privilegePolicy: privilege-policy.json
  awsCredentials: aws-e2e

pipeline:
  checkOut:
    action: vc/git:checkout
    Origin:
      URL: https://github.com/viant/smirror.git

  package:
    action: exec:run
    target: $target
    checkError: true
    commands:
      - cd ${appPath}
      - unset GOPATH
      - export GO111MODULE=on
      - go mod vendor
      - export GOOS=linux
      - export GOARCH=amd64
      - cd aws
      - go build -o smirror
      - zip -j smirror.zip smirror

  deploy:
    action: aws/lambda:deploy
    credentials: $awsCredentials
    functionname: $functionName
    runtime:  go1.x
    handler: smirror
    timeout: 360
    environment:
      variables:
        LOGGING: 'true'
        CONFIG: s3://${s3ConfigBucket}/StorageMirror/mirror.json
    code:
      zipfile: $LoadBinary(${codeZip})
    rolename: lambda-${functionName}-executor
    define:
      - policyname: kms-${functionName}-role
        policydocument: $Cat('${privilegePolicy}')
    attach:
      - policyarn: arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole

  setTrigger:
    action: aws/s3:setupBucketNotification
    sleepTimeMs: 20000
    bucket: ${s3TriggerBucket}
    lambdaFunctionConfigurations:
      - functionName: $functionName
        id: ObjectCreatedEvents
        events:
          - s3:ObjectCreated:*
        filter:
          prefix:
            - data
```

Where lambda uses permissions defined in [@privilege-policy.json](usage/deploy/aws/privilege-policy.json)

- With **aws cli** 
-[Serverless-deploying](https://docs.aws.amazon.com/lambda/latest/dg/with-userapp.html)

- With **sam cli**
-[Serverless-deploying](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-deploying.html)






###### Encrypting AWS credentials with GCP KMS 

In our example s3-cred.json.enc is encrypted version of [@s3-cred.json](usage/gs_to_s3/s3-cred.json) storing AWS credentials.

The following step can be used to encrypt a file.

- With **endly cli**


```bash
endly encryt
```

[@encrypt.yaml](usage/deploy/gcp/encrypt.yaml)
```yaml
init:
  gsConfigBucket: ${gsConfigBucket}
  
pipeline:
  secure:
    deployKey:
      action: gcp/kms:deployKey
      credentials: gcp-e2e
      ring: my_ring
      key: my_key
      logging: false
      purpose: ENCRYPT_DECRYPT
      bindings:
        - role: roles/cloudkms.cryptoKeyEncrypterDecrypter
          members:
            - serviceAccount:$gcp.serviceAccount

    keyInfo:
      action: print
      message: 'Deployed Key: ${deployKey.Name}'

    encrypt:
      action: gcp/kms:encrypt
      logging: false
      ring: my_ring
      key: my_key
      source:
        URL: s3-cred.json
      dest:
        credentials: gcp-e2e
        URL: gs://$gsConfigBucket/StorageMirror/secrets/s3-cred.json.enc
    info:
      action: print
      message: ${encrypt.CipherBase64Text}
```

Where [gcp-credentials](https://github.com/viant/endly/tree/master/doc/secrets#gc) is service account based GCP secrets
stored in ~/.secret/gcp-credentials.json


- With **gcloud cli**

```bash

gcloud kms keyrings create my_ring --location us-central1
gcloud kms keys create my_key --location us-central1 \
  --keyring my_ring --purpose encryption

## Encrypt s3-cred.json

gcloud kms encrypt \
  --location=us-central1  \
  --keyring=my_ring \
  --key=my_key \
  --version=1 \
  --plaintext-file=s3-cred.json \
  --ciphertext-file=s3-cred.json.enc

## Upload encrypted version to google storage

gsutil cp s3-cred.json.enc gs://sourceBucket/secret/s3-cred.json.enc

```

###### Encrypting AWS credentials with GCP KMS 

In our example AWS System Manager  'smirror.gcp' parameters is encrypted version of [@gcp-cred.json](usage/s3_to_gs/gcp-cred.json) Google Secrets.

The following step can be used to encrypt a google secrets.

- With **endly cli**

[@encrypt.yaml](usage/s3_to_gs/encrypt.yaml)
```yaml
init:
  awsCredentials: aws-e2e
  gcpSecrets: $Cat(gcp-cred.json)

pipeline:
  secure:
    credentials: $awsCredentials
    action: aws/kms:setupKey
    aliasName: alias/storagemirror

  encrypt:
    action: aws/ssm:setParameter
    name: smirror.gcp
    '@description': Google Storage credentials
    type: SecureString
    keyId: alias/storagemirror
    value: $gcpSecrets
```


- With **aws cli**


```bash
- aws kms create-key  
- aws kms create-alias --alias-name=smirror --target-key-id=KEY_ID
- aws ssm put-parameter \
    --name "smirror.gs" \
    --value 'CONTENT OF GOOGLE SECRET HERE' \
    --type SecureString \
    --key-id alias/storagemirror

```





## End to end testing


### Prerequisites:

  - [Endly e2e runner](https://github.com/viant/endly/releases) or [endly docker image](https://github.com/viant/endly/tree/master/docker)
  - [Google secrets](https://github.com/viant/endly/tree/master/doc/secrets#google-cloud-credentials) for dedicated e2e project  ~/.secret/gcp-e2e.json 
  - [AWS secrets](https://github.com/viant/endly/tree/master/doc/secrets#asw-credentials) for dedicated e2e account ~/.secret/aws-e2e.json 


```bash
git clone https://github.com/viant/smirror.git
cd smirror/e2e
### Update mirrors bucket for both S3, Google Storage in e2e/run.yaml (gsTriggerBucket, s3TriggerBucket)
endly 
```


## Monitoring and limitation


## Code Coverage

[![GoCover](https://gocover.io/github.com/viant/smirror)](https://gocover.io/github.com/viant/smirror)

	
<a name="License"></a>
## License

The source code is made available under the terms of the Apache License, Version 2, as stated in the file `LICENSE`.

Individual files may be made available under their own specific license,
all compatible with Apache License, Version 2. Please see individual files for details.


<a name="Credits-and-Acknowledgements"></a>

## Credits and Acknowledgements

**Library Author:** Adrian Witas

