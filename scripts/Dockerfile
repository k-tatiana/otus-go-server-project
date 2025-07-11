FROM golang:1.24

RUN apt-get update && apt-get install -y postgresql-client

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

ENV PATH="/go/bin:${PATH}"