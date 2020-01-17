.PHONEY: all build run clean

all: build run

build:
	@echo "Building:"
	@echo "Setting up database"
	docker run --name dnd_db -e POSTGRES_PASSWORD=password -d -p 5432:5432 postgres:latest
	$(bash ./scrips/db_setup.sh)
	env GOOS=linux GOARCH=amd64 go build --tags netgo --ldflags '-extldflags "-lm -lstdc++ -static"'

run:
	@echo "Starting serviceRestExample.."
	./serviceRestExample

clean:
	@echo "Cleaning serviceRestExample.."
	go clean