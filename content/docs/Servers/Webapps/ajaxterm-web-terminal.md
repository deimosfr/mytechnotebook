---
weight: 999
url: "/ajaxterm-utiliser-un-terminal-en-web/"
title: "Ajaxterm: Using a Terminal via Web"
description: "How to set up Ajaxterm to access a Linux terminal through a web browser with proper security settings."
categories: ["Linux", "Server", "Web", "Security"]
date: "2011-04-22T13:44:00+02:00"
lastmod: "2011-04-22T13:44:00+02:00"
tags: ["Terminal", "Web", "Ajaxterm", "Lighttpd", "SSH", "Proxy"]
toc: true
---

## Introduction

[Ajaxterm](https://antony.lesuisse.org/software/ajaxterm/) allows you to have a terminal through a web page. For example, if you're at a restricted site that doesn't allow you to use SSH, there's a good chance you can still use HTTPS. This is where Ajaxterm becomes handy, as it allows you to connect to the machine hosting Ajaxterm.

The downside is security. If you don't configure this properly, your server can quickly end up in other hands than yours. I will address this point here as well to prevent that from happening.

## Installation

For the installation, I kept it simple:

```bash
wget http://antony.lesuisse.org/software/ajaxterm/files/Ajaxterm-0.10.tar.gz
tar -xzvf Ajaxterm-0.10.tar.gz
mv Ajaxterm-0.10 /usr/share/ajaxterm
```

Next, I want Ajaxterm to start automatically at boot as a daemon with the www-data user (or the web user if you're not on Debian) which corresponds to PID 33. I'll add this line to the following file:

```bash
# In /etc/rc.local
...
python /usr/share/ajaxterm/ajaxterm.py -u33 -d
...
```

For the more courageous, you can create an init script, which would be cleaner.

## Configuration

There's no special configuration needed for Ajaxterm; by default everything is nice. However, you'll need to configure a web proxy to redirect from port 443 to 8022 on localhost.

### Lighttpd

If you use Lighttpd, we'll configure it to load the proxy module:

```perl
# In /etc/lighttpd/conf-available/10-proxy.conf
 ...
 server.modules   += ( "mod_proxy" )
 ...
```

Then we'll create a configuration for Ajaxterm that will allow us access if:
* The host corresponds to www.deimos.fr
* The URL corresponds to term or terminal

Additionally, some security points have been added, such as:
* If port 80 is used, we redirect to port 443
* The user must enter the correct credentials for a htaccess file

Here's the configuration:

```perl
### Ajaxterm
$HTTP["host"] =~ "www\.deimos\.fr" {
    $HTTP["url"] =~ "^/(term|terminal|ajaxterm)" {
        $SERVER["socket"] == ":80" {
            url.redirect = ( "" => "https://www.deimos.fr/ajaxterm/" )
        }
        auth.require = ( "" =>
            (
            "method" => "digest",
            "realm" => "Authorized users only",
            "require" => "valid-user"
            )
        )
        proxy.server = ( "" =>
            (
                ( "host" => "127.0.0.1",
                  "port" => 8022
                )
            )
        ),
    }
}
```

So to summarize, the security is pretty good, but we're not far from the vital minimum. Ideally, you would add even more security options, but I won't cover them here.

To finish, we'll enable our new modules:

```bash
lighty-enable-mod proxy
lighty-enable-mod ajaxterm
```

All that's left is to restart Lighttpd and connect: https://www.deimos.fr/ajaxterm :-)
