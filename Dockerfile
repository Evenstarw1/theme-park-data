FROM golang:1.23-alpine AS builder

# Install git for go get command
RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/restapi

# Use a minimal base image to package the binary
FROM alpine:3.14

# Add ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder image
COPY --from=builder /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the compiled binary
CMD ["./main"]
