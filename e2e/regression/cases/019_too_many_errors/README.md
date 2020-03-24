#### Scenario:

Source file is in a CSV with **Schema** option set with CSV format, only valid line will be 

 
#### Input:

Configuration:

* Global Config: [@config,json](../../../config/gs.json)
* Rule

[@rule.yaml](rule.yaml)
```yaml
Source:
  Prefix: "/data/p17"
  Suffix: ".json"
Dest:
  URL: gs://${gsDestBucket}/
Schema:
  Format: JSON
  Fields:
    - Name: id
      DataType: int
    - Name: Name
      DataType: string
    - Name: Timestamp
      DataType: time
      SourceDateFormat: 'yyyy-MM-dd hh:mm:ss'
    - Name: segement
      DataType: int

OnSuccess:
  - Action: delete
```



Data:
- [inout data files](data/prepare)


Output:
- [mirrored data files](data/expect)
