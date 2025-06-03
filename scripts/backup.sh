#!/bin/bash
BACKUP_DIR="/opt/goline-looking-glass/backups/$(date +%Y%m%d_%H%M%S)"
mkdir -p $BACKUP_DIR
cp config.json index.html $BACKUP_DIR/
cp -r logs/ $BACKUP_DIR/ 2>/dev/null
tar -czf $BACKUP_DIR.tar.gz -C $BACKUP_DIR .
rm -rf $BACKUP_DIR
echo "Backup created: $BACKUP_DIR.tar.gz"
