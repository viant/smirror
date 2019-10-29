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
        "URL": "s3:///${s3OpsBucket}/StorageMirror/Errors/"
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
- s3://${s3TriggerBucket}/data/p9/events.csv


Output:
- gs://${gsDestBucket}/data/p9/events.csv