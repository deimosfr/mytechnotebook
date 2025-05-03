---
weight: 999
url: "/IP_\\:_La_commande_de_gestion_de_sa_carte_réseau/"
title: "IP: Network Interface Management Command"
description: "How to use the ip command to manage network interfaces, which is gradually replacing ifconfig due to its enhanced functionality."
categories: ["Linux", "RHEL"]
date: "2012-05-03T14:47:00+02:00"
lastmod: "2012-05-03T14:47:00+02:00"
tags: ["Networking", "Linux", "Command Line", "System Administration"]
toc: true
---

## Introduction

The "ip" command is gradually replacing the ifconfig command due to its enhanced functionality. With the arrival of RHEL6, this transition is becoming more evident. In this article, we'll see how to use this command effectively.

## Usage

- View the status of all interfaces:

```bash
> ip link show
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 16436 qdisc noqueue state UNKNOWN 
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP qlen 1000
    link/ether 08:00:27:9a:b6:35 brd ff:ff:ff:ff:ff:ff
```

- View the status of a specific interface (eth0 in this case):

```bash
> ip link show eth0
2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP qlen 1000
    link/ether 08:00:27:9a:b6:35 brd ff:ff:ff:ff:ff:ff
```

- View the status of a specific interface (eth0) with IPv4 information only (-4):

```bash
> ip -4 addr show eth0
2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP qlen 1000
    inet 10.0.1.9/24 brd 10.0.1.255 scope global eth0
```

- View detailed statistics for interface eth0:

```bash
> ip -s link show eth0
2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP qlen 1000
    link/ether 08:00:27:9a:b6:35 brd ff:ff:ff:ff:ff:ff
    RX: bytes  packets  errors  dropped overrun mcast   
    28627      284      0       0       0       0      
    TX: bytes  packets  errors  dropped carrier collsns 
    40732      128      0       0       0       0
```

- Create a virtual IP address (VIP):

```bash
ip addr add 192.168.0.1/24 dev eth0 label eth0:0
```

- Remove a virtual IP address (VIP):

```bash
ip addr del 192.168.0.1/24 dev eth0
```

- Add a VLAN and VIP on the VLAN (VLAN 90 on bond0):

```bash
ip link add link bond0 name bond0:90 type vlan id 90
ip link set dev bond0:90 up
ip addr add 192.168.0.1/24 dev bond0:90
```
