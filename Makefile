

build-backend:
	docker build -t go-e2ee-app ./backend
run-backend:
	docker run -p 8080:8080 --env-file ./backend/.env go-e2ee-app
