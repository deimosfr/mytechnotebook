---
weight: 999
url: "/Bindgraph_\\:_Avoir_des_stats_et_des_graphs_des_requêtes_DNS/"
title: "Bindgraph: Get Statistics and Graphs of DNS Queries"
description: "How to install and configure Bindgraph to visualize DNS query statistics in graphical form"
categories: ["DNS", "Monitoring", "Statistics"]
date: "2010-02-02T21:16:00+02:00"
lastmod: "2010-02-02T21:16:00+02:00"
tags: ["Bind", "DNS", "Monitoring", "Graphs", "Statistics", "Lighttpd"]
toc: true
---

## Introduction

This is the kind of software I really like. It installs quickly and is very practical.

## Installation

```bash
aptitude install bindgraph
```

## Configuration

### Lighttpd

Create a configuration file for bindgraph in the /etc/lighttpd/conf-available/ directory:

```bash {linenos=table}
# Alias for phpMyAdmin directory
alias.url += (
    "/bindgraph" => "/usr/lib/cgi-bin/bindgraph.cgi",
)

$HTTP["url"] =~ "^/bindgraph*", {

}
```

Now, create the necessary symlink:

```bash
cd /etc/lighttpd/conf-enabled
ln -s /etc/lighttpd/conf-available/50-bindgraph.conf .
```

All that's left is to restart the web server and it's accessible via:

http://myserver/bindgraph

## FAQ

### I have a blank page, why?

In the Lighttpd error logs (/var/log/lighttpd/error.log), you might see something like this:

```
"-T" is on the #! line, it must also be used on the command line at /usr/lib/cgi-bin/bindgraph.cgi line 1.
```

Edit the file /usr/lib/cgi-bin/bindgraph.cgi and remove the letter 'T' from the first line so it looks like:

```bash
#!/usr/bin/perl -w
```

That's it :-)
