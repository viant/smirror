#### Scenario:

Mirror data from gs://${gsTriggerBucket}/data/p8 and suffixed *.csv to s3://${s3DestBucket}/data

#### Input:

Configuration:

* Global Config: [@config,json](../../../config/gs.json)

* Route:

[@routes,json](rule.json)
```json
[
  {
    "Source": {
      "Prefix": "/data/p8",
      "Suffix": ".csv"
    },
    "Replace": {
        "10": "33333333"
    },
    "Dest": {
      "URL": "s3://${s3DestBucket}/data",
      "Credentials": {
        "URL": "gs://${gsConifgBucket}/Secrets/s3-mirror.json.enc",
        "Key": "projects/${gcpProject}/locations/us-central1/keyRings/gs_mirror_ring/cryptoKeys/gs_mirror_key"
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
  - CONFIG: gs://${gsConifgBucket}/StorageMirror/config.json
 


Data:
- gs://${gsTriggerBucket}/data/p8/events.csv


Output:
- s3://${s3DestBucket}/data/p8/events.csv