#### Scenario:

Mirror suffixed *.csv data from gs://${s3ExternalBucket}/data/p1 to gs://${gsTriggerBucket}/data

#### Input:


**StorageMirrorCron**

Configuration:

* Global Config: [@config,json](../../../config/s3Cron.json)

* Rule:  [@cron,json](cron.json)




**StorageMirror**

Configuration:

* Global Config: [@config,json](../../../config/s3.json)

* Rule:  [@route,json](rule.json)
```json
[
  {
    "Prefix": "/data/p7",
    "Suffix": ".csv",
    "Source": {

      "CustomKey": {
        "Parameter": "storagemirror.customer001CustomKey",
        "Key": "alias/storagemirror"
      }
    },
    "Dest": {
      "URL": "gs://${gsDestBucket}/data",
      "Credentials": {
        "Parameter": "storagemirror.gcp",
        "Key": "alias/storagemirror"
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
    ],
    "PreserveDepth": 1
  }

]
```
