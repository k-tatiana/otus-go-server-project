#!/bin/sh
set -e

# setup master db
# Wait for the database to be ready
until PGPASSWORD=postgres pg_isready -h pgmaster -p 5432 -U postgres; do
  echo "Waiting for postgres..."
  sleep 2
done

# Run goose migrations
GOOSE_DRIVER=postgres GOOSE_DBSTRING="user=postgres dbname=postgres host=pgmaster port=5432 password=postgres sslmode=disable" goose -dir /migrations up

# # setup slave 1 db
# # Wait for the database to be ready
# until PGPASSWORD=postgres pg_isready -h pgslave -p 5432 -U postgres; do
#   echo "Waiting for postgres..."
#   sleep 2
# done

# # Run goose migrations
# GOOSE_DRIVER=postgres GOOSE_DBSTRING="user=postgres dbname=postgres host=pgslave port=5432 password=postgres sslmode=disable" goose -dir /migrations up

# # setup slave 2 db
# # Wait for the database to be ready
# until PGPASSWORD=postgres pg_isready -h pgasyncslave -p 5432 -U postgres; do
#   echo "Waiting for postgres..."
#   sleep 2
# done

# # Run goose migrations
# GOOSE_DRIVER=postgres GOOSE_DBSTRING="user=postgres dbname=postgres host=pgasyncslave port=5432 password=postgres sslmode=disable" goose -dir /migrations up