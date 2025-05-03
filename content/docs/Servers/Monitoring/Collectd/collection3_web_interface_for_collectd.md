---
weight: 999
url: "/Collection3_\\:_Une_interface_web_pour_Collectd/"
title: "Collection3: A Web Interface for Collectd"
description: "How to install and configure Collection3, a web interface for Collectd that enables viewing gathered statistics through a browser"
categories: ["Servers", "Monitoring", "Web"]
date: "2010-06-09T21:58:00+02:00"
lastmod: "2010-06-09T21:58:00+02:00"
tags: ["collectd", "collection3", "monitoring", "web interface", "apache"]
toc: true
---

## Introduction

[Collectd]({{< ref "docs/Servers/Monitoring/Collectd/collectd_installation_and_configuration.md">}}) is great, but with a functional web interface, it's even better. We're going to see how to install [Collection3](https://collectd.org/wiki/index.php/Collection3), an interface that isn't necessarily pretty, but is very functional.

## Installation

You must have a web server like Apache and CGI enabled. Here's what you need to install if you choose Apache:

```bash
aptitude install apache2 librrds-perl libconfig-general-perl libhtml-parser-perl libregexp-common-perl
```

I haven't mentioned it, but it's obvious that you need [Collectd]({{< ref "docs/Servers/Monitoring/Collectd/collectd_installation_and_configuration.md">}}) installed.

## Configuration

We'll configure the Apache2 part:

```apache
ScriptAlias /collectd/bin/ /usr/share/doc/collectd/examples/collection3/bin/
Alias /collectd/ /usr/share/doc/collectd/examples/collection3/

<Directory /usr/share/doc/collectd/examples/collection3/>
    AddHandler cgi-script .cgi
    DirectoryIndex bin/index.cgi
    Options +ExecCGI
    Order Allow,Deny
    Allow from all
</Directory>
```

Then restart Apache2 or reload it.

You can now access your data via the following address: http://collectd-server/collectd/

## Resources

- http://collectd.org/wiki/index.php/Collection3
- http://collectd.org/wiki/index.php/List_of_front-ends
