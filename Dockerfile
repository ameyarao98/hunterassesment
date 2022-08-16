FROM golang:1.18-buster

WORKDIR /app

COPY . ./

RUN go mod download

ENTRYPOINT [ "go", "run", "main.go" ]