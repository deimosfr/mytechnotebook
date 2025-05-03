---
weight: 999
url: "/Systemd\\:_how_to_debug_on_boot_fail/"
title: "Systemd: How to Debug on Boot Failure"
description: "Guide on how to debug systemd boot failures including checking failed services, getting service status information and troubleshooting techniques."
categories: ["Linux", "Debian"]
date: "2014-07-27T12:16:00+02:00"
lastmod: "2014-07-27T12:16:00+02:00"
tags: ["systemd", "boot", "debugging", "Debian", "Linux"]
toc: true
---

![Systemd](/images/poweredbylinux.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 208-6 |
| **Operating System** | Debian 8 |
| **Last Update** | 27/07/2014 |
{{< /table >}}

## Introduction

systemd[^1] is a system and service manager for Linux, compatible with SysV and LSB init scripts. Systemd provides aggressive parallelization capabilities, uses socket and D-Bus activation for starting services, offers on-demand starting of daemons, keeps track of processes using Linux control groups, supports snapshotting and restoring of the system state, maintains mount and automount points and implements an elaborate transactional dependency-based service control logic.

When I decided to migrate to Systemd on Debian, it unfortunately worked at the first time. That's why I need to deep dive into Systemd issues and understand why it wasn't working.

## Usage

First of all, you have to search for failed services:

```bash
> systemctl --state=failed
```

To get more information on a service:

```bash
systemctl status <servicename>
```

Try to know more on the specific PID (may not work in some cases):

```bash
> journalctl -b _PID=
```

Generally, it is because of modules that don't load properly or shouldn't load. Make your changes, try to start the problematic services:

```bash
> systemctl start <servicename>
```

Now things should be ok

## References

[^1]: https://wiki.archlinux.org/index.php/Systemd
