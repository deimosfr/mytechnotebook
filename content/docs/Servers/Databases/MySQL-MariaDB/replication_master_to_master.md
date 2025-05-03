---
weight: 999
url: "/Replication_Master_to_Master/"
title: "Replication Master to Master"
description: "How to set up MySQL master-master replication for high-availability with step-by-step configuration guide and troubleshooting tips"
categories:
  - Debian
  - Database
  - Linux
date: "2013-06-21T07:25:00+02:00"
lastmod: "2013-06-21T07:25:00+02:00"
tags:
  - MySQL
  - Replication
  - High Availability
  - Database
  - Servers
  - Network
toc: true
---

![MySQL Logo](/images/mysql_logo.avif)

## Introduction

After completing the installation (See: [Installation and Configuration]({{< ref "docs/Servers/Databases/MySQL-MariaDB/mysql_installation_and_configuration.md" >}})), we are going to set up MySQL master-master replication. We need to replicate MySQL servers to achieve high-availability (HA). In my case I need two masters that are synchronized with each other so that if one of them goes down, the other could take over and no data is lost. Similarly, when the first one goes up again, it will still be used as slave for the live one.

![MySQL replication](/images/mysql_replication.avif)

Here is a basic step by step tutorial, that will cover the MySQL master and slave replication and also will describe the MySQL master and master replication.

**Notions:** we will call system 1 as master1 and slave2, and system2 as master2 and slave 1.

---

Install MySQL on master 1 and slave 1. Configure network services on both systems, like:

- **Master 1/Slave 2 IP: 10.8.0.1**
- **Master 2/Slave 1 IP: 10.8.0.6**

# Master 1 configuration

## Step 1

On Master 1, make changes in my.cnf:

```ini
...
[mysqld]
log-bin=/var/log/mysql/mysql-bin.log
binlog-do-db=<database name>  # input the database which should be replicated
binlog-ignore-db=mysql        # input the database that should be ignored for replication
binlog-ignore-db=test
expire_logs_days=14
# binlog_cache_size = 64K # Enable this if binlog_cache_disk_use increase in 'show global status'
sync_binlog = 1 # Reduce performances but essential for data integrity between slave <- master
slave_compressed_protocol = 1 # Good for wan replication to reduce network I/O at low cost of extra CPU

server-id=1
bind-address = 0.0.0.0
...
```

Notes: Italic lines are nearly optional and bold lines need to be set

If you want to sync all databases without setting exclusions, you could do without adding '\*ignore-db' values:

```ini
...
[mysqld]
log-bin=/var/log/mysql/mysql-bin.log
expire_logs_days=14
# binlog_cache_size = 64K # Enable this if binlog_cache_disk_use increase in 'show global status'
sync_binlog = 1 # Reduce performances but essential for data integrity between slave <- master
slave_compressed_protocol = 1 # Good for wan replication to reduce network I/O at low cost of extra CPU

server-id=1
bind-address = 0.0.0.0
...
```

Restart MySQL:

```bash
/etc/init.d/mysql restart
```

## Step 2

On master 1, connect to it and create a replication slave account in MySQL:

```bash
mysql> create user 'replication'@'10.8.0.1' identified by 'password';
```

- replication: username for replication
- password: password for replication user

```bash
mysql> grant replication slave on <database name>.* to 'replication'@'10.8.0.1';
```

or

```bash
mysql> grant replication slave on *.* to 'replication'@'10.8.0.1';
```

Flush the privileges and database tables with read lock:

```bash
flush privileges;
flush tables with read lock;
show master status;
+------------------+----------+--------------+------------------+
| File             | Position | Binlog_Do_DB | Binlog_Ignore_DB |
+------------------+----------+--------------+------------------+
| mysql-bin.000173 |    50937 | wikidb       |                  |
+------------------+----------+--------------+------------------+
1 row in set (0.00 sec)
```

Note, if you have full InnoDB databases, you don't need to flush tables with read locks.

## Step 3

Backup the desired database and upload them to the slave. If you're on InnoDB:

```bash
mysqldump -uroot -p --opt --add-drop-table --routines --triggers --events --single-transaction --master-data=2 -B wikidb > wikidb.sql
```

