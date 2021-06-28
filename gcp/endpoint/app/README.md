# GCP Smirror standalone app

This stand alone application use Pubsub topic to receive smirror storage notifiction

## Deployment



The following scenario uses



```bash
  export DEBUG_MSG=true
  export GCLOUD_PROJECT=my-gcp-project
  export GOOGLE_APPLICATION_CREDENTIALS=myGoogle.secret
  
  ### subscriber config
  export APP_CONFIG = '{"Topic":"my-topic"})'
  ### original smirror config (shared with lambda)
  export CONFIG = 'gs://myConfigBucket/StorageMirror/config.json'
  nohup ./subscriber &
```
