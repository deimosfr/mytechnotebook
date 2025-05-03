---
weight: 999
url: "/Xen_avec_Bonding_+_Vlan_Tagging/"
title: "Xen with Bonding + VLAN Tagging"
description: "This guide explains how to configure Xen with network bonding and VLAN tagging to maximize ethernet interfaces usage with fault tolerance and load balancing."
categories: ["Linux", "Networking", "Virtualization"]
date: "2008-01-27T11:48:00+02:00"
lastmod: "2008-01-27T11:48:00+02:00"
tags: ["Xen", "VLAN", "Bonding", "Networking", "Bridge", "Debian"]
toc: true
---

## Introduction

As powerful and flexible as it may be, Xen's network configuration often feels like an obstacle course as soon as you stray a bit from the beaten path.

Recently, I had to set up a Xen-based solution in a somewhat special network topology. Each server where Xen was deployed had 2 ethernet adapters, each with 2 gigabit ethernet interfaces. The goal was to maximize the use of the different ethernet interfaces while providing both fault tolerance and load balancing. To make things even more interesting, the network was segmented into VLANs, each VLAN segmenting the network by function (DMZ, LAN, ADMIN, etc.).

The different points that interest us here at the system administration level are:

- [Channel Bonding]({{< ref "docs/Linux/Network/network_creating_bonding.md">}})
- [VLAN Tagging]({{< ref "docs/Servers/Network/setting_up_vlan.md" >}})
- Bridging

After some research on the web, you quickly realize that you will need to create your own network script for Xen. Shame on me, I didn't really feel capable of doing that (too complicated/I don't know python/it's still too obscure for me). So I decided to move the problem from Xen to the host operating system, with most of the configuration happening at the OS level (in this case a Debian Etch), with Xen simply using what already exists.

The final goal is as follows:

![Xenvlanbond.png](/images/xenvlanbond.avif)

Don't worry about bond0, we're only interested in bond1, eth1, and eth2 here.

## Channel Bonding Configuration

