#!/bin/bash
set -e

# we wait for Postgres to start
until pg_isready -U postgres; do
  echo "Waiting for postgres..."
  sleep 2
done

# we create savanna_test if not exists
psql -U postgres <<'EOF'
DO $$
BEGIN
   IF NOT EXISTS (
      SELECT FROM pg_database WHERE datname = 'savanna_test'
   ) THEN
      CREATE DATABASE savanna_test;
   END IF;
END
$$;
EOF

echo "✅ savanna_test database is ready."
