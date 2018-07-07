.PHONY: up down reapp

# Start up everything
up:
	@docker-compose up --build -d

# Tear down everything
down:
	@docker-compose down

# Rebuild and start just the go app. Useful when making code changes.
reapp:
	@docker-compose up --build -d app