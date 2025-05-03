---
weight: 999
url: "/Trunking_\\:_Créer_du_trunking_(bonding)_sur_OpenBSD/"
title: "Trunking: Creating Trunking (bonding) on OpenBSD"
description: "This guide explains how to configure network interface trunking (bonding) on OpenBSD systems to combine multiple physical interfaces into a single virtual interface."
categories: ["Network", "Servers"]
date: "2007-11-15T15:55:00+02:00"
lastmod: "2007-11-15T15:55:00+02:00"
tags: ["OpenBSD", "Network", "Trunking", "Bonding"]
toc: true
---

## Introduction

The following is used on OpenBSD since version 3.9 to combine two physical interfaces (fxp1, fxp2) into a single virtual interface (trunk0). This method allows one to take the feeds from a traditional two-output tap and present a single virtual interface to NSM applications.

## Configuration

Modify the interface and put yours:

```bash
ifconfig fxp1 up
ifconfig fxp2 up
ifconfig trunk0 trunkport fxp1 up
ifconfig trunk0 trunkport fxp2 up
ifconfig trunk0 trunkproto roundrobin up
```

### Modes

If you don't need roundrobin, choose the mode that you would like:

- **roundrobin**: Distributes outgoing traffic using a round-robin scheduler through all active ports and accepts incoming traffic from any active port.
- **failover**: Sends and receives traffic only through the master port. If the master port becomes unavailable, the next active port is used. The first interface added is the master port; any interfaces added after that are used as failover devices.
- **loadbalance**: Balances outgoing traffic across the active ports based on hashed protocol header information and accepts incoming traffic from any active port. The hash includes the Ethernet source and destination address, and, if available, the VLAN tag, and the IP source and destination address.
- **broadcast**: Sends frames to all ports of the trunk and receives frames on any port of the trunk.
- **none**: This protocol is intended to do nothing: it disables any traffic without disabling the trunk interface itself.

To make this configuration permanent between reboots:

```bash
echo "up" > /etc/hostname.fxp1
echo "up" > /etc/hostname.fxp2
echo "trunkproto roundrobin trunkport fxp1 trunkport fxp2 172.16.1.100 netmask 255.255.255.0" > /etc/hostname.trunk0
```

Remember to replace fxp1 and fxp2 with the network interfaces on your OpenBSD system (e.g., em0, xl0, rl0, etc.). Don't forget to add your good IP address.

OpenBSD 3.8 only supported the round robin trunk protocol.

## References

[OpenBSD trunk man](https://www.openbsd.org/cgi-bin/man.cgi?query=trunk&apropos=0&sektion=4&manpath=OpenBSD+Current&arch=i386&format=html)
