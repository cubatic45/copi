# Build
FROM golang:1.22.1-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o copi

# Deploy
FROM debian:stretch-slim

WORKDIR /app
COPY --from=builder /app/copi .

# Copy all system certificates from the builder stage
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

EXPOSE 8081
ENTRYPOINT ["./copi", "-token", "", "-token_url", ""]

