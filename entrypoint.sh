#!/bin/sh

# Exit immediately if any command fails
set -e

# Check if required environment variables are set
if [ -z "$DB_CONNECTION" ] || [ -z "$DB_USERNAME" ] || [ -z "$DB_PASSWORD" ] || [ -z "$DB_HOST" ] || [ -z "$DB_PORT" ] || [ -z "$DB_DATABASE" ]; then
  echo "Error: One or more required environment variables are not set."
  exit 1
fi

# Construct the database URL using environment variables
DB_URL="$DB_CONNECTION://$DB_USERNAME:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_DATABASE?sslmode=disable"

# Wait for the database to be ready
sleep 5

# Apply migrations using the constructed database URL
migrate -path=/app/migrations -database="$DB_URL" -verbose up

# Start go-pos service
exec ./gopos
