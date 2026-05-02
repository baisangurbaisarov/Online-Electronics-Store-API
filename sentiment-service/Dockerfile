FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o sentiment .

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/sentiment .

EXPOSE 9090

CMD ["./sentiment"]