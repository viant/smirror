#### Scenario:

Mirror suffixed *.csv data from gs://${gsTriggerBucket}/data/p1 to s3://${s3TriggerBucket}/data

#### Input:


**StorageMirrorCron**

Configuration:

* Global Config: [@config,json](../../../config/s3Cron.json)

* Route:

[@route,json](routes.json)




**StorageMirror**

Configuration:

* Global Config: [@config,json](../../../config/s3.json)

* Route:

[@route,json](routes.json)
```json
[
  {
    "Prefix": "/data/p7",
    "Suffix": ".csv",
    "Source": {

      "CustomKey": {
        "Parameter": "smirror.partnerXKey",
        "Key": "alias/smirror"
      }
    },
    "Dest": {
      "URL": "gs://${gsDestBucket}/data",
      "Credentials": {
        "Parameter": "smirror.gs",
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
        "URL": "s3:///${s3OpsBucket}/e2e-mirror/errors/"
      }
    ],
    "FolderDepth": 1
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