---
weight: 999
url: "/MysqlTuner_\\:_Optimiser_votre_serveur_MySQL/"
title: "MysqlTuner: Optimizing Your MySQL Server"
description: "A guide on how to install and use MysqlTuner to optimize MySQL server performance, including cache management, query analysis, and table optimization techniques."
categories: ["Database", "MySQL", "Storage", "Performance", "Debian"]
date: "2013-06-06T09:19:00+02:00"
lastmod: "2013-06-06T09:19:00+02:00"
tags:
  [
    "MySQL",
    "Performance",
    "Database Optimization",
    "Tuning",
    "mysqltuner",
    "Patch",
  ]
toc: true
---

![MySQL Logo](/images/mysql_logo.avif)

## Introduction

MySQL is widely used in business environments but can become very complex when managing databases of significant size. We need to check cache, memory usage, and look for optimization opportunities in various areas.

This documentation will not only cover MysqlTuner, but also how to optimize your MySQL server in general.

## Installation

Installing MysqlTuner is fairly simple:

```bash
apt-get install mysqltuner
```

## Usage

Simply launch mysqltuner with the correct connection rights to your database:

```bash
 >>  MySQLTuner 1.0.1 - Major Hayden <major@mhtx.net>
 >>  Bug reports, feature requests, and downloads at http://mysqltuner.com/
 >>  Run with '--help' for additional options and output filtering
Please enter your MySQL administrative login: root
Please enter your MySQL administrative password:

-------- General Statistics --------------------------------------------------
[--] Skipped version check for MySQLTuner script
[OK] Currently running supported MySQL version 5.1.39-log
[OK] Operating on 64-bit architecture

-------- Storage Engine Statistics -------------------------------------------
[--] Status: -Archive -BDB -Federated +InnoDB -ISAM -NDBCluster
[--] Data in MyISAM tables: 2G (Tables: 35)
[--] Data in InnoDB tables: 7M (Tables: 41)
[!!] Total fragmented tables: 43

-------- Performance Metrics -------------------------------------------------
[--] Up for: 16h 38m 24s (5K q [0.089 qps], 22 conn, TX: 20M, RX: 3M)
[--] Reads / Writes: 65% / 35%
[--] Total buffers: 812.0M global + 1.2M per thread (151 max threads)
[OK] Maximum possible memory usage: 991.3M (1% of installed RAM)
[!!] Slow queries: 29% (1K/5K)
[OK] Highest usage of available connections: 1% (3/151)
[OK] Key buffer size / total MyISAM indexes: 128.0M/167.5M
[OK] Key buffer hit rate: 99.1% (59K cached / 557 reads)
[!!] Query cache efficiency: 0.1% (4 cached / 3K selects)
[OK] Query cache prunes per day: 0
[OK] Sorts requiring temporary tables: 0% (0 temp sorts / 6 sorts)
[OK] Temporary tables created on disk: 9% (8 on disk / 87 total)
[OK] Thread cache hit rate: 86% (3 created / 22 connections)
[!!] Table cache hit rate: 17% (64 open / 376 opened)
[OK] Open file limit used: 16% (121/755)
[OK] Table locks acquired immediately: 100% (4K immediate / 4K locks)
[!!] Connections aborted: 13%
[OK] InnoDB data size / buffer pool: 7.9M/512.0M

-------- Recommendations -----------------------------------------------------
General recommendations:
    Run OPTIMIZE TABLE to defragment tables for better performance
    MySQL started within last 24 hours - recommendations may be inaccurate
    Increase table_cache gradually to avoid file descriptor limits
    Your applications are not closing MySQL connections properly
Variables to adjust:
    query_cache_limit (> 1M, or use smaller result sets)
    table_cache (> 64)
```

Here we can see specific recommendations.

## Auditing and Improving Performance

### Checking Global Status

```sql
mysql> SHOW global STATUS;
| Binlog_cache_disk_use             | 0         |
| Binlog_cache_use                  | 16        |
```

- Binlog_cache_disk_use should remain at 0 for best performance (no direct disk writing, going through the cache first)
- Using Binlog_cache_use for binary logs is good when this number is not 0.

