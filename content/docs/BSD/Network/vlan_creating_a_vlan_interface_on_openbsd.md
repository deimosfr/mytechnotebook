---
weight: 999
url: "/VLAN_\\:_Créer_une_interface_VLAN_sous_OpenBSD/"
title: "VLAN: Creating a VLAN Interface on OpenBSD"
description: "How to create and configure VLAN interfaces on OpenBSD systems, both dynamically and statically."
categories: ["Linux", "Network"]
date: "2010-08-04T07:24:00+02:00"
lastmod: "2010-08-04T07:24:00+02:00"
tags: ["Servers", "Network", "OpenBSD", "VLAN"]
toc: true
---

## Introduction

A virtual local area network, commonly called [VLAN (for Virtual LAN)](https://fr.wikipedia.org/wiki/VLAN), is an independent logical computer network. Many VLANs can coexist on the same network switch.

## Configuration

### Dynamique

To create a VLAN interface on the fly:

```bash
ifconfig vlan110 create
```

Then assign a tag (VLAN ID) to this VLAN as well as the physical interface on which it should be created:

```bash
ifconfig vlan110 vlan 110 vlandev sis1
```

And finally, assign it a specific IP address:

```bash
ifconfig vlan110 inet 192.168.110.254 netmask 255.255.255.0
```

### Statique

Now, to keep this configuration active when the machine restarts (`/etc/hostname.vlan110`):

```bash
inet 192.168.110.254 255.255.255.0 NONE vlan 110 vlandev sis1
```
