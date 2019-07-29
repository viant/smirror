# smirror - Serverless Cloud Storage Mirror

[![GoReportCard](https://goreportcard.com/badge/github.com/viant/smirror)](https://goreportcard.com/report/github.com/viant/smirror)
[![GoDoc](https://godoc.org/github.com/viant/smirror?status.svg)](https://godoc.org/github.com/viant/smirror)

This library is compatible with Go 1.11+

Please refer to [`CHANGELOG.md`](CHANGELOG.md) if you encounter breaking changes.


- [Motivation](#motivation)
- [Introduction](#introduction)
- [Usage](#usage)
   * [GS to S3 Mirror](#gs-to-s3-mirror)
   * [S3 to G3 Mirror](#s3-to-gs-mirror)
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

### GS to S3 mirror

To mirror data from google storage that match /data/ prefix and '.csv.gz' suffix to s3://destBucket/data
preserving parent folder (folderDepth:1) the following configuration can be used with Mirror cloud function

[@gs://sourceBucket/config/config.json](usage/gs_to_s3/config.json)
```json
{
  "Routes": [
    {
      "Prefix": "/data/",
      "Suffix": ".csv.gz",
      "DestURL": "s3://destBucket/data",
      "OnCompletion": {
        "OnSuccess": [
          {
            "Action": "delete"
          }
        ],
        "OnError": [
          {
            "Action": "move",
            "URL": "gs://sourceBucket/data/errors/"
          }
        ]
      },
      "Codec": "gzip",
      "FolderDepth": 1
    },
    {
      "Prefix": "/large/data/",
      "Suffix": ".csv.gz",
      "DestURL": "s3://destBucket/data/chunks/",
      "Split": {
        "MaxLines": 10000,
        "Template": "%s_%05d"
      },
      "OnCompletion": {
        "OnSuccess": [
          {
            "Action": "delete"
          }
        ],
        "OnError": [
          {
            "Action": "move",
            "URL": "gs://sourceBucket/data/errors/"
          }
        ]
      },
      "Codec": "gzip",
      "FolderDepth": 1
    }
  ],
  "Secrets": [
    {
      "Provider": "gcp",
      "TargetScheme": "s3",
      "URL": "gs://sourceBucket/secret/s3-cred.json.enc",
      "Key": "projects/my_project/locations/us-central1/keyRings/my_ring/cryptoKeys/my_key"
    }
  ]
}
```

###### Encrypting AWS credentials with GCP KMS 

In our example s3-cred.json.enc is encrypted version of [@s3-cred.json](usage/gs_to_s3/s3-cred.json) storing AWS credentials.

The following step can be used to encrypt a file.

- With **endly cli**


```bash
endly encryt
```

[@encrypt.yaml](usage/gs_to_s3/encrypt.yaml)
```yaml
pipeline:
  secure:
    deployKey:
      action: gcp/kms:deployKey
      credentials: gcp-credentials
      ring: my_ring
      key: my_key
      logging: false
      purpose: ENCRYPT_DECRYPT
      bindings:
        - role: roles/cloudkms.cryptoKeyEncrypterDecrypter
          members:
            - serviceAccount:$gcp.serviceAccount

    encrypt:
      action: gcp/kms:encrypt
      logging: false
      ring: my_ring
      key: my_key
      source:
        URL: config.json
      dest:
        credentials: gcp-credentials
        URL: gs://sourceBucket/secret/s3-cred.json.enc
        
    info:
      action: print
      message: ${encrypt.CipherData}
```

Where [gcp-credentials](https://github.com/viant/endly/tree/master/doc/secrets#gc) is service account based GCP secrets
stored in ~/.secret/gcp-credentials.json

- With **gcloud cli**

```bash

## Crate symmetric key

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

###### Deploying cloud function


- With **endly cli**

```bash
endly deploy
```

[@deploy.yaml](usage/gs_to_s3/deploy.yaml)
```yaml
init:
  appPath: $Pwd(../..)
  gsBucket: e2etst

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
    '@name': MyGsBucketToS3Mirror
    entryPoint: Fn
    runtime: go111
    availableMemoryMb: 512
    eventTrigger:
      eventType: google.storage.object.finalize
      resource: projects/_/buckets/${gsBucket}
    environmentVariables:
      LOGGING: 'true'
      CONFIG: gs://gsBucket/mirror/config/gs.json
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
    --trigger-resource e2etst 
    --trigger-event google.storage.object.finalize \
    --set-env-vars=LOGGING=true,CONFIG=gs://gsBucket/mirror/config/gs.json \
    --memory=512M \
    --timeout=500s \
    --runtime=go111 
```

### S3 to GS mirror







## End to end testing

### Prerequisites:

  - [Endly e2e runner](https://github.com/viant/endly/releases) or [endly docker image](https://github.com/viant/endly/tree/master/docker)
  - [Google secrets](https://github.com/viant/endly/tree/master/doc/secrets#google-cloud-credentials) for dedicated e2e project  ~/.secret/gcp-e2e.json 
  - [AWS secrets](https://github.com/viant/endly/tree/master/doc/secrets#asw-credentials) for dedicated e2e account ~/.secret/aws-e2e.json 





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

