---
weight: 999
url: "/Réseau_\\:_créer_un_bonding/"
title: "Network: Creating Bonding"
description: "Learn how to configure network bonding in Linux for load balancing and fault tolerance with different operation modes and configuration methods."
categories: ["RHEL", "Red Hat", "Debian", "Linux", "Network"]
date: "2012-04-18T18:44:00+02:00"
lastmod: "2012-04-18T18:44:00+02:00"
tags:
  [
    "Bonding",
    "Network",
    "High Availability",
    "Load Balancing",
    "Fault Tolerance",
    "Linux",
    "Configuration",
  ]
toc: true
---

## Introduction

Channel bonding, also known as Port Trucking, allows applying a series of predefined policies on a group of network links combined into one. This group of network interfaces gives you the ability to perform load balancing and especially fault tolerance. Seven operation modes are available to combine necessity and flexibility. In this article, we will explore the possibilities offered and their implementation.

## Installation

### Prerequisites

Two prerequisites appear at the network switch level, where the interfaces are connected:

- Support and configuration of "port trunking" mode on the used ports
- Support for IEEE 802.3ad standard

For your Linux system, you need:
Network cards (preferably compatible with ethtools and miitools).

Bonding module for the kernel.

- Kernel 2.4.x

```text
[*] Network device support
<M> Bonding driver support
```

- Kernel 2.6.x

```text
[*] Networking support
<M> Bonding driver support
```

```bash
aptitude install ifenslave-2.6
```

## ethtools and miitools

