---
weight: 999
url: "/MacFuse_+_NTFS-3G_\\:_Lecture_et_écriture_de_partitions_NTFS_sur_Mac_OS_X/"
title: "MacFuse + NTFS-3G: Reading and Writing NTFS Partitions on Mac OS X"
description: "Tutorial for installing and configuring MacFuse and NTFS-3G to enable read and write access to NTFS partitions on Mac OS X"
categories: ["MacOS"]
date: "2007-10-31T06:24:00+02:00" 
lastmod: "2007-10-31T06:24:00+02:00"
tags: ["MacFuse", "NTFS", "MacOS X", "Disk Management", "Filesystem"]
toc: true
---

## Introduction

Even with the release of Leopard (10.5), everyone thought they would have NTFS write support. We had already missed out on user-friendly native ZFS...

In short, if you want to be able to write to NTFS partitions, you need to install MacFuse and [NTFS-3G](https://www.ntfs-3g.org/). Here's how to proceed.

## Prerequisites

* MacFuse: http://code.google.com/p/macfuse/
* MacPorts: http://www.macports.org/
* X11: Mac OS X DVD
* XCode: Mac OS X DVD

## Installation

It's fairly simple - download the MacPorts and MacFuse packages and install them. For the rest, everything is on the Leopard DVD.

## Configuration

### MacPorts

If this is your first MacPorts installation, run this command:

```bash
export PATH=/opt/local/bin:/opt/local/sbin:$PATH
```

Then we'll install what we need. But first, let's get the MacPorts package list:

```bash
sudo port -d selfupdate
```

Next, install pkgconfig and ntfs-3g:

```bash
sudo port install pkgconfig ntfs-3g
```

## Usage

If you have Bootcamp installed or if your NTFS partition is already mounted, check the corresponding device using the "df" command:

```bash
mac% df
Filesystem    512-blocks     Used Available Capacity  Mounted on
/dev/disk0s2   127664128 60996656  66155472    48%    /
devfs                212      212         0   100%    /dev
fdesc                  2        2         0   100%    /dev
map -hosts             0        0         0   100%    /net
map auto_home          0        0         0   100%    /home
/dev/disk0s3    67035608 28905504  38130104    44%    /Volumes/Untitled
```

Here, `/dev/disk0s3` corresponds to Windows. So I'll unmount the partition:

```bash
sudo umount /Volumes/Untitled
```

Next, I need to create a folder called Vista for example in Volumes, then mount my device in this folder:

```bash
sudo mkdir /Volumes/Vista
sudo ntfs-3g /dev/disk0s3 /Volumes/Vista -o ping_diskarb,volname="Vista"
```

Then if you run df, it should appear and you can now access it :-)

```bash
mac% df
Filesystem    512-blocks     Used Available Capacity  Mounted on
/dev/disk0s2   127664128 60997024  66155104    48%    /
devfs                223      223         0   100%    /dev
fdesc                  2        2         0   100%    /dev
map -hosts             0        0         0   100%    /net
map auto_home          0        0         0   100%    /home
/dev/disk0s3    67035608 28905504  38130104    44%    /Volumes/Vista
```
