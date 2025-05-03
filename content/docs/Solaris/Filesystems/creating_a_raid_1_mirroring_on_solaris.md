---
weight: 999
url: "/Creation_d'un_Raid_1_(mirroring)_sous_Solaris/"
title: "Creating a RAID 1 (mirroring) on Solaris"
description: "Guide to setting up RAID 1 mirroring for filesystems in Solaris using DiskSuite and ZFS."
categories: ["Backup", "Linux"]
date: "2009-02-11T16:15:00+02:00"
lastmod: "2009-02-11T16:15:00+02:00"
tags: ["Solaris", "Storage", "Mirroring", "RAID", "UFS", "ZFS"]
toc: true
---

## Introduction

The Solaris system includes the DiskSuite package that allows RAID 1 mirroring of a UFS filesystem using LVM. This tutorial explains how to achieve this.
It goes without saying that you need two disks of the same capacity.

## Procedure

Here are the necessary steps:
Once we're ready to mirror a disk, we display its partitions with the format command:

```bash
Format
 
Searching for disks...done
 
 
AVAILABLE DISK SELECTIONS:
       0. c1t0d0 <SUN146G cyl 14087 alt 2 hd 24 sec 848>  root
          /pci@0/pci@0/pci@2/scsi@0/sd@0,0
       1. c1t1d0 <HITACHI-H101414SCSUN146G-SA25-136.73GB>
          /pci@0/pci@0/pci@2/scsi@0/sd@1,0
```

