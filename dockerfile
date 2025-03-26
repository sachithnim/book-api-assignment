# Use the official Golang image as the base image
FROM golang:1.24-alpine

# Set the working directory in the container
WORKDIR /app

# Copy the Go module files and download the dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the entire project into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the port the app will run on
EXPOSE 8000

# Run the Go application
CMD ["./main"]
