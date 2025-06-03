# Changelog

All notable changes to the GOLINE Looking Glass project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release preparation
- Multi-vendor router support architecture
- RESTful API design
- Responsive web interface framework

## [1.0.0] - 2024-06-03

### Added
- Initial public release of GOLINE Looking Glass
- Multi-vendor router support (Juniper, Huawei, Cisco)
- RESTful API for network commands
- Responsive web interface with navy blue theme
- Docker containerization support
- Comprehensive installation and configuration documentation

### Features
- BGP route lookups and analysis
- Ping and traceroute functionality  
- IPv4 and IPv6 dual-stack support
- SSH key authentication for secure router access
- Rate limiting and security features
- Customizable company branding and themes
- Real-time command execution with progress indicators
- Mobile-responsive design
- Professional Apache reverse proxy setup

### Security
- SSH key-based authentication system
- Input validation and sanitization
- Rate limiting implementation
- Secure command execution with timeouts
- User privilege separation for system service
- SSL/TLS encryption support with Let's Encrypt

### Documentation
- Complete installation guide with automated scripts
- Router configuration examples for all supported vendors
- Security best practices and hardening guide
- Troubleshooting guide with common solutions
- API reference documentation
- Contributing guidelines for developers
- Docker deployment documentation

### Infrastructure
- Systemd service integration
- Log rotation configuration
- Apache virtual host templates
- SSL certificate automation
- Backup and restore procedures

---

**Author:** Paolo Caparrelli (GOLINE SA - AS202032)  
**Contact:** noc@goline.ch  
**Project:** https://github.com/paolokappa/goline-looking-glass
