---
weight: 999
url: "/Nginx_Git_et_Gitweb/"
title: "Nginx Git and Gitweb"
description: "Guide to make Git over HTTP(S) and Gitweb coexist in an Nginx setup"
categories: ["Nginx", "Linux"]
date: "2013-01-04T17:28:00+02:00"
lastmod: "2013-01-04T17:28:00+02:00"
tags: ["Nginx", "Git", "Gitweb", "Web Servers", "Development"]
toc: true
---

I spent a lot of time figuring out how to make Git over HTTP(S) and Gitweb coexist, but finally got it working.

{{< alert context="info" text="Consider using the <a href=\"/Gitweb_:_Installation_et_configuration_d'une_interface_web_pour_git/#Nginx\">Gitweb</a> method only if you don't need Git over HTTP(S)" />}}

## Configuration Method

Here's the method I used:

```nginx {linenos=table}
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

With this configuration, Git over HTTPS is working (HTTP is redirected to HTTPS) and Gitweb is working too since everything matching gitweb.cgi is correctly routed.

## Repository Configuration

For the Git part, we need to authorize the repositories we want to expose. For this, we need to rename a file in our repository and run a command:

```bash
cd /var/cache/git/myrepo.git
hooks/post-update{.sample,}
su - www-data -c 'cd /var/cache/git/myrepo.git && /usr/lib/git-core/git-update-server-info'
```

Replace www-data with the user that has permissions on the repository. Use www-data so that Nginx has the necessary rights.

## Using the Repository

After this setup, you can clone your repository:

```
git clone http://www.deimos.fr/git/deimosfr.git deimosfr
```
