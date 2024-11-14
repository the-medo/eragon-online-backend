-include app.env
-include app.env.local
export

DB_URL=postgresql://root:secret@localhost:5432/talebound?sslmode=disable
PROTO_SERVICES := $(shell find proto/services -mindepth 1 -maxdepth 1 -type d | sort | awk -F/ '{print "proto/services/" $$NF "/*.proto"}' | tr '\n' ' ')

network-create:
	-docker network inspect talebound-backend >nul 2>&1 || docker network create talebound-backend

rm-postgres:
	docker stop postgres15
	docker rm postgres15

rm-redis:
	docker stop redis
	docker rm redis

postgres:
	docker run --name postgres15 --network talebound-backend -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root talebound

wait-for-createdb:
	timeout 4

dropdb:
	docker exec -it postgres15 dropdb talebound

migrateup:
	go run github.com/the-medo/golang-migrate-objects -mpath=$(MIGRATION_URL) -obj_path=$(MIGRATION_OBJECTS_URL) -db_source=$(DB_SOURCE) -co_filename=$(MIGRATION_CREATE_OBJECTS_FILENAME) -do_filename=$(MIGRATION_DROP_OBJECTS_FILENAME) -up

migrateup1:
	go run github.com/the-medo/golang-migrate-objects -mpath=$(MIGRATION_URL) -obj_path=$(MIGRATION_OBJECTS_URL) -db_source=$(DB_SOURCE) -co_filename=$(MIGRATION_CREATE_OBJECTS_FILENAME) -do_filename=$(MIGRATION_DROP_OBJECTS_FILENAME) -up -step=1

migratedown:
	go run github.com/the-medo/golang-migrate-objects -mpath=$(MIGRATION_URL) -obj_path=$(MIGRATION_OBJECTS_URL) -db_source=$(DB_SOURCE) -co_filename=$(MIGRATION_CREATE_OBJECTS_FILENAME) -do_filename=$(MIGRATION_DROP_OBJECTS_FILENAME) -down

migratedown1:
	go run github.com/the-medo/golang-migrate-objects -mpath=$(MIGRATION_URL) -obj_path=$(MIGRATION_OBJECTS_URL) -db_source=$(DB_SOURCE) -co_filename=$(MIGRATION_CREATE_OBJECTS_FILENAME) -do_filename=$(MIGRATION_DROP_OBJECTS_FILENAME) -down -step=1

sqlc-generate:
	go run github.com/the-medo/golang-migrate-objects -mpath=$(MIGRATION_URL) -obj_path=$(MIGRATION_OBJECTS_URL) -db_source=$(DB_SOURCE) -co_filename=$(MIGRATION_CREATE_OBJECTS_FILENAME) -do_filename=$(MIGRATION_DROP_OBJECTS_FILENAME) -sumfile
	docker run --rm -v "$(CURDIR):/src" -w /src kjconroy/sqlc generate

test:
	go test -v -cover -short ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/the-medo/talebound-backend/db/sqlc Store
	mockgen -package mockwk -destination worker/mock/distributor.go github.com/the-medo/talebound-backend/worker TaskDistributor

db_docs:
	dbdocs password --set secret --project talebound

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

proto_delete_win:
	del /Q pb\**\**\*.pb.go
	del /Q pb\**\**\*.pb.gw.go
	del /Q doc\swagger\*.swagger.json

proto_delete_linux:
	rm -rf pb/**/*.pb.go
	rm -f pb/**/**/*.pb.go
	rm -f pb/**/**/*.pb.gw.go
	rm -f doc/swagger/*.swagger.json

#proto_without_clean:
#	protoc --proto_path=proto --go_out=pb --go-grpc_out=pb --grpc-gateway_out=pb --go_opt=module=pb --go-grpc_opt=module=pb --openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=talebound $(PROTO_SERVICES);
#	statik -src=./doc/swagger -dest=./doc

proto_without_clean:
	protoc --proto_path=proto --go_out=paths=import:pb --go-grpc_out=paths=import:pb --grpc-gateway_out=pb --openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=talebound $(PROTO_SERVICES);
	statik -src=./doc/swagger -dest=./doc
	mv pb/github.com/the-medo/talebound-backend/pb/* pb/
	rm -rf pb/github.com

proto_without_clean2:
	protoc --proto_path=proto --go_out=pb --go-grpc_out=pb --grpc-gateway_out=pb --openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=talebound $(PROTO_SERVICES);
	statik -src=./doc/swagger -dest=./doc

proto_win: proto_delete_win	proto_without_clean

proto: proto_delete_linux proto_without_clean

evans:
	evans --host localhost --port 9090 -r repl

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

prepare: network-create postgres redis wait-for-createdb createdb wait-for-createdb migrateup server

.PHONY: network-create rm-postgres createdb dropdb postgres migrateup migratedown sqlc-generate test mock migrateup1 migratedown1 db_docs db_schema proto_win proto_linux evans proto_without_clean redis new_migration