#### Scenario:

Mirror data from gs://${gsTriggerBucket}/data/p1 and suffixed *.csv to s3://${s3DestBucket}/data

#### Input:

**Configuration**:

* Global Config: [@config,json](../../../config/gs.json)
* Rule: [@rule,json](rule.json)

```json
[
  {
    "Prefix": "/data/p1",
    "Suffix": ".csv",
    "Dest": {
      "URL": "s3://${s3DestBucket}/data",
      "Credentials": {
        "URL": "gs://${gsConfigBucket}/Secrets/s3-mirror.json.enc",
        "Key": "projects/${gcpProject}/locations/us-central1/keyRings/${gsPrefix}_ring/cryptoKeys/${gsPrefix}_key"
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
        "URL": "gs:///${gsOpsBucket}/StorageMirror/Errors/"
      }
    ],
    "PreserveDepth": 1
  }
]
```
 

* **Trigger**:

    * event Type: google.storage.object.finalize
    * resource: projects/_/buckets/${gsTriggerBucket}
    * entryPoint: StorageMirror
    * environmentVariables:
      - LOGGING: 'true'
      - CONFIG: gs://${gsConfigBucket}/StorageMirror/config.json
 


**Data**:
- gs://${gsTriggerBucket}/data/p1/events.csv


#### Output

**Data**
- s3://${s3DestBucket}/data/p1/events.csv