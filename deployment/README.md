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
    "BaseURL": "gs://${configBucket}/StorageMirror/Rules",
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
        "URL": "gs://${opsBucket}/Errors/"
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
    |      |- xxxxx
                |- dataXXX.csv.gz 
```    



##### Mirrored bucket 

This bucket stores all data that was mirrored from other cloud storage 

**${destBucket}**


```bash
    /
    | - data
    |      |- xxxxx
                |- dataXXX.csv.gz 
```    


## Deployment


### Google Cloud Platform 

###### StorageMirror

To deploy StorageMirror cloud function with **endly** automation runner use the following workflow:

```bash
git chckout https://github.com/viant/smirror.git
cd smirror/deployment/smirror/gcp
endly deploy.yaml authWith=myGCPSecretFile.json
```

_where:_
- [@deploy.yaml](mirror/gcp/deploy.yaml)




To deploy StorageMirror cloud function with **gcloud** cli use the following commands:

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

###### Post deployment


######  Securing AWS credentials

To deploy StorageMirror cloud function with **endly** automation runner use the following workflow:

```bash
git chckout https://github.com/viant/smirror.git
cd smirror/deployment/smirror/gcp
endly secure.yaml authWith=myGCPSecretFile.json awsSecrets=awsSecret.json
```
_where:_
- [@secure.yaml](mirror/gcp/secure.yaml)
- [@awsSecret.json](mirror/gcp/awsSecrets.json)



To deploy secret aws secrets with **gcloud** cli use the following commands:

```bash

gcloud kms keyrings create my_ring --location us-central1
gcloud kms keys create my_key --location us-central1 \
  --keyring my_ring --purpose encryption

## Encrypt s3-cred.json

gcloud kms encrypt \
  --location=us-central1  \
  --keyring=my_ring \
  --key=my_key \
  --version=1 \
  --plaintext-file=s3-cred.json \
  --ciphertext-file=s3-cred.json.enc

## Upload encrypted version to google storage

gsutil cp s3-cred.json.enc gs://sourceBucket/secret/s3-cred.json.enc

```

######  Securing Slack credentials

To deploy StorageMirror cloud function with **endly** automation runner use the following workflow:

```bash
git chckout https://github.com/viant/smirror.git
cd smirror/deployment/smirror/gcp
endly secure_slack.yaml authWith=myGCPSecretFile.json slackSecrets=slackSecrets.json
```
_where:_
- [@secure_slack.yaml](mirror/gcp/secure_slack.yaml)
- [@slackSecrets.json](mirror/gcp/slackSecrets.json)



###### Test Rule With s3 destination 

```bash
git chckout https://github.com/viant/smirror.git
cd smirror/deployment/smirror/test
endly test.yaml authWith=myGCPSecrets.json awsCredentials=myAWSSecrets.json
```

_where:_
- [@rule.yaml](mirror/gcp/test/rule.json)
- [@test.yaml](mirror/gcp/test/test.yaml)


###### Test Rule With pubsub destination

```bash
git chckout https://github.com/viant/smirror.git
cd smirror/deployment/smirror/test
endly test_pubsub authWith=myGoogleSecrets.json
```
_where:_
- [@rule_pubsub.yaml](mirror/gcp/test/rule_pubsub.json)
- [@test_pubsub.yaml](mirror/gcp/test/test_pubsub.yaml)
 

###### StorageMirror Subscriber

```bash
## Set topic to bucket event notification
endly set_notification.yaml authWith=myGCPSecretFile bucket=${prefix}_pubsub_trigger topic=${prefix}_storage_mirror_trigger

## Deploy cloud function 
endly deploy_subscriber.yaml authWith=myGCPSecretFile topic=${prefix}_storage_mirror_trigger
```


_where:_
- [@set_notification.yaml](mirror/gcp/set_notification.yaml)
- [@deploy_subscriber.yaml](mirror/gcp/deploy_subscriber.yaml)
- _prefix_:  project name



###### StorageMonitor



###### StorageReplay




### Amazon Web Service 


##### StorageMirror

To deploy with endly automation runner use the following workflow:

```bash
endly deploy.yaml authWith=myAWSSecretFile
```

_where:_
- [@deploy.yaml](mirror/aws/deploy.yaml)
- [@privilege-policy.json](mirror/aws/privilege-policy.json)


##### Post deployment


#####  Secure GCP credentials (google secrets)

To secure google secrets with endly runner use the following workflow:

```bash
git chckout https://github.com/viant/smirror.git
cd smirror/deployment/mirror/aws
endly secure.yaml authWith=aws-e2e gcpSecrets=gcp-e2e
```

where:
[@secure.yaml](mirror/aws/secure.yaml)

To secure google secrets with aws cli  use following commands:

```bash
- aws kms create-key  
- aws kms create-alias --alias-name=smirror --target-key-id=KEY_ID
- aws ssm put-parameter \
    --name "storagemirror.gcp" \
    --value 'CONTENT OF GOOGLE SECRET HERE' \
    --type SecureString \
    --key-id alias/storagemirror

