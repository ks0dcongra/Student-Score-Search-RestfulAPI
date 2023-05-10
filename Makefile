postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:14-alpine

migrateup:
	migrate -path database/migration -database "postgresql://postgres:postgres@localhost:5432/gexample?sslmode=disable" -verbose up 1

migratedown:
	migrate -path database/migration -database "postgres://postgres:@127.0.0.1:5432/example?sslmode=disable" -verbose down 1
    
hello:
	echo "Hello"

.PHONY:postgres migrateup migratedown hello