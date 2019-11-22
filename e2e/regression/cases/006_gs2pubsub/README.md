#### Scenario:

Mirror suffixed *.csv data from gs://${gsTriggerBucket}/data/p1 to s3://${s3DestBucket}/data

#### Input:

Configuration:

* Global Config: [@config,json](../../../config/gs.json)
* Rule

[@rule.json](rule.json)
```json
[
  {
    "Prefix": "/data/p6",
    "Suffix": ".csv",
    "Dest": {
      "Topic": "${gsPrefix}_storage_mirror"
    },
    "Split": {
      "MaxLines": 10
    },
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

* event Type: s3:ObjectCreated:*
* resource: ${s3TriggerBucket}
* environmentVariables:
  - LOGGING: 'true'
  - CONFIG: s3://${s3ConfigBucket}/StorageMirror/config.json
 
Data:
- s3://${s3TriggerBucket}/data/p6/events.csv


Output:
- gs://${gsDestBucket}/data/p8/events.csv