```sql
| Created_tmp_disk_tables           | 8         |
| Created_tmp_files                 | 5         |
| Created_tmp_tables                | 54        |
```

- Created_tmp_disk_tables: should be 0 as often as possible (except when having blob or text fields, not much can be done)
- Created_tmp_files: when Created_tmp_disk_tables is not sufficient and more files need to be created on disk
- Created_tmp_tables: number of times temporary tables have been created

Here, I'd like to improve this a bit, although it's not really necessary given the small number of disk tables created. Let's check the current maximum allowed size:

```sql
mysql> SHOW global VARIABLES LIKE 'tmp_table_size';
+----------------+----------+
| Variable_name  | VALUE    |
+----------------+----------+
| tmp_table_size | 16777216 |
+----------------+----------+
1 ROW IN SET (0.00 sec)
```

The maximum size of a temporary table in memory is 16MB. To increase it to 32MB on-the-fly:

```sql
mysql> SET global tmp_table_size=32*1024*1024;
```

And don't forget to make it permanent in the configuration file:

```bash
...
tmp_table_size = 32M
...
```

Next, to make it effective, we need to increase the max_heap_table_size. First, let's check the current size:

```sql
mysql> SHOW global VARIABLES LIKE 'max_heap_table_size';
+---------------------+----------+
| Variable_name       | VALUE    |
+---------------------+----------+
| max_heap_table_size | 16777216 |
+---------------------+----------+
1 ROW IN SET (0.00 sec)
```

Increasing max_heap_table_size to at least the size of tmp_table_size (allocating more memory for memory tables). To increase it on-the-fly:

```sql
mysql> SET global max_heap_table_size=32*1024*1024;
```

To make it permanent, add to the configuration:

```bash
...
max_heap_table_size = 32M
...
```

### Enabling Slow Queries

To monitor slow queries, we need to enable slow_query and set a maximum value before queries are logged (here 1 second). First, let's see if it's enabled:

```sql
mysql> SHOW global VARIABLES LIKE '%log%';
+---------------------------------+------------+
| Variable_name                   | VALUE      |
+---------------------------------+------------+
...
| log_slow_queries                | OFF        |
...
+---------------------------------+------------+
26 ROWS IN SET (0.00 sec)
```

We can see it's disabled. Starting from MySQL 5.1, we can enable it dynamically:

```sql
mysql> SET global slow_query_log=1;
Query OK, 0 ROWS affected (0.02 sec)

mysql> SET global long_query_time=1;
Query OK, 0 ROWS affected (0.00 sec)
```

And to make it persistent:

```bash
...
slow_query_log=1
long_query_time=1
...
```

### General Status

```bash
> mysqladmin -uroot -p status
Enter password:
Uptime: 57774  Threads: 3  Questions: 5029  Slow queries: 1411  Opens: 373  Flush tables: 1  Open tables: 63  Queries per second avg: 0.87
```

- Questions: number of queries
- Slow queries: slow queries
- Opens: file openings
- Queries per second avg: queries per second (average since server start)

### Cache Management

```sql
mysql> SHOW global STATUS LIKE 'Qc%';
+-------------------------+-----------+
| Variable_name           | VALUE     |
+-------------------------+-----------+
| Qcache_free_blocks      | 1         |
| Qcache_free_memory      | 134208784 |
| Qcache_hits             | 0         |
| Qcache_inserts          | 207       |
| Qcache_lowmem_prunes    | 0         |
| Qcache_not_cached       | 2825      |
| Qcache_queries_in_cache | 0         |
| Qcache_total_blocks     | 1         |
+-------------------------+-----------+
8 ROWS IN SET (0.00 sec)
```

- Qcache_free_memory: Free memory to add queries in the query cache (storing the query + result (130 MB here))
- Qcache_hits: number of cache hits
- Qcache_inserts: Qcache_hits should be higher than Qcache_inserts for optimal cache performance (here since there isn't much traffic, this doesn't increase)

For better performance, determine if the query cache should be kept. If:

