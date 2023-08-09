export POSTGRES_PASSWORD := $(shell echo $$POSTGRES_P)

postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -d postgres:12-alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:POSTGRES_PASSWORD@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:POSTGRES_PASSWORD@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: createdb, postgres, dropdb, migrateup, migratedown test server
