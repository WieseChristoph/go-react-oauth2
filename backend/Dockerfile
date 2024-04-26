FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o main ./cmd/app/main.go

FROM alpine:3.14

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

CMD ["./main"]
