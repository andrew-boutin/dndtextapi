# Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

.PHONY: all up down reapp test fmt lint vet updatedeps installdeps inttests

all: down test vet lint fmt up

# Start up the app related containers
up:
	@docker-compose up --build -d app

# Tear down all containers
down:
	@docker-compose down

# Rebuild and start just the go app container. Useful when making code changes but you want to keep the current db data.
reapp:
	@docker-compose up --build -d app

# Run all unit tests in the project
test:
	@./unit-test.sh

# Lint for style suggestions
lint:
	@$(GOLINT_PATH) `go list ./... | grep -v /vendor/`

# Vet the code to look for potential issues
vet:
	@go vet -shadow -shadowstrict ./...

# Format the files using Go standard best practices
fmt:
	@go fmt ./...

# Use govendor to update the dependencies in vendor/vendor.json
updatedeps:
	@$(GOVENDOR_PATH) update +v

# Use govendor to install the dependencies listedin vendor/vendor.json
installdeps:
	@$(GOVENDOR_PATH) install

# Run the integration tests
inttests:
	@docker-compose rm --force inttest
	@docker-compose up --build --exit-code-from inttest inttest