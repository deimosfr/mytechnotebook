---
weight: 999
url: "/Port_forwarding_depuis_dom0_vers_bridged_domU_avec_IPVS/"
title: "Port forwarding from dom0 to bridged domU with IPVS"
description: "How to forward ports from a dom0 to a bridged domU in Xen using IPVS, overcoming Netfilter issues"
categories: ["Debian", "Linux"]
date: "2007-09-16T21:32:00+02:00"
lastmod: "2007-09-16T21:32:00+02:00"
tags: ["Network", "Xen", "IPVS", "Virtualization", "Linux"]
toc: true
---

## Introduction

Those following zone0's adventures know: Netfilter sucks badly when it comes to forwarding simple ports from a dom0 to a bridged domU. That's just how it is, we don't know where it comes from, maybe from 64-bit, maybe from BSD domUs, maybe from the Xen kernel, maybe, maybe, maybe. Anyway, after many hours of tweaking, debugging, tcpdump and so on, we decided on [IPVS](https://www.linuxvirtualserver.org/software/ipvs.html). I've posted the results of our experiences here, so if you also want to set up a simple port forwarding between Xen domains, you won't have to waste an entire Saturday and miss the techno-parade.

I'll add that [this excellent tutorial on IPVS](https://www.ultramonkey.org/papers/lvs_tutorial/html/) will allow you to quickly familiarize yourself with the tool.

By the way, I know, IPVS is OLD.

## Mise en place

We, an OSS advocacy group, setup a Xen 3.1 machine composed of:

- a 64 bits dom0 running Debian stable amd64
- 2 hvm domUs running OpenBSD amd64
- 2 hvm domUs running NetBSD i386

This machine is to be hosted and reachable from the Internet, but it will only have one public IP. Naturally, our first tought was to port-forward using iptables / netfilter. We didn't really though it would be an issue... and that was a mistake :) We tried many options, read many hints, even on this list, but no matter what, the port-forwarding, using a ultra-classic PREROUTING / FORWARD rule, was given a TCP RST in the best scenario. We read here stories about activating NAT / masquerading on the domU to fix (???) this issue, but as the machine is meant to be hosted, that was not the cleanest approach.

And then we took a look at IPVS (http://www.linuxvirtualserver.org/software/ipvs.html), an opensource Linux kernel module initially meant to act as a loadbalancer. We thought that providing a unique real server (the domU) to the VIP would do the trick... and it did! Here's a quick example of a working configuration:

- dom0 has a public IP address, no services but ssh available
- domU has a RFC1918 address, linked to a bridge on the second ethernet interface of the dom0

We want to redirect the port 2222 of the dom0 to the port 22 of the domU:

- Install ipvsadm on the dom0 (apt-get install ipvsadm on debian)
- Setup the VIP:

```
ipvsadm -A -t <public_ip>:2222 -s rr
```

We choosed the Round-Robin algorithm, but obviously this has no effect for us as there will be only one real server behind the loadbalancer

- Insert domU's private IP on the VIP:

```
ipvsadm -a -t <public_ip>:2222 -r <domU_private_ip>:22 -m
```

Here we use the simple masquerading mode of IPVS

- See the output:

```bash
$ root@dom0:~# ipvsadm -L
IP Virtual Server version 1.2.1 (size=4096)
Prot LocalAddress:Port Scheduler Flags
  -> RemoteAddress:Port           Forward Weight ActiveConn InActConn
TCP  dom0:2222 rr
  -> shells:ssh                   Masq    1      0          0
```

And finally, from an outside machine:

```bash
$ imil@tatooine:~$ ssh -p 2222 dom0
imil@dom0_public_ip's password:
Last login: Sun Sep 16 01:15:40 2007 from somewhere_else
OpenBSD 4.1 (GENERIC) #874: Sat Mar 10 19:09:51 MST 2007

imil@shells
~$
```

It Works!

Hope this method can save time to some of you, for us it's now the perfect solution as it provides us also the ability to loadbalance services on other domU's.

## References

http://www.gcu.info/2411/2007/09/16/si-tu-casses-ta-cuiller-prend-une-pellle-sin-youn-beol-prophete/
