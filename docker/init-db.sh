#!/bin/bash
set -e

# unix socket directly - avoided tcp localhost issues
SOCKET_DIR="/var/run/postgresql"

echo "⏳ Waiting for PostgreSQL to become ready..."
until pg_isready -U "$POSTGRES_USER" -h "$SOCKET_DIR"; do
  sleep 1
done
echo "✅ PostgreSQL is ready."

# this is to chheck if 'savanna_test' exists
EXISTS=$(psql -U "$POSTGRES_USER" -h "$SOCKET_DIR" -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname = 'savanna_test';")

if [ "$EXISTS" != "1" ]; then
  echo "🔧 Creating 'savanna_test' database..."
  createdb -U "$POSTGRES_USER" -h "$SOCKET_DIR" savanna_test
  echo "✅ 'savanna_test' created."
else
  echo "ℹ️  'savanna_test' already exists."
fi
