---
weight: 999
url: "/Firefox_Sync_Server_create_your_own_Sync_Server/"
title: "Firefox Sync Server: Create Your Own Sync Server"
description: "How to install and configure your own Firefox Sync Server for synchronizing Firefox data across multiple devices"
categories: ["Servers", "Firefox", "Web"]
date: "2013-09-05T14:34:00+02:00"
lastmod: "2013-09-05T14:34:00+02:00"
tags: ["Firefox", "Sync", "MariaDB", "Debian", "Nginx"]
toc: true
---

![Firefox Sync Server](/images/firefox_sync_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Operating System** | Debian 7 |
| **Website** | [Firefox Sync Server Website](https://docs.services.mozilla.com/howtos/run-sync.html) |
| **Last Update** | 05/09/2013 |
| **Others** | MariaDB 5.5 |
{{< /table >}}

## Introduction

Firefox Sync, originally branded Mozilla Weave, is a browser synchronization feature that allows users to partially synchronize bookmarks, browsing history, preferences, passwords, filled forms, add-ons and the last 25 opened tabs across multiple computers.

It keeps user data on Mozilla servers, but the data is encrypted in such a way that no third party, not even Mozilla, can access user information.

Firefox Sync was originally an add-on for Mozilla Firefox 3.x and SeaMonkey 2.0, but it has been a built-in feature since Firefox 4.0 and SeaMonkey 2.1.[^1]

## Prerequisite

First of all, we need to install those dependencies:

```bash
aptitude install python-dev mercurial python-virtualenv make gcc
```

Then we will install MariaDB.

To install MariaDB, it's unfortunately not embedded in Debian, so we'll add a repository. First of all, install a python tool to get aptkey:

```bash
aptitude install python-software-properties
```

Then let's add this repository (https://downloads.mariadb.org/mariadb/repositories/):

```bash
apt-key adv --recv-keys --keyserver keyserver.ubuntu.com 0xcbcb082a1bb943db
add-apt-repository 'deb http://mirrors.linsrv.net/mariadb/repo/10.0/debian wheezy main'
```

We're now going to change apt pinning to prioritize MariaDB's repository (`/etc/apt/preferences.d/mariadb`):

```bash
Package: *
Pin: release o=MariaDB
Pin-Priority: 1000
```

Then install MariaDB:

```bash
aptitude update
aptitude install mariadb-server
```

Now we will need that package also to be able to build the server.

```bash
aptitude install libmariadbclient-dev
```

## Installation

We can now get the sources:

```bash
cd /usr/share
hg clone https://hg.mozilla.org/services/server-full firefox_sync
cd firefox_sync
```

And launch the build:

```bash
> make build
[...]
Building the app
  Checking the environ   [ok]
  Updating the repo   [ok]
  Building Services dependencies
    Getting server-core     [ok]
    Getting server-reg     [ok]
    Getting server-storage     [ok]  [ok]
  Building External dependencies   [ok]
  Now building the app itself   [ok]
[done]
```

Let's install the latest Python module and Guinicorn (WSGI HTTP Server):

```bash
./bin/pip install Mysql-Python
./bin/pip install gunicorn
```

We are going to create a dedicated user for this application and reset rights:

```bash
groupadd firefoxsync
useradd -d /usr/share/firefox_sync -g firefoxsync -r -s /bin/bash firefoxsync
chown -Rf firefoxsync. /usr/share/firefox_sync
```

## Configuration

Alright, all the installation is now finished. Let's configure everything.

### MariaDB

You need to create a database and user (fit with your informations):

```sql
CREATE DATABASE firefox_syncdb;
CREATE USER 'firefox_sync'@'localhost' IDENTIFIED BY 'password';
GRANT ALL ON firefox_syncdb.* TO 'firefox_sync'@'localhost' IDENTIFIED BY 'password';
FLUSH privileges;
```

### Sync Server

Edit the configuration to set database informations (`etc/sync.conf`):

```ini {linenos=table,hl_lines=[9,22,27,31],anchorlinenos=true}
[captcha]
use = true
public_key = xxxxxxxxxxxxxxxxxxxxxxxx
private_key = xxxxxxxxxxxxxxxxxxxxxxxx
use_ssl = false

[storage]
backend = syncstorage.storage.sql.SQLStorage
sqluri = mysql://firefox_sync:password@localhost:3306/firefox_syncdb
standard_collections = false
# Set quota and size
use_quota = true
quota_size = 5120
pool_size = 100
pool_recycle = 3600
reset_on_return = true
display_config = true
create_tables = true

[auth]
backend = services.user.sql.SQLUser
sqluri = mysql://firefox_sync:password@localhost:3306/firefox_syncdb
pool_size = 100
pool_recycle = 3600
create_tables = true
# Uncomment the next line to disable creation of new user accounts.
#allow_new_users = false

[nodes]
# You must set this to your client-visible server URL.
fallback_node = http://firefoxsync.deimos.fr:5000/

[smtp]
host = localhost
port = 25
sender = firefoxsync@mycompany.com

[cef]
use = true
file = syslog
vendor = mozilla
version = 0
device_version = 1.3
product = weave
```

{{< alert context="warning" text="It's preferable to use SSL connection. If you have autosigned certificates, open manually the URL with firefox to accept them and avoiding errors" />}}

If you're not going to use Nginx, check that your firewall port is open on 5000 port number:

```bash
iptables -t filter -A INPUT -p tcp --dport 5000 -j ACCEPT
```

For more security and if you're going to use a web server like Nginx, it's better to listen only on localhost. In addition, you need to change the 'use' parameter from http to gunicorn. And to finish, you also need to change the log path (`development.ini`):

```ini
...
[server:main]
use = egg:gunicorn
host = 0.0.0.0
port = 5000
workers = 2
timeout = 60
...
[handler_syncserver_errors]
class = handlers.RotatingFileHandler
args = ('/var/log/firefoxsync/sync-error.log',)
level = ERROR
formatter = generic
```

Then let's create those folders:

```bash
mkdir -p /var/log/firefoxsync /var/run/firefoxsync
chown firefoxsync. /var/log/firefoxsync /var/run/firefoxsync
```

You can now try to manually launch the server if you want and sync a user:

```bash
su - firefoxsync -c '/usr/share/firefox_sync/bin/gunicorn_paster /usr/share/firefox_sync/development.ini &'
```

Then kill it once you've tested it as we're going to add an init script for it.

### Nginx

Adapt the configuration to your needs:

```bash
server {
    include listen_port.conf;
    listen 443 ssl;

    ssl_certificate /etc/nginx/ssl/server.crt;
    ssl_certificate_key /etc/nginx/ssl/server.key;
    ssl_session_timeout 5m;

    # Force SSL
    if ($scheme = http) {
        return 301 https://$host$request_uri;
    }

    server_name firefoxsync.deimos.fr;

    access_log /var/log/nginx/firefoxsync.deimos.fr_access.log;
    error_log /var/log/nginx/firefoxsync.deimos.fr_error.log;

    location / {
        proxy_pass_header Server;
        proxy_set_header Host $http_host;
        proxy_redirect off;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Scheme $scheme;
        proxy_connect_timeout 10;
        proxy_read_timeout 10;
        proxy_pass http://localhost:5000/;
    }

    # Drop config
    include drop.conf;
}
```

Then enable it:

```bash
ln -s /etc/nginx/sites-available/firefoxsync.deimos.fr /etc/nginx/sites-enabled/
service nginx reload
```

### Debian

As there is no init script to launch it automatically on boot, we're going to change that (`/etc/init.d/firefoxsync`):

```bash
#!/bin/bash

### BEGIN INIT INFO
# Provides:          paster
# Required-Start:    $all
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: starts the paster server
# Description:       starts paster
### END INIT INFO
DESC="Mozilla Sync Server"
PROJECT=/usr/share/firefox_sync
VIRTUALENV=$PROJECT
PID_DIR=/var/run/firefoxsync
PID_FILE=$PID_DIR/firefoxsync.pid
LOG_FILE=/var/log/firefoxsync/firefoxsync.log
USER=firefoxsync
GROUP=firefoxsync
PROD_FILE=$PROJECT/development.ini
RET_VAL=0

# Load the VERBOSE setting and other rcS variables
. /lib/init/vars.sh

# Define LSB log_* functions.
# Depend on lsb-base (>= 3.2-14) to ensure that this file is present
# and status_of_proc is working.
. /lib/lsb/init-functions

# Activate Python virtual environment
source $VIRTUALENV/bin/activate

# Change directory to project
cd $PROJECT

case "$1" in
    start)
        log_daemon_msg "Starting $DESC"
        mkdir -p /var/run/firefoxsync
        chown firefoxsync. /var/run/firefoxsync
        paster serve \
        --daemon \
        --pid-file=$PID_FILE \
        --log-file=$LOG_FILE \
        --user=$USER \
        --group=$GROUP \
        $PROD_FILE \
        start
    ;;
    stop)
        log_daemon_msg "Stopping $DESC"
        paster serve \
        --daemon \
        --pid-file=$PID_FILE \
        --log-file=$LOG_FILE \
        --user=$USER \
        --group=$GROUP \
        $PROD_FILE \
        stop
    ;;
    restart|force-reload)
        log_daemon_msg "Restarting $DESC"
        paster serve \
        --daemon \
        --pid-file=$PID_FILE \
        --log-file=$LOG_FILE \
        --user=$USER \
        --group=$GROUP \
        $PROD_FILE \
        restart
    ;;
    status)
        paster serve \
        --daemon \
        --pid-file=$PID_FILE \
        --log-file=$LOG_FILE \
        --user=$USER \
        --group=$GROUP \
        status
    ;;
    *)
        echo $"Usage: $0 {start|stop|status|restart|force-reload}"
        exit 3
esac
exit $RET_VAL
```

Then update it on runlevels and start it:

```bash
cd /etc/init.d
chmod 755 firefoxsync
update-rc.d firefoxsync defaults
/etc/init.d/firefoxsync start
```

### Logrotate

You will see logs are verbose enough to install a logrotate script (`/etc/logrotate.d/firefoxsync`):

```bash
/var/log/firefoxsync/*.log {
    weekly
    missingok
    rotate 5
    compress
    delaycompress
    notifempty
    create 644 firefoxsync firefoxsync
    sharedscripts
    postrotate
        chown firefoxsync. /var/log/firefoxsync/firefox-sync.log
        chmod 644 /var/log/firefoxsync/firefox-sync.log
    endscript
}
```

## Upgrade

To upgrade the sever, simply run those commands as root:

```bash
/etc/init.d/firefox_sync stop
cd /usr/share/firefox_sync
su - firefoxsync
hg pull
hg update
make build
exit
/etc/init.d/firefox_sync start
```

## Client

On the client side, there is one account to create and specify the url of the server. Then you could associate all your device to this account.

## References

[^1]: http://en.wikipedia.org/wiki/Firefox_Sync
