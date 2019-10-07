# Deployment

- [Cloud storage layout](#cloud-storage-layout)
    - [Configuration Bucket](#configuration-bucket)
    - [Operation bucket](#operation-bucket)
    - [Trigger bucket](#trigger-bucket)
    - [Mirrored bucket](#mirrored-bucket)
-[Deployment](#deployment)
    - [Google Cloud Function](#google-cloud-function)

The following document describes global shared storage mirror deployments for various data transfer processes, with one
SMirror lambda and cloud functions per bucket.


### Cloud storage layout

The following google storage layout is used for universal deployment

- [Configuration Bucket](#configuration-bucket)
- [Operational bucket](#operational-bucket)
- [Trigger bucket](#trigger-bucket)
- [Mirrored bucket](#mirrored-bucket)


##### Configuration Bucket

This bucket stores all configuration files:

**${ConfigBucket}:**  

```bash
    /
    | - StorageMirror
    |      |- config.json
    |      |- dataflow
    |      |     | - route_rule1.json
    |      |     | - route_ruleN.json        
        
```            



Where:

[@smirror.json](usage/gcp/config.json)

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


##### Operation bucket

This bucket stores all processed, error files. 

**${OpsBucket}:**

```bash
    /
    | - StorageMirror
    |      |- errors
    |      |- processed
    |      |- replayed
       
```            

##### Trigger bucket 

This bucket stores all data that needs to be mirror 

**${TriggerBucket}**



```bash
    /
    | - data
    |      |- idfa
                |- dataXXX.csv.gz 
```    




##### Mirrored bucket 

This bucket stores all data that was mirrored from other cloud storage 

**${DestBucket}**


```bash
    /
    | - data
    |      |- idfa
                |- dataXXX.csv.gz 
```    


## Deployment

#### Google Cloud Function 

###### deploy with endly cli

```bash
endly authWith=myGCPSecretFile deploy.yaml
```
_where:_
- [@deploy.yaml](gcp/deploy.yaml)



###### deploy with gcloud

```bash
git chckout https://github.com/viant/smirror.git
cd smirror
unset GOPATH
export GO111MODULE=on
go mod vendor

gcloud functions deploy MyGsBucketToS3Mirror --entry-point StorageMirror \ 
    --trigger-resource ${gsTriggerBucket} 
    --trigger-event google.storage.object.finalize \
    --set-env-vars=LOGGING=true,CONFIG=gs://gsTriggerBucket/mirror/config/gs.json \
    --memory=512M \
    --timeout=540s \
    --runtime=go111 
```