- Qcache_hits < Qcache_inserts: disable the cache
- Qcache_not_cached < (Qcache_hits + Qcache_inserts): Try increasing cache size and query size limit. If Qcache_not_cached continues to rise, disable the cache
- Qcache_lowmem_prunes: if no more space in cache, old queries will be replaced by new ones. Increase cache if you have too many Qcache_lowmem_prunes

#### Disabling the Query Cache

To disable the Query cache:

```sql
 SET global query_cache_type = 0
```

And to make it permanent, add to the configuration file:

```bash
query_cache_type = 0
```

#### Changing the Query Cache Size

Let's get the current cache size for a query and its result:

```sql
mysql> SHOW global VARIABLES LIKE 'query%';
+------------------------------+----------+
| Variable_name                | VALUE    |
+------------------------------+----------+
| query_alloc_block_size       | 8192     |
| query_cache_limit            | 1048576  |
| query_cache_min_res_unit     | 4096     |
| query_cache_size             | 16777216 |
| query_cache_type             | ON       |
| query_cache_wlock_invalidate | OFF      |
| query_prealloc_size          | 8192     |
+------------------------------+----------+
7 ROWS IN SET (0.00 sec)
```

The maximum size for a query cache is 1 MB, we can increase it to 2MB:

```sql
 SET global query_cache_limit=2*1024*1024;
```

Also changing the total cache size. For example, increase it to 32 MB:

```sql
 SET global query_cache_size = 32*1024*1024;
```

And finally, let's add this to the configuration:

```ini
 query_cache_limit = 2M
 query_cache_size = 32M
```

### Auditing Query Performance

Let's take a query that's interesting:

To see execution time of a query:

```sql
mysql> SET profiling =1;
mysql> SELECT  * FROM ma_table;
+---+-------------------------------------------+
| a | b                                         |
+---+-------------------------------------------+
| 1 | *0F105C1BB64CDADBA3E0AB29141550D4EDBDADCD |
| 2 | *B552D2C3751A26A345932AA1196C0D04BC9DC909 |
| 3 | *6F23566A5409E9809512C079C8D4EC7EF82AB8A1 |
| 4 | *65109C8FC01571CB9897AD479FF605F73DCD4752 |
| 5 | *7B9EBEED26AA52ED10C0F549FA863F13C39E0209 |
+---+-------------------------------------------+
5 ROWS IN SET (0.00 sec)
```

When the query cache is used (very good):

```sql
mysql> SHOW profile;
+--------------------------------+----------+
| STATUS                         | Duration |
+--------------------------------+----------+
| starting                       | 0.000035 |
| checking query cache FOR query | 0.000008 |
| checking privileges ON cached  | 0.000006 |
| sending cached RESULT TO clien | 0.000028 |
| logging slow query             | 0.000003 |
| cleaning up                    | 0.000003 |
+--------------------------------+----------+
6 ROWS IN SET (0.01 sec)
```

When it's not used:

```sql
mysql> SHOW profile;
+--------------------+----------+
| STATUS             | Duration |
+--------------------+----------+
| starting           | 0.000084 |
| Opening TABLES     | 0.000015 |
| System LOCK        | 0.000006 |
| TABLE LOCK         | 0.000010 |
| init               | 0.000019 |
| optimizing         | 0.000006 |
| statistics         | 0.000015 |
| preparing          | 0.000010 |
| executing          | 0.000004 |
| Sending DATA       | 0.000067 |
| END                | 0.000005 |
| query END          | 0.000003 |
| freeing items      | 0.000026 |
| logging slow query | 0.000003 |
| cleaning up        | 0.000003 |
+--------------------+----------+
15 ROWS IN SET (0.00 sec)
```

Execution time difference:

```sql
mysql> SHOW profiles;
+----------+------------+---------------------------------------------------+
| Query_ID | Duration   | Query                                             |
+----------+------------+---------------------------------------------------+
(...)
|       27 | 0.00008200 | SELECT  * FROM query_cache                        |
|       28 | 0.00027475 | SELECT SQL_NO_CACHE * FROM query_cache            |
+----------+------------+---------------------------------------------------+
```

To see how many times queries had to wait for another to finish before executing:

