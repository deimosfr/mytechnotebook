---
weight: 999
url: "/Php-syslog-ng_\\:_Interprétation_des_logs_Syslog-ng_dans_une_interfaçe_web/"
title: "PHP-syslog-ng: Interpretation of Syslog-ng logs in a web interface"
description: "How to install and configure PHP-syslog-ng to analyze and interpret syslog logs in a web interface"
categories: ["Linux", "Apache", "MySQL"]
date: "2008-07-23T14:03:00+02:00"
lastmod: "2008-07-23T14:03:00+02:00"
tags: ["Servers", "Network", "Monitoring", "Logs"]
toc: true
---

## Introduction

[php-syslog-ng](https://code.google.com/p/php-syslog-ng/) is a web application that allows you to format, search, and interpret logs. For searching, it requires logs to be in an SQL database, and for log interpretation, it's specifically designed for Cisco logs.

Note: Before continuing, you'll need a web server like Apache with PHP module installed. You'll also need the MySQL module for PHP.

## Installation

Let's use the latest version:

```bash
cd /var/www
wget http://php-syslog-ng.googlecode.com/files/php-syslog-ng-2.9.8.tgz
```

Now let's extract it:

```bash
tar -xzvf php-syslog-ng-2.9.8.tgz
```

If we want graphs, we need to install Microsoft fonts:

```bash
apt-get install msttcorefonts
```

## Configuration

Simply go to the page [http://localhost/php-syslog-ng/html/](http://localhost/php-syslog-ng/html/) and fill in the correct information. The installer will prepare your MySQL database. After installation, you can delete the "installation" folder and modify the configuration file located at `/var/www/php-syslog-ng/html/config/config.php` whenever you want.

### Log Rotation

Let's edit the configuration file `/var/www/php-syslog-ng/html/config/config.php` and modify this line to keep 6 months of logs:

```php
define('LOGROTATERETENTION', 180);
```

Here, we want all logs older than 180 days to be deleted. However, we need to make sure this is configured in root's crontab. Edit it and add these lines:

```bash
@daily php /var/www/php-syslog-ng/scripts/logrotate.php >> /var/log/php-syslog-ng/logrotate.log
@daily find /var/www/php-syslog-ng/html/jpcache/ -atime 1 -exec rm -f '{}' ';'
*/5 * * * * php /var/www/php-syslog-ng/scripts/reloadcache.php >> /var/log/php-syslog-ng/reloadcache.log
```

Also create the log folder if it doesn't exist:

```bash
mkdir -p /var/log/php-syslog-ng/
```
