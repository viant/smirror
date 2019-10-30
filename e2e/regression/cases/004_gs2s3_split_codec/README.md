#### Scenario:

Mirror data suffixed *.csv from gs://${gsTriggerBucket}/data/p1  to s3://${s3DestBucket}/data with compressed chunks (each no more than 10 rows)

#### Input:

Configuration:

* Global Config: [@config,json](../../../config/gs.json)

* Rule

[@config,json](../../../config/gs.json)
```json
  {
       "Prefix": "/data/p4",
       "Suffix": ".csv",
       "Dest": {
         "URL": "s3://${s3TriggerBucket}/data",
         "Credentials": {
           "URL": "gs://${gsConfigBucket}/Secrets/s3-mirror.json.enc",
           "Key": "projects/${gcpProject}/locations/us-central1/keyRings/${gsPrefix}_ring/cryptoKeys/${gsPrefix}_key"
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
             "URL": "gs:///${gsOpsBucket}/StorageMirror/Errors/"
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
  - CONFIG: gs://${gsConfigBucket}/StorageMirror/config.json
 


Data:
- gs://${gsTriggerBucket}/data/p1/events.csv


Output:
- s3://${s3TriggerBucket}/data/p1/events.csv