```sql
mysql> SHOW global STATUS LIKE 'Table%';
+-----------------------+-------+
| Variable_name         | VALUE |
+-----------------------+-------+
| Table_locks_immediate | 4549  |
| Table_locks_waited    | 0     |
+-----------------------+-------+
2 ROWS IN SET (0.00 sec)
```

Here Table_locks_waited is 0, which is perfect. This number generally increases when you have high load and are not using the right engine (typically using MyISAM instead of InnoDB).

### Improving the Engine Based on Usage

Let's first see what we have:

```sql
mysql> SELECT TABLE_SCHEMA,ENGINE,SUM(TABLE_ROWS),ENGINE,SUM(DATA_LENGTH),SUM(INDEX_LENGTH) FROM INFORMATION_SCHEMA.TABLES GROUP BY ENGINE,TABLE_SCHEMA ORDER BY TABLE_SCHEMA;
+--------------------+--------+-----------------+--------+------------------+-------------------+
| TABLE_SCHEMA       | ENGINE | SUM(TABLE_ROWS) | ENGINE | SUM(DATA_LENGTH) | SUM(INDEX_LENGTH) |
+--------------------+--------+-----------------+--------+------------------+-------------------+
| information_schema | MEMORY |            NULL | MEMORY |                0 |                 0 |
| information_schema | MyISAM |            NULL | MyISAM |                0 |              4096 |
| jahia              | InnoDB |              10 | InnoDB |           311296 |             98304 |
| jahia              | MyISAM |            4268 | MyISAM |           661004 |            661504 |
| jira               | MyISAM |         1726671 | MyISAM |        167178929 |          73067520 |
| mysql              | MyISAM |            1774 | MyISAM |           465108 |             68608 |
| sugarcrm           | MyISAM |               4 | MyISAM |              719 |              7168 |
| sugarcrm           | InnoDB |          351405 | InnoDB |         89718784 |         131317760 |
+--------------------+--------+-----------------+--------+------------------+-------------------+
8 ROWS IN SET (2.60 sec)
```

Here we mainly see tables in MyISAM and InnoDB.

To check the size of the index cache for MyISAM:

```sql
mysql> SHOW global STATUS LIKE 'KEY%';
+------------------------+---------+
| Variable_name          | VALUE   |
+------------------------+---------+
| Key_blocks_not_flushed | 0       |
| Key_blocks_unused      | 14497   |
| Key_blocks_used        | 12268   |
| Key_read_requests      | 2413795 |
| Key_reads              | 67366   |
| Key_write_requests     | 106173  |
| Key_writes             | 104947  |
+------------------------+---------+
7 ROWS IN SET (0.00 sec)
```

There is still space (in KB) for Key_blocks_unused corresponding to unused space.

To check the size of the index cache for InnoDB:

```sql
mysql> SHOW engine innodb STATUS\G;
...
----------------------
BUFFER POOL AND MEMORY
----------------------
Total memory allocated 18804332; IN additional pool allocated 1048576
Buffer pool SIZE   512
Free buffers       0
DATABASE pages     511
Modified db pages  33
Pending reads 0
Pending writes: LRU 0, FLUSH list 0, single page 0
Pages READ 41821, created 371, written 102582
0.18 reads/s, 0.00 creates/s, 0.00 writes/s
Buffer pool hit rate 1000 / 1000
...
```

Here the buffer is (512\*16) 8MB and there is no free space (Free buffers). So to address the problem, we need to increase the buffer pool. A reminder:

```sql
mysql> SELECT TABLE_SCHEMA,ENGINE,SUM(TABLE_ROWS),ENGINE,SUM(DATA_LENGTH),SUM(INDEX_LENGTH) FROM INFORMATION_SCHEMA.TABLES GROUP BY ENGINE,TABLE_SCHEMA ORDER BY TABLE_SCHEMA;
+--------------------+--------+-----------------+--------+------------------+-------------------+
| TABLE_SCHEMA       | ENGINE | SUM(TABLE_ROWS) | ENGINE | SUM(DATA_LENGTH) | SUM(INDEX_LENGTH) |
+--------------------+--------+-----------------+--------+------------------+-------------------+
...
| sugarcrm           | InnoDB |          351405 | InnoDB |         89718784 |         131317760 |
...
```

