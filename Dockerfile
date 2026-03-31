# Stage 1: Build the Go binary
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
# Build the binary
RUN go build -o go-app .

# Stage 2: Create a lightweight image
FROM alpine:3.19
WORKDIR /app
# Copy binary and static files from builder
COPY --from=builder /app/go-app .
COPY --from=builder /app/static ./static

EXPOSE 8080

# Run the executable
CMD ["./go-app"]