_Note: If there is InnoDB/MyISAM, add --lock-all-tables!_

```bash
scp wikidb.sql 10.8.0.6:~/
```

Then we will unlock tables:

```bash
mysql -uroot -p
unlock tables;
quit;
```

# Slave 1 configuration

## Step 1

Let's create the database:

```bash
mysql -uroot -p
create database wikidb;
quit;
```

Then we import the database:

```bash
mysql -uroot -p wikidb < wikidb.sql
```

## Step 2

Now edit my.cnf on Slave1 or Master2:

```bash
[mysqld]
...
server-id=2
bind-address = 0.0.0.0
# sync_binlog = 1
expire_log_days = 14
...
```

Then restart MySQL:

```bash
/etc/init.d/mysql restart
```

## Step 3

Grab information from the Master:

```bash
head -n100 database.sql
```

Now we're going to inform the slave database:

```bash
mysql -uroot -p
stop slave;
reset slave;
change master to master_host='10.8.0.1', master_user='replication', master_password='password', master_log_file='mysql-bin.000173', master_log_pos=50937;
```

Let's start the slave 1:

```bash
start slave;
show slave status\G;

************************** 1. row ***************************
             Slave_IO_State: Waiting for master to send event
                Master_Host: 192.168.16.4
                Master_User: replica
                Master_Port: 3306
              Connect_Retry: 60
            Master_Log_File: MASTERMYSQL01-bin.000009
        Read_Master_Log_Pos: 4
             Relay_Log_File: MASTERMYSQL02-relay-bin.000015
              Relay_Log_Pos: 3630
      Relay_Master_Log_File: MASTERMYSQL01-bin.000009
           Slave_IO_Running: Yes
          Slave_SQL_Running: Yes
            Replicate_Do_DB:
        Replicate_Ignore_DB:
         Replicate_Do_Table:
     Replicate_Ignore_Table:
    Replicate_Wild_Do_Table:
Replicate_Wild_Ignore_Table:
                 Last_Errno: 0
                 Last_Error:
               Skip_Counter: 0
        Exec_Master_Log_Pos: 4
            Relay_Log_Space: 3630
            Until_Condition: None
             Until_Log_File:
              Until_Log_Pos: 0
         Master_SSL_Allowed: No
         Master_SSL_CA_File:
         Master_SSL_CA_Path:
            Master_SSL_Cert:
          Master_SSL_Cipher:
             Master_SSL_Key:
      Seconds_Behind_Master: 1519187

1 row in set (0.00 sec)
```

The highlighted rows above must indicate related log files and **Slave_IO_Running and Slave_SQL_Running must be set to YES.**

- Seconds_Behind_Master: is the difference time between master and slave data sync (in progress...).
- Master_Host: This is the address of the master server. The slave will connect to it to get the logs that have to be replayed
- Master_User: The user that we will use for the replication
- Master_Port: The port where the slave will connect to on the master
- Master_Log_File: This file is the binary log that the slave will have to replay. The binary log is located on the master server. This file contains every statement that makes a change to the database in a binary format.
- Read_Master_Log_Pos: This is the position number on the binary log file that is being read by the "Slave_I/O Process" (the process that established the connection to the master and ensures the sync)
- Relay_Log_File: Same as the binary log but this file is located on the slave. This is a relay file of the master binary log
- Relay_Log_Pos: This is the position number on the relay log file that is being read by the "Slave_SQL Process" (the process that replays the changes to the slave database)
- Relay_Master_Log_File: It contains all the statements that have been processed by the slave.
- Slave_IO_Running: The status of the Slave_I/O process. When it's "Yes" that means that the slave is properly connected to the master. If it's "No" that means something is wrong with the connectivity between the 2 nodes
- Slave_SQL_Running: The status of the Slave_SQL process. When it's "Yes" that means that the slave is able to process statements and that all is working correctly. If it's "No" that means something is wrong while reading the relay logs

# Master 1 verification

```bash
mysql> show master status;
+------------------------+----------+--------------+------------------+
| File                   | Position | Binlog_Do_DB | Binlog_Ignore_DB |
+------------------------+----------+--------------+------------------+
|MysqlMYSQL01-bin.000008 |      410 | adam         |                  |
+------------------------+----------+--------------+------------------+
1 row in set (0.00 sec)
```

