# --- builder ---
FROM golang:1.25 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . ./


RUN go build \
    -ldflags="-s -w" \
    -o /app/server .

# --- runtime ---
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/server /app/server

EXPOSE ${AUTH_APP_PORT}

ENV AUTH_APP_PORT=${AUTH_APP_PORT}
CMD ["/app/server"]

