FROM golang:1.24 as builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o daggerfall ./cmd/server/main.go

FROM alpine:3.19

RUN apk add --no-cache ca-certificates
RUN adduser -D -g '' daggeruser

WORKDIR /home/daggeruser

COPY --from=builder /app/daggerfall .

USER daggeruser

CMD ["./daggerfall"]
