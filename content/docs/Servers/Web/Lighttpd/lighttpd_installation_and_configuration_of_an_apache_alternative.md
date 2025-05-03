---
weight: 999
url: "/Lighttpd_\\:_Installation_et_configuration_d'une_alternative_d'Apache/"
title: "Lighttpd: Installation and configuration of an Apache alternative"
description: "A guide for installing and configuring Lighttpd as a fast and flexible alternative to Apache with various configuration examples"
categories: ["Debian", "Linux", "Servers"]
date: "2011-05-10T11:05:00+02:00"
lastmod: "2011-05-10T11:05:00+02:00"
tags: ["Web Server", "Lighttpd", "PHP", "SSL", "FastCGI"]
toc: true
---

## Introduction

LigHTTPd (or "lighty") is a secure, fast, and flexible HTTP server.

Its speed comes from having a smaller memory footprint than other HTTP servers, as well as intelligent CPU load management.

Many languages like PHP, Perl, Ruby, and Python are supported via FastCGI.

The main disadvantage of LigHTTPd is that it has only one configuration file and doesn't support .htaccess files (though there are alternatives): directives must be in the configuration file. This is also an advantage, as the server administrator only has one file to manage.

In April 2007, LigHTTPd entered the Top 5 of the most widely used web servers.

## Installation

To install it, it's simple:

```bash
apt-get install lighttpd
```

### Installation of the PHP module

```bash
apt-get install php5-cgi php5
```

To enable it, do:

```bash
lighty-enable-mod fastcgi
```

We'll also add these few lines at the end of the configuration file:

```perl
...
fastcgi.server = ( ".php" => ((
                     "bin-path" => "/usr/bin/php5-cgi",
                     "socket" => "/tmp/php.socket"
                 )))
```

Then reload Lighttpd:

```bash
/etc/init.d/lighttpd force-reload
```

## Configuration

### Forcing SSL redirections

It can be useful to redirect part of your traffic to SSL. This is my case for mediawiki which doesn't take SSL into account during authentication (or I don't know how to do it, but it doesn't matter, the idea here is to have a small exercise):

```perl
# Mediawiki secure auth
$SERVER["socket"] == ":80" {
    $HTTP["host"] =~ "(.*)" {
        url.redirect = ( "^/(.*Connexion\&returnto.*)" => "https://%1/$1" )
    }
}
```

Here I'm telling it:

- Everything that comes in on port 80
- From any host
- If it arrives on a page containing 'Connexion&returnto'
- Then it is redirected to https at the same location as the previous one

### Differentiating logs

For example, I want to differentiate the logs of my blog as well as the rest of my domain and my wiki to be able to process them correctly with awstats. Here's how to do it:

```perl
...
# Deimos Domain
$HTTP["host"] =~ "deimos\.fr" {
    server.document-root = "/var/www/deimos.fr"
    # Default logs
    accesslog.filename = "/var/log/lighttpd/access-blog_and_co.log"
    # Wiki logs
    $HTTP["url"] =~ "^/blocnotesinfo/" {
        accesslog.filename = "/var/log/lighttpd/access-deimos-wiki.log"
    }
}
...
```

### Blocking access to your site via Internet Explorer

It is possible to block access to all kinds of browsers. If like me, you're not friends with IE (which breaks your PNGs in version 6, doesn't respect standards, breaks CSS, etc.), it might be useful to block it and kindly indicate to the user to download Firefox as soon as possible.

I did it under [Apache]({{< ref "docs/Servers/Web/Apache/apache_2_installation_and_configuration.md#blocking-access-to-your-site-via-internet-explorer" >}}) and here I couldn't pass up the opportunity to do it for Lighttpd, so here's the solution:

```perl
$HTTP["useragent"] =~ ".*(MSIE|Opera).*" {
    $HTTP["url"] !~ "^/ie.html" {
        url.redirect = ( "" => "/ie.html" )
    }
}
```

All that's left is to create the ie.html file and put your nice text in it (you can also make a simple text file). Here's what I use (`/var/www/ie.html`):

