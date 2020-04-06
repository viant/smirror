#### Scenario:

In this scenario source data in XLSX format is automatically transcoder into AVRO.

#### Input:

* Rule

[@rule.yaml](rule/rule.yaml)
```yaml
Source:
  Prefix: "/data/p022/"
  Suffix: ".xlsx"
Dest:
  URL: gs://${gsDestBucket}/data
Transcoder:
  Source:
    Format: XLSX
  Dest:
    Format: AVRO
Autodetect: true
OnSuccess:
  - Action: delete
PreserveDepth: 1
Info:
  Workflow: e2e test case 1
```

Input:
- s3://${s3TriggerBucket}/data/p022/book.xlsx


Output:
- gs://${gsDestBucket}/data/p022/book.avro