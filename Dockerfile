FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o foundry_api ./cmd/api.go


FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/foundry_api .

CMD ["./foundry_api"]