#### Scenario:

Mirror suffixed *.csv data from s3://${s3TriggerBucket}/data/p1 to ${s3Queue}

#### Input:

Configuration:

* Global Config: [@config,json](../../../config/s3.json)
* Rule

[@routes,json](rule.json)
```json
[
  {
    "Source": {
      "Prefix": "/data/p11",
      "Suffix": ".csv"
    },
    "Dest": {
      "Queue": "$s3Queue"
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
        "URL": "s3:///${s3OpsBucket}/StorageMirror/errors/"
      }
    ]
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
- s3://${s3TriggerBucket}/data/p1/events.csv