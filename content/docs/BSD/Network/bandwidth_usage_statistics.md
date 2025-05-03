---
weight: 999
url: "/Statistiques_sur_la_bande_passante_occupée/"
title: "Bandwidth Usage Statistics"
description: "A simple shell script to calculate and monitor bandwidth usage on external network interfaces for BSD and Linux systems."
categories:
  - Linux
date: "2007-10-07T10:25:00+02:00"
lastmod: "2007-10-07T10:25:00+02:00"
tags:
  - Linux
  - Network
  - Servers
  - BSD
toc: true
---

## Introduction

Here is a simple shell script to calculate the bandwidth usage on the external interface of a BSD or Linux box. Netstat bandwidth summary works well on OpenBSD 4.1, but colleagues have mentioned 3.9 may not work. Linux should work without issue. Also, remember that the netstat stats will reset on reboot of the box.

One could use this script to keep track of Internet bandwidth in case their ISP accused them of using too much bandwidth. Comcast for example will call foul if you use more than 90 to 150 gigabytes of download per month. We can only guess that the upload limit is the same. Verizon says they do not have a limit, but they will contact bandwidth abusers. Your ISP might have different rules so check with them. Then use this simple tool to make sure you know what you are using.

This is what the report of my system looks like...

```bash
  External interface bandwidth usage:
   uptime           16 days
   ExtIf in total   13 GBytes
   ExtIf out total  16 GBytes
   ExtIf in/day     831 MBytes/day
   ExtIf out/day    986 MBytes/day
   ExtIf in/30day   24 GBytes/month
   ExtIf out/30day  29 GBytes/month
```

## Script

You could put the executable line into `/etc/daily` on the 13th line. This way you will get an email in the "daily output" email the BSD box sends and includes the above stats.

Here is the script called "calomel_interface_stats.sh":

```bash
#!/usr/local/bin/bash

SECS=`uptime | awk '{print $3}'`
EXT_IN=`netstat -I em0 -b | tail -1 | awk '{print $5}'`
EXT_OUT=`netstat -I em0 -b | tail -1 | awk '{print $6}'`

echo " "
echo "External interface bandwidth usage:"
echo " uptime          " $(($SECS/86400)) "days"
echo " ExtIf in total  " $(($EXT_IN/1000033000)) "GBytes"
echo " ExtIf out total " $(($EXT_OUT/1000033000)) "GBytes"
echo " ExtIf in/day    " $(($EXT_IN*86400/SECS/1000033)) "MBytes/day"
echo " ExtIf out/day   " $(($EXT_OUT*86400/SECS/1000033)) "MBytes/day"
echo " ExtIf in/30day  " $(($EXT_IN*86400*30/SECS/1000033000)) "GBytes/month"
echo " ExtIf out/30day " $(($EXT_OUT*86400*30/SECS/1000033000)) "GBytes/month"
```
