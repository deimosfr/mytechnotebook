---
weight: 999
url: "/ZRM_\\:_Sauvegardes_automatisées_et_restaurations_faciles/"
title: "ZRM: Automated Backups and Easy Restorations"
description: "A guide to using ZRM for MySQL for automated database backups and point-in-time recovery for all MySQL storage engines"
categories: ["Backup", "Linux", "Database"]
date: "2008-04-19T11:16:00+02:00"
lastmod: "2008-04-19T11:16:00+02:00"
tags: ["MySQL", "Backup", "Recovery", "Database", "ZRM"]
toc: true
---

## Introduction

ZRM for MySQL is a powerful, flexible and robust backup and recovery solution for MySQL databases for all storage engines. With ZRM for MySQL a Database Administrator can automate logical or raw backup to a local or remote disk. In this How To, we attempt to explain how to recover from an user error at any given point in time.

## Our Scenario

At approximately 2:30pm there were 5 tablespaces added into the MovieID table. At 3pm you get a call from a user, who was trying to delete some unused tablespaces but ended up deleting the last 5 that he just added. The last full backup you performed was the night before at 7pm. How do you recover back to right before the last 5 tablespaces were deleted? In this case we will demonstrate a point in time restore.

### Note

* More information on ZRM for MySQL is available [here](https://www.zmanda.com/backup-mysql.html).
* To know more about configuring ZRM for MySQL you can look at [https://www.zmanda.com/quick-mysql-backup.html](https://www.zmanda.com/quick-mysql-backup.html).
* For this How To, we use [ZRM for MySQL version 1.1.1](https://www.zmanda.com/downloads.html#ZRM).
* To perform incremental backups the binary logging option must be turned on. Edit the `/etc/my.cnf` file and add the following line under the [mysqld] section:

```bash
log-bin=/var/lib/mysql/mysql-bin.log
```

You will have restart your MySQL daemon before this goes into effect.

### Verifications

We will verify that the last full backup ran successfully last night:

```bash
$ mysql-zrm-reporter -show restore-info --where backup-set=dailyrun

backup_set backup_date backup_level backup_directory
----------------------------------------------------------------------------------------------------------
dailyrun Wed 18 Oct 2006 07:07:08 PM PDT 0 /var/lib/mysql-zrm/dailyrun/20061018190708
```

We can see above that the last full did run successfully last night at 7:07pm.

### Backup

So now we will run an incremental backup manually to record all the changes between the last backup and now, by typing:

```bash
mysql-zrm-scheduler --now --backup-set dailyrun --backup-level 1
```

### Parsing logs

Next we will parse through the binary logs backed up in the last incremental backup

```bash
$ mysql-zrm --action parse-binlogs --source-directory /var/lib/mysql-zrm/dailyrun/20061019151937 --backup-set dailyrun
 
------------------------------------------------------------
Log filename | Log Position | Timestamp | Event Type | Event
------------------------------------------------------------
/var/lib/mysql-zrm/dailyrun/20061019151937/mysql-bin.000002 | 4 | 06-11-19 14:09:58 | Start: binlog v 4, server v 5.0.22-log created 061019 14:09:58 |
/var/lib/mysql-zrm/dailyrun/20061019151937/mysql-bin.000002 | 98 | 06-11-19 14:34:27 | Query | use movies; INSERT INTO `MovieID` (`MovieID`, `Year`, `MovieTitle`) VALUES ('17786', '1999', 'Sopranos: Season 1 Disc 1');
/var/lib/mysql-zrm/dailyrun/20061019151937/mysql-bin.000002 | 272 | 06-11-19 14:35:46 | Query | INSERT INTO `MovieID` (`MovieID`, `Year`, `MovieTitle`) VALUES ('17787', '1999', 'Sopranos: Season 1 Disc 2');
/var/lib/mysql-zrm/dailyrun/20061019151937/mysql-bin.000002 | 446 | 06-11-19 14:36:02 | Query | INSERT INTO `MovieID` (`MovieID`, `Year`, `MovieTitle`) VALUES ('17788', '1999', 'Sopranos: Season 1 Disc 3');
/var/lib/mysql-zrm/dailyrun/20061019151937/mysql-bin.000002 | 620 | 06-11-19 14:36:36 | Query | INSERT INTO `MovieID` (`MovieID`, `Year`, `MovieTitle`) VALUES ('17789', '1999', 'Sopranos: Season 1 Disc 4');
/var/lib/mysql-zrm/dailyrun/20061019151937/mysql-bin.000002 | 794 | 06-11-19 14:36:53 | Query | INSERT INTO `MovieID` (`MovieID`, `Year`, `MovieTitle`) VALUES ('17790', '1999', 'Sopranos: Season 1 Disc 5');
/var/lib/mysql-zrm/dailyrun/20061019151937/mysql-bin.000002 | 968 | 06-11-19 14:56:15 | Query | DELETE FROM `MovieID` WHERE `MovieID`.`MovieID` = 17786 LIMIT 1;
/var/lib/mysql-zrm/dailyrun/20061019151937/mysql-bin.000002 | 1096 | 06-11-19 14:56:15 | Query | DELETE FROM `MovieID` WHERE `MovieID`.`MovieID` = 17787 LIMIT 1;
/var/lib/mysql-zrm/dailyrun/20061019151937/mysql-bin.000002 | 1224 | 06-11-19 14:56:15 | Query | DELETE FROM `MovieID` WHERE `MovieID`.`MovieID` = 17788 LIMIT 1;
/var/lib/mysql-zrm/dailyrun/20061019151937/mysql-bin.000002 | 1352 | 06-11-19 14:56:15 | Query | DELETE FROM `MovieID` WHERE `MovieID`.`MovieID` = 17789 LIMIT 1;
/var/lib/mysql-zrm/dailyrun/20061019151937/mysql-bin.000002 | 1480 | 06-11-19 14:56:15 | Query | DELETE FROM `MovieID` WHERE `MovieID`.`MovieID` = 17790 LIMIT 1;
/var/lib/mysql-zrm/dailyrun/20061019151937/mysql-bin.000002 | 1608 | 06-11-19 15:19:37 | Rotate to mysql-bin.000003 pos: 4 |
------------------------------------------------------------
INFO: Removing all of the uncompressed/unencrypted data
```

### Restoration

So now we will restore the database to what it looked like at approximately 2:45pm. Since the tables were added at 2:30pm and accidentally deleted at 3pm. Since we want the database back to the state it was at right before the delete.

```bash
$ mysql-zrm --action restore --source-directory /var/lib/mysql-zrm/dailyrun/20061019151937 --backup-set dailyrun --stop-datetime "20061019144500"
 
INFO: ZRM for MySQL Community Edition - version 1.1
INFO: Mail address: dba@zmanda.com is ok
INFO: Input Parameters Used {
INFO: verbose=1
INFO: retention-policy=10D
INFO: backup-level=1
INFO: mailto=dba@zmanda.com
INFO: databases=movies
INFO: source-directory=/var/lib/mysql-zrm/dailyrun/20061019151937
INFO: html-reports=backup-status-info
INFO: password=******
INFO: backup-mode=logical
INFO: compress-plugin=/usr/bin/gzip
INFO: compress=/usr/bin/gzip
INFO: user=backup-user
INFO: stop-datetime=20061019144500
INFO: Getting mysql variables
INFO: mysqladmin --user=backup-user --password=***** variables
INFO: datadir is /var/lib/mysql/
INFO: mysql_version is 5.0.22-log
INFO: log_bin=ON
INFO: Uncompressing backup
INFO: Command used is 'cat "/var/lib/mysql-zrm/dailyrun/20061019151937/backup-data" | "/usr/bin/gzip" -d | tar --same-owner -xpsC "/var/lib/mysql-zrm/dailyrun/20061019151937" 2>/tmp/HId0KZkvcS'
INFO: Restoring incremental to tmpfile
INFO: mysqlbinlog --user=backup-user --password=***** --stop-datetime=20061019144500 --database=movies -r /tmp/NNqSZFZa8R "/var/lib/mysql-zrm/dailyrun/20061019151937"/mysql-bin.[0-9]*
INFO: restoring using command mysql --user=backup-user --password=***** -e "source /tmp/NNqSZFZa8R;"
INFO: Incremental restore done for database movies
INFO: Shutting down MySQL
INFO: Removing all of the uncompressed/unencrypted data
INFO: Restore done in 2 seconds. MySQL server has been shutdown. Please restart after verification.
```

Once you restart the MySQL services you'll notice that the database has been restored to what it looked at 2:45pm. Which means it has the last 5 tablespaces that were accidentally deleted at 3pm.

## Resources
- [MySQL Backups Using ZRM For MySQL 2.0](/pdf/mysql_backups_using_zrm_for_mysql_2_0.pdf)
