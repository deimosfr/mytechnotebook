---
weight: 999
url: "/Sauvegardes,_restaurations_et_transferts/"
title: "MySQL: Backups, Restorations and Transfers"
description: "Guide to backup, restore, and transfer MySQL databases with different methods including mysqldump, mysqlhotcopy, and LVM snapshots."
categories: ["MySQL", "Database", "Linux"]
date: "2013-11-15T12:46:00+02:00"
lastmod: "2013-11-15T12:46:00+02:00"
tags: ["MySQL", "Backup", "Restore", "LVM", "mysqldump", "Shell Scripts"]
toc: true
---

![MySQL Logo](/images/mysql_logo.avif)

## Backups

### mysqldump

In my opinion, this is THE best method. It's provided as standard and works very well. **However, be careful with large databases**. If you have large databases, then look at [XtraBackup](./xtrabackup_:_optimiser_ses_backups_mysql.html).

Its use with a shell is as follows (for small databases):

```bash
mysqldump -uUSER -pPASSWORD -e --single-transaction --opt DATABASE > backup-`date +%y%m%d`.sql
```

To back up all databases at once:

```bash
mysqldump -uUSER -pPASSWORD --all-databases > backup-`date +%y%m%d`.sql
```

There are also other options:

- --no-data: allows you to generate a file that contains all the instructions for creating your database tables (CREATE DATABASE, CREATE TABLE...), but not the data itself.
- --no-createdb and --no-create-info allow you to automatically comment out schema creation instructions. So you can insert data from a backup into a database with existing tables.
- **--lock-tables: allows you to place a LOCK on an entire database. This gives consistency between tables, but not between databases.**
- **--lock-all-tables: allows you to place a LOCK on all tables. Be aware of the impacts of such a LOCK on a large non-transactional database like MyISAM.** Putting all queries on hold during the backup could make visitors think that the site's server has crashed. At the end of the dump, all pending queries will be executed (no losses, unless the queue length has been exceeded).
- -e: allows you to accumulate creation of elements of the same type (so will be faster on reimport)
- **--single-transaction: allows you to run backups in a transaction. It therefore guarantees that InnoDB tables (transactional engine) will be in a consistent state between them. Tables that do not support transactions will also be backed up, but consistency will not be guaranteed (the BEGIN instruction is simply ignored for these tables).** The advantage of this option is that it does not use a LOCK. In fact, they will be locked inside the transaction, without preventing other clients from making changes to the database.
- --master-data: allows you to have the master's positions for replication directly in the dump.
- --opt: adds drop-tables if they exist.

#### Quickly backup each database

To back up each database simply and quickly, here is a solution:

```bash
for I in `echo "show databases;" | mysql -uroot -p | grep -v Database`; do mysqldump -uroot -p -e --single-transaction --opt $I > "$I.sql"; done
```

You'll need to hardcode the password (create a backup user with only backup rights for improved security).

### mysqlhotcopy

This command is a PERL script delivered with MySQL that essentially performs a raw file copy. However, this method is quite efficient.

This command allows you to do:

- cp: if it's for a local backup
- scp: if it's for remote

It places lock files on the tables to be backed up. There is also the **--record_log_pos option which allows you to record the position in the binary logs when the server is in master or slave mode (to quickly create a new slave).** A few small examples for a local backup:

```bash
mysqlhotcopy user password /var/lib/mysql/ma_base
```

And for a remote backup:

```bash
mysqlhotcopy --user=user --password=pass user user@host:/home/mon_backup
```

A small disadvantage, but one that is important: **mysqlhostcopy only works with MyISAM and ARCHIVE**. There is a paid tool (HOT Backup) that allows you to do the same thing but with InnoDB as well.

### LVM Snapshot

I won't go into how to use LVM here, but it's the method to use for large databases (16 GB for example).

