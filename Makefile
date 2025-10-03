GIT_TAG := $(shell git describe --tags --always)
BUILD_TIME := $(shell TZ='America/Los_Angeles' date -u +"%c")


run:
	GIT_TAG="$(GIT_TAG)" BUILD_TIME="$(BUILD_TIME)" air

.PHONY: docker
docker:
	podman-compose -f docker/docker-compose-dev.yml up -d

compile_sql:
	sqlc generate

gen_docs:
	swag init --dir cmd/api --parseInternal --parseDepth 10