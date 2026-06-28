FROM golang:1.25-alpine AS builder

WORKDIR /src
COPY contracts ./contracts
COPY microservices/requests-service ./microservices/requests-service
WORKDIR /src/microservices/requests-service
RUN go build -o /bin/requests-service ./cmd

FROM alpine:3.22
WORKDIR /app
COPY --from=builder /bin/requests-service /app/requests-service
EXPOSE 50053
CMD ["/app/requests-service"]
