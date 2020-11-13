FROM golang:alpine

ENV GIN_MODE=release \
    APP_ENV=release \
    PORT=8080

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

EXPOSE $PORT

ENTRYPOINT ["/app/go-google-scraper"]
