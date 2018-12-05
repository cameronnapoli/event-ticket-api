FROM golang:1.11

ADD ./go/src/festival_ticket_api/api

RUN go get github.com/gorilla/mux
RUN go get github.com/go-redis/redis

ENTRYPOINT go run main.go
