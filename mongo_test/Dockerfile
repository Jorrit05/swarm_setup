# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory
WORKDIR /app

# Copy the source code into the container
COPY . .

# Download and install any necessary dependencies
RUN go mod download

# Build the application
RUN go build -o main .

# Expose port 8080 to the host
EXPOSE 8080

# Start the application
CMD ["./main"]
