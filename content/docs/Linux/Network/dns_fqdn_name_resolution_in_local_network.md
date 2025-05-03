---
weight: 999
url: "/Résolution_noms_DNS_FQDN_en_local/"
title: "DNS FQDN Name Resolution in Local Network"
description: "How to solve DNS FQDN name resolution issues on Ubuntu systems by modifying the nsswitch.conf file."
categories: ["Linux", "Ubuntu", "Network"]
date: "2008-12-17T23:26:00+02:00"
lastmod: "2008-12-17T23:26:00+02:00"
tags: ["DNS", "Ubuntu", "Network", "Configuration", "Troubleshooting"]
toc: true
---

## Problem

You may have noticed that on Ubuntu, there is a problem resolving FQDN (Fully Qualified Domain Names) for your local DNS servers.

## Solution

Simply edit the `/etc/nsswitch.conf` file and replace this line:

```bash
hosts:          files mdns4_minimal [NOTFOUND=return] dns mdns4
```

With:

```bash
hosts:          files dns mdns4
```

And just like magic, everything works! :-)
