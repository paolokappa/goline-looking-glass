# ?? Router Configuration Guide

This comprehensive guide covers configuration for all supported router vendors.

## ?? Overview

The Looking Glass supports three major router vendors:

| Vendor | OS | Tested Versions | SSH Key Support | Notes |
|--------|----|-----------------|--------------------|-------|
| **Juniper** | Junos | 15.1+, 18.4+, 20.4+ | ? Yes | Full feature support |
| **Huawei** | VRP | V200R010+, V800R011+ | ? Yes | BGP and basic commands |
| **Cisco** | IOS/IOS-XE | 15.0+, 16.0+, 17.0+ | ? Yes | IOS and IOS-XE supported |

## ?? Security Best Practices

### 1. Create Dedicated User Account

**Never use admin/root accounts for Looking Glass access.**

#### Juniper (Junos)
```bash
# Create user with minimal privileges
set system login user looking-glass uid 2001
set system login user looking-glass class read-only
set system login user looking-glass authentication ssh-rsa "ssh-rsa AAAAB3..."

# Custom class with specific permissions (recommended)
set system login class looking-glass permissions view
set system login class looking-glass permissions view-configuration
set system login class looking-glass allow-commands "show route"
set system login class looking-glass allow-commands "show bgp"
set system login class looking-glass allow-commands "ping"
set system login class looking-glass allow-commands "traceroute"
set system login class looking-glass deny-commands "configure"
set system login class looking-glass deny-commands "request"

# Apply custom class to user
set system login user looking-glass class looking-glass
```

#### Huawei (VRP)
```bash
# Create user group with limited privileges
user-group looking-glass

# Create user with SSH key
local-user looking-glass
 password irreversible-cipher YourSecurePassword123!
 user-group looking-glass
 service-type ssh
 authorization-attribute user-role network-operator

# SSH key configuration (preferred)
ssh user looking-glass
ssh user looking-glass authentication-type rsa
ssh user looking-glass assign rsa-key "ssh-rsa AAAAB3..."
```

#### Cisco (IOS/IOS-XE)
```bash
# Create privilege level for looking-glass
privilege exec level 5 show ip route
privilege exec level 5 show ip bgp
privilege exec level 5 ping
privilege exec level 5 traceroute
privilege exec level 5 show version

# Create user with limited privileges
username looking-glass privilege 5
username looking-glass secret YourSecurePassword123!

# SSH key configuration
ip ssh rsa keypair-name looking-glass
username looking-glass
 ssh-key-type rsa
 ssh-key-data "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ..."
```

## ?? Supported Commands

| Command Type | Juniper | Huawei | Cisco | Description |
|-------------|---------|--------|-------|-------------|
| BGP Route | ? | ? | ? | Show specific BGP route |
| BGP Summary | ? | ? | ? | BGP session overview |
| BGP Neighbors | ? | ? | ? | BGP peer information |
| Advertised Routes | ? | ? | ? | Routes advertised to peer |
| Ping IPv4 | ? | ? | ? | ICMP ping test |
| Ping IPv6 | ? | ? | ? | ICMPv6 ping test |
| Traceroute IPv4 | ? | ? | ? | IPv4 path trace |
| Traceroute IPv6 | ? | ? | ? | IPv6 path trace |

## ?? Configuration Examples

### Juniper Configuration Example

```json
{
  "name": "Juniper Core Router",
  "display_name": "AMS-Core-01",
  "hostname": "ams-core-01.yournet.com",
  "ip": "192.168.1.10",
  "type": "juniper",
  "username": "${JUNIPER_USER}",
  "ssh_key": "/opt/looking-glass/keys/juniper_rsa",
  "port": 22,
  "timeout": 30,
  "enabled": true
}
```

### Huawei Configuration Example

