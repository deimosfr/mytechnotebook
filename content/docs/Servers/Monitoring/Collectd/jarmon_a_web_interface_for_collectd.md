---
weight: 999
url: "/Jarmon_\\:_Une_interface_web_pour_Collectd/"
title: "Jarmon: A Web Interface for Collectd"
description: "This article explains how to install and configure Jarmon, a web interface for Collectd that allows for clear visualization and zooming of monitoring data."
categories:
  - "Linux"
  - "Servers"
date: "2010-07-27T14:12:00+02:00"
lastmod: "2010-07-27T14:12:00+02:00"
tags:
  - "Collectd"
  - "Monitoring"
  - "Web interface"
  - "Apache"
toc: true
---

## Introduction

Jarmon is another interface for Collectd. I like this one because it's clean and allows zooming. However, at the moment it only works on a host-by-host basis. This means that in the configuration, you need to specify a particular host.

## Installation

For the installation, we need a web server:

```bash
aptitude install apache2 bzr
```

## Configuration

Let's retrieve the source code:

```bash
mkdir -p /var/www/
cd /var/www
bzr branch lp:~richardw/jarmon/trunk
mv trunk jarmon
chown -Rf www-data. jarmon
cd jarmon
```

Now we create a symbolic link to the directory containing the RRD files of the machine we want to monitor:

```bash
ln -s /var/lib/collectd/rrd/localhost data
```

Now you can access the web interface: http://collectd-server/jarmon
