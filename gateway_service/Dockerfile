# Use the official Golang image as the base image
FROM golang:1.19.5-alpine3.17

# Set the working directory
WORKDIR /app

# Copy the source code into the container
COPY . .

# Download and install any necessary dependencies
RUN go mod download

# Build the application
RUN go build -o main .

EXPOSE 3000
# Start the application
CMD ["./main"]
