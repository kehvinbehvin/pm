# Use an official Go image as the base image
FROM golang:1.20-alpine

# Install Zsh and necessary tools
RUN apk add --no-cache zsh

# Set Zsh as the default shell
RUN echo "exec /bin/zsh" >> ~/.bashrc

# Create a directory for your Go project
WORKDIR /usr/src/app

# Start with Zsh when the container is started
CMD ["zsh"]
