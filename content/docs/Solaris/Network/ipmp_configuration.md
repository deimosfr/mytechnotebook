---
weight: 999
url: "/Configuration_de_l'IPMP/"
title: "IPMP Configuration"
description: "Learn how to set up and configure IP Network Multipathing (IPMP) on Solaris systems for improved network reliability and load balancing"
categories: ["Solaris", "Network", "Servers"]
date: "2011-04-06T09:41:00+02:00"
lastmod: "2011-04-06T09:41:00+02:00"
tags: ["Solaris", "IPMP", "networking", "high availability", "failover"]
toc: true
---

## Introduction

IP Network Multipathing (IPMP) enables you to detect interface failures and transparently switch network access for a system with multiple interfaces on the same IP link. IPMP also allows load balancing of packets for systems with multiple interfaces.

The equivalent on Linux is called bonding and on BSD, it's called trunking.

IPMP improves the reliability, availability, and network performance of systems with multiple physical interfaces. Sometimes a physical interface or the networking hardware connected to that interface fails or requires maintenance. Traditionally, it is then impossible to contact the system through any of the IP addresses associated with the failed interface. Additionally, any existing connections to the system using those IP addresses are disrupted.

By using IPMP, you can configure one or more physical interfaces into an IPMP group. After IPMP configuration, the system automatically monitors the interfaces in the group for failure. If an interface in the group fails or is removed for maintenance, IPMP automatically migrates, or fails over, the failed interface's IP addresses. The recipient of these addresses is a functioning interface in the failed interface's IPMP group. The failover feature of IPMP preserves connectivity and prevents disruption of any existing connections. Additionally, IPMP improves overall network performance by spreading outbound network traffic across all interfaces in the IPMP group. This process is called load spreading.

## Configuration

### Prerequisites

IPMP is built into the Solaris operating system and does not require any special hardware. Any interface supported by Solaris can be used with IPMP. However, your network configuration and topology must meet the following IPMP-related requirements:

* All interfaces in an IPMP group must have unique MAC addresses. Note that by default, network interfaces of SPARC-based systems share the same MAC address. Therefore, you need to explicitly change the default address to use IPMP on SPARC-based systems.
* All interfaces in an IPMP group must be of the same media type (e.g., Ethernet with Ethernet, Fiber with Fiber, but not mixed).
* All interfaces in an IPMP group must be on the same IP link (same subnet).
* Depending on your requirements for failure detection, you'll either need to use specific types of network interfaces or configure additional IP addresses on each network interface.

### /etc/hosts

You'll need to configure the hosts file to specify the IPs of the machines, test IPs, and virtual IPs:

```bash
#
# Internet host table
#
::1             localhost
127.0.0.1       localhost
192.168.0.72    sun-node1
192.168.0.74    sun-node1-if2
192.168.0.73    sun-node1-test-e1000g0
192.168.0.77    vip1
192.168.0.78    vip2
```

### /etc/netmasks

Now you'll need to set the network and subnet in the /etc/netmasks file:

```bash
#
# The netmasks file associates Internet Protocol (IP) address
# masks with IP network numbers.
#
#       network-number  netmask
#
# The term network-number refers to a number obtained from the Internet Network
# Information Center.
#
# Both the network-number and the netmasks are specified in
# "decimal dot" notation, e.g:
#
#               128.32.0.0 255.255.255.0
#
192.168.0.0     255.255.255.0
```

### Creating an IPMP Group

You'll need to create an IPMP group and add network cards to it. If you want fault tolerance, you'll need to activate a test IP. This IP will not be part of the VIPs (Virtual Private Interface). Configure as follows:

```bash
192.168.0.72 netmask + broadcast + group ipmp0 up
addif 192.168.0.73 deprecated -failover netmask + broadcast + up
```

```bash
192.168.0.74 netmask + broadcast + group ipmp0 standby up
```

Here are the meanings:

* deprecated: Indicates that the test address is not used for outgoing packets.
* failover: Indicates that the test address does not failover when the interface fails.
* standby: Marks the interface as the standby interface.

