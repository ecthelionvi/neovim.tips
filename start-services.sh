#!/bin/bash

# Navigate to the Go backend directory
cd /neovim-tips/backend

# Start the Go backend
RUN ./neovim-tips &

# Navigate to the Next.js frontend directory
cd /neovim-tips/frontend/neovim-tips

# Start the Next.js application
npm start &

# Start Nginx in the foreground
nginx -g 'daemon off;'
