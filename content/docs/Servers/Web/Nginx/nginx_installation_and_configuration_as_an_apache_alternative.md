---
weight: 999
url: "/Nginx_\\:_Installation_et_configuration_d'une_alternative_d'Apache/"
title: "Nginx: Installation and Configuration as an Apache Alternative"
description: "A comprehensive guide to installing and configuring Nginx as an alternative to Apache, including various optimizations and application configurations."
categories: ["Nginx", "Debian", "Security"]
date: "2014-08-29T15:33:00+02:00"
lastmod: "2014-08-29T15:33:00+02:00"
tags:
  ["nginx", "web server", "php-fpm", "https", "ssl", "wordpress", "mediawiki"]
toc: true
---

![Nginx Logo](/images/nginx-logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 1.2.1 |
| **Operating System** | Debian 7 |
| **Website** | [Nginx Website](https://nginx.org/) |
| **Last Update** | 29/08/2014 |
{{< /table >}}

## Introduction

[Nginx](https://en.wikipedia.org/wiki/Nginx) [engine x] is a web server (or HTTP) software written by Igor Sysoev, whose development began in 2002 for the needs of a high-traffic Russian site.

I've been looking for an alternative to Apache for some time because it's too resource-hungry. When discovering lighttpd, I realized that other servers existed besides Apache and IIS. It was time to dig deeper into this question.

## Installation

For the installation on Debian, it's always simple:

```bash
aptitude install nginx
```

Then we'll start it:

```bash
/etc/init.d/nginx start
```

## Configuration

### php-fpm

This is the most optimal solution for multi-threaded machines (i.e., several cores).

#### Installation

```bash
aptitude install php5-fpm
```

#### Configuration

For the configuration, we will create a basic configuration (if possible, delete the default configuration). We will change it later, but this gives you an idea of a minimal working configuration:

```bash {linenos=table}
server {
    listen       80;
    server_name  www.deimos.fr;
    root   /var/www/deimos.fr;
    index  index.html index.htm index.php;

    location ~ \.php$ {
        fastcgi_index  index.php;
        fastcgi_pass   unix:/var/run/php5-fpm.sock;
        # fastcgi_pass    127.0.0.1:9000;
        fastcgi_param  SCRIPT_FILENAME  $document_root$fastcgi_script_name;
        include fastcgi_params;
    }
}
```

We will change the listening type from TCP to sockets:

```ini {linenos=table}
; The address on which to accept FastCGI requests.
; Valid syntaxes are:
;   'ip.add.re.ss:port'    - to listen on a TCP socket to a specific address on
;                            a specific port;
;   'port'                 - to listen on a TCP socket to all addresses on a
;                            specific port;
;   '/path/to/unix/socket' - to listen on a unix socket.
; Note: This value is mandatory.
listen = /var/run/php5-fpm.sock
```

Then apply the configuration:

```bash
ln -s /etc/nginx/sites-available/www.deimos.fr /etc/nginx/sites-enabled/www.deimos.fr
```

Then restart php-fpm and nginx!

##### PHP-FPM Status

If you need to monitor PHP-FPM, you'll certainly need information about your PHP-FPM. For this, edit the following file and uncomment these lines:

```bash {linenos=table}
[...]
pm.status_path = /fpm-status
[...]
```

Then edit your Nginx configuration and add:

```bash {linenos=table}
[...]
    # PHP-FPM Status
    location /fpm-status {
        include fastcgi_params;
        fastcgi_pass unix:/var/run/php5-fpm.sock;
        # User access
        auth_basic "Please logon";
        auth_basic_user_file /etc/nginx/access/htaccess;
    }
[...]
```

To prevent anyone from accessing it, I've placed a htaccess, but you can choose the security method you prefer.

Restart PHP-FPM and Nginx, then access your server like this:

- http://server/fpm-status: general information
- http://server/fpm-status?full: detailed general information
- http://server/fpm-status?json&full: json export
- http://server/fpm-status?html&full: html export
- http://server/fpm-status?xml&full: xml export

This will give you something like this in the non-detailed version:

```
pool:                 www
process manager:      dynamic
start time:           18/Dec/2013:19:00:41 +0100
start since:          52972
accepted conn:        5268
listen queue:         0
max listen queue:     0
listen queue len:     0
idle processes:       3
active processes:     1
total processes:      4
max active processes: 4
max children reached: 0
```

### PHP Fast CGI

This method works but is not optimal. Use [php-fpm](#php-fpm) if you can.

#### Installation

To make Nginx support PHP, we'll install:

```bash
aptitude install php5-cgi
```

Then we will create an init service for this fast cgi:

```bash {linenos=table}
#! /bin/sh
### BEGIN INIT INFO
# Provides:          php-fcgi
# Required-Start:    $all
# Required-Stop:     $all
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: Start and stop php-cgi in external FASTCGI mode
# Description:       Start and stop php-cgi in external FASTCGI mode
### END INIT INFO

# Author: Kurt Zankl <kz@xon.uni.cc>
# Modified by Deimos <xxx@mycompany.com>

# Do NOT "set -e"

BIND=127.0.0.1:9000
USER=www-data
PHP_FCGI_CHILDREN=15
PHP_FCGI_MAX_REQUESTS=1000

PATH=/usr/bin:/sbin:/usr/sbin:/bin
PHP_CGI=/usr/bin/php-cgi
PHP_CGI_NAME=`basename $PHP_CGI`
PHP_CGI_ARGS="- USER=$USER PATH=$PATH PHP_FCGI_CHILDREN=$PHP_FCGI_CHILDREN PHP_FCGI_MAX_REQUESTS=$PHP_FCGI_MAX_REQUESTS $PHP_CGI -b $BIND"
RETVAL=0

# Load the VERBOSE setting and other rcS variables
. /lib/init/vars.sh

# Define LSB log_* functions.
# Depend on lsb-base (>= 3.0-6) to ensure that this file is present.
. /lib/lsb/init-functions

do_start()
{
        echo -n "Starting PHP FastCGI : "
        start-stop-daemon --start --quiet --background --chuid "$USER" --exec /usr/bin/env -- $PHP_CGI_ARGS
        RETVAL=$?
        echo "$PHP_CGI_NAME."
}

do_stop()
{
        echo -n "Stopping PHP FastCGI : "
        killall -q -w -u $USER $PHP_CGI
        RETVAL="$?"
        echo "$PHP_CGI_NAME."
}

case "$1" in
   start)
      do_start
   ;;
   stop)
      do_stop
   ;;
   restart)
      do_stop
      do_start
   ;;
   *)
      echo "Usage: php-fcgi {start|stop|restart}"
      exit 3
    ;;
esac
exit $RETVAL
```

Then we will update the default RCs and start everything:

```bash
cd /etc/init.d
chmod 755 /etc/init.d/php-fcgi
update-rc.d php-fcgi defaults
/etc/init.d/php-fcgi start
```

#### Configuration

To configure PHP support, edit the base Nginx file and adapt:

```bash
server {
       listen   80;
       server_name  localhost;

       access_log  /var/log/nginx/localhost.access.log;

       location / {
               root   /var/www/nginx-default;
               index  index.html index.htm index.php;
       }

       location /doc {
               root   /usr/share;
               autoindex on;
               allow 127.0.0.1;
               deny all;
       }

       location /images {
               root   /usr/share;
               autoindex on;
       }

       #error_page  404  /404.html;

       # redirect server error pages to the static page /50x.html
       #
       error_page   500 502 503 504  /50x.html;
       location = /50x.html {
               root   /var/www/nginx-default;
       }

       # proxy the PHP scripts to Apache listening on 127.0.0.1:80
       #
       #location ~ \.php$ {
               #proxy_pass   http://127.0.0.1;
       #}

       # pass the PHP scripts to FastCGI server listening on 127.0.0.1:9000
       #
       location ~ \.php$ {
               root /var/www/nginx-default;
               include /etc/nginx/fastcgi_params;
               fastcgi_pass   127.0.0.1:9000;
               fastcgi_index  index.php;
               fastcgi_param  SCRIPT_FILENAME  /var/www/nginx-default/$fastcgi_script_name;
       }

       # deny access to .htaccess files, if Apache's document root
       # concurs with nginx's one
       #
       #location ~ /\.ht {
               #deny  all;
       #}
}
```

### nginx.conf

Here is the global server configuration. I added comments on important lines:

```bash {linenos=table,hl_lines=["2-4","8-9","25-26","31-37","53-61"]}
user www-data;
# Number of working process
worker_processes 2;
worker_rlimit_nofile    2000;
pid /var/run/nginx.pid;

events {
    # Maximum connection number
    worker_connections 2048;
    use epoll;
    # multi_accept on;
}

http {
	##
	# Basic Settings
	##

	sendfile on;
	tcp_nopush on;
	tcp_nodelay on;
	keepalive_timeout 65;
	types_hash_max_size 2048;
    map_hash_bucket_size 64;
    # Security to hide version
	server_tokens off;

	server_names_hash_bucket_size 64;
	# server_name_in_redirect off;

    # Grow this 2 values if you get 502 error message
    fastcgi_buffers 256 16k;
    fastcgi_buffer_size 32k;
    # Nginx cache to boost performances
    fastcgi_cache_path /usr/share/nginx/cache levels=1:2
        keys_zone=mycache:10m
        inactive=1h max_size=256m;

	include /etc/nginx/mime.types;
	default_type application/octet-stream;

	##
	# Logging Settings
	##

	access_log /var/log/nginx/access.log;
	error_log /var/log/nginx/error.log;

	##
	# Gzip Settings
	##

	gzip on;
	gzip_disable "msie6";
        gzip_static on;
	gzip_vary on;
	gzip_proxied any;
	gzip_comp_level 6;
	gzip_buffers 16 8k;
	gzip_http_version 1.1;
	gzip_types text/plain text/css application/json application/x-javascript text/xml application/xml application/xml+rss text/javascript;

	##
	# nginx-naxsi config
	##
	# Uncomment it if you installed nginx-naxsi
	##

	#include /etc/nginx/naxsi_core.rules;

	##
	# nginx-passenger config
	##
	# Uncomment it if you installed nginx-passenger
	##

	#passenger_root /usr;
	#passenger_ruby /usr/bin/ruby;

	##
	# Virtual Host Configs
	##

	include /etc/nginx/conf.d/*.conf;
	include /etc/nginx/sites-enabled/*;
}
```

### drop.conf

This file will allow us to not log non-essential information and prevent access to potentially sensitive files:

```bash {linenos=table}
# Do not log robots.txt if not found
location = /robots.txt  { access_log off; log_not_found off; }
# Do not log favicon.ico if not found
location = /favicon.ico { access_log off; log_not_found off; }
# Do not give access to hidden files
location ~ /\.          { access_log off; log_not_found off; deny all; }
# Do not give access to vim backuped files
location ~ ~$           { access_log off; log_not_found off; deny all; }
```

Then in your configuration files, just add these 2 lines:

```bash {linenos=table,hl_lines=[3,4]}
server {
[...]
    # Drop config
    include drop.conf;
}
```

### VirtualHosts

Just like with Apache, you can create virtual hosts. For example, let's create the deimos.fr blog:

```bash
server {
       listen   80;
       server_name  blog.deimos.fr;

       access_log  /var/log/nginx/localhost.access.log;

       location / {
               root   /var/www/nginx-default/blog;
               index  index.html index.htm index.php;
       }

       location /doc {
               root   /usr/share;
               autoindex on;
               allow 127.0.0.1;
               deny all;
       }

       location /images {
               root   /usr/share;
               autoindex on;
       }
}
```

Here, blog.deimos.fr will be accessible on port 80.

Let's activate this config:

```bash
ln -s /etc/nginx/sites-available/blog /etc/nginx/sites-enabled/blog
```

All that remains is to reload the server and it works :-)

### htaccess

To restrict access to some of your sites, here's one of the oldest but effective solutions - htaccess! To enable them, we will need the htpasswd binary:

```bash
aptitude install apache2-utils
```

We will now generate a file containing the logins and passwords with a first user:

```bash
htpasswd -c /etc/nginx/htaccess deimos
```

Enter the password and you're set. Let's declare it in the config, in my previously created virtualhost:

```bash
server {
       listen   80;
       server_name  blog.deimos.fr;

       access_log  /var/log/nginx/localhost.access.log;

       location / {
               root   /var/www/nginx-default/blog;
               index  index.html index.htm index.php;
               auth_basic "Restricted";
               auth_basic_user_file /etc/nginx/htaccess;
       }

       location /doc {
               root   /usr/share;
               autoindex on;
               allow 127.0.0.1;
               deny all;
       }

       location /images {
               root   /usr/share;
               autoindex on;
       }
}
```

All that remains is to reload the server and it works :-)

### Default Port

To specify the default port, you can use an include in your configuration files that will call a file containing the port. This way it will be very easy to modify the default port in one go:

```bash {linenos=table}
listen 80;
```

And in the configuration files:

```bash {linenos=table}
server {
    include listen_port.conf;
[...]
```

### SSL

#### Installation

For SSL installation, compared to Apache which doesn't natively and simply allow SSL on VirtualHosts, it's much easier with Nginx:

```bash
aptitude install openssl
```

#### Configuration

First, we'll generate SSL keys:

```bash
mkdir -p /etc/nginx/ssl
cd /etc/nginx/ssl
openssl req -new -x509 -nodes -out server.crt -keyout server.key
```

I put it in an ssl folder with the nginx config, but since you can have multiple certificates, I encourage you to create a hierarchy:

```bash {linenos=table,hl_lines=[2,"5-14"]}
 server {
        listen   443 ssl;
        server_name  localhost;

        ssl_certificate /etc/nginx/ssl/server.crt;
        ssl_certificate_key /etc/nginx/ssl/server.key;
        # Resumption
        ssl_session_cache shared:SSL:10m;
        # Timeout
        ssl_session_timeout 10m;

        # Security options
        ssl_ciphers ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-AES256-SHA384:ECDHE-RSA-AES256-SHA384:ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES128-SHA256:ECDHE-RSA-AES128-SHA256:ECDHE-ECDSA-RC4-SHA:ECDHE-RSA-RC4-SHA:ECDH-ECDSA-RC4-SHA:ECDH-RSA-RC4-SHA:ECDHE-RSA-AES256-SHA:RC4-SHA;
        ssl_prefer_server_ciphers on;

        # HSTS (force users to come in SSL if they've already been once)
        add_header Strict-Transport-Security "max-age=31536000; includeSubdomains";

        access_log  /var/log/nginx/localhost.access.log;

        location / {
                root   /var/www/nginx-default/webmail;
                index  index.html index.htm index.php;
        }

        # pass the PHP scripts to FastCGI server listening on 127.0.0.1:9000
        #
        location ~ \.php$ {
                root /var/www/nginx-default;
                include /etc/nginx/fastcgi_params;
                fastcgi_pass   unix:/var/run/php5-fpm.sock;
                fastcgi_index  index.php;
                fastcgi_param  SCRIPT_FILENAME  /var/www/nginx-default/$fastcgi_script_name;
        }
 }
```

Let's enable this config:

```bash
ln -s /etc/nginx/sites-available/webmail /etc/nginx/sites-enabled/webmail
```

And restart the Nginx server for it to take effect.

You can test the security of your SSL server: https://www.ssllabs.com/ssltest/index.html

##### Force SSL Connection

It is possible to configure your virtualhost in the normal way but to absolutely want a redirect to SSL. For this, it's very simple, just add:

```bash {linenos=table}
server {
    listen      80;
    server_name www.deimos.fr;
    rewrite     ^   https://$server_name$request_uri? permanent;
}
```

### FastCGI Cache

FastCGI cache is useful to increase the speed of your website. We'll allocate a few MB for the cache and put it in a tmpfs to speed things up even more. In the main configuration, we'll create this cache:

```bash {linenos=table}
[...]
http {
    [...]
    fastcgi_cache_path /usr/share/nginx/cache levels=1:2
        keys_zone=mycache:10m
        inactive=1h max_size=256m;
    [...]
}
[...]
```

Here I'm creating a cache called "mycache" with a size of 256MB.

In your VirtualHost configurations, add these lines:

```bash {linenos=table,hl_lines=[3,4,5]}
[...]
        location ~ \.php$ {
                fastcgi_cache mycache;
                fastcgi_cache_key $request_method$host$request_uri;
                fastcgi_cache_valid any 1h;
                include fastcgi_params;
                fastcgi_pass unix:/var/run/php5-fpm.sock;
                fastcgi_intercept_errors on;
        }
[...]
```

Next we'll create a tmpfs cache folder ([read this doc for more info](https://wiki.deimos.fr/Tmpfs_:_un_filesystem_en_ram_ou_comment_%C3%A9crire_en_ram)):

```bash {linenos=table}
[...]
tmpfs /usr/share/nginx/cache    tmpfs   defaults,size=256m   0   0
[...]
```

Then we'll create the cache directory, mount it and restart nginx:

```bash
mkdir /usr/share/nginx/cache
mount /usr/share/nginx/cache
/etc/init.d/nginx restart
```

### Maintenance Page

Here's an elegant way to create a maintenance page:

```bash {linenos=table}
server {
[...]
   location / {
      try_files /maintenance.html $uri $uri/ @maintenance;
   }

   location @maintenance {
      return 503;
   }
}
```

All you need to do is insert a maintenance.html page at the root of your site and this page will be displayed while you make your changes. You can also use this method:

```bash {linenos=table}
server {
[...]
    location / {
        if (-f $document_root/maintenance.html) {
             return 503;
        }
        [...]
    }

    error_page 503 @maintenance;
    location @maintenance {
        rewrite ^(.*)$ /maintenance.html break;
    }
}
```

### Limit Flooding and Sniffers

It can happen that a malicious person tries to bring your server to its knees by making many requests that typically saturate PHP processes and make your CPUs go to 100%. To avoid this, it is possible to limit the number of requests per IP by using the limit_req module[^1]. To enable this module, insert this line in your Nginx core (global) configuration:

```bash {linenos=table}
http {
[...]
    # Flood/DoS protection
    limit_req_zone $binary_remote_addr zone=limit:10m rate=5r/s;
    limit_req_log_level notice;
[...]
}
```

Here are the options used:

- zone=limit: IPs will be stored in a zone called limit of 10M. Provide enough space in this zone to avoid 503 errors if there is no more space. A 1M zone can contain 16,000 binary_remote_addr type entries.
- rate: we allow 5 requests per second.
- log_level: you can remove this line if you don't want to know in the logs if the limits are reached. Here I want to track this kind of event.

Next, we need to inform on which virtualhost or location we want to apply this security measure:

```bash {linenos=table,hl_lines=[2]}
    location ~ \.php$ {
        limit_req zone=limit burst=5 nodelay;
        include fastcgi_params;
        fastcgi_pass unix:/var/run/php5-fpm.sock;
        fastcgi_intercept_errors on;
    }
```

Again we find options:

- zone: I use the limit zone previously declared in the http (core) configuration
- burst: I allow the burst to 5 maximum connections
- nodelay: if you want requests to still be served even if they are slowed down, use the nodelay parameter.

All that remains is to restart nginx.

## Application Configurations

By default, I use www.deimos.fr to redirect to other VirtualHosts:

```bash {linenos=table}
server {
    include listen_port.conf;
    listen 443 ssl;

    server_name www.deimos.fr;

    access_log /var/log/nginx/www.deimos.fr_access.log;
    error_log /var/log/nginx/www.deimos.fr_error.log;

    # Blog redirect
    rewrite ^/$ $scheme://blog.deimos.fr permanent;
    rewrite ^/blog(/.*)?$ $scheme://blog.deimos.fr$1 permanent;

    # Blocnotesinfo redirect
    rewrite ^/blocnotesinfo(/.*)?$ $scheme://wiki.deimos.fr$1 permanent;

    # Piwik
    rewrite ^/piwik(/.*)?$ $scheme://piwik.deimos.fr$1 permanent;

    # Gitweb
    rewrite ^/gitweb(/.*)?$ $scheme://git.deimos.fr$1 permanent;

    # Drop config
    include drop.conf;
}
```

### Server Status

By default, it's possible to get information about the server status, like this:

```bash {linenos=table,hl_lines=["10-49"]}
server {
    include listen_port.conf;
    listen 443 ssl;

    server_name www.deimos.fr;

    access_log /var/log/nginx/www.deimos.fr_access.log;
    error_log /var/log/nginx/www.deimos.fr_error.log;

    # Nginx status
    location /server-status {
        # Turn on nginx stats
        stub_status on;
        # I do not need logs for stats
        access_log off;
        # Allow from localhost
        allow 127.0.0.1;
        # Deny others
        deny all;
    }

    # PHP-FPM status
    location /php-fpm_status {
        include fastcgi_params;
        fastcgi_pass unix:/var/run/php5-fpm.sock;
        # I do not need logs for stats
        access_log off;
        # Allow from localhost
        allow 127.0.0.1;
        # Deny others
        deny all;
    }

    # Cache APC
    location /apc-cache {
        root /usr/share/doc/php-apc;
    }

    location ~ \.php$ {
        include fastcgi_params;
        fastcgi_pass unix:/var/run/php5-fpm.sock;
        fastcgi_intercept_errors on;
        # I do not need logs for stats
        access_log off;
        # Allow from localhost
        allow 127.0.0.1;
        # Deny others
        deny all;
    }

    # Blog redirect
    rewrite ^/$ $scheme://blog.deimos.fr permanent;
    rewrite ^/blog(/.*)?$ $scheme://blog.deimos.fr$1 permanent;

    # Blocnotesinfo redirect
    rewrite ^/blocnotesinfo(/.*)?$ $scheme://wiki.deimos.fr$1 permanent;

    # Piwik
    rewrite ^/piwik(/.*)?$ $scheme://piwik.deimos.fr$1 permanent;

    # Gitweb
    rewrite ^/gitweb(/.*)?$ $scheme://git.deimos.fr$1 permanent;

    # Drop config
    include drop.conf;
}
```

Configure the authorized addresses properly or place [authentication via htpasswd](#htaccess).

Here's the result when accessing the page (http://www.deimos.fr/server-status):

```
Active connections: 11
server accepts handled requests
 8109 8109 58657
Reading: 0 Writing: 1 Waiting: 10
```

### Wordpress

For the configuration of Wordpress under Nginx, here's an example:

```bash {linenos=table}
server {
    include listen_port.conf;
    listen 443 ssl;

    ssl_certificate /etc/nginx/ssl/deimos.fr/server-unified.crt;
    ssl_certificate_key /etc/nginx/ssl/deimos.fr/server.key;
    ssl_session_timeout 5m;

    server_name blog.deimos.fr;
    root /usr/share/nginx/www/deimos.fr/blog;
    index index.php;

    access_log /var/log/nginx/blog.deimos.fr_access.log;
    error_log /var/log/nginx/blog.deimos.fr_error.log;

    location / {
        try_files $uri $uri/ /index.php?$args;
    }

    location ~ \.php$ {
        fastcgi_cache mycache;
        fastcgi_cache_key $request_method$host$request_uri;
        fastcgi_cache_valid any 1h;
        include fastcgi_params;
        fastcgi_pass unix:/var/run/php5-fpm.sock;
        fastcgi_intercept_errors on;
    }

    # Drop config
    include drop.conf;

    # BEGIN W3TC Browser Cache
    gzip on;
    gzip_types text/css application/x-javascript text/x-component text/richtext image/svg+xml text/plain text/xsd text/xsl text/xml image/x-icon;
    location ~ \.(css|js|htc)$ {
        expires 31536000s;
        add_header Pragma "public";
        add_header Cache-Control "public, must-revalidate, proxy-revalidate";
        add_header X-Powered-By "W3 Total Cache/0.9.2.4";
    }

    location ~ \.(html|htm|rtf|rtx|svg|svgz|txt|xsd|xsl|xml)$ {
        expires 3600s;
        add_header Pragma "public";
        add_header Cache-Control "public, must-revalidate, proxy-revalidate";
        add_header X-Powered-By "W3 Total Cache/0.9.2.4";
    }

    location ~ \.(asf|asx|wax|wmv|wmx|avi|bmp|class|divx|doc|docx|eot|exe|gif|gz|gzip|ico|jpg|jpeg|jpe|mdb|mid|midi|mov|qt|mp3|m4a|mp4|m4v|mpeg|mpg|mpe|mpp|otf|odb|odc|odf|odg|odp|ods|odt|ogg|pdf|png|pot|pps|ppt|pptx|ra|ram|svg|svgz|swf|tar|tif|tiff|ttf|ttc|wav|wma|wri|xla|xls|xlsx|xlt|xlw|zip)$ {
        expires 31536000s;
        add_header Pragma "public";
        add_header Cache-Control "public, must-revalidate, proxy-revalidate";
        add_header X-Powered-By "W3 Total Cache/0.9.2.4";
    }
    # END W3TC Browser Cache

    # BEGIN W3TC Minify core
    rewrite ^/wp-content/w3tc/min/w3tc_rewrite_test$ /wp-content/w3tc/min/index.php?w3tc_rewrite_test=1 last;
    rewrite ^/wp-content/w3tc/min/(.+\.(css|js))$ /wp-content/w3tc/min/index.php?file=$1 last;
    # END W3TC Minify core
}
}
```

### Mediawiki

For setting up Mediawiki with Nginx and short URLs, here's the configuration to adopt. I've also added SSL and forced redirects from the login page to SSL:

```bash {linenos=table}
server {
    include listen_port.conf;
    listen 443 ssl;

    ssl_certificate /etc/nginx/ssl/deimos.fr/server-unified.crt;
    ssl_certificate_key /etc/nginx/ssl/deimos.fr/server.key;
    ssl_session_timeout 5m;

    server_name wiki.deimos.fr wiki.m.deimos.fr;
    root /usr/share/nginx/www/deimos.fr/blocnotesinfo;

    client_max_body_size 5m;
    client_body_timeout 60;

    access_log /var/log/nginx/wiki.deimos.fr_access.log;
    error_log /var/log/nginx/wiki.deimos.fr_error.log;

    location / {
        rewrite ^/$ $scheme://$host/index.php permanent;
        # Short URL redirect
        try_files $uri $uri/ @rewrite;
    }

    location @rewrite {
        if (!-f $request_filename){
            rewrite ^/(.*)$ /index.php?title=$1&$args;
        }
    }

    # Force SSL Login
    set $ssl_requested 0;
    if ($arg_title ~ Sp%C3%A9cial:Connexion) {
        set $ssl_requested 1;
    }
    if ($scheme = https) {
        set $ssl_requested 0;
    }
    if ($ssl_requested = 1) {
        return 301 https://$host$request_uri;
    }

    # Drop config
    include drop.conf;

    # Deny direct access to specific folders
    location ^~ /(maintenance|images)/ {
        return 403;
    }

    location ~ \.php$ {
        fastcgi_cache mycache;
        fastcgi_cache_key $request_method$host$request_uri;
        fastcgi_cache_valid any 1h;
        include fastcgi_params;
        fastcgi_pass unix:/var/run/php5-fpm.sock;
    }

    location = /_.gif {
        expires max;
        empty_gif;
    }

    location ^~ /cache/ {
        deny all;
    }

    location /dumps {
        root /usr/share/nginx/www/deimos.fr/blocnotesinfo/local;
        autoindex on;
    }

    # BEGIN W3TC Browser Cache
    gzip on;
    gzip_types text/css application/x-javascript text/x-component text/richtext image/svg+xml text/plain text/xsd text/xsl text/xml image/x-icon;
    location ~ \.(css|js|htc)$ {
        expires 31536000s;
        add_header Pragma "public";
        add_header Cache-Control "public, must-revalidate, proxy-revalidate";
        add_header X-Powered-By "W3 Total Cache/0.9.2.4";
    }

    location ~ \.(html|htm|rtf|rtx|svg|svgz|txt|xsd|xsl|xml)$ {
        expires 3600s;
        add_header Pragma "public";
        add_header Cache-Control "public, must-revalidate, proxy-revalidate";
        add_header X-Powered-By "W3 Total Cache/0.9.2.4";
    }

    location ~ \.(asf|asx|wax|wmv|wmx|avi|bmp|class|divx|doc|docx|eot|exe|gif|gz|gzip|ico|jpg|jpeg|jpe|mdb|mid|midi|mov|qt|mp3|m4a|mp4|m4v|mpeg|mpg|mpe|mpp|otf|odb|odc|odf|odg|odp|ods|odt|ogg|pdf|png|pot|pps|ppt|pptx|ra|ram|svg|svgz|swf|tar|tif|tiff|ttf|ttc|wav|wma|wri|xla|xls|xlsx|xlt|xlw|zip)$ {
        expires 31536000s;
        add_header Pragma "public";
        add_header Cache-Control "public, must-revalidate, proxy-revalidate";
        add_header X-Powered-By "W3 Total Cache/0.9.2.4";
        try_files $uri $uri/ @rewrite;
    }
    # END W3TC Browser Cache
}
```

### Gitweb

With Nginx, you need to allow certain extensions in php-fpm:

```bash {linenos=table}
[...]
security.limit_extensions = .php .php3 .php4 .php5 .cgi
[...]
```

And here's a configuration example to adapt to your needs:

```bash {linenos=table}
server {
    listen 80;
    listen 443 ssl;

    ssl_certificate /etc/nginx/ssl/deimos.fr/server-unified.crt;
    ssl_certificate_key /etc/nginx/ssl/deimos.fr/server.key;
    ssl_session_timeout 5m;

    server_name git.deimos.fr;
    root /usr/share/gitweb/;

    access_log /var/log/nginx/git.deimos.fr_access.log;
    error_log /var/log/nginx/git.deimos.fr_error.log;

    index gitweb.cgi;

    location /gitweb.cgi {
        fastcgi_cache mycache;
        fastcgi_cache_key $request_method$host$request_uri;
        fastcgi_cache_valid any 1h;
        include fastcgi_params;
        fastcgi_pass  unix:/run/fcgiwrap.socket;
    }
}
```

Next we'll need a cgi wrapper:

```bash
aptitude install fcgiwrap
```

Create the link, then reload the configuration:

```bash
cd /etc/nginx/sites-enabled
ln -s /etc/nginx/sites-available/git.deimos.fr .
/etc/init.d/nginx reload
/etc/init.d/fcgiwrap restart
/etc/init.d/php5-fpm reload
```

### Git

I spent quite a bit of time making Git over http(s) and Gitweb coexist, but I got it working.

{{< alert context="info" text="Prefer the <a href=\"https://wiki.deimos.fr/Gitweb_:_Installation_et_configuration_d'une_interface_web_pour_git#Nginx\">Gitweb</a> method only if you don't need git over http(s)" />}}

Here's the method I used:

```bash {linenos=table}
server {
    listen 80;
    listen 443 ssl;

    ssl_certificate /etc/nginx/ssl/deimos.fr/server-unified.crt;
    ssl_certificate_key /etc/nginx/ssl/deimos.fr/server.key;
    ssl_session_timeout 5m;

    server_name git.deimos.fr;
    root /usr/share/gitweb/;

    access_log /var/log/nginx/git.deimos.fr_access.log;
    error_log /var/log/nginx/git.deimos.fr_error.log;

    index gitweb.cgi;

    # Drop config
    include drop.conf;

    # Git over https
    location /git/ {
        alias /var/cache/git/;
        if ($scheme = http) {
            rewrite ^ https://$host$request_uri permanent;
        }
    }

    # Gitweb
    location ~ gitweb\.cgi {
        fastcgi_cache mycache;
        fastcgi_cache_key $request_method$host$request_uri;
        fastcgi_cache_valid any 1h;
        include fastcgi_params;
        fastcgi_pass  unix:/run/fcgiwrap.socket;
    }
}
```

Here I have my git over https working (http is redirected to https) and my gitweb also since everything that matches gitweb.cgi is matched. Now, for the git part, we'll need to authorize the repositories we want. For this, we'll need to rename a file in our repository and run a command:

```bash
cd /var/cache/git/myrepo.git
hooks/post-update{.sample,}
su - www-data -c 'cd /var/cache/git/myrepo.git && /usr/lib/git-core/git-update-server-info'
```

Replace www-data with the user who has rights to the repository. Use www-data so that nginx has the rights. Then, you have permission to clone:

```bash
git clone http://www.deimos.fr/git/deimosfr.git deimosfr
```

### Piwik

For Piwik, here's the configuration:

```bash {linenos=table}
server {
    include listen_port.conf;
    listen 443 ssl;

    ssl_certificate /etc/nginx/ssl/deimos.fr/server-unified.crt;
    ssl_certificate_key /etc/nginx/ssl/deimos.fr/server.key;
    ssl_session_timeout 5m;

    server_name piwik.deimos.fr;
    root /usr/share/nginx/www/deimos.fr/piwik;
    index index.php;

    access_log /var/log/nginx/piwik.deimos.fr_access.log;
    error_log /var/log/nginx/piwik.deimos.fr_error.log;

    # Drop config
    include drop.conf;

    location / {
        try_files $uri $uri/ /index.php?$args;
    }

    location ~ \.php$ {
        fastcgi_cache mycache;
        fastcgi_cache_key $request_method$host$request_uri;
        fastcgi_cache_valid any 1h;
        include fastcgi_params;
        fastcgi_pass unix:/var/run/php5-fpm.sock;
        fastcgi_intercept_errors on;
    }

    location ~* \.(js|css|png|jpg|jpeg|gif|ico)$ {
        expires max;
        log_not_found off;
    }
}
```

### ownCloud

#### ownCloud 4.X

Here's the configuration for ownCloud 4.X with Nginx:

```bash {linenos=table}
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

```bash {linenos=table}
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

### Firefox Sync Server

For Firefox Server Sync, it's proxification. We forward requests to another service:

```bash {linenos=table,hl_lines=["5-6","14-17"]}
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

### Selfoss

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

### Seafile

```bash {linenos=table}
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

### Gogs

For Gogs, it's simple because all you need to set is a proxy pass:

```bash {linenos=table}
server {
    include listen_port.conf;
    server_name git.deimos.fr;
    # Force redirect http to https
    rewrite ^ https://$http_host$request_uri? permanent;
}

server {
    include ssl/deimos.fr_ssl.conf;

    server_name git.deimos.fr;

    access_log /var/log/nginx/git.deimos.fr_access.log;
    error_log /var/log/nginx/git.deimos.fr_error.log;

    location / {
        proxy_pass    http://127.0.0.1:3000;
    }
}
```

## FAQ

### How to Analyze Rules

You can see the behavior of your rules (e.g., 301 redirects) using curl:

```bash
> curl -I "http://wiki.deimos.fr"
HTTP/1.1 301 Moved Permanently
Server: nginx
Content-Type: text/html
Location: http://wiki.deimos.fr/index.php
Content-Length: 178
Accept-Ranges: bytes
Date: Wed, 16 Jan 2013 22:17:05 GMT
X-Varnish: 2038764987
Age: 0
Via: 1.1 varnish
Connection: keep-alive
```

This allows you to analyze how it works and detect potential problems.

### client intended to send too large body

If you have this kind of message in your logs:

```
nginx client intended to send too large body
```

It simply means that Nginx doesn't allow you to send files as large as you want. You'll need to use this variable in the 'server' part:

```
client_max_body_size 10M;
```

Replace 10M with the maximum value you want to allow.

### Access denied.

If you get this type of error message, it's due to PHP-FPM having issues with extension management. You can either add the extension in question:

```php {linenos=table}
security.limit_extensions = .php .php3 .php4 .php5 .cgi
```

or disable this security (not recommended):

```php {linenos=table}
security.limit_extensions = false
```

## Resources

- [Documentation on installing Nginx HTTP Server](/pdf/nginx_http_server_fast-cgi_and_xcache.pdf)
- http://www.if-not-true-then-false.com/2011/nginx-and-php-fpm-configuration-and-optimizing-tips-and-tricks/
- [Nagios check php-fpm](https://github.com/regilero/check_phpfpm_status)
- [^1]: http://nginx.org/en/docs/http/ngx_http_limit_req_module.html
