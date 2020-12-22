.PHONY: test env-setup

build-dependencies:
	go get github.com/mattn/goreman@v0.3.7
	go get github.com/cosmtrek/air@v1.15.1

env-setup:
	docker-compose -f docker-compose.dev.yml up -d

start-dev: env-setup
	goreman start

test: env-setup
	go test -v -p 1 ./...
