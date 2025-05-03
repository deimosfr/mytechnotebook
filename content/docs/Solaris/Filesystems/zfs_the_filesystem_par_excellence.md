---
weight: 999
url: "/ZFS_\\:_Le_FileSystem_par_excellence/"
title: "ZFS: The Filesystem Par Excellence"
description: "Complete guide on using ZFS filesystem, including creating and managing zpools, partitions, swap space management, and advanced usage techniques"
categories:
  - Linux
date: "2013-01-02T13:30:00+02:00"
lastmod: "2013-01-02T13:30:00+02:00"
tags:
  - Servers
  - Storage
  - Solaris
  - Filesystems
toc: true
---

![ZFS](/images/zfs_logo.avif)

## Introduction

ZFS or Z File System is an open-source filesystem under the CDDL license. The 'Z' doesn't officially stand for anything specific, but has been referred to in various ways in the press, such as Zettabyte (from the English unit zettabyte for data storage), or ZFS as "the last word in filesystems".

Produced by Sun Microsystems for Solaris 10 and above, it was designed by Jeff Bonwick's team. Announced for September 2004, it was integrated into Solaris on October 31, 2005, and on November 16, 2005, as a feature of OpenSolaris build 27. Sun announced that ZFS was integrated into the Solaris update dated June 2006, a year after the opening of the OpenSolaris community.

The characteristics of this filesystem include its very high storage capacity, the integration of all previous concepts related to filesystems and volume management into a single product. It integrates on-disk structure, is lightweight, and easily allows setting up a storage management platform.

## Locating Your Disks

Use your usual tools to identify your disks. For example, under Solaris:

```bash
bash-3.00# format
Searching for disks...done


AVAILABLE DISK SELECTIONS:
       0. c0t600A0B80005A2CAA000004104947F51Ed0 <SUN-LCSM100_F-0670-10.00GB>
          /scsi_vhci/disk@g600a0b80005a2caa000004104947f51e
       1. c0t600A0B80005A2CB20000040B4947F57Fd0 <DEFAULT cyl 1303 alt 2 hd 255 sec 63>
          /scsi_vhci/disk@g600a0b80005a2cb20000040b4947f57f
       2. c1t1d0 <DEFAULT cyl 5098 alt 2 hd 255 sec 63>
          /pci@0,0/pci8086,25e2@2/pci8086,3500@0/pci8086,3510@0/pci1000,3150@0/sd@1,0
       3. c2t1d31 <DEFAULT cyl 17 alt 2 hd 64 sec 32>
          /pci@0,0/pci8086,25f8@4/pci1077,143@0/fp@0,0/disk@w203400a0b85a2caa,1f
       4. c3t3d31 <DEFAULT cyl 17 alt 2 hd 64 sec 32>
          /pci@0,0/pci8086,25f8@4/pci1077,143@0,1/fp@0,0/disk@w203500a0b85a2caa,1f
```

## Zpool

A Zpool is similar to a VG (Volume Group) for those familiar with LVM. The issue is that at the time of writing this article, you cannot reduce the size of a zpool, but you can increase it. You can use a zpool directly as a filesystem since it's based on ZFS. However, you can create ZFS filesystems (that's what they're called, I know it's confusing, but think of them more like LVs (Logical Volumes) or partitions). You can also create other filesystems containing other filesystem types (NTFS, EXT4, REISERFS...).

### Creating a Zpool

To create a zpool, follow these steps:

```bash
zpool create zpool_name c0t600A0B80005A2CAA000004104947F51Ed0
```

- zpool_name: specify the name you want for the zpool
- c0t600A0B80005A2CAA000004104947F51Ed0: this is the device name displayed by the format command

### Listing Zpools

#### Simple Zpool

To know which pools exist on the machine:

```bash
zpool list
```

#### Raid-Z

A Raid-Z is like a Raid 5, but without one major problem: no parity resynchronization or loss during a power outage. Here's how to do it:

```bash
zpool create my_raidz raidz c1t0d0 c2t0d0 c3t0d0 c4t0d0 /dev/dsk/c5t0d0
```

### Mounting a Zpool

By default, zpools have the mount point /zpool_name. To mount a zpool, we'll use the zfs command:

```bash
zfs mount zpool_name
```

It will remember where it should be mounted, as this information is stored in the filesystem.

### Unmounting a Zpool

This is super simple, as usual:

```bash
umount /zpool_name
```

Just use umount followed by the mount point.

### Deleting a Zpool

To delete a zpool:

