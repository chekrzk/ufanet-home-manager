FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY contracts ./contracts
COPY microservices/auth-service ./microservices/auth-service

WORKDIR /src/microservices/auth-service

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/auth-service ./cmd

FROM alpine:3.21

WORKDIR /app
RUN apk add --no-cache tzdata

COPY --from=builder /bin/auth-service /app/auth-service

EXPOSE 50051

ENTRYPOINT ["/app/auth-service"]
