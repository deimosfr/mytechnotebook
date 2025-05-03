---
weight: 999
url: "/Iperf_\\:_Tester_sa_bande_passante_de_bout_en_bout/"
title: "Iperf: Testing End-to-End Bandwidth"
description: "How to use Iperf to measure bandwidth and network performance between endpoints"
categories: ["Linux", "Debian", "Ubuntu", "Network"]
date: "2008-02-04T10:44:00+02:00"
lastmod: "2008-02-04T10:44:00+02:00"
tags: ["Network", "Servers", "Windows", "Linux"]
toc: true
---

## Iperf

Iperf is a computer software tool for measuring various variables of an IP network connection. Iperf is developed by the National Laboratory for Applied Network Research ([NLANR](https://www.nlanr.net/)). Based on a client/server architecture and available on different operating systems, Iperf is an important tool for network administrators.

## How to Get It

Iperf is included in most Linux distributions by default. However, you can [follow this link](https://dast.nlanr.net/Projects/Iperf/#download) to download the version that corresponds to your operating system (Windows and MacOS X versions are also available).

For Debian / Ubuntu:

```bash
apt-get install iperf
```

For Fedora:

```bash
yum install iperf
```

## How Does Iperf Work?

Iperf must be run on two machines located at opposite ends of the network to be tested. The first machine runs Iperf in "server mode" (with the -s option), the second in "client mode" (with the -c option). By default, the network test is done using the TCP protocol (but it is also possible to use UDP mode with the -u option).

## How to Use It

Let's take the example of a network test between machine A and machine B.

On machine A, run the following command:

```bash
# iperf -s
```

Then on machine B, run the command:

```bash
# iperf -c <IP address of machine A>
```

The following result will be displayed:

```bash
------------------------------------------------------------
Client connecting to 192.168.29.1, TCP port 5001
TCP window size: 65.0 KByte (default)
------------------------------------------------------------
[ 3] local 192.168.29.157 port 50675 connected with 192.168.29.1 port 5001
[ ID] Interval       Transfer     Bandwidth
[ 3]  0.0-10.0 sec   110 MBytes   92.6 Mbits/sec
```

This gives us the actual throughput between machine A and machine B. Using the -i option, we can get other types of information such as transit delay or network jitter.

It is possible to evaluate your Internet connection via a public Iperf server located on the Internet:

![Iperf](/images/iperf.avif)

Here are 3 suggested command lines for the server on Linux:

- TCP 5001: $ iperf -s -m -w 500K -i 5
- TCP 4662: $ iperf -s -m -w 500K -i 5 -p 4662
- UDP 5001: $ iperf -s -i 5 -u

Here are 5 suggested command lines for the client on Linux:

- Upload only: $ iperf -c 212.27.33.25 -m -w 500K -i 5 -t 30
- Upload + download: $ iperf -c 212.27.33.25 -m -w 500K -i 5 -t 30 -r
- Simultaneous upload + download: $ iperf -c 212.27.33.25 -m -w 500K -i 5 -t 30 -d -P 2
- Upload + download on port 4662: $ iperf -c 212.27.33.25 -m -w 500K -i 5 -t 30 -p 4662 -r
- Upload + download in UDP at 80 Mb/s: $ iperf -c 212.27.33.25 -i 5 -t 30 -r -u -b 80M

The -w parameter is very important; it specifies the "TCP window size" as the default value is too small.
The window size value cannot exceed that of the operating system's TCP/IP stack.

![Iperf4662](/images/iperf4662.avif)

## IPERF in Multicast

Iperf can work in multicast mode (-B). Launch it as follows:

On the server:

- $ iperf -s -u -B 225.0.1.2

On the client:

- $ iperf -c 225.0.1.2 -u -b 3M

This generates a UDP multicast stream (on address 225.0.1.2) of 3 Mb/sec.

## IPERF with Linux 2.6.21 and Later

Starting with kernel 2.6.21, "The new high resolution timer option in the kernel causes usleep(0) to be a nop so the thread keeps running (until its quanta is exhausted)." => IPERF consumes much more CPU for the same throughput.

On a Pentium IV 2.8 GHz HT (bogomips = 5600):

- Kernel 2.6.20:

```bash
$ iperf -c 127.0.0.1
------------------------------------------------------------
Client connecting to 127.0.0.1, TCP port 5001
TCP window size: 49.4 KByte (default)
------------------------------------------------------------
[  3] local 127.0.0.1 port 54566 connected with 127.0.0.1 port 5001
[  3]  0.0-10.0 sec  5.33 GBytes  4.58 Gbits/sec
```

- Kernel 2.6.22:

```bash
$ iperf -c 127.0.0.1
------------------------------------------------------------
Client connecting to 127.0.0.1, TCP port 5001
TCP window size: 49.4 KByte (default)
------------------------------------------------------------
[  3] local 127.0.0.1 port 44642 connected with 127.0.0.1 port 5001
[  3]  0.0-10.1 sec    273 MBytes    228 Mbits/sec
```

A patch exists: http://dast.nlanr.net/Projects/Iperf2.0/patch-iperf-linux-2.6.21.txt  
A patched i686 binary is available: http://lafibre.info/images/iperf/iperf

## Resources

- [Official website](https://dast.nlanr.net/Projects/Iperf/) (in English)
- [User manual and batch files for simple use under Windows](https://lafibre.info/iperf) (in French)