```html
<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd">
<html>
  <head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <title>Access Forbidden with Internet Explorer</title>
  </head>
  <body>
    Dear Internet User,<br />
    <br />
    This site is not accessible using Internet Explorer that you are using.<br />
    You must now understand that times are changing.<br />
    <br />
    You are currently using a browser (Internet Explorer) that does not
    respect<br />
    the <a href="http://www.w3.org">standards</a> and has the monopoly, due to
    its mandatory presence in your<br />
    dear OS (Windows). That said, you may be at work and don't have a choice of
    OS.<br />
    <br />
    However, Internet Explorer should no longer be used when there are all kinds
    of<br />
    other free, open source browsers that respect standards!<br />
    But since you don't seem to be aware, that's okay, let me help you.<br />
    <br />
    To start, you can download a clean browser such as
    <a href="http://www.mozilla.org">Firefox</a>.<br />
    This would already help you get on the right track and will also allow you
    to access<br />
    my site.<br />
    <br />
    Still in your new quest to the light side of the force, you should switch
    to<br />
    a free and open source OS (<a href="http://www.ubuntu.com">Ubuntu</a> for
    example) which will surely make you happy.<br />
    <br />
    I invite you to take matters into your own hands as soon as possible.<br />
    <br />
    Sincerely,<br />
    Pierre (aka Deimos)
  </body>
</html>
```

### Access restriction by login and password

Htaccess files don't exist in Lighttpd, but there is an equivalent. Check before you start that the mod_auth module is loaded.
First, we'll generate (with -c for the first time, like htaccess) a file containing the credentials for authorization to view a particular site:

```bash
htdigest -c /etc/lighttpd/.passwd 'Authorized users only' deimos
```

Here I'm creating the user deimos. The realm (here 'Authorized users only') will allow us to differentiate between the different login/password files we'll have since we can only specify one for the entire server.

Then add these lines to the global lighttpd configuration:

```perl
auth.backend = "htdigest"
auth.backend.htdigest.userfile = "/etc/lighttpd/.passwd"
auth.debug = 2
```

Then I add protection at the location I'm interested in:

```perl
auth.require = ( "/docs/" =>
   (
      "method" => "digest",
      "realm" => "Authorized users only",
      "require" => "valid-user"
   )
)
```

Restart lighty and you're done.

### Adding compression

To speed up the display of your website, it's recommended to enable compression. By default, not enough elements are compressed, so you'll need to define more in the standard configuration.
_Note: Enabling compression will put a bit more load on your server's CPU/RAM._

Edit your lighty configuration to configure the compress module:

```perl
server.modules = (
            "mod_alias",
            "mod_compress",
            "mod_rewrite",
            "mod_redirect",
            "mod_cgi",
#           "mod_usertrack",
#           "mod_expire",
#           "mod_flv_streaming",
#           "mod_evasive"
)
...
#### compress module
compress.cache-dir          = "/var/cache/lighttpd/compress/"
#compress.filetype           = ("text/plain", "text/html", "application/x-javascript", "text/css")
compress.filetype           = ("application/x-javascript", "application/javascript", "text/javascript", "text/x-js", "text/css", "text/html", "text/plain")
...
```

PHP5 compression will save us precious seconds on page display. To enable it, edit the following file and set the parameter to "on":

```apache
; Transparent output compression using the zlib library
; Valid values for this option are 'off', 'on', or a specific buffer size
; to be used for compression (default is 4KB)
; Note: Resulting chunk size may vary due to nature of compression. PHP
;   outputs chunks that are few hundreds bytes each as a result of
;   compression. If you prefer a larger chunk size for better
;   performance, enable output_buffering in addition.
; Note: You need to use zlib.output_handler instead of the standard
;   output_handler, or otherwise the output will be corrupted.
; http://php.net/zlib.output-compression
zlib.output_compression = on
```

Then restart your web server for the configuration to be applied.

## FAQ

### Symbol `FamErrlist' has different size in shared object, consider re-linking

If you get this kind of error when restarting Lighttpd:

```
Stopping web server: lighttpd.
Starting web server: lighttpd/usr/sbin/lighttpd: Symbol `FamErrlist' has different size in shared object, consider re-linking
2010-08-20 21:58:02: (network.c.345) can't bind to port: :: 80 Address already in use
 failed!
