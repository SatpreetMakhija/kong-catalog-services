# Service Catalog Server

## Assumptions
We make the following assumptions when defining the data model and the API endpoints for Catalog Service:
- The name of a service is not a unique identifier. There can exist two services with the same name.
- The version value is of the form `vx.y.z` where `x,y,z` are integer values. Example: `v1.2.3`.
- The description of the service is less than or around 1500 words. The choice of persistence layer depends on this assumption. Note, longer descriptions can still be stored as exceptions but if all services' descriptions are longer, we might need to consider other solutions such as blob storage.

## API Definition
We define the following endpoints to serve our use case.
#### Fetch a particular service
```GET /api/v1/services/:id```
Here, the `id` is the unique ID of the service. The response is a JSON object with the following schema:
```
{
    "id": string,
    "name": string,
    "description": string,
    "version": string
}
```
#### Search service(s)
```GET /api/v1/services?```

The following URL parameters are supported when searching for services:
- `version`: Matches exact version of services. Eg: "v1.2.3"
- `name`: Matches approximate names. Eg: If service name is "Monitoring", "Monitor" will also yield the result.
- `keyword`: Keyword based search. Eg: "monitoring services with metadata"
- `sort`: Comma separated values to sort the result. Eg: `sort=name,-version` will first sort the result by `name` and then `version`. Note, addition of char `-` before a sort value means sort in descending order else ascending order. Supported values: `name`, `version`
- `page`: Enables pagination. Default value `1`.
- `page_size`: Number or results per page. Default value `10`.

The response is a JSON object with the following schema:
```
{
    "items" [array]{"id": string, "name": string, "description": string, "version": string}
    "page": int,
    "page_size" int 
}
```
Note, in the response object `page` and `page_size` is returned to help make the request for the next page. Increase the `page` value by 1 and keep `page_size` same in the next request.

## Solution overview
The catalog service server runs as an HTTP server written in golang. We've used PostgresSQL to store the catalog services. The service catalog item easily forms into a tabular form. Therefore, a tabular persistence layer is ideal. Since, keyword based search is the norm when querying for catalog services, we've created a column in `services` table of type [tsvector](https://www.postgresql.org/docs/current/datatype-textsearch.html#DATATYPE-TSVECTOR) to help with full text search.



## Guide to run the server

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
## Future work
- For full text search, we're using lexical logic not semantic logic. If the users' search dictates more free flowing search where exact keywords don't match, then, we can provide semantic search where the description of the service will be used to create vector embeddings and searches will be based on vector embeddings' difference.
- Define CORS, TLS security, authentication layer in the server.
- Write unit tests and integration tests.
Given the time constraint, we've skipped the ones mentioned above for now.
