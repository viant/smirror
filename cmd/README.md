## StorageMirror - command line transfer/transformer

Stand alone data transfer/transformer.

### Introduction

Storage Mirror transfers data using  [transfer rule](../../smirror/README.md#rule).
For each source datafile an event is triggered to local smirror process.
Storage Mirror client supports all feature of serverless smirror, on top of that it also support constant data streaming option.

When transferring data smirror process manages history for all successfully process file, to avoid processing the same file more than once.
By default only streaming mode stores history file in file:///${env.HOME}/.smirror location, otherwise memory filesystem is used.

### Installation

##### OSX

```bash
wget https://github.com/viant/smirror/releases/download/v1.0.0/smirror_osx_1.0.0.tar.gz
tar -xvzf smirror_osx_1.0.0.tar.gz
cp smirror /usr/local/bin/
```

##### Linux

```bash
wget https://github.com/viant/smirror/releases/download/v1.0.0/smirror_linux_1.0.0.tar.gz
tar -xvzf smirror_linux_1.0.0.tar.gz
cp smirror /usr/local/bin/
```

### Usage  


**Data ingestion rule validation**

To validate rule use -V option.

```bash
smirror -r='myRuleURL' -p=myProject -V
smirror -r='gs://MY_CONFIG_BUCKET/smirror/Rules/myrule.yaml' -V
smirror -s=mydatafile -d='myProject:mydataset.mytable' -V
```

##### Simple data transfer

```bash
smirror --s=sourceURL --d=destURL
```

##### Split data into smaller data files 

```bash
smirror -r=split.yaml -s=/tmp/data  
```

where 
- [split.yaml](usage/split.yaml)
```yaml
Source:
  Prefix: "/tmp/data/"
  Suffix: ".csv"
Dest:
  URL: '/tmp/transformed'
Split:
  MaxLines: 10
  Template: "%s_%05d"
OnSuccess:
  - Action: delete
PreserveDepth: 1
```

##### Transfer with data quality check

```bash
smirror -r=schema.yaml -s=/tmp/data  
```

where 
- [schema.yaml](usage/schema.yaml)
```yaml
Source:
  Prefix: "/tmp/data/"
  Suffix: ".csv"
Dest:
  URL: '/tmp/transformed'
Schema:
  MaxBadRecords: 10
  Format: CSV
  Fields:
    - Name: id
      Position: 0
      DataType: int
    - Name: Name
      Position: 1
      DataType: string
    - Name: Timestamp
      Position: 2
      DataType: time
      SourceDateFormat: 'yyyy-MM-dd hh:mm:ss'
    - Name: segement
      Position: 3
      DataType: int
PreserveDepth: 0
OnSuccess:
  - Action: delete
```


##### Transcoder transfer

```bash
smirror -r=avro.yaml -s=/tmp/data  
```

where

- [avro.yaml](usage/avro.yaml) is the rule
```yaml
Source:
  Prefix: "/tmp/data/"
  Suffix: ".csv"
Dest:
  URL: '/tmp/transformed'
Transcoder:
  Source:
    Format: CSV
    Fields:
      - id
      - event_type
      - timestamp
    HasHeader: true
  Dest:
    Format: AVRO
    SchemaURL: schema.avsc
    RecordPerBlock: 10
OnSuccess:
  - Action: delete
PreserveDepth: 1
```
- [schema.avsc](schema.avsc)
```avroidl
{
    "namespace": "my.namespace.com",
    "type":	"record",
    "name": "root",
    "fields": [
        { "name": "id", "type": "int"},
        { "name": "event_type", "type": "string"},
        { "name": "timestamp", "type":["null",{"type":"long","logicalType":"timestamp-millis"}]}
    ]
}
```


### Publishing data to message bus

```yaml
 smirror -r=pubsub.yaml -s=/tmp/data 
```

where
-[pubsub.yaml](usage/pubsub.yaml)
```yaml
Source:
  Prefix: "/tmp/data/"
  Suffix: ".csv"
Dest:
  Topic: "mytopic"
Split:
  MaxLines: 10
OnSuccess:
  - Action: delete
PreserveDepth: 1
```

### Transferring zip/tar assets

```yaml
 smirror -r=codec.yaml -s=/tmp/data 
```
where

- [codec.yaml](usage/codec.yaml)
```yaml
Source:
  Prefix: "/tmp/data/"
  Suffix: ".zip"
Uncompress: true
Dest:
  URL: /tmp/transformed/
Codec: gzip
PreserveDepth: -1
```

### Authentication

smirror client can use one the following auth method

1. With Google Service Account Secrets

```bash
export GOOGLE_APPLICATION_CREDENTIALS=myGoogle.secret
```

2. [AWS SDK config file](https://docs.aws.amazon.com/sdk-for-go/api/aws/session/)  

```bash
export AWS_SDK_LOAD_CONFIG=true
```




Help: 

```bash
smirror -h
```


