# Use the official Golang image as the base image
FROM golang:1.22

# Set the working directory inside the container
WORKDIR /app

# Copy the Go source code into the container
COPY . .

# Ensure dependencies are downloaded
RUN go mod tidy

# Expose port 9090
EXPOSE 9090

# Run the Go application
CMD ["go", "run", "main.go"]
