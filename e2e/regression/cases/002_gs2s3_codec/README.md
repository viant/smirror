#### Scenario:

Mirror compressed data from gs://${gsTriggerBucket}/data/p2 and suffixed *.csv.gz to s3://${s3DestBucket}/data

#### Input:

**Configuration**:

* Global Config: [@config,json](../../../config/gs.json)
* Rule: [@rule,json](rule.json)

```json
   {
       "Prefix": "/data/p2",
       "Suffix": ".csv.gz",
       "Dest": {
           "URL": "s3://${s3TriggerBucket}/data",
           "Credentials": {
              "URL": "gs://${gsConifgBucket}/Secrets/s3-mirror.json.enc",
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
             "URL": "gs:///${gsTriggerBucket}/StorageMirror/errors/"
           }
         ]
       },
       "PreserveDepth": 1
     }
```
* **Trigger**:

    * event Type: google.storage.object.finalize
    * resource: projects/_/buckets/${gsTriggerBucket}
    * entryPoint: StorageMirror
    * environmentVariables:
      - LOGGING: 'true'
      - CONFIG: gs://${gsConifgBucket}/StorageMirror/config.json
 


**Data**:
- gs://${gsTriggerBucket}/data/p2/events.csv


#### Output

**Data**
- s3://${s3DestBucket}/data/p1/events.csv