#!/bin/bash

case "$1" in
  "up")
    docker-compose up -d
    ;;
  "down")
    docker-compose down
    ;;
  "createdb")
    docker-compose exec postgres createdb --username=postgres --owner=postgres postgres
    ;;
  "dropdb")
    docker-compose exec postgres dropdb postgres
    ;;
  "postgres")
    docker-compose exec postgres psql -U postgres
    ;;
  "migrateup")
    migrate -path server/db/migrations -database "postgresql://postgres:password@localhost:5433/postgres?sslmode=disable" -verbose up
    ;;
  "migratedown")
    migrate -path server/db/migrations -database "postgresql://postgres:password@localhost:5433/postgres?sslmode=disable" -verbose down
    ;;
  *)
    echo "Usage: ./db.sh [up|down|createdb|dropdb|postgres|migrateup|migratedown]"
    exit 1
    ;;
esac
