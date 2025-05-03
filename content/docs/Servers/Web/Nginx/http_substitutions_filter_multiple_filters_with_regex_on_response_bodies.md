---
weight: 999
url: "/Http_substitutions_filter_\\:_multiple_filters_with_regex_on_response_bodies/"
title: "HTTP Substitutions Filter: Multiple Filters with Regex on Response Bodies"
description: "Learn how to implement HTTP substitutions filter for Nginx to perform multiple regex filters on response bodies. This guide covers compilation, installation, and configuration on Debian systems."
categories: ["Linux", "Debian", "Nginx"]
date: "2013-04-26T15:22:00+02:00"
lastmod: "2013-04-26T15:22:00+02:00"
tags: ["Nginx", "HTTP Filters", "Web Server", "Regex", "Module"]
toc: true
---

![Nginx Logo](/images/nginx-logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 1.2.1 |
| **Operating System** | Debian 7 |
| **Website** | [nginx_substitutions_filter Website](https://github.com/yaoweibin/ngx_http_substitutions_filter_module) |
| **Last Update** | 26/04/2013 |
{{< /table >}}

## Introduction

nginx_substitutions_filter is a filter module which can do both regular expression and fixed string substitutions on response bodies. This module is quite different from the Nginx's native Substitution Module. It scans the output chains buffer and matches string line by line, just like Apache's mod_substitute.[^1]

I've played with classic substitution module but due to limitations (only one pattern, no regex), it wasn't easy to do all I wanted to do. That's why I've searched a better module and to add it on my Debian server. Unfortunately, it's not currently available in Nginx packages on Debian. That's why I needed to create new packages with it built in.

To get a quick understand of what this module is able to do, we're going to take an example. Let's say we need to add a CSS style on every pages on our website. You've got multiple solution to do it and most of them need to modify the application code. The problem is, each time you will get an application upgrade, you'll need to think about this modifications. The other solution (if possible) is to create an extension for you application, but it's boring to manage.

With that solution, you can change the end of the head banner and add any style or JS that you want. It works with everything in fact. As described, it's a string replacement solution. This extension is really powerfull. I'll explain in the next days how I use it for my own usage.

## Prerequisites

We first need to download package source and install all dependencies to recompile Nginx:

```bash
mkdir ~/nginx_new ; cd ~/nginx_new
aptitude install devscripts dch
apt-get build-dep nginx-extras
apt-get source nginx-extras
```

As I use nginx-extras package, I take this one but you can take only nginx if you want.

## Compilation

Now let's get sources from the official site:

```bash
cd ~/nginx_new/nginx-1.2.1/debian/modules/
git clone git://github.com/yaoweibin/ngx_http_substitutions_filter_module.git
```

Then add this line to make it compiled on the next step (`~/nginx_new/nginx-1.2.1/debian/rules`):

```bash {linenos=table,hl_lines=[45]}
config.status.extras: config.env.extras config.sub config.guess
    cd $(BUILDDIR_extras) && CFLAGS="$(CFLAGS)" CORE_LINK="$(LDFLAGS)" ./configure  \
        --prefix=/etc/nginx \
        --conf-path=/etc/nginx/nginx.conf \
        --error-log-path=/var/log/nginx/error.log \
        --http-client-body-temp-path=/var/lib/nginx/body \
        --http-fastcgi-temp-path=/var/lib/nginx/fastcgi \
        --http-log-path=/var/log/nginx/access.log \
        --http-proxy-temp-path=/var/lib/nginx/proxy \
        --http-scgi-temp-path=/var/lib/nginx/scgi \
        --http-uwsgi-temp-path=/var/lib/nginx/uwsgi \
        --lock-path=/var/lock/nginx.lock \
        --pid-path=/var/run/nginx.pid \
        --with-pcre-jit \
        --with-debug \
        --with-http_addition_module \
        --with-http_dav_module \
        --with-http_flv_module \
        --with-http_geoip_module \
        --with-http_gzip_static_module \
        --with-http_image_filter_module \
        --with-http_mp4_module \
        --with-http_perl_module \
        --with-http_random_index_module \
        --with-http_realip_module \
        --with-http_secure_link_module \
        --with-http_stub_status_module \
        --with-http_ssl_module \
        --with-http_sub_module \
        --with-http_xslt_module \
        --with-ipv6 \
        --with-sha1=/usr/include/openssl \
        --with-md5=/usr/include/openssl \
        --with-mail \
        --with-mail_ssl_module \
        --add-module=$(MODULESDIR)/nginx-auth-pam \
        --add-module=$(MODULESDIR)/chunkin-nginx-module \
        --add-module=$(MODULESDIR)/headers-more-nginx-module \
        --add-module=$(MODULESDIR)/nginx-development-kit \
        --add-module=$(MODULESDIR)/nginx-echo \
        --add-module=$(MODULESDIR)/nginx-http-push \
        --add-module=$(MODULESDIR)/nginx-lua \
        --add-module=$(MODULESDIR)/nginx-upload-module \
        --add-module=$(MODULESDIR)/nginx-upload-progress \
        --add-module=$(MODULESDIR)/ngx_http_substitutions_filter_module \
        --add-module=$(MODULESDIR)/nginx-upstream-fair \
        --add-module=$(MODULESDIR)/nginx-dav-ext-module \
            $(CONFIGURE_OPTS) >$@
    touch $@
```

```bash
dch -l local 'New version with nginx_substitutions_filter included'
```

Now you're ready to compile and package it automatically:

```bash
cd ../..
debuild -us -uc
```

## Installation

Now it's easy, let's install:

```bash
cd ~/nginx_new
dpkg -i nginx-extras_1.2.1-2.2local1_amd64.deb nginx-common_1.2.1-2.2local1_all.deb
```

## Configuration

Let's take the default configuration for instance and then add to your location this kind of string replacements:

```bash
[...]
    location / {
        # First attempt to serve request as file, then
        # as directory, then fall back to displaying a 404.
        try_files $uri $uri/ /index.html;
        subs_filter_types text/html text/css text/xml;
        subs_filter st(\d*).example.com $1.example.com ir;
        subs_filter a.example.com s.example.com;
        # Uncomment to enable naxsi on this location
        # include /etc/nginx/naxsi.rules
    }
[...]
```

As you can see, regex works and simple remplacement strings works as well.

## References

[^1]: https://github.com/yaoweibin/ngx_http_substitutions_filter_module