```json
{
  "name": "Huawei Core Router",
  "display_name": "FRA-Core-01",
  "hostname": "fra-core-01.yournet.com",
  "ip": "192.168.1.20",
  "type": "huawei",
  "username": "${HUAWEI_USER}",
  "password": "${HUAWEI_PASS}",
  "port": 22,
  "timeout": 30,
  "enabled": true
}
```

### Cisco Configuration Example

```json
{
  "name": "Cisco Core Router",
  "display_name": "LON-Core-01", 
  "hostname": "lon-core-01.yournet.com",
  "ip": "192.168.1.30",
  "type": "cisco",
  "username": "${CISCO_USER}",
  "password": "${CISCO_PASS}",
  "port": 22,
  "timeout": 30,
  "enabled": true
}
```

## ?? SSH Key Authentication

### 1. Generate SSH Key Pair

```bash
# Generate SSH key for looking-glass service
sudo -u looking-glass ssh-keygen -t rsa -b 4096 -f /opt/looking-glass/keys/looking-glass_rsa -N ""

# Set proper permissions
sudo chmod 600 /opt/looking-glass/keys/looking-glass_rsa
sudo chmod 644 /opt/looking-glass/keys/looking-glass_rsa.pub
sudo chown looking-glass:looking-glass /opt/looking-glass/keys/*
```

### 2. Deploy Public Key to Routers

```bash
# Copy public key content
cat /opt/looking-glass/keys/looking-glass_rsa.pub
```

**Add to each router following vendor-specific instructions above.**

### 3. Update Configuration

```json
{
  "routers": [
    {
      "name": "Example Router",
      "username": "looking-glass",
      "ssh_key": "/opt/looking-glass/keys/looking-glass_rsa",
      "password": ""
    }
  ]
}
```

## ?? Testing Router Connectivity

### Manual SSH Test

```bash
# Test SSH connectivity manually
sudo -u looking-glass ssh -i /opt/looking-glass/keys/looking-glass_rsa -o ConnectTimeout=10 looking-glass@router1.yournet.com

# Test specific command
sudo -u looking-glass ssh -i /opt/looking-glass/keys/looking-glass_rsa looking-glass@router1.yournet.com "show version"
```

### API Test

```bash
# Test router list API
curl -s http://localhost:3002/api/routers | jq .

# Test command execution
curl -s -X POST http://localhost:3002/api/execute \
  -H "Content-Type: application/json" \
  -d '{
    "query": "ping",
    "addr": "8.8.8.8", 
    "router": "router1",
    "protocol": "IPv4"
  }' | jq .
```

## ?? Troubleshooting

### Common Issues

#### 1. SSH Connection Failed

**Solutions:**
```bash
# Check SSH service on router
ssh -v looking-glass@router.yournet.com

# Verify SSH key permissions
ls -la /opt/looking-glass/keys/
sudo -u looking-glass ssh-add -l

# Test network connectivity
sudo -u looking-glass ping -c 3 router.yournet.com
sudo -u looking-glass telnet router.yournet.com 22
```

#### 2. Authentication Failed

**Solutions:**
```bash
# Verify user exists on router
# Check SSH key is properly installed
# Verify username in configuration

# Test with verbose SSH output
sudo -u looking-glass ssh -vvv looking-glass@router.yournet.com
```

#### 3. Command Execution Failed

**Solutions:**
```bash
# Check user privileges on router
# Verify command syntax for router OS
# Test commands manually via SSH

# Check Looking Glass logs
sudo journalctl -u looking-glass.service -f
sudo tail -f /opt/looking-glass/logs/lg.log
```

## ?? Support

For router configuration support:
- ?? **Email**: [noc@goline.ch](mailto:noc@goline.ch)
- ?? **Documentation**: [GitHub Wiki](https://github.com/paolokappa/goline-looking-glass/wiki)
- ?? **Issues**: [GitHub Issues](https://github.com/paolokappa/goline-looking-glass/issues)

---

**Router Configuration Guide** | GOLINE Looking Glass  
*Paolo Caparrelli - GOLINE SA*
