# Author: Rodrigo Salazar

# Help
help:
	@echo "\033[1;34mAvailable commands:\033[0m"
	@echo "\033[1;32m  backend:\033[0m Run the backend service that will be deployed in ECS"
	@echo "\033[1;32m  mock:\033[0m Run a mock service that simulates the WS API Gateway + Connect/Disconnect Lambdas"
	@echo "\033[1;32m  client:\033[0m Run the client service"

# Shortcuts
backend: build-backend stop-backend run-backend
mock: build-mock stop-mock run-mock
client: run-client

# Commands
create-network:
	@docker network inspect e2ee-network > /dev/null 2>&1 || docker network create e2ee-network
build-backend:
	docker build -t go-e2ee-app ./backend
run-backend: create-network
	docker run -p 8080:8080 --network e2ee-network --env-file ./backend/.env --name local-e2ee-app go-e2ee-app
stop-backend:
	@docker ps -q -f name=local-e2ee-app && docker stop local-e2ee-app || echo "Backend container not running"
	@docker ps -aq -f name=local-e2ee-app && docker rm local-e2ee-app || echo "Backend container not found"
shell-backend:
	docker run -it --network e2ee-network --env-file ./backend/.env --name local-e2ee-app go-e2ee-app /bin/sh
build-mock:
	docker build -t go-e2ee-mock ./mock
run-mock: create-network
	docker run -p 8082:8082 --network e2ee-network --name mock-gateway go-e2ee-mock
stop-mock:
	@docker ps -q -f name=mock-gateway && docker stop mock-gateway || echo "Mock container not running"
	@docker ps -aq -f name=mock-gateway && docker rm mock-gateway || echo "Mock container not found"
run-client:
	cd ./client/e2ee_client && go run .´

push-backend:
	docker tag go-e2ee-app tux550/go-e2ee-app
	docker push tux550/go-e2ee-app
