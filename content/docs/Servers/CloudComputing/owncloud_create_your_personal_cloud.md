---
weight: 999
url: "/OwnCloud_\\:_créer_son_cloud_personnel/"
title: "OwnCloud: Create Your Personal Cloud"
description: "How to set up your own personal cloud server using OwnCloud with Nginx and MariaDB"
categories: ["Nginx", "Debian", "Database"]
date: "2013-02-26T08:03:00+02:00"
lastmod: "2013-02-26T08:03:00+02:00"
tags: ["Cloud Storage", "Nginx", "MariaDB", "Self-Hosted", "OwnCloud"]
toc: true
---

![ownCloud](/images/owncloud-logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 4.5.6 |
| **Operating System** | Debian 7 |
| **Website** | [ownCloud Website](https://owncloud.org/) |
| **Last Update** | 26/02/2013 |
| **Others** | MariaDB 5.5 |
{{< /table >}}

## Introduction

[ownCloud](https://owncloud.org/)[^1] is an open source implementation of online storage services and various applications (cloud computing).

This tutorial bases the installation of ownCloud on [Nginx]({{< ref "docs/Servers/Web/Nginx/nginx_installation_and_configuration_as_an_apache_alternative.md">}}) and [MariaDB]({{< ref "docs/Servers/Databases/MySQL-MariaDB/mysql_installation_and_configuration.md" >}}).

## Server Installation

### MariaDB

We need to have [MariaDB]({{< ref "docs/Servers/Databases/MySQL-MariaDB/mysql_installation_and_configuration.md" >}}) installed. Then, we must create a database and an account:

```sql
> mysql -uroot -p
CREATE DATABASE owncloud;
CREATE USER 'owncloud_user'@'localhost' IDENTIFIED BY 'owncloud_password';
GRANT USAGE ON * . * TO 'owncloud_user'@'localhost' IDENTIFIED BY 'owncloud_password';
GRANT ALL ON `owncloud`.* TO 'owncloud_user'@'localhost';
FLUSH privileges;
```

### ownCloud

Let's first install the dependencies:

```bash
aptitude install php5 php5-gd php-xml-parser php5-intl php5-mysql smbclient curl libcurl3 php5-curl
```

We'll download the latest version:

```bash
cd /usr/share/nginx/www/
wget http://mirrors.owncloud.org/releases/owncloud-4.5.6.tar.bz2
tar -xjf owncloud-4.5.6.tar.bz2
rm -f owncloud-4.5.6.tar.bz2
chown -Rf www-data. owncloud
```

## Configuration

### Nginx

#### ownCloud 4.X

Here is the ownCloud 4.X configuration for Nginx:

```bash
server {
    include listen_port.conf;
    listen 443 default ssl;
    ssl on;

    ssl_certificate /etc/nginx/ssl/deimos.fr/server-unified.crt;
    ssl_certificate_key /etc/nginx/ssl/deimos.fr/server.key;

    server_name owncloud.deimos.fr;
    root /usr/share/nginx/www/deimos.fr/owncloud;
    index index.php;
    client_max_body_size 1024M;

    access_log /var/log/nginx/cloud.deimos.fr_access.log;
    error_log /var/log/nginx/cloud.deimos.fr_error.log;

    # Force SSL
    if ($scheme = http) {
        return 301 https://$host$request_uri;
    }

    # deny direct access
    location ~ ^/(data|config|\.ht|db_structure\.xml|README) {
        deny all;
    }

    # default try order
    location / {
        try_files $uri $uri/ @webdav;
    }

    # owncloud WebDAV
    location @webdav {
        fastcgi_cache mycache;
        fastcgi_cache_key $request_method$host$request_uri;
        fastcgi_cache_valid any 1h;
        fastcgi_split_path_info ^(.+\.php)(/.*)$;
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        fastcgi_param HTTPS on;
        include fastcgi_params;
        fastcgi_pass unix:/var/run/php5-fpm.sock;
        fastcgi_intercept_errors on;
    }

    location ~ \.php$ {
        try_files $uri = 404;
        fastcgi_cache mycache;
        fastcgi_cache_key $request_method$host$request_uri;
        fastcgi_cache_valid any 1h;
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        fastcgi_param HTTPS on;
        include fastcgi_params;
        fastcgi_pass unix:/var/run/php5-fpm.sock;
        fastcgi_intercept_errors on;
    }

    # Drop config
    include drop.conf;
}
```

#### ownCloud 5.X

And for version 5:

```bash
server {
    include listen_port.conf;
    listen 443 default ssl;
    ssl on;

    ssl_certificate /etc/nginx/ssl/deimos.fr/server-unified.crt;
    ssl_certificate_key /etc/nginx/ssl/deimos.fr/server.key;

    server_name cloud.deimos.fr;
    root /usr/share/nginx/www/deimos.fr/owncloud;
    index index.php;
    client_max_body_size 1024M;

    access_log /var/log/nginx/cloud.deimos.fr_access.log;
    error_log /var/log/nginx/cloud.deimos.fr_error.log;

    # Force SSL
    if ($scheme = http) {
        return 301 https://$host$request_uri;
    }

    rewrite ^/caldav((/|$).*)$ /remote.php/caldav$1 last;
    rewrite ^/carddav((/|$).*)$ /remote.php/carddav$1 last;
    rewrite ^/webdav((/|$).*)$ /remote.php/webdav$1 last;

    error_page 403 = /core/templates/403.php;
    error_page 404 = /core/templates/404.php;

    location ~ ^/(data|config|\.ht|db_structure\.xml|README) {
            deny all;
    }

    location / {
        rewrite ^/.well-known/host-meta /public.php?service=host-meta last;
        rewrite ^/.well-known/host-meta.json /public.php?service=host-meta-json last;
        rewrite ^/.well-known/carddav /remote.php/carddav/ redirect;
        rewrite ^/.well-known/caldav /remote.php/caldav/ redirect;
        rewrite ^(/core/doc/[^\/]+/)$ $1/index.html;
        try_files $uri $uri/ index.php;
    }

    location ~ ^(?<script_name>.+?\.php)(?<path_info>/.*)?$ {
        try_files $script_name = 404;
        fastcgi_cache_valid any 1h;
        include fastcgi_params;
        fastcgi_pass unix:/var/run/php5-fpm.sock;
    }

    location ~* ^.+.(jpg|jpeg|gif|bmp|ico|png|css|js|swf)$ {
        expires 30d;
        # Optional: Don't log access to assets
        access_log off;
    }

    # Drop config
    include drop.conf;
}
```

Restart Nginx and access your cloud: http://owncloud.myserver.com

## Clients

There are desktop clients available, allowing automatic file synchronization. You can [download these versions here](https://owncloud.org/sync-clients/). If you need a higher level of logging, you can also launch the owncloud client like this:

```bash
owncloud --logwindow
```

It's also possible to adjust the polling interval time:

```bash
[ownCloud]
remotePollInterval=30000
localPollInterval=10000
PollTimerExceedFactor=10
maxLogLines=20000
[...]
```

For more information about these options, see [the associated documentation](https://owncloud.org/support/sync-clients/)[^2].

## FAQ

### csync could not create the lock file

If you have this kind of message on the client:

```
csync could not create the lock file
```

It's simply that a lock file persists. Quit the application, delete the lock file:

```bash
rm ~/.local/share/data/ownCloud/lock
```

And restart ownCloud :-)

## References

[^1]: http://owncloud.org/
[^2]: http://owncloud.org/support/sync-clients/
