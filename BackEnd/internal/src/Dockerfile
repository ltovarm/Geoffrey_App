# Use the latest Golang image as the base image
FROM golang:latest

# Set the working directory in the container
WORKDIR /app

# Copy the Go program to the container
COPY . .

# Build the Go program
RUN go build -o main .

# Specify the command to run when the container starts
CMD ["./main"]
