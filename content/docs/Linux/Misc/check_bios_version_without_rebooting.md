---
weight: 999
url: "/Connaitre_sa_version_du_Bios_sans_rebooter/"
title: "How to Check BIOS Version Without Rebooting"
description: "A simple method to check your system's BIOS version from Linux without having to reboot"
categories: ["Linux", "Hardware", "System Administration"] 
date: "2007-08-29T19:46:00+02:00"
lastmod: "2007-08-29T19:46:00+02:00"
tags: ["BIOS", "hardware", "system", "Linux"]
toc: true
---

Sometimes when you have old machines, you might want to know if flashing the BIOS would allow you to use a larger disk. From a remote system, the following command (run as root) will give you plenty of information:

```bash
dd if=/dev/mem bs=32k skip=31 count=1 | strings -n 8 | grep -i bios
1+0 records in
1+0 records out
32768 bytes transferred in 0.011551 seconds (2836813 bytes/sec)
Award SoftwareIBM COMPATIBLE 486 BIOS COPYRIGHT Award Software Inc.oftware Inc. Aw
Award Modular BIOS v4.51PG
```

## Explanation

On an x86 system, the BIOS is traditionally accessible in the last 64KB of the first MB of memory. The dd command is instructed to read from RAM starting at the first byte, skipping 31 blocks of 32KB each, and displaying the 32nd block.

This tip was found on comp.os.linux.misc (from marcello).

Another example looking at the last 32KB of the first MB of memory:

```bash
sudo dd if=/dev/mem bs=32k skip=30 count=1 | strings -n 8 | grep -i bios
1+0 records read
1+0 records written
32768 bytes (33 kB) copied, 0.000149 seconds, 220 MB/s
Phoenix cME FirstBIOS Notebook Pro
```
