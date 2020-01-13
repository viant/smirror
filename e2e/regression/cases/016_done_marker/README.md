#### Scenario:

Source is CSV file, that can be only transfer when done marker file is present.
Once done marker file is present the whole folder is replayed except marker file.  

#### Input:

Configuration:

* Global Config: [@config,json](../../../config/gs.json)
* Rule
[@rule.yaml](rule.yaml)
```yaml
Source:
  Prefix: "/data/p016"
  Suffix: ".csv"
Dest:
  URL: gs://${gsDestBucket}/
DoneMarker: sucess.txt
OnSuccess:
  - Action: delete
```

