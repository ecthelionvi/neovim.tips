# Guide to Setting Up a Systemd Service and Configuring Nginx

This guide provides step-by-step instructions for creating a Systemd service for a project (here named `neovim-tips`) and configuring Nginx as a web server to serve the project.

## Setting Up a Systemd Service

Systemd services allow you to manage and control background services on your Linux system. Here's how to set up a service for `neovim-tips`.

1. **Create a New Service File**:
   Use `nano` (or your preferred text editor) to create a new service file under `/etc/systemd/system/`.
   ```
   sudo nano /etc/systemd/system/neovim-tips.service
   ```

   ```
   [Unit]
   Description=neovim-tips

   [Service]
   ExecStart=//home/ubuntu/neovim.tips/neovim-tips
   WorkingDirectory=/home/ubuntu/neovim.tips
   User=ubuntu
   Group=ubuntu
   Restart=always
   RestartSec=5s

   [Install]
   WantedBy=multi-user.target
   ```

2. **Check Service Status**:
   To check the current status of your newly created service, use:
   ```
   sudo systemctl status neovim-tips.service
   ```

3. **Stop the Service**:
   If the service is running and you need to stop it, use:
   ```
   sudo systemctl stop neovim-tips.service
   ```

4. **Restart the Service**:
   To restart the service after making changes, use:
   ```
   sudo systemctl restart neovim-tips.service
   ```

5. **Reload System Daemon**:
   After making changes to the service file, reload the systemd daemon to apply changes:
   ```
   sudo systemctl daemon-reload
   ```

## Configuring Nginx

Nginx is a powerful and efficient web server. Here’s how to configure it to serve your project.

1. **Install Nginx**:
   Install Nginx using the package manager:
   ```
   sudo apt install nginx -y
   ```

2. **Edit Nginx Configuration**:
   Open the default Nginx configuration file:
   ```
   sudo nano /etc/nginx/sites-enabled/default
   ```
   ```
   server {
        listen 80 default_server;
        listen [::]:80 default_server;

        root /var/www/html;

        server_name neovim.tips www.neovim.tips;

        location /api/ {
                proxy_pass http://localhost:8080; # Forward requests to the Go application
                proxy_http_version 1.1;
                proxy_set_header Upgrade $http_upgrade;
                proxy_set_header Connection 'upgrade';
                proxy_set_header Host $host;
                proxy_cache_bypass $http_upgrade;
                proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
                proxy_set_header X-Forwarded-Proto $scheme;

                # Allow CORS
                if ($request_method = 'OPTIONS') {
                        add_header 'Access-Control-Allow-Origin' '*';
                        add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS';
                        add_header 'Access-Control-Allow-Headers' 'DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type>
                        add_header 'Access-Control-Max-Age' 1728000;
                        add_header 'Content-Type' 'text/plain charset=UTF-8';
                        add_header 'Content-Length' 0;
                        return 204;
                }

                if ($request_method = 'POST') {
                        add_header 'Access-Control-Allow-Origin' '*';
                        add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS';
                        add_header 'Access-Control-Allow-Headers' 'DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type>
                }

                if ($request_method = 'GET') {
                        add_header 'Access-Control-Allow-Origin' '*';
                        add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS';
                        add_header 'Access-Control-Allow-Headers' 'DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type>
                }

        }

        location / {
                # First attempt to serve request as file, then
                # as directory, then fall back to displaying a 404.
                try_files $uri $uri/ =404;
        }
   }
   
   ```

   Replace or modify the contents to suit your project's needs. Key configurations include:
   - `root /var/www/html;`: Sets the root directory for your site.
   - `server_name`: Specifies the domain name of your server.
   - `location /api/`: Configures the handling of the API endpoint.
   - `proxy_pass`: Forwards requests to your application, running perhaps on `localhost:8080`.
   - `add_header`: Sets various headers, including CORS (Cross-Origin Resource Sharing) settings.

3. **Test Nginx Configuration**:
   After making changes, test your Nginx configuration for syntax errors:
   ```
   sudo nginx -t
   ```

4. **Reload Nginx**:
   Apply your changes by reloading Nginx:
   ```
   sudo systemctl reload nginx
   ```

Following these steps should set up a Systemd service for `neovim-tips` project and configure Nginx to serve it efficiently.
