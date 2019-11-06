#### Scenario:

Mirror suffixed *.csv data from s3://${s3TriggerBucket}/data/p9 to gs://${gsDestBucket}/
Source file is larger then lambda memory, but streming threshold is met

#### Input:

Configuration:

* Global Config: [@config,json](../../../config/s3.json)
```json
{
  "Streaming": {
    "PartSizeMb": 5,
    "ThresholdMb": 10
  }
}
```


* Rule

[@routes,json](rule.json)
```json
[
  {
    "Source": {
      "Prefix": "/data/p9"

    },
    "Dest": {
      "URL": "gs://${gsDestBucket}/",
      "Credentials": {
        "Parameter": "storagemirror.gcp",
        "Key": "alias/smirror"
      }
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

* event Type: s3:ObjectCreated:*
* resource: ${s3TriggerBucket}
* environmentVariables:
  - LOGGING: 'true'
  - CONFIG: s3://${s3ConfigBucket}/StorageMirror/config.json
 
Data:
- s3://${s3TriggerBucket}/data/p9/events.csv


Output:
- gs://${gsDestBucket}/data/p9/events.csv/p9/events.csv