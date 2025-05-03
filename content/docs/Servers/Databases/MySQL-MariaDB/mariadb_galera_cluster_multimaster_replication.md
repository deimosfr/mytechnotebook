---
weight: 999
url: "/MariaDB_Galera_Cluster_\\:_la_réplication_multi_maitres/"
title: "MariaDB Galera Cluster: Multi-Master Replication"
description: "Learn how to set up and manage MariaDB Galera Cluster for multi-master replication in a database environment with synchronous replication across multiple nodes."
categories: ["Debian", "Storage", "Networking"]
date: "2014-04-19T05:51:00+02:00"
lastmod: "2014-04-19T05:51:00+02:00"
tags:
  [
    "MariaDB",
    "Galera",
    "Database",
    "Clustering",
    "Replication",
    "High Availability",
    "MySQL",
  ]
toc: true
---

![MariaDB](/images/galeracluster_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 5.5 |
| **Operating System** | Debian 7 |
| **Website** | [MariaDB Website](https://mariadb.org/) |
| **Last Update** | 19/04/2014 |
| **Others** | Galera 23.2.4(r147) |
{{< /table >}}

## Introduction

If you have already set up MySQL replication, you know that you're limited to 2 master nodes maximum. Many actions must be done manually, making it difficult to have something that's both scalable and accessible for simultaneous writes, since by default writes are synchronous.

Note: If you need advice or support for MySQL/MariaDB/Galera, I recommend [OceanDBA](https://www.oceandba.fr/).

There's a tool called [Galera](https://www.codership.com) that integrates with MariaDB or MySQL (through recompilation in both cases) and allows multi-master replication (3 nodes minimum). There are several products that use Galera:

- MariaDB Galera Cluster
- Percona Cluster
- MySQL Galera Cluster

For those wondering, this is really different from MySQL Cluster which knows how to scale writes. Here it's only multi-threaded. And unlike MHA which is an asynchronous solution, Galera is synchronous.

Galera **only works with the InnoDB engine** and allows:

- Synchronous replication
- Active multi-master replication
- Simultaneous read/write on multiple nodes
- Automatic detection when a node fails
- Automatic node reintegration
- No lag on slaves
- No lost transactions
- Lower client latency

Although this seems perfect on paper, there are [some limitations](https://www.codership.com/wiki/doku.php?id=limitations):

- Only supports InnoDB
- All tables must have primary keys
- DELETE only works on tables with primary keys
- LOCK/UNLOCK/GET_LOCK/RELEASE_LOCK doesn't work in multi-master
- Query logs can only be sent to files, not tables
- XA transactions are not supported

To help you understand a typical architecture:

![Galera](/images/galera.avif)

There's also an online tool to help build this kind of infrastructure: [Galera Configurator](https://www.severalnines.com/galera-configurator/)

## Installation

To install MariaDB, it's unfortunately not embedded in Debian, so we'll add a repository. First of all, install a python tool to get aptkey:

```bash
aptitude install python-software-properties
```

Then let's add this repository (https://downloads.mariadb.org/mariadb/repositories/):

```bash
apt-key adv --recv-keys --keyserver keyserver.ubuntu.com 0xcbcb082a1bb943db
add-apt-repository 'deb http://mirrors.linsrv.net/mariadb/repo/10.0/debian wheezy main'
```

We're now going to change apt pinning to prioritize MariaDB's repository:

```bash
Package: *
Pin: release o=MariaDB
Pin-Priority: 1000
```

{{< alert context="info" text="At the time of writing, there is a small dependency issue with some packages. You'll need to download them in advance and install them. You don't need to follow the rest of this installation if you haven't encountered errors with the previous steps." />}}

```bash
wget http://ftp.fr.debian.org/debian/pool/main/o/openssl/libssl0.9.8_0.9.8o-4squeeze14_amd64.deb
dpkg -i libssl0.9.8_0.9.8o-4squeeze14_amd64.deb
```

Now, we'll install MariaDB and Galera:

```bash
aptitude update
aptitude install mariadb-galera-server galera rsync openntpd
```

{{< alert context="info" text="The rsync package is not mandatory, but necessary if you're going to use this transfer method later." />}}

## Configuration

### MariaDB

Before we start changing the configuration, we need to delete some log files or MariaDB won't start. We'll need to stop the MariaDB service:

```bash
service mysql stop
```

Then we'll apply this MariaDB configuration:

```ini
# MariaDB database server configuration file.
# Pierre Mavro / Deimosfr
#
# You can copy this file to one of:
# - "/etc/mysql/my.cnf" to set global options,
# - "~/.my.cnf" to set user-specific options.
#
# One can use all long options that the program supports.
# Run program with --help to get a list of available options and with
# --print-defaults to see which it would actually understand and use.
#
# For explanations see
# http://dev.mysql.com/doc/mysql/en/server-system-variables.html

# This will be passed to all mysql clients
# It has been reported that passwords should be enclosed with ticks/quotes
# escpecially if they contain "#" chars...
# Remember to edit /etc/mysql/debian.cnf when changing the socket location.
[client]
port		= 3306
socket		= /var/run/mysqld/mysqld.sock

# Here is entries for some specific programs
# The following values assume you have at least 32M ram

# This was formally known as [safe_mysqld]. Both versions are currently parsed.
[mysqld_safe]
socket		= /var/run/mysqld/mysqld.sock
nice		= 0

[mysqld]
#
# * Basic Settings
#
user		= mysql
pid-file	= /var/run/mysqld/mysqld.pid
socket		= /var/run/mysqld/mysqld.sock
port		= 3306
basedir		= /usr
datadir		= /var/lib/mysql
innodb_log_group_home_dir = /var/lib/mysql
tmpdir		= /tmp
lc_messages_dir	= /usr/share/mysql
lc_messages	= en_US
skip-external-locking

character_set_server = utf8
collation_server = utf8_general_ci

#
# Instead of skip-networking the default is now to listen only on
# localhost which is more compatible and is not less secure.
bind-address		= 0.0.0.0
#
# * Fine Tuning
#
max_connections		= 500
connect_timeout		= 5
wait_timeout		= 600
max_allowed_packet	= 16M
thread_cache_size       = 128
sort_buffer_size	= 16M
bulk_insert_buffer_size	= 16M
tmp_table_size		= 32M
max_heap_table_size	= 64M
net_buffer_length	= 4k

#
# * MyISAM
#
# This replaces the startup script and checks MyISAM tables if needed
# the first time they are touched. On error, make copy and try a repair.
myisam_recover          = BACKUP
key_buffer_size		= 128M
#open-files-limit	= 2000
table_open_cache	= 400
myisam_sort_buffer_size	= 512M
concurrent_insert	= 2
read_buffer_size	= 2M
read_rnd_buffer_size	= 1M
#
# * Query Cache Configuration
#
# Cache only tiny result sets, so we can fit more in the query cache.
query_cache_limit		= 128K
query_cache_size		= 64M
# for more write intensive setups, set to DEMAND or OFF
#query_cache_type		= DEMAND
#
# * Logging and Replication
#
# Both location gets rotated by the cronjob.
# Be aware that this log type is a performance killer.
# As of 5.1 you can enable the log at runtime!
#general_log_file        = /var/log/mysql/mysql.log
#general_log             = 1
#
# Error logging goes to syslog due to /etc/mysql/conf.d/mysqld_safe_syslog.cnf.
#
# we do want to know about network errors and such
log_warnings		= 2
#
# Enable the slow query log to see queries with especially long duration
#slow_query_log[={0|1}]
slow_query_log_file	= /var/log/mysql/mariadb-slow.log
long_query_time = 10
#log_slow_rate_limit	= 1000
log_slow_verbosity	= query_plan

#log-queries-not-using-indexes
#log_slow_admin_statements
#
# The following can be used as easy to replay backup logs or for replication.
# note: if you are setting up a replication slave, see README.Debian about
#       other settings you may need to change.
#server-id		= 1
#report_host		= master1
#auto_increment_increment = 2
#auto_increment_offset	= 1
log_bin			= /var/log/mysql/mariadb-bin
log_bin_index		= /var/log/mysql/mariadb-bin.index
# not fab for performance, but safer
#sync_binlog		= 1
expire_logs_days	= 10
max_binlog_size         = 100M
# slaves
#relay_log		= /var/log/mysql/relay-bin
#relay_log_index	= /var/log/mysql/relay-bin.index
#relay_log_info_file	= /var/log/mysql/relay-bin.info
#log_slave_updates
#read_only
#
# If applications support it, this stricter sql_mode prevents some
# mistakes like inserting invalid dates etc.
#sql_mode		= NO_ENGINE_SUBSTITUTION,TRADITIONAL
#
# * InnoDB
#
# InnoDB is enabled by default with a 10MB datafile in /var/lib/mysql/.
# Read the manual for more InnoDB related options. There are many!
default_storage_engine	= InnoDB
# you can't just change log file size, requires special procedure
#innodb_log_file_size			= 50M
innodb_buffer_pool_size			= 256M
innodb_log_buffer_size			= 8M
innodb_log_file_size			= 256M
thread_concurrency			= 64
innodb_thread_concurrency		= 64
innodb_read_io_threads			= 16
innodb_write_io_threads			= 16
innodb_flush_log_at_trx_commit 		= 2
innodb_file_per_table			= 1
innodb_open_files			= 400
innodb_io_capacity			= 600
innodb_lock_wait_timeout 		= 60
innodb_flush_method			= O_DIRECT
innodb_doublewrite 			= 0
innodb_additional_mem_pool_size		= 20M
innodb_buffer_pool_restore_at_startup	= 500
innodb_file_per_table
#
# * Security Features
#
# Read the manual, too, if you want chroot!
# chroot = /var/lib/mysql/
#
# For generating SSL certificates I recommend the OpenSSL GUI "tinyca".
#
# ssl-ca=/etc/mysql/cacert.pem
# ssl-cert=/etc/mysql/server-cert.pem
# ssl-key=/etc/mysql/server-key.pem

[mysqldump]
quick
quote-names
max_allowed_packet	= 16M

[mysql]
#no-auto-rehash	# faster start of mysql but no tab completition

[isamchk]
key_buffer		= 16M

[mysqlhotcopy]
interactive-timeout

#
# * IMPORTANT: Additional settings that can override those from this file!
#   The files must end with '.cnf', otherwise they'll be ignored.
#
!includedir /etc/mysql/conf.d/
```

Now we'll delete the log files because we just changed their configuration (or we won't be able to start the instance) and be able to start our MariaDB service:

```bash
rm /var/lib/mysql/ib_logfile*
service mysql start
```

### Galera

Here's the configuration to apply for a cluster on the same site:

```ini {linenos=table,hl_lines=["19-27"]}
# MariaDB-specific config file.
# Read by /etc/mysql/my.cnf

[client]
# Default is Latin1, if you need UTF-8 set this (also in server section)
#default-character-set = utf8

[mysqld]
#
# * Character sets
#
# Default is Latin1, if you need UTF-8 set all this (also in client section)
#
#character-set-server  = utf8
#collation-server      = utf8_general_ci
#character_set_server   = utf8
#collation_server       = utf8_general_ci

# Load Galera Cluster
wsrep_provider = /usr/lib/galera/libgalera_smm.so
wsrep_cluster_name='mariadb_cluster'
wsrep_node_name=node2
wsrep_node_address="10.0.0.2"
wsrep_cluster_address = 'gcomm://10.0.0.1,10.0.0.2,10.0.0.3,10.0.0.4'
wsrep_retry_autocommit = 0
wsrep_sst_method = rsync
wsrep_provider_options="gcache.size = 1G; gcache.name = /tmp/galera.cache"
#wsrep_replication_myisam = 1
#wsrep_sst_receive_address = <x.x.x.x>
#wsrep_notify_cmd="script.sh"

# Other mysqld options
binlog_format = ROW
innodb_autoinc_lock_mode = 2
innodb_flush_log_at_trx_commit = 2
innodb_locks_unsafe_for_binlog = 1
```

- wsrep_cluster_name: the name of the Galera cluster. Use this especially if you have multiple Galera clusters in the same subnet to prevent nodes from entering the wrong cluster.
- wsrep_node_name: the name of the machine where this configuration file is located. You'll understand that **you absolutely must avoid duplicates** (especially for debugging ;-))
- wsrep_node_address: IP address of the current node (same warning as the previous line)
- wsrep_cluster_address: list of cluster members that can be master (separated by commas).
- wsrep_provider_options: enables [additional options](https://www.codership.com/wiki/doku.php?id=galera_parameters_0.8).
  - gcache allows storing data to transfer to other nodes. Default is 128M, it's advised to increase this value.
- wsrep_retry_autocommit: defines the number of times a query should be retried in case of conflict.
- wsrep_sst_method: the data exchange method. Rsync is currently the fastest.
- wsrep_replication_myisam: enables replication of MyISAM data (**no transaction management...so avoid it!**)
- wsrep_sst_receive_address: forces the use of a certain address for remote hosts to connect (resolves VIP problems)
- wsrep_notify_cmd: allows executing a script for each Galera event (node state change)
- binlog_format: defines the log format in ROW mode
- innodb_autoinc_lock_mode: changes the lock behavior
- innodb_flush_log_at_trx_commit: performance optimization

{{< alert context="info" text="Adapt the <i>wsrep_node_name</i> and <i>wsrep_cluster_address</i> lines to your respective machines." />}}

The above configuration is applicable to all nodes except the master (node 1). The master should have the identical configuration with the only difference being this line:

```bash
wsrep_cluster_address = 'gcomm://'
```

{{< alert context="warning" text="It's important that <b>only one machine</b> has the 'gcomm://' configuration, as this initializes the cluster." />}}

### Geo cluster

It's possible to set up a geo cluster, but the disadvantage is that when there's a network outage exceeding the specified timeouts, one of the clusters will have to completely resynchronize. Here's the line to configure to modify these timeouts:

```bash
[...]
wsrep_provider_options = "evs.keepalive_period = PT3S; evs.inactive_check_period = PT10S; evs.suspect_timeout = PT30S; evs.inactive_timeout = PT1M; evs.install_timeout = PT1M"
[...]
```

If you use the 'wsrep_sst_receive_address' option, you'll need to add a parameter to this line (ist.recv_addr) with the same IP as the 'wsrep_sst_receive_address' option:

```bash
[...]
wsrep_provider_options = "evs.keepalive_period = PT3S; evs.inactive_check_period = PT10S; evs.suspect_timeout = PT30S; evs.inactive_timeout = PT1M; evs.install_timeout = PT1M; ist.recv_addr = <x.x.x.x>"
[...]
```

### Replication methods

There are several solutions for data transfers between nodes. In the example above, we used rsync. When a machine requests a donor (another machine) to receive data, **transactions are blocked for the donor during the data exchange period with the receiver!!!**

{{< alert context="warning" text="If you use a load balancer, you'll need to remove the donor node during this period." />}}

This block can be limited by using [the xtrabackup method](#xtrabackup).

#### mysqldump

SST (State Snapshot Transfer) will allow complete exchanges (only complete, not incremental). So you need to create a user on all machines:

```sql
grant all on *.* to 'sst_user'@'%' identified by 'sst_password';
```

Then modify the configuration so that the user and password match:

```bash
[...]
wsrep_sst_auth = 'sst_user:sst_password'
[...]
```

So when we add a new node, we can see that a node has become a donor:

```sql {linenos=table,hl_lines=[7]}
MariaDB [(NONE)]> SHOW global STATUS LIKE 'wsrep%stat%';
+---------------------------+--------------------------------------+
| Variable_name             | VALUE                                |
+---------------------------+--------------------------------------+
| wsrep_local_state_uuid    | 3e10ea72-f2b9-11e2-0800-4b821a7d26d5 |
| wsrep_local_state         | 2                                    |
| wsrep_local_state_comment | Donor/Desynced                       |
| wsrep_cluster_state_uuid  | 3e10ea72-f2b9-11e2-0800-4b821a7d26d5 |
| wsrep_cluster_status      | PRIMARY                              |
+---------------------------+--------------------------------------+
5 ROWS IN SET (0.00 sec)
```

We also see that the wsrep_local_state value changes to 4 when the process is complete.

{{< alert context="info" text="The problem with this method is that it doesn't handle IST (Incremental State Transfer). Everything must be transferred in case of problems, not just the incremental data. You need to look at rsync or xtrabackup methods to be able to do IST." />}}

#### Rsync

This is a very efficient method! The major disadvantage is that it doesn't allow hot transfers. Here's how to configure SST:

```bash
[...]
wsrep_sst_method = rsync
[...]
```

The data must no longer be accessed by MariaDB for this to work properly. Another prerequisite is the presence of rsync on the servers.

#### XtraBackup

**XtraBackup is the best method today**. It allows transfers while minimizing lock time to just a few seconds. To use this method, we need to install XtraBackup.

To install it on Debian, it's very simple, we'll add its repository:

```bash
apt-key adv --keyserver keys.gnupg.net --recv-keys 1C4CBDCDCD2EFD2A
```

Create this file to add the repository:

```
deb http://repo.percona.com/apt VERSION main
deb-src http://repo.percona.com/apt VERSION main
```

Update and install:

```bash
aptitude update
aptitude install xtrabackup
```

Then we can configure our nodes to use this method:

```bash
[...]
wsrep_sst_method = xtrabackup
[...]
```

#### Specifying a donor

For replication, if you want to dedicate a machine as a donor (and possibly for backups) from a server:

```sql
SET global wsrep_sst_donor=<node3>
```

This solves the load balancer issue mentioned above and avoids blocks during a backup.

Otherwise, you can specify it directly for node integration:

```bash
mysqld --wsrep_cluster_address='gcomm://<node1>' --wsrep_sst_donor='<node1>'
```

{{< alert context="info" text="If you start MariaDB twice in a row without running a Galera sync at least once, the second time it will do a complete sync." />}}

## Usage

Before using, you need to understand the principle. When we start our MariaDB instances, we'll end up in this configuration (I deliberately didn't draw all the communication arrows to avoid clutter, but all nodes talk to each other):

![Galera1](/images/galera1.avif)

- Node 1 initializes the cluster with the empty gcomm value.
- The other nodes connect to node 1 and exchange their data to have the same data level everywhere

### Creating the cluster

Go to **node 1** to create the cluster. We'll start it with an empty cluster address, which will indicate its creation:

```bash
service mysql start --wsrep_cluster_address='gcomm://'
```

or

```bash
mysqld --wsrep_cluster_address='gcomm://'
```

### Adding nodes to the cluster

To join nodes to the newly created cluster, it's simple:

```bash
service mysql start --wsrep_cluster_address='gcomm://<ip_of_node_1>'
```

or

```bash
mysqld --wsrep_cluster_address='gcomm://<ip_of_node_1>'
```

Put the IP of node 1 to connect to it. You can also simply start the MariaDB service, as we have a working configuration:

```bash
service mysql start
```

If you can't reach the master, run this command in mysql on the master to make sure it's properly started:

```sql
SET global wsrep_cluster_address='gcomm://';
```

{{< alert context="info" text="You can use IPs or DNS names." />}}

### Checking the cluster status

To check the cluster status, here's the command to run in MariaDB:

```bash {linenos=table,hl_lines=[6,14]}
MariaDB [(none)]> SHOW STATUS LIKE 'wsrep_%';
+--------------------------+----------------------+
| Variable_name            | Value                |
+--------------------------+----------------------+
| wsrep_cluster_conf_id    | 18446744073709551615 |
| wsrep_cluster_size       | 0                    |
| wsrep_cluster_state_uuid |                      |
| wsrep_cluster_status     | Disconnected         |
| wsrep_connected          | OFF                  |
| wsrep_local_index        | 18446744073709551615 |
| wsrep_provider_name      |                      |
| wsrep_provider_vendor    |                      |
| wsrep_provider_version   |                      |
| wsrep_ready              | ON                   |
+--------------------------+----------------------+
10 rows in set (0.01 sec)
```

Here I only have my main server running. No other nodes have joined the cluster yet. But when I add some:

```bash {linenos=table,hl_lines=[36,44]}
> mysql -uroot -p -e "SHOW STATUS LIKE 'wsrep_%';"
+----------------------------+----------------------------------------------------------------------------+
| Variable_name              | Value                                                                      |
+----------------------------+----------------------------------------------------------------------------+
| wsrep_local_state_uuid     | 9e9f8568-a025-11e2-0800-be0dc874ac98                                       |
| wsrep_protocol_version     | 4                                                                          |
| wsrep_last_committed       | 0                                                                          |
| wsrep_replicated           | 0                                                                          |
| wsrep_replicated_bytes     | 0                                                                          |
| wsrep_received             | 12                                                                         |
| wsrep_received_bytes       | 639                                                                        |
| wsrep_local_commits        | 0                                                                          |
| wsrep_local_cert_failures  | 0                                                                          |
| wsrep_local_bf_aborts      | 0                                                                          |
| wsrep_local_replays        | 0                                                                          |
| wsrep_local_send_queue     | 0                                                                          |
| wsrep_local_send_queue_avg | 0.000000                                                                   |
| wsrep_local_recv_queue     | 0                                                                          |
| wsrep_local_recv_queue_avg | 0.000000                                                                   |
| wsrep_flow_control_paused  | 0.000000                                                                   |
| wsrep_flow_control_sent    | 0                                                                          |
| wsrep_flow_control_recv    | 0                                                                          |
| wsrep_cert_deps_distance   | 0.000000                                                                   |
| wsrep_apply_oooe           | 0.000000                                                                   |
| wsrep_apply_oool           | 0.000000                                                                   |
| wsrep_apply_window         | 0.000000                                                                   |
| wsrep_commit_oooe          | 0.000000                                                                   |
| wsrep_commit_oool          | 0.000000                                                                   |
| wsrep_commit_window        | 0.000000                                                                   |
| wsrep_local_state          | 4                                                                          |
| wsrep_local_state_comment  | Synced                                                                     |
| wsrep_cert_index_size      | 0                                                                          |
| wsrep_causal_reads         | 0                                                                          |
| wsrep_incoming_addresses   | 10.0.0.1:3306,10.0.0.2:3306,10.0.0.3:3306,10.0.0.4:3306                    |
| wsrep_cluster_conf_id      | 2                                                                          |
| wsrep_cluster_size         | 4                                                                          |
| wsrep_cluster_state_uuid   | 9e9f8568-a025-11e2-0800-be0dc874ac98                                       |
| wsrep_cluster_status       | Primary                                                                    |
| wsrep_connected            | ON                                                                         |
| wsrep_local_index          | 0                                                                          |
| wsrep_provider_name        | Galera                                                                     |
| wsrep_provider_vendor      | Codership Oy <info@codership.com>                                          |
| wsrep_provider_version     | 23.2.4(r147)                                                               |
| wsrep_ready                | ON                                                                         |
+----------------------------+----------------------------------------------------------------------------+
```

Here I can clearly see my 4 master nodes :-)

## Garbd (quorum)

To avoid split brains (cluster inconsistencies), it's advisable to use a tool provided with the Galera cluster that acts as a cluster quorum, **especially if you're in 2-node mode** (and that's generally the main advantage). Here's a use case:

```
                    ,---------.
                    |  garbd  |
                    `---------'
         ,---------.     |     ,---------.
         | clients |     |     | clients |
         `---------'     |     `---------'
                    \    |    /
                     \ ,---. /
                     ('     `)
                    (   WAN   )
                     (.     ,)
                     / `---' \
                    /         \
         ,---------.           ,---------.
         |  node1  |           |  node2  |
         |  node3  |           |  node4  |
         `---------'           `---------'
        Data Center 1         Data Center 2
```

You'll need to use the garbd service. It's installed by default but simply not activated. To configure it, we'll edit its configuration:

```bash {linenos=table,hl_lines=[5,8]}
# Copyright (C) 2012 Coedership Oy
# This config file is to be sourced by garb service script.

# A space-separated list of node addresses (address[:port]) in the cluster
GALERA_NODES="10.0.0.1:4567 10.0.0.2:4567 10.0.0.3:4567 10.0.0.4:4567"

# Galera cluster name, should be the same as on the rest of the nodes.
GALERA_GROUP="mariadb_cluster"

# Optional Galera internal options string (e.g. SSL settings)
# see http://www.codership.com/wiki/doku.php?id=galera_parameters
# GALERA_OPTIONS=""

# Log file for garbd. Optional, by default logs to syslog
# LOG_FILE=""
```

Adapt GALERA_NODES with the list of all your nodes and GALERA_GROUP with the name of the Galera cluster. Now we just need to activate it at machine startup and start the service:

```bash
update-rc.d -f garb defaults
service garb start
```

Now, on one of your nodes, you'll see there's a new node, which is actually just the quorum:

```mysql
MariaDB [(none)]> SHOW STATUS LIKE 'wsrep_cluster_size';
+--------------------+-------+
| Variable_name      | Value |
+--------------------+-------+
| wsrep_cluster_size | 5     |
+--------------------+-------+
```

## Backups and restorations

For backups, there are several methods and [Xtrabackup](../xtrabackup_:_optimiser_ses_backups_mysql/) is again one of the favorites.

{{< alert context="info" text="Just like when a new node joins the cluster and blocks the donor's transactions, it's the same for backups, but only for MyISAM!" />}}

If you only use InnoDB and use Xtrabackup, there will be no transaction locks and therefore no special node needed for backups!

### Installation

To install on Debian, it's very simple, we'll add its repository:

```bash
apt-key adv --keyserver keys.gnupg.net --recv-keys 1C4CBDCDCD2EFD2A
```

Create this file to add the repository:

```
deb http://repo.percona.com/apt VERSION main
deb-src http://repo.percona.com/apt VERSION main
```

Update and install:

```bash
aptitude update
aptitude install xtrabackup
```

### Usage

#### Backup

To back up a Galera cluster and allow incremental backup while blocking the node that will perform the backups for only a few seconds:

```bash
innobackupex --galera-info --user=xxxxx --password=xxxx <directory_to_store_backup>
```

The 'galera-info' option prevents problems during the uuid request (which will always return 0). If this option is not specified, incremental restorations won't be possible.

{{< alert context="info" text="For databases containing only InnoDB tables, it's possible to have no blocking at all during the lock by adding the '--no-lock' option." />}}

#### Restoration

During restoration, one of the nodes **will have blocked transactions after the Xtrabackup restoration, while the differential is being applied. To restore, copy the backup files to the MariaDB directory:**

```bash
cp -Rf <backup_directory> /var/lib/mysql/
chown -Rf mysql. /var/lib/mysql/
```

[You can specify the donor](#specifying-a-donor) (optional), then check the backup position:

```bash
> cat <backup_directory>/xtrabackup_galera_info
cfa9b8f1-f37b-11e2-0800-b37f8ac5092c:1
```

And integrate the node into the cluster by specifying the position to avoid it getting a complete backup:

```bash
service mysql start --wsrep_cluster_address='gcomm://<ip_of_node_1>' --wsrep_start_position="cfa9b8f1-f37b-11e2-0800-b37f8ac5092c:1"
```

or

```bash
mysqld --wsrep_cluster_address='gcomm://<ip_of_node_1>' --wsrep_start_position="cfa9b8f1-f37b-11e2-0800-b37f8ac5092c:1"
```

## Recovery and maintenance

### Automatic method

You don't need to worry about replication if a node other than the master (node1) fails. Once repaired and powered on, it will automatically reconnect to node 1 and catch up. However, in case of a problem with node 1:

![Galera2](/images/galera2.avif)

The other nodes will continue to communicate with each other and wait for the master to return. Once the master is turned back on, you'll need to tell it another node to which it should connect to continue synchronization:

![Galera3](/images/galera3.avif)

Whether you want to force a reconnection or perform maintenance on the master node, it's advisable to redirect the other servers to another master to avoid outages:

```mysql
SET GLOBAL wsrep_cluster_address='gcomm://10.0.0.2';
```

You can then check the master node on your MariaDB instances:

```mysql
MariaDB [(none)]> SHOW VARIABLES LIKE 'wsrep_cluster_address';
+-----------------------+-----------------------+
| Variable_name         | Value                 |
+-----------------------+-----------------------+
| wsrep_cluster_address | gcomm://10.0.0.2      |
+-----------------------+-----------------------+
```

I've also tested violently shutting down any node and turning it back on. Once integrated into the cluster, it properly retrieves all the differential information. I had no corruption problems. The only issues I encountered were with MariaDB startup and lock problems as explained in the [FAQ](#FAQ).

### Manual method

It's possible to update a node from a delta between an up-to-date version and one that's behind. To do this, check the version in which one of the up-to-date nodes is:

```sql {linenos=table,hl_lines=[5,7]}
SHOW global STATUS LIKE 'wsrep%';
+----------------------------+---------------------------------------+
| Variable_name              | VALUE                                 |
+----------------------------+---------------------------------------+
| wsrep_local_state_uuid     | 3e10ea72-f2b9-11e2-0800-4b821a7d26d5  |
| wsrep_protocol_version     | 4                                     |
| wsrep_last_committed       | 6                                     |
```

Here we can see the uuid number and the position of the last commit (wsrep_last_committed).

{{< alert context="warning" text="The following command should only be used on a turned off server or you risk losing data on it!!!" />}}

With the server off, it's possible to retrieve the position of the last commit:

```bash {linenos=table,hl_lines=[12]}
> mysqld --wsrep_recover=1
130722 15:44:53 InnoDB: The InnoDB memory heap is disabled
130722 15:44:53 InnoDB: Mutexes and rw_locks use GCC atomic builtins
130722 15:44:53 InnoDB: Compressed tables use zlib 1.2.7
130722 15:44:53 InnoDB: Using Linux native AIO
130722 15:44:53 InnoDB: Initializing buffer pool, size = 256.0M
130722 15:44:53 InnoDB: Completed initialization of buffer pool
130722 15:44:53 InnoDB: highest supported file format is Barracuda.
130722 15:44:53  InnoDB: Waiting for the background threads to start
130722 15:44:54 Percona XtraDB (http://www.percona.com) 1.1.8-29.3 started; log sequence number 1603853
130722 15:44:54 [Note] Plugin 'FEEDBACK' is disabled.
130722 15:44:54 [Note] WSREP: Recovered position: 3e10ea72-f2b9-11e2-0800-4b821a7d26d5:4
130722 15:44:54  InnoDB: Starting shutdown...
130722 15:44:55  InnoDB: Shutdown completed; log sequence number 1603853
130722 15:44:55 [Note] mysqld: Shutdown complete
```

Then we can restart with the delta from this last commit:

```bash
mysqld --wsrep_start_position=<uuid>:<position>
```

or

```bash
mysqld --wsrep_start_position=3e10ea72-f2b9-11e2-0800-4b821a7d26d5:4
```

The delta is then performed and the server is now at position 6.

### Force a node to resynchronize

If you really don't know what a node has and you want to completely resynchronize it because you're unsure about its data, just delete the MariaDB content and restart it:

```bash
service mysql stop
rm -Rf /var/lib/mysql/*
service mysql start
```

All data will then resynchronize.

{{< alert context="warning" text="This can take some time if your database is large or if the bandwidth between nodes is low." />}}

### Split brain

When you have one or more nodes in a split brain state, it's possible to continue using the cluster and discard all changes from the other down nodes so they'll do a complete sync when they start back up. On a 'primary' server:

```sql
SET global wsrep_provider_options = 'pc.bootstrap=1';
SET global wsrep_provider_options = 'pc.ignore_quorum=0';
```

- pc.bootstrap: takes control of the other nodes in the cluster and indicates that it's the master
- pc.ignore_quorum: allows splitting nodes and having a split brain

Then restart your other servers so they synchronize.

{{< alert context="warning" text="Once the cluster is restored, you <b>MUST</b> set these variables back to False to avoid future split brains that you can't recover from!!!" />}}

## FAQ

### One of my MariaDB services refuses to start after shutting it down

It can happen that when you turn off a MariaDB service, the lock files aren't properly released and the rsync service is still running. To clean everything up (without completely restarting the machine), follow these steps:

1. Make sure MariaDB is no longer running:

```bash
service mysql stop
ps aux | grep mysql
```

2. If it's still running, kill the process!
3. Check that the rsync process is no longer running and kill it if it is
4. Delete the lock files:

```bash
rm -f /var/run/mysqld/mysqld.sock
```

5. Check that the directory for storing the pid exists, otherwise create it:

```bash
if [ ! -d /var/run/mysqld ]  ; then mkdir /var/run/mysqld  ; chown mysql. /var/run/mysqld  ; fi
```

6. You can now start the service, it should work:

```bash
service mysql start
```

Otherwise check the logs (/var/log/syslog)

### /dev/stderr: Permission denied

If you have this problem, it's due to a Galera bug that incorrectly redirects its error output:

```
/usr//bin/wsrep_sst_common: line 94: /dev/stderr: Permission denied
```

To fix the problem, simply set the correct permissions on the error output:

```bash
chmod 777 /proc/self/fd/2
```

## References

1. http://www.codership.com/wiki/doku.php?id=galera_deployment
2. http://www.severalnines.com/galera-configurator/
3. http://www.severalnines.com/blog/understanding-gcache-galera
4. http://www.codership.com/wiki/doku.php?id=sst_mysql
5. http://www.codership.com/wiki/doku.php?id=galera_node_fsm
6. http://www.sebastien-han.fr/blog/2012/04/01/mysql-multi-master-replication-with-galera/

- [MariaDB: high performances](https://www.amazon.fr/MariaDB-High-Performance-Pierre-MAVRO/dp/1783981601)
- [MariaDB MySQL Advanced](/pdf/mariadb_mysql_avance.pdf)
- File:Galera1_src.vsdx - File:Galera2_src.vsdx - File:Galera3_src.vsdx
