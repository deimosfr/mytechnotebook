---
weight: 999
url: "/Pootle_\\:_simple_translation_tool/"
title: "Pootle: Simple Translation Tool"
description: "A guide on how to install and configure Pootle, a simple online translation tool that makes the translation process easier and allows crowd-sourced translations."
categories: ["Debian", "Development", "Servers"]
date: "2013-09-06T15:45:00+02:00"
lastmod: "2013-09-06T15:45:00+02:00"
tags: ["Pootle", "Translation", "Apache", "Python", "Web Applications"]
toc: true
---

![Pootle](/images/pootle_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 2.5.0 |
| **Operating System** | Debian 6 |
| **Website** | [Pootle Website](https://pootle.translatehouse.org/) |
| **Last Update** | 06/09/2013 |
{{< /table >}}

## Introduction

[Pootle](https://pootle.translatehouse.org/)[^1] is an online tool that makes the process of translating so much simpler. It allows crowd-sourced translations, easy volunteer contribution and gives statistics about the ongoing work.

Pootle is built using the powerful API of the Translate Toolkit and the Django framework:

![Pootle screenshot](/images/pootle_screenshot.avif)

## Installation

You can install Poole directly from the packages, but you won't have the latest version. To get it we'll need to install required packages and download all dependencies via PIP:

```bash
aptitude install python-pip gcc python2.6-dev libxslt1-dev python-virtualenv
```

Then install and create the python virtual environment:

```bash
virtualenv /var/www/pootle/env/
source /var/www/pootle/env/bin/activate
```

And finally install Pootle:

```bash
pip install pootle
```

## Configuration

### Pootle

Then you can initiate the configuration:

```bash
pootle init /var/www/pootle/pootle.conf
```

Edit the configuration to change allowed host to access to the web frontend (`/var/www/pootle/pootle.conf`):

```bash
# A list of strings representing the host/domain names that this Pootle server
# can serve. This is a Django's security measure. More details at
# https://docs.djangoproject.com/en/dev/ref/settings/#allowed-hosts
ALLOWED_HOSTS = ['127.0.0.1']
# Allow all
# ALLOWED_HOSTS = ['*']
```

{{< alert context="warning" text="Edit the configuration and setup the MySQL database instead of the SQLite by default" />}}

And launch it:

```bash
pootle --config=/var/www/pootle/pootle.conf start
```

You'll now get an access to the web interface: http://127.0.0.1:8000  
Credentials are: admin/admin

{{< alert context="info" text="The first launch will take a few minutes as it populate the database" />}}

### Apache

To avoid typing port number on the URL, you can use Apache mod proxy.

```bash
aptitude install apache2 apache2-utils apache2.2-common libapache2-mod-proxy-html
```

Then activate modules:

```bash
a2enmod proxy_connect
a2enmod proxy_http
a2enmod proxy_html
```

And restart Apache.

Then configure your apache (`/etc/apache2/sites-enabled/pootle`):

```apache
<VirtualHost *:80>
    ServerAdmin webmaster@localhost
    ServerName pootle.deimos.fr

    # Pootle
    <Location />
    	order deny,allow
    	allow from all
    	ProxyPass http://server:8000/
    	ProxyPassReverse http://server:8000/
    </Location>
</VirtualHost>
```

You'll need to enable this new site and restart Apache.

## References

[^1]: http://pootle.translatehouse.org/
