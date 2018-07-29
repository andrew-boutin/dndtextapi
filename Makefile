# Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

# So we can see what commands get ran from the command line output.
SHELL = sh -xv

.PHONY: all up down reapp test fmt lint vet updatedeps installdeps inttests fetchdeps

all: down test vet lint fmt up

# Start up the app related containers.
up:
	@docker-compose up --build -d app

# Tear down all containers.
down:
	@docker-compose down

# Rebuild and start just the go app container.
# Useful when making code changes but you want to keep the current db data.
reapp:
	@docker-compose up --build -d app

# Run all unit tests in the project.
test:
	@./unit-test.sh

# Lint for style suggestions.
lint:
	@$(GOLINT_PATH) `go list ./... | grep -v /vendor/`

# Vet the code to look for potential issues.
vet:
	@go vet -shadow -shadowstrict ./...

# Format the files using Go standard best practices.
fmt:
	@go fmt ./...

# Use govendor to update the dependencies in vendor/vendor.json
updatedeps:
	@$(GOVENDOR_PATH) update +v

# Use govendor to install the dependencies listed in vendor/vendor.json
installdeps:
	@$(GOVENDOR_PATH) install

# Use govendor to fetch the dependencies listed in vendor/vendor.json.
# Requred to run this before most other Go commands that operate outside of the Docker containers.
fetchdeps:
	@$(GOVENDOR_PATH) fetch +out

# Run the integration tests.
# Requires that the app and database are already running.
# You'll probably want the DNDTEXTAPI_ENV variable set to `int` as well.
inttests:
	@docker-compose rm --force inttest
	@docker-compose up --build --exit-code-from inttest inttest