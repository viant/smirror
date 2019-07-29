#### Scenario:

Data from gs://${gsBucket}/data/p1 and suffixed *.csv is chunked (by max rows=10) to s3://${s3Bucket}/data

#### Input:

Configuration:

* Route:

[@config,json](../../../config/gs_to_s3.json)
```json
 {
      "Prefix": "/data/p3",
      "Suffix": ".csv",
      "DestURL": "s3://${s3Bucket}/data",
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
            "URL": "gs:///${gsBucket}/e2e-mirror/errors/"
          }
        ]
      },
      "FolderDepth": 1
    }
```
 
* Secret:

```json
    {
      "Provider": "gcp",
      "TargetScheme": "s3",
      "URL": "gs://${gsBucket}/e2e-mirror/secret/s3-mirror.json.enc",
      "Key": "projects/${gcpProject}/locations/us-central1/keyRings/gs_mirror_ring/cryptoKeys/gs_mirror_key"
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