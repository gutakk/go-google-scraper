.PHONY: test env-setup

build-dependencies:
	go get github.com/mattn/goreman@v0.3.7
	go get github.com/cosmtrek/air@v1.15.1

env-setup:
	docker-compose -f docker-compose.dev.yml up -d

start-dev: env-setup
	goreman start

test:
	docker-compose -f docker-compose.test.yml up -d
	go test -v -p 1 -count=1 ./...
	docker-compose -f docker-compose.test.yml down

