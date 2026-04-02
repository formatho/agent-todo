# Agent-Todo Scripts

This directory contains operational scripts for maintaining the Agent-Todo platform.

## Scripts

### backup-database.sh

Creates compressed, timestamped backups of the PostgreSQL database.

**Usage:**
```bash
./backup-database.sh
```

**Configuration:**
- Backup directory: `/home/deploy/backups`
- Retention: 30 days
- Compression: gzip

**Cron Setup (Daily at 2 AM):**
```bash
crontab -e
# Add this line:
0 2 * * * /home/deploy/todo/scripts/backup-database.sh >> /home/deploy/backups/cron.log 2>&1
```

### restore-database.sh

Restores the database from a backup file.

**Usage:**
```bash
./restore-database.sh /path/to/backup.sql.gz
```

## Backup Strategy

### Retention Policy
- **Daily:** Last 30 days
- **Weekly:** Last 4 weeks (TODO: implement)
- **Monthly:** Last 12 months (TODO: implement)

### Off-Site Storage (TODO)
- Configure S3/Backblaze B2 upload
- Encrypt backups before upload
- Verify uploads

### Monitoring (TODO)
- Set up backup success/failure alerts
- Monitor backup storage size
- Weekly backup verification

## Disaster Recovery

### Recovery Time Objective (RTO)
Target: < 1 hour

### Recovery Point Objective (RPO)
Target: < 24 hours (with daily backups)

### Recovery Procedure

1. **Stop the application:**
   ```bash
   cd /home/deploy/todo
   docker-compose stop backend
   ```

2. **Restore from backup:**
   ```bash
   ./scripts/restore-database.sh /path/to/backup.sql.gz
   ```

3. **Restart the application:**
   ```bash
   docker-compose start backend
   ```

4. **Verify:**
   - Check logs: `docker logs agent-todo-backend`
   - Test API: `curl https://todo.formatho.com/api/health`
   - Check data integrity

## Notes

- Always test backups regularly
- Document any changes to backup strategy
- Keep recovery documentation up to date
