.PHONY: test env-setup

build-dependencies:
	go get github.com/ddollar/forego
	go get github.com/cosmtrek/air

env-setup:
	docker-compose -f docker-compose.dev.yml up -d

start-dev: env-setup
	forego start

test: env-setup
	go test -v -p 1 -count=1 ./...
