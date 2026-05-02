# ── Build stage ──────────────────────────────────────────────────────────────
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# ── Run stage ─────────────────────────────────────────────────────────────────
FROM alpine:3.20

WORKDIR /app

# ca-certificates needed for HTTPS calls (e.g. sentiment API)
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server .
COPY migrations ./migrations

EXPOSE 8080

CMD ["./server"]
