---
weight: 999
url: "/Resetter_completement_la_configuration_d'un_Cisco/"
title: "Completely Reset a Cisco Configuration"
description: "How to completely reset a Cisco device configuration with simple commands"
categories:
  - Network
date: "2007-05-23T13:40:00+02:00" 
lastmod: "2007-05-23T13:40:00+02:00"
tags:
  - Cisco
  - Network
  - Servers
toc: true
---

## Reset a Cisco PIX Configuration

To completely reset the configuration of a Cisco PIX, here is the solution:

```bash
clear config all
```

```bash
write erase
```

The installation can then restart:

```bash
  -----------------------------------------------------------------------
                               ||        ||
                               ||        ||
                              ||||      ||||
                          ..:||||||:..:||||||:..
                         c i s c o S y s t e m s
                        Private Internet eXchange
  -----------------------------------------------------------------------
                        Cisco PIX Firewall

Cisco PIX Firewall Version 6.3(5)
Licensed Features:
Failover:                    Enabled
VPN-DES:                     Enabled
VPN-3DES-AES:                Disabled
Maximum Physical Interfaces: 6
Maximum Interfaces:          10
Cut-through Proxy:           Enabled
Guards:                      Enabled
URL-filtering:               Enabled
Inside Hosts:                Unlimited
Throughput:                  Unlimited
IKE peers:                   Unlimited

This PIX has a Failover Only (FO) license.
```
