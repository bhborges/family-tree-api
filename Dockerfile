# Use golang version 1.20 as the base image
FROM golang:1.20-alpine

# Set the working directory to /app
WORKDIR /app

# Copy the local package files to the container's workspace
COPY go.mod ./
COPY go.sum ./

# Build the Go app
RUN go mod download

COPY . .

RUN go build ./...

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./familytree"]
