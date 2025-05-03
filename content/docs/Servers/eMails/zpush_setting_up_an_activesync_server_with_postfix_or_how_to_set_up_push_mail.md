---
weight: 999
url: "/Z-Push_\\:_Avoir_un_serveur_ActiveSync_avec_Postfix_(ou_comment_faire_du_push_mail)/"
title: "Z-Push: Setting Up an ActiveSync Server with Postfix (or How to Set Up Push Mail)"
description: "Learn how to set up a push mail server using Z-Push with Postfix as an alternative to Microsoft Exchange, compatible with iPhone and Windows Mobile devices."
categories: ["Apache", "Linux"]
date: "2008-07-24T08:12:00+02:00"
lastmod: "2008-07-24T08:12:00+02:00"
tags: ["PHP", "Servers", "Email"]
toc: true
---

## Introduction

Push server technology is very trendy nowadays, especially with the iPhone which now allows connections (just like Windows Mobile) to an Exchange-like push server. The problem is that for the open-source world, Exchange is not an option. I found a well-developed project on SourceForge called [Z-Push](https://z-push.sourceforge.net) that works perfectly with Postfix.

![](/images/1204022285.avif)

## Installation

We'll download the latest version from [https://z-push.sourceforge.net](https://z-push.sourceforge.net) and extract it to `/var/www`:

```bash
tar zxvf z-push-<version>.tar.gz -C /var/www
```

Now we'll apply the proper permissions:

```bash
chmod 777 /var/www/z-push/state
chmod 755 /var/www/z-push/state
chown www-data. /var/www/z-push
```

## Configuration

### Apache

We need to configure Apache to redirect /Microsoft-Server-ActiveSync to /var/www/z-push/index.php. There are two options:

* Using alias:

```apache
Alias /Microsoft-Server-ActiveSync /var/www/z-push/index.php
```

* Using VirtualHost

Add this to your virtualhost configuration (`/etc/apache2/sites-enabled/000-default`):

```apache
       <Location /z-push>
               Options Indexes FollowSymLinks MultiViews
               Order allow,deny
               allow from all
               RedirectMatch ^/Microsoft-Server-ActiveSync /var/www/z-push/index.php
       </Location
```

### PHP

Now we need to make some modifications to your PHP configuration:

```ini
php_flag magic_quotes_gpc off
register_globals off
magic_quotes_runtime off
short_open_tag on
```

Restart the Apache server:

```bash
/etc/init.d/apache2 restart
```

### Z-Push

Now we just need to edit a few fields in the configuration:

```php
...
date_default_timezone_set("Europe/Paris"
...
$BACKEND_PROVIDER = "BackendIMAP";
...
define('IMAP_SERVER', 'deimos.fr');
...
```

Now, you just need to test your configuration by connecting to your server: http://<serverip>/Microsoft-Server-ActiveSync  
If you get a login/password prompt, enter your IMAP account credentials, and if you get this message, it means it's working :-)

```
GET not supported
This is the z-push location and can only be accessed by Microsoft ActiveSync-capable devices.
```

## FAQ

### I have problems but don't know where they're coming from, how can I debug?

Simply create a debug file:

```bash
touch /var/www/z-push/debug.txt
chmod 777 /var/www/z-push/debug.txt
```