```bash
zpool destroy zpool_name
```

### Expanding a Zpool

To expand a zpool, we'll use the zpool name and the additional device:

```bash
zpool add zpool_name device1 device2...
```

### Modifying Zpool Parameters

#### Changing the Mount Point of a Zpool

By default, zpools are mounted in /, to change this:

```bash
zfs set mountpoint=/mnt/datas my_zpool
```

- /mnt/datas: the desired mount point
- my_zpool: name of the zpool

### Importing All Zpools

To import all zpools:

```bash
zpool import -f -a
```

- -f: force (optional and may be dangerous in some cases)
- -a: will import all zpools

### Renaming a Zpool

Renaming a zpool is actually not very complicated:

```bash
zpool export mon_zpool
zpool import mon_zpool mon_nouveau_nom_de_zpool
```

And that's it, the zpool is renamed :-).

### Using in a Cluster Environment

In a cluster environment, you'll need to mount and unmount zpools quite regularly. If you use Sun Cluster (at the time of writing this in version 3.2), you are forced to use zpools for mounting partitions. Filesystems cannot be mounted and unmounted from one node to another when they belong to the same zpool.

You'll therefore need to unmount the zpool, export the information in ZFS, then import it on the other node. Imagine the following scenario:

- sun-node1 (node 1)
- sun-node2 (node 2)
- 1 disk array with 1 LUN of 10 GB

The LUN is a zpool created as described earlier in this document on sun-node1. Now you need to switch this zpool to sun-node2. Since ZFS is not a cluster filesystem, you must unmount it, export the information, then import it. Unmounting is not mandatory since the export will do it, but if you want to do things properly, then let's proceed on sun-node1:

```bash
umount /zpool_name
export zpool_name
```

Now, list the available zpools, you should no longer see it. Let's move to sun-node2:

```bash
zpool import zpool_name
```

There you go, you'll find your files, normally it's automatically mounted, but if for some reason it's not, you can do it manually (see above).

## ZFS

### Creating a ZFS Partition

To create a ZFS partition, it's extremely simple! You obviously have your Zpool created first, then you execute:

```bash
zfs create zpool/partition
```

Then you can specify options with -o and see all available options with:

```bash
zfs get all zpool/partition
```

### Renaming a Partition

To rename a ZFS partition. If you want to rename a ZFS partition, nothing could be simpler:

```bash
zfs rename zpool/partitionold zpool/partitionnew
```

## Managing Swap on ZFS

On Solaris, you can use multiple combined swap spaces (partitions + files indifferently).

To list the swap spaces used:

```bash
> swap -l
swapfile                  dev     swaplo blocks   free
/dev/zvol/dsk/rpool/swap1 181,3       8  4194296  4194296
/dev/zvol/dsk/rpool/swap  181,2       8  62914552 62914552
```

To know a bit more, you can also use:

```bash
> swap -s
total: 283264k bytes allocated + 258412k reserved = 541676k used, 31076108k available
```

Here we have two ZFS volumes used as swap. By default, when creating a ZPOOL, a swap space is created "rpool/swap".

### Adding SWAP

You just need to increase the size of the ZFS associated with the swap.

#### Adding a Swap

Let's verify the number of assigned swaps:

```bash
> swap -l
swapfile             dev  swaplo blocks   free
/dev/dsk/c1t0d0s1   30,65      8 8401984 8401984
```

Now we add a ZFS:

```bash
zfs create -V 30G rpool/swap1
```

Here we've just created the swap at 30G. Then we declare this new partition as swap:

```bash
swap -a /dev/zvol/dsk/rpool/swap1
```

Now, when I display the list of active partitions, I can see the new one:

```bash
> swap -l
swapfile             dev  swaplo blocks   free
/dev/dsk/c1t0d0s1   30,65      8 8401984 8401984
/dev/dsk/c1t0d0s5   30,69      8 146801960 146801960
```

If you get a message like this:

```
/dev/zvol/dsk/rpool/swap is in use for live upgrade -. Please see ludelete(1M).
```

You'll need to use the following command to activate it:

```bash
/sbin/swapadd
```

#### Expanding a Swap

When the machine is running and the swap space is being used, you can increase the size of the swap so the system can use it. This will require deactivation and reactivation for the new space to be taken into account. For this, we'll expand the zfs:

```bash
zfs set volsize=72G rpool/swap
zfs set refreservation=72G rpool/swap
```

Now we'll deactivate the swap:

