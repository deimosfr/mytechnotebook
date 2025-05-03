---
weight: 999
url: "/Gestion_des_packages_Solaris/"
title: "Solaris Package Management"
description: "A guide on how to manage packages in Solaris, including installation, verification, and removal of packages using various tools."
categories: ["Unix", "Solaris"]
date: "2012-02-15T18:06:00+02:00"
lastmod: "2012-02-15T18:06:00+02:00"
tags: ["Solaris", "Package Management", "System Administration", "Unix"]
toc: true
---

## Introduction

Packages, like in all distributions, are a simple way to install software. Here we will examine all the ways to manage these packages.

## Locations

To find out what is installed on your system, simply look at the `/var/sadm/install/contents` file:

```bash
more /var/sadm/install/contents
```

```
(output edited for brevity)
/bin=./usr/bin s none SUNWcsr
/dev d none 0755 root sys SUNWcsr SUNWcsd
/dev/allkmem=../devices/pseudo/mm@0:allkmem s none SUNWcsd
/dev/arp=../devices/pseudo/arp@0:arp s none SUNWcsd
/etc/ftpd/ftpusers e ftpusers 0644 root sys 198 16387 1094222536 SUNWftpr
/etc/passwd e passwd 0644 root sys 580 48298 1094222123 SUNWcsr
```

To find where a specific software is installed on the system, you can do this:

```bash
pkgchk -l -P showrev
```

```
Pathname: /usr/bin/showrev
Type: regular file
Expected mode: 0755
Expected owner: root
Expected group: sys
Expected file size (bytes): 29980
Expected sum(1) of contents: 57864
Expected last modification: Dec 14 06:17:58 AM 2004
Referenced by the following packages:
        SUNWadmc       
Current status: installed

Pathname: /usr/share/man/man1m/showrev.1m
Type: regular file
Expected mode: 0644
Expected owner: root
Expected group: root
Expected file size (bytes): 3507
Expected sum(1) of contents: 35841
Expected last modification: Dec 10 10:42:54 PM 2004
Referenced by the following packages:
        SUNWman        
Current status: installed
```

## Identification

When you have downloaded a package and want to check if it's for SPARC or x86, here's what to do:

First, verify that it's a package:

```bash
file SUNWrsc.pkg
```

```
SUNWrsc.pkg:	package datastream
```

Then display the package header:

```bash
head SUNWrsc.pkg
```

```
# PaCkAgE DaTaStReAm
SUNWrsc 1 3266
# end of header
SUNW_PRODVERS=2.2.1
SUNW_PKGVERS=1.0
PKG=SUNWrsc
NAME=Remote System Control
DESC=Sun Remote System Control system software
ARCH=sparc
VENDOR=Sun Microsystems, Inc.
```

## The Tools

{{< table "table-hover table-striped" >}}
| Tools | Descriptions |
|-------|-------------|
| pkgtrans | Transform packages from one format to another |
| pkgadd | Install a package to the system |
| pkgrm | Remove a package from the system |
| pkginfo | Display information about a package |
| pkgchk | Verify the installation state of a package |
{{< /table >}}

### pkgtrans

This will transform the system package into the "data stream" format:

```bash
pkgtrans /var/tmp /tmp/SUNWrsc.pkg SUNWrsc
```

You will get a pkg. If you want to do the reverse:

```bash
pkgtrans SUNWrsc.pkg .
```

### pkginfo

Here's an example:

```bash
pkginfo -l SUNWman
```

```
   PKGINST:  SUNWman
      NAME:  On-Line Manual Pages
  CATEGORY:  system
      ARCH:  sparc
   VERSION:  43.0,REV=67.0
   BASEDIR:  /usr
    VENDOR:  Sun Microsystems, Inc.
      DESC:  System Reference Manual Pages
    PSTAMP:  2004.09.01.17.00
  INSTDATE:  Sep 24 2004 12:32
   HOTLINE:  Please contact your local service provider
    STATUS:  completely installed
     FILES:    11383 installed pathnames
                   8 shared pathnames
                  97 directories
              119848 blocks used (approx)
```

