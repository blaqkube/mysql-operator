# MySQL Agent

The MySQL agent is used by the Operator to perform various operations. What
follows show how to test those operations from outside.

## Database Backup

The commands below perform a backup and send it to a S3 bucket:

```shell
export AWS_ACCESS_KEY_ID=AKIAXXXXXXXXXXXXXXXXXXX
export AWS_SECRET_ACCESS_KEY=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
export AWS_REGION=us-east-1

curl -H 'Accept: application/json' -H 'Content-Type: application/json' -X POST localhost:8080/backup \
  -d'{"s3access":{"bucket":"docs.blaqkube.io","path":"/location","credentials":
  {"aws_access_key_id": "'$AWS_ACCESS_KEY_ID'", "region": "'$AWS_REGION'",
  "aws_secret_access_key":"'$AWS_SECRET_ACCESS_KEY'"}}}' | tee output.json

BACKUP=$(jq -r ".timestamp" output.json)

curl -H 'Accept: application/json' localhost:8080/backup/$BACKUP |jq
```

