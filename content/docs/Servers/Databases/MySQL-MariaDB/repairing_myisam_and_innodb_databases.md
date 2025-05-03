---
weight: 999
url: "/Réparer_des_bases_MyISAM_et_InnoDB/"
title: "Repairing MyISAM and InnoDB Databases"
description: "A comprehensive guide on how to repair and fix corrupted MyISAM and InnoDB databases in MySQL, with solutions for common database issues."
categories: ["Programming", "Storage", "Database"]
date: "2009-01-28T02:43:00+02:00"
lastmod: "2009-01-28T02:43:00+02:00"
tags: ["MySQL", "Database", "Recovery", "InnoDB", "MyISAM", "Repair"]
toc: true
---

## Introduction

So... your shiny MySQL database is no longer running and you want to fix it?

You've come to the right place!

I've assembled a list of 7 ways to fix your MySQL database when a simple restart doesn't do the trick, or when you have corrupt tables.

Simple MySQL restart:

```bash
/usr/local/mysql/bin/mysqladmin -uUSERNAME -pPASSWORD shutdown
/usr/local/mysql/bin/mysqld_safe &
```

## Solutions

### Corrupt MyISAM tables

MySQL database allows you to define a different MySQL storage engine for different tables. The storage engine is the engine used to store and retrieve data. Most popular storage engines are MyISAM and InnoDB.

MyISAM tables -will- get corrupted eventually. This is a fact of life.

Luckily, in most cases, MyISAM table corruption is easy to fix.

To fix a single table, connect to your MySQL database and issue a:

```bash
repair TABLENAME
```

To fix everything, go with:

```bash
/usr/local/mysql/bin/mysqlcheck --all-databases -uUSERNAME -pPASSWORD -r
```

A lot of times, MyISAM tables will get corrupt and you won't even know about it unless you review the log files.

I highly suggest you add this line to your `/etc/my.cnf` config file. It will automatically fix MyISAM tables as soon as they become corrupt:

```bash
...
[mysqld]
myisam-recover=backup,force
...
```

If this doesn't help, there are a few additional tricks you can try.

### Multiple instances of MySQL

This is pretty common. You restart MySQL and the process immediately dies.

Reviewing the log files will tell you another instance of MySQL may be running.

To stop all instances of MySQL:

```bash
/usr/local/mysql/bin/mysqladmin -uUSERNAME -pPASSWORD shutdown
killall mysql
killall mysqld
```

Now you can restart the database and you will have a single running instance

### Changed InnoDB log settings

Once you have a running InnoDB MySQL database, you should never ever change these lines in your `/etc/my.cnf` file:

```ini
datadir = /usr/local/mysql/data
innodb_data_home_dir = /usr/local/mysql/data
innodb_data_file_path = ibdata1:10M:autoextend
innodb_log_group_home_dir = /usr/local/mysql/data
innodb_log_files_in_group = 2
innodb_log_file_size = 5242880
```

InnoDB log file size cannot be changed once it has been established. If you change it, the database will refuse to start.

### Disappearing MySQL host tables

I've seen this happen a few times. Probably some kind of freakish MyISAM bug.

Easily fixed with:

```bash
/usr/local/bin/mysql_install_db
```

### MyISAM bad auto_increment

If the auto_increment count goes haywire on a MyISAM table, you will no longer be able to INSERT new records into that table.

You can typically tell the auto_increment counter is malfunctioning, by seeing an auto_increment of -1 assigned to the last inserted record.

To fix - find the last valid auto_increment id by issuing something like:

```bash
SELECT max(id) from tablename
```

And then update the auto_increment counter for that table

```bash
ALTER TABLE tablename AUTO_INCREMENT = id+1
```

### Too many connections

Your database is getting hit with more connections than it can handle and now you cannot even connect to the database yourself.

First, stop the database:

```bash
/usr/local/mysql/bin/mysqladmin -uUSERNAME -pPASSWORD shutdown
```

If that doesn't help you can try "killall mysql" and "killall mysqld"

