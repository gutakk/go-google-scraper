.PHONY: test env-setup

build-dependencies:
	go get github.com/mattn/goreman@v0.3.7
	go get github.com/cosmtrek/air@v1.15.1

env-setup:
	docker-compose -f docker-compose.dev.yml up -d

start-dev: env-setup
	goreman start

test:
	go test -v -p 1 ./...

test-env-setup:
	docker-compose -f docker-compose.test.yml up -d

test-env-destroy:
	docker-compose -f docker-compose.test.yml up -d destroy
