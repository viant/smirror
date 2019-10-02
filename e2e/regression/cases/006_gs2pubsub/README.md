#### Scenario:

Mirror suffixed *.csv data from gs://${gsTriggerBucket}/data/p1 to s3://${s3TriggerBucket}/data

#### Input:

Configuration:

* Global Config: [@config,json](../../../config/gs.json)
* Route:

[@routes,json](rule.json)
```json
[
  {
    "Prefix": "/data/p6",
    "Suffix": ".csv",
    "Dest": {
      "Topic": "smirrorTopic"
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
        "URL": "gs:///${gsTriggerBucket}/e2e-mirror/errors/"
      }
    ],
    "PreserveDepth": 1
  }
]
```

* Trigger:

* event Type: google.storage.object.finalize
* resource: projects/_/buckets/${gsTriggerBucket}
* entryPoint: Fn
* environmentVariables:
  - LOGGING: 'true'
  - CONFIG: gs://${gsTriggerBucket}/e2e-mirror/config/mirror.json
 


Data:
- gs://${gsTriggerBucket}/data/p1/events.csv


Output:
- s3://${gsTriggerBucket}/data/p1/events.csv