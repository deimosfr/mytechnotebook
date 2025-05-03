---
weight: 999
url: "/OpenNTPd_\\:_Mise_en_place_d'un_serveur_NTP_sous_OpenBSD/"
title: "OpenNTPd: Setting up an NTP server on OpenBSD"
description: "A guide for installing and configuring OpenNTPd on OpenBSD"
categories: ["Linux"]
date: "2009-05-21T11:40:00+02:00"
lastmod: "2009-05-21T11:40:00+02:00"
tags: ["Servers", "Network", "OpenBSD", "NTP"]
toc: true
---

## Introduction

Having your machines synchronized to the correct time is very practical! Especially when you're trying to read logs from two different machines.

For those who need microsecond-level precision, OpenNTPd is not the right solution - use the traditional NTP server instead.

## Installation

Nothing to do :-)

## Configuration

Open the `/etc/ntpd.conf` file and configure it as follows:

```bash
# $OpenBSD: ntpd.conf,v 1.8 2007/07/13 09:05:52 henning Exp $
# sample ntpd configuration file, see ntpd.conf(5)

# Addresses to listen on (ntpd does not listen by default)
# listen on *

# sync to a single server
server 0.pool.ntp.org
server 1.pool.ntp.org
server 2.pool.ntp.org
server ntp1.jussieu.fr

# use a random selection of 8 public stratum 2 servers
# see http://support.ntp.org/bin/view/Servers/NTPPoolServers
servers pool.ntp.org
```

This is a basic configuration, but it works well.

Next, edit the `/etc/rc.conf.local` file to indicate that the NTP server should start at boot time, and modify the line as follows:

```bash
ntpd_flags="-s"
```

Now, if you want to start the daemon to see how it works:

```bash
ntpd
```

## Verification

To ensure that everything is working, use this command:

```bash
tail -f /var/log/daemon
```

It can take up to about 4 minutes to synchronize your machine.

## Manual Synchronization

If you want to perform a manual synchronization, use this command:

```bash
rdate -ncv 0.pool.ntp.org
```

## Resources
- http://www.openntpd.org/
- http://www.openbsd101.com/tipstricks.html
