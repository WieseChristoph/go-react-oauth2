FROM golang:1.22-alpine

RUN apk add --no-cache make

RUN go install github.com/cosmtrek/air@latest

CMD ["air"]