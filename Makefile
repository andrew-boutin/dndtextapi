.PHONY: all up down reapp test lint

all: down test lint up

# Start up everything
up:
	@docker-compose up --build -d

# Tear down everything
down:
	@docker-compose down

# Rebuild and start just the go app. Useful when making code changes.
reapp:
	@docker-compose up --build -d app

# Run unit tests
test:
	@./unit-test.sh

# Lint
lint:
	@$(GOLINT_PATH) ./...