# Stage 1: Build the Go binary
FROM golang:1.18-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o livekit-server .

# Stage 2: Create a minimal runtime image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/livekit-server .
EXPOSE 7880
ENTRYPOINT ["./livekit-server"]
