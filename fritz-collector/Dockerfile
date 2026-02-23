# Phase 1: Kompilieren
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o fritz-monitor

# Phase 2: Ausführen (schlankes Image)
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/fritz-monitor .
# Wir brauchen Zertifikate für HTTPS-Anfragen
RUN apk add --no-cache ca-certificates
CMD ["./fritz-monitor"]