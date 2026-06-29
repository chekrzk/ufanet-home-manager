FROM golang:1.25-alpine AS builder

WORKDIR /src
COPY contracts ./contracts
COPY microservices/notification-service ./microservices/notification-service
WORKDIR /src/microservices/notification-service
RUN go build -o /bin/notification-service ./cmd

FROM alpine:3.22
WORKDIR /app
RUN apk add --no-cache tzdata
COPY --from=builder /bin/notification-service /app/notification-service
EXPOSE 50054
CMD ["/app/notification-service"]
