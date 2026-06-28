FROM golang:1.25-alpine AS builder

WORKDIR /src
COPY contracts ./contracts
COPY microservices/news-service ./microservices/news-service
WORKDIR /src/microservices/news-service
RUN go build -o /bin/news-service ./cmd

FROM alpine:3.22
WORKDIR /app
COPY --from=builder /bin/news-service /app/news-service
EXPOSE 50052
CMD ["/app/news-service"]
