
run:
	air

.PHONY: docker
docker:
	podman-compose -f docker/docker-compose-dev.yml up -d

compile_sql:
	sqlc generate
