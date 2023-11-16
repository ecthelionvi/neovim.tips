# Use the official Ubuntu base image
FROM ubuntu:latest

# Use the official Node.js 21 base image for linux/amd64
FROM --platform=linux/amd64 node:21 as builder

# Install Nginx, Git, Go, Node.js, npm, and other necessary tools
RUN apt-get update && \
    apt-get install -y nginx git wget curl && \
    curl -fsSL https://deb.nodesource.com/setup_21.x | bash - && \
    apt-get install -y nodejs

# Install Go 1.21.4
RUN wget https://go.dev/dl/go1.21.4.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.21.4.linux-amd64.tar.gz && \
    rm go1.21.4.linux-amd64.tar.gz

ENV PATH="${PATH}:/usr/local/go/bin"

# Clone the repository
RUN git clone https://github.com/ecthelionvi/neovim.tips.git /neovim-tips

# Build Backend
WORKDIR /neovim-tips/backend
RUN go build -o neovim-tips-backend

# Build Frontend
WORKDIR /neovim-tips/frontend/neovim-tips
RUN npm install && npm run build

# Set the working directory back to the root of your project
WORKDIR /neovim-tips

# Copy the Nginx configuration file into the container
COPY default /etc/nginx/sites-enabled/default

# Copy the .env file into the container (if necessary)
COPY .env /neovim-tips/backend/.env

# Expose the port that Nginx is listening on
EXPOSE 80

# Use a script to start both Nginx, the backend, and serve the frontend
COPY start-services.sh /start-services.sh
RUN chmod +x /start-services.sh
CMD ["/start-services.sh"]
