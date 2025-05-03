---
weight: 999
url: "/Parted_\\:_résoudre_les_problèmes_de_partionnnement_sur_gros_filesystems/"
title: "Parted: Solving Partitioning Problems on Large Filesystems"
description: "Guide on how to use the Parted tool to solve partitioning issues with large filesystems and disks over 2TB, with commands and examples for proper partitioning."
categories: ["Linux", "Debian", "FreeBSD"]
date: "2013-02-13T12:54:00+02:00"
lastmod: "2013-02-13T12:54:00+02:00"
tags: ["parted", "storage", "partitioning", "filesystem", "disk management", "large disks"]
toc: true
---

![Parted logo](/images/parted_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|---|---|
|**Software version**|2.1|
|**Operating System**|Debian 6<br>Red Hat 6.3|
|**Last Update**|13/02/2013|
{{< /table >}}

## Introduction

GNU Parted is a program for creating, destroying, resizing, checking, and copying partitions, and the file systems on them. This is useful for creating space for new operating systems, reorganizing hard disk usage, copying data between hard disks, and disk imaging. It was written by Andrew Clausen and Lennert Buytenhek.

It consists of a library, libparted, and a command-line frontend, parted, that also serves as reference implementation.

Currently, Parted runs only under Linux, GNU/Hurd, FreeBSD and BeOS.

## Problem Statement

You may have just acquired a disk array (SATA, SAS...) and are encountering partitioning issues. For example, you might not be able to create partitions larger than 99GB or 2TB. This is what happened to me. And in your boot logs (dmesg), you might see errors like:

```bash
sdb : very big device. try to use READ CAPACITY(16).
Losing some ticks... checking if CPU frequency changed.
SCSI device sdb: 8784445440 512-byte hdwr sectors (4497636 MB)
sdb: Write Protect is off
sdb: Mode Sense: 1f 00 10 08
SCSI device sdb: drive cache: write through w/ FUA
sdb: sdb1
sdb : very big device. try to use READ CAPACITY(16).
SCSI device sdb: 8784445440 512-byte hdwr sectors (4497636 MB)
sdb: Write Protect is off
sdb: Mode Sense: 1f 00 10 08
SCSI device sdb: drive cache: write through w/ FUA
sdb:
sdb : very big device. try to use READ CAPACITY(16).
SCSI device sdb: 8784445440 512-byte hdwr sectors (4497636 MB)
sdb: Write Protect is off
sdb: Mode Sense: 1f 00 10 08
SCSI device sdb: drive cache: write through w/ FU
```

You likely partitioned with fdisk or cfdisk. And that's where your error lies! Apparently, these tools still have difficulty handling large capacities. That's why I'm recommending parted :-)

## Installation

Quite simple:

```bash
apt-get install parted
```

## Partitioning

### Wizard

My partitioning is fairly simple; I want to create a single 4.5TB disk. First, I check that I don't have any existing partitions and I begin:

```bash
$ parted 
GNU Parted 1.7.1
Using /dev/sda
Welcome to GNU Parted! Type 'help' to view a list of commands.
(parted) p                                                                

Disk /dev/sda: 145GB
Sector size (logical/physical): 512B/512B
Partition Table: msdos

Number  Start   End     Size    Type     File system  Flags
 1      32,3kB  41,1MB  41,1MB  primary  fat16             
 2      41,1MB  173MB   132MB   primary  ext3         boot 
 3      173MB   145GB   145GB   primary               lvm
```

Here we can see that I'm on my first disk (sda) and not on my array (sdb). Let's switch to it:

```bash
(parted) select /dev/sdb                                                  
Using /dev/sdb
```

Let's create a label:

```bash
(parted) mklabel
New disk label type?  [msdos]? gpt
```

Next, create the partition:

```bash
(parted) mkpart                                                           
Partition name?  []? san           # Give any name you want
File system type?  [ext2]? ext3    # The filesystem        
Start? 0                           # The beginning of your disk                        
End? -1                            # -1 corresponds to the end
```

For those who want to use LVM, simply add a flag:

```bash
set 1 lvm on
```

### Command line

If you want to do everything from the command line, for example to automate the process, here's how. I've made a small script that creates a partition taking up the entire disk and optimizes it (via disk alignment) for the best performance:

```bash
datas_device=/dev/sdb
parted -s -a optimal $datas_device mklabel gpt
parted -s -a optimal $datas_device mkpart primary ext4 0% 100%
parted -s $datas_device set 1 lvm on
```

* Line 1: we create a gpt-type label for large partitions (greater than 2TB)
* Line 2: we create a partition that takes up the entire disk
* Line 3: we indicate that this partition will be LVM type

## Validation

We can see the result:

```bash
(parted) p                                                                

Disk /dev/sdb: 4498GB
Sector size (logical/physical): 512B/512B
Partition Table: gpt

Number  Start   End     Size    File system  Name  Flags
 1      17,4kB  4498GB  4498GB  ext3         san
```

Everything looks good. I can exit:

```bash
(parted) q
```

And verify once more:

```bash
$ cat /proc/partitions 
major minor  #blocks  name

   8     0  142082048 sda
   8     1      40131 sda1
   8     2     128520 sda2
   8     3  141910177 sda3
   8    16 4392222720 sdb
   8    17 4392222686 sdb1
 254     0    2097152 dm-0
 254     1   20971520 dm-1
 254     2   20971520 dm-2
 254     3   20971520 dm-3
 254     4   76894208 dm-4
```

## Formatting

All that's left is to format it in the desired format (ext3 in this case):

```bash
mkfs.ext3 /dev/sdb1
```

## Maximum Filesystem Sizes

{{< table "table-hover table-striped" >}}
|Filesystem|Maximum partition size|Maximum file size|
|----------|----------------------|-----------------|
|ext2|8 TB|2 GB|
|ext3|8 TB|2 TB|
|ext4|1024 PB||
|Resiser4||8 TB|
|ZFS|16 EB|16 EB|
{{< /table >}}
