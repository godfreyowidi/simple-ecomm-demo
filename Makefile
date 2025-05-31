.PHONY: up down rebuild migrate test init-db drop-test-db reset-test-db test-migrate

# Docker Commands
up:
	docker compose up -d

down:
	docker compose down -v

rebuild:
	docker compose down -v
	docker compose build
	docker compose up -d

# Run migrations using migrate container (on savanna_test DB)
migrate:
	docker compose run --rm migrate

# Init script is already baked into docker-entrypoint; this is optional now.
init-db:
	@echo "Test DB 'savanna_test' will be created automatically via docker-compose entrypoint."

# Drop test DB manually
drop-test-db:
	docker exec -i savanna-db psql -U postgres -c "DROP DATABASE IF EXISTS savanna_test;"

# Reset test DB (drop + recreate)
reset-test-db: drop-test-db up migrate

# Run unit tests against savanna_test
test:
	@echo "Running tests on: ${TEST_DATABASE_URL}"
	@TEST_DATABASE_URL=${TEST_DATABASE_URL} go test ./internal/repo/... -v

# Manually run migrations on savanna_test (host-based)
test-migrate:
	docker run --rm \
	  -v ${PWD}/migrations:/migrations \
	  migrate/migrate \
	  -path=/migrations \
	  -database "postgres://postgres:securepass@localhost:5434/savanna_test?sslmode=disable" \
	  up
