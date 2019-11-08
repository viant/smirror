#### Scenario:

Mirror data from gs://${gsTriggerBucket}/data/p12/archive.zip  to gs://${gsTriggerBucket}/data/p12/archive.zip/*.gz

Since **Uncompress** option is specified all files from archive are individually mirror to destination.

#### Input:

**Configuration**:

* Global Config: [@config,json](../../../config/gs.json)
* Rule: [@rule,json](rule.json)

```json
[
  {
    "Source": {
      "Prefix": "/data/p12/",
      "Suffix": ".zip"
    },
    "Dest": {
      "URL": "s3://${s3DestBucket}/data",
      "Credentials": {
        "URL": "gs://${gsConfigBucket}/Secrets/s3-mirror.json.enc",
        "Key": "projects/${gcpProject}/locations/us-central1/keyRings/${gsPrefix}_ring/cryptoKeys/${gsPrefix}_key"
      }
    },
    "Uncompress": true,
    "Codec": "gzip",
    "OnSuccess": [
      {
        "Action": "move",
        "URL": "gs://${gsOpsBucket}/StorageMirror/processed/"
      },
      {
        "Action": "notify",
        "Title": "Transfer $SourceURL done",
        "Message": "success !!",
        "Channels": [
          "#e2e"
        ],
        "Body": "$Response"
      }
    ],
    "PreserveDepth": -1,
    "Info": {
      "Workflow": "e2e test case 1"
    }
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
- gs://${gsTriggerBucket}/data/p12/archive.zip


#### Output

**Data**
- s3://${s3DestBucket}/data/p12/archive.zip/events.csv.gz