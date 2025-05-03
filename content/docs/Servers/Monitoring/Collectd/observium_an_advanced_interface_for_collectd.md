---
weight: 999
url: "/Observium_\\:_Une_interface_évoluée_pour_Collectd/"
title: "Observium: An Advanced Interface for Collectd"
description: "Learn how to set up Observium to provide an advanced graphical interface for Collectd monitoring data, including installation, configuration, and integration with SNMP."
categories: ["RHEL", "Debian", "Database"]
date: "2012-03-22T15:56:00+02:00"
lastmod: "2012-03-22T15:56:00+02:00"
tags: ["Collectd", "SNMP", "Monitoring", "RRD", "Apache", "MySQL"]
toc: true
---

![Observium](/images/observium-logo.avif)

## Introduction

[Observium](https://www.observium.org) is a graphically oriented monitoring tool similar to [Munin](). The advantage of this solution is that it manages [RRD]() graphs from [Collectd]().

Since today, a high-performance graphical interface is really something that's missing from [Collectd](), it would be a shame to have such a powerful solution without an equally capable interface. That's why after much research, I finally came across [Observium](https://www.observium.org) which meets this need.

## Installation

### Prerequisites

We'll use Debian 6 for this tutorial. Here is the list of required packages (use non-free repositories for the snmp-mibs-downloader package):

```bash
aptitude install libapache2-mod-php5 php5-cli php5-mysql php5-gd php5-snmp php-pear snmp graphviz subversion mysql-server mysql-client rrdtool fping imagemagick whois mtr-tiny nmap ipmitool snmp-mibs-downloader php-net-ipv4
```

Then we'll install a necessary PHP IPv6 library that unfortunately isn't packaged:

```bash
pear install Net_IPv6
```

### Observium

Now we'll download and extract Observium, then create some necessary folders:

```bash
cd /var/www/
wget http://www.observium.org/observium-latest.tar.gz
tar -xzf observium-latest.tar.gz
rm -f observium-latest.tar.gz
cd /var/www/observium
mkdir graphs rrd
chown -Rf www-data. .
```

We create a basic configuration file:

```bash
cp config.php.default config.php
```

### MySQL

Let's move on to the database. We need to create it with a user (don't forget to replace the password):

```bash
> mysql -uroot -p
mysql> CREATE DATABASE observium;
mysql> GRANT ALL PRIVILEGES ON observium.* TO 'observium'@'localhost' IDENTIFIED BY '<password>';
mysql> flush privileges;
```

Let's change the Observium location in the configuration file, add an entry for fping, and modify the SQL fields to adapt them with the correct data (`/var/www/observium/config.php`):

```php {linenos=table,hl_lines=[6,7,12,13]}
<?php

## Have a look in defaults.inc.php for examples of settings you can set here. DO NOT EDIT defaults.inc.php!

### Database config
$config['db_host'] = "localhost";
$config['db_user'] = "USERNAME";
$config['db_pass'] = "PASSWORD";
$config['db_name'] = "observium";

### Locations
$config['fping'] = "/usr/bin/fping";
$config['install_dir']  = "/var/www/observium";
...
```

Next, we'll import the schema for building the database:

```bash
mysql -uobservium -p observium < database-schema.sql
```

And we'll update the database with some updates:

```bash
scripts/update-sql.php database-update-pre1000.sql
scripts/update-sql.php database-update-pre1435.sql
scripts/update-sql.php database-update-pre2245.sql
scripts/update-sql.php database-update.sql
```

### Apache

Then we'll create the Apache configuration for this new site (don't forget to make the DNS record) (`/etc/apache2/sites-available/observium`):

```apache {linenos=table,hl_lines=[3]}
<VirtualHost *:80>
    DocumentRoot /var/www/observium/html/
    ServerName  observium.mycompany.com
    CustomLog /var/log/apache2/access.log combined
    ErrorLog /var/log/apache2/error.log
    <Directory "/var/www/observium/html/">
        Options Indexes FollowSymLinks MultiViews
        AllowOverride All
        Order allow,deny
        allow from all
    </Directory>
</VirtualHost>
```

Then we activate this new configuration, make sure the rewrite module is active, and reload Apache:

```bash
a2enmod rewrite
a2ensite observium
/etc/init.d/apache2 reload
```

### SNMP

For SNMP, we'll edit the configuration file to add a line containing information about the MIBs (`/etc/snmp/snmp.conf`):

```bash {linenos=table,hl_lines=[5,6]}
#
# As the snmp packages come without MIB files due to license reasons, loading
# of MIBs is disabled by default. If you added the MIBs you can reenable
# loaging them by commenting out the following line.
# mibs :
mibdirs /var/www/observium/mibs
```

