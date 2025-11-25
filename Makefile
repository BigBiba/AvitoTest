build:
	docker build -t pr .

up:
	docker compose up -d

down:
	docker compose down

migrate_up:
	echo "make migrations"

migrate_down:
	echo "delete migrations"
