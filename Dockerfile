# Stage 1: Build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o go-app .

# Stage 2: Minimal runtime
FROM scratch
WORKDIR /app
COPY --from=builder /app/go-app .
COPY --from=builder /app/static ./static
EXPOSE 8080
CMD ["./go-app"]
