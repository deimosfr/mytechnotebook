---
weight: 999
url: "/Configurer_le_réseau_sur_Solaris/"
title: "Configure Network on Solaris"
description: "A comprehensive guide to configuring network interfaces, routes, DNS, and VIPs on Solaris systems"
categories: ["Solaris", "Network", "System Administration"]
date: "2013-05-06T14:01:00+02:00"
lastmod: "2013-05-06T14:01:00+02:00"
tags: ["Solaris", "network", "configuration", "interfaces", "routing", "VIP"]
toc: true
---

## Introduction

Solaris can be challenging! Especially when you come from a Linux/BSD world and find that all the network commands are strange.

## Detecting Network Cards

To detect network cards that are present and especially connected, here's a very useful command:

```bash
$ dladm show-dev
igb0            link: up        speed: 1000  Mbps       duplex: full
igb1            link: up        speed: 1000  Mbps       duplex: full
igb2            link: up        speed: 1000  Mbps       duplex: full
igb3            link: up        speed: 1000  Mbps       duplex: full
clprivnet0              link: unknown   speed: 0     Mbps       duplex: unknown
```

Or you can use `dladm show-links` to get more detailed information about the connection status.

For a link aggregation (IPMP) of type 802.3ad:

```bash
> dladm show-aggr -s -i 2 1

key:1		ipackets   rbytes       opackets   obytes       %ipkts  %opkts
 	Total	355021     531533375    60288      4944021     
	nxge0	166090     249992028    0          0              46.8     0.0  
	nxge1	120638     179830318    0          0              34.0     0.0  
	nxge4	16         1172         25728      2109696         0.0    42.7  
	nxge5	68277      101709857    34560      2834325        19.2    57.3  

key:1		ipackets   rbytes       opackets   obytes       %ipkts  %opkts
 	Total	344131     513180425    47543      3900596     
	nxge0	167398     250160702    12         1672           48.6     0.0  
	nxge1	95286      142041090    8          1330           27.7     0.0  
	nxge4	17         1320         21601      1771571         0.0    45.4  
	nxge5	81430      120977313    25922      2126023        23.7    54.5
```

## Basic Network Configuration

Here are some useful commands to reset a Solaris configuration, especially the network:

```bash
ifconfig -a 	Display interfaces with IP and MAC addresses
show-devs 	Display peripherals
prtconf -vD    Display all peripherals
```

There are two ways to reconfigure the network. The wizard and manual:

* The wizard:

```bash
sys-unconfig 	This will reset ALL network configuration! Reboot required
sysidconfig 	Not tested, but also supposed to reconfigure the network
```

* And then manually, here are the files to modify:

```
/etc/hostname.x (x corresponds to the network interface)
/etc/nodename
/etc/defaultrouter
/etc/netmasks
```

### Dynamic Configuration

Finally, a "reboot" or "boot net" should do the trick. However, you may not be able to reboot, so here's the solution:

```bash
ifconfig e1000g0 plumb
ifconfig e1000g0 192.168.0.1 netmask 255.255.255.0
ifconfig e1000g0 up
ifconfig -a
```

### Persistent Configuration

To make it always active:

```bash
192.168.0.1 broadcast + netmask 255.255.255.0 + up
```

This is an example.

#### DHCP

How do we set up DHCP? Very simple:

```bash
ifconfig e1000g1 dhcp start
```

To make it permanently active, create a file:

```bash
touch /etc/dhcp.e1000g1
touch /etc/hostname.e1000g0
```

## Routing

The folks at Sun who can't do things like everyone else have their own `route` command. So to list the present routes:

```bash
netstat -rn
```

And to add or delete routes:

```bash
netstat -rn   # show current routes
netstat -rnv  # show current routes and their masks
route add destIP gatewayIP
route add destIP -netmask 255.255.0.0 gatewayIP
route delete destIP gatewayIP
route delete destIP -netmask 255.255.0.0 gatewayIP
```

## DNS

You'll now need to define the nsswitch to add name resolution:

```bash
...
hosts:      files dns
...
```

Then edit the resolv.conf file and insert the DNS servers like this:

```bash
nameserver 212.27.40.241
nameserver 212.27.40.240
```

## Routes

Here's how to add routes. To make them persistent, add the `-p` option:

```bash
route -p add 192.168.15.0 192.168.15.1
```

Persistent routes are stored in `/etc/inet/static_routes`.

## VIP

Here's how to manually set up a VIP in Solaris:

```bash
ifconfig <network_interface> addif <vip>/<mask> up
```

For example:

```bash
ifconfig nge1 addif 192.168.0.1/24 up
```

## Resources
- [https://www.sunsolarisadmin.com/solaris-10/dladm-display-link-statusspeedduplexstatisticsmtu/](https://www.sunsolarisadmin.com/solaris-10/dladm-display-link-statusspeedduplexstatisticsmtu/)