### pkgadd

To install a specific package, do this:

```bash
pkgadd -d . SUNWvts
```

```
Processing package instance <SUNWvts> from
 </cdrom/sol_10_sparc_4/Solaris_10/ExtraValue/CoBundled/SunVTS_6.0/Packages>

SunVTS Framework(sparc) 6.0,REV=2004.08.18.12.00
Copyright 2004 Sun Microsystems, Inc.  All rights reserved.
Use is subject to license terms.
Using </opt> as the package base directory.
## Processing package information.
## Processing system information.
## Verifying package dependencies.
## Verifying disk space requirements.
## Checking for conflicts with packages already installed.
## Checking for setuid/setgid programs.

This package contains scripts which will be executed with super-user
permission during the process of installing this package.

Do you want to continue with the installation of <SUNWvts> [y,n,?] y

Installing SunVTS Framework as <SUNWvts>

## Installing part 1 of 1.
9213 blocks

Installation of <SUNWvts> was successful.
```

To install all packages in data stream format:

```bash
pkgadd -d /tmp/SUNWrsc.pkg all
```

```
Processing package instance <SUNWrsc> from </tmp/SUNWrsc.pkg>

Remote System Control(sparc) 2.2.1,REV=2002.02.11
Copyright 2001 Sun Microsystems, Inc. All rights reserved.
Using </> as the package base directory.
## Processing package information.
## Processing system information.
   15 package pathnames are already properly installed.
## Verifying disk space requirements.
## Checking for conflicts with packages already installed.
## Checking for setuid/setgid programs.

Installing Remote System Control as <SUNWrsc>

## Installing part 1 of 1.
10499 blocks

Installation of <SUNWrsc> was successful.
```

If the package is on a website:

```bash
pkgadd -d http://instructor/packages/SUNWrsc.pkg all
```

```
## Downloading...
..............25%..............50%..............75%..............100%
## Download Complete


Processing package instance <SUNWrsc> from 
<http://instructor/packages/SUNWrsc.pkg>

Remote System Control(sparc) 2.2.1,REV=2002.02.11
Copyright 2001 Sun Microsystems, Inc. All rights reserved.
Using </> as the package base directory.
## Processing package information.
## Processing system information.
   15 package pathnames are already properly installed.
## Verifying disk space requirements.
## Checking for conflicts with packages already installed.
## Checking for setuid/setgid programs.

Installing Remote System Control as <SUNWrsc>

## Installing part 1 of 1.
10499 blocks

Installation of <SUNWrsc> was successful.
```

#### Spool

The spool is where packages go (`/var/spool/pkg`). If, for example, on a Sun CD, we want to install a package and also have it copied to the spool directory, here's an example:

```bash
kgadd -d /cdrom/cdrom0/s0/Solaris_10/Product -s spool SUNWauda
```

```
Transferring <SUNWauda> package instance
```

Let's check:

```bash
ls -al /var/spool/pkg
```

```
total 6
drwxrwxrwt   3 root     bin          512 Oct  1 14:26 .
drwxr-xr-x  12 root     bin          512 Sep 30 20:03 ..
drwxrwxr-x   5 root     root         512 Oct  1 14:26 SUNWauda
```

To install it later:

```bash
pkgadd SUNWauda
```

### pkgchk

This command allows you to verify if a package is properly installed (path, checksum...).

* To list the contents of a package:

```bash
pkgchk -v SUNWladm
```

```
/usr
/usr/sadm
/usr/sadm/lib
/usr/sadm/lib/localeadm
/usr/sadm/lib/localeadm/Locale_config_S10.txt
/usr/sadm/lib/localeadm/admin
/usr/sbin
/usr/sbin/localeadm
```

* To check if a file has changed from its original state (in the package):

```bash
pkgchk -p /etc/shadow
```

