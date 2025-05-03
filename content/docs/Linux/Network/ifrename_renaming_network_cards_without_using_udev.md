---
weight: 999
url: "/Ifrename_\\:_renomer_ses_cartes_réseaux_sans_utiliser_udev/"
title: "Ifrename: Renaming Network Cards Without Using Udev"
description: "How to use Ifrename to rename network interfaces based on MAC addresses without using udev"
categories:
  - Linux
date: "2007-04-05T15:49:00+02:00"
lastmod: "2007-04-05T15:49:00+02:00"
tags:
  - Linux
  - Networking
  - Configuration
toc: true
---

## Introduction

If you have multiple ethernet devices on a system, it's useful to make sure they are always given the device names that you expect. This can be helpful when you're managing upgrades - or for situations where you accidentally setup a system with eth1 plugged into a switch rather than eth0.

There are several different ways of managing the naming of devices if you're using a dynamic /dev system such as udev or hotplug - but the simplest system which works for most cases is provided by the ifrename package.

## Installation

To install it:

```bash
apt-get install ifrename libiw28
```

## Configuration

Once installed this package will let you rename devices based upon something that shouldn't change - their MAC addresses. (Finding MAC addresses of an ethernet device is simple.)

Once installed you may create a new file `/etc/iftab` to define the mapping between your ethernet device's MAC addresses and the interface names.

The contents of this file should look similar to this:

```bash
eth0 mac 00:17:31:56:BC:2D
eth1 mac 00:16:3E:2F:0E:9C
```

With this configuration file in place when you reboot next you'll discover that regardless of your kernel upgrading, that the network card with MAC address "00:17:31:56:BC:2D" will be setup as eth0, and that the card with MAC address "00:16:3E:2F:0E:9C" will be known as eth1.

(The actual renaming will happen automatically via the addition of `/etc/init.d/ifrename`.)

## Udev

You can edit this file too to change the card interfaces:

```bash
/etc/udev/rules.d/z25_persistent-net.rules
```
