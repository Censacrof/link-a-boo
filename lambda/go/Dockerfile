FROM golang:latest
WORKDIR /usr/src/app

RUN mkdir -p /asset

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go test ./cmd/... ./pkg/...

RUN mkdir -p bin
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin ./cmd/...
ARG CMD_NAME
RUN mv ./bin/${CMD_NAME} /asset/bootstrap