First, we need to configure and load the bonding module. In Debian, go to `/etc/modprobe.d/arch/i386` (even if you're on an x86_64 architecture):

```bash
# Channel bonding configuration
#
# mode=802.3ad
#        Use 802.3ad for link aggregation (require LACP compatible switch and
#        additionnal switch configurations)
# miimon= 100
#       MII monitoring frequency in ms
# use_carrier=0
#       Use an alternative and deprecated method (ethtool ioctls) to monitor
#       link status. Works better with this mode (maybe network drivers issues)
# lacp_rate=fast
#       Transmit LACPDU packet each second instead of each 30 seconds
# max_bonds=2
#       Indicate that there's 2 bonding device (bond0 and bond1) even if
#       only bond0 is explicitely configured. With max_bonds options there's
#       no way to configure separately each bond device.
#       NOTE : -o bond# (better) seems not to be supported
alias bond0 bonding
options bond0 mode=802.3ad miimon=100 use_carrier=0 lacp_rate=fast max_bonds=2
```

Yes, I know I'm only specifying the configuration for bond0 here. However, by changing the value of max_bonds, you configure the number of interfaces you want. Despite what is specified in the bonding.txt.gz file of the Linux documentation (aptitude install linux-doc), I was not able to have a per-interface configuration. The kernel is the default one provided with Debian Etch.

### Configuration of the bond1 Interface

Next, we need to edit the `/etc/network/interfaces` file as follows:

```bash
iface bond1 inet manual
       pre-up /sbin/ifconfig eth0 up
       pre-up /sbin/ifconfig eth3 up
       pre-up /sbin/ifconfig bond1 up
       pre-up /sbin/ifenslave bond1 eth0
       pre-up /sbin/ifenslave bond1 eth3
       post-down /sbin/ifenslave -d bond1 eth3
       post-down /sbin/ifenslave -d bond1 eth0
       post-down /sbin/ifconfig eth0 down
       post-down /sbin/ifconfig eth3 down
       post-down /sbin/ifconfig bond1 down
```

Nothing too complicated here. Note a few things though:

- You need to install the ifenslave-2.6 package.
- I use the pre-up directive instead of up which is more commonly found. I had some weird behaviors with the latter (like interfaces not coming up).
- No IP configuration is specified. We're staying at layer 2 of the OSI model.

At this point, you can start playing with bringing your bond1 interface up and down:

```
ifup bond1
ifdown bond1
```

## VLAN

On top of this bond1 interface, we'll now set up one interface per VLAN. First, you need to install the vlan package:

```bash
# vlan 3 = LAN
iface vlan3 inet manual
        vlan-raw-device bond1
        pre-up /sbin/ifup bond1
        pre-up /sbin/ifconfig bond1 up
        pre-up /sbin/modprobe 8021q

# vlan 2 = DMZ
iface vlan2 inet manual
        vlan-raw-device bond1
        pre-up /sbin/ifup bond1
        pre-up /sbin/ifconfig bond1 up
        pre-up /sbin/modprobe 8021q
```

Note again the use of the pre-up directive (it's really the only one that worked for me) and the loading of the module handling VLAN (802.1Q standard).

## Creating Bridges for DomUs

First, it's important to understand that a software bridge under Linux works exactly the same way as a physical switch. In everyday life (of an IT person), what we're doing here is equivalent to:

- Taking a switch
- Plugging network cables into it
- Obviously, each cable is connected to a computer

In my particular case, I need to create 2 bridges: one over the DMZ VLAN, the other over the LAN VLAN:

```bash
# xen bridge for hosts which needs LAN network
auto xenbrlan
iface xenbrlan inet dhcp
        pre-up /sbin/ifup vlan3
        pre-up /sbin/ifconfig vlan3 up
        pre-up /usr/sbin/brctl addbr xenbrlan
        pre-up /usr/sbin/brctl addif xenbrlan vlan3
        pre-up /usr/sbin/brctl setfd xenbrlan 0
        post-down /usr/sbin/brctl delif xenbrlan vlan3
        post-down /usr/sbin/brctl delbr xenbrlan
        post-down /sbin/ifdown vlan3

# xen bridge for hosts which needs DMZ network
auto xenbrdmz
iface xenbrdmz inet manual
        pre-up /sbin/ifup vlan2
        pre-up /sbin/ifconfig vlan2 up
        pre-up /usr/sbin/brctl addbr xenbrdmz
        pre-up /usr/sbin/brctl addif xenbrdmz vlan2
        pre-up /usr/sbin/brctl setfd xenbrdmz 0
        post-down /usr/sbin/brctl delif xenbrdmz vlan2
        post-down /usr/sbin/brctl delbr xenbrdmz
        post-down /sbin/ifdown vlan2
```

Note that like any respectable switch, our xenbrlan and xenbrdmz bridges don't need an IP address (it's a layer 2 device) to function. At this point, you can play with bringing your bridges up and down on the fly:

```
ifup xenbrlan xenbrdmz
ifdown xenbrlan xenbrdmz
```

## Configuring Xen

The only small difficulty here is to create a little script that will allow the use of 2 different bridges. To do this, we'll create a file `/etc/xen/script/my-network-bridge` as follows:

```bash
#!/bin/sh
dir=$(dirname "$0")

"$dir/network-bridge" "$@" vifnum=0 bridge=xenbrlan
"$dir/network-bridge" "$@" vifnum=1 bridge=xenbrdmz
```

You'll also need to tell Xen to use our script instead of the default one. To do this, go to the `xend-config.sxp` file and replace `(network-script network-bridge)` with `(network-script 'my-network-bridge')`.

All that's left is to specify in the configuration of each DomU the bridge you want to use:

```
vif     = ['mac=02:00:00:00:00:01, bridge=xenbrdmz']
```

You can now start your VMs :-)

## Resources
- [https://anothergeekwebsite.com/fr/2007/06/xen-vlan-et-bonding-oui-oui-tout-ca](https://anothergeekwebsite.com/fr/2007/06/xen-vlan-et-bonding-oui-oui-tout-ca)