```
ERROR: /etc/shadow
    modtime <09/03/04 03:35:24 PM> expected <09/30/04 08:06:14 PM> actual
    file size <296> expected <309> actual
    file cksum <20180> expected <21288> actual
```

* The -l option lists the information about the package contents:

```bash
pkgchk -l -p /usr/bin/showrev
```

```
Pathname: /usr/bin/showrev
Type: regular file
Expected mode: 0755
Expected owner: root
Expected group: sys
Expected file size (bytes): 29656
Expected sum(1) of contents: 31261
Expected last modification: Sep 02 09:21:11 2004
Referenced by the following packages:
       SUNWadmc       
Current status: installed
```

### pkgrm

To remove a package:

```bash
pkgrm SUNWapchr
```

```
The following package is currently installed:
   SUNWapchr       Apache Web Server (root)
                   (sparc) 11.10.0,REV=2004.08.20.02.37

Do you want to remove this package? [y,n,?,q] y

## Removing installed package instance <SUNWapchr>
## Verifying package dependencies.
WARNING:
    The <SUNWapchu> package depends on the package
    currently being removed.
WARNING:
    The <SUNWapchd> package depends on the package
    currently being removed.
WARNING:
    The <SUNWipplr> package depends on the package
    currently being removed.
WARNING:
    The <SUNWserweb> package depends on the package
    currently being removed.
Dependency checking failed.

Do you want to continue with the removal of this package [y,n,?,q] y
## Processing package information.
## Removing pathnames in class <initd>
/etc/rcS.d/K16apache
/etc/rc3.d/S50apache
/etc/rc2.d/K16apache

(output ommited for brevity)

/etc/apache/httpd.conf-example
/etc/apache/README.Solaris
/etc/apache <shared pathname not removed>
/etc <shared pathname not removed>
## Updating system information.

Removal of <SUNWapchr> was successful.
```

#### Spool

To remove a package and delete it from the spool:

```bash
pkgrm -s spool SUNWauda
```

```
The following package is currently spooled:
   SUNWauda        Audio Applications
                   (sparc) 11.10.0,REV=2004.09.03.08.15

Do you want to remove this package? [y,n,?,q] y

Removing spooled package instance <SUNWauda>
```

### pkgtrans

To transform packages into a stream package, it's quite simple:

```bash
pkgtrans -s  Product  /var/tmp/stream.pkg SUNWzlib SUNWftpr SUNWftpu
```

```
Transferring <SUNWzlib> package instance
Transferring <SUNWftpr> package instance
Transferring <SUNWftpu> package instance
```

Let's verify:
```bash
file /var/tmp/stream.pkg
```

Let's look at the header:
```bash
head -5 /var/tmp/stream.pkg
```

```
# PaCkAgE DaTaStReAm
SUNWzlib 1 186
SUNWftpr 1 70
SUNWftpu 1 300
# end of header
```

Now, for installation:

```bash
pkgadd -d /var/tmp/stream.pkg
```

```
The following packages are available:
  1  SUNWftpr     FTP Server, (Root)
                  (sparc) 11.10.0,REV=2004.12.11.01.30
  2  SUNWftpu     FTP Server, (Usr)
                  (sparc) 11.10.0,REV=2004.12.11.01.30
  3  SUNWzlib     The Zip compression library
                  (sparc) 11.10.0,REV=2004.12.10.05.25

Select package(s) you wish to process (or 'all' to process
all packages). (default: all) [?,??,q]: q
```

### prodreg

This provides a graphical interface for managing packages, similar to Solaris installation:

![Prodreg](/images/prodreg.avif)

## Package Locations

To avoid risking damage to the main system, packages are installed in specific locations:

{{< table "table-hover table-striped" >}}
| Files or Folders | Description |
|------------------|-------------|
| `/var/sadm/install/contents` | List of all system packages |
| `/opt/pkgname` | Path for most installed packages |
| `/opt/pkgname/bin` or `/opt/bin` | Binaries for most installed packages |
| `/var/opt/pkgname` or `/etc/opt/pkgname` | Logs for most installed packages |
{{< /table >}}
