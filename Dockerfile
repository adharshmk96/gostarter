# Start from a base image
FROM alpine:latest

# Set the working directory in the container
WORKDIR /app

# Copy the binary file from your host to your present location (PWD) in the image
COPY ./bin/sharepoint_sync .
COPY .sharepoint_sync.yaml .

# Command to run the executable
ENTRYPOINT ["./sharepoint_sync", "serve"]