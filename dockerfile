# Stage 1: Build the binary
FROM golang:1.23rc1-alpine3.20 AS builder

# Create and change to the app directory.
WORKDIR /app

# Copy go mod and sum files.
COPY go.mod go.sum ./

# Download dependencies.
RUN go mod download

# Copy the source code.
COPY . .

# Build the binary.
RUN go build -o main main.go

# Stage 2: Run the binary
FROM alpine:latest

# Copy the binary from the builder stage.
COPY --from=builder /app/main /app/main
COPY ./internal/database/sql/creation_script.sql ./app/internal/database/sql/creation_script.sql


# Set the working directory.
WORKDIR /app

# Command to run the binary.
CMD ["./main"]

