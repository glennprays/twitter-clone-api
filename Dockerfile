# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Set the environment variables from the .env file
ARG DB_USER
ARG DB_PASSWORD
ARG DB_HOST
ARG DB_PORT
ENV DB_USER=${DB_USER}
ENV DB_PASSWORD=${DB_PASSWORD}
ENV DB_HOST=${DB_HOST}
ENV DB_PORT=${DB_PORT}

# Build the Go application
RUN go build -o app .


# Expose any necessary ports
EXPOSE 8080

# Define the command to run the Go application
CMD ["./app"]
