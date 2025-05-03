---
weight: 999
url: "/Modifier_les_paramètres_des_cartes_réseaux/"
title: "Modifying Network Card Parameters"
description: "A guide to modifying network card parameters in Solaris using the dladm utility, including how to view and modify link properties."
categories: ["Solaris", "Network"]
date: "2009-11-28T16:20:00+02:00"
lastmod: "2009-11-28T16:20:00+02:00"
tags: ["Solaris", "Network", "dladm", "Configuration"]
toc: true
---

## Introduction

Project Brussels from the OpenSolaris project revamped how link properties are managed, and their push to get rid of ndd and device-specific properties is now well underway!

## Show properties

Link properties are actually pretty cool, and they can be displayed with the dladm utilities "show-linkprop" option:

```bash
$ dladm show-linkprop e1000g0

LINK         PROPERTY        PERM VALUE          DEFAULT        POSSIBLE
e1000g0      speed           r-   0              0              --
e1000g0      autopush        --   --             --             --
e1000g0      zone            rw   --             --             --
e1000g0      duplex          r-   half           half           half,full
e1000g0      state           r-   down           up             up,down
e1000g0      adv_autoneg_cap rw   1              1              1,0
e1000g0      mtu             rw   1500           1500           --
e1000g0      flowctrl        rw   bi             bi             no,tx,rx,bi
e1000g0      adv_1000fdx_cap r-   1              1              1,0
e1000g0      en_1000fdx_cap  rw   1              1              1,0
e1000g0      adv_1000hdx_cap r-   0              1              1,0
e1000g0      en_1000hdx_cap  r-   0              1              1,0
e1000g0      adv_100fdx_cap  r-   1              1              1,0
e1000g0      en_100fdx_cap   rw   1              1              1,0
e1000g0      adv_100hdx_cap  r-   1              1              1,0
e1000g0      en_100hdx_cap   rw   1              1              1,0
e1000g0      adv_10fdx_cap   r-   1              1              1,0
e1000g0      en_10fdx_cap    rw   1              1              1,0
e1000g0      adv_10hdx_cap   r-   1              1              1,0
e1000g0      en_10hdx_cap    rw   1              1              1,0
e1000g0      maxbw           rw   --             --             --
e1000g0      cpus            rw   --             --             --
e1000g0      priority        rw   high           high           low,medium,high
```

As you can see in the above output, the typical speed, duplex, mtu and flowctrl properties are listed. In addition to those, the "maxbw" and "cpus" properties that were introduced with the recent crossbow putback are visible. The "maxbw" property is especially useful, since it allows you to limit how much bandwidth is available to an interface. Here is an example that caps bandwidth for an interface at 2Mb/s:

```bash
dladm set-linkprop -p maxbw=2m e1000g0
```

To see how this operates, you can use your favorite data transfer client:

```bash
$ scp techtalk1* 192.168.1.10:
Password:
techtalk1.mp3 5% 2128KB 147.0KB/s 04:08 ETA
```

## Modify properties

The read/write link properties can be changed on the fly with dladm, so increasing the "maxbw" property will allow the interface to consume additional bandwidth:

```bash
$ dladm set-linkprop -p maxbw=10m e1000g0
```

Once the bandwidth is increased, you can immediately see this reflected in the data transfer progress:

```bash
techtalk1.mp3 45% 17MB 555.3KB/s 00:38 ETA
```

Clearview rocks, and it's awesome to see that link properties are going to be managed in a standard uniform way going forward! Nice!

## Resources
- [https://prefetch.net/blog/index.php/2009/03/15/viewing-network-device-properties-on-solaris-hosts/](https://prefetch.net/blog/index.php/2009/03/15/viewing-network-device-properties-on-solaris-hosts/)
