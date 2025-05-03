---
weight: 999
url: "/Lancer_plusieurs_instances_de_MySQL_sur_Solaris/"
title: "Running Multiple MySQL Instances on Solaris"
description: "Guide on how to install and configure multiple MySQL instances on Solaris, including configuration files, initialization, startup scripts, and SMF file creation."
categories: ["Storage", "Database", "MySQL", "Solaris"]
date: "2010-01-29T06:12:00+02:00"
lastmod: "2010-01-29T06:12:00+02:00"
tags:
  ["MySQL", "Solaris", "Database", "Multiple Instances", "SMF", "Configuration"]
toc: true
---

## Installation

Download MySQL

```bash
wget http://dev.mysql.com/get/Downloads/MySQL-5.1/mysql-5.1.41-solaris10-x86_64.pkg.gz/from/http://mir2.ovh.net/ftp.mysql.com/
```

Create the MySQL user

```bash
groupadd mysql
useradd -g mysql -s /bin/false -d /var/empty mysql
```

Verify that the user is properly created

```bash
# finger mysql
Login name: mysql
Directory: /var/empty                   Shell: /bin/sh
Never logged in.
No unread mail
No Plan.
```

Decompress the MySQL archive and install it

```bash
gunzip mysql-5.1.41-solaris10-x86_64.pkg.gz
pkgadd -d mysql-5.1.41-solaris10-x86_64.pkg
```

## Configuration

### my.cnf

Create two distinct configuration files by changing:

- The file name
- The port
- The socket
- The datadir

Example: `/etc/my-prod.cnf`

```ini
# The MySQL server
[mysqld]
port            = 3307
socket          = /var/lib/mysql/mysql-prod.sock
datadir = /mnt/ulprod-ld_mysql/databases
log-bin = /mnt/ulprod-ld_mysql/logs/mysql-bin
pid-file = /mnt/ulprod-ld_mysql/mysql-prod.pid
skip-locking
key_buffer = 64M
max_allowed_packet = 1M
table_cache = 512
sort_buffer_size = 2M
read_buffer_size = 2M
read_rnd_buffer_size = 8M
myisam_sort_buffer_size = 128M
thread_cache_size = 8
query_cache_size = 128M
query_cache_limit = 2M
thread_concurrency = 8
skip-name-resolve
innodb_data_file_path=ibdata1:10M:autoextend
innodb_buffer_pool_size = 2048M
innodb_additional_mem_pool_size = 20M
innodb_data_home_dir = /mnt/ulprod-ld_mysql/
innodb_log_group_home_dir = /mnt/ulprod-ld_mysql/
innodb_file_per_table
innodb_log_file_size = 256M
innodb_log_buffer_size = 8M
innodb_flush_log_at_trx_commit = 1
default-storage_engine=InnoDB
thread_cache_size = 16
slow_query_log = 1
long_query_time = 1
skip-federated
server-id       = 1
#log-bin=mysql-bin
#sync_binlog = 1

[mysqldump]
quick
max_allowed_packet = 16M

[mysql]
no-auto-rehash

[isamchk]
key_buffer = 256M
sort_buffer_size = 256M
read_buffer = 2M
write_buffer = 2M

[myisamchk]
key_buffer = 256M
sort_buffer_size = 256M
read_buffer = 2M
write_buffer = 2M

[mysqlhotcopy]
interactive-timeout
```

Example: `/etc/my-dr.cnf`

```ini
# The MySQL server
[mysqld]
port            = 3306
socket          = /var/lib/mysql/mysql-dr.sock
datadir = /mnt/ulprod-pa_mysql/databases
log-bin = /mnt/ulprod-pa_mysql/logs/mysql-bin
pid-file = /mnt/ulprod-pa_mysql/mysql-dr.pid
skip-locking
key_buffer = 64M
max_allowed_packet = 1M
table_cache = 512
sort_buffer_size = 2M
read_buffer_size = 2M
read_rnd_buffer_size = 8M
myisam_sort_buffer_size = 128M
thread_cache_size = 8
query_cache_size = 128M
query_cache_limit = 2M
thread_concurrency = 8
skip-name-resolve
innodb_data_file_path=ibdata1:10M:autoextend
innodb_buffer_pool_size = 2048M
innodb_additional_mem_pool_size = 20M
innodb_data_home_dir = /mnt/ulprod-pa_mysql/
innodb_log_group_home_dir = /mnt/ulprod-pa_mysql/
innodb_file_per_table
innodb_log_file_size = 256M
innodb_log_buffer_size = 8M
innodb_flush_log_at_trx_commit = 1
default-storage_engine=InnoDB
thread_cache_size = 16
slow_query_log = 1
long_query_time = 1
skip-federated
server-id       = 1
#log-bin=mysql-bin
#sync_binlog = 1

[mysqldump]
quick
max_allowed_packet = 16M

[mysql]
no-auto-rehash

[isamchk]
key_buffer = 256M
sort_buffer_size = 256M
read_buffer = 2M
write_buffer = 2M

[myisamchk]
key_buffer = 256M
sort_buffer_size = 256M
read_buffer = 2M
write_buffer = 2M

[mysqlhotcopy]
interactive-timeout
```

Then create the client configuration file to force the use of TCP protocol instead of local socket

Example: `/etc/my.cnf`

```bash
[client]
protocol        = tcp
```

### Initializing the Databases

Create directories for databases and logs

```bash
mkdir /mnt/ulprod-ld_mysql/{databases,logs}
mkdir /mnt/ulprod-pa_mysql/{databases,logs}
```

Initialize the MySQL databases

