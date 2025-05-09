FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build migrate binary
RUN go build -o /app/migrate ./cmd/migrate
RUN go build -o /app/server ./cmd/server

# Final runtime image
FROM alpine:latest

# Copy binaries from builder
COPY --from=builder /app/migrate /app/migrate
COPY --from=builder /app/server /app/server

EXPOSE 8080
CMD ["/app/server"]
