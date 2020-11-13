env-setup:
	docker-compose -f docker-compose.dev.yml up -d

test:
	go test -v ./...