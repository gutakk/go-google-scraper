.PHONY: test env-setup

build-dependencies:
	go get github.com/ddollar/forego@v0.16.1
	go get github.com/cosmtrek/air@v1.15.1

env-setup:
	docker-compose -f docker-compose.dev.yml up -d

start-dev: env-setup
	forego start

test: env-setup
	go test -v -p 1 ./...
