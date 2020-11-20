.PHONY: test env-setup

env-setup:
	docker-compose -f docker-compose.dev.yml up -d

start-dev: env-setup
	forego start

test: env-setup
	go test -v ./...