#!/bin/bash
set -e

echo "Installing GOLINE Looking Glass..."

# Check if running as root
if [[ $EUID -eq 0 ]]; then
   echo "Error: Don't run this script as root" 
   exit 1
fi

# Detect OS
if [[ -f /etc/debian_version ]]; then
    OS="debian"
elif [[ -f /etc/redhat-release ]]; then
    OS="redhat"
else
    echo "Error: Unsupported OS"
    exit 1
fi

# Install dependencies
echo "Installing dependencies..."
if [[ "$OS" == "debian" ]]; then
    sudo apt update
    sudo apt install -y golang-go git apache2 ssl-cert logrotate curl wget
    sudo a2enmod proxy proxy_http ssl rewrite headers expires
elif [[ "$OS" == "redhat" ]]; then
    sudo yum install -y golang git httpd mod_ssl logrotate curl wget
    sudo systemctl enable httpd
fi

# Create user and directories
echo "Creating looking-glass user..."
sudo useradd -r -s /bin/false looking-glass || true
sudo mkdir -p /opt/looking-glass/{logs,config,public/images}
sudo chown -R looking-glass:looking-glass /opt/looking-glass

echo "Installation complete!"
echo "Next steps:"
echo "1. Configure routers: cp config/routers.example.json config/routers.json"
echo "2. Configure company: cp config/company.example.json config/company.json"
echo "3. Build and install: make build && sudo make install"
