.PHONY: help

help:
	@echo "\033[32mAvailable targets:\033[0m"
	@grep "^[a-zA-Z\-]*:" Makefile | grep -v "grep" | sed -e 's/^/make /' | sed -e 's/://'

cp-env:
	cp .env.example .env
up:
	docker compose up -d
db:
	docker compose exec -it db mysql -u myuser -pmypassword -D mydatabase
logs:
	docker compose logs -f
orm:
	sqlboiler mysql -c config/database.toml -o model -p model --no-tests --wipe
migrate:
	migrate create -ext sql -dir migrations -seq add_timestamps_to_users_and_todos
mup:
	migrate -path migrations -database 'mysql://myuser:mypassword@tcp(127.0.0.1:3306)/mydatabase' -verbose up
mdown:
	migrate -path migrations -database 'mysql://myuser:mypassword@tcp(127.0.0.1:3306)/mydatabase' -verbose down