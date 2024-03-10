FROM golang:latest
WORKDIR /usr/src/app

RUN mkdir -p /asset

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod/ \
	go mod download && go mod verify

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod/ \
	--mount=type=cache,target=/root/.cache/go-build/ \
	go test -v ./cmd/... ./pkg/...

RUN mkdir -p bin
RUN --mount=type=cache,target=/go/pkg/mod/ \
	--mount=type=cache,target=/root/.cache/go-build/ \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin ./cmd/...

ARG CMD_NAME
RUN mv ./bin/${CMD_NAME} /asset/bootstrap
