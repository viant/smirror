#### Scenario:

Mirror data from gs://${gsTriggerBucket}/data/p8 and suffixed *.csv to s3://${s3DestBucket}/data

#### Input:

Configuration:

* Global Config: [@config,json](../../../config/gs.json)

* Rule

[@rule.json](rule.json)
```json
[
  {
    "Source": {
      "Prefix": "/data/p8",
      "Suffix": ".csv.gz"
    },
    "Replace": [
      {
        "From": "10",
        "To": "33333333"
      }
    ],
    "Dest": {
      "URL": "s3://${s3DestBucket}/data",
      "Credentials": {
        "URL": "gs://${gsConfigBucket}/Secrets/s3-mirror.json.enc",
        "Key": "projects/${gcpProject}/locations/us-central1/keyRings/${gsPrefix}_ring/cryptoKeys/${gsPrefix}_key"
      }
    },
    "Codec": "gzip",
    "OnSuccess": [
      {
        "Action": "delete"
      }
    ],
    "OnFailure": [
      {
        "Action": "move",
        "URL": "gs:///${gsOpsBucket}/StorageMirror/errors/"
      }
    ],
    "PreserveDepth": 1
  }
]
```
 


* Trigger:

* event Type: google.storage.object.finalize
* resource: projects/_/buckets/${gsTriggerBucket}
* entryPoint: StorageMirror
* environmentVariables:
  - LOGGING: 'true'
  - CONFIG: gs://${gsConfigBucket}/StorageMirror/config.json
 


Data:
- gs://${gsTriggerBucket}/data/p8/events.csv


Output:
- s3://${s3DestBucket}/data/p8/events.csv