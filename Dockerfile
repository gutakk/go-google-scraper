FROM golang:alpine

ENV GIN_MODE=release \
    APP_PORT=8080

WORKDIR /app

RUN apk add --update nodejs npm

# Install JS dependencies
COPY package.json package-lock.json ./
RUN npm install

# Install Go dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN npm run build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

WORKDIR /app/workers

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

EXPOSE 8080

ENTRYPOINT ["/app/go-google-scraper"]
