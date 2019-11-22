#### Scenario:

Mirror suffixed *.csv data from s3://${s3TriggerBucket}/data/p5 to gs://${gsDestBucket}/data

#### Input:

Configuration:

* Global Config: [@config,json](../../../config/s3.json)

* Rule

[@rule.json](rule.json)
```json
 [
   {
     "Prefix": "/data/p5",
     "Suffix": ".csv",
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
         "URL": "s3:///${s3OpsBucket}/StorageMirror/errors/"
       }
     ],
     "PreserveDepth": 1
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