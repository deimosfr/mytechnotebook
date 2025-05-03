---
weight: 999
url: "/Mise_en_place_de_Nagios_NRPE_sur_Solaris_10/"
title: "Setting up Nagios NRPE on Solaris 10"
description: "A guide on how to set up Nagios NRPE on Solaris 10 with detailed steps for installation and configuration."
categories: ["Linux", "Servers", "Network"]
date: "2008-04-05T10:22:00+02:00"
lastmod: "2008-04-05T10:22:00+02:00"
tags: ["Solaris", "Nagios", "NRPE", "Monitoring", "Unix"]
toc: true
---

## Introduction

Traditionally Linux is the platform of choice for running Nagios however Solaris takes some beating for performance and reliability. Installing Nagios on Solaris is pretty straightforward however there are some tricks to make the install smoother.

Usually our installation comprises Nagios 2.0, Nagios Plugins 1.4.2 and NRPE 2.4 so we'll cover these here.

Our target system is a Sun Ultra 60 so we're using the SPARC version of Solaris 10. The same principles should apply to Solaris on Intel though.

## Installing Solaris 10

We started off with a vanilla install of Solaris 10 (01/06) and opted to install all packages except OEM. You can get away with less if you're running the system headless.

Many of the development tools and libraries we'll need for Nagios come on the Solaris Companion DVD so the next step was to install those. Again we opted for the standard list of components.

Our Ultra 60 only has a CD-ROM. To get around this we used NFS to mount the DVD from another server.

## Preparing Solaris

In order to make the Nagios install go smoothly we performed some tweaks. These are made from the command shell.

### Setting default shell to Bash

```bash
$ usermod -s /usr/bin/bash
```

When you log in you'll now get the Bash shell.

If you have a strong preference for another shell then free to use that instead :-)

### Correcting Path

The default path misses packages installed from the Solaris Companion DVD, we can fix this by typing the following command:

```bash
$ export PATH=/usr/sfw/bin:/usr/ccs/bin:/opt/sfw/bin/:/usr/bin:/usr/ucb:/usr/sbin
```

### Setting Compiler flags

These are optional but will make Nagios run faster on your system. Setting -mcpu and -mtune to 'ultrasparc' should be fine for most Sun UltraSPARC based systems. Using -pipe speeds up compilation but may cause problems if you are tight on memory.

```bash
$ export CFLAGS="-O3 -pipe -mcpu=ultrasparc -mtune=ultrasparc"
```

More information can be found at gcc.gnu.org.

### Making changes stick

We to make our changes permenant we appended the following lines to .profile in our home directory:

```bash
PATH=/usr/sfw/bin:/usr/ccs/bin:/opt/sfw/bin/:/usr/bin:/usr/ucb:/usr/sbin
CFLAGS="-O3 -pipe -mcpu=ultrasparc -mtune=ultrasparc"
export PATH CFLAGS
```

## Nagios Prerequisites

### GD Graphics Libraries

We found it was necessary to install the GD Graphics libraries for Nagios to compile, we used version 2.0.33. A standard 'configure', 'make' and 'make install' was sufficient, remember to install as root (or configure sudo). Software can be found at [https://www.boutell.com/gd/](https://www.boutell.com/gd/)

### NRPE

We found that NRPE couldn't find the required SSL libraries to compile so we downloaded the package 'openssl-0.9.8a-sol10-sparc-local.gz' from Sun Freeware - [https://www.sunfreeware.com/](https://www.sunfreeware.com/)

#### Installation steps

```bash
$ gzip -d openssl-0.9.8a-sol10-sparc-local.gz
$ su -
# /usr/sbin/pkgadd -d openssl-0.9.8a-sol10-sparc-local
$ exit
```

#### Compilation and Installation

Both Nagios and Nagios Plugins compiled and installed as normal. NRPE required the following configuration parameters to correctly build:

```bash
$ ./configure --with-ssl-lib=/usr/local/ssl/lib/
```

Now your system is ready to configure...
