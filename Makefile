IMAGE_NAME=friends-management-app
CONTAINER_NAME=friends-management-app-container
POSTGRES_CONTAINER=my-postgres
DOCKER_NETWORK=my_network

DB_HOST=host.docker.internal
DB_USER=postgres
DB_PASSWORD=admin
DB_NAME=friends_management
DB_PORT=5432
PORT=:8080

build:
	docker build -t $(IMAGE_NAME) .

build-docker-network:
	docker network create ${DOCKER_NETWORK} || true

run-postgres: build-docker-network
	docker run --name ${POSTGRES_CONTAINER} \
		--network ${DOCKER_NETWORK} \
		-e POSTGRES_USER=${DB_USER} \
		-e POSTGRES_PASSWORD=${DB_PASSWORD} \
		-e POSTGRES_DB=${DB_NAME} \
		-p ${DB_PORT}:5432 \
		-d postgres:15


run:
	docker run -d \
		--name ${CONTAINER_NAME} \
		--network ${DOCKER_NETWORK} \
		-p ${PORT}:8080 \
		-e DB_HOST=${DB_HOST} \
		-e DB_PORT=${DB_PORT} \
		-e DB_USER=${DB_USER} \
		-e DB_PASSWORD=${DB_PASSWORD} \
		-e DB_NAME=${DB_NAME} \
		${IMAGE_NAME}

start-app: build build-docker-network run-postgres run
	@echo "Application is now running!"
	@echo "Run 'make logs' to see the application logs"

clean:
	@echo "Cleaning up..."
	docker rm -f $(CONTAINER_NAME) || true
	docker rm -f $(POSTGRES_CONTAINER) || true
	docker rmi -f $(IMAGE_NAME) || true
	docker rmi -f $(IMAGE_NAME) || true
	docker network rm ${DOCKER_NETWORK} || true

logs:
	docker logs -f $(CONTAINER_NAME)

migrate:
	go run ./cmd/migrate

test:
	go test -v ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out
	@echo "Test coverage report generated: coverage.html"
	@echo "Run 'open coverage.html' to view the report in your browser"

