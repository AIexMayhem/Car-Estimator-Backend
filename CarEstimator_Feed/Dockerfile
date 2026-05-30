FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOPATH=/go

WORKDIR /go/src/github.com/icecarti/feed_service

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /go/bin/feed_service ./main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /go/bin/feed_service /app/feed_service
COPY --from=builder /go/src/github.com/icecarti/feed_service/database/migrations ./database/migrations

USER appuser

EXPOSE 50052

ENTRYPOINT ["./feed_service"]