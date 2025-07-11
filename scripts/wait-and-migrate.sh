#!/bin/sh
set -e

# Wait for the database to be ready
until PGPASSWORD=postgres pg_isready -h go-server-db -p 5432 -U postgres; do
  echo "Waiting for postgres..."
  sleep 2
done

# Run goose migrations
GOOSE_DRIVER=postgres GOOSE_DBSTRING="user=postgres dbname=postgres host=db password=postgres sslmode=disable" goose -dir /migrations up