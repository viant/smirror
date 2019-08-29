#### Scenario:

Mirror suffixed *.csv data from gs://${gsBucket}/data/p1 to s3://${s3Bucket}/data

#### Input:

Configuration:

* Route:

[@config,json](../../../config/s3.json)
```json
     {
         "Prefix": "/data/p5",
         "Suffix": ".csv",
         "Dest": {
           "URL": "gs://${s3Bucket}/data",
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
               "URL": "s3:///${gsBucket}/e2e-mirror/errors/"
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