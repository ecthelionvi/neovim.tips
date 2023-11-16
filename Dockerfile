# Use the official Ubuntu base image
FROM ubuntu:latest

# Install Nginx, Git, Go, and other necessary tools
RUN apt-get update && \
  apt-get install -y nginx git wget

# Install Go 1.23
RUN wget https://go.dev/dl/go1.21.4.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.21.4.linux-amd64.tar.gz
ENV PATH="${PATH}:/usr/local/go/bin"

# Clone the repository
RUN git clone https://github.com/ecthelionvi/neovim.tips.git /neovim-tips

# Set the working directory to the application directory
WORKDIR /neovim-tips

# Assuming there's a main.go in the root of your repository:
RUN go build -o neovim-tips

# Copy the Nginx configuration file into the container
COPY default /etc/nginx/sites-enabled/default

# Copy the .env file into the container
COPY .env /neovim-tips/.env

# Expose the port that Nginx is listening on
EXPOSE 80

# Use a multi-command CMD or a script to start both Nginx and your application
CMD nginx -g 'daemon off;' & ./neovim-tips