First, we need to place a lock (yes, we're still obliged to, but it only lasts a few hundredths of a second):

```bash
flush tables with read lock;
```

Then we create our snapshot:

```bash
lvcreate --snapshot -n snap -L 16G /dev/volumegroupe/snapshot_mysql
```

Since the operation is instantaneous, we can remove the lock:

```bash
mysql> unlock tables;
```

## Restorations

Restoration is done by directly using the main program:

```bash
mysql -u<user> -p<password> <database_name> < backup-`date +%y%m%d`.sql
```

If you want to restore only one database from a dump where there are all databases:

```bash
mysql -u<user> -p<password> --one-database <database_name> < backup-`date +%y%m%d`.sql
```

## Transfers

Finally, if the goal is to transfer the database from one machine to another, you can combine the two calls on a single line:

```bash
mysqldump -uUSER -pPASSWORD DATABASE | ssh user:password@IP "cat > myfile.sql"
```

## Scripts

### Backup by mail

Here is a small script to adapt according to your needs:

```bash
#!/bin/sh
## MySQL Backup Script
## Made by Pierre Mavro

## MYSQL INFO ##
LOGIN=root # Login for mysql
PASS=toto # Pass for mysql

## Email ##
ADMINMAIL=xxx@mycompany.com # This is the mail where the saves will be sent

## Vars ##
TMPDIR=/tmp/baksql # Temp directory

## DO NOT MODIFY NOW ##
mkdir $TMPDIR
# Backuping databases
for databases in "cacti" "information_schema" "mysql" "snort" "wikidb" ; do
       mysqldump -u$LOGIN -p$PASS $databases > $TMPDIR/$databases-`date +%y%m%d`.sql && baksql=$baksql`echo "Backup of database $databases - OK ; "` || baksql=$baksql`echo "Backup of database $databases - FAILED ; "`
done

# Compressing and emailing
tar -czf $TMPDIR/mysql_backup.tgz $TMPDIR/*.sql && echo "`echo $baksql`" | mutt -x -a $TMPDIR/mysql_backup.tgz -s "MySQL backup - `date +%d\-%m\-%Y` - fire" $ADMINMAIL
rm -Rf $TMPDIR
```

### Backup and compression of all databases

This more tedious method has the advantage of backing up database by database, which allows you to restore only the database you're interested in if there's a problem.

Additionally, it includes on-the-fly compression of your database. However, you'll need to install 7zip beforehand (I chose 7zip for better compression).

Let's create a script that we'll place in `/etc/scripts` for example:

```bash
#!/bin/bash
user='root'
password='password'
destination='/tmp/backups_sql'
mail='xxx@mycompany.com'

mkdir -p $destination
for i in `echo "show databases;" | mysql -u$user -p$password | grep -v Database`; do
       mysqldump -u$user -p$password --opt --add-drop-table --routines --triggers --events --single-transaction --master-data=2 -B $i | 7z a -t7z -mx=9 -si $destination/$i.sql.7z
done

problem_text=''
problem=0
for i in `ls $destination/*` ; do
    size=`du -sk $i | awk '{ print $1 }'`
    if [ $size -le 4 ] ; then
        problem_text="$problem_text- $i database. Backupped database size is equal or under 4k ($size)\n"
        problem=1
    fi
done

if [ $problem -ne 0 ] ; then
    echo -e "Backups problem detected on:\n\n$problem_text" | mail -s "$HOSTNAME - MySQL backup problem" $mail
fi
```

The problem here is the clear text password. So make sure to properly restrict access to the user who will be backing up:

```bash
chmod 700 /etc/scripts/backup_mysql_databases.sh
```

### Transfers of a MySQL database to another via SSH

```bash
#!/bin/sh
## MySQL Backup Script
## Made by Pierre Mavro

## MYSQL INFO ##
LOGIN=root # Login for mysql
PASS=toto # Pass for mysql

## SSH SERVER FOR TRANSFER & REINJECTION ##
MYSQLSRV2="192.168.0.1"

## Email ##
ADMINMAIL=xxx@mycompany.com # This is the mail where transfer failures will be notified

## Vars ##
TMPDIR=/tmp/baksql # Temp directory

## DO NOT MODIFY NOW ##
mkdir $TMPDIR
# Backuping databases
for databases in "bugzilla" "cacti" "mysql" "networkdb" "wiki" ; do
       mysqldump -u$LOGIN -p$PASS $databases > $TMPDIR/$databases-`date +%y%m%d`.sql && baksql=$baksql`echo "Backup of database $databases - OK ; "` || baksql=$baksql`echo "Backup of database $databases - FAILED ; "`
done

# Transfering databases to the other SQL server (need ssh key)
tar -czf $TMPDIR/mysql_backup.tgz $TMPDIR/*.sql
scp $TMPDIR/mysql_backup.tgz $MYSQLSRV2:~/
ssh $MYSQLSRV2 tar -xzvf ~/mysql_backup.tgz
for databases in "bugzilla" "cacti" "mysql" "networkdb" "wiki" ; do
        ssh $MYSQLSRV2 "mysql -u$LOGIN -p$PASS $databases < ~/tmp/baksql/$databases-`date +%y%m%d`.sql" && baksql=$baksql`echo "Reinjection of database $databases - OK ; "` || baksql=$baksql`echo "Reinjection of database $databases - FAILED ; "`
done

# Delete temps
ssh $MYSQLSRV2 rm -Rf /root/tmp/baksql ~/mysql_backup.tgz ~/tmp
rm -Rf $TMPDIR

# Send mail
echo $baksql | mail -s "Mysql transfers" $ADMINMAIL
```

## Resources
- [How To Back Up MySQL Databases With mylvmbackup](/pdf/how_to_back_up_mysql_databases_with_mylvmbackup.pdf)
