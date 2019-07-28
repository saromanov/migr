# migr [![Go Report Card](https://goreportcard.com/badge/github.com/saromanov/migr)](https://goreportcard.com/report/github.com/saromanov/migr)

Tool for SQL migrations

For the start of working with migr, create a new directory for migrations
with Postgres driver

```
migr --run ss --driver postgres --dbname postgres --username postgres --password postgres --username postgres --host 127.0.0.1 --port 54
```

For applying migrations, please run following command

```
migr --run ss --driver postgres --dbname postgres --username postgres --password postgres --username postgres --host 127.0.0.1 --port 5432
```

For downgrade migrations

```
migr --down --driver postgres --dbname postgres --username postgres --password postgres --username postgres --host 127.0.0.1 --port 5432
```
