---
weight: 999
url: "/NTP_\\:_Création_d'un_serveur_NTP/"
title: "NTP: Creating an NTP Server"
description: "Guide for setting up and configuring an NTP server on various operating systems including Debian and Solaris."
categories: ["Debian", "Linux", "Solaris"]
date: "2011-10-06T11:10:00+02:00"
lastmod: "2011-10-06T11:10:00+02:00"
tags: ["NTP", "Network", "Servers", "Linux", "Debian", "Solaris", "Windows"]
toc: true
---

## Introduction

The NTP server is a time server. It provides time based on atomic clock references.

## Installation

### Debian

Super quick:

```bash
apt-get install ntp
```

### Solaris

On Solaris, NTP is normally installed by default. Otherwise, search for packages containing *ntpd*.
You will need to copy the server configuration file:

```bash
cp /etc/inet/ntp.server /etc/inet/ntp.conf
```

## Configuration

Here is some information before diving into configuration files:

* stratum: indicates that it's not very reliable data
* 127.127.1.0: time given by the machine (BIOS)
* peer: defines a server of the same stratum as the server (same level/same geographic location), which allows for more reliability. Using 1 server + peers provides a more refined configuration than having only geographically distant servers.

### Debian

Here's the configuration file for a time server (`/etc/inet/ntp.conf`):

```bash
# /etc/ntp.conf, configuration for ntpd

# ntpd will use syslog() if logfile is not defined
logfile /var/log/ntpd

driftfile /var/lib/ntp/ntp.drift
statsdir /var/log/ntpstats/

statistics loopstats peerstats clockstats
filegen loopstats file loopstats type day enable
filegen peerstats file peerstats type day enable
filegen clockstats file clockstats type day enable
# Servers to use to update
server 0.fr.pool.ntp.org prefer
server 1.fr.pool.ntp.org
server 2.fr.pool.ntp.org

# ... and use the local system clock as a reference if all else fails
server 127.127.1.0
# fudge 127.127.1.0 stratum 13

# Local users may interrogate the ntp server more closely.
restrict 127.0.0.1

# Allow local network to synchronize to this server
restrict 192.168.0.0 mask 255.255.255.0 nomodify notrap

broadcastdelay 0.008
```

After configuring and restarting the service, **you must wait a few minutes for the server to synchronize**.

### Solaris

Here's the configuration file for Solaris (`/etc/inet/ntp.conf`):

```bash
#
# Copyright 1996-2003 Sun Microsystems, Inc.  All rights reserved.
# Use is subject to license terms.
#
# /etc/inet/ntp.server
#
# An example file that could be copied over to /etc/inet/ntp.conf and
# edited; it provides a configuration template for a server that
# listens to an external hardware clock, synchronizes the local clock,
# and announces itself on the NTP multicast net.
#

# This is the external clock device.  The following devices are
# recognized by xntpd 3-5.93e:
#
# XType Device    RefID          Description
# -------------------------------------------------------
#  1    local     LCL            Undisciplined Local Clock
#  2    trak      GPS            TRAK 8820 GPS Receiver
#  3    pst       WWV            PSTI/Traconex WWV/WWVH Receiver
#  4    wwvb      WWVB           Spectracom WWVB Receiver
#  5    true      TRUE           TrueTime GPS/GOES Receivers
#  6    irig      IRIG           IRIG Audio Decoder
#  7    chu       CHU            Scratchbuilt CHU Receiver
#  8    parse     ----           Generic Reference Clock Driver
#  9    mx4200    GPS            Magnavox MX4200 GPS Receiver
# 10    as2201    GPS            Austron 2201A GPS Receiver
# 11    arbiter   GPS            Arbiter 1088A/B GPS Receiver
# 12    tpro      IRIG           KSI/Odetics TPRO/S IRIG Interface
# 13    leitch    ATOM           Leitch CSD 5300 Master Clock Controller
# 15    *         *              TrueTime GPS/TM-TMD Receiver
# 17    datum     DATM           Datum Precision Time System
# 18    acts      ACTS           NIST Automated Computer Time Service
# 19    heath     WWV            Heath WWV/WWVH Receiver
# 20    nmea      GPS            Generic NMEA GPS Receiver
# 22    atom      PPS            PPS Clock Discipline
# 23    ptb       TPTB           PTB Automated Computer Time Service
# 24    usno      USNO           USNO Modem Time Service
# 25    *         *              TrueTime generic receivers
# 26    hpgps     GPS            Hewlett Packard 58503A GPS Receiver
# 27    arc       MSFa           Arcron MSF Receiver
#
# * All TrueTime receivers are now supported by one driver, type 5.
#   Types 15 and 25 will be retained only for a limited time and may
#   be reassigned in future.
#
# Some of the devices benefit from "fudge" factors.  See the xntpd
# documentation.

# Either a peer or server.  Replace "XType" with a value from the
# table above.
# Servers to use to update
server 0.fr.pool.ntp.org prefer
server 1.fr.pool.ntp.org
server 2.fr.pool.ntp.org
# fudge 127.127.XType.0 stratum 0

broadcast 224.0.1.1 ttl 4

enable auth monitor
driftfile /var/ntp/ntp.drift
statsdir /var/ntp/ntpstats/
filegen peerstats file peerstats type day enable
filegen loopstats file loopstats type day enable
filegen clockstats file clockstats type day enable

keys /etc/inet/ntp.keys
trustedkey 0
requestkey 0
controlkey 0
```

Create a drift file as indicated in the configuration:

```bash
touch /etc/inet/ntp.drift
```

Then enable the service:

```bash
svcadm enable svc:/network/ntp:default
```

## NTP Clients

### Linux

To synchronize a Unix/Linux machine, use this command:

```bash
ntpdate my_time_server
```

### Solaris

It's quite simple on Solaris:

```bash
cp /etc/inet/ntp.client /etc/inet/ntp.conf
svcadm enable svc:/network/ntp:default
```

### Windows

From a Windows machine:

```bash
net time /setsntp:my_time_server
net time /querysntp
net stop w32time && net start w32time
```

## FAQ

### How to change the timezone?

On Debian, there's the **tzconfig** command! It's very practical:

```bash
tzconfig
```

If this doesn't work, try:

```bash
dpkg-reconfigure tzdata
```

### "no server suitable for synchronization found"

If you encounter this error, there could be two reasons:

* Your NTP server needs to synchronize before it can provide time to other servers. This can sometimes take some time (~30 min).
* One of the servers in your list is unreachable, which makes your server unavailable for updates. Comment it out temporarily.

Once the following line appears in the logs ("clock is now synced"), everything is OK:

```
May 24 16:48:05 chronos ntpd[14262]: ntp engine ready
May 24 16:48:24 chronos ntpd[14262]: peer 91.121.45.45 now valid
May 24 16:48:24 chronos ntpd[14262]: peer 81.93.183.116 now valid
May 24 16:48:27 chronos ntpd[14262]: peer 94.23.220.143 now valid
May 24 17:27:13 chronos ntpd[14262]: clock is now synced
```