The bonding driver can only use one interface if a cable is disconnected. Since hardware failures are the most common on a server, this provides high availability. Note that there are two small tools for querying the ETHTOOL and MII registers of network cards. The tools are mii-tool (from the net-tools package) and ethtool. Not all interfaces support these registers; in this case, we cannot detect disconnection (there is a solution using ARP requests, but that's beyond the scope of this article). To use this bonding function, you need to use modes such as active-backup; for more info, read the documentation :).

### MII

```bash
> mii-tool eth0
eth0: negotiated 100baseTx-FD, link ok
```

```bash
> mii-tool eth1
eth1: negotiated 100baseTx-FD flow-control, link ok
```

### Ethtool

With the Intel Pro dual card:

```bash
> ethtool eth2
Settings for eth2:
      Supported ports:  TP MII
      Supported link modes:   10baseT/Half 10baseT/Full
                              100baseT/Half 100baseT/Full
      Supports auto-negotiation: Yes
      Advertised link modes:  10baseT/Half 10baseT/Full
                              100baseT/Half 100baseT/Full
      Advertised auto-negotiation: Yes
      Speed: 100Mb/s
      Duplex: Full
      Port: Twisted Pair
      PHYAD: 1
      Transceiver: internal
      Auto-negotiation: on
      Supports Wake-on: pg
      Wake-on: d
      Link detected: yes
```

With the 3Com card, it doesn't work:

```bash
> ethtool eth0
Settings for eth0:
No data available
```

## The Different Modes

In this chapter, we'll see the different modes offered by the bonding module:

{{< table "table-hover table-striped" >}}
| Name | Mode | Description |
|------|------|-------------|
| **balance-rr**: load balancing | 0 | With load balancing, packets travel on one active network card, then on another, sequentially. The bandwidth is increased. If one of the network cards fails, the load balance skips this card and continues to rotate cyclically. |
| **active-backup**: active backup | 1 | This mode is simple redundancy with failover. Only one interface is active. As soon as its failure is detected, another interface is activated and takes over. Your bandwidth does not change. |
| **balance-xor**: Xor balance | 2 | An interface is assigned to send to the same MAC address. Thus transfers are parallelized and the choice of the interface follows the rule: (Source MAC address XOR Destination MAC address) modulo number of interfaces. |
| **broadcast**: broadcast | 3 | No particularity in this case, all data is transmitted on all active interfaces. No other rule. |
| **802.3ad**: 802.3ad standard | 4 | The 802.3ad standard allows link aggregation, dynamically expanding bandwidth. Groups are created dynamically based on common parameters. |
| **balance-tlb**: TLB balance | 5 | "TLB" for Traffic Load Balancing. Outgoing traffic is distributed according to the current load (calculated relative to the speed) of each interface. Incoming traffic is received by the current interface. If the receiving interface becomes inactive, another interface takes the MAC address of the inactive interface. |
| **balance-alb**: ALB balance | 6 | "ALB" for Adaptive load balancing. This is an extended mode of tlb balance, which includes receiving load balancing. Receive load balancing is performed at the ARP response level. The module intercepts ARP responses and changes the MAC address to that of one of the interfaces. |
{{< /table >}}

## Implementation

### Presentation with a Script

**WARNING: I use mode0!** This is only possible if **each interface is on a different switch**, otherwise beware of packet duplication.

Let's see a minimal script:

```bash
#!/bin/bash
#Loading the bonding module.
modprobe bonding mode=0 miimon=100
# We close the two ethernet interfaces
ifconfig eth0 down
ifconfig eth1 down
# We create a bond0 interface, which receives a dummy Mac address (this will change)
ifconfig bond0
# We start the interface, as if it were an ethernet interface
ifconfig bond0 192.168.0.4 netmask 255.255.255.0
# We add the interfaces to bond0
ifenslave bond0 eth0
ifenslave bond0 eth1
```

Let's look closer:

- Module loading:

Two options were included here: mode (load balancing round-robin) and miimon (Interface monitoring frequency)
**It's possible to put in the /etc/modules.conf file the loading of the bonding module or in /etc/modprobe.d/arch/i386:**

```text
alias bond0 bonding
options bond0 mode=0 miimon=100
```

This first example shows the case of a single bonding. It is indeed possible to have multiple bondings with different modes.
For this, let's modify the /etc/modules.conf accordingly:

```text
alias bond0 bonding
options bond0 -o bond0 mode=2 miimon=100 primary=eth1 max_bonds=2
```

- primary=eth1: forces eth1 interface to be primary
- max_bonds=2: allows having 2 bonds (so 4 interfaces but 2 bonds). Just increase this number if you want several bonds.
- -o bond0: this is essential if you want to use multiple bonds! You will then need to specify bond0, bond1, bond2...

### Adding Interfaces to bond0

The ifenslave command allows us to add ethernet interfaces to bond0 (then called slave interfaces). When adding the first network interface to bond0, the latter takes the MAC address of that interface. The other interfaces then lose their MAC addresses, covered by that of the bond0 interface.

To remove an interface, simply run the command:

```text
ifenslave -d bondx ethx
```

The bonding module then gives ethX back its real MAC address. If we remove the first MAC address (the one used by bond0), then bond0 retrieves that of eth1.

To verify our bonding:

```bash
> cat /proc/net/bonding/bond0
Ethernet Channel Bonding Driver: v2.5.0 (December 1, 2003)

Bonding Mode: load balancing (round-robin)
MII Status: up
MII Polling Interval (ms): 0
Up Delay (ms): 0
Down Delay (ms): 0

Slave Interface: eth0
MII Status: up
Link Failure Count: 0

Slave Interface: eth1
MII Status: up
Link Failure Count: 0
```

### Automatic Configuration Files

{{< tabs tabTotal="3">}}
{{% tab tabName="Debian" %}}

Here are the files to modify to load your configuration at startup, add to /etc/modules.conf or /etc/modprobe.d/arch/i386:

```text
alias bond0 bonding
options bond0 mode=0 miimon=100
```

In `/etc/network/interfaces` for active-backup:

```text
auto bond0
iface bond0 inet static
    address 10.31.1.5
    netmask 255.255.255.0
    network 10.31.1.0
    gateway 10.31.1.254
    slaves eth0 eth1
    bond_mode active-backup
    bond_miimon 100
    bond_downdelay 200
    bond_updelay 200
```

Or for 802.3ad in `/etc/network/interfaces` :

```text
auto lo
iface lo inet loopback

auto bond0
iface bond0 inet static
   address 192.168.1.1
   netmask 255.255.255.0
   gateway 192.168.1.254
   network 192.168.1.0
   bond-slaves enp1s0 enp2s0
   bond-mode 4
   bond-miimon 100
   bond-downdelay 200
   bond-updelay 200
```

{{% /tab %}}
{{% tab tabName="RedHat" %}}

If you are on RHEL 6, in the file /etc/modprobe.d/bonding.conf, create the following line:

```text
alias bond0 bonding
```

If you are on RHEL <6, it will be in the /etc/modprobes.conf file that you will need to insert this same content.

In /etc/sysconfig/network-script/, create the ifcfg-bond0 file to enter the bonding configuration:

```text
DEVICE=bond0
IPADDR=192.168.0.104
NETMASK=255.255.255.0
GATEWAY=192.168.0.1
ONBOOT=yes
BOOTPROTO=static
BONDING_OPTS="mode 1"
```

Then we will move on to our network interfaces to tell them that they will be used for bonding:

```text
DEVICE=eth0
MASTER=bond0
ONBOOT=yes
SLAVE=yes
BOOTPROTO=static
```

Do the same for interface 1:

```text
DEVICE=eth1
MASTER=bond0
ONBOOT=yes
SLAVE=yes
BOOTPROTO=static
```

{{% /tab %}}
{{< /tabs >}}

### Bonding Module Options

Here is a series of the most commonly used options for the bonding module that allow fine-tuning the operation of your bonding:

{{< table "table-hover table-striped" >}}
| Parameter | Description |
|-----------|-------------|
| Primary | Only for active-backup. Prioritizes a slave interface. It will become active again as soon as it can, even if another interface is active. |
| updelay | (0 by default) Latency time between discovering the reconnection of an interface and its re-use. |
| downdelay | (0 by default) Latency time between discovering the disconnection of an interface and its deactivation from bond0. |
| miimon | (0 by default) Frequency of monitoring interfaces by Mii or ethtool. The recommended value is 100. |
| use_carrier | (1 by default) Specifies the use of interface monitoring by miitool or by the network card itself (requires integrated instructions). |
| arp_interval (in ms) | ARP monitoring system, avoiding the use of miitool and ethtool. If no frame arrives during the arp_interval, up to 16 ARP requests are sent through this interface to 16 IP addresses. If no response is obtained, the interface is deactivated. |
| arp_ip_target | List of IP addresses, separated by a comma, used by ARP monitoring. If none is specified, ARP monitoring is inactive. |
{{< /table >}}

If you want the complete list of options, you can get it this way:

```bash
> modinfo bonding
filename:       /lib/modules/2.6.32-220.el6.x86_64/kernel/drivers/net/bonding/bonding.ko
author:         Thomas Davis, tadavis@lbl.gov and many others
description:    Ethernet Channel Bonding Driver, v3.6.0
version:        3.6.0
license:        GPL
srcversion:     B956376CB253D2B7312733C
depends:        ipv6
vermagic:       2.6.32-220.el6.x86_64 SMP mod_unload modversions
parm:           max_bonds:Max number of bonded devices (int)
parm:           tx_queues:Max number of transmit queues (default = 16) (int)
parm:           num_grat_arp:Number of gratuitous ARP packets to send on failover event (int)
parm:           num_unsol_na:Number of unsolicited IPv6 Neighbor Advertisements packets to send on failover event (int)
parm:           miimon:Link check interval in milliseconds (int)
parm:           updelay:Delay before considering link up, in milliseconds (int)
parm:           downdelay:Delay before considering link down, in milliseconds (int)
parm:           use_carrier:Use netif_carrier_ok (vs MII ioctls) in miimon; 0 for off, 1 for on (default) (int)
parm:           mode:Mode of operation; 0 for balance-rr, 1 for active-backup, 2 for balance-xor, 3 for broadcast, 4 for 802.3ad, 5 for balance-tlb, 6 for balance-alb (charp)
parm:           primary:Primary network device to use (charp)
parm:           primary_reselect:Reselect primary slave once it comes up; 0 for always (default), 1 for only if speed of primary is better, 2 for only on active slave failure (charp)
parm:           lacp_rate:LACPDU tx rate to request from 802.3ad partner; 0 for slow, 1 for fast (charp)
parm:           ad_select:803.ad aggregation selection logic; 0 for stable (default), 1 for bandwidth, 2 for count (charp)
parm:           xmit_hash_policy:balance-xor and 802.3ad hashing method; 0 for layer 2 (default), 1 for layer 3+4, 2 for layer 2+3 (charp)
parm:           arp_interval:arp interval in milliseconds (int)
parm:           arp_ip_target:arp targets in n.n.n.n form (array of charp)
parm:           arp_validate:validate src/dst of ARP probes; 0 for none (default), 1 for active, 2 for backup, 3 for all (charp)
parm:           fail_over_mac:For active-backup, do not set all slaves to the same MAC; 0 for none (default), 1 for active, 2 for follow (charp)
parm:           all_slaves_active:Keep all frames received on an interfaceby setting active flag for all slaves; 0 for never (default), 1 for always. (int)
parm:           resend_igmp:Number of IGMP membership reports to send on link failure (int)
```

## Some Examples

A simple redundancy with priority on eth1:

```bash
#!/bin/bash
modprobe bonding -o bond0 mode=1 miimon=100 priority=eth1
ifconfig eth0 down
ifconfig eth1 down
ifconfig bond0
ifconfig bond0 192.168.0.4 netmask 255.255.255.0
ifenslave bond0 eth0
ifenslave bond0 eth1
```

Load balancing with ARP request verification:

```bash
#!/bin/bash
modprobe bonding -o bond1 mode=0 arp_interval=1000 arp_ip_target=x.x.x.x
ifconfig eth0 down
ifconfig eth1 down
ifconfig eth2 down
ifconfig bond1
ifconfig bond1 192.168.1.5 netmask 255.255.255.0
ifenslave bond1 eth0
ifenslave bond1 eth1
ifenslave bond1 eth2
```

## Summary Table of Modes

Here are the different modes offered by the bonding module with load balancing support:

{{< table "table-hover table-striped" >}}
| Mode | Load-Balancing |
|------|---------------|
| 0 | incoming |
| 1 | none |
| 2 | incoming |
| 3 | none |
| 4 | none |
| 5 | incoming |
| 6 | incoming and outgoing |
{{< /table >}}

## Conclusion

Network card redundancy is often perceived as useless and expensive despite a simple installation. With these 7 modes of operation, you have the ability to obtain a high availability network interface, even responding to specific needs for load balancing. So why not take advantage of it?

## Other Documentation

[Bonding documentation on Ubuntu](/pdf/bonding_ubuntu.pdf)  
http://www.linux-foundation.org/en/Net:Bonding  
[NIC Bonding On Debian Lenny](/pdf/nic_bonding_on_debian_lenny.pdf)  
http://bisscuitt.blogspot.fr/2007/10/redhat-linux-rhel-4-nic-bonding.html
