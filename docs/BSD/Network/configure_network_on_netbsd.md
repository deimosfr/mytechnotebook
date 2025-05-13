---
title: "Configure Network on NetBSD"
slug: configure-network-on-netbsd/
description: "How to configure network interfaces, static IPs, and gateways on NetBSD"
categories: ["BSD", "Network", "System Administration"]
date: "2013-03-15T22:32:00+01:00"
lastmod: "2013-03-15T22:32:00+01:00"
tags: ["NetBSD", "network", "configuration", "interfaces"]
---

![NetBSD](../../static/images/netbsd-logo.avif)


|||
|-|-|
| **Operating System** | 6.0.1 |
| **Website** | [NetBSD Website](https://www.netbsd.org/) |


## Introduction

The network is an essential part of system configuration. I will cover some aspects of it here.

## Configuration

### Display Interfaces and Associated IPs

Here's the command:

```bash
ifconfig
```

### Static IPs

We can declare interfaces persistently:

```bash
[...]
# Add local overrides below
#
ifconfig_vr0="inet 192.168.0.254 netmask 255.255.255.0" # Wan
ifconfig_vr1="inet 192.168.1.254 netmask 255.255.255.0" # DMZ
[...]
```

You add an ifconfig entry with the interface name, followed by its parameters.

### Gateway

It's possible to set the gateway like this:

```bash
[...]
defaultroute="192.168.0.138"
[...]
```

Or alternatively:

```bash
192.168.0.138
```

### Restart Network Services

```bash
/etc/rc.d/network restart
```
