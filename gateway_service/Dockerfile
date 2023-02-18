FROM golang:1.19.5-alpine3.17
ARG NAME
ENV NAME=$NAME

# Set the working directory
WORKDIR /app

# Copy the source code into the container
COPY . .

# Download and install any necessary dependencies
RUN go mod download

# Build the application
RUN go build -o $NAME .

EXPOSE 3000

# Start the application
CMD  [ "/bin/sh", "-c", "/app/$NAME" ]
