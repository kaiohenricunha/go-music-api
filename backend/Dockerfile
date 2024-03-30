# Use a lightweight base image; here we're using the official Go image.
FROM golang:1.21-alpine as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go mod and sum files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o backend .

# Use a Docker multi-stage build to keep the final image lean
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/backend .

# Command to run
CMD ["./backend"]
