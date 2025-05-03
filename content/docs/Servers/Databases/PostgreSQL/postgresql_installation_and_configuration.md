
---

weight: 999
url: "/Installation_et_configuration_de_PostgreSQL/"
title: "PostgreSQL Installation and Configuration"
description: "A comprehensive guide to installing, configuring, and managing PostgreSQL databases including user management, backups, and basic SQL operations."
categories: ["PostgreSQL", "Database"]
date: "2012-11-16T08:55:00+02:00"
lastmod: "2012-11-16T08:55:00+02:00"
tags: ["PostgreSQL", "Database", "SQL", "Backup", "User Management"]
toc: true

---

![PostgreSQL](/images/bases_de_donnees_icon.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 8.3+ |
| **Operating System** | Debian 6 |
| **Website** | [PostgreSQL Website](https://www.postgresql.org/) |
| **Last Update** | 16/11/2012 |
{{< /table >}}

## Introduction

PostgreSQL is an object-relational database management system (ORDBMS). It is a free tool available under a BSD-like license.

This system competes with other database management systems, both free (like MySQL and Firebird) and proprietary (like Oracle, Sybase, DB2, and Microsoft SQL Server). Like the Apache and Linux projects, PostgreSQL is not controlled by a single company but is based on a global community of developers and companies.

## Installation

To install PostgreSQL:

```bash
apt-get install postgresql postgresql-client
```

This creates the postgres user who has rights over the database.
To initialize a PostgreSQL database in the postgres user's $HOME (/var/lib/postgres on Debian) (the installation on Debian does this automatically):

```bash
> su postgres
> cd
> /usr/lib/postgresql/8.3/bin/initdb -D data
The files in this cluster will belong to user "postgres".
The server process must also be owned by this user.

The cluster will be initialized with locale fr_FR@euro.
The default database encoding has been set accordingly to LATIN9.

creating directory data... ok
creating subdirectories... ok
selecting default max_connections... 100
selecting default shared_buffers/max_fsm_pages... 24MB/153600
creating configuration files... ok
creating template1 database in data/base/1... ok
initializing pg_authid... ok
initializing dependencies... ok
creating system views... ok
loading system objects' descriptions... ok
creating conversions... ok
initializing access privileges on built-in objects... ok
creating information schema... ok
vacuuming database template1... ok
copying template1 to template0... ok
copying template1 to postgres... ok

WARNING: enabling "trust" authentication for local connections
You can change this by editing pg_hba.conf or using the -A option the
next time you run initdb.

Success. You can now start the database server using:

    /usr/lib/postgresql/bin/postgres -D data
or
    /usr/lib/postgresql/bin/pg_ctl -D data -l logfile start
```

Note: You need to have correct permissions on /tmp/

On Debian 5 (lenny), to initialize the main postgres database, run the following command as root:

```bash
pg_createcluster 8.3 main
```

Configuration files are located in $HOME/data, particularly pg_hba.conf for access rights management and postgresql.conf for general service configuration. On Debian, these files are symbolic links to files in /etc/postgresql/.

## Configuration

### User Management

#### Admin

First, let's change the password:

```bash
passwd postgres
```

This will modify the postgres password on the machine, but not on the database. You don't have to do this if you don't need to (e.g., if you always go through root then postgres).
Now, let's define a password for the database:

```bash
psql -d template1 -c "alter user postgres with password 'password'"
```

or

```bash
\password postgres
```

This will make the user admin of the template1 database.

Note: To use web clients or graphical clients presented in the following chapters, you must define a password for "postgres".

#### Authentication

By default, to access PostgreSQL, you need to connect as the "postgres" user.
To create a new user, you would need to create a system account for them first, which may not be desirable.
To change this configuration, modify the "/etc/postgresql/pg_hba.conf" file and replace "ident sameuser" with "password" on the following two lines:

```bash
[...]
local  all     all                                        ident sameuser
host   all     all    127.0.0.1         255.255.255.255   ident sameuser
```

This becomes:

```bash
[...]
local  all     all                                        password
host   all     all    127.0.0.1         255.255.255.255   password
```

**This modification avoids using system accounts and requires a password for each connection.**

To allow another machine to connect to this Postgres server, add a line like this to the "/etc/postgresql/pg_hba.conf" file:

```bash
host   all     all    192.168.0.1       255.255.255.255   password
```

{{< alert context="warning" text="WARNING: The order of insertion of the lines is very important! The first line that matches will be the one that takes the rule. So be careful about the insertion order in this file." />}}

And modify the "/etc/postgresql/postgresql.conf" file by adding the following option so that postgres listens on all its addresses and not just "localhost":

```bash
listen_addresses = '*'
```

Restart PostgreSQL:

```bash
/etc/init.d/postgresql restart
```

#### Creation

To create a user:

```bash
> createuser toto
Is the new role a superuser? (y/n) n
Can the new user create databases? (y/n) y
Can the new user create users? (y/n) n
CREATE USER
```

By default, this user has no password. To assign one:

```bash
psql -d template1 -c "alter user toto with password <password>"
```

#### Suppression

To delete a user:

```bash
drop user toto
```

### Database Management

#### Creation

The following command creates the "mybase" database for the user "toto" using the "UNICODE" encoding:

```sql
CREATE DATABASE <database_name> owner <username>;
```

Be careful with the table encoding: LATIN9, LATIN1, UNICODE, etc.

#### Suppression

To delete a database:

```sql
DROP DATABASE <mybase>
```

#### Listing

To list existing databases:

```bash
psql -l
```

#### Backing Up a Database

To back up a database:

```bash
pg_dump DATABASE_NAME > FILE_NAME
```

#### Backing Up All Databases

To back up all databases at once:

```bash
pg_dumpall > FILE_NAME
```

#### Backups on Very Large Databases

It seems that by passing certain parameters, it is easier to restore databases:

```bash
pg_dump -Ft DATABASE_NAME > FILE_NAME
```

[https://www.postgresql.org/docs/8.1/static/app-pgdump.html](https://www.postgresql.org/docs/8.1/static/app-pgdump.html)
Note: Be aware of some limitations in "hot" backups.

#### Restoration

To restore a database, [first create an empty database](#creation), then import:

```bash
psql <newbase> < FILE_NAME
```

### Displaying All Queries

You can increase the log level to see all queries that pass through. These queries will be recorded in a file.

{{< alert context="warning" text="This can slow down your system if you have a lot of activity" />}}

To set up a higher log level, edit these lines in the PostgreSQL configuration:

```bash
[...]
logging_collector = on
log_directory = 'pg_log'
log_filename = 'postgresql-%Y-%m-%d_%H%M%S.log'
log_statement = 'all'
[...]
```

Just restart your PostgreSQL to see logs appear in /var/lib/postgresql/<postgresql_version>/main/pg_log

## Usage

You can now "use" your database with the PostgreSQL command line client:

```bash
> psql -h hostname -U user -d database
mybase=#
```

Here are some useful commands to remember:

```
\l = list databases
\d = list tables
\q = quit
\h = help
SELECT version(); = PostgreSQL version
SELECT current_date; = current date
\i file.sql = read instructions from file.sql
\d table = describe a table (like DESCRIBE in MySQL)
```

### Creating and Deleting Tables

Here are the different data types for table fields:

```sql
CHAR(n)
VARCHAR(n)
INT
REAL
DOUBLE PRECISION
DATE
TIME
TIMESTAMP
INTERVAL
```

Note: You can also define your own data types.

The creation syntax:

```sql
CREATE TABLE my_table (col1 TYPE, [...], coln TYPE);
```

Deletion:

```sql
DROP TABLE my_table;
```

As a small example taken from the PostgreSQL documentation:

```sql
CREATE TABLE weather (
    city            VARCHAR(80),
    temp_lo         INT,           -- low temperature
    temp_hi         INT,           -- high temperature
    prcp            REAL,          -- precipitation
    DATE            DATE
);
```

{{< alert context="info" text="Two dashes -- introduce comments..." />}}

### Data Extraction

Nothing beats examples:

```sql
SELECT * FROM weather;
SELECT city, (temp_hi+temp_lo)/2 AS temp_avg, DATE FROM weather;
SELECT * FROM weatherWHERE city = 'San Francisco' AND prcp > 0.0;
SELECT DISTINCT city FROM weather ORDER BY city;
```

With joins:

```sql
SELECT * FROM weather, cities WHERE city = name;
SELECT weather.city, weather.temp_lo, cities.location 
FROM weather, cities WHERE cities.name = weather.city;
SELECT * FROM weather INNER JOIN cities ON (weather.city = cities.name);
SELECT * FROM weather LEFT OUTER JOIN cities ON 
(weather.city = cities.name);
SELECT * FROM weather w, cities c WHERE w.city = c.name;
```

With functions (Aggregate Functions):

```sql
SELECT MAX(temp_lo) FROM weather;
```

Note that "Aggregate Functions" cannot be used in the WHERE clause.
So the following query is incorrect:

```sql
SELECT city FROM weather WHERE temp_lo = MAX(temp_lo);
```

You should do instead:

```sql
SELECT city FROM weather WHERE 
temp_lo = (SELECT MAX(temp_lo) FROM weather);
```

You can, of course, use "GROUP BY ...", "HAVING ...", etc.

### Data Updates

Still with an example:

```sql
UPDATE weather SET temp_hi = temp_hi - 2, 
temp_lo = temp_lo - 2 WHERE DATE > '1994-11-28';
```

### Data Deletion

Again with an example:

```sql
DELETE FROM weather WHERE city = 'Hayward';
```

To delete all data from a table:

```sql
DELETE FROM weather;
```

### Querying a Database Size

```sql
SELECT pg_size_pretty(pg_database_size('database_name'));
```

## FAQ

### Tsearch

If an application requires tsearch2 for example, you need to install a package:

```bash
aptitude install postgresql-contrib-8.2
```

Use your corresponding version number at the end of the package.

And finally, you need to patch the database in question:

```bash
psql wikidb < /usr/share/postgresql/8.2/contrib/tsearch2.sql
```

Here wikidb is my database, and 8.2 is the postgres version again.

### UTF8 Encoding Doesn't Match the Locale

The annoying error that bugged me. This happens when you want to create a database that's not in the server's current encoding:

```
createdb: database creation failed: ERROR: UTF8 encoding does not match locale fr_FR@euro of the server
DETAIL: The LC_CTYPE server parameter requires the LATIN9 encoding.
```

To fix the problem, first check that the locales are good at the OS level:

```bash
dpkg-reconfigure locales
```

Then, let's say I want UTF-8, I set my environment variables properly:

```bash
export LC_ALL=fr_FR.UTF-8
export LANG=fr_FR.UTF-8
export PGCLIENTENCODING=fr_FR.UTF-8
```

Now I can list the postgres environment variable:

```bash
$ postgres=# show lc_ctype;
  lc_ctype   
 -------------
 fr_FR.UTF-8
(1 line)
```

Normally if you create a new database, the encoding will now be UTF-8. If that still doesn't change anything, you must run this, but it will put all your future databases in this format:

```bash
initdb -E UTF-8
```

Now try creating a database again.

## Resources
- [PostgreSQL Official Website](https://www.postgresql.org/)
- [https://www.postgresql.org/docs/7.3/static/multibyte.html](https://www.postgresql.org/docs/7.3/static/multibyte.html)
- [System Views in PostgreSQL 8.3](/pdf/vues_system_postgres_unixgarden.pdf)
