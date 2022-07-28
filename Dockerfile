# syntax=docker/dockerfile:1

FROM golang:1.18.4-alpine3.16

LABEL maintainer = "Petkanych Bogdan <iosifsuzuki@gmail.com>"

ENV GO111MODULE=on

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . ./
RUN go build -o main ./cmd
EXPOSE 8080
CMD ["/app/main"]