Restart afterward to apply the configuration and check that it's properly taken into account at reboot.

### Modifying the Probing Time

If you want to change the interval at which the system will detect a disconnection of an interface, edit the following file:

```bash
#
#pragma ident   "@(#)mpathd.dfl 1.2     00/07/17 SMI"
#
# Time taken by mpathd to detect a NIC failure in ms. The minimum time
# that can be specified is 100 ms.
# 
FAILURE_DETECTION_TIME=10000
#
# Failback is enabled by default. To disable failback turn off this option
#
FAILBACK=yes
#
# By default only interfaces configured as part of multipathing groups 
# are tracked. Turn off this option to track all network interfaces 
# on the system
#
TRACK_INTERFACES_ONLY_WITH_GROUPS=yes
```

## Testing

To test now, it's quite simple. Look at the current configurations:

```bash
lo0: flags=2001000849<UP,LOOPBACK,RUNNING,MULTICAST,IPv4,VIRTUAL> mtu 8232 index 1
        inet 127.0.0.1 netmask ff000000 
e1000g0: flags=1000843<UP,BROADCAST,RUNNING,MULTICAST,IPv4> mtu 1500 index 2
        inet 192.168.0.72 netmask ffffff00 broadcast 192.168.0.255
        groupname ipmp0
        ether 0:1e:68:49:ae:98 
e1000g0:1: flags=9040843<UP,BROADCAST,RUNNING,MULTICAST,DEPRECATED,IPv4,NOFAILOVER> mtu 1500 index 2
        inet 192.168.0.73 netmask ffffff00 broadcast 192.168.0.255
e1000g0:2: flags=1000843<UP,BROADCAST,RUNNING,MULTICAST,IPv4> mtu 1500 index 2
        inet 192.168.0.74 netmask ffffff00 broadcast 192.168.0.255
e1000g1: flags=69040842<BROADCAST,RUNNING,MULTICAST,DEPRECATED,IPv4,NOFAILOVER,STANDBY,INACTIVE> mtu 0 index 3
        inet 0.0.0.0 netmask 0 
        groupname ipmp0
        ether 0:1e:68:49:ae:99
```

Here you can clearly see the IPMP interfaces as well as the interface that is in standby and inactive.
Now, if we disconnect the first interface, the failover will happen automatically :-)

```bash
lo0: flags=2001000849<UP,LOOPBACK,RUNNING,MULTICAST,IPv4,VIRTUAL> mtu 8232 index 1
        inet 127.0.0.1 netmask ff000000 
e1000g0: flags=19000802<BROADCAST,MULTICAST,IPv4,NOFAILOVER,FAILED> mtu 0 index 2
        inet 0.0.0.0 netmask 0 
        groupname ipmp0
        ether 0:1e:68:49:ae:98 
e1000g0:1: flags=19040803<UP,BROADCAST,MULTICAST,DEPRECATED,IPv4,NOFAILOVER,FAILED> mtu 1500 index 2
        inet 192.168.0.73 netmask ffffff00 broadcast 192.168.0.255
e1000g1: flags=21040842<BROADCAST,RUNNING,MULTICAST,DEPRECATED,IPv4,STANDBY> mtu 1500 index 3
        inet 192.168.0.72 netmask ffffff00 broadcast 192.168.0.255
        groupname ipmp0
        ether 0:1e:68:49:ae:99 
e1000g1:1: flags=21000843<UP,BROADCAST,RUNNING,MULTICAST,IPv4,STANDBY> mtu 1500 index 3
        inet 192.168.0.74 netmask ffffff00 broadcast 192.168.0.255
```

Now the interface has gone into failed state. And there was no interruption during this period. The interface that was in standby is now no longer inactive. So now it's functional, and if we reconnect the interface, we return to the previous configuration.

## Resources
- http://docs.sun.com/app/docs/doc/820-2982/ipmptm-1?l=fr&a=view
- http://www.eng.auburn.edu/~doug/howtos/multipathing.html
