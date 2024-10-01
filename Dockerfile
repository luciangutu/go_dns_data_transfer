FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o dns-server


FROM alpine:latest

RUN apk --no-cache add bind-tools
WORKDIR /app
RUN mkdir data
COPY --from=builder /app/dns-server .
EXPOSE 53/udp

ENTRYPOINT ["./dns-server", "-ip"]
