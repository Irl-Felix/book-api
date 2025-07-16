# --- Build stage ---
FROM golang:latest AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY app/ ./app/


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./app/main.go

# --- Runtime stage ---
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]