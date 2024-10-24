include app.env

postgres:
	docker run --name postgres_container --network simplebank-network -p 5432:5432 -e POSTGRES_USER=$(POSTGRES_USER) -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) -d postgres

createdb: 
	docker exec -it postgres_container createdb --username=$(POSTGRES_USER) --owner=$(POSTGRES_USER) simple_bank

dropdb:
	docker exec -it postgres_container dropdb simple_bank

migrateup:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose down 1

db_docs: 
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --potsgres -o doc/schema.sql doc/db.dbml

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server: 
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go  github.com/jnka9755/go-07SIMPLE-BANK/db/sqlc Store

proto:
	rm -f pb/*.pb.go
	rm -f doc/swagger/*.swagger.js
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
    proto/*.proto
	statik -src=./doc/swagger -dest=./doc

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 db_docs db_schema sqlc test server mock proto
