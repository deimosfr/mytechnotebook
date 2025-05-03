---
weight: 999
url: "/Configuration_d'un_réseau_local_sur_Debian_et_Ubuntu/"
title: "Configuring a Local Network on Debian and Ubuntu"
description: "A simple guide to configuring static and dynamic network interfaces on Debian and Ubuntu systems"
categories: ["Linux", "Network", "Debian", "Ubuntu"]
date: "2007-04-18T14:03:00+02:00"
lastmod: "2007-04-18T14:03:00+02:00"
tags: ["network", "interfaces", "ethernet", "IP", "DNS", "DHCP"]
toc: true
---

## Introduction

Not always easy to remember what to put in these damn files, right? Here are some simple examples.

## Configuring Network Interfaces

Edit the file `/etc/network/interfaces`.

### Static IP Addressing

For example, here's how I configure my eth0 interface:

```bash
auto lo eth0
iface lo inet loopback

# Interface eth0
iface eth0 inet static
       address 192.168.0.1
       netmask 255.255.255.0
       broadcast 192.168.0.255
       network 192.168.0.0
       gateway 192.168.0.138
```

### Dynamic IP Addressing

If I want dynamic addressing, it's even simpler:

```bash
auto lo eth0
iface lo inet loopback

# Interface eth0 
iface eth0 inet dhcp
```

## Configuring DNS Servers

Edit the file `/etc/resolv.conf`:

```bash
search deimos.fr local
nameserver 192.168.0.1
nameserver 192.168.0.2
nameserver 212.27.32.5
nameserver 212.27.32.6
```