### Crontab

You'll need to [create a user](#add-a-user) such as admin for example and [add a host](#add-a-machine-to-monitor) before initializing and retrieving the first data.

Once done, we run the data collection for the first time:

```bash
./discovery.php -h all
./poller.php -h all
```

And finally we create a cron file for this service (`/etc/cron.d/observium`):

```bash
33 */6 * * *   root    cd /var/www/observium/ && ./discovery.php -h all >> /dev/null 2>&1
*/5 * * * *   root    cd /var/www/observium/ && ./discovery.php -h new >> /dev/null 2>&1
*/5 * * * *   root    cd /var/www/observium/ && ./poller.php -h all >> /dev/null 2>&1
```

And finally, we reload the service:

```bash
/etc/init.d/cron reload
```

## Configuration

### Add a User

This is normally only for the admin. After that, it's better to have an LDAP backend:

```bash
/var/www/observium/adduser.php <username> <password> 10
```

- 10: This is the permission level going from 0 to 10 (10 being admin).

### Add a Machine to Monitor

Here's how to add a machine to monitor via SNMP:

```bash
/var/www/observium/addhost.php <hostname> <community> v2c
```

Here's an example:

```bash
> /var/www/observium/addhost.php localhost public v2c
Created host : localhost (id:1) (os:linux)
```

### Collectd

If you want to add [Collectd]() to your server graphs, it's possible. You need to edit the server configuration and add this:

```bash
...
### Locations
$config['collectd_dir'] = '/var/lib/collectd/rrd';
...
```

On the client side, you need to ensure that the Hostname variable is properly configured, such as:

```bash
Hostname "localhost"
```

All servers existing in Collectd must exist on Observium. That's why we're going to automatically add all machines:

```bash
cd /var/lib/collectd/rrd/ ; for i in * ; do /var/www/observium/addhost.php $i lookullink v2c ; done
```

#### Uppercase Hostname Hack

Even though I contacted the Observium developers, you'll have issues if you want to integrate Collectd with hostnames that have uppercase names. Indeed, since it makes sense to have lowercase hostnames, they decided not to handle uppercase. And if, like me, your hostnames in Collectd are uppercase, it's impossible to perform an integration. I searched for 1000 ways to properly work around this problem, and each time it required code rewriting rather than configuration on the Collectd or Observium side.

So, I ended up modifying the code, but very, very, very slightly. In fact, I simply took care of the host creation script so it creates symbolic links (I know it's not super clean) (`addhost.php`):

```php {linenos=table,hl_lines=["24-26","56-64"]}
#!/usr/bin/env php
<?php

/* Observium Network Management and Monitoring System
 * Copyright (C) 2006-2011, Observium Developers - http://www.observium.org
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * See COPYING for more details.
 */

include("includes/defaults.inc.php");
include("config.php");
include("includes/functions.php");

if (isset($argv[1]) && $argv[1])
{
  $host    = strtolower($argv[1]);
  $community = $argv[2];
  $snmpver   = strtolower($argv[3]);
  $host_upper = strtoupper($host);
  $host_fullpath = $config['collectd_dir'] . '/' . $host;
  $host_fullpath_upper = $config['collectd_dir'] . '/' . $host_upper;

  if (is_numeric($argv[4]))
  {
    $port = $argv[4];
  }
  else
  {
    $port = 161;
  }

  if (@!$argv[5])
  {
    $transport = 'udp';
  }
  else
  {
    $transport = $argv[5];
  }

  if (!$snmpver) $snmpver = "v2c";

  if ($community)
  {
    unset($config['snmp']['community']);
    $config['snmp']['community'][] = $community;
  }

  addHost($host, $community, $snmpver, $port = '161', $transport = 'udp');

    # Add symlink to correct uppercase hostname problems
    if (is_dir($host_fullpath) or is_link($host_fullpath))
    {
        exit(0);
    }
    elseif (is_dir($host_fullpath_upper))
    {
        symlink($host_fullpath_upper, $host_fullpath);
    }

} else {

print Console_Color::convert("

Observium v".$config['version']." Add Host Tool

Usage: ./addhost.php <%Whostname%n> [community] [v1|v2c] [port] [" . join("|",$config['snmp']['transports']) . "]

%rRemeber to discover the host afterwards.%n

");
}

?>
```

## Resources
- http://www.observium.org/wiki/Ubuntu_SVN_Installation
- http://www.outsidaz.org/blog/2012/02/09/deploying-observium-on-rhel6-with-selinux/
