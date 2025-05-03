---
weight: 999
url: "/Kernel_\\:_Compilation_des_modules/"
title: "Kernel: Module Compilation"
description: "Guide on compiling Linux kernel modules with focus on Iptables firewall requirements."
categories: ["Linux"]
date: "2006-08-27T22:24:00+02:00"
lastmod: "2006-08-27T22:24:00+02:00"
tags: ["Linux", "Kernel", "Iptables", "Security", "System Administration"]
toc: true
---

Iptables is nowadays the Linux firewall of choice. However, when you're a beginner, it's not always easy to know what each module corresponds to.

## Minimum Requirements

What do you need to recompile at minimum for the kernel?

```bash
CONFIG_PACKET - Direct communication with network interfaces
CONFIG_NETFILTER - Kernel management, necessary for Netfilter
CONFIG_IP_NF_CONNTRACK - Necessary for NAT and Masquerade
CONFIG_IP_NF_NETFILTER - Adds NETFILTER table
CONFIG_IP_NF_IPTABLES - Required for iptables user space utility
CONFIG_IP_NF_MANGLE - Adds MANGLE table
CONFIG_IP_NF_NAT - Adds NAT table
```

**Rule not to add:**

```bash
CONFIG_NET_FASTROUTE - Fast routing bypasses NETFILTER entry points
```

## Legacy Firewall Compatibility

Here are the modules that will provide compatibility with previous firewalls:

```bash
CONFIG_IP_NF_COMPAT_IPCHAINS
CONFIG_IP_NF_COMPAT_IPFWADM
```

## Service-Specific Modules

This is a list of modules needed according to the services you want to use:

```bash
IP_CONNTRACK_AMANDA - Amanda is a backup software
IP_CONNTRACK_FTP - FTP is used for file transfers
IP_CONNTRACK_IRC - IRC (Internet Relay Chat)
IP_CONNTRACK_TFTP - Trivial FTP
```