```


###### Test Mirror Rule 

```bash
git chckout https://github.com/viant/smirror.git
cd smirror/deployment/mirror/aws/test
endly test.yaml authWith=aws-e2e  gcpSecrets=gcp-e2e
```

where:
- [@test.yaml](mirror/aws/test/test.yaml)
- [@rule.json](mirror/aws/test/rule.json)




##### StorageMirror SQS Proxy

To deploy with endly automation runner use the following workflow:

```bash
git chckout https://github.com/viant/smirror.git
cd smirror/deployment/mirror/aws/sqs
endly deploy.yaml authWith=myAWSSecretFile
```
_where:_
- [@deploy.yaml](mirror/aws/sqs/deploy.yaml)
- [@privilege-policy.json](mirror/aws/sqs/privilege-policy.json)

###### Test Mirror Rule 

```bash
git chckout https://github.com/viant/smirror.git
cd smirror/deployment/mirror/aws/test
endly test_sqs.yaml authWith=aws-e2e  gcpSecrets=gcp-e2e
```

where:
- [@test_sqs.yaml](mirror/aws/test/test_sqs.yaml)
- [@rule.json](mirror/aws/test/rule.json)


##### StorageMirror SNS Proxy

To deploy with endly automation runner use the following workflow:

```bash
git chckout https://github.com/viant/smirror.git
cd smirror/deployment/mirror/aws/sns
endly deploy.yaml authWith=myAWSSecretFile
```
_where:_
- [@deploy.yaml](mirror/aws/sns/deploy.yaml)
- [@privilege-policy.json](mirror/aws/sns/privilege-policy.json)


##### Test Mirror Rule

```bash
git chckout https://github.com/viant/smirror.git
cd smirror/deployment/mirror/aws/test
endly test_sns.yaml authWith=aws-e2e  gcpSecrets=gcp-e2e
```

where:
- [@test_sns.yaml](mirror/aws/test/test_sqs.yaml)
- [@rule.json](mirror/aws/test/rule.json)




##### StorageMirror Cron

To deploy with endly automation runner use the following workflow:

```bash
git chckout https://github.com/viant/smirror.git
cd smirror/deployment/mirror/aws/sqs
endly deploy.yaml authWith=aws-e2e gcpSecrets=viant-e2e
```
_where:_
- [@deploy.yaml](mirror/aws/cron/deploy.yaml)
- [@privilege-policy.json](mirror/aws/cron/privilege-policy.json)

###### Test Cron Mirror Rule 

```bash
git chckout https://github.com/viant/smirror.git
cd smirror/deployment/mirror/aws/test
endly test_sqs.yaml authWith=aws-e2e  gcpSecrets=gcp-e2e
```

where:
- [@test_cron.yaml](mirror/aws/test/test_cron.yaml)
- [@cron_rule.json](mirror/aws/test/cron.json)
- [@rule.json](mirror/aws/test/rule_cron.json)



##### StorageMonitor

To deploy with endly automation runner use the following workflow:

```bash
git chckout https://github.com/viant/smirror.git
cd smirror/deployment/monitor/aws
endly deploy.yaml authWith=myAWSSecretFile
```
_where:_
- [@deploy.yaml](monitor/aws/deploy.yaml)
- [@privilege-policy.json](monitor/aws/privilege-policy.json)




##### StorageReplay

To deploy with endly automation runner use the following workflow:

```bash
git chckout https://github.com/viant/smirror.git
cd smirror/deployment/replay/aws
endly deploy.yaml authWith=myAWSSecretFile
```
_where:_
- [@deploy.yaml](replay/aws/deploy.yaml)
- [@privilege-policy.json](replay/aws/privilege-policy.json)


