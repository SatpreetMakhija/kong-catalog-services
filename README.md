# Service Catalog Server

## Developer Guide

### Prerequisites
- Docker (We use docker to setup PostgresSQL instance)
- psql (We use psql to setup seed data in the PostgresSQL instance)

### How to start server locally?

The catalog service uses a PostgresSQL instance as its database. Use the commands below to first, spin up a local PostgresSQL instance and second, start catalog service server.
1. Start local Postgres container. Default username is `postgres` and password is `admin`.
```
make init-dev-db
```
2. Create required database, tables in Postgres instance and populate seed data.
```
make seed-dev-db
```
3. Start server locally
```
make start-dev-server
```

