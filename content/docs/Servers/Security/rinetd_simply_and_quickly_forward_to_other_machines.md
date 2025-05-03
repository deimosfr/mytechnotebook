---
weight: 999
url: "/Rinetd_\\:_Forwarder_simplement_et_rapidement_vers_d'autres_machines/"
title: "Rinetd: Simply and Quickly Forward to Other Machines"
description: "A guide on how to use Rinetd to set up TCP forwarding between machines without complex firewall rules."
categories: ["Linux", "Network", "Debian"]
date: "2008-09-25T12:07:00+02:00"
lastmod: "2008-09-25T12:07:00+02:00"
tags: ["Debian", "Network", "Servers", "Linux"]
toc: true
---

## Introduction

In the past we've examined the use of firewall rules for forwarding incoming connections from one machine to another. But there is a simpler approach using the rinetd package. Read on to learn about this tool.

The rinetd package contains a simple tool which may be configured to listen for connections upon a machine, and silently redirect them to a new destination. In short it acts as a simple to configure TCP proxy.

## Installation

You may install this package via:

```bash
apt-get update
apt-get install rinetd
```

(If you prefer you may use "aptitude update; aptitude install rinetd" - old habits die hard with me!)

## Configuration

Once installed you'll find a configuration file located at `/etc/rinetd.conf`. This file is used to tell the deamon which ports it should listen for connections upon, and what it should do when they arrive.

By default no ports are configured for forwarding, and so the file will consist entirely of comments. A default configuration file would look something like this, to give you an idea of the configuration:

```bash
#
# forwarding rules come here
#
# you may specify allow and deny rules after a specific forwarding rule
# to apply to only that forwarding rule
#
# bindadress    bindport  connectaddress  connectport


# logging information
logfile /var/log/rinetd.log

# uncomment the following line if you want web-server style logfile format
# logcommon
```

Note: There are more details about allowed options in the manpage which you may view by running "man rinetd".

To demonstrate how the forwarding is configured and used we'll make a simple example. Assume that you have a machine with the IP address 1.2.3.4 which has been running Apache, and that you'd like to move that to the IP address 4.3.2.1.

You've already updated DNS to point visitors to the new IP address, but you want to ensure that people connecting to the old IP still continue to receive service.

To handle this case you should update the `/etc/rinetd.conf` file to read:

```bash
# bindadress    bindport  connectaddress  connectport
1.2.3.4         80        4.3.2.1         80
1.2.3.4         443       4.3.2.1         443
```

Once you restart rinetd all incoming connections on port 80 and 443 will be seamlessly redirected from the old IP to the new one - although you will need to restart rinetd after making the change to your configuration file:

```bash
$ /etc/init.d/rinetd restart
Stopping internet redirection server: rinetd.
Starting internet redirection server: rinetd.
```

rinetd is a very small, stable, and simple program, and you might find it simpler to understand than the matching generic iptables TCP proxy solution.

The only downside to using rinetd is that there is no support for UDP connections, and no support for redirecting FTP access - because of the complex nature of FTP.

## Resources
- http://www.debian-administration.org/articles/601
- [Port-Forwarding With rinetd](/pdf/port-forwarding_with_rinetd_on_debian_etch.pdf)
