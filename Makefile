.PHONY: up down build help


build: ## Build the Docker image
	docker-compose build --no-cache


up: ## Build and start the containers
	docker-compose up --build -d


down: ## Stop and remove the containers
	docker-compose down

help: ## Show help for Makefile commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
