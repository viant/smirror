# StorageMirror AWS proxy

Storage Mirror proxy propagates bucket events to destination.

Destination can be one of the following

1. **StorageMirror lambda function** - (sqs/sns) source bucket storage events is used to invoke StorageMirror Lambda.
2  **Trigger URL** - in this case source resource is copied over to trigger bucket, original bucket name is
prefixes destination object key. Assume that storage event was triggered from _s3://my-sqs-bucket/folder/asset1.txt_
and destination trigger is  _s3://xxx_, then source object is copy to _s3://xxx/my-sqs-bucket/folder/asset1.txt_
    

   