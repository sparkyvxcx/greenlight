# greenlight

A movie api service

## Getting started

Set up Database connection string:

```shell
export GREENLIGHT_DSN="postgres://greenlight:pa55word@10.5.0.105/greenlight?sslmode=disable
```

Bring up the service:

```shell
go run ./cmd/api/ -dsn $GREENLIGHT_DSN -limiter-switch=false
```

Create example movies:

```shell
bash create_records.sh
```
