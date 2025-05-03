---
weight: 999
url: "/Installer_Debian_sur_un_Mac_en_single_boot/"
title: "Installing Debian on a Mac in Single Boot"
description: "This guide explains how to install Debian on a Mac in single boot mode, including tips for shortening boot time, setting up auto-restart after power loss, and enabling Wake on LAN."
categories: ["Linux", "Debian"]
date: "2012-05-02T20:36:00+02:00"
lastmod: "2012-05-02T20:36:00+02:00"
tags: ["Mac", "Debian", "Boot", "Installation"]
toc: true
---

![Debian](/images/debian_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Operating System** | Debian 6 |
| **Website** | [Debian](https://www.debian.org/) |
| **Last Update** | 05/02/2012 |
{{< /table >}}

## Introduction

There are many methods for installing Debian on a Mac, some of which are less elegant than others. One approach that worked for me, but is quite messy, involves having an HFS+ partition with refit and Linux partitions that need to be synchronized with EFI (MBR -> EFI) with each new partition.

In short, these are cumbersome methods. That's why I looked into other solutions like gparted and target boot.

## Installing Debian

First, you can either boot from a gparted live CD, or like me, boot your Mac in target mode (via Firewire 1 or 2 port, holding the T key during boot). Then launch gparted on the machine connected to the Firewire, and:

- Delete the current partition table
- Create a new one in **MBR** format, not GPT
- Create the partitions you want and apply the configuration
- Reboot on the Debian CD and perform a normal installation

That's it, no more problems with flashing folders or similar issues.

## Shortening Boot Time

You'll quickly notice that it's annoying to wait 30 seconds for the Mac to find the right boot partition. We can speed up this process by specifying which partition to use. To do this, you'll need to boot from the Mac OS X installation DVD. You can use the Disk Utility to view your partitions and identify the one containing your Debian /boot directory. Then, open a terminal and run this command:

```bash
bless --device /dev/disk0s1 --setBoot --legacy --verbose
```

## Automatically Restarting After a Power Outage

To make your Mac automatically restart after a power outage, add this to the rc.local file (`/etc/rc.local`):

```bash
[...]
# PPC Mac Mini
echo server_mode=1 > /proc/pmu/options

# Intel Mac Mini
setpci -s 0:1f.0 0xa4.b=0

# nVidia Mac Mini
setpci -s 00:03.0 0x7b.b=0x19

# Unibody Mac Mini
setpci -s 0:3.0 -0x7b=20
```

Use the line that works for your Mac.

## Wake On LAN

If you want to enable Wake on LAN (WOL), here are the commands to add to rc.local (choose one of the three setpci lines, the one that works for your system) (`/etc/rc.local`):

```bash
[...]
## Wake on Lan
# Choose one of the 3 lines (use the working one)
setpci -d 8086:27b9 0xa4.b=0
setpci -s 00:03.0 0xa4.b=0
setpci -s 00:03.0 0x7b.b=19
ethtool -s eth0 wol g
```

You'll need the ethtool command. Install the ethtool package:

```bash
aptitude install ethtool
```

## Resources
- http://doc.ubuntu-fr.org/installation_macbook_sans_macosx
- http://blog.dhampir.no/content/wake-on-lan-on-a-n-intel-mac-mini-with-linux
