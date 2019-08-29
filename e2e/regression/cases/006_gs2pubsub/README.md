#### Scenario:

Mirror suffixed *.csv data from gs://${gsBucket}/data/p1 to smirrorTopic topic, where each message is max 10 lines.

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
            "URL": "gs:///${gsBucket}/e2e-mirror/errors/"
          }
        ]
      },
      "FolderDepth": 1
    }
```


* Trigger:

* event Type: google.storage.object.finalize
* resource: projects/_/buckets/${gsBucket}
* entryPoint: Fn
* environmentVariables:
  - LOGGING: 'true'
  - CONFIG: gs://${gsBucket}/e2e-mirror/config/mirror.json
 
Data:
- gs://${gsBucket}/data/p1/events.csv

Output:
- s3://${gsBucket}/data/p1/events.csv

