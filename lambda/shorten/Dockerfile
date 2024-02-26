FROM golang:latest
WORKDIR /usr/src/app

RUN mkdir -p /asset

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /asset ./cmd/...
