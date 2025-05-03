---
weight: 999
url: "/ZFS_On_Linux_\\:_Mise_en_place_de_ZFS_sous_Linux/"
title: "ZFS On Linux: Setting up ZFS on Linux"
description: "A guide to installing and configuring ZFS filesystem on Linux distributions, including Ubuntu and Debian, with instructions for installation, configuration, and troubleshooting common issues."
categories: ["Linux", "FreeBSD", "Ubuntu"]
date: "2011-12-06T20:57:00+02:00"
lastmod: "2011-12-06T20:57:00+02:00"
tags: ["ZFS", "Linux", "Ubuntu", "Debian", "Filesystem"]
toc: true
---

![ZFS on Linux Logo](/images/zfs-linux.avif)

## Introduction

If like me you're a fan of this filesystem and find it a shame that it's not natively available on Linux, there are currently several solutions to have this filesystem:

- [Solaris](https://www.oracle.com)/OpenSolaris: This is where ZFS comes from, but it remains a proprietary OS
- [FreeBSD](https://www.freebsd.org/): The first port of ZFS appeared on FreeBSD, but we're looking to use Linux here
- [Kfreebsd](https://www.debian.org/ports/kfreebsd-gnu/): not really Linux (although Debian), but a FreeBSD kernel that allows ZFS to run with a Debian-style layer on top
- [ZFS on Fuse](https://zfs-fuse.net/): works on Linux, slow (because it runs on FUSE) but historically the first to be released for Linux (so supposedly the most mature)
- [ZFS on Linux](https://zfsonlinux.org): newer, but has the advantage of running as a Linux kernel module

I chose this last solution because I wanted to keep a Linux machine (Debian/Ubuntu) and have ZFS.

## Installation

### Prerequisites

To install ZFS on Linux, we'll need some dependencies:

```bash
aptitude install build-essential gawk alien fakeroot linux-headers-$(uname -r) install zlib1g-dev uuid-dev libblkid-dev libselinux1-dev
```

Once these dependencies are installed, you'll need to either get the list of packages and install them one by one, or use the Ubuntu repository (which is what we'll do here).

### Ubuntu

If you're on Ubuntu, run this command to add the repository:

```bash
add-apt-repository ppa:zfs-native/stable
```

If it doesn't work, you might be missing this package:

```bash
apt-get install python-software-properties
```

### Debian

For Debian, you'll need to add the following sources to your sources.list file (`/etc/apt/sources.list`):

```bash
deb http://ppa.launchpad.net/dajhorn/zfs/ubuntu natty main 
deb-src http://ppa.launchpad.net/dajhorn/zfs/ubuntu natty main 
```

### ZFS

To install ZFS, all you need to do is:

```bash
aptitude update
aptitude install ubuntu-zfs
```

And that's it, it's installed! I'll let you check the references of this page for how to use ZFS.

## FAQ

### My server crashed because there was no more RAM available

I already talk about this [here](./zfs_:_le_filesystem_par_excellence.html), it's because of the [ZFS cache](./zfs_:_le_filesystem_par_excellence.html#le_cache_arc_zfs) that needs to be customized. Here we'll set it to 512 MB (`/etc/modprobe.d/zfs.conf`):

```bash
options zfs zfs_arc_max=536870912
```

## References

http://www.deimos.fr/blocnotesinfo/index.php?title=ZFS_:_Le_FileSystem_par_excellence
