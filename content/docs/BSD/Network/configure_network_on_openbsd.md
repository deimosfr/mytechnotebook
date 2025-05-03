---
weight: 999
url: "/Configurer_le_réseau_sous_OpenBSD/"
title: "Configure Network on OpenBSD"
description: "A guide to basic network configuration on OpenBSD systems including interface setup, static IPs, gateways, DNS, and performance tuning"
categories: ["BSD", "Network", "System Administration"]
date: "2008-03-17T07:22:00+01:00"
lastmod: "2008-03-17T07:22:00+01:00"
tags: ["OpenBSD", "network", "configuration", "interfaces", "routing"]
toc: true
---

## Introduction

I won't go into too much detail, as that's what the OpenBSD website is for. This is just a quick reference for recalling how to rapidly configure networking.

## Basic Configuration

### ifconfig

To list available network cards and view defined IP addresses:

```bash
$ ifconfig
lo0: flags=8049<UP,LOOPBACK,RUNNING,MULTICAST> mtu 33224
        inet 127.0.0.1 netmask 0xff000000
        inet6 ::1 prefixlen 128
        inet6 fe80::1%lo0 prefixlen 64 scopeid 0x5
lo1: flags=8008<LOOPBACK,MULTICAST> mtu 33224
fxp0: flags=8843<UP,BROADCAST,RUNNING,SIMPLEX,MULTICAST> mtu 1500
        address: 00:04:ac:dd:39:6a
        media: Ethernet autoselect (100baseTX full-duplex)
        status: active
        inet 10.0.0.38 netmask 0xffffff00 broadcast 10.0.0.255
        inet6 fe80::204:acff:fedd:396a%fxp0 prefixlen 64 scopeid 0x1
pflog0: flags=0<> mtu 33224
pfsync0: flags=0<> mtu 2020
sl0: flags=c010<POINTOPOINT,LINK2,MULTICAST> mtu 296
sl1: flags=c010<POINTOPOINT,LINK2,MULTICAST> mtu 296
ppp0: flags=8010<POINTOPOINT,MULTICAST> mtu 1500
ppp1: flags=8010<POINTOPOINT,MULTICAST> mtu 1500
tun0: flags=10<POINTOPOINT> mtu 3000
tun1: flags=10<POINTOPOINT> mtu 3000
enc0: flags=0<> mtu 1536
bridge0: flags=0<> mtu 1500
bridge1: flags=0<> mtu 1500
vlan0: flags=0<> mtu 1500
        address: 00:00:00:00:00:00
vlan1: flags=0<> mtu 1500
        address: 00:00:00:00:00:00
gre0: flags=9010<POINTOPOINT,LINK0,MULTICAST> mtu 1450
carp0: flags=0<> mtu 1500
carp1: flags=0<> mtu 1500
gif0: flags=8010<POINTOPOINT,MULTICAST> mtu 1280
gif1: flags=8010<POINTOPOINT,MULTICAST> mtu 1280
gif2: flags=8010<POINTOPOINT,MULTICAST> mtu 1280
gif3: flags=8010<POINTOPOINT,MULTICAST> mtu 1280
```

Here's a listing of available interfaces with their explanations:

- lo - Loopback Interface
- pflog - Packet Filter Logging Interface
- sl - SLIP Network Interface
- ppp - Point to Point Protocol
- tun - Tunnel Network Interface
- enc - Encapsulating Interface
- bridge - Ethernet Bridge Interface
- vlan - IEEE 802.1Q Encapsulation Interface
- gre - GRE/MobileIP Encapsulation Interface
- gif - Generic IPv4/IPv6 Tunnel Interface
- carp - Common Address Redundancy Protocol Interface

### Setting a Permanent IP

Simply create a file `/etc/hostname.interface_name` and insert the IP and subnet mask:

```bash
echo "inet 172.16.1.100 255.255.255.0 NONE" > /etc/hostname.fxp0
```

### Default Gateway

To set your default gateway (or default route), create the file `/etc/mygate` and specify the desired IP:

```bash
echo "172.16.1.254" > /etc/mygate
```

### DNS

To set DNS servers, it's quite simple, as in all Unix-like systems, edit the resolv.conf file:

```bash
search example.com
nameserver 125.2.3.4
nameserver 125.2.3.5
```

### Hostname

To define a machine's hostname, simply insert the desired name in the `/etc/myname` file:

```bash
echo "server" > /etc/myname
```

### Applying Changes Without Rebooting

Obviously, you don't need to restart the machine, but you do need to restart the network services! To restart network services:

```bash
$ sh /etc/netstart
writing to routing socket: File exists
add net 127: gateway 127.0.0.1: File exists
writing to routing socket: File exists
add net 224.0.0.0: gateway 127.0.0.1: File exists
```

### Routes

To display routes:

```bash
route show
```

or

```bash
netstat -rn
```

Finally, to add static routes, here's an example of an `/etc/hostname.if` file (adapt to your needs):

```bash
inet 192.168.254.254 255.255.255.0 NONE !
route add 10.0.0.0 192.168.254.1
```

## Network Tuning

### Improving Performance on High Bandwidth

If you want to increase bandwidth performance, you can activate it via sysctl. To test temporarily:

```bash
sysctl net.inet.tcp.recvspace=65536
sysctl net.inet.tcp.sendspace=65536
```

To make it permanent, add these lines to `/etc/sysctl.conf`:

```bash
echo "# Increase performance" >> /etc/sysctl.conf
echo "net.inet.tcp.recvspace=65536" >> /etc/sysctl.conf
echo "net.inet.tcp.sendspace=65536" >> /etc/sysctl.conf
```

## Resources

- [OpenBSD Networking FAQ](https://www.openbsd.org/faq/faq6.html)
- [OpenBSD hostname man](https://www.openbsd.org/cgi-bin/man.cgi?query=hostname.if&sektion=5)
