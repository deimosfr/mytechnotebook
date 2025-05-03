---
weight: 999
url: "/Webdav_avec_Lighttpd/"
title: "WebDAV with Lighttpd"
description: "A guide on how to set up WebDAV with Lighttpd server for uploading files to a web server."
categories: ["Linux"]
date: "2009-08-19T07:08:00+02:00"
lastmod: "2009-08-19T07:08:00+02:00"
tags: ["Servers", "Lighttpd", "WebDAV", "Authentication"]
toc: true
---

## Introduction

WebDAV is great once you've tried it. We'll see how to set it up with Lighttpd. I need it to upload images for my photo album.

## Installation

For the installation, it's straightforward:

```bash
apt-get install lighttpd-mod-webdav
```

## Configuration

### Lighttpd

Let's enable our newly installed module:

```bash
lighty-enable-mod auth
lighty-enable-mod webdav
```

Edit the Lighttpd configuration file and uncomment this section:

```perl
server.modules              <nowiki>=</nowiki> (
...
            "mod_webdav",
...
)
```

We'll add the folder we want to have WebDAV access to below.

* If you want to restrict by IP, use this example (`/etc/lighttpd/lighttpd.conf`):

```perl
...
$HTTP["host"] =~ "^webdav\.(deimos\.fr)" {
    server.document-root = "/var/www/photos/galleries"
    alias.url = ( "(.*)" => "/var/www/photos/galleries" )
    $HTTP["remoteip"] != "192.168.0.0/24" {
        webdav.activate = "enable"
        webdav.is-readonly = "disable"
        auth.backend = "htpasswd"
        auth.backend.htpasswd.userfile = "/var/www/photos/galleries/.htpasswd"
        auth.require = ( "" => ( "method" => "basic",
            "realm" => "webdav",
            "require" => "valid-user" ) ) 
    }   
}
...
```

* If you prefer something simpler, use this example (`/etc/lighttpd/lighttpd.conf`):

```perl
$["remoteip"] == "^webdav\.(deimos\.fr)" {
    alias.url += ( "/webdav" => "/var/www/photos/galleries" )
    $HTTP["url"] =~ "^webdav($|/)" {
        dir-listing.activate = "enable"
        webdav.activate = "enable"
        webdav.is-readonly = "disable"
        auth.backend = "htpasswd"
        auth.backend.htpasswd.userfile = "/var/www/photos/galleries/.htpasswd"
        auth.require = ("" => "method" => "basic",
             "realm" => "webdav",
             "require" => "valid-user" ) )
    }
}
```

Finally, restart Lighttpd.

Of course, if you want to go a bit further with all this, you can also connect it to an LDAP directory, for example.

### WebDAV

Now let's configure the file we're interested in (here `/var/www/photos`). First, we'll create an htpasswd file:

```bash
htpasswd -c /var/www/photos/galleries/.htpasswd www-data
```

I deliberately chose www-data here. But it's preferable to use a specific user. Let's also set the correct permissions:

```bash
chown www-data. /var/www/photos/galleries/.htpasswd
```

## Resources
- [How To Set Up WebDAV With Lighttpd](/pdf/how_to_set_up_webdav_with_lighttpd.pdf)
