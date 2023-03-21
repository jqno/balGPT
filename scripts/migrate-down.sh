#!/usr/bin/env sh

migrate -path db/migrations -database "postgres://$DB_USER:$DB_PASS@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable" down
