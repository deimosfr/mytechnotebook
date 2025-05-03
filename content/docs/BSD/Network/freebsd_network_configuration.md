---
weight: 999
url: "/Configurer_le_réseau_sous_FreeBSD/"
title: "Configuring Network on FreeBSD"
description: "A guide to network configuration on FreeBSD, including interface configuration, static IP settings, DHCP, routing, and more"
categories: ["FreeBSD", "Networking", "System Administration"]
date: "2012-07-02T10:07:00+02:00"
lastmod: "2012-07-02T10:07:00+02:00"
tags: ["FreeBSD", "Networking", "System Configuration", "DHCP", "Routing"]
toc: true
---

![FreeBSD](/images/poweredbyfreebsd.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Operating System** | FreeBSD 9 |
| **Website** | [FreeBSD Website](https://www.freebsd.org) |
| **Last Update** | 02/07/2012 |
{{< /table >}}

## Introduction

The network is an essential part of system configuration, so I'll cover some aspects of it here.

## Configuration

### Display Interfaces and Associated IPs

The command is always the same:

```bash
ifconfig
```

### Declare Interfaces

We can declare the interfaces to manage at startup by simply listing the interfaces separated by spaces:

```bash
# Network
network_interfaces="lo0 vr0 vr1 vr2"
ifconfig_lo0="inet 127.0.0.1"
```

Here I've declared 4 interfaces and configured lo0.

### DHCP

If you want to set an interface to use DHCP, it's very simple:

```bash
# Network
ifconfig_vr0="DHCP"
```

Here my vr0 interface is configured with DHCP.

### Static IPs

If you want to set a static IP address to an interface, it's very simple:

```bash
# Network
ifconfig_vr0="inet 192.168.10.254 netmask 255.255.255.0"
```

Here my vr0 interface is configured with a static IP.

### Default Gateway

To configure the default gateway:

```bash
# Network
defaultrouter="192.168.10.138"
```

### Display Routes

To display routes:

```bash
netstat -rn
```

### Add a Route

To add a route, simply define one or more route names and define them line by line:

```bash
static_routes="route1 route2"
route_route1="-net 222.2.90.0/24 222.2.30.1"
route_route2="-net 222.2.100.0/24 222.2.30.1"
```

### Restart Network Services

To restart network services:

```bash
/etc/rc.d/netif restart
```

And for routing services:

```bash
/etc/rc.d/routing restart
```

## References

http://www.freebsd.org/doc/fr/articles/ppp/chap3.html  
http://www.cyberciti.biz/faq/freebsd-setup-default-routing-with-route-command/
