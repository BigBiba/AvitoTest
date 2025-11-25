build:
	docker build -t pr .

up:
	docker compose up -d

down:
	docker compose down

migrate_up:
	echo "create migration"

migrate_down:
	echo "delete migrations"
