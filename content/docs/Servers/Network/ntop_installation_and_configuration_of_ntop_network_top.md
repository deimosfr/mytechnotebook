---
weight: 999
url: "/NTOP\\:_Installation_et_configuration_d'NTOP_(Network_TOP)/"
title: "NTOP: Installation and Configuration of NTOP (Network TOP)"
description: "A guide for installing and configuring NTOP, a tool for collecting network information with simple configurations and web interface for viewing network statistics."
categories: ["Linux", "CentOS", "Network"]
date: "2009-02-05T06:07:00+02:00"
lastmod: "2009-02-05T06:07:00+02:00"
tags: ["Networking", "Monitoring", "Network", "Servers"]
toc: true
---

## Introduction

This is a tool to collect network information with simple configurations. Users may use a web browser to access the current network views that include charts and statistics.

## Installation

ntop is available on most Linux distributions. Below are the steps tested on Ubuntu 8.10 Linux. I am sure that ntop is available for RPM based distro like CentOS.

```bash
apt-get install ntop
/etc/init.d/ntop -i eth0 start
```

Replace eth0 with the network interface name. The system will maintain a database of the network information in the `/var/lib/ntop` directory. Create the default admin user with:

```bash
ntop
```

* Enter the admin user password.
* View ntop

The ntop results can be viewed with a web browser pointing to any of the URL address:

* http://localhost:3000
* https://localhost:3001 (this may need to be configured)

The ntop screen should have the following main menus:

* About: What is ntop, credits, documentations and configurations.
* Summary: Traffic, Hosts, Network load, network flows
* All protocols: Traffic, Throughput, Activity
* IP: Summary, Traffic directions, Local
* Utils: Data dump, View log
* Plugins: Lots of plugins to enable or configure
* Admin: Configure, Shutdown

### Startup settings

To change the ntop server startup settings, Select Admin -> Configure -> Startup options in the web interface.

Edit the changes as needed and save.

## References

http://www.howtoforge.com/installing-and-configuring-ntop
