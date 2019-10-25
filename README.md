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
- [Monitoring ](#monitoring)
- [Replay ](#replay)
- [Limitation](#limitation)
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


To mirror data from google storage that matches /data/ prefix and '.csv.gz' suffix to s3://destBucket/data
the following rule can be used 

[@gs://sourceBucket/config/config.json](usage/gs_to_s3/rule.json)
```json
[
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
    "Codec": "gzip",
    "Info": {
      "Workflow": "My workflow name here",
      "Description": "my description",
      "ProjectURL": "JIRA/WIKi or any link referece",
      "LeadEngineer": "my@dot.com"
    }
  }
]
```


### S3 to Google Storage

[![Google storage to S3 mirror](images/s3to_gs_mirror.png)](images/s3to_gs_mirror.png)


To mirror data from S3 that matches /myprefix/ prefix and '.csv.gz' suffix to gs://destBucket/data
splitting source file into maximum 8 MB files in destination you can use the following


[@gs://sourceBucket/config/config.json](usage/s3_to_gs/rule.json)
```json
[
  {
    "Source": {
      "Prefix": "/myprefix/",
      "Suffix": ".csv.gz"
    },
    "Dest": {
      "URL": "gs://${destBucket}",
      "Credentials": {
        "Parameter": "StorageMirror.GCP-DestProject",
        "Key": "smirror"
      }
    },
    "Split": {
      "MaxSize": 8388608,
      "Template": "%s_%05d"
    },
    "Codec": "gzip"
  }
]
```



 ### Google Storage To Pubsub

[![Google storage to Pubsub](images/g3ToPubsub.png)](images/g3ToPubsub.png)


To mirror data from google storage that match /myprefix/ prefix and '.csv' suffix to pubsub 'myTopic' topic
you can use the following rule

[@gs://sourceBucket/config/config.json](usage/gs_to_pubsub/rule.json)
```json
[
  {
    "Source": {
      "Prefix": "/myprefix/",
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
        "URL": "gs:///${opsBucket}/StorageMirror/Errors/"
      }
    ],
    "PreserveDepth": 1
  }
]
```


## Post actions

Each mirror rule accepts collection on OnSuccess and OnFailure post actions.

The follwoing action are supported:

- **delete**: remove source (trigger asset)


```json
{
   "OnSuccess": [{"Action": "delete"}]
}
```

- **move**: move source to specified destination

```json
{
  "OnSuccess": [{
          "Action": "move",
          "URL": "gs:///${opsBucket}/StorageMirror/Processed/"
 }],
 "OnFailure": [{
        "Action": "move",
        "URL": "gs:///${opsBucket}/StorageMirror/Errors/",
 }]
}
```

- **notify**: notify slack

```json
{
  "OnSuccess": [{
        "Action": "notify",
        "Title": "Transfer $SourceURL done",
        "Message": "success !!",
        "Channels": ["#e2e"],
        "Body": "$Response"
  }]
}
```



## Deployment

The following are used by storage mirror services:

**Prerequisites**

- _$configBucket_: bucket storing storage mirror configuration and mirror rules
- _$triggerBucket_: bucket storing data that needs to be mirror, event triggered by GCP
- _$opsBucket_: bucket string error, processed mirrors
-  Mirrors.BaseURL: location storing routes rules as JSON Array

The following [Deployment](deployment/mirror/README.md) details storage mirror generic deployment.


## Monitoring 

[StorageMonitor](mon) can be used to monitor trigger and error buckets.


**On Google Cloud Platform:**

```bash
curl -d @monitor.json -X POST  -H "Content-Type: application/json"  $monitorEndpoint
```

[@monitor.json](usage/monitor.json)
```json
{
  "ConfigURL":    "gs://${configBucket}/StorageMirror/config.json",
  "TriggerURL":   "gs://${triggerBucket}",
  "UnprocessedDuration": "2hours",
  "ErrorURL":     "gs://${opsBucket}/StorageMirror/Errors/",
  "ErrorRecency": "3hours"
}
```

_where:_
- **UnprocessedDuration** - check for any unprocessed data file over specified time
- **ErrorRecency** - specified errors within specified time


On Amazon Web Service cloud


```endly monitor.yaml authWith=aws-e2e```
[@monitor.yaml](usage/monitor.yaml)

```yaml
init:
  '!awsCredentials': $params.authWith
  bucketPrefix: ms-dataflow
  configBucket: ${bucketPrefix}-config
  triggerBucket: ${bucketPrefix}-trigger
  opsBucket: ${bucketPrefix}-operation

  monitor:
    ConfigURL: s3://${configBucket}/StorageMirror/config.json
    TriggerURL: s3://${triggerBucket}
    ErrorURL:  gs://${opsBucket}/StorageMirror/Errors/


pipeline:
  trigger:
    action: aws/lambda:call
    credentials: $awsCredentials
    functionname: StorageMonitor
    payload: $AsJSON($monitor)
```


## Replay

Sometimes during regular operation cloud function or lambda may terminate with error, leaving unprocess file. 
Replay function will move data back and forth between trigger and replay bucket triggering another event.
Each replayed file leaves trace in replay bucket  to control no more then one replay per file.


```bash
curl -d @replay.json -X POST  -H "Content-Type: application/json"  $replayEndpoint
```

[@replay.json](usage/replay.json)
```json
{
  "TriggerURL": "gs://${triggerBucket}",
  "ReplayBucket":"${replayBucket}",
  "UnprocessedDuration": "1hour"
}
```

_where:_
- **UnprocessedDuration** - check for any unprocessed data file over specified time


## Limitation

This project uses serverless stack, so any transfer exceeded memory  

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
