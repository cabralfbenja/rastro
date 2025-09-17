# Define the default target
.PHONY: default
default: help

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  seeder FILENAME     Process the specified seeder file name"

# application repository and binary file name
NAME=Rastro

# application repository path
REPOSITORY=github.com/cabralfbenja/rastro

run-dev:
	echo "Starting application in dev mode"
	go run ./cmd/api

migrate:
	echo "Running migrations up"
	go run ./internal/database/migrate_up.go

install:
	go mod download

#E.G make seeder FILENAME=category
.PHONY: seeder
seeder:
ifdef FILENAME
	echo "Seedinng : $(FILENAME).go"
	go run "./internal/database/seeders/$(FILENAME)_seeder.go"
else
	echo "Error: FILENAME is not specified. Please provide the filename using 'make seeder FILENAME=<filename>'"
	exit 1
endif

migrate_fresh:
	echo "Running fresh migrations"
	go run ./internal/database/migrate_fresh.go

migrate_down:
	go run ./migrations/drop_migrations.go

test:
	 gotest ./tests/... -v