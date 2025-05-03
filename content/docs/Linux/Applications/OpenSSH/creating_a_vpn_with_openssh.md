---
weight: 999
url: "/Créer_un_VPN_avec_OpenSSH/"
title: "Creating a VPN with OpenSSH"
description: "Learn how to create a VPN connection using OpenSSH to encrypt all traffic between two machines."
categories: ["Linux"]
date: "2008-07-26T12:20:00+02:00"
lastmod: "2008-07-26T12:20:00+02:00"
tags: ["Servers", "Network", "Verifications", "Mac OS X", "Divers", "Windows"]
toc: true
---

## 1 Introduction

Most SSH clients have the ability to perform local and remote port forwarding. This is a pretty neat use of SSH if you haven't ever seen it before. OpenSSH can take it one step further and provide a full VPN solution encrypting all network traffic on all ports between two machines. This is pretty powerful stuff. This is useful for a quick-and-dirty way to encrypt all traffic between two machines. For a longer term solution, you might want to check out how to configure IPsec or use OpenVPN. All three solutions have some really cool features and benefits.

OpenSSH is the most widely deployed open source SSH client/server solution today. Most Linux/BSD hosts I have encountered will use this as the client/server by default. Sun's SSH packages are based off of the OpenSSH distribution with some tweaks and modifications, but it's pretty close to OpenSSH's implementation.

## 2 Configuration

### 2.1 Verifications

Anyways, to create a VPN tunnel between two machines, two variables in sshd_config need to be tweaked as well as the presence of the tun/tap kernel module. This kernel module is available on most Linux/BSD distributions. It may have to be compiled and inserted into the Solaris kernel, or you can download it here.

```bash
$ uname -a
Linux locutus 2.6.25-14.fc9.i686 #1 SMP Thu May 1 06:28:41 EDT 2008 i686 i686 i386 GNU/Linux
$ lsmod | grep tun
tun                    11776  2
$ modinfo tun
filename:       /lib/modules/2.6.25-14.fc9.i686/kernel/drivers/net/tun.ko
alias:          char-major-10-200
license:        GPL
author:         (C) 1999-2004 Max Krasnyansky <maxk@qualcomm.com>
description:    Universal TUN/TAP device driver
srcversion:     12C02361DF16200902CDE64
depends:
vermagic:       2.6.25-14.fc9.i686 SMP mod_unload 686 4KSTACKS
```

### 2.2 SSH

So, the tun module has already been inserted into my running kernel. Next, set these two variables in sshd_config and have sshd re-read its configuration files…

```bash
...
PermitRootLogin yes
PermitTunnel yes
...
```

```bash
$ service sshd reload
Reloading sshd:                                            [  OK  ]
```

Make sure you don't mess up sshd_config and reload the daemon as your only way to access the machine! Console access is always a good thing.

Now all we need to do is to open the VPN tunnel itself. Here, I open a VPN tunnel to localhost (not really useful) but you can get the idea…

```bash
ssh -w any:any root@localhost
```

The any:any defines the local:remote "tun" device. We could have put 0:1 here, (tun0 as the local, tun1 as the remote) but any:any takes care of it for us in case there are any pre-existing tun devices in use.) or 0:0 if we were accessing a real remote machine. (I can't define two "tun0" devices on localhost)

### 2.3 Network

So, after my SSH session connects, sure enough the tun devices exist..

```bash
$ ifconfig tun0
tun0      Link encap:UNSPEC  HWaddr 00-00-00-00-00-00-00-00-00-00-00-00-00-00-00-00
POINTOPOINT NOARP MULTICAST  MTU:1500  Metric:1
RX packets:0 errors:0 dropped:0 overruns:0 frame:0
TX packets:0 errors:0 dropped:0 overruns:0 carrier:0
collisions:0 txqueuelen:500
RX bytes:0 (0.0 b)  TX bytes:0 (0.0 b)

$ ifconfig tun1
tun1      Link encap:UNSPEC  HWaddr 00-00-00-00-00-00-00-00-00-00-00-00-00-00-00-00
POINTOPOINT NOARP MULTICAST  MTU:1500  Metric:1
RX packets:0 errors:0 dropped:0 overruns:0 frame:0
TX packets:0 errors:0 dropped:0 overruns:0 carrier:0
collisions:0 txqueuelen:500
RX bytes:0 (0.0 b)  TX bytes:0 (0.0 b)
```

Now all we have to do is assign them some network addresses and let them know that its a point-to-point connection between the two…

```bash
$ ifconfig tun0 10.0.0.10 pointopoint 10.0.0.11
$ ifconfig tun1 10.0.0.11 pointopoint 10.0.0.10

$ ifconfig tun0
tun0      Link encap:UNSPEC  HWaddr 00-00-00-00-00-00-00-00-00-00-00-00-00-00-00-00
inet addr:10.0.0.10  P-t-P:10.0.0.11  Mask:255.255.255.255
UP POINTOPOINT RUNNING NOARP MULTICAST  MTU:1500  Metric:1
RX packets:0 errors:0 dropped:0 overruns:0 frame:0
TX packets:0 errors:0 dropped:0 overruns:0 carrier:0
collisions:0 txqueuelen:500
RX bytes:0 (0.0 b)  TX bytes:0 (0.0 b)

$ ifconfig tun1
tun1      Link encap:UNSPEC  HWaddr 00-00-00-00-00-00-00-00-00-00-00-00-00-00-00-00
inet addr:10.0.0.11  P-t-P:10.0.0.10  Mask:255.255.255.255
UP POINTOPOINT RUNNING NOARP MULTICAST  MTU:1500  Metric:1
RX packets:0 errors:0 dropped:0 overruns:0 frame:0
TX packets:0 errors:0 dropped:0 overruns:0 carrier:0
collisions:0 txqueuelen:500
RX bytes:0 (0.0 b)  TX bytes:0 (0.0 b)
```

Nice! Now all network traffic between the machines using those newly created addresses will be tunneled through SSH! Could this be a solution for NFSv3 and firewalls? ;-)

If there are multiple machines on the "other side" of the VPN that you would want to connect to, you will also need to add a route..

```bash
$ route add -net 10.0.0.0 netmask 255.255.255.0 gw 10.0.0.10 tun0
$ netstat -rn
Kernel IP routing table
Destination     Gateway         Genmask         Flags   MSS Window  irtt Iface
10.0.0.11       0.0.0.0         255.255.255.255 UH        0 0          0 tun0
10.0.0.10       0.0.0.0         255.255.255.255 UH        0 0          0 tun1
10.0.0.0        10.0.0.10       255.255.255.0   UG        0 0          0 tun0
```

So any traffic destined for 10.0.0.0/24 is gonna go out tun0 to be routed onwards and upwards.

We've also got to set up an arp entry..

```bash
$ arp -sD 10.0.0.11 eth0 pub
$ arp -an
? (10.0.0.11) at * PERM PUP on eth0
```

Ryan McGuire wrote a [pretty cool blog entry](https://wiki.enigmacurry.com/OpenSSH) about not only this feature, but some other really neat things with OpenSSH. I based a lot of this article after learning them from his site. Check out [his python script](https://www.enigmacurry.com/blog-post-files/vpn-up.py) that will automate a lot of this for you. Thanks Ryan!

## 3 Resources

http://prefetch.net/blog/index.php/2008/06/26/opensshs-vpn/
