---
weight: 999
url: "/Modifier_les_serveurs_de_synchronisation_NTP_sous_Windows/"
title: "Modifying NTP Synchronization Servers on Windows"
description: "How to configure NTP synchronization servers on Windows systems using w32tm commands"
categories: ["Windows", "Server"]
date: "2009-10-30T07:34:00+02:00"
lastmod: "2009-10-30T07:34:00+02:00"
tags: ["Windows", "NTP", "Server", "Time", "Synchronization"]
toc: true
---

## Introduction

I have a Windows Server 2008 DC and I wanted to use my internal time server on a Linux box running ntpd.

After a little hunting around, I found the command required to set Windows up to use the correct time peer.

## Configuration

```bash
w32tm /config /update /manualpeerlist:"0.pool.ntp.org,0x8 1.pool.ntp.org,0x8" /syncfromflags:MANUAL
```

After making this change, you need to restart the Windows Time Service by issuing the following 2 commands:

```bash
net stop w32time
net start w32time
```

If you have problems, first make sure the Windows Time Service is enabled.

Note: This works with Windows XP, Windows Vista, Windows Server 2003 and Windows Server 2008.

## Resources
- http://blogs.msdn.com/w32time/default.aspx
- http://www.meinberg.de/german/sw/ntp.htm
