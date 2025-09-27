# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Ensure content directories exist even if repo doesn't include them
RUN mkdir -p content/posts content/pages

# Build the application (modernc.org/sqlite is pure-Go; keep CGO disabled)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy templates and static files
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/content ./content
COPY --from=builder /app/data ./data

# Create directories for user content
RUN mkdir -p content/posts content/pages data

# Expose port
EXPOSE 3000

# Run the application
CMD ["./main"]