- SUM(DATA_LENGTH): Data size without indexes
- SUM(INDEX_LENGTH): Size of indexes without data

That's: (89718784 + 131317760) =~ 100MB
So we'll increase the buffer pool to 128MB (more than necessary):

```bash
 ...
 innodb_buffer_pool_size = 128M
 ..
```

Let's also change the default parameters with best practices. Note that changing the log size will require some modifications before restarting. Modify your MySQL configuration like this:

```ini
...
default-storage-engine = innodb
innodb_buffer_pool_size = 128M
innodb_data_file_path=ibdata1:10M:autoextend
innodb_additional_mem_pool_size = 20M
innodb_file_per_table
# Warning : changing this needs stop of mysql, removal (backup of ib_log* files), and mysql startup
innodb_log_file_size = 256M
innodb_log_buffer_size = 8M
innodb_flush_log_at_trx_commit = 1
...
```

Then follow these steps:

- Backup the 'ib_log\*' files
- Stop the MySQL server
- Delete the 'ib_logs\*' files
- Restart the server

Note: on Debian, you might hit the timeout of the service, but it will still start (a small pgrep -f mysql will show where it is, as well as the logs).

### Checking Open Files

It can be important to see which tables cannot be cached (and therefore files on disk):

```sql
mysql> SHOW global STATUS LIKE 'Open%';
+---------------+-------+
| Variable_name | VALUE |
+---------------+-------+
| Open_files    | 130   |
| Open_streams  | 0     |
| Open_tables   | 64    |
| Opened_tables | 14581 |
+---------------+-------+
4 ROWS IN SET (0.00 sec)
```

Here I have a few open tables and a bunch of tables that have been opened.

To address this issue, let's change the size of the table cache which is currently at 64MB:

```sql
mysql> SHOW global VARIABLES LIKE 'table%';
+-------------------------+--------+
| Variable_name           | VALUE  |
+-------------------------+--------+
| table_cache             | 64     |
| table_lock_wait_timeout | 50     |
| table_type              | MyISAM |
+-------------------------+--------+
3 ROWS IN SET (0.00 sec)
```

Let's increase it to 128 MB. However, we won't do it on-the-fly (side effects almost guaranteed) and it's not recommended to exceed 4GB (here we have plenty of room).

Let's modify the configuration:

```sql
 ...
 table_cache = 128M
 ...
```

### Resetting Statistics

To reset global status, you need to run this command:

```bash
flush status;
```

Now you can see it:

```bash
show global status;
```

If you restart your MySQL server, it does the same thing.

### Increasing Performance for Temporary File Access

The idea is to point temporary files to a RAM filesystem to store temporary files. This will significantly increase performance. First, build your temporary filesystem:

- [Documentation for Linux]({{< ref "docs/Linux/FilesystemsAndStorage/tmpfs_ram_filesystem_or_how_to_write_to_ram.md" >}})
- [Documentation for Solaris]({{< ref "docs/Solaris/Filesystems/tmpfs_mounting_a_ram_filesystem_on_solaris.md" >}})

Next, modify your MySQL configuration file to change the tmpdir variable and point it to your new temporary space:

```bash {linenos=table,hl_lines=[3]}
...
# The MySQL server
[mysqld]
tmpdir = /mnt/mysql_tmpfs
...
```

Restart MySQL and check that the new parameter is correctly applied:

```bash {linenos=table,hl_lines=[5]}
mysql> show variables like 'tmpdir';
+---------------+------------------+
| Variable_name | Value            |
+---------------+------------------+
| tmpdir        | /mnt/mysql_tmpfs |
+---------------+------------------+
1 row in set (0.00 sec)
```

And there you go, now it's fast :-)

### Analyzing Columns

It's important to know what can be optimized in your table structure for all tables in the database. There is a command to analyze a database and provide optimization hints per column:

```sql
SELECT * FROM my_table PROCEDURE ANALYSE();
```

You can then compare with the current table structure:

```sql
DESCRIBE my_table;
```

## Table Optimization