The above scenario is for master-slave, now we will create a slave master scenario for the same systems and it will work as master-master.

# Master 2 configuration

## Step 1

On Master 2, make changes in my.cnf:

```ini
...
[mysqld]
log-bin=/var/log/mysql/mysql-bin.log
binlog-do-db=<database name>  # input the database which should be replicated
binlog-ignore-db=mysql        # input the database that should be ignored for replication
binlog-ignore-db=test
expire_logs_days=14
# binlog_cache_size = 64K # Enable this if binlog_cache_disk_use increase in 'show global status'
sync_binlog = 1 # Reduce performances but essential for data integrity between slave <- master
slave_compressed_protocol = 1 # Good for wan replication to reduce network I/O at low cost of extra CPU

server-id=2
bind-address = 0.0.0.0
# sync_binlog = 1
expire_log_days = 14
...
```

Restart MySQL:

```bash
/etc/init.d/mysql restart
```

## Step 2

On master 1, connect to it and create a replication slave account in MySQL:

```bash
mysql> create user 'replication'@'10.8.0.6' identified by 'password';
```

- replication: username for replication
- password: password for replication user

```bash
mysql> grant replication slave on <database name>.* to 'replication'@'10.8.0.6';
```

or

```bash
mysql> grant replication slave on *.* to 'replication'@'10.8.0.6';
```

Flush the privileges:

```bash
flush privileges;
show master status;
+------------------+----------+--------------+------------------+
| File             | Position | Binlog_Do_DB | Binlog_Ignore_DB |
+------------------+----------+--------------+------------------+
| mysql-bin.000173 |    50937 | wikidb       |                  |
+------------------+----------+--------------+------------------+
1 row in set (0.00 sec)
```

# Slave 1 configuration

Now we're going to inform the slave database:

```bash
mysql -uroot -p
slave stop;
reset slave;
change master to master_host='10.8.0.6', master_user='replication', master_password='password', master_log_file='mysql-bin.000173', master_log_pos=50937;
```

Let's start the slave 2:

```bash
start slave;
show slave status\G;

************************** 1. row ***************************
             Slave_IO_State: Waiting for master to send event
                Master_Host: 10.8.0.6
                Master_User: replica
                Master_Port: 3306
              Connect_Retry: 60
            Master_Log_File: MASTERMYSQL01-bin.000009
        Read_Master_Log_Pos: 4
             Relay_Log_File: MASTERMYSQL02-relay-bin.000015
              Relay_Log_Pos: 3630
      Relay_Master_Log_File: MASTERMYSQL01-bin.000009
           Slave_IO_Running: Yes
          Slave_SQL_Running: Yes
            Replicate_Do_DB:
        Replicate_Ignore_DB:
         Replicate_Do_Table:
     Replicate_Ignore_Table:
    Replicate_Wild_Do_Table:
Replicate_Wild_Ignore_Table:
                 Last_Errno: 0
                 Last_Error:
               Skip_Counter: 0
        Exec_Master_Log_Pos: 4
            Relay_Log_Space: 3630
            Until_Condition: None
             Until_Log_File:
              Until_Log_Pos: 0
         Master_SSL_Allowed: No
         Master_SSL_CA_File:
         Master_SSL_CA_Path:
            Master_SSL_Cert:
          Master_SSL_Cipher:
             Master_SSL_Key:
      Seconds_Behind_Master: 1519187

1 row in set (0.00 sec)
```

The highlighted rows above must indicate related log files and **Slave_IO_Running and Slave_SQL_Running must be set to YES.**

# Master 2 verification

```bash
mysql> show master status;
+------------------------+----------+--------------+------------------+
| File                   | Position | Binlog_Do_DB | Binlog_Ignore_DB |
+------------------------+----------+--------------+------------------+
|MysqlMYSQL01-bin.000008 |      410 | adam         |                  |
+------------------------+----------+--------------+------------------+
1 row in set (0.00 sec)
```

The above scenario is for master-slave, now we will create a slave master scenario for the same systems and it will work as master-master.

