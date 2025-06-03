# ?? Installation Guide

Complete installation guide for GOLINE Looking Glass.

## ?? System Requirements

### Minimum Requirements
- **OS**: Ubuntu 20.04+, CentOS 8+, RHEL 8+, Debian 11+
- **RAM**: 512MB (1GB recommended)
- **CPU**: 1 vCPU (2+ recommended)
- **Disk**: 1GB free space
- **Network**: SSH access to routers

### Software Dependencies
- Go 1.21+ (for building from source)
- Apache 2.4+ or Nginx 1.18+ (for SSL termination)
- OpenSSH client
- Git

## ?? Installation Methods

### Method 1: Automated Script (Recommended)

```bash
# Download and run installation script
curl -sSL https://raw.githubusercontent.com/paolokappa/goline-looking-glass/main/scripts/install.sh | bash
```

The script will:
1. ? Install all dependencies
2. ? Create system user and directories
3. ? Build the Go application
4. ? Install systemd service
5. ? Configure basic Apache setup

### Method 2: Manual Installation

#### Step 1: Install Dependencies

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install -y golang-go git apache2 ssl-cert logrotate curl wget build-essential

# Enable Apache modules
sudo a2enmod proxy proxy_http ssl rewrite headers expires
```

**CentOS/RHEL:**
```bash
sudo yum update
sudo yum install -y golang git httpd mod_ssl logrotate curl wget gcc

# Enable and start Apache
sudo systemctl enable httpd
sudo systemctl start httpd
```

#### Step 2: Create System User

```bash
# Create dedicated user for security
sudo useradd -r -s /bin/false -d /opt/looking-glass looking-glass

# Create directories
sudo mkdir -p /opt/looking-glass/{logs,config,public/images,backups}
sudo chown -R looking-glass:looking-glass /opt/looking-glass
sudo chmod 755 /opt/looking-glass
```

#### Step 3: Build Application

```bash
# Clone repository
git clone https://github.com/paolokappa/goline-looking-glass.git
cd goline-looking-glass

# Build for production
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o looking-glass main.go

