# Go Google Scraper
#### Gathering Google search information with your own keywords

### [Project Backlog](https://github.com/gutakk/go-google-scraper/projects/1)

## Prerequisite
* [Go - 1.15](https://golang.org/doc/go1.15)
* [Docker](https://docs.docker.com/get-docker/)
* [Docker Compose](https://docs.docker.com/compose/install/)

## Usage
#### Setup and boot the Docker containers
```sh
make env-setup
```

#### Run the Go application for development
```go
go run main.go
```
To visit app locally: `localhost:8080`

#### Run tests
```go
make test
```

## About
This project is created to complete **Web Certification Path** using **Golang** at [Nimble](https://nimblehq.co)
