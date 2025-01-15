postgres_create:
	docker run --name postgres-container -e POSTGRES_USER=root -e POSTGRES_PASSWORD=admin -p 5432:5432 -d postgres

postgres_start:
	docker start postgres-container

create_db:
	docker exec -it postgres-container createdb --username=root --owner=root go-bank

drop_db:
	docker exec -it postgres-container dropdb go-bank

migrate_up:
	migrate -path ./db/migration -database "postgresql://root:admin@localhost:5432/go-bank?sslmode=disable" -verbose up

migrate_down:
	migrate -path ./db/migration -database "postgresql://root:admin@localhost:5432/go-bank?sslmode=disable" -verbose down

test:
	go test -v -cover ./...

.PHONY: postgres_start postgres_create create_db drop_db migrate_up migrate_down test