```bash
swap -d /dev/zvol/dsk/rpool/swap
```

You now need to delete or comment out the entry in `/etc/vfstab` that corresponds to the swap, as it will be automatically created in the next step:

```bash
#/dev/zvol/dsk/rpool/swap    -   -   swap    -   no  -
```

Then reactivate it so the new size is taken into account:

```bash
swap -a /dev/zvol/dsk/rpool/swap
```

You can check the swap size:

```bash
> swap -l
swapfile             dev  swaplo blocs   libres
/dev/zvol/dsk/rpool/swap 181,1       8 150994936 150994936
```

## Advanced Usage

### The ZFS ARC Cache

The problem with ZFS is that it's very RAM-hungry (about 1/8 of the total + swap). This can quickly become problematic on machines with a lot of RAM. Here's some explanation.

#### Available Memory mdb -k and ZFS ARC I/O Cache

The `mdb -k` command with the `::memstat` option provides a global view of available memory on a Solaris machine:

```bash
echo ::memstat | mdb -k

Page Summary                Pages                MB  %Tot
------------     ----------------  ----------------  ----
Kernel                     587481              2294    7%
Anon                       180366               704    2%
Exec and libs                6684                26    0%
Page cache                   7006                27    0%
Free (cachelist)            13192                51    0%
Free (freelist)           7591653             29654   91%

Total                     8386382             32759
Physical                  8177488             31943
```

In the example above, this is a machine with 32 GB of physical memory.

ZFS uses a kernel cache called ARC for I/Os. To know the size of the I/O cache at a given time used by ZFS, use the `kmastat` option with the `mdb -k` command and look for the Total [zio_buf] statistic:

```bash
echo ::kmastat | mdb -k

cache                        buf    buf    buf    memory     alloc alloc
name                        size in use  total    in use   succeed  fail
------------------------- ------ ------ ------ --------- --------- -----
 ...
Total [zio_buf]                                1157632000   1000937     0
 ...
```

In the example above, the ZFS I/O cache uses 1.1 GB in memory.

#### Limiting the ZFS ARC Cache (zfs_arc_max and zfs_arc_min)

For machines with a very large amount of memory, it's better to limit the ZFS I/O cache to avoid any memory overflow on other applications. In practice, this cache increases and decreases dynamically based on the needs of applications installed on the machine, but it's better to limit it to prevent any risk. The `zfs_arc_max` parameter (in bytes) in the `/etc/system` file allows limiting the amount of memory for the ZFS I/O cache. Below is an example where the ZFS I/O cache is limited to 4GB:

```bash
...
set zfs:zfs_arc_max = 4294967296
...
```

#### Statistics on the ZFS ARC Cache (kstat zfs)

Similarly, you can specify the minimum amount of memory to allocate to the ZFS I/O cache with the `zfs_arc_min` parameter in the `/etc/system` file.

The `kstat` command with the `zfs` option provides detailed statistics on the ZFS ARC cache (hits, misses, size, etc.) at a given time: you'll find the maximum possible value (c_max) for this cache, the current size (size) in the output of this command. In the example below, the `zfs_arc_max` parameter hasn't yet been applied, which explains why the maximum possible size corresponds to the physical memory of the machine.

```bash
> kstat zfs
module:zfs                             instance: 0
name:  arcstats                        class:    misc
       c                               33276878848
       c_max                           33276878848
       c_min                           4159609856
       crtime                          121.419237623
       deleted                         497690
       demand_data_hits                14319099
       demand_data_misses              6491
       demand_metadata_hits            45356553
       demand_metadata_misses          33470
       evict_skip                      2004
       hash_chain_max                  4
       hash_chains                     1447
       hash_collisions                 1807933
       hash_elements                   40267
       hash_elements_max               41535
       hdr_size                        6992496
       hits                            60821130
       l2_abort_lowmem                 0
       l2_cksum_bad                    0
       l2_evict_lock_retry             0
       l2_evict_reading                0
       l2_evict_reading                0
       l2_feeds                        0
       l2_free_on_write                0
       l2_hdr_size                     0
       l2_hits                         0
       l2_io_error                     0
       l2_misses                       0
       l2_rw_clash                     0
       l2_size                         0
       l2_writes_done                  0
       l2_writes_error                 0
       l2_writes_hdr_miss              0
       l2_writes_sent                  0
       memory_throttle_count           0
       mfu_ghost_hits                  3387
       mfu_hits                        53995731
       misses                          48704
       mru_ghost_hits                  1180
       mru_hits                        5891117
       mutex_miss                      0
       p                               21221559296
       prefetch_data_hits              237031
       prefetch_data_misses            3520
       prefetch_metadata_hits          908447
       prefetch_metadata_misses        5223
       recycle_miss                    0
       size                            1362924368
       snaptime                        14013729.1668961

module:zfs                             instance: 0
name:  vdev_cache_stats                class:    misc
       crtime                          121.419271852
       delegations                     4453
       hits                            27353
       misses                          9753
       snaptime                        14013729.1677954
```

