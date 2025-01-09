postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:latest

createdb:
	docker exec -it postgres createdb --username=root --owner=root finance_bot

dropdb:
	docker exec -it postgres dropdb finance_bot

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/finance_bot?sslmode=disable" -verbose up

migratedown:
	migrate	-path db/migration -database "postgresql://root:secret@localhost:5432/finance_bot?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc
