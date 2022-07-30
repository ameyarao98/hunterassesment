FROM golang:1.18-buster AS build

WORKDIR /app

COPY . ./

RUN go mod download

ENTRYPOINT [ "go", "run", "main.go" ]