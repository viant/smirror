#### Scenario:

Mirror data suffixed *.csv from gs://${gsTriggerBucket}/data/p1  to s3://${s3DestBucket}/data with compressed chunks (each no more than 10 rows)

#### Input:

Configuration:

* Global Config: [@config,json](../../../config/gs.json)

* Route:

[@config,json](../../../config/gs.json)
```json
  {
       "Prefix": "/data/p4",
       "Suffix": ".csv",
       "Dest": {
         "URL": "s3://${s3TriggerBucket}/data",
         "Credentials": {
           "URL": "gs://${gsConifgBucket}/Secrets/s3-mirror.json.enc",
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
             "URL": "gs:///${gsTriggerBucket}/StorageMirror/errors/"
           }
         ]
       },
       "PreserveDepth": 1
     }
```
 

* Trigger:

* event Type: google.storage.object.finalize
* resource: projects/_/buckets/${gsTriggerBucket}
* entryPoint: StorageMirror
* environmentVariables:
  - LOGGING: 'true'
  - CONFIG: gs://${gsConifgBucket}/StorageMirror/config.json
 


Data:
- gs://${gsTriggerBucket}/data/p1/events.csv


Output:
- s3://${gsTriggerBucket}/data/p1/events.csv