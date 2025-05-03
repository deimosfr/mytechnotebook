---
weight: 999
url: "/Monitorer_Windows_avec_Munin/"
title: "Monitoring Windows with Munin"
description: "How to set up monitoring for Windows systems using Munin through SNMP"
categories: ["Debian", "Linux", "Windows", "Monitoring"]
date: "2007-06-27T07:41:00+02:00"
lastmod: "2007-06-27T07:41:00+02:00"
tags: ["Munin", "Windows", "Monitoring", "SNMP", "Servers", "Network"]
toc: true
---

## Introduction

Windows and its complications... it's so much easier on Linux! Anyway, there's no getting around it, it's still good when Munin runs on Windows, so here's the documentation...

## Configuring SNMP on Windows

SNMP is a Windows service that can be installed as follows:

- Control Panel
- Add/Remove Programs
- Add or Remove Windows Components
- Management and Analysis Tools
- SNMP

(You'll need the Windows CD-ROM)

Then, configure the SNMP service:

- Control Panel
- Administrative Tools
- Services
- SNMP Service

### Agent Tab

Check all the relevant boxes

### Traps Tab

- Community name: public
- Trap destinations: the server(s) that collect(s) SNMP information

### Security Tab

- Modify the Public community by setting rights to READ CREATE
- Check "Accept SNMP packets from any host"
- OK
- Restart the service

## Installing Munin on Linux Debian

Nothing complicated:

```bash
apt-get install munin libwww-perl libnet-snmp-perl
```

Remember to open ports 4949 (tcp) as well as 161 (tcp) and 162 (tcp) for SNMP support.
Do a small search for the SNMP capabilities of the machine in question:

```bash
munin-node-configure-snmp windows.mydomain
```

Which should result in something like this:

```bash
ln -s /usr/share/munin/plugins/snmp__df /etc/munin/plugins/snmp_windows.mydomain_df
ln -s /usr/share/munin/plugins/snmp__if_err_ /etc/munin/plugins/snmp_windows.mydomain_if_err_16777219
ln -s /usr/share/munin/plugins/snmp__if_ /etc/munin/plugins/snmp_windows.mydomain_if_16777219
ln -s /usr/share/munin/plugins/snmp__processes /etc/munin/plugins/snmp_windows.mydomain_processes
ln -s /usr/share/munin/plugins/snmp__users /etc/munin/plugins/snmp_windows.mydomain_users
```

Copy and paste these commands.

## Configuring Munin

Add the following lines in `/etc/munin/munin/munin.conf`:

```bash
[windows.mydomain]
  address 127.0.0.1
  use_node_name no
```

Note that the name in brackets must correspond to the name used for the links created previously.
127.0.0.1 is not an error because it's the local SNMP server that manages data from remote machines.

Then restart Munin:

```bash
/etc/init.d/munin-node restart
```

## FAQ

### It's not graphing!!!

In these cases... only one solution (don't forget we're on Windows): reboot the Windows machine.

## References

[https://munin.projects.linpro.no/wiki/HowToMonitorWindows](https://munin.projects.linpro.no/wiki/HowToMonitorWindows)  
[https://www.debian-administration.org/articles/380](https://www.debian-administration.org/articles/380)  
[https://www.skolelinux.no/~klaus/sarge/x3579.html](https://www.skolelinux.no/~klaus/sarge/x3579.html)
