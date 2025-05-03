---
weight: 999
url: "/PhpPgAdmin_\\:_Installation_et_configuration/"
title: "PhpPgAdmin: Installation and Configuration"
description: "Learn how to install and configure PhpPgAdmin to manage PostgreSQL databases through a web interface, including multi-server setup."
categories: ["PostgreSQL", "Linux", "Apache"]
date: "2007-08-10T09:03:00+02:00"
lastmod: "2007-08-10T09:03:00+02:00"
tags:
  ["phppgadmin", "postgresql", "database", "web interface", "administration"]
toc: true
---

## Introduction

PhpPgAdmin is equivalent to phpMyAdmin for those who know it. It allows you to administer PostgreSQL through a small web interface.

## Installation

To install it:

```bash
apt-get install phppgadmin
```

Then, to choose which web server to run on (otherwise we have to do it manually):

```bash
dpkg-reconfigure phppgadmin
```

## Configuration

To allow connections from an external machine, modify the file `/etc/phppgadmin/apache.conf` and add:

```bash
allow from all
```

This is not the ideal configuration; it's better to add only your IP:

```bash
allow from 192.168.0.2
```

This is a bit better, but may not necessarily meet your needs.

Now restart Apache:

```bash
/etc/init.d/apache2 restart
```

For security reasons, the "postgres" login is not allowed to access the database via "phppgadmin". To allow it, you need to modify the file `/etc/phppgadmin/config.inc.php`. In this file, set "false" on the following line:

```php
$conf['extra_login_security'] = false
```

You also need to add the local IP address on the following line:

```php
$conf['servers'][0]['host'] = '127.0.0.1';
```

### Multi-Server

In a multi-server configuration, here are the lines to add:

```php
$conf['servers'][1]['desc'] = 'Server2'; // Don't forget to change the number (1) to indicate the server, then the server name to display
$conf['servers'][1]['host'] = '192.168.0.87'; // The IP of the server
$conf['servers'][1]['port'] = 5432;
$conf['servers'][1]['sslmode'] = 'allow';
$conf['servers'][1]['defaultdb'] = 'template1';
$conf['servers'][1]['pg_dump_path'] = '/usr/bin/pg_dump';
$conf['servers'][1]['pg_dumpall_path'] = '/usr/bin/pg_dumpall';
$conf['servers'][1]['slony_support'] = false;
$conf['servers'][1]['slony_sql'] = '/usr/share/postgresql';

$conf['servers'][2]['desc'] = 'Server3';
$conf['servers'][2]['host'] = '192.168.0.89';
$conf['servers'][2]['port'] = 5432;
$conf['servers'][2]['sslmode'] = 'allow';
$conf['servers'][2]['defaultdb'] = 'template1';
$conf['servers'][2]['pg_dump_path'] = '/usr/bin/pg_dump';
$conf['servers'][2]['pg_dumpall_path'] = '/usr/bin/pg_dumpall';
$conf['servers'][2]['slony_support'] = false;
$conf['servers'][2]['slony_sql'] = '/usr/share/postgresql';
```
