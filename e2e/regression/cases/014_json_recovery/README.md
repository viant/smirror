#### Scenario:

Source file is in new line delimitered JSON, with **Recover** option set with JSON format, only valid JSON line will be mirrored to destination.
  
 
#### Input:

Configuration:

* Global Config: [@config,json](../../../config/gs.json)
* Rule

[@rule.json(rule.json)
```json
[
  {
    "Source": {
      "Prefix": "/data/p14",
      "Suffix": ".json"
    },
    "Dest": {
      "URL": "gs://${gsDestBucket}/"
    },
    "Recover": {
        "Format": "JSON"
    },
    "OnSuccess": [
      {
        "Action": "delete"
      }
    ]
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
- gs://${gsTriggerBucket}/data/p14/events.csv


Output:
- gs://${gsDestBucket}/data/p14/events.csv