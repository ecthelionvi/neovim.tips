# Use the official Ubuntu base image
FROM ubuntu:latest

# Install Nginx, Git, wget, curl, and other necessary tools
RUN apt-get update && \
    apt-get install -y nginx git wget curl ca-certificates gnupg

# Install Node.js
RUN mkdir -p /etc/apt/keyrings && \
    curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg && \
    NODE_MAJOR=20 && \
    echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_$NODE_MAJOR.x nodistro main" | tee /etc/apt/sources.list.d/nodesource.list && \
    apt-get update && \
    apt-get install -y nodejs

# Install Python
RUN apt-get install -y python3 python3-pip

# Install Go 1.21.4 for ARM64
RUN wget https://go.dev/dl/go1.21.4.linux-arm64.tar.gz && \
    tar -C /usr/local -xzf go1.21.4.linux-arm64.tar.gz && \
    rm go1.21.4.linux-arm64.tar.gz

# Set the Go environment variables
ENV PATH="${PATH}:/usr/local/go/bin"

# Clone the repository
RUN git clone https://github.com/ecthelionvi/neovim.tips.git /neovim-tips

# Build Backend
WORKDIR /neovim-tips/backend
RUN go build

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
COPY start_services.py /start_services.py
RUN chmod +x /start_services.py
CMD ["python3", "/start_services.py"]
