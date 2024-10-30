# Start from a base image
FROM alpine:latest

# Set the working directory in the container
WORKDIR /app

# Copy the binary file from your host to your present location (PWD) in the image
COPY ./bin/gostarter .
COPY .gostarter.yaml.docker .gostarter.yaml

# Command to run the executable
ENTRYPOINT ["./gostarter", "serve"]