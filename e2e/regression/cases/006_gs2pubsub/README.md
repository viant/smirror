#### Scenario:

Mirror suffixed *.csv data from gs://${gsTriggerBucket}/data/p1 to s3://${s3DestBucket}/data

#### Input:

Configuration:

* Global Config: [@config,json](../../../config/gs.json)
* Rule

[@routes,json](rule.json)
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
        "URL": "gs:///${gsTriggerBucket}/StorageMirror/Errors/"
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
- gs://${gsTriggerBucket}/data/p1/events.csv


Output:
- s3://${gsTriggerBucket}/data/p1/events.csv