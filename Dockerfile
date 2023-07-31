# Use the official Go image as the base image
FROM golang:1.17

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application source code
COPY . .

# Build the Go application
RUN go build -o ProjectHash

# Set the entry point for the container
ENTRYPOINT ["./ProjectHash"]
