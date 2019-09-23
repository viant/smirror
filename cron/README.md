# Scheduled mirroring

This package implement scheduler that pull external URL to simulate cloud storage event
for new file arrival. It can be used with 3rd party partners where there is 
no direct cloud event integration.

# Configuration

[@config.json](usage/basic.json)
```json
{
  "MetaURL": "s3://myopsBucket/smirror/cron/meta.json",
  "TimeWindow": {
    "DurationInSec": 720
  },
  "Resources": [
    {
      "URL": "s3://externalBucket/data/2019/",
      "DestFunction": "myLambdaFun",
      "CustomKey": {
        "Parameter": "smirror.partnerX.customKey",
        "Key": "smirror"
      }
    }
  ]
}
```

- **MetaURL** stores all process files. In example above all files with modified time within 2 * 720 sec
- **Resources** list of resources with source base URL and DestFunction function
- **CustomKey** kms key name and ssm parameters storing [AES256Key](../config/key.go) encrypted value.
- **Credentials**  kms key name and ssm parameters storing encrypted credentials

