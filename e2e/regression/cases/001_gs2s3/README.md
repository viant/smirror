#### Scenario:

Mirror data from gs://${gsTriggerBucket}/data/p1 and suffixed *.csv to s3://${s3DestBucket}/data

#### Input:

**Configuration**:

* Global Config: [@config,json](../../../config/gs.json)
* Rule: [@rule,yaml](rule.yaml)

```yaml
Source:
  Prefix: "/data/p1/"
  Suffix: ".csv"
Dest:
  URL: s3://${s3DestBucket}/data
  Credentials:
    URL: gs://${gsConfigBucket}/Secrets/s3-mirror.json.enc
    Key: projects/${gcpProject}/locations/us-central1/keyRings/${gsPrefix}_ring/cryptoKeys/${gsPrefix}_key
OnSuccess:
  - Action: move
    URL: gs://${gsOpsBucket}/StorageMirror/processed/
  - Action: notify
    Title: Transfer $SourceURL done
    Message: success !!
    Channels:
      - "#e2e"
    Body: "$Response"
OnFailure:
  - Action: move
    URL: gs://${gsOpsBucket}/StorageMirror/errors/
PreserveDepth: 1
Info:
  Workflow: e2e test case 1
```
 

* **Trigger**:

    * event Type: google.storage.object.finalize
    * resource: projects/_/buckets/${gsTriggerBucket}
    * entryPoint: StorageMirror
    * environmentVariables:
      - LOGGING: 'true'
      - CONFIG: gs://${gsConfigBucket}/StorageMirror/config.json
 


**Data**:
- gs://${gsTriggerBucket}/data/p1/events.csv


#### Output

**Data**
- s3://${s3DestBucket}/data/p1/events.csv