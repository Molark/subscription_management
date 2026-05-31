FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN swag init -g cmd/main.go

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .

COPY --from=builder /app/docs ./docs

EXPOSE 8080
CMD ["./main"]