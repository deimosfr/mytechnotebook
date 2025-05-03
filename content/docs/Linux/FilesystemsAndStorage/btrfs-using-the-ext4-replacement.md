---
weight: 999
url: "/BTRFS_\\:_Utilisation_du_remplaçant_de_l'Ext4/"
title: "BTRFS: Using the Ext4 Replacement"
description: "Learn how to use BTRFS filesystem, the replacement for ExtX, including creating partitions, subvolumes, snapshots, compression, and RAID configurations."
categories: ["Linux", "Filesystem", "Storage"]
date: "2012-07-05T21:08:00+02:00"
lastmod: "2012-07-05T21:08:00+02:00"
tags: ["BTRFS", "Filesystem", "Linux", "Ext4", "Storage", "ZFS"]
toc: true
---

![BTRFS](/images/btrfs_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 0.19 |
| **Operating System** | Debian 7 |
| **Website** | [BTRFS Website](https://btrfs.wiki.kernel.org) |
| **Last Update** | 05/07/2012 |
| **Others** | Kernel used:<br>3.2.0-2-amd64 |
{{< /table >}}

## Introduction

[BTRFS](https://btrfs.wiki.kernel.org/) is the perfect replacement for the aging ExtX filesystem. For those familiar with [the ZFS filesystem](zfs:_le_filesystem_par_excellence/), BTRFS draws heavily from it.

BTRFS, like Ext4, is based on the concept of extents. This is a contiguous area (which can reach several hundred MB, unlike the clusters of some older formats) reserved each time a file is saved on the hard drive. This allows, in case of writing at the end of a file (append) or a complete rewrite, to often add the new data directly to the extent rather than in another area of the hard disk, which would increase fragmentation. Large files are thus stored more efficiently through a greater disk space occupation, but at a cost that has decreased considerably. BTRFS stores the data of very small files directly in the extent of the directory file, not in a separate extent.

BTRFS manages a notion of "subvolumes" allowing, within the filesystem, to have a separate tree (including the root) containing directories and files, giving the possibility to have various trees simultaneously, and therefore a greater independence from the main system. This also allows for better separation of data and imposing different quotas on different subvolumes. The most practical use of this system concerns snapshots. A snapshot offers the possibility to "take a photograph" at a given moment of a filesystem to back it up. This snapshot under BTRFS is a subvolume, which allows it to be modified afterward. Having a snapshot accessible in write mode is of obvious interest for high-availability online databases.

To exploit these subvolumes and snapshots, BTRFS uses the classic technique of "Copy-on-write". If data is written to a memory block, then the block will be copied to another location in the filesystem and the new data will be recorded on the copy instead of on the original. Then the metadata pointing to the block is automatically modified to take into account the new data. We thus have a transactional mechanism distinct from the journaling present in Ext3. Before each write, taking a snapshot of the system would allow, in case of a problem, to return to the snapshot, but this seems to pose, if not performance problems, at least questions: should you take a snapshot at each write, or for a certain volume of data? This also raises the question of time lost at each creation/destruction of a snapshot. The use of snapshots for this purpose is not emphasized by the developers.

BTRFS has its own data protection techniques: the use of back references (i.e., knowing, from a data block, which metadata points to the block) allows identification of system corruptions. If a file claims to belong to a set of blocks and these blocks claim to be related to another file, this indicates that the consistency of the system is altered. BTRFS also performs checksums on all data and stored metadata to detect all kinds of corruptions on the fly, repair some of them, and thus offer a better level of reliability.

It allows hot resizing of the filesystem size (including shrinking it) while maintaining excellent protection of metadata that is duplicated in several places for security. The operation is simple: btrfsctl -r +2g /mnt adds 2 GiB to your filesystem. This function is not intended to be redundant with what the Linux logical volume manager offers but claims to technically complement it.

Checking the filesystem through the btrfsck program is error-tolerant and presented as extremely fast by its design. The use of B-trees allows exploring the disk structure at a speed essentially limited by the disk's read speed. The price to pay is a strong memory footprint since btrfsck uses three times more memory than e2fsck.

BTRFS respects the hierarchy of Linux's functional "layers". For example, while offering functions to complement it, it tries as much as possible not to rewrite the whole volume management system proposed as standard by LVM.

Google's lightweight and fast Snappy compression algorithm was added in January 2012, allowing for faster data access compared to LZO compression (around 10%) and no compression (around 15%).

It was followed by the LZ4 compression algorithm in February 2012, which further improves performance compared to Snappy (by about 20-30%).

## Installation

We'll need to install the BTRFS tools:

```bash
aptitude install btrfs-tools parted
```

## Usage

In my use cases below, many examples will be in relation to a configuration of this type:

```bash
$ fdisk -l

Disk /dev/sda: 10.7 GB, 10737418240 bytes
255 heads, 63 sectors/track, 1305 cylinders, total 20971520 sectors
Units = sectors of 1 * 512 = 512 bytes
Sector size (logical/physical): 512 bytes / 512 bytes
I/O size (minimum/optimal): 512 bytes / 512 bytes
Disk identifier: 0x0007f48a

   Device Boot      Start         End      Blocks   Id  System
/dev/sda1   *        2048     5859327     2928640   83  Linux
/dev/sda2         5859328     7813119      976896   82  Linux swap / Solaris
/dev/sda3         7813120     9910271     1048576   83  Linux
/dev/sda4         9910272    16201727     3145728   83  Linux
```

sda3 and sda4 are partitions we'll be working with.

### Creating a BTRFS Partition

I have a 3GB partition here (sda4). I'll format it as BTRFS:

```bash
$ mkfs.btrfs /dev/sda4

WARNING! - Btrfs Btrfs v0.19 IS EXPERIMENTAL
WARNING! - see http://btrfs.wiki.kernel.org before using

fs created label (null) on /dev/sda4
	nodesize 4096 leafsize 4096 sectorsize 4096 size 3.00GB
Btrfs Btrfs v0.19
```

**Note: It is strongly recommended to create a BTRFS partition on LVM for future hot resizing!**

My partition is ready to be mounted:

```bash
mount -t btrfs /dev/sda4 /mnt/
```

And we can see that the partition is correctly mounted:

```bash
$ mount | grep sda4
/dev/sda4 on /mnt type btrfs (rw,relatime,space_cache)
```

### Subvolumes

Just like ZFS, it's possible to create subvolumes. That is, in a formatted partition (think of it as a VG under LVM), it's possible to create subvolumes (LV under LVM) but whose use allows great data flexibility.

Let's create our first subvolume:

```bash
btrfs subvolume create /mnt/volume1
```

A subvolume is materialized at the directory tree level by a directory present at the root of the volume's mount point.

```bash
$ btrfs subvolume list /mnt/
ID 145 top level 3 path volume1
```

The number 145 uniquely identifies our subvolume. The volume1 path is also indicated. The volume1 path, which is also the name of the subvolume, is relative to the root mount of our btrfs volume.

To mount a subvolume at the same location:

```bash
mount -o subvol=volume1 /dev/sda4 /mnt/
```

You can also mount your subvolume from its identifier:

```bash
mount -o subvolid=145 /dev/sda4 /mnt
```

### Converting an extX Partition to BTRFS

It's possible to convert an ext3 or ext4 partition to btrfs! In this case, I'll convert sda3, which is already in ext4. We'll use the btrfs-convert command:

```bash
$ btrfs-convert /dev/sda3
creating btrfs metadata.
creating ext2fs image file.
cleaning up system chunk.
conversion complete.
```

And there it is, that's all :-). It's simple, right! I can now mount it:

```bash
mount -t btrfs /dev/sda3 /mnt/
```

And verify that everything is good:

```bash
$ mount | grep sda3
/dev/sda3 on /mnt type btrfs (rw,relatime,space_cache)
```

### Resizing a Partition

#### Method 1: Filesystem Expansion

To increase the size of a partition on the fly, it's very simple as long as we're on LVM or if we have the partition that physically has at least the desired size. Otherwise, we'll need to do a cold operation to expand the partition, then the filesystem.
The difference with an ext-type filesystem is that we can specify that the partition size can take a size x without taking the whole (like resize2fs).

```bash
$ btrfs filesystem resize +1G /mnt/
Resize '/mnt/' of '+1G'
```

#### Method 2: Adding a Device

We have another possibility in case our current volume becomes too small... we can add a device to an existing volume. First, let's get a status of the btrfs volumes present on our system:

```bash
$ btrfs filesystem show
failed to read /dev/sr0
Label: none  uuid: 9fd825a0-3bee-44b6-88e9-b8e4bc554e82
	Total devices 1 FS bytes used 92.00KB
	devid    1 size 3.00GB used 343.12MB path /dev/sda4

Label: none  uuid: 0773e361-0342-4960-8f8e-5a26db8bc93e
	Total devices 1 FS bytes used 49.40MB
	devid    1 size 1.00GB used 1.00GB path /dev/sda3

Btrfs Btrfs v0.19
```

If I look at my mounted partitions:

```bash
$ df -h
[...]
/dev/sda3                                               1.0G   50M  640M   8% /mnt
```

I have my sda3 which is 1G. We'll add sda4 to it:

```bash
btrfs device add /dev/sda4 /mnt/
```

Let's check that /mnt has been enlarged:

```bash
> df -h
[...]
/dev/sda3                                               4.0G   50M  3.7G   2% /mnt
```

#### Method 3: RAID 0

There is a solution like RAID 0, but with data distribution and metadata replication of the filesystem on all disks:

```bash
$ btrfs filesystem balance /mnt
$ btrfs filesystem df /mnt/
Data, RAID0: total=3.71GB, used=256.00KB System,
RAID1: total=8.00MB, used=8.00KB
System: total=4.00MB, used=0.00
Metadata, RAID1: total=344.34MB, used=32.00KB
```

The data is in RAID0 between the different partitions while the system and metadata are RAID1. This means that even if you lose one of the system partitions, you will still be able to mount the remaining partition. However, since the data is not replicated but distributed, you will have lost the data present on the disappeared partition.

### Reducing a Partition

To reduce a partition on the fly, it's really very simple, just specify the size, then the mount point on which you want to remove the size:

```bash
$ btrfs filesystem resize -1G /mnt/
Resize '/mnt/' of '-1G'
```

And that's it :-)

### RAID 1

It's possible to do software RAID 1. I remind you that you need disks of the same size, or it will be the size of the smallest disk that will be used. Let's start by initializing our RAID:

```bash
mkfs.btrfs -m raid1 -d raid1 /dev/sda3 /dev/sda4
```

Then you can mount sda3 or sda4, the replication is done :-)

### Compression

#### Cold Method

It's possible to have a compressed filesystem. For this, nothing could be simpler:

```bash
mount -t btrfs -o compress=lzo /dev/sda3 /mnt/
```

If we add the compress-force option, the compression on files that btrfs will be greater. By default, btrfs doesn't compress well, because for large files it can lead to a lot of I/O. The behavior of btrfs's on-the-fly compression algorithm therefore tries to spare the processor when it determines according to its first operations if a file can be difficult to compress:

```bash
mount -o compress=zlib,compress-force /dev/sda3 /mnt/
```

#### Hot Method

If you want to perform the same compression activation operation directly from a mounted btrfs filesystem, we can use the following command which will activate the compression option and compress the data already present on the disk:

```bash
btrfs filesystem defragment -czlib /mnt
```

### Snapshot

There are several complementary tools that allow managing snapshots such as [Snapper](https://fr.opensuse.org/Portal:Snapper) or yum-plugin-fs-snapshot on Fedora/RedHat. But for now, we'll see how to manage snapshots the standard btrfs way.

We create from a volume the snapshot that will allow us to make modifications on this file tree:

```bash
btrfs subvolume nsnapshot /mnt /mnt/snapshot
```

We now unmount the current volume to mount our snapshot instead in which we create a new file:

```bash
umount /mnt/
mount -o subvol=snapshot /dev/sda4 /mnt/
```

#### Revert

If we want to cancel our modifications, it will be very simple, we will unmount our snapshot and mount the old one:

```bash
umount /mnt
mount /dev/sda4 /mnt
```

#### Merge

If we want to merge the data, we need to retrieve the ID of our subvolume. Then the set-default order, followed by the subvolume ID followed by the original volume's mount point allows declaring a new default volume:

```bash
$ btrfs subvolume list /mnt
ID 146 top level 5 path snapshot
$ btrfs subvolume set-default 146 /mnt
```

## References

[https://fr.wikipedia.org/wiki/Btrfs](https://fr.wikipedia.org/wiki/Btrfs)  
[https://www.funtoo.org/wiki/BTRFS_Fun](https://www.funtoo.org/wiki/BTRFS_Fun)  
[https://www.rashardkelly.com/extending-a-btrfs-filesystem-2/](https://www.rashardkelly.com/extending-a-btrfs-filesystem-2/)
