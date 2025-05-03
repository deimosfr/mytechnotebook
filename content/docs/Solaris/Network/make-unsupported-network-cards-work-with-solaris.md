---
weight: 999
url: "/make-unsupported-network-cards-work-with-solaris/"
title: "Make Unsupported Network Cards Work with Solaris"
description: "How to make unsupported network cards work with Solaris by installing additional drivers"
categories: ["Solaris", "Hardware", "Networking"]
date: "2008-12-21T16:12:00+02:00"
lastmod: "2008-12-21T16:12:00+02:00"
tags: ["Solaris", "Network Cards", "Driver", "Hardware"]
toc: true
---

## Introduction

I recently faced an issue with network card recognition for my home server. I have two D-Link DGE-530T gigabit network cards, and unfortunately, they don't work out of the box with Solaris. No need to panic though - as the drivers exist under BSD license, they have been ported and some even packaged.

This was my case for these DLINK cards. I'll share the approach I followed to get my two cards working, and I'll provide reference links in case you have other network cards that need to be recognized.

## Installation

First, let's remove the old package containing the drivers:

```bash
pkgrm SK98sol
```

Next, we'll check if there are any driver aliases and then remove them if they exist:

```bash
grep sk98 /etc/driver_aliases
```

Remove all the lines that appear from this command.

Now, let's proceed with installing the driver. [Download](https://www.skd.de/e_en/support/driver_searchresults.html?navanchor=10013&term=typ.treiber+bs.SUN_Solaris+produkt.SK-9821V2.0&produkt=produkt.SK-9821V2.0&typ=typ.treiber&system=bs.SUN_Solaris) the version corresponding to your architecture of the D-Link driver, then decompress and install the package:

```bash
gtar -xzvf skge*.tar.Z
pkgadd -d . SKGEsol
```

Follow the instructions. At the end, update the list so that updated aliases are created:

```bash
update_drv -a -i "pci1186,4b01" skge
```

Finally, we need to make the hardware detected at boot:

```bash
touch /reconfigure
```

All that's left is to reboot the machine and you're good to go :-)

## FAQ

### I have problems with VirtualBox recognizing my card, why?

In the VirtualBox logs, you might find something like:

```
vboxflt:vboxNetFltSolarisOpenStream Failed to open '/dev/skge0' rc=19 pszName='skge0'
```

And that's the problem - in /dev/, there's only skge and not skge0 as indicated in /etc/hostname.skge0.

Unfortunately, I didn't find a quick solution. I preferred to change the network card instead. Sorry for those who were hoping to find a solution to this problem.

## References

- [List of network cards and additional drivers for Solaris](https://homepage2.nifty.com/mrym3/taiyodo/eng/)
- [https://www.sun.com/bigadmin/hcl/data/components/details/2729.html](https://www.sun.com/bigadmin/hcl/data/components/details/2729.html)
- [https://www.skd.de/e_en/support/driver_searchresults.html?navanchor=10013&term=typ.treiber+bs.SUN_Solaris+produkt.SK-9821V2.0&produkt=produkt.SK-9821V2.0&typ=typ.treiber&system=bs.SUN_Solaris](https://www.skd.de/e_en/support/driver_searchresults.html?navanchor=10013&term=typ.treiber+bs.SUN_Solaris+produkt.SK-9821V2.0&produkt=produkt.SK-9821V2.0&typ=typ.treiber&system=bs.SUN_Solaris)
