---
title: "PowerDNS: Creating a DNS Cache Server"
slug: powerdns-creating-a-dns-cache-server/
description: "How to install and configure PowerDNS as a caching DNS server on Debian 6"
categories: ["Debian", "Linux", "Network", "Servers"]
date: "2012-05-15T14:46:00+02:00"
lastmod: "2012-05-15T14:46:00+02:00"
tags: ["PowerDNS", "DNS", "Caching", "Debian", "Network"]
---

![PowerDNS](../../../static/images/powerdns_logo.avif)


|||
|-|-|
| **Software version** | 2.9.22 |
| **Operating System** | Debian 6 |
| **Website** | [PowerDNS Website](https://www.powerdns.com) |


## Introduction

PowerDNS is (as its name suggests) a DNS server. It's a direct competitor to Bind. It aims to be less RAM-intensive and offers more flexible configuration options than Bind.

PowerDNS is divided into several roles:
- Master
- Cache

Here we will cover the cache aspect. If you want to set up a PowerDNS master server, I invite you to [follow this link](./powerdns_creating_a_master_dns_server.md).

## Installation

To install PowerDNS:

```bash
aptitude install pdns-recursor
```

## Configuration

Once installed, the cache server is functional for the local server. All you need to do is configure the listening address to enable it for the rest of your network:

``` bash hl_lines="5 9"
[...]
#################################
# allow-from    If set, only allow these comma separated netmasks to recurse
#
allow-from=
[...]
#################################
# local-address IP addresses to listen on, separated by spaces or commas. Also accepts ports.
#
local-address=0.0.0.0
```

And restart the service to activate it:

```bash
/etc/init.d/pdns-recursor restart
```

Now all you need to do is point your machines to this new server :-)

## References

* http://www.adminsehow.com/2009/05/how-to-install-a-caching-only-dns-server-using-powerdns-on-debian-lenny/
