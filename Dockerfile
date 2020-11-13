FROM golang:1.15-buster as builder

ENV GO111MOD=on

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

# Release stage
FROM alpine as release

RUN apk update \
    && apk upgrade \
    && apk add --no-cache \
    ca-certificates \
    && update-ca-certificates 2>/dev/null || true

ENV APP_ENV=release \
    GIN_MODE=release

COPY --from=builder /app/go-google-scraper /app/

EXPOSE 8080

ENTRYPOINT ["/app/go-google-scraper"]
