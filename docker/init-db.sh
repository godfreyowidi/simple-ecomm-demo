#!/bin/bash
set -e

echo "⏳ Waiting for PostgreSQL to become ready..."
until pg_isready -U "$POSTGRES_USER"; do
  sleep 1
done
echo "✅ PostgreSQL is ready."

# Create savanna_test if it doesn't exist
EXISTS=$(psql -U "$POSTGRES_USER" -tc "SELECT 1 FROM pg_database WHERE datname = 'savanna_test';" | tr -d '[:space:]')

if [ "$EXISTS" != "1" ]; then
  echo "🔧 Creating 'savanna_test' database..."
  createdb -U "$POSTGRES_USER" savanna_test
  echo "✅ 'savanna_test' created."
else
  echo "ℹ️  'savanna_test' already exists."
fi
