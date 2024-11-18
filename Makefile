

build-backend:
	docker build -t go-e2ee-app ./backend
run-backend:
  # With name local-e2ee-app
	docker run -p 8080:8080 --env-file ./backend/.env --name local-e2ee-app go-e2ee-app
stop-backend:
	docker stop local-e2ee-app
	docker rm local-e2ee-app