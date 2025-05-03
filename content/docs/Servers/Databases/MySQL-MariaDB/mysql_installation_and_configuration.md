---
weight: 999
url: "/MySQL_\\:_Installation_et_configuration/"
title: "MySQL: Installation and Configuration"
description: "A comprehensive guide for MySQL installation, configuration, and common management tasks including user management, database operations, and troubleshooting."
categories: ["Database", "MySQL", "Server"]
date: "2014-04-09T16:39:00+02:00"
lastmod: "2014-04-09T16:39:00+02:00"
tags: ["MySQL", "Database", "SQL", "Server Configuration"]
toc: true
---

![MySQL Logo](/images/mysql_logo.avif)

## Introduction

[MySQL](https://fr.wikipedia.org/wiki/Mysql) is a relational SQL database server developed with a focus on high performance. It is multi-threaded, robust, and multi-user. It is open-source software developed under a dual license depending on its use: in an open-source product or in a proprietary product. In the latter case, the license is paid; otherwise, it's free.

## Installation

To install it, nothing could be simpler:

```bash
apt-get install mysql
```

## Usage

I strongly recommend a small built-in MySQL utility to configure MySQL simply and securely. Let's start:

```bash
$ mysql_secure_installation

NOTE: RUNNING ALL PARTS OF THIS SCRIPT IS RECOMMENDED FOR ALL MySQL
      SERVERS IN PRODUCTION USE!  PLEASE READ EACH STEP CAREFULLY!


In order to log into MySQL to secure it, we'll need the current
password for the root user.  If you've just installed MySQL, and
you haven't set the root password yet, the password will be blank,
so you should just press enter here.

Enter current password for root (enter for none):
```

**Just press "Enter"** because there is no default password

```
OK, successfully used password, moving on...

Setting the root password ensures that nobody can log into the MySQL
root user without the proper authorisation.

You already have a root password set, so you can safely answer 'n'.

Change the root password? [Y/n]
```

**Answer "y"** and change the password

```
By default, a MySQL installation has an anonymous user, allowing anyone
to log into MySQL without having to have a user account created for
them.  This is intended only for testing, and to make the installation
go a bit smoother.  You should remove them before moving into a
production environment.

Remove anonymous users? [Y/n]
```

**Answer "y"** to remove anonymous users

```
Normally, root should only be allowed to connect from 'localhost'.  This
ensures that someone cannot guess at the root password from the network.

Disallow root login remotely? [Y/n]
```

**Answer "y"** to not allow root to connect remotely

```
By default, MySQL comes with a database named 'test' that anyone can
access.  This is also intended only for testing, and should be removed
before moving into a production environment.

Remove test database and access to it? [Y/n]
```

**Answer "y"** to remove the test database

```
Reloading the privilege tables will ensure that all changes made so far
will take effect immediately.

Reload privilege tables now? [Y/n]
```

**And finally answer "y"** to reload the privilege tables

```
Cleaning up...



All done!  If you've completed all of the above steps, your MySQL
installation should now be secure.

Thanks for using MySQL!
```

Everything is finished; you can now use your database.

### Creating a database

To create a database, here is the command to execute:

```bash
create database database_name;
```

### Creating a user

To create a user:

```bash
create user 'user'@'localhost' identified by 'password';
GRANT USAGE ON * . * TO 'user'@'localhost' IDENTIFIED BY 'password';
grant SELECT, INSERT, UPDATE, DELETE on `database_name` .* to 'user'@'localhost';
flush privileges;
```

### Changing a user's password

To change a user's password, it's simple:

```mysql
SET PASSWORD FOR 'debian-sys-maint'@'localhost' = PASSWORD('newpass');
```

This command includes flush privileges, so you don't need to type it afterward :-)

### Modifying user rights

If, for example, I want to change the connection hostname for all users:

```sql
UPDATE mysql.USER SET host = '10.0.0.%' WHERE host = 'localhost' AND USER != 'root';
UPDATE mysql.db SET host = '10.0.0.%' WHERE host = 'localhost' AND USER != 'root';
FLUSH PRIVILEGES;
```

### Deleting a user

Before deleting a user, it is advisable to list the privileges to revoke current rights:

```sql
SHOW grants FOR <user>;
```

Then revoke:

```sql
REVOKE ALL privileges FROM <user>;
```

And finally, delete the user:

```sql
DROP USER 'user'@'localhost';
```

### Listing running processes

```mysql
mysql> show processlist;
+----+-------------+-----------------+------+---------+------+-----------------------------------------------------------------------------+------------------+
| Id | User        | Host            | db   | Command | Time | State                                                                       | Info             |
+----+-------------+-----------------+------+---------+------+-----------------------------------------------------------------------------+------------------+
|  1 | system user |                 | NULL | Connect |  601 | Connecting to master                                                        | NULL             |
|  2 | system user |                 | NULL | Connect |  601 | Slave has read all relay log; waiting for the slave I/O thread to update it | NULL             |
|  8 | root        | localhost:36538 | NULL | Query   |    0 | NULL                                                                        | show processlist |
+----+-------------+-----------------+------+---------+------+-----------------------------------------------------------------------------+------------------+
3 rows in set (0.00 sec)
```

### Renaming a database

```bash
#!/bin/sh
olddb=olddb
newdb=newdb
user=root
pass='password'
port=3306

echo "######### COPY / PASTE THOSE LINES TO RENAME DATABASE #########"
echo ""
gettables=`mysql -u\$user \$pass -P\$port -e "show tables from \$olddb;" | grep -v "Tables_in_\$olddb" | grep -v "+" | grep -v "^mysql" | awk '{print \$1}'`
for i in `echo \$gettables | tr '\n' ' '`  ; do
echo "RENAME TABLE \$olddb.\$i TO \$newdb.\$i;"
done
mysql -u\$user \$pass -P\$port -e "use \$newdb" > /dev/null 2>&1 || echo "############################ WARNING ##########################\nTHE NEW DATABASE \"\$newdb\" DOES NOT EXIST. please create it first\n###############################################################"
```

Run this script and it will rename table by table to finally create the new database. This is the method recommended by MySQL.

### Knowing the size of a database

To get the size of all databases in MB:

```sql
SELECT table_schema,round(SUM(data_length+index_length)/1024/1024,4) AS "Size (MB)"
FROM information_schema.TABLES
GROUP BY table_schema;
```

If you want the size of a single table, you need to specify the table name (table_name field) and the database name (table_schema field):

```sql
SELECT table_schema,round(SUM(data_length+index_length)/1024/1024,4)
FROM information_schema.TABLES
WHERE table_schema = 'mysql'
AND TABLE_NAME = 'user';
```

## Connecting with default credentials

It can be useful to be able to connect simply without having to enter credentials. Here is a very simple method that consists of entering your credentials in a file in your home:

```ini
[client]
user=root
password=password
```

Then apply the right permissions:

```bash
chmod 600 ~/.my.cnf
```

Connect without credentials :-)

