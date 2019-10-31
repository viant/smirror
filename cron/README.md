# Scheduled mirroring

This package implement scheduler that pull external URL to simulate cloud storage event
for new file arrival. It can be used with 3rd party partners where there is 
no direct cloud event integration.

# Configuration

**Config:**

[@rule.json](usage/config.json)
```json
{
  "MetaURL": "s3://myopsBucket/smirror/cron/meta.json",
  "TimeWindow": {
    "DurationInSec": 720
  },
  "Resources": {
    "CheckInMs": 60000,
    "BaseURL": "s3://myopsBucket/smirror/cron/rules"
  }
}

```

**Resource Rule:**

[@rule.json](usage/rule.json)
```json
[
  {
    "Source": {
      "URL": "s3://externalBucket/data/2019/",
      "Suffix": ".zip",
      "CustomKey": {
        "Parameter": "smirror.partnerX.customKey",
        "Key": "storagemirror"
      }
    },
    "Dest": {
      "URL": "s3://triggerBucket/data/2019/"
    }
  }
]
```

- **MetaURL** stores all process files. In example above all files with modified time within 2 * 720 sec
- **Resources** baseURL for resource rules or list of resources rules with source base URL and Dest function
- **CustomKey** kms key name and ssm parameters storing [AES256Key](../config/key.go) encrypted value.
- **Credentials**  kms key name and ssm parameters storing encrypted credentials