```bash
/opt/mysql/mysql/scripts/mysql_install_db --user=mysql --ldata=/mnt/ulprod-ld_mysql/databases
/opt/mysql/mysql/scripts/mysql_install_db --user=mysql --ldata=/mnt/ulprod-pa_mysql/databases
```

### Startup Script

Make a copy of the startup script for each instance

```bash
mkdir /opt/mysql/mysql/share/mysql
cp /opt/mysql/mysql/support-files/mysql.server /opt/mysql/mysql/share/mysql/mysql.server-prod
cp /opt/mysql/mysql/support-files/mysql.server /opt/mysql/mysql/share/mysql/mysql.server-dr
```

Edit each of these scripts and modify the line:

```bash
$bindir/mysqld_safe --datadir=$datadir --pid-file=$server_pid_file $other_args >/dev/null 2>&1 &
```

Specify the configuration file to use, and remove "--datadir=$datadir" and "--pid-file=$server_pid_file". For "mysql.server-prod" file:

```bash
$bindir/mysqld_safe --defaults-file=/etc/my-prod.cnf $other_args >/dev/null 2>&1 &
```

And for "mysql.server-dr":

```bash
$bindir/mysqld_safe --defaults-file=/etc/my-dr.cnf $other_args >/dev/null 2>&1 &
```

Also comment out the "datadir=" line:

```bash
#datadir=/var/lib/mysql
```

Define the "server_pid_file=" variable by specifying the PID file defined in the configuration. For example:

```bash
server_pid_file=/mnt/ulprod-pa_mysql/mysql-uat.pid
```

Define the "extra_args=" variable by specifying the configuration file. For example:

```bash
extra_args="-c /etc/my-uat.cnf"
```

### SMF File

Create 2 SMF files with different names and modify the "exec" line to call the correct script:

Example: `/var/svc/manifest/application/database/mysql-prod.xml`

```xml
<?xml version="1.0"?>
<!DOCTYPE service_bundle SYSTEM "/usr/share/lib/xml/dtd/service_bundle.dtd.1">
<!--
    Copyright 2005 Sun Microsystems, Inc.  All rights reserved.
    Use is subject to license terms.
    MySQL.xml : MySQL manifest, Scott Fehrman, Systems Engineer
    updated: 2005-09-16
-->

<service_bundle type='manifest' name='MySQL Prod'>
<service name='application/database/mysql-prod' type='service' version='1'>

<single_instance />

    <dependency
        name='filesystem'
        grouping='require_all'
        restart_on='none'
        type='service'>
        <service_fmri value='svc:/system/filesystem/local' />
    </dependency>

    <exec_method
        type='method'
        name='start'
        exec='/opt/mysql/mysql/share/mysql/mysql.server-prod start'
        timeout_seconds='120' />

    <exec_method
        type='method'
        name='stop'
        exec='/opt/mysql/mysql/share/mysql/mysql.server-prod stop'
        timeout_seconds='120' />

    <instance name='default' enabled='false' />

    <stability value='Unstable' />

    <template>
        <common_name>
            <loctext xml:lang='C'>MySQL Prod RDBMS 5.0.19</loctext>
        </common_name>
        <documentation>
            <manpage title='mysql' section='1' manpath='/opt/mysql/mysql/man' />
        </documentation>
    </template>

</service>
</service_bundle>
```

Example: `/var/svc/manifest/application/database/mysql-dr.xml`

```xml
<?xml version="1.0"?>
<!DOCTYPE service_bundle SYSTEM "/usr/share/lib/xml/dtd/service_bundle.dtd.1">
<!--
    Copyright 2005 Sun Microsystems, Inc.  All rights reserved.
    Use is subject to license terms.
    MySQL.xml : MySQL manifest, Scott Fehrman, Systems Engineer
    updated: 2005-09-16
-->

<service_bundle type='manifest' name='MySQL DR'>
<service name='application/database/mysql-dr' type='service' version='1'>

<single_instance />

    <dependency
        name='filesystem'
        grouping='require_all'
        restart_on='none'
        type='service'>
        <service_fmri value='svc:/system/filesystem/local' />
    </dependency>

    <exec_method
        type='method'
        name='start'
        exec='/opt/mysql/mysql/share/mysql/mysql.server-dr start'
        timeout_seconds='120' />

    <exec_method
        type='method'
        name='stop'
        exec='/opt/mysql/mysql/share/mysql/mysql.server-dr stop'
        timeout_seconds='120' />

    <instance name='default' enabled='false' />

    <stability value='Unstable' />

    <template>
        <common_name>
            <loctext xml:lang='C'>MySQL DR RDBMS 5.0.19</loctext>
        </common_name>
        <documentation>
            <manpage title='mysql' section='1' manpath='/opt/mysql/mysql/man' />
        </documentation>
    </template>

</service>
</service_bundle>
```

### Validating SMF Files

Validate the services with the following command:

```bash
svccfg validate /var/svc/manifest/application/database/mysql-prod.xml
svccfg validate /var/svc/manifest/application/database/mysql-dr.xml
```

If the command returns nothing, everything is perfect!

Next, import the XML scripts:

```bash
svccfg import /var/svc/manifest/application/database/mysql-prod.xml
svccfg import /var/svc/manifest/application/database/mysql-dr.xml
```

To validate, check that the services are present:

```bash
svcs mysql-prod
svcs mysql-dr
```

Finally, set the proper permissions:

```bash
chown -Rf mysql: /mnt/ulprod-pa_mysql
chown -Rf mysql: /mnt/ulprod-ld_mysql
chmod -Rf 700 /mnt/ulprod-pa_mysql
chmod -Rf 700 /mnt/ulprod-ld_mysql
```

### Starting the Services

All that's left is to enable the service:

```bash
svcadm enable mysql-prod
svcadm enable mysql-dr
```
