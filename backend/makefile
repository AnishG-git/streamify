build:
	docker compose build

run:
	docker compose up

down:
	docker compose down

# sample usage: make create-migration migration-name=add_user_table
create-migration:
	goose -dir ./migrations create $(migration-name) sql

goose-up:
	goose -dir ./migrations postgres "user=postgres password=postgres host=localhost port=5434 dbname=streamify sslmode=disable" up

goose-down:
	goose -dir ./migrations postgres "user=postgres password=postgres host=localhost port=5434 dbname=streamify sslmode=disable" down

# only works if you have psql installed. If so, password is postgres
psql:
	psql -p 5434 -U postgres -d streamify