Once the database stopped, edit your `/etc/my.cnf` file and increase the number of connections. Don't go crazy with this number or you'll bring your entire machine down.

On a dedicated database machine we typically use:

```text
max_connections = 200
wait_timeout = 100
```

Try restarting the database and see if that helps.

If you're getting bombarded with queries and you need to be able to connect to the database to make some table changes, set a different port number in your `/etc/my.cnf` file, start the database, make any changes, then update the port back to normal (master-port = 3306) and restart.

### Corrupt InnoDB tables

InnoDB tables are my favorite. Transactional, reliable and unlike MyISAM, InnoDB supports concurrent writes into the same table.

InnoDB's internal recovery mechanism is pretty good. If the database crashes, InnoDB will attempt to fix everything by running the log file from the last timestamp. In most cases it will succeed and the entire process is transparent.

Unfortunately if InnoDB fails to repair itself, the -entire- database will not start. MySQL will exit with an error message and your entire database will be offline. You can try to restart the database again and again, but if the repair process fails - the database will refuse to start.

This is one reason why you should always run a master/master setup when using InnoDB - have a redundant master if one fails to start.

Before you go any further, review MySQL log file and confirm the database is not starting due to InnoDB corruption.

There are tricks to update InnoDB's internal log counter so that it skips the queries causing the crash, but in our experience this is not a good idea. You lose data consistency and will often break replication.

Once you have corrupt InnoDB tables that are preventing your database from starting, you should follow this five step process:

- Add this line to your `/etc/my.cnf` configuration file:

```text
 ...
 [mysqld]
 innodb_force_recovery = 4
 ...
```

- Restart MySQL. Your database will now start, but with innodb_force_recovery, all INSERTs and UPDATEs will be ignored.
- Dump all tables
- Shutdown database and delete the data directory. Run mysql_install_db to create MySQL default tables
- Remove the innodb_force_recovery line from your `/etc/my.cnf` file and restart the database. (It should start normally now)
- Restore everything from your backup

### InnoDB Corruption

Recently I was faced with the daunting task of reparing an InnoDB database gone bad. The database would not start due to corruption.

First step was turning-on InnoDB force-recovery mode, where InnoDB starts but ignores all UPDATEs and INSERTs.

Add this line to `/etc/my.cnf`:

```bash
innodb_force_recovery = 2
```

Now we can restart the database:

```bash
/usr/local/bin/mysqld_safe &
```

(Note: If MySQL doesn't restart, keep increasing the innodb_force_recovery number until you get to innodb_force_recovery = 8)

Save all data into a temporary alldb.sql (this next command can take a while to finish):

```bash
mysqldump --force --compress --triggers --routines --create-options -uUSERNAME -pPASSWORD --all-databases > /usr/alldb.sql
```

Shutdown the database again:

```bash
mysqladmin -uUSERNAME -pPASSWORD shutdown
```

Delete the database directory. (Note: In my case the data was under /usr/local/var. Your setup may be different. Make sure you're deleting the correct directory)

```bash
rm -Rfd /usr/local/var
```

Recreate the database directory and install MySQL basic tables

```bash
mkdir /usr/local/var
chown -R mysql:mysql /usr/local/var
/usr/local/bin/mysql_install_db
chown -R mysql:mysql /usr/local/var
```

Remove innodb_force_recovery from `/etc/my.cnf` and restart database:

```bash
/usr/local/bin/mysqld_safe &
```

Import all the data back (this next command can take a while to finish):

```bash
mysql -uroot --compress < /usr/alldb.sql
```

And finally - flush MySQL privileges (because we're also updating the MySQL table)

```bash
/usr/local/bin/mysqladmin -uroot flush-privileges
```

Note: For best results, add port=8819 (or any other random number) to `/etc/my.cnf` before restarting MySQL and then add --port=8819 to the mysqldump command. This way you avoid the MySQL database getting hit with queries while the repair is in progress.

## Resources
- http://www.softwareprojects.com/resources/programming/t-how-to-fix-mysql-database-myisam-innodb-1634.html
