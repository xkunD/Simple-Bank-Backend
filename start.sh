#!/bin/sh

# Make sure the script exit inmediatly if an error is detected.
set -e 

# Run migrations up
echo "run db migration"
/usr/bin/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

# Take all parameters pass to the script and run them.
echo "start the app"
exec "$@"