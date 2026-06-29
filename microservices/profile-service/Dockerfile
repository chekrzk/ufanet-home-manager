FROM golang:1.25-alpine AS builder

WORKDIR /src
COPY contracts ./contracts
COPY microservices/profile-service ./microservices/profile-service
WORKDIR /src/microservices/profile-service
RUN go build -o /bin/profile-service ./cmd

FROM alpine:3.22
WORKDIR /app
RUN apk add --no-cache tzdata
COPY --from=builder /bin/profile-service /app/profile-service
EXPOSE 50055
CMD ["/app/profile-service"]
