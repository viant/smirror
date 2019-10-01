#### Scenario:

Mirror suffixed *.csv data from gs://${gsTriggerBucket}/data/p1 to smirrorTopic topic, where each message is max 10 lines.

#### Input:

Configuration:

* Route:

[@config,json](../../../config/gs.json)
```json
  {
      "Prefix": "/data/p6",
      "Suffix": ".csv",
      "Dest": {
        "Topic": "smirrorTopic"
      },
      "Split": {
        "MaxLines": 10
      },
      "OnCompletion": {
        "OnSuccess": [
          {
            "Action": "delete"
          }
        ],
        "OnError": [
          {
            "Action": "move",
            "URL": "gs:///${gsTriggerBucket}/e2e-mirror/errors/"
          }
        ]
      },
      "FolderDepth": 1
    }
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

