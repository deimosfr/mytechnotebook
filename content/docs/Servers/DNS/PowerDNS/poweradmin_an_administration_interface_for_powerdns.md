---
weight: 999
url: "/PowerAdmin_\\:_Une_interface_d'administration_pour_PowerDNS/"
title: "PowerAdmin: An Administration Interface for PowerDNS"
description: "How to install and configure PowerAdmin, a PHP interface for managing PowerDNS"
categories: ["MySQL", "Debian", "Linux"]
date: "2012-05-15T15:02:00+02:00"
lastmod: "2012-05-15T15:02:00+02:00"
tags: ["PowerDNS", "DNS", "PHP", "Web Interface", "Administration"]
toc: true
---

![PowerAdmin](/images/poweradmin_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 2.1.6 |
| **Operating System** | Debian 6 |
| **Website** | [PowerAdmin Website](https://www.poweradmin.org) |
| **Last Update** | 15/05/2012 |
{{< /table >}}

## Introduction

PowerAdmin is an administration interface for PowerDNS. It's built with PHP and allows you to manage much more than the basic functionality of PowerDNS.

It's very useful if you have a MySQL backend which is not very practical for adding records/zones or modifying the PowerDNS configuration.

## Installation

We'll need a MySQL server, Apache and PHP:

```bash
aptitude install apache2 php5 php5-mysql libapache2-mod-php5 php5-mcrypt mysql-server
```

Then we'll get the latest version of PowerAdmin:

```bash
cd /var/www
wget --no-check-certificate -O poweradmin.tgz "https://github.com/poweradmin/poweradmin/tarball/v2.1.6"
tar -xzf poweradmin.tgz
mv poweradmin-poweradmin-* poweradmin
rm -f poweradmin.tgz
chown -Rf www-data. poweradmin
```

## Configuration

Use the wizard by pointing to your server: http://<poweradmin_server>/poweradmin/install/

During the installation, it will generate a configuration that you'll need to copy to the right location (`/var/www/poweradmin/inc/config.inc.php`):

```php
<?php

$db_host                = '127.0.0.1';
$db_user                = 'pdns';
$db_pass                = 'password';
$db_name                = 'pdns';
$db_port                = '3306';
$db_type                = 'mysql';
$db_layer               = 'PDO';

$session_key            = '=^EfdfdfdsfsdfezzezcdfdsfezeJ_';

$iface_lang             = 'en_EN';

$dns_hostmaster         = 'deimos.deimos.fr';
$dns_ns1                = 'dns1.deimos.fr';
$dns_ns2                = 'dns2.deimos.fr';

?>
```

Then we'll remove the installation directory once completed:

```bash
rm -Rf /var/www/poweradmin/install
```

Make sure you have the correct MySQL permissions:

```mysql
GRANT SELECT, INSERT, UPDATE, DELETE ON `pdns`.* TO 'pdns'@'localhost';
```

Log in to the interface:

* URL: http://<poweradmin_server>/poweradmin/
* Login: admin
* Password: password (the one you entered in the installer)
