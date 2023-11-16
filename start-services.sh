#!/bin/bash

# Start Nginx
nginx -g 'daemon off;' &

# Start the Go backend
/neovim-tips/backend/neovim-tips-backend &

# Serve the Next.js frontend
cd /neovim-tips/frontend/neovim-tips
npm start
