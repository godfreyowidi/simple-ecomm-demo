include .env
export

.PHONY: up down rebuild migrate migrate-test test init-db drop-test-db reset-test-db reset

# Docker Compose Controls
up:
	docker compose up -d

down:
	docker compose down -v

rebuild:
	docker compose down -v
	docker compose build
	docker compose up -d

# Run dev migrations
migrate:
	docker compose run --rm migrate

# Run test migrations
migrate-test:
	docker compose run --rm migrate-test

# Reset everything (DB, app, migrations)
reset:
	docker compose down -v
	docker compose build
	docker compose up -d

# Drop test DB manually
drop-test-db:
	docker exec -i savanna-db psql -U postgres -c "DROP DATABASE IF EXISTS savanna_test;"

# Reset test DB only
reset-test-db: drop-test-db up migrate-test

# Run tests
test:
	@echo "Running tests on: ${TEST_DATABASE_URL}"
	@TEST_DATABASE_URL=${TEST_DATABASE_URL} go test ./internal/repo/... -v
