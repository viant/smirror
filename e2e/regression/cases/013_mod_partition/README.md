#### Scenario:

Mirror suffixed *.csv data from gs://${gsTriggerBucket}/data/p13 to pubsub topic suffixed by partition
Source data in partition on 0 indexed filed, with modulo 2 
 
#### Input:

Configuration:

* Global Config: [@config,json](../../../config/gs.json)
* Rule

[@routes,json](rule.json)
```json
[
  {
    "Source": {
      "Prefix": "/data/p13",
      "Suffix": ".csv"
    },
    "Dest": {

      "Topic": "${gsPrefix}_storage_mirror_p$partition"
    },
    "Split": {
      "MaxSize": 1048576,
      "Partition": {
          "FieldIndex": 0,
          "Mod": 2
      },
      "Template": "$name_$chunk_$partition"
    },
    "OnSuccess": [
      {
        "Action": "delete"
      }
    ],
    "OnFailure": [
      {
        "Action": "move",
        "URL": "gs:///${gsOpsBucket}/StorageMirror/errors/"
      }
    ],
    "PreserveDepth": 1
  }
]
```

* Trigger:

* event Type: s3:ObjectCreated:*
* resource: ${s3TriggerBucket}
* environmentVariables:
  - LOGGING: 'true'
  - CONFIG: s3://${s3ConfigBucket}/StorageMirror/config.json
 
Data:
- s3://${s3TriggerBucket}/data/p13/events.csv


Output:
- **Topic**: ${gsPrefix}_storage_mirror_p$partition
     