build-backend:
	docker build -t go-e2ee-app ./backend
run-backend: create-network
	docker run -p 8080:8080 --network e2ee-network --env-file ./backend/.env --name local-e2ee-app go-e2ee-app
stop-backend:
	docker stop local-e2ee-app
	docker rm local-e2ee-app
build-mock:
	docker build -t go-e2ee-mock ./mock
run-mock: create-network
	docker run -p 8082:8082 --network e2ee-network --name mock-gateway go-e2ee-mock
stop-mock:
	docker stop mock-gateway
	docker rm mock-gateway
build-all: build-backend build-mock

create-network:
	@docker network inspect e2ee-network > /dev/null 2>&1 || docker network create e2ee-network

run-client:
	cd ./client/e2ee_client && go run .