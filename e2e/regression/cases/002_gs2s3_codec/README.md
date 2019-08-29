#### Scenario:

Mirror compressed data from gs://${gsBucket}/data/p2 and suffixed *.csv.gz to s3://${s3Bucket}/data

#### Input:

Configuration:

* Route:

[@config,json](../../../config/gs.json)
```json
   {
       "Prefix": "/data/p2",
       "Suffix": ".csv.gz",
       "Dest": {
           "URL": "s3://${s3Bucket}/data",
           "Credentials": {
              "URL": "gs://${gsBucket}/e2e-mirror/secret/s3-mirror.json.enc",
              "Key": "projects/${gcpProject}/locations/us-central1/keyRings/gs_mirror_ring/cryptoKeys/gs_mirror_key"
           }
       },
       "Codec": "gzip",
       "OnCompletion": {
         "OnSuccess": [
           {
             "Action": "delete"
           }
         ],
         "OnError": [
           {
             "Action": "move",
             "URL": "gs:///${gsBucket}/e2e-mirror/errors/"
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
- gs://${gsBucket}/data/p2/events.csv.gz


Output:
- s3://${gsBucket}/data/p2/events.csv.gz