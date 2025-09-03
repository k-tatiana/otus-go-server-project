#!/bin/sh
set -e

setup citus master
Wait for the database to be ready
until PGPASSWORD= pg_isready -h citus_master -p 5432 -U postgres; do
  echo "Waiting for postgres..."
  sleep 2
done

# Run goose migrations
GOOSE_DRIVER=postgres GOOSE_DBSTRING="user=postgres dbname=postgres host=citus_master port=5432 password= sslmode=disable" goose -dir /migrations up
