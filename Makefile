DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable

postgres:
	docker run --name postgres12 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -p 5432:5432 -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedownlast:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateuplast:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1


sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/kiyuu10/simplebank/db/sqlc Store

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto

.PHONY: postgres createdb dropdb migrateup migratedown migratedownlast migrateuplast sqlc test server mock proto