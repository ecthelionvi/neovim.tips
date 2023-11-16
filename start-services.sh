#!/bin/bash

# Start the Go backend
/neovim-tips/backend/neovim-tips-backend &

# Navigate to the Next.js frontend directory
cd /neovim-tips/frontend/neovim-tips

# Start the Next.js application
npm start &

# Start Nginx in the foreground
nginx -g 'daemon off;'
