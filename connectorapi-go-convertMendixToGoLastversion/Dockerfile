# Stage 1: Build stage
FROM harbor-private.aeonth.com/container-images/library/golang:1.23.2-alpine AS builder
 
# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum for dependency management
COPY go.mod go.sum ./

# Download and verify dependencies
RUN go mod download && go mod verify

# Copy the entire project into the container
COPY . .

# Change to the directory containing main.go and build the binary
WORKDIR /app/cmd/server

# Build app (generate the `server` binary in the `cmd/server` folder)
RUN go build -o server .

# Stage 2: Run stage
FROM alpine:latest  

# Set the working directory inside the container
WORKDIR /root/

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/cmd/server/server .

# Copy the certificate files from local machine to the container
#COPY certs/CA-sit.cer /etc/ssl/certs/CA-sit.cer
COPY certs/ /etc/ssl/certs/

# Copy the config folder into the container
COPY configs/ /root/configs/

# Expose port for the application
EXPOSE 8080

# Command to run the executable with the port as an argument
CMD ["./server"]