I created a small Perl script to automate database optimization. To use it, simply use a user who has select and insert rights on all databases. You can also exclude certain databases. To ensure you have the latest version of this script, go to [my git](https://www.deimos.fr/gitweb/?p=git_deimosfr.git;a=tree). Here's the beginning of the code:

```perl
#!/usr/bin/perl
#===============================================================================
#
#         FILE:  optimizer.pl
#
#        USAGE:  ./optimizer.pl
#
#  DESCRIPTION:  MySQL Automatic Tables Optimizer
#
#      OPTIONS:  ---
# REQUIREMENTS:  ---
#         BUGS:  ---
#        NOTES:  ---
#       AUTHOR:  Pierre Mavro (), xxx@mycompany.com
#      COMPANY:
#      VERSION:  0.1
#      CREATED:  08/04/2010 17:50:23
#     REVISION:  ---
#===============================================================================

use strict;
use warnings;
use DBI;

# MySQL configuration
# This user should have insert and select rights
my $host = 'localhost';
my $database = 'mysql';
my $user = 'optimizer';
my $pw = 'optimizer';
my @exclude_databases = qw/information_schema database1 database2/;

########## DO NOT TOUCH NOW ##########

# Vars
my (@all_databases,@all_tables);
my ($aref, $cur_database, $cur_table);

# Connect to the BDD
my $dbdetails = "DBI:mysql:$database;host=$host";
my $dbh = DBI->connect($dbdetails, $user, $pw) or die "Could not connect to database: $DBI::errstr\n";

# Get all databases
my $sth=$dbh->prepare(q{SHOW DATABASES}) or die "Unable to prepare show databases: ". $dbh->errstr."\n";
$sth->execute or die "Unable to exec show databases: ". $dbh->errstr."\n";
while ($aref = $sth->fetchrow_arrayref)
{
    push(@all_databases,$aref->[0]);
}
$sth->finish;
# Disconnect the BDD
$dbh->disconnect();

# Optimize all tables of all databases
my $unwanted_found=0;
foreach $cur_database (@all_databases)
{
    # Exclude optimization on unwanted databases
    foreach (@exclude_databases)
    {
        if ($cur_database eq $_)
        {
            $unwanted_found=1;
            last;
        }
    }
    if ($unwanted_found == 1)
    {
        $unwanted_found=0;
        next;
    }
    # Connect on database
    my $dbdetails = "DBI:mysql:$cur_database;host=$host";
    my $dbh = DBI->connect($dbdetails, $user, $pw) or die "Could not connect to database: $DBI::errstr\n";
    # Get the List of tables
    @all_tables = $dbh->tables;
    # Optimize all tables
    foreach $cur_table (@all_tables)
    {
        $dbh->do("optimize table $cur_table");
    }
    $dbh->disconnect();
}
```

## Patch mysqltuner

Here's a small patch I wrote to fix a bug when trying to connect to a port other than the default:

```diff
*** mysqltuner.pl.old   mar. avr  10 09:26:01 2012
--- mysqltuner.pl       mar. avr  10 09:28:14 2012
***************
*** 274,280 ****
        }
        # Did we already get a username and password passed on the command line?
        if ($opt{user} ne 0 and $opt{pass} ne 0) {
!               $mysqllogin = "-u $opt{user} -p'$opt{pass}'".$remotestring;
                my $loginstatus = `mysqladmin ping $mysqllogin 2>&1`;
                if ($loginstatus =~ /mysqld is alive/) {
                        goodprint "Logged in using credentials passed on the command line\n";
--- 274,280 ----
        }
        # Did we already get a username and password passed on the command line?
        if ($opt{user} ne 0 and $opt{pass} ne 0) {
!               $mysqllogin = "-u $opt{user} -p'$opt{pass}' -P $opt{port}".$remotestring;
                my $loginstatus = `mysqladmin ping $mysqllogin 2>&1`;
                if ($loginstatus =~ /mysqld is alive/) {
                        goodprint "Logged in using credentials passed on the command line\n";
```

I submitted it, but since there's no maintainer anymore...

## Resources

- http://dev.mysql.com/doc/refman/5.0/en/temporary-files.html
