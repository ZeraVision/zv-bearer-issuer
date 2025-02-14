# Use the official Go image as a parent image
FROM golang:1.20 AS builder

# Setup go private repo download
ENV GOPRIVATE=github.com/*
# Set up the SSH key
ARG SSH_KEY
RUN mkdir -p /root/.ssh && echo "$SSH_KEY" > /root/.ssh/id_rsa && chmod 600 /root/.ssh/id_rsa
RUN ssh-keyscan github.com >> /root/.ssh/known_hosts

# Set destination for COPY
WORKDIR /app

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
# Refer to .dockerignore for excluded files
COPY . .

# Download Go modules
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o bearer-issuer .

# Create a smaller final image
FROM alpine:3.18
# Install curl
RUN apk --no-cache add curl

WORKDIR /app
COPY --from=builder /app/bearer-issuer .
# Expose the port the application will run on
EXPOSE 80

# Command to run the executable
CMD ["./bearer-issuer"]
