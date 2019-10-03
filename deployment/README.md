# Deployment


The following document describes global shared storage mirror deployments for various data transfer processes, with one
SMirror lambda and cloud functions per bucket.

### Google Storage layout:

The following google storage layout is used:


##### Storage Mirror bucket

This bucket stores all configuration files:



**${configBucket}:**

```bash
    /
    | - StorageMirror
    |      |- config.json
    |      |- dataflow
    |      |     | - route_rule1.json
    |      |     | - route_ruleN.json        
        
```            


Where:

[@smirror.json](usage/gcp/smirror.json)

```json
{

  "CheckInMs": 60000,
  "BaseURL": "gs://${configBucket}/StorageMirror/dataflow/"
}
```

and routes files store JSON array with process routes.

[@route_rule1.json](usage/gcp/route_rule1.json)
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

**XXXXX-ops:**


##### Trigger bucket (inbound) 

This bucket stores all data that needs to be mirror 

**storagemirror-inbound**


##### Mirrored bucket (outbound) 

This bucket stores all data that was mirrored from other cloud storage 

**storagemirror-outbound**



