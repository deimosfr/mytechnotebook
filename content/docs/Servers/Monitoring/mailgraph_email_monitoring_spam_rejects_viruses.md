---
weight: 999
url: "/Mailgraph\_\\:\_Surveillance\_des\_mails\_(Spams,\_rejects,\_virus...)/"
title: "Mailgraph: Email Monitoring (Spam, Rejects, Viruses...)"
description: "A guide on installing and configuring Mailgraph to monitor email traffic, spam, viruses and rejected emails on your mail server."
categories: ["Linux", "Debian", "Nginx"]
date: "2010-02-03T06:33:00+02:00"
lastmod: "2010-02-03T06:33:00+02:00"
tags:
  [
    "email",
    "monitoring",
    "mailgraph",
    "spam",
    "postfix",
    "lighttpd",
    "visualization",
  ]
toc: true
---

## Introduction

Mailgraph is a software tool that generates graphs of email statistics such as spam, viruses, and more.

It provides a good overview of what's happening on your mail server.

[Mailgraph Official Site](https://people.ee.ethz.ch/~dws/software/mailgraph/)

## Installation

The installation is really simple:

```bash
apt-get install mailgraph
```

Now you just need to access the interface. On Apache:

```
http://mysite/cgi-bin/mailgraph.cgi
```

And on Lighttpd:

```
http://mysite/mailgraph
```

## Mailgraph without CGI

Whether for performance, security, or simplicity reasons, it's quite common not to have a CGI module on a server (installing CGI with nginx is tedious for example). However, the mailgraph stats tool is only designed to run in CGI. Here is a small script that allows you to generate mailgraph graphs without CGI:

```bash
#!/bin/sh
MAILGRAPH_PATH=/usr/lib/cgi-bin/mailgraph.cgi # Debian
#MAILGRAPH_PATH=/usr/local/www/cgi-bin/mailgraph.cgi # FreeBSD
#MAILGRAPH_PATH=/usr/local/lib/mailgraph/mailgraph.cgi # OpenBSD

MAILGRAPH_DIR=/var/www/mailgraph

umask 022

mkdir -p $MAILGRAPH_DIR

$MAILGRAPH_PATH | sed '1,2d ; s/mailgraph.cgi?//' > $MAILGRAPH_DIR/index.html

for i in 0-n 0-e 1-n 1-e 2-n 2-e 3-n 3-e; do
        QUERY_STRING=$i $MAILGRAPH_PATH | sed '1,3d' > $MAILGRAPH_DIR/$i
done
```

This script can be added to crontab, which allows regular saving of the generated graphs. Tested on Debian, FreeBSD and OpenBSD (MAILGRAPH_PATH variable should be adapted).

## FAQ

### No graphs displaying under Lighttpd

I encountered this annoying bug where no graphs are displayed. To work around this issue (not in an elegant way), you must edit the file `/etc/lighttpd/conf-enabled/50-mailgraph.conf` and modify the 2nd line:

```perl
# Alias for phpMyAdmin directory
alias.url += (
    "/mailgraph.cgi" => "/usr/lib/cgi-bin/mailgraph.cgi",
)

$HTTP["url"] =~ "^/mailgraph*", {

}
```

Here "/mailgraph" has been replaced with "/mailgraph.cgi". Reload your lighttpd configuration and it should work.

## Resources
- [Mailgraph documentation](/pdf/mailgraph.pdf)
