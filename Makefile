PROJECT_NAME = spy-cats

docker-compose = docker compose -f .docker/docker-compose.yaml

run:  ## run app
	go run ./cmd/app
env: ## copy env
	cp ./config/.env.example ./config/.env
setup: ## first launch
	make env
	make start
start:
	$(docker-compose) up -d --remove-orphans
stop:
	$(docker-compose) down
restart:
	make stop
	make start
