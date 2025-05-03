---
weight: 999
url: "/Utilisation_avancé_de_Bind/"
title: "Advanced usage of Bind"
description: "This guide explains advanced techniques for managing BIND DNS servers including zone updates, round-robin load balancing, and troubleshooting common issues."
categories: 
  - Linux
date: "2013-05-07T11:33:00+02:00"
lastmod: "2013-05-07T11:33:00+02:00"
tags:
  - Network
  - Servers
  - DNS
toc: true
---

## Introduction

Bind is good, but sometimes it becomes a bit complex. Especially when you want to manage DNS servers like providers do. Here are some tips I found to improve the beast.

## Force zone updates

### Method 1

You may quickly need to update your zones without waiting for Bind to do it itself (see SOA, refresh, etc... for each zone). For this, you must freeze Bind updates on the zones in question:

```bash
rndc freeze deimos.fr in internalview
```

* deimos.fr: the zone
* internalview: the view in which the zone is located

You can remove "in <view>" if you don't have a defined view. If the action went well, you'll find this in the logs:

```
Oct  2 19:05:42 star1 named[8403]: freezing zone 'deimos.fr/IN' internalview: success
```

Now, make the changes you want on your zone file.

Once finished, run these commands:

```bash
rndc reload deimos.fr in internalview
rndc thaw deimos.fr in internalview
```

* reload: reload configuration file and zones.
* thaw: Enable updates to a frozen dynamic zone and reload it.

If all goes well, you'll see something like this in the logs:

```
Oct  2 19:09:50 star1 named[8403]: zone deimos.fr/IN/internalview: loaded serial 2008100208
Oct  2 19:09:50 star1 named[8403]: zone deimos.fr/IN/internalview: sending notifies (serial 2008100208)
Oct  2 19:09:50 star1 named[8403]: client 192.168.0.27#47874: view internalview: transfer of 'deimos.fr/IN': AXFR-style IXFR started
Oct  2 19:09:50 star1 named[8403]: client 192.168.0.27#47874: view internalview: transfer of 'deimos.fr/IN': AXFR-style IXFR ended
Oct  2 19:09:54 star1 named[8403]: unfreezing zone 'deimos.fr/IN' internalview: success
```

### Method 2

Here is a second method, to do the same thing:

```bash
rndc retransfer deimos.fr
```

## Round Robin: Load balancer

First of all, you should know that you can only load balance with A records (except apparently for BIND v4 where you can do it with CNAME). Here's an example of how to manage a domain configuration:

```bash
; Round Robin / Load Balancing
www    60   IN  A      x.x.x.1
www    60   IN  A      x.x.x.2
```

* 60: corresponds to the TTL, and it is very important because it will decide when to switch
* x.x.x.x: IP addresses of servers. Unfortunately, you cannot use DNS names.

You can also choose the type of load balancing among these:

* fixed - records are returned in the order they are defined in the zone file
* random - records are returned in a random order
* cyclic - records are returned in a round-robin fashion

To do this, add this type of lines and adapt according to your needs:

```bash
...
rrset-order { order cyclic; };
...
```

For more information, see the reference sites below.

## FAQ

### journal rollforward failed: journal out of sync with zone

If you have this error, it's due to a zone synchronization problem. Look in your logs for the zone(s) causing problems, then delete the sync file:

```bash
rm /etc/bind/db.deimos.fr.jnl
```

All that's left is to reload or restart your bind and the sync will restart correctly.

## Resources
- http://jon.netdork.net/2008/08/21/bind-dynamic-zones-and-updates
- http://www.zytrax.com/books/dns/
- http://www.zytrax.com/books/dns/ch7/queries.html#rrset-order
- http://www.zytrax.com/books/dns/info/ttl.html
- http://www.zytrax.com/books/dns/ch9/rr.html
