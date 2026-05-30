FROM golang:1.24 AS builder

WORKDIR /build

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o application main.go

FROM alpine:latest AS runner

WORKDIR /car_estimator_api_gateway

COPY --from=builder /build/application .

EXPOSE 4242

CMD [ "./application" ]