I also recommend the excellent [arc_summary](https://cuddletech.com/arc_summary/) which provides very precise information or [arcstat](https://blogs.sun.com/realneel/resource/arcstat.pl).

### Not Mounting All Zpools at Boot

If you encounter errors when booting your machine during zpool mounting (in my case, a continuous reboot of the machine), there is a solution that allows it to forget all those that were imported during the last session (ZFS remembers imported zpools and automatically reimports them during a boot, which is convenient but can be constraining in some cases).

To do this, you'll need to boot in single user or multi-user mode (if it doesn't work, try failsafe mode for Solaris (chrooted for Linux)), then we'll remove the ZFS cache:

```bash
mv /etc/zfs/zpool.cache /etc/zfs/zpool.cache.`date "+%Y-%m-%d"`
```

If your OS is installed on ZFS, and you're in failsafe mode, you'll need to repopulate the cache (the /a corresponds to the root under Solaris):

```bash
cd /
bootadm update-archive -R /a
umount /a
```

Then reboot and reimport the desired zpools.

## FAQ

### FAULTED

I got a little **FAULTED** as you can see below without really knowing why:

```bash
bash-3.00# zpool list
NAME             SIZE   USED  AVAIL    CAP  HEALTH  ALTROOT
test1      -      -      -      -  FAULTED  -
```

Here we need to debug a bit. For that, we use the following command:

```bash
bash-3.00# zpool status -x
  pool: test1
 state: UNAVAIL
status: One or more devices could not be opened.  There are insufficient
        replicas for the pool to continue functioning.
action: Attach the missing device and online it using 'zpool online'.
   see: http://www.sun.com/msg/ZFS-8000-3C
 scrub: none requested
config:

        NAME                                     STATE     READ WRITE CKSUM
                test1                           UNAVAIL      0     0     0  insufficient replicas
         c4t600A0B800048A9B6000005B84A8293C9d0  UNAVAIL      0     0     0  cannot open
```

That doesn't look good! The error messages are scary. Yet the solution is simple: **just export and reimport (forcing if necessary) the defective zpools.**

Then, you can check the status of your filesystem via a scrub:

```bash
zpool scrub test1
zpool status
```

### How to Repair Grub After Zpool Upgrade

It appears that zpool upgrade can break the grub bootloader.

To fix this, we need to reinstall grub on the partition. Proceed as follows:

- Unplug all fiber cables from the server (or other disks than the OS)
- Boot on Solaris 10 Install DVD
- On boot Select option 6: "Start Single User Shell"
- It will scan for existing zpool containing OS installation then ask you if you want to mount your rpool to /a. Answer "yes"
- When you get the prompt, launch this to check the status of the rpool:

```bash
> zpool status rpool
  pool: rpool
 state: ONLINE
 scrub: none requested
config:

        NAME        STATE     READ WRITE CKSUM
        rpool       ONLINE       0     0     0
          c3t0d0s0  ONLINE       0     0     0

errors: No known data errors
```

- This way we can see all the disks involved in the zpool (here c3t0d0s0). We will reinstall grub on all the disks in the zpool with this command:

```bash
> installgrub -m /boot/grub/stage1 /boot/grub/stage2 /dev/rdsk/c3t0d0s0

Updating master boot sector destroys existing boot managers (if any).
continue (y/n)? y
stage1 written to partition 1 sector 0 (abs 96406065)
stage2 written to partition 1, 267 sectors starting at 50 (abs 96406115)
stage1 written to master boot sector
```

- Unmount the zpool

```bash
zpool export rpool
```

- Plug back the fiber cables
- Reboot

```bash
init 6
```

That's it!

## Resources

- http://fr.wikipedia.org/wiki/ZFS
- [ZFS Admin Documentation](/pdf/zfsadmin.pdf)
- http://jean-francois.im/2008/04/faulted-argh.html
- http://www.sqlpac.com/referentiel/docs/unix-solaris-10-zfs-oracle.htm#L6480