# Monitoring

Now you need to monitor your replication. You need to look at this when running 'show slave status\G':

```bash
             ...
             Slave_IO_Running: Yes
             Slave_SQL_Running: Yes
             Seconds_Behind_Master: 0
             ...
```

What is interesting to investigate are these lines:

```bash
             ...
             Last_IO_Error: error connecting to master 'replication@192.168.25.9:3306' - retry-time: 60  retries: 86400
             Last_SQL_Error:
             ...
```

You also need to look at binlog size:

```bash
mysql> show global status like '%binlog%';
+------------------------+---------+
| Variable_name          | Value   |
+------------------------+---------+
| Binlog_cache_disk_use  | 2       |
| Binlog_cache_use       | 4333061 |
| Com_binlog             | 0       |
| Com_show_binlog_events | 0       |
| Com_show_binlogs       | 14      |
+------------------------+---------+
5 rows in set (0.00 sec)
```

Upgrade Binlog_cache_size from 32 to 128k if Binlog_cache_disk_use is not equal to 0.

# FAQ

## After a reboot the slave is not synchronized

On the master get information:

```bash
mysql -uroot -p
show master status;
```

Try to do this on the slave (adapt to your configuration):

```bash
mysql -uroot -p
slave stop;
reset slave;
change master to master_host='10.8.0.6', master_user='replication', master_password='password', master_log_file='mysql-bin.000173', master_log_pos=50937;
start slave;
show slave status\G;
```

The slave synchronization should be repaired :-)

## Watch current process to determine syncro problem

With this command, we can see if there is a MySQL replication problem:

```bash
$ mysqladmin -uroot -p processlist

+------+-------------+------------------------+----+-------------+--------+-----------------------------------------------------------------------+------------------+
| Id   | User        | Host                   | db | Command     | Time   | State                                                                 | Info             |
+------+-------------+------------------------+----+-------------+--------+-----------------------------------------------------------------------+------------------+
| 1    | system user |                        |    | Connect     | 264546 | Waiting for master to send event                                      |                  |
| 2    | system user |                        |    | Connect     | 1798   | Has read all relay log; waiting for the slave I/O thread to update it |                  |
| 4796 | replication | 10.8.0.6:60523         |    | Binlog Dump | 4733   | Has sent all binlog to slave; waiting for binlog to be updated        |                  |
| 5020 | root        | localhost              |    | Query       | 0      |                                                                       | show processlist |
+------+-------------+------------------------+----+-------------+--------+-----------------------------------------------------------------------+------------------+
```

## What if I have binlogs corrupted?

If you've got this kind of problem, a tool called mk-table-checksum available [here](https://www.maatkit.org/doc/mk-table-checksum.html) allows you to check consistency between 2 servers. It also can help to recreate a consistency database from the master.

## Repair replication

You may need to repair your replication because something went wrong and you can see that in the slave status (Last_Error). The way to correct the problem is to stop the slave:

```bash
stop slave;
```

Then analyze the error message and correct it on the slave. Then skip the error message:

```bash
set global sql_slave_skip_counter=1;
```

Then start the slave again:

```bash
start slave;
```

All should now be ok :-)

# Resources

Here is another documentation if you encounter problems:

- [MySQL Database Scale-out and Replication for High Growth Businesses (Very good documentation)](/pdf/mysql_database_scale-out_and_replication_for_high_growth_businesses.pdf)
- [Other MySQL masters replication documentation](/pdf/mysql_mastertomaster.pdf)
- [Master to Master replication MySQL 5 on Debian Etch](/pdf/mastertomaster_replication_with_mysql_5.pdf)
- [Master to Master replication MySQL 5 on Fedora](/pdf/master-master_replication_with_mysql_5_on_fedora_8.pdf)
- [How To Repair MySQL Replication](/pdf/how_to_repair_mysql_replication.pdf)
- [Setting Up Master-Master Replication On Four Nodes](/pdf/setting_up_master-master_replication_on_four_nodes_with_mysql_5_on_debian_etch.pdf)
- [MariaDB MySQL Advanced](/pdf/mariadb_mysql_avance.pdf)
