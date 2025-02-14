.PHONY: up down env

down:
	docker compose down

up:
	docker compose up -d

# ex:
# 	export $(grep -v '^#' .env | xargs)

db_login:
	psql ${DB_URL}

db_create_migrate:
	migrate create -ext sql -dir migrations $(name)
	
db_migrate:
	migrate -database ${DB_URL} -path migrations up 