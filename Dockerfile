# Use the official Node.js base image for the frontend build
FROM node:latest as frontend-builder

# Set the working directory for the frontend build
WORKDIR /app/frontend

# Copy the frontend application files
COPY ./neovim-tips/frontend/neovim-tips /app/frontend

# Install frontend dependencies and build the static files
RUN npm install && npm run build --verbose

# Use the official Ubuntu base image for the backend and final image
FROM ubuntu:latest

# Install Nginx, Git, and other necessary tools
RUN apt-get update && \
    apt-get install -y nginx git wget curl

# Install the latest version of Go
RUN wget https://go.dev/dl/go1.18.1.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.18.1.linux-amd64.tar.gz && \
    rm go1.18.1.linux-amd64.tar.gz

ENV PATH="${PATH}:/usr/local/go/bin"

# Clone the repository if needed, or copy your local backend code directly
# RUN git clone https://github.com/ecthelionvi/neovim-tips.git /neovim-tips
COPY ./neovim-tips/backend /neovim-tips/backend

# Build Backend
WORKDIR /neovim-tips/backend
RUN go build -o neovim-tips-backend

# Copy the frontend build from the frontend-builder stage
COPY --from=frontend-builder /app/frontend/.next /neovim-tips/frontend/.next
COPY --from=frontend-builder /app/frontend/node_modules /neovim-tips/frontend/node_modules
COPY --from=frontend-builder /app/frontend/public /neovim-tips/frontend/public
COPY --from=frontend-builder /app/frontend/package.json /neovim-tips/frontend/package.json

# Set the working directory back to the root of your project
WORKDIR /neovim-tips

# Copy the Nginx configuration file into the container
COPY default /etc/nginx/sites-enabled/default

# Copy the .env file and other necessary scripts into the container
COPY .env /neovim-tips/backend/.env
COPY start-services.sh /start-services.sh
RUN chmod +x /start-services.sh

# Expose the port that Nginx is listening on
EXPOSE 80

# Start services
CMD ["/start-services.sh"]
