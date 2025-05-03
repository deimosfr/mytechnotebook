---
weight: 999
url: "/IPv6_\\:_Suppression_complète,_IPv4_seulement/"
title: "IPv6: Complete Removal, IPv4 Only on Debian"
description: "How to completely disable IPv6 support and run with IPv4 only on a Debian system."
categories: ["Linux", "Debian"]
date: "2007-08-21T21:10:00+02:00"
lastmod: "2007-08-21T21:10:00+02:00"
tags: ["IPv6", "IPv4", "Networking", "System Configuration"]
toc: true
---

## Introduction

When installing a new Debian 4.0 distribution, IPv6 support is enabled by default. This can cause some problems or even simply slow down your system. Indeed, all applications will use IPv6 support for name resolution before or after trying with IPv4.

## Method 1

This is the case for mplayer used to read an audio or video stream:

```bash
[...]
Resolving live.radio-gresivaudan.org for AF_INET6...
Couldn't resolve name for AF_INET6: live.radio-gresivaudan.org
Resolving live.radio-gresivaudan.org for AF_INET...
Connecting to server live.radio-gresivaudan.org[217.117.157.190]: 8000...
[...]
```

There is a simple method to disable IPv6 support. You just need to prevent the system from loading the corresponding module. Edit the file `/etc/modprobe.d/blacklist` and add the line:

```
blacklist ipv6
```

Then edit your `/etc/hosts` file to remove IPv6 entries.

## Method 2

Another solution (less elegant in my opinion) is to modify `/etc/modprobe.d/aliases` from:

```
alias net-pf-10 ipv6
```

to:

```
alias net-pf-10 off
alias ipv6 off
```

You just need to restart the system (this is not strictly necessary, it is possible to do otherwise, but it is much simpler, especially after an installation). A quick ifconfig will confirm that IPv6 is no longer managed (no more inet6 addr:).
