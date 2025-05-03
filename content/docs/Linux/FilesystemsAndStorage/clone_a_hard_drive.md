---
weight: 999
url: "/Cloner_un_disque_dur/"
title: "Clone a Hard Drive"
description: "Various methods to clone hard drives in Linux systems including dd, cat, partimage, and over network connections"
categories: ["Linux", "System Administration", "Storage"]
date: "2013-05-07T09:10:00+02:00"
lastmod: "2013-05-07T09:10:00+02:00"
tags: ["clone", "disk", "hard drive", "dd", "partimage", "backup"]
toc: true
---

## Introduction

To clone a disk under Windows, you need to pull out all the tools, and if it's bootable, then hold on tight... how much does it cost? Obviously, Windows purist pirates will tell me they download a cracked version.

Why use such tools and break the law when free and amazing tools exist? Let's take a look at some options...

## Solutions

### dd

dd is the ultimate solution. To duplicate a disk with a progress bar:

```bash
dd if=/dev/sda2 of=/dev/sdb2 bs=4096 conv=notrunc,noerror | bar -s 500g
```

Here I'm copying a 500g hard drive. For those who don't want to use the bar command:

```bash
dd if=/dev/sda2 of=/dev/sdb2 bs=4096 conv=notrunc,noerror &
watch -n5 -- pkill -USR1 ^dd$
```

To clone a disk remotely:

```bash
dd if=/dev/vgname/lvname bs=1M | ssh root@new-server 'dd of=/dev/vgname/lvname bs=1M'
```

#### The Partition Table

* Method 1

You can, if you wish, simply back up the partition table / MBR (sector 0):

```bash
dd if=/dev/sda of=~/sda.sector0 count=1
```

Then to restore:

```bash
dd if=~/sda.sector0 of=/dev/sda count=1
```

* Method 2

Here's another method to save the partition table:

```bash
sfdisk -d /dev/sda > ~/sda.ptbl
```

And to restore it:

```bash
sfdisk /dev/sda < ~/sda.ptbl
```

#### Across the Network

##### Netcat

On the target machine:

```bash
nc -l -p 1234 | dd of=/dev/sda1 bs=4k
```

On the source server:

```bash
dd if=/dev/sda1 bs=4k | nc 1234
```

##### SSH

To do a dd via SSH:

```bash
dd if=/dev/sda1 | ssh user@destination-srv 'dd of=/dev/sda1'
```

### Cat

Here's the simplest solution for copying an entire disk (partitions, boot sectors...):

```bash
cat /dev/hdx > /dev/hdy
```

hdx: the source disk  
hdy: the destination disk

### Partimage

There's the wonderful [Part image](https://www.partimage.org/Index.fr.html) software which also allows cloning and even over the network into a disk image :-)

## Verifying Disk Integrity

Once the disk is cloned, it's best to check the integrity of its data (force a check disk at reboot), for example, in ext3:

```bash
touch /forcefsck
```

Then restart the machine and at the next boot, it will force the check. After that, you can use it without issues.

### CloneZilla

CloneZilla works a bit like Symantec (Norton) Ghost, it allows you to have a server and create copies over the network:  
[Back Up Restore Hard Drives And Partitions With CloneZilla Live](/pdf/back_up_restore_hard_drives_and_partitions_with_clonezilla_live.pdf)
