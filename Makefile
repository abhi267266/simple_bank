postgres:
	docker run --name postgres12 -p 5050:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5050/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5050/simple_bank?sslmode=disable" -verbose down

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

.PHONY:postgres createdb dropdb migrateup migratedown sqlc test