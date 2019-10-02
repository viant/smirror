# Deployment


The following document describes global shared storage mirror deployments for various data transfer processes, with one
SMirror lambda and cloud functions per bucket.

### Google Storage layout:

The following google storage layout is used:


##### Storage Mirror bucket

This bucket stores all configuration files:

**${gcp.projectID}-smirror:**


```bash
    /
    | - config
    |      |- smirror.json
    |      |- processes
    |      |     | - process1_routes.json
    |      |     | - processN_routes.json        
        
```            


Where:

[@smirror.json](usage/gcp/smirror.json)

```json
{

  "CheckInMs": 60000,
  "BaseURL": "gs://${gcp.projectID}-smirror/config/routes/"
}
```

and routes files store JSON array with process routes.

[@process1_routes.json](usage/gcp/process1_routes.json)
```json
[
  {
    "Source": { 
      "Prefix": "/data/",
      "Suffix": ".csv.gz"
    },
    "Dest": {
      "URL": "s3://destBucket/data",
      "Credentials": {
        "URL": "gs://sourceBucket/secret/s3-cred.json.enc",
        "Key": "projects/my_project/locations/us-central1/keyRings/my_ring/cryptoKeys/my_key"
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
        "URL": "gs://myProject-smirro-ops/data/errors/"
      }
    ],
    "Codec": "gzip",
    "PreserveDepth": 1
  }
]
```


##### Operational bucket

This bucket stores all processed, error files. 

**${gcp.projectID}-smirror-ops:**


##### Inbound mirror bucket 

This bucket stores all data that needs to be ingested to Big Query, 

**${gcp.projectID}-smirror-inbound**


# Deployment

You can deploy the described infrastructure with SMirror cloud function with [endly](https://github.com/viant/endly/) automation runner.


TODO add endly workflow
```bash

```
