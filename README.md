# Go Google Scraper
#### Gathering Google search information with your own keywords

### [Project Backlog](https://github.com/gutakk/go-google-scraper/projects/1)

## Prerequisite
* [Go - 1.15](https://golang.org/doc/go1.15)
* [Docker](https://docs.docker.com/get-docker/)
* [Docker Compose](https://docs.docker.com/compose/install/)
* [NodeJS](https://nodejs.org/en/download/package-manager/)

## Create necessary dot env files
- Create `.env` file
- Create env dependent files (depend on your env)
  - `.env.release`
  - `.env.debug` (for development)
  - `.env.test`
- Add values from `.env.example` (for env dependent variables eg. `DB_NAME` must add to dependent env  files)
## Usage
#### Run the Go application for development
[**`.env` and `.env.debug` files are required**](#create-necessary-dot-env-files)

This command will hot reload both assets build and app build
```sh
make start-dev
```
To visit app locally: `localhost:8080`

#### Build assets
```sh
npm run build
```

#### Run tests
```sh
make test
```

## About
This project is created to complete **Web Certification Path** using **Golang** at [Nimble](https://nimblehq.co)
