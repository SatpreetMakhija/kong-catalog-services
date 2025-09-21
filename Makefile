init-dev-db:
	docker run -d -e POSTGRES_PASSWORD=admin -p 5432:5432 postgres

seed-dev-db:
	psql "postgres://postgres:admin@localhost:5432/postgres?sslmode=disable" -f ./sql/setup.sql

start-dev-server:
	DB_HOST=localhost DB_PORT=5432 DB_USERNAME=postgres DB_PASSWORD=admin DB_DATABASE=service_catalog go run main.go	
