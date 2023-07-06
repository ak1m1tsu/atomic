# Atomic

## `Makefile` commands

| Command      | Description                                    |
|--------------|------------------------------------------------|
| `make lint`  | analysis code base with `golangci-lint`        |
| `make build` | builds service                                 |
| `make run`   | builds and runs service                        |
| `make clean` | removes `bin` folder                           |
| `make test`  | runs tests                                     |
| `make race`  | runs tests with `-race` flag                   |
| `make mock`  | generates mock data                            |
| `make up`    | ups service with database using docker compose |
| `make down`  | downs service with database                    |

## Configure for the local development

First, set the path to the configuration file in the environment variables

```bash
$ export CONFIG_PATH=./config/local.yaml
```

Next, add migrations to the database with [migrate](https://github.com/golang-migrate/migrate)

```bash
$ migrate -database DATABASE_DSN -path migrations up
```

> Make sure that the values in the `local.yaml` correspond to the values of your database

## Configure for the production

### Railway

The first thing to do is to create a project in the [Railway](https://railway.app/). Then add fork of this repository and Postgres DB, then in service settings set environment variables.

```bash
CONFIG_PATH = /config.yaml
HTTP_SERVER_PASSWORD = mypass
POSTGRES_DB = railway db name
POSTGRES_HOST = railway db host
POSTGRES_PASSWORD = railway db password
POSTGRES_PORT = railway db port
POSTGRES_USER = railway db user
```
