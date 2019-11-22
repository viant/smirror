#### Scenario:

Source is CSV file, with **Recover** option set with CSV format,all record will be adjusted to specified filed count and mirrored to destination. 
  
 
#### Input:

Configuration:

* Global Config: [@config,json](../../../config/gs.json)
* Rule

[@rule.json](rule.json)
```json
[
  {
    "Source": {
      "Prefix": "/data/p015",
      "Suffix": ".json"
    },
    "Dest": {
      "URL": "gs://${gsDestBucket}/"
    },
    "Recover": {
        "Format": "CSV",
        "FieldCount": 3
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
- gs://${gsTriggerBucket}/data/p015/events.csv


Output:
- gs://${gsDestBucket}/data/p015/events.csv