```

This is due to a bug with IPv6 on Debian. To solve the problem, you need to install libfam0 which also causes these kinds of bugs:

```bash
aptitude install libfam0
```

and in the Lighttpd configuration, configure these lines like this:

```
#include_shell "/usr/share/lighttpd/use-ipv6.pl"
server.socket = "[::]"
```

All that's left is to restart Lighttpd :-)

### sockets disabled, connection limit reached

If your server stops abruptly and you get this kind of message:

```
2011-05-06 14:18:24: (server.c.1398) [note] sockets disabled, connection limit reached
2011-05-06 14:19:27: (server.c.1512) server stopped by UID = 0 PID = 19214
```

It's simply that you've exceeded the maximum number of file descriptors. You just need to increase the current value. To find out which one is currently applied (default 1024), just run this command:

```bash
> cat /proc/`ps ax | grep lighttpd | grep -v grep | awk -F " " '{print $1}'`/limits |grep "Max open files"
Max open files            1024                 1024                 files
```

Here it is limited to 1024. To increase this value, edit the lighty conf file and increase the value:

```
server.max-fds = 2048
```

All that's left is to restart lighty.

### backend is overloaded; we'll disable it for 1 seconds and send the request to another backend instead: reconnects

If you get this kind of message with a nice 500 error:

```
2011-05-09 09:06:07: (mod_fastcgi.c.2764) fcgi-server re-enabled:  0 /tmp/php.socket
2011-05-09 09:06:07: (mod_fastcgi.c.2764) fcgi-server re-enabled:  0 /tmp/php.socket
2011-05-09 09:06:07: (mod_fastcgi.c.2764) fcgi-server re-enabled:  0 /tmp/php.socket
2011-05-09 09:06:07: (mod_fastcgi.c.2764) fcgi-server re-enabled:  0 /tmp/php.socket
2011-05-09 09:06:10: (mod_fastcgi.c.3001) backend is overloaded; we'll disable it for 1 seconds and send the request to another backend instead: reconnects: 0 load: 521
2011-05-09 09:06:10: (mod_fastcgi.c.3001) backend is overloaded; we'll disable it for 1 seconds and send the request to another backend instead: reconnects: 1 load: 521
2011-05-09 09:06:10: (mod_fastcgi.c.3001) backend is overloaded; we'll disable it for 1 seconds and send the request to another backend instead: reconnects: 2 load: 521
2011-05-09 09:06:10: (mod_fastcgi.c.3001) backend is overloaded; we'll disable it for 1 seconds and send the request to another backend instead: reconnects: 3 load: 521
```

It's because your fastcgi limits have been exceeded. To solve the problem, here's the configuration to modify:

```
fastcgi.server = (
    ".php" => ((
        "bin-path" => "/usr/bin/php5-cgi",
        "socket" => "/tmp/php.socket",
        "max-procs" => 10,
        "bin-environment" => (
            "PHP_FCGI_CHILDREN" => "10",
             "PHP_FCGI_MAX_REQUESTS" => "500"
        )
    )
))
```

Then adjust the lines above to accommodate the maximum number of clients you want, as well as the number of instances. You can test your new configuration with the "ab" command.

More info on this page: [https://redmine.lighttpd.net/wiki/1/Docs:PerformanceFastCGI](https://redmine.lighttpd.net/wiki/1/Docs:PerformanceFastCGI)

## Resources

- [Documentation on installing Lighttpd With PHP5 And MySQL Support](/pdf/installing_lighttpd_with_php5_and_mysql_support.pdf)
- [Documentation on reducing Apache's load with lighttpd](/pdf/reduce_apache's_load_with_lighttpd.pdf)
- [Optimizing Lighttpd server load](/pdf/lighttpd-web-server.pdf)
- [https://redmine.lighttpd.net/wiki/lighttpd/Docs](https://redmine.lighttpd.net/wiki/lighttpd/Docs)
- [https://redmine.lighttpd.net/wiki/lighttpd/Docs:ConfigurationOptions](https://redmine.lighttpd.net/wiki/lighttpd/Docs:ConfigurationOptions)
- [Installing Lighttpd With PHP5 And MySQL Support On Debian Lenny](/pdf/installing_lighttpd_with_php5_and_mysql_support_on_debian_lenny.pdf)
