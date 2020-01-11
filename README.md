# serviceRestExample

Example of a Rest API

requires:

 - docker

## Quickrun

    docker run --name dnd_db -e POSTGRES_PASSWORD=password -d -p 5432:5432 postgres:latest
    go build
    ./serviceRestExample
    
## Cleanup
    go clean
    