#### Scenario:

Data from gs://${gsBucket}/data/p1 and suffixed *.csv is mirrored to s3://${s3Bucket}/data

#### Input:

Configuration:

* Route:

[@config,json](../../../config/s3_to_gs.json)
```json
    {
      "Prefix": "/data/p1",
      "Suffix": ".csv",
      "DestURL": "gs://${s3Bucket}/data",
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
 
* Secret:

```json
   {
     "Provider": "aws",
     "TargetScheme": "s3",
     "Parameter": "smirror.gs",
     "Key": "alias/smirror"
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