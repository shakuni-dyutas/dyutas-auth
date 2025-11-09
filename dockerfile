# --- builder ---
FROM golang:1.25 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . ./


RUN go build \
    -ldflags="-s -w" \
    -o /app/server ./cmd/dyutas-auth

# --- runtime ---
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/server /app/server

CMD ["/app/server"]

