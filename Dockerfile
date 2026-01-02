FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN ["go", "build", "-o", "gollab-backend", "./cmd"]

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/gollab-backend .

EXPOSE 8080

CMD ["./gollab-backend"]