We choose the 1st disk (if it's the one to be duplicated), then:

```bash
format> partition
```

And we display the partition table with "p":

```bash
partition> p
Volume:  root
Current partition table (original):
Total disk cylinders available: 14087 + 2 (reserved cylinders)
 
Part      Tag    Flag     Cylinders         Size            Blocks
  0       root    wm       0 - 14086      136.71GB    (14087/0/0) 286698624
  1 unassigned    wm       0                0         (0/0/0)             0
  2     backup    wm       0 - 14086      136.71GB    (14087/0/0) 286698624
  3 unassigned    wm       0                0         (0/0/0)             0
  4 unassigned    wm       0                0         (0/0/0)             0
  5 unassigned    wm       0                0         (0/0/0)             0
  6 unassigned    wm       0                0         (0/0/0)             0
  7 unassigned    wm       0                0         (0/0/0)             0
```

Note: In Solaris, like in BSD, the 3rd partition (no. 2) is actually the entire disk.

So now we have a view of our existing partitions...

### UFS

Important: We need to create a small partition of about 20MB that will host the "metadata" concerning the RAID 1. This metadata will be used by DiskSuite.
The first step will be to copy the partition table from the 1st disk to the 2nd.
Then we will create databases for the metadata.
Then we will manually decide which partition will be mirrored by creating sub-mirrors.
We will change the vfstab (the file that indicates which partition mounts where).
We will attach the sub-mirrors to a mirror.
We will create aliases for the mirrors.
We will add this alias to the "boot-device".

### ZFS

We will simply create an identical partition to the one on the master disk and set it as root.
Go directly to [Copying the partition table to the 2nd disk](#copying-the-partition-table-to-the-2nd-disk)

## Creating a small partition for the metadata

```
partition> p
Volume:  root
Current partition table (original):
Total disk cylinders available: 14087 + 2 (reserved cylinders)

Part      Tag    Flag     Cylinders         Size            Blocks
 0       root    wm       0 - 14086      136.71GB    (14087/0/0) 286698624
 1 unassigned    wm       0                0         (0/0/0)             0
 2     backup    wm       0 - 14086      136.71GB    (14087/0/0) 286698624
 3 unassigned    wm       0                0         (0/0/0)             0
 4 unassigned    wm       0                0         (0/0/0)             0
 5 unassigned    wm       0                0         (0/0/0)             0
 6 unassigned    wm       0                0         (0/0/0)             0
 7 unassigned    wm       0                0         (0/0/0)             0
```

Enter the number of the partition to edit and press "Enter"
Choose the tag "unassigned"
flag: vm
size: 20mb

Then exit:

```
label
```

## Copying the partition table to the 2nd disk

```bash
prtvtoc /dev/rdsk/c0t0d0s2 | fmthard -s - /dev/rdsk/c0t1d0s2
```

Use the 2nd slice to indicate the entire disk.

## Creating metadata database for DiskSuite

```bash
metadb -a -f -c2 /dev/dsk/c0t0d0s3 /dev/dsk/c0t1d0s3
```

CAUTION: Choose the correct partition letter on both disks (the small one we created)

## Creating sub-mirrors

### UFS

Let's assume we want to mirror the 6 partitions of the disk (except the swap), for example / /usr /var /opt /home and /etc

Let's start with / (root partition):

```bash
 metainit -f d10 1 1 c0t0d0s0 
 metainit -f d20 1 1 c0t1d0s0
 metainit d0 -m d10
 metaroot d0 (Use this command only on the root slice!)
```

CAUTION: Be sure to enter the correct disk names.
So here we have associated the partition containing / on the 1st disk with the mirror partition that will be on the 2nd disk, then we indicated that the 1st partition will be the master, then we specified that it was the root partition.

You need to do this for each partition (except the last command)

For /usr:

```bash
 metainit -f d11 1 1 c0t0d0s1
 metainit -f d21 1 1 c0t1d0s1
 metainit d1 -m d11
```

For /var:

```bash
 metainit -f d14 1 1 c0t0d0s4
 metainit -f d24 1 1 c0t1d0s4
 metainit d4 -m d14
```

For /opt:

```bash
 metainit -f d15 1 1 c0t0d0s5
 metainit -f d25 1 1 c0t1d0s5
 metainit d5 -m d15
```

For /etc:

```bash
 metainit -f d16 1 1 c0t0d0s6
 metainit -f d26 1 1 c0t1d0s6
 metainit d6 -m d16
```

For /home:

```bash
 metainit -f d17 1 1 c0t0d0s7
 metainit -f d27 1 1 c0t1d0s7
 metainit d7 -m d17
```

We can view the metadata with the command:

```bash
metastat
```

### ZFS

```bash
zpool attach -f rpool c0t0d0s0 c0t1d0s0
```

Once finished, that's it! You can stop here, it's done for ZFS.

## Editing the vfstab file

```bash
vi /etc/vfstab
```

From now on, vfstab will no longer point to a disk but to a cluster. Here are the lines to edit:

Before, for /: /dev/md/dsk/d30 /dev/md/rdsk/d30        /       ufs     1       no      logging
After, for /: /dev/md/dsk/d0 /dev/md/rdsk/d0        /       ufs     1       no      logging

d0 will be the partition for /, then d1, d2, d3, etc...

At this point we can restart with these two commands in succession:

```bash
lockfs -fa
init 6
```

## Attaching mirrors to sub-mirrors

```bash
 metattach d0 d20
 metattach d1 d21
 metattach d4 d24
 metattach d5 d25
 metattach d6 d26
 metattach d7 d27
```

These commands will start the synchronization of mirrors and sub-mirrors with each other. You can see the progress with "metastat".

Then we change the crash dump:

```bash
dumpadm -d `swap -l | tail -1 | awk '{print $1}'`
```

## Creating mirror aliases

We need to know the absolute path of the mirrored disk:

```bash
ls -l /dev/dsk/c0t1d0s0
lrwxrwxrwx   1 root     root          50 Jan 16 10:20 /dev/rdsk/c0t1d0s0 -> ../../devices/pci@1f,0/pci@1,1/ide@3/dad@1,0:a
```

With this, we will create an alias for the mirror, replacing "dad" with "disk":

```bash
eeprom "nvramrc=devalias mirror /pci@1f,0/pci@1,1/ide@3/disk@1,0:a"
eeprom "use-nvramrc?=true"
```

## Adding mirrors to the boot device

```bash
eeprom "boot-device=disk mirror net"
```

Then if we only have 2 disks, we need to add this line to the /etc/system file:

```
set md:mirrored_root_flag = 1
```

## Resources
- http://www.brandonhutchinson.com/Mirroring_disks_with_DiskSuite.html
