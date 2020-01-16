# serviceRestExample

Example of a Rest API

requires:

 - docker

## Quickrun

    docker run --name dnd_db -e POSTGRES_PASSWORD=password -d -p 5432:5432 postgres:latest
    go build
    ./serviceRestExample
    
## Connect to database
    
    export PGPASSWORD='password'
    psql -h localhost -p 5432 -U postgres -a -f ./schema/setup.sql
    psql -h localhost -p 5432 -U postgres -a -f ./schema/seed.sql
    
## Cleanup
    go clean
    