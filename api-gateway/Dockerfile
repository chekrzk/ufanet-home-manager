FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY contracts ./contracts
COPY api-gateway ./api-gateway

WORKDIR /src/api-gateway

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/api-gateway ./cmd

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /bin/api-gateway /app/api-gateway

EXPOSE 8086

ENTRYPOINT ["/app/api-gateway"]

