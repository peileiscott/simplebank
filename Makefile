.PHONY: postgres createdb dropdb migrateup migratedown

postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=password -d postgres

createdb:
	docker exec -it postgres createdb --username=postgres simplebank

dropdb:
	docker exec -it postgres dropdb --username=postgres simplebank

migrateup:
	migrate -path db/migration -database "postgresql://postgres:password@localhost:5432/simplebank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:password@localhost:5432/simplebank?sslmode=disable" -verbose down