# Install binary
sudo cp looking-glass /opt/looking-glass/
sudo chown looking-glass:looking-glass /opt/looking-glass/looking-glass
sudo chmod +x /opt/looking-glass/looking-glass
```

#### Step 4: Install Configuration Files

```bash
# Copy configuration templates
sudo cp -r config/* /opt/looking-glass/config/
sudo cp -r public/* /opt/looking-glass/public/
sudo chown -R looking-glass:looking-glass /opt/looking-glass/config
sudo chown -R www-data:www-data /opt/looking-glass/public
```

#### Step 5: Install Systemd Service

```bash
# Create systemd service file
sudo tee /etc/systemd/system/looking-glass.service <<'EOF'
[Unit]
Description=GOLINE Looking Glass - Network Diagnostic Tool
Documentation=https://github.com/paolokappa/goline-looking-glass
After=network.target
Wants=network.target

[Service]
Type=simple
User=looking-glass
Group=looking-glass
WorkingDirectory=/opt/looking-glass
ExecStart=/opt/looking-glass/looking-glass
ExecReload=/bin/kill -HUP $MAINPID
KillMode=mixed
KillSignal=SIGTERM
TimeoutStopSec=30
RestartSec=10
Restart=always

# Security settings
NoNewPrivileges=yes
ProtectSystem=strict
ReadWritePaths=/opt/looking-glass/logs
ProtectHome=yes
PrivateTmp=yes

# Environment
Environment=GIN_MODE=release
Environment=PORT=3002
Environment=CONFIG_PATH=/opt/looking-glass/config

[Install]
WantedBy=multi-user.target
EOF

# Enable and start service
sudo systemctl daemon-reload
sudo systemctl enable looking-glass.service
sudo systemctl start looking-glass.service
```

### Method 3: Docker Installation

#### Using Docker Compose (Recommended)

```bash
# Clone repository
git clone https://github.com/paolokappa/goline-looking-glass.git
cd goline-looking-glass

# Configure environment
cp config/company.example.json config/company.json
cp config/routers.example.json config/routers.json

# Edit configurations
nano config/company.json
nano config/routers.json

# Start with Docker Compose
docker-compose up -d

# Check status
docker-compose ps
docker-compose logs -f looking-glass
```

#### Manual Docker Build

```bash
# Build image
docker build -t looking-glass:latest .

# Run container
docker run -d \
  --name looking-glass \
  -p 3002:3002 \
  -v $(pwd)/config:/root/config:ro \
  -v $(pwd)/logs:/root/logs \
  -e GIN_MODE=release \
  looking-glass:latest
```

## ?? Post-Installation Configuration

### 1. Configure Company Information

```bash
# Edit company configuration
sudo nano /opt/looking-glass/config/company.json
```

Example configuration:
```json
{
  "company": {
    "name": "Your ISP Name",
    "as_number": "AS64512",
    "domain": "yourisp.com",
    "support_email": "noc@yourisp.com",
    "logo_path": "/images/logo.png",
    "theme": {
      "background_color": "#1e3a8a",
      "primary_color": "#667eea",
      "secondary_color": "#764ba2"
    }
  }
}
```

### 2. Configure Routers

```bash
# Edit router configuration
sudo nano /opt/looking-glass/config/routers.json
```

See [Router Configuration Guide](router-config.md) for detailed examples.

### 3. Setup SSL Certificate

```bash
# Using Let's Encrypt (recommended)
./scripts/setup-ssl.sh yourdomain.com noc@yourdomain.com

# Or configure manually
sudo nano /etc/apache2/sites-available/looking-glass-ssl.conf
```

### 4. Configure Log Rotation

```bash
# Install logrotate configuration
sudo tee /etc/logrotate.d/looking-glass <<'LOGROTATE_EOF'
/opt/looking-glass/logs/*.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    create 0644 looking-glass looking-glass
    copytruncate
}
LOGROTATE_EOF

# Test logrotate configuration
sudo logrotate -d /etc/logrotate.d/looking-glass
```

## ? Verification

### 1. Check Service Status

```bash
# Service status
sudo systemctl status looking-glass.service

# View logs
sudo journalctl -u looking-glass.service -f

# Application logs
sudo tail -f /opt/looking-glass/logs/lg.log
```

### 2. Test HTTP Endpoint

```bash
# Test local endpoint
curl http://localhost:3002/

# Test API endpoints
curl http://localhost:3002/api/routers
curl http://localhost:3002/health
```

### 3. Test SSL Configuration

```bash
# Test SSL certificate
openssl s_client -connect yourdomain.com:443 -servername yourdomain.com

# Test SSL configuration
curl -I https://yourdomain.com/
```

## ?? Firewall Configuration

### UFW (Ubuntu)

```bash
# Allow SSH (if not already configured)
sudo ufw allow ssh

# Allow HTTP and HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Allow application port (if accessed directly)
sudo ufw allow 3002/tcp

# Enable firewall
sudo ufw enable
```

### Firewalld (CentOS/RHEL)

```bash
# Allow HTTP and HTTPS
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https

# Allow application port (if needed)
sudo firewall-cmd --permanent --add-port=3002/tcp

# Reload firewall
sudo firewall-cmd --reload
```

## ?? Troubleshooting

### Common Issues

#### 1. Service Won't Start

```bash
# Check service status and logs
sudo systemctl status looking-glass.service
sudo journalctl -u looking-glass.service --no-pager

# Check binary permissions
ls -la /opt/looking-glass/looking-glass

# Check configuration files
sudo -u looking-glass /opt/looking-glass/looking-glass --test-config
```

#### 2. Permission Issues

```bash
# Fix ownership
sudo chown -R looking-glass:looking-glass /opt/looking-glass
sudo chown -R www-data:www-data /opt/looking-glass/public

# Fix permissions
sudo chmod 755 /opt/looking-glass
sudo chmod +x /opt/looking-glass/looking-glass
sudo chmod 644 /opt/looking-glass/config/*.json
```

#### 3. Network Issues

```bash
# Check if port is listening
sudo netstat -tlnp | grep :3002

# Test firewall
sudo ufw status
sudo iptables -L

# Check Apache proxy
sudo apache2ctl configtest
sudo systemctl status apache2
```

## ?? Performance Tuning

### 1. System Limits

```bash
# Increase file descriptor limits for the user
sudo tee -a /etc/security/limits.conf <<'LIMITS_EOF'
looking-glass soft nofile 65536
looking-glass hard nofile 65536
LIMITS_EOF

# Or use systemd service limits
sudo mkdir -p /etc/systemd/system/looking-glass.service.d
sudo tee /etc/systemd/system/looking-glass.service.d/limits.conf <<'SERVICE_LIMITS_EOF'
[Service]
LimitNOFILE=65536
SERVICE_LIMITS_EOF

sudo systemctl daemon-reload
sudo systemctl restart looking-glass.service
```

### 2. Apache Performance

```bash
# Enable compression
sudo a2enmod deflate
sudo systemctl restart apache2

# Configure caching for static assets
sudo tee -a /etc/apache2/sites-available/looking-glass.conf <<'APACHE_PERF_EOF'
<Directory "/opt/looking-glass/public">
    # Enable compression
    <IfModule mod_deflate.c>
        AddOutputFilterByType DEFLATE text/plain
        AddOutputFilterByType DEFLATE text/html
        AddOutputFilterByType DEFLATE text/css
        AddOutputFilterByType DEFLATE application/javascript
    </IfModule>
    
    # Set cache headers
    <IfModule mod_expires.c>
        ExpiresActive On
        ExpiresByType image/png "access plus 1 month"
        ExpiresByType text/css "access plus 1 month"
        ExpiresByType application/javascript "access plus 1 month"
        ExpiresDefault "access plus 2 days"
    </IfModule>
</Directory>
APACHE_PERF_EOF
```

## ?? Backup Configuration

```bash
# Create backup script
sudo tee /opt/looking-glass/backup.sh <<'EOF'
#!/bin/bash
BACKUP_DATE=$(date +%Y%m%d_%H%M%S)
tar -czf /opt/looking-glass/backups/backup_$BACKUP_DATE.tar.gz \
    -C /opt/looking-glass \
    config/ public/ looking-glass
echo "Backup created: backup_$BACKUP_DATE.tar.gz"
EOF

sudo chmod +x /opt/looking-glass/backup.sh
```

## ?? Support

For installation support:
- ?? **Email**: [noc@goline.ch](mailto:noc@goline.ch)
- ?? **Documentation**: [GitHub Wiki](https://github.com/paolokappa/goline-looking-glass/wiki)
- ?? **Issues**: [GitHub Issues](https://github.com/paolokappa/goline-looking-glass/issues)

## ?? Next Steps

1. ?? [Configure your routers](router-config.md)
2. ?? [Customize the appearance](customization.md)
3. ?? [Review security settings](security.md)
4. ?? [Setup monitoring](monitoring.md)

---

**Installation Guide** | GOLINE Looking Glass  
*Paolo Caparrelli - GOLINE SA*
