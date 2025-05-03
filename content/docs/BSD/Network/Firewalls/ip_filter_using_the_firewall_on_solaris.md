---
weight: 999
url: "/IP_Filter_\\:_Utilisation_du_firewall_sous_Solaris/"
title: "IP Filter: Using the Firewall on Solaris"
description: "A guide on how to configure and use IP Filter, a firewall solution for Solaris operating systems."
categories: ["Unix", "Solaris", "Security"]
date: "2010-12-28T07:44:00+02:00"
lastmod: "2010-12-28T07:44:00+02:00"
tags: ["Solaris", "Firewall", "IPFilter", "Security", "Network"]
toc: true
---

## Introduction

IPFilter (commonly referred to as ipf) is an open source software package that provides firewall services and network address translation (NAT) for many UNIX-like operating systems. The author and software maintainer is Darren Reed. IPFilter supports both IPv4 and IPv6 protocols, and is a stateful firewall.

IPFilter is delivered with FreeBSD, NetBSD and Solaris 10. It used to be a part of OpenBSD, but it was removed in May 2001 due to problems with the license of IP Filter, after negotiations between Theo de Raadt and Reed broke down. At first glance, the license looks a lot like BSD Licenses, but does not allow redistribution of modified versions. Reed came back with another proposal but it was already too late. The software was removed from OpenBSD.

IPFilter can be installed as a runtime-loadable kernel module or directly incorporated into the operating system kernel, depending on the specifics of each kernel and user preferences. The software's documentation recommends the module approach, if possible.

## Usage

- To activate Solaris IP Filter:

```bash
svcadm enable network/ipfilter
```

- To enable IPF:

```bash
ipf -E
```

- To disable IPF:

```bash
ipf -D
```

- Reload configuration

```bash
ipf -f config_file
```

- Activate Nat (optional):

```bash
ipfnat -f config_file
```

- Remove active rule set from the kernel:

```bash
ipf -Fa
```

- Remove incoming packet filtering rules:

```bash
ipf -Fi
```

- Remove outgoing packet filtering rules:

```bash
ipf -Fo
```

- Get stats:

```bash
ipfstat -io
```

or

```bash
ipfstat
```

## Configuration

### Files locations

The default configurations files are located in `/etc/ipf`:

- `ipf.conf`: Containing the main configuration
- `ipnat.conf`: Containing Nat configuration
- `ippool.conf`: Define server pool

If files are named like this, they will be loaded at boot time. If you don't want, rename them in an other name.

### Redirect all incoming connections to a specific IP on a specific port

This is an example to forward any incoming connection to the port 4175:

```bash
svcadm enable svc:/network/ipv4-forwarding:default
ipf -E
rdr e1000g0 192.168.76.0/24 port 4175 -> 192.168.15.30 port 4175 tcp
rdr e1000g1 192.168.76.0/24 port 4175 -> 192.168.15.30 port 4175 tcp
map e1000g0 from any to 192.168.15.30 port = 4175 -> 0/32
map e1000g1 from any to 192.168.15.30 port = 4175 -> 0/32
```

192.168.76.0 is the subnet of the "forwarder", 192.168.15.30 is the destination ip.
