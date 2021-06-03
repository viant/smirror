# AWS Smirror standalone app

This stand alone application use AWS SQS to receive smirror storage notifiction

## Deployment



The following scenario uses

- ~/.aws/credentials
- ~/.aws/config


```bash
  export AWS_SDK_LOAD_CONFIG=true
  ### subscriber config
  export APP_CONFIG = '{"Queue":"my-overflow-queue"})'
  ### original smirror config (shared with lambda)
  export CONFIG = 's3://myConfigBucket/StorageMirror/config.json'
  ### Fake Lambda identity
  export AWS_SDK_LOAD_CONFIG='SMirror'
  nohup ./subscriber &
```
