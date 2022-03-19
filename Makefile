postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123456 -d postgres:alpine
createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres dropdb simple_bank
migrateup:
	migrate -path db/migrations -database "postgres://root:123456@localhost:5432/simple_bank?sslmode=disable" up
migratedown:
	migrate -path db/migrations -database "postgres://root:123456@localhost:5432/simple_bank?sslmode=disable" down
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run .
mock:
	mockgen --destination db/mock/store.go --package mockdb simplebank-app/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock