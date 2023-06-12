
##
startdb:
	docker start practiceDB

stopdb:
	docker stop practiceDB

## Database
postgres:
	docker run --name practiceDB -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=admin -d postgres:14-alpine

createdb:
	docker exec -it practiceDB createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it practiceDB dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:admin@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:admin@localhost:5432/simple_bank?sslmode=disable" -verbose down

## sqlc
sqlc:
	sqlc generate

# tests
test:
	go test -v -cover ./...

.PHONY:
	startdb stopdb postgres createdb dropdb migrateup migratedown sqlc test
