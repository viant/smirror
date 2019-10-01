#### Scenario:

Mirror data suffixed *.csv from gs://${gsTriggerBucket}/data/p1  to s3://${s3TriggerBucket}/data with compressed chunks (each no more than 10 rows)

#### Input:

Configuration:

* Route:

[@config,json](../../../config/gs.json)
```json
  {
       "Prefix": "/data/p4",
       "Suffix": ".csv",
       "Dest": {
         "URL": "s3://${s3TriggerBucket}/data",
         "Credentials": {
           "URL": "gs://${gsTriggerBucket}/e2e-mirror/secret/s3-mirror.json.enc",
           "Key": "projects/${gcpProject}/locations/us-central1/keyRings/gs_mirror_ring/cryptoKeys/gs_mirror_key"
         }
       },
       "Codec": "gzip",
       "Split": {
         "MaxLines": 10,
         "Template": "%s_%05d"
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
             "URL": "gs:///${gsTriggerBucket}/e2e-mirror/errors/"
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