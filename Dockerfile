FROM golang:alpine

ENV GIN_MODE=release
ENV APP_ENV=release

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

EXPOSE 8080

ENTRYPOINT ["/app/go-google-scraper"]
