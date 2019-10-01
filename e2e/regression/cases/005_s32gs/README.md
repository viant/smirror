#### Scenario:

Mirror suffixed *.csv data from gs://${gsTriggerBucket}/data/p1 to s3://${s3TriggerBucket}/data

#### Input:

Configuration:

* Route:

[@config,json](../../../config/s3.json)
```json
     {
         "Prefix": "/data/p5",
         "Suffix": ".csv",
         "Dest": {
           "URL": "gs://${s3TriggerBucket}/data",
           "Credentials": {
             "Parameter": "smirror.gs",
             "Key": "alias/smirror"
           }
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
               "URL": "s3:///${gsTriggerBucket}/e2e-mirror/errors/"
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