# Integration test

Integration tests are part of the set of tests you should run in order to
validate your application. For the purpose of GCP, you should create a `.env`
file like the one below:

```conf
BACKUP_BUCKET=mybucket
BACKUP_LOCATION=demo.txt
GOOGLE_APPLICATION_CREDENTIALS='{"type": "service_account", "project_id": ...}'
```

To run the tests, execute the command like the one below:

```shell
go test . -v -tags=integration
```