---
weight: 999
url: "/Seafile\\:\_make\_your\_personal\_storage\_cloud\_easily/"
title: "Seafile: Make Your Personal Storage Cloud Easily"
description: "Learn how to install and configure Seafile cloud storage system on Debian 7 with Nginx and MariaDB/MySQL."
categories: ["Nginx", "Debian", "Storage"]
date: "2014-06-22T05:10:00+02:00"
lastmod: "2014-06-22T05:10:00+02:00"
tags: ["Cloud Storage", "Seafile", "MariaDB", "MySQL", "Nginx", "Webdav"]
toc: true
---

![Seafile](/images/seafile_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | Seafile 3.0 |
| **Operating System** | Debian 7 |
| **Website** | [Seafile Website](https://seafile.com) |
| **Last Update** | 22/06/2014 |
{{< /table >}}

## Introduction

Seafile[^1] is a next-generation open source cloud storage system, with advanced support for file syncing, privacy protection and teamwork.

Collections of files are called libraries, and each library can be synced separately. A library can be encrypted with a user chosen password. This password is not stored on the server, so even the server admin can't view your file contents.

Seafile lets you create groups with file syncing, wiki, and discussion enabling easy collaboration around documents within a team.

For this documentation, you need to have:

1. A web server (here [Nginx](/Nginx_:_Installation_et_configuration_d'une_alternative_d'Apache/))
2. A [MariaDB]({{< ref "docs/Servers/Databases/MySQL-MariaDB/mariadb_migration_from_mysql.md">}}) / [MySQL server]({{< ref "docs/Servers/Databases/MySQL-MariaDB/mysql_installation_and_configuration.md">}})
3. Python installed

## Installation

The installation of Seafile is not complicated but a little bit long. First of all, you need o grab the latest seafile version:

```bash
mkdir -p /usr/share/nginx/www/seafile
cd /usr/share/nginx/www/seafile
wget https://bitbucket.org/haiwen/seafile/downloads/seafile-server_3.0.0_x86-64.tar.gz
tar -xzf seafile-server_3.0.0_x86-64.tar.gz
chown -Rf www-data. seafile-server-3.0.0
```

Then we need to install dependencies:

```bash
aptitude install python-simplejson python-setuptools python-mysqldb python-imaging
```

## Configuration

### MariaDB

The first thing to do is to configure MariaDB accounts and create databases (adapt those lines with your needs):

```mysql
CREATE DATABASE `seafile\_ccnet`;
CREATE DATABASE `seafile\_db`;
CREATE DATABASE `seafile\_hub`;
CREATE USER 'seafile_user'@'127.0.0.1' identified by 'password';
GRANT ALL PRIVILEGES ON `seafile\_ccnet` . * TO 'seafile_user'@'127.0.0.1';
GRANT ALL PRIVILEGES ON `seafile\_db` . * TO 'seafile_user'@'127.0.0.1';
GRANT ALL PRIVILEGES ON `seafile\_hub` . * TO 'seafile_user'@'127.0.0.1';
FLUSH PRIVILEGES;
```

Then you have to launch the installer:

```bash
> ./setup-seafile-mysql.sh
[...]
-------------------------------------------------------
Please choose a way to initialize seafile databases:
-------------------------------------------------------

[1] Create new ccnet/seafile/seahub databases
[2] Use existing ccnet/seafile/seahub databases
```

Choose to use existing databases and credentials when it is asked to you. The 1 choice can be used but as I encountered issues, that's why I did it manually.

### Ccnet

Regarding the ccnet configuration, there is nothing to do especially instead of if you want SSL instead. If you want https, modify the SERVICE_URL like this (`ccnet/ccnet.conf`):

```ini
[General]
SERVICE_URL = https://seafile.deimos.fr
...
```

### Seahub

For SSL as well, update this file by adding on top of the file this kind of line (`seahub_settings.py`):

```ini
HTTP_SERVER_ROOT = 'https://seafile.deimos.fr/seafhttp'
```

### Webdav

To implement webdav for clients that can't use Seafile client, add this configuration (`/usr/share/nginx/www/seafile/conf/seafdav.conf`):

```bash
[WEBDAV]
enabled = true
port = 8080
fastcgi = true
share_name = /seafdav
```

### Nginx

Now to finish, here is an example of Nginx configuration for SSL purpose:

```bash
server {
    include listen_port.conf;
    server_name seafile.deimos.fr;
    # Force redirect http to https
    rewrite ^ https://$http_host$request_uri? permanent;
}

server {
    include ssl/deimos.fr_ssl.conf;

    include pagespeed.conf;
    server_name seafile.deimos.fr;

    root /usr/share/nginx/www/deimos.fr/seafile;

    access_log /var/log/nginx/seafile.deimos.fr_access.log;
    error_log /var/log/nginx/seafile.deimos.fr_error.log;

    # Max upload size
    client_max_body_size 1G;

    location / {
        fastcgi_pass    127.0.0.1:8090;
        fastcgi_param   SCRIPT_FILENAME     $document_root$fastcgi_script_name;
        fastcgi_param   PATH_INFO           $fastcgi_script_name;
        fastcgi_param   SERVER_PROTOCOL     $server_protocol;
        fastcgi_param   QUERY_STRING        $query_string;
        fastcgi_param   REQUEST_METHOD      $request_method;
        fastcgi_param   CONTENT_TYPE        $content_type;
        fastcgi_param   CONTENT_LENGTH      $content_length;
        fastcgi_param   SERVER_ADDR         $server_addr;
        fastcgi_param   SERVER_PORT         $server_port;
        fastcgi_param   SERVER_NAME         $server_name;
        fastcgi_param   REMOTE_ADDR         $remote_addr;
        fastcgi_param   HTTPS               on;
        fastcgi_param   HTTP_SCHEME         https;
    }

    # Webdav
    location /seafdav {
        fastcgi_pass    127.0.0.1:8080;
        fastcgi_param   SCRIPT_FILENAME     $document_root$fastcgi_script_name;
        fastcgi_param   PATH_INFO           $fastcgi_script_name;

        fastcgi_param   SERVER_PROTOCOL     $server_protocol;
        fastcgi_param   QUERY_STRING        $query_string;
        fastcgi_param   REQUEST_METHOD      $request_method;
        fastcgi_param   CONTENT_TYPE        $content_type;
        fastcgi_param   CONTENT_LENGTH      $content_length;
        fastcgi_param   SERVER_ADDR         $server_addr;
        fastcgi_param   SERVER_PORT         $server_port;
        fastcgi_param   SERVER_NAME         $server_name;

        fastcgi_param   HTTPS               on;
    }

    location /seafhttp {
        rewrite ^/seafhttp(.*)$ $1 break;
        proxy_pass http://127.0.0.1:8082;
        client_max_body_size 0;
    }

    location /media {
        root /var/www/deimos.fr/seafile/seafile-server-latest/seahub;
    }
}
```

### Init script

Regarding the init script, add it and adapt the highlighted lines (`/etc/init.d/seafile-server`):

```bash {linenos=table,hl_lines=[14,17,23,25]}
#!/bin/sh

### BEGIN INIT INFO
# Provides:          seafile-server
# Required-Start:    $local_fs $remote_fs $network
# Required-Stop:     $local_fs
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: Starts Seafile Server
# Description:       starts Seafile Server
### END INIT INFO

# Change the value of "user" to your linux user name
user=www-data

# Change the value of "script_path" to your path of seafile installation
seafile_dir=/usr/share/nginx/www/deimos.fr/seafile
script_path=${seafile_dir}/seafile-server-latest
seafile_init_log=${seafile_dir}/logs/seafile.init.log
seahub_init_log=${seafile_dir}/logs/seahub.init.log

# Change the value of fastcgi to true if fastcgi is to be used
fastcgi=true
# Set the port of fastcgi, default is 8000. Change it if you need different.
fastcgi_port=8090

case "$1" in
        start)
                sudo -u ${user} ${script_path}/seafile.sh start >> ${seafile_init_log}
                if [  $fastcgi = true ];
                then
                        sudo -u ${user} ${script_path}/seahub.sh start-fastcgi ${fastcgi_port} >> ${seahub_init_log}
                else
                        sudo -u ${user} ${script_path}/seahub.sh start >> ${seahub_init_log}
                fi
        ;;
        restart)
                sudo -u ${user} ${script_path}/seafile.sh restart >> ${seafile_init_log}
                if [  $fastcgi = true ];
                then
                        sudo -u ${user} ${script_path}/seahub.sh restart-fastcgi ${fastcgi_port} >> ${seahub_init_log}
                else
                        sudo -u ${user} ${script_path}/seahub.sh restart >> ${seahub_init_log}
                fi
        ;;
        stop)
                sudo -u ${user} ${script_path}/seafile.sh $1 >> ${seafile_init_log}
                sudo -u ${user} ${script_path}/seahub.sh $1 >> ${seahub_init_log}
        ;;
        *)
                echo "Usage: /etc/init.d/seafile {start|stop|restart}"
                exit 1
        ;;
esac
```

I've changed the default port as mine was already used.

Then activate it on boot:

```bash
chmod 755 /etc/init.d/seafile-server
update-rc.d seafile-server defaults
```

You can now reload your Nginx configuration and start Seafile:

```bash
service nginx reload
service seafile-server start
```

You should have access to Seafile now :-)

## References

[^1]: [https://seafile.com](https://seafile.com)
