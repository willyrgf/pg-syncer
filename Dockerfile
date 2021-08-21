FROM golang:1.17.0-alpine3.14

WORKDIR /app

COPY . .
RUN go build .
