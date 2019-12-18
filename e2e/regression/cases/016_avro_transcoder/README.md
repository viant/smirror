#### Scenario:

Transcode CSV to AVRO format

#### Input:

**Configuration**:

* Global Config: [@config,json](../../../config/gs.json)
* Rule: [@rule,json](rule.json)

```json
[
  [
    {
      "Source": {
        "Prefix": "/data/p16/",
        "Suffix": ".csv"
      },
      "Dest": {
        "URL": "gs://${gsDestBucket}/data"
      },
      "Transcoder": {
        "Source": {
          "Format": "CSV",
          "Fields": [
            "id",
            "event_type",
            "timestamp"
          ],
          "HasHeader": true
        },
        "Dest": {
          "Format": "AVRO",
          "SchemaURL": "gs://${gsConfigBucket}/StorageMirror/Rules/case_${parentIndex}/schema.avsc",
          "RecordPerBlock": 10
        }
      },
      "OnSuccess": [
        {
          "Action": "move",
          "URL": "gs://${gsOpsBucket}/StorageMirror/processed/"
        }
      ],
      "OnFailure": [
        {
          "Action": "move",
          "URL": "gs://${gsOpsBucket}/StorageMirror/errors/"
        }
      ],
      "PreserveDepth": 1,
      "Info": {
        "Workflow": "e2e test case 1"
      }
    }
  ]
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
- gs://${gsTriggerBucket}/data/p16/events.csv


#### Output

**Data**
- gs://${gsDestBucket}/data/p16/events.avro