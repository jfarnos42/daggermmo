FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o daggerfall ./cmd/server/main.go

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

RUN adduser --disabled-password --gecos '' daggeruser

WORKDIR /home/daggeruser

COPY --from=builder /app/daggerfall .
COPY --from=builder /app/cert.pem /home/daggeruser/
COPY --from=builder /app/key.pem /home/daggeruser/

RUN chown -R daggeruser:daggeruser /home/daggeruser

USER daggeruser

CMD ["./daggerfall"]
