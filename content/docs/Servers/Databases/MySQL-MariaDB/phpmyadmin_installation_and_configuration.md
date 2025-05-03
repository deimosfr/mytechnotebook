---
weight: 999
url: "/PhpMyAdmin_\\:_Installation_et_configuration/"
title: "PhpMyAdmin: Installation and Configuration"
description: "Guide for installing and configuring PhpMyAdmin, a web interface for MySQL database management"
categories: ["MySQL", "Database", "Linux"]
date: "2009-05-09T08:16:00+02:00"
lastmod: "2009-05-09T08:16:00+02:00"
tags:
  [
    "PhpMyAdmin",
    "MySQL",
    "Apache",
    "Lighttpd",
    "Web Interface",
    "Database Management",
  ]
toc: true
---

## Introduction

phpMyAdmin (pronounced "p h p my admin") is a free, user-friendly interface written in PHP for the MySQL DBMS to facilitate the management of MySQL databases on a server, and is distributed under the GNU GPL license.

It is one of the most famous interfaces for managing a MySQL database on a PHP server. Many hosting providers, whether free or paid, offer it, which means the user doesn't have to install it. Among them, we can mention Free or Lycos Multimania.

This practical interface allows you to easily execute many queries without extensive knowledge in the database field, such as table creation, insertions, updates, deletions, and modifications to database structure. This system is very convenient for backing up a database as an .sql file, making it easy to transfer data. It also accepts SQL queries directly in SQL language, which allows you to test your queries, for example when creating a website, saving precious time.

## Installation

To install it, you must first have installed [MySQL Server]({{< ref "docs/Servers/Databases/MySQL-MariaDB/mysql_installation_and_configuration.md" >}})! Once that's done, let's install phpMyAdmin:

```bash
apt-get install phpmyadmin
```

## Configuration

By default, no configuration is needed if you have installed it on the same machine as MySQL Server. Otherwise, you need to specify the address you want to connect to.

### config.inc.php

#### Mono server

Edit the file `/etc/phpmyadmin/config.inc.php`, then modify this line:

```bash
//$cfg['Servers'][$i]['host']          = 'localhost'; // MySQL hostname or IP address
```

to:

```bash
$cfg['Servers'][$i]['host']          = 'IP_OF_SERVER'; // MySQL hostname or IP address
```

Uncomment the line by removing the "//" and put the IP of the server you're interested in.

Then, you must authorize a person to connect remotely to the server. For simplicity here, I chose root who has access to all databases, but choosing another username is not a bad idea! On the server, adapt and execute these lines:

```bash
create user 'root'@'IP_OF_PHPMYADMIN' identified by '**********';
grant all privileges on * . * to 'root'@'IP_OF_PHPMYADMIN' identified by '**********' with grant option max_queries_per_hour 0 max_connections_per_hour 0 max_updates_per_hour 0 max_user_connections 0;
grant all privileges on `root_%` . * to 'root'@'IP_OF_PHPMYADMIN';
flush privileges;
```

Replace "IP_OF_PHPMYADMIN" with the IP address of the machine where phpMyAdmin is located, then replace "****\*\*****" with the password you want to assign to the user.

Now, connect via the web interface and you'll have permission to connect.

#### Multi servers

For a basic multi-server configuration, edit the file `/etc/phpmyadmin/config.inc.php`, then add these lines:

```bash
// Add server 2
$i++;
$cfg['Servers'][$i]['host']          = 'server2';

// Add server 3
$i++;
$cfg['Servers'][$i]['host']          = 'server3';
```

Then follow the same procedure for each server as in the single server section above.

### Permanent 80% font size

In the latest versions of phpMyAdmin, the font is quite large (pseudo web2.0 effect?!), which is not very cool for smaller resolutions, and anyway... it's simply ugly!

A small CSS fix will solve this. Edit the left.php file at the root and the header.inc.php file in the libraries. Add the following code in the head tags of these pages:

```html
<style type="text/css">
  <!--
  body { font-size: 0.8em ! important; }
  -->
</style>
```

### With Lighttpd

For a phpMyAdmin configuration with Lighttpd, here's what you need (`/etc/lighttpd/conf-available/50-phpmyadmin.conf`):

```perl
# Alias for phpMyAdmin directory
alias.url += (
    "/phpmyadmin" => "/usr/share/phpmyadmin",
)

# Disallow access to libraries
$HTTP["url"] =~ "^/phpmyadmin/libraries" {
    url.access-deny = ( "" )
}

# Limit access to setup script
$HTTP["url"] =~ "^/phpmyadmin/scripts/setup.php" {
    auth.backend = "plain"
    auth.backend.plain.userfile = "/etc/phpmyadmin/htpasswd.setup"
    auth.require = (
        "/" => (
            "method" => "basic",
            "realm" => "phpMyAdmin Setup",
            "require" => "valid-user"
        )
    )
}
```
