---
weight: 999
url: "/Purger_les_mysql-bin_logs/"
title: "Purging MySQL Binary Logs"
description: "This guide explains how to manage slow query logs and binary logs in MySQL, including how to flush, rotate, and purge logs to free up disk space."
categories: ["MySQL", "Database", "Linux"]
date: "2011-12-02T08:06:00+02:00"
lastmod: "2011-12-02T08:06:00+02:00"
tags: ["MySQL", "Binary logs", "Maintenance", "Database administration", "Replication"]
toc: true
---

## Introduction

Slow query logs and binary logs are essential for several purposes in MySQL. In this documentation, we'll see how to manage them. This is particularly useful when you're running low on disk space.

## Slow Query Logs

### Introduction

Slow query logs are simply standard logs. They can be managed like standard server logs.

### Flush Logs

To flush logs, simply run this command:

```bash
echo "" > slow-query.log
```

### Rotate Logs

```bash
mv slow-query.log slow-query.log.old
mysqladmin flush-logs
```

## Binary Logs

### Introduction

You might find a large number of bin files in the MySQL data directory called "server-bin.n" or mysql-bin.00000n, where n is an incrementing number. Usually `/var/lib/mysql` stores the binary log files. The binary log contains all statements that update data or potentially could have updated it. For example, a DELETE or UPDATE which matched no rows. Statements are stored in the form of events that describe the modifications. The binary log also contains information about how long each statement that updated data took to execute.

**Warning: do not delete these files manually or you won't be able to restart your database!**

You'll get this kind of error messages:

```
061031 17:38:48  mysqld started
061031 17:38:48  InnoDB: Started; log sequence number 14 1645228884
/usr/libexec/mysqld: File '/var/lib/mysql/mysql-bin.000017' not found
(Errcode: 2)
061031 17:38:48 [ERROR] Failed to open log (file
'/var/lib/mysql/mysql-bin.000017', errno 2)
061031 17:38:48 [ERROR] Could not open log file
061031 17:38:48 [ERROR] Can't init tc log
061031 17:38:48 [ERROR] Aborting

061031 17:38:48  InnoDB: Starting shutdown...

061031 17:38:51  InnoDB: Shutdown completed; log sequence number 14 1645228884

061031 17:38:51 [Note] /usr/libexec/mysqld: Shutdown complete
061031 17:38:51  mysqld ended
```

### The Purpose of MySQL Binary Log

The binary log has two important purposes:

* **Data Recovery**: It may be used for data recovery operations. After a backup file has been restored, the events in the binary log that were recorded after the backup was made are re-executed. These events bring databases up to date from the point of the backup.
* **High availability / replication**: The binary log is used on master replication servers as a record of the statements to be sent to slave servers. The master server sends the events contained in its binary log to its slaves, which execute those events to make the same data changes that were made on the master.

#### Disable MySQL Binlogging

**If you are not replicating**, you can disable binlogging by changing your my.ini or my.cnf file. Open your my.ini or `/etc/my.cnf` (`/etc/mysql/my.cnf`), find a line that reads "log_bin" and remove or comment the following lines:

```bash
#log_bin                        = /var/log/mysql/mysql-bin.log
#expire_logs_days        = 10
#max_binlog_size         = 100M
```

Now restart MySQL.

#### Purge Master Logs

**If you ARE replicating**, then you need to periodically RESET MASTER or PURGE MASTER LOGS to clear out the old logs as those files are necessary for the proper operation of replication. Use the following command to purge master logs:

```bash
mysql -u root -p 'MyPassword' -e "PURGE MASTER LOGS TO 'mysql-bin.03';"
```

or

```bash
mysql -u root -p 'MyPassword' -e "PURGE MASTER LOGS BEFORE '2008-12-15 10:06:06';"
```

Note: I recommend using the penultimate MySQL master log.

You may want MySQL to do it automatically. For that, change the configuration:

```bash
expire_logs_days = 14
```

Then restart the MySQL server.

## Resources
- http://www.cyberciti.biz/faq/what-is-mysql-binary-log/
- http://dev.mysql.com/doc/refman/5.0/en/binary-log.html
- http://forums.mysql.com/read.php?10,78659,78660#msg-78660
