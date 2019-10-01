#### Scenario:

Mirror suffixed *.csv data from s3://${s3TriggerBucket}/data/p5 to gs://${gsDestBucket}/data

#### Input:

Configuration:

* Global Config: [@config,json](../../../config/s3.json)

* Route:

[@routes,json](routes.json)
```json
 [
   {
     "Prefix": "/data/p5",
     "Suffix": ".csv",
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