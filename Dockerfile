FROM golang:alpine

ENV GIN_MODE=release \
    APP_ENV=release

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY .env ./
COPY .env.release ./

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

EXPOSE 8080

ENTRYPOINT ["/app/go-google-scraper"]