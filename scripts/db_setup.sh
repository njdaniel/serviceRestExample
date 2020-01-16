#!/bin/bash
echo "Setting up database.."
export PGPASSWORD='password'
docker run --name dnd_db -e POSTGRES_PASSWORD=password -d -p 5432:5432 postgres:latest
psql -h localhost -p 5432 -U postgres -a -f ./schema/setup.sql
psql -h localhost -p 5432 -U postgres -a -f ./schema/seed.sql
echo "Database completed"