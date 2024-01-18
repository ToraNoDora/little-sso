# Variables
APP_PATH := ./sso
APP_PORT := 44044

MAIN_PACKAGE_PATH := ./cmd/sso/main.go

SECRET_PATH := ./sso/.env
BINARY_NAME := little_sso
OUTPUT_PATH := ./tmp/bin

TEST_PATH := ./tests

STORE_DB := sso_dev
STORE_HOST := 127.0.0.1
STORE_PORT := 5432
STORE_USER := postgres
STORE_PSW := chinchi

REDIS_PORT := 6379

MIGRATIONS_PATH := ./migrations/schema

DOCKER_FILE := ./Dockerfile
DOCKER_IMAGE_NAME := little_sso
DOCKER_REGISTRY := docker-registry-test
DOCKER_REPO := little-sso

REDIS_NAME := redis_sso_dev
PSQL_NAME := psql_sso_dev


# base install
.PHONY: install
install: install run.psql migrate.up run.redis
	cd ${APP_PATH} && \
		go mod download


# Quality control
.PHONY: tidy
tidy:
	cd ${APP_PATH} && \
		go fmt ./... && \
		go mod tidy -v


# Tests
.PHONY: test
test:
	cd ${APP_PATH} && \
		go test ${TEST_PATH}

test.detail:
	cd ${APP_PATH} && \
		go test -v ${TEST_PATH}

test.cover:
	cd ${APP_PATH} && \
		go tool cover -func=coverage.out


# Build the application
.PHONY: build
build:
	cd ${APP_PATH} && \
		go build -o=${OUTPUT_PATH}/${BINARY_NAME} ${MAIN_PACKAGE_PATH}


# Run the application
.PHONY: run
run: build
	cd ${APP_PATH} && \
		${OUTPUT_PATH}/${BINARY_NAME}

# Reloading on file changes
.PHONY: run.live
run.live:
	cd ${APP_PATH} && \
		go run ${MAIN_PACKAGE_PATH}


# migrations
install.migrate:
	go install github.com/golang-migrate/migrate/cli@latest

migrate.create:
	migrate create -ext sql -dir ${MIGRATIONS_PATH} -seq $(name)

migrate.up:
	migrate -path ${MIGRATIONS_PATH} -database 'postgres://${STORE_USER}:${STORE_PSW}@${STORE_HOST}:${STORE_PORT}/${STORE_DB}?sslmode=disable' up

migrate.down:
	migrate -path ${MIGRATIONS_PATH} -database 'postgres://${STORE_USER}:${STORE_PSW}@${STORE_HOST}:${STORE_PORT}/${STORE_DB}?sslmode=disable' down


# Docker build/run/push
docker.build:
	docker build --no-cache -f ${DOCKER_FILE} -t ${DOCKER_REGISTRY}/${DOCKER_REPO}:$(tag) .

docker.run:
	docker run \
		--env-file ${SECRET_PATH} \
		-p ${APP_PORT}:${APP_PORT} \
		-d --name ${DOCKER_IMAGE_NAME} \
		--rm -ti ${DOCKER_REGISTRY}/${DOCKER_REPO}:$(tag)

docker.push:
	docker push ${DOCKER_REGISTRY}/${DOCKER_REPO}:$(tag)


# additionally
.PHONY: run.psql
run.psql:
	docker run \
		--name ${PSQL_NAME} \
		-p ${STORE_PORT}:${STORE_PORT} \
		-e POSTGRES_PASSWORD=${STORE_PSW} \
		-e POSTGRES_DB=${STORE_DB} -d -v "$(pwd)":/docker-entrypoint-initdb.d \
		-d postgres

run.redis:
	docker run \
		-p ${REDIS_PORT}:${REDIS_PORT} \
		--name ${REDIS_NAME} \
		-d redis && \
		docker exec -it ${REDIS_NAME} redis-server

