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

[@smirror.json](mirror/usage/gcp/config.json)

```json
{
  "Mirrors": {
    "BaseURL": "gs://${configBucket}/StorageMirror/dataflow",
    "CheckInMs": 60000
  }
}
```

and routes files store JSON array with process routes.

[@route_rule1.json](mirror/usage/gcp/route_rule1.json)
```json
[
  {
    "Prefix": "/data/",
    "Suffix": ".csv.gz",
    "Dest": {
      "URL": "s3://${destBucket}/data",
      "Credentials": {
        "URL": "gs://${configBucket}/StorageMirror/Secrets/s3-cred.json.enc",
        "Key": "projects/my_project/locations/us-central1/keyRings/my_ring/cryptoKeys/my_key"
      }
    },
    "OnSuccess": [
      {
        "Action": "move",
        "URL": "gs://${opsBucket}/processed/"
      }
    ],
    "OnFailure": [
      {
        "Action": "move",
        "URL": "gs://${opsBucket}/errors/"
      }
    ],
    "Codec": "gzip",
    "PreserveDepth": 1
  }
]
```


##### Operation bucket

This bucket stores all processed, error files. 

**${opsBucket}:**

```bash
    /
    | - StorageMirror
    |      |- errors
    |      |- processed
    |      |- replayed
       
```            

##### Trigger bucket 

This bucket stores all data that needs to be mirror 

**${triggerBucket}**


```bash
    /
    | - data
    |      |- idfa
                |- dataXXX.csv.gz 
```    



##### Mirrored bucket 

This bucket stores all data that was mirrored from other cloud storage 

**${destBucket}**


```bash
    /
    | - data
    |      |- idfa
                |- dataXXX.csv.gz 
```    


## Deployment



### Google Cloud Function 

###### StorageMirror

To deploy with endly automation runner use the following workflow:

```bash
endly deploy.yaml authWith=myGCPSecretFile
```

_where:_
- [@deploy.yaml](mirror/gcp/deploy.yaml)


To deploy with gcloud cli use the following commands:

```bash
git chckout https://github.com/viant/smirror.git
cd smirror
unset GOPATH
export GO111MODULE=on
go mod vendor

gcloud functions deploy MyGsBucketToS3Mirror --entry-point StorageMirror \ 
    --trigger-resource ${triggerBucket} 
    --trigger-event google.storage.object.finalize \
    --set-env-vars=LOGGING=true,CONFIG=gs://triggerBucket/mirror/config/gs.json \
    --memory=512M \
    --timeout=540s \
    --runtime=go111 
```


Testing deployment:

```bash
endly deploy.yaml authWith=myGCPSecretFile
```

_where:_
- [@deploy.yaml](mirror/gcp/deploy.yaml)


To deploy with gcloud cli use the following commands:

```bash
git chckout https://github.com/viant/smirror.git
cd smirror
unset GOPATH
export GO111MODULE=on
go mod vendor

gcloud functions deploy MyGsBucketToS3Mirror --entry-point StorageMirror \ 
    --trigger-resource ${triggerBucket} 
    --trigger-event google.storage.object.finalize \
    --set-env-vars=LOGGING=true,CONFIG=gs://triggerBucket/mirror/config/gs.json \
    --memory=512M \
    --timeout=540s \
    --runtime=go111 
```



###### StorageMirror Subscriber




###### StorageMonitor

###### StorageReplay



### AWS Lambda function 


###### StorageMirror

To deploy with endly automation runner use the following workflow:

```bash
endly deploy.yaml authWith=myAWSSecretFile
```

_where:_
- [@deploy.yaml](mirror/aws/deploy.yaml)
- [@privilege-policy.json](mirror/aws/privilege-policy.json)


**Testing deployment:**

-  Secure google storage credentials (google secrets)


```bash
cd $SmirrorRoot
cd deployment/mirror/aws
endly secure.yaml authWith=aws-e2e gcpSecrets=gcp-e2e
```

where:
[@secure.yaml](mirror/aws/secure.yaml)

-  **Test Mirror Rule** 

```bash
cd $SmirrorRoot
cd deployment/mirror/aws/rule
endly test.yaml authWith=aws-e2e  gcpSecrets=gcp-e2e
```

where:
- [@test.yaml](mirror/aws/rule/test.yaml)
- [@rule.json](mirror/aws/rule/rule.json)



###### StorageMirror SQS Proxy

To deploy with endly automation runner use the following workflow:

```bash
endly deploy.yaml authWith=myAWSSecretFile
```
_where:_
- [@deploy.yaml](mirror/aws/sqs/deploy.yaml)
- [@privilege-policy.json](mirror/aws/sqs/privilege-policy.json)

-  **Test Mirror Rule** 

```bash
cd $SmirrorRoot
cd deployment/mirror/aws/rule
endly test_sqs.yaml authWith=aws-e2e  gcpSecrets=gcp-e2e
```

where:
- [@test_sqs.yaml](mirror/aws/rule/test_sqs.yaml)
- [@rule.json](mirror/aws/rule/rule.json)



###### StorageMirror SNS Proxy


To deploy with endly automation runner use the following workflow:

```bash
endly deploy.yaml authWith=myAWSSecretFile
```
_where:_
- [@deploy.yaml](mirror/aws/sns/deploy.yaml)
- [@privilege-policy.json](mirror/aws/sns/privilege-policy.json)

**Test Mirror Rule**

```bash
cd $SmirrorRoot
cd deployment/mirror/aws/rule
endly test_sns.yaml authWith=aws-e2e  gcpSecrets=gcp-e2e
```

where:
- [@test_sns.yaml](mirror/aws/rule/test_sqs.yaml)
- [@rule.json](mirror/aws/rule/rule.json)


###### StorageMirror Cron



###### StorageMonitor


To deploy with endly automation runner use the following workflow:

```bash
endly deploy.yaml authWith=myAWSSecretFile
```
_where:_
- [@deploy.yaml](monitor/aws/deploy.yaml)
- [@privilege-policy.json](monitor/aws/privilege-policy.json)


###### StorageReplay


To deploy with endly automation runner use the following workflow:

```bash
endly deploy.yaml authWith=myAWSSecretFile
```
_where:_
- [@deploy.yaml](replay/aws/deploy.yaml)
- [@privilege-policy.json](replay/aws/privilege-policy.json)

