#### Scenario:

Data from gs://${gsTriggerBucket}/data/p1 and suffixed *.csv is chunked (by max rows=10) to s3://${s3DestBucket}/data

#### Input:

Configuration:

* Global Config: [@config,json](../../../config/gs.json)

* Rule

[@rule.yaml](rule.yaml)
```yaml
Source:
  Prefix: "/data/p3"
  Suffix: ".csv"
Dest:
  URL: s3://${s3DestBucket}/data
  Credentials:
    URL: gs://${gsConfigBucket}/Secrets/s3-mirror.json.enc
    Key: projects/${gcpProject}/locations/us-central1/keyRings/${gsPrefix}_ring/cryptoKeys/${gsPrefix}_key
Split:
  MaxLines: 10
  Template: "%s_%05d"
OnSuccess:
  - Action: delete
OnFailure:
  - Action: move
    URL: gs:///${gsOpsBucket}/StorageMirror/errors/
PreserveDepth: 1

````
 

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