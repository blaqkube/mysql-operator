# Integration test

Integration tests are part of the set of tests you should run in order to
validate your application. For the purpose of S3, you should create a `.env`
file like the one below:

```conf
BACKUP_AWS_ACCESS_KEY_ID=
BACKUP_AWS_SECRET_ACCESS_KEY=
BACKUP_AWS_REGION=us-east-1
BACKUP_BUCKET=mybucket
BACKUP_LOCATION=/tmp/demo.txt
```

To run the tests, execute the command like the one below:

```shell
go test . -v -tags=integration
```