## FAQ

### How to reset your root password when you've lost it?

Have you ever forgotten the root password on one of your MySQL servers? No? Well maybe I'm not as perfect as you. This is a quick h00tow (how to) reset your MySQL root password. It does require root access on your server. If you have forgotten that password wait for another article:

First things first. Log in as root and stop the mysql daemon. Now let's start up the mysql daemon and skip the grant tables which store the passwords.

```bash
mysqld_safe --skip-grant-tables --skip-networking &
```

You should see mysqld start up successfully. If not, well you have bigger issues. Now you should be able to connect to mysql without a password.

```bash
$ mysql --user=root mysql

update user set Password=PASSWORD('new-password') WHERE user = 'root';
flush privileges;
exit;
```

Now kill your running mysqld, then restart it normally. You should be good to go. Try not to forget your password again.

### Fatal error: Can't open and lock privilege tables: Table 'mysql.host' doesn't exist

If you get this kind of message when booting mysql:

```
Fatal error: Can't open and lock privilege tables: Table 'mysql.host' doesn't exist
```

You need to rebuild the missing databases like this:

```bash
mysql_install_db --user=mysql --ldata=/new-data-location
mysqld_safe --datadir=/new-data-location --user=mysql &
```

## Resources
- [Optimising MySQL under Sun](/pdf/optimising_mysql_under_sun.pdf)
- [MySQL Utils: beautiful cacti graphs for Monitoring MySQL](https://faemalia.net/mysqlUtils/)
- [Setting Changing And Resetting MySQL Root Passwords](/pdf/setting_changing_and_resetting_mysql_root_passwords.pdf)
- [Monolith-toolkit: easy tools for complex MySQL administration](https://code.google.com/p/monolith-toolkit/)
