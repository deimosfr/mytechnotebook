---
weight: 999
url: "/Selfoss_\\:_An_elegant_RSS_reader/"
title: "Selfoss: An Elegant RSS Reader"
description: "A guide to install and configure Selfoss, an elegant self-hosted RSS reader solution, on Debian with Nginx and MariaDB."
categories: ["Nginx", "Debian", "Database"]
date: "2013-06-14T21:13:00+02:00"
lastmod: "2013-06-14T21:13:00+02:00"
tags: ["RSS", "Selfoss", "Nginx", "MariaDB", "Database", "Web server"]
toc: true
---

![Selfoss](/images/selfoss_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 2.7 |
| **Operating System** | Debian 7 |
| **Website** | [Selfoss Website](https://selfoss.aditu.de/) |
| **Last Update** | 14/06/2013 |
| **Others** | Nginx |
{{< /table >}}

## Introduction

[Selfoss](https://selfoss.aditu.de/)[^1] is a self-hosted RSS reader solution. It could be compared to Tiny Tinys RSS or Feedly.

## Installation

First, you need to download the zip file and uncompress it to the good directory:

```bash
mkdir /usr/share/nginx/www/selfoss
cd /usr/share/nginx/www/selfoss
wget http://selfoss.aditu.de/selfoss-2.7.zip
unzip selfoss-2.7.zip
chown -Rf www-data. .
rm -f selfoss-2.7.zip
```

## Configuration

### Database server

We're going to setup the MariaDB part:

```sql
CREATE DATABASE selfoss;
CREATE USER 'selfoss_user'@'localhost' IDENTIFIED BY 'selfoss_password';
GRANT ALL ON selfoss.* TO 'selfoss_user'@'localhost' IDENTIFIED BY 'selfoss_password';
FLUSH privileges;
```

Replace login and password with your needs.

### Web server

We need to configure the web server:

```bash {linenos=table}
server {
    include listen_port.conf;

    server_name feed.deimos.fr;
    root /usr/share/nginx/www/selfoss;
    index index.php;

    access_log /var/log/nginx/feed.deimos.fr_access.log;
    error_log /var/log/nginx/feed.deimos.fr_error.log;

    location ~* \ (gif|jpg|png) {
        expires 30d;
    }

    location ~ ^/favicons/.*$ {
        try_files $uri /data/$uri;
    }

    location ~* ^/(data\/logs|data\/sqlite|config\.ini|\.ht) {
        deny all;
    }

    location / {
        try_files $uri /public/$uri /index.php$is_args$args;
    }

    location ~ \.php$ {
        client_body_timeout 360;
        send_timeout 360;
        include fastcgi_params;
        fastcgi_pass unix:/var/run/php5-fpm.sock;
        fastcgi_intercept_errors on;
    }

    # Drop config
    include drop.conf;
}
```

Then create the symlink:

```bash
ln -s /etc/nginx/sites-available/feed.deimos.fr /etc/nginx/sites-enabled/
```

Then add it to the www-data crontab:

```bash
curl -L -s -k "https://rss2.deimos.fr/update"
```

Copy and past the default configuration:

```bash
cp defaults.ini config.ini
```

Then modify the fields corresponding to your new database setup and remove unneeded lines, like in this example:

```bash
; see http://selfoss.aditu.de for more information about
; the configuration parameters
[globals]
db_type=mysql
db_host=localhost
db_database=selfoss
db_username=selfoss_user
db_password=selfoss_password
db_port=3306
logger_level=ERROR
items_perpage=50
items_lifetime=30
base_url=http://field.deimos.fr
username=deimos
password=xxxxxxxxxxxxxxxx
;salt=lkjl1289
salt=Tarn0twif
rss_title=selfoss feed
rss_max_items=10000
```

You can generate a password in adding "/password" at the end of the url (ex. http://feed.deimos.fr/password)

Then connect to the interface to finish the setup :-)

## References

[^1]: http://selfoss.aditu.de/
