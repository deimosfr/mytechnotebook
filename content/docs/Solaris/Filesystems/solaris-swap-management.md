---
weight: 999
url: "/ajout-de-swap-sous-solaris/"
title: "Adding Swap Space on Solaris"
description: "A guide on how to add, manage, and configure swap space in Solaris using different methods including UFS, ZFS, and swap files."
categories: ["Solaris", "System Administration", "Storage"]
date: "2012-01-30T12:09:00+02:00"
lastmod: "2012-01-30T12:09:00+02:00"
tags:
  [
    "Solaris",
    "Swap",
    "UFS",
    "ZFS",
    "Memory Management",
    "System Administration",
  ]
toc: true
---

## Introduction

The purpose of this documentation is to quickly describe how to create swap space on Solaris. There are other documentation sources such as [Disk Management in Solaris]({{< ref "docs/Solaris/Filesystems/disk_management_in_solaris.md">}}) that describe in depth how disks work on Solaris, but that's not the goal here.

We'll look at three methods to add swap space:

- On UFS
- On ZFS
- Using a swap file

## Swap on UFS

If you want to add swap space on a UFS disk, you'll need to create a partition. Let's first see what we have:

```bash
> swap -l
swapfile             dev  swaplo blocks   free
/dev/dsk/c1t0d0s1   30,65      8 8401984 8401984
```

Next, we'll run the format command then select the disk we want to work with:

```bash {linenos=table,hl_lines=[5,12]}
> Searching for disks...done


AVAILABLE DISK SELECTIONS:
       0. c1t0d0 <DEFAULT cyl 53499 alt 2 hd 255 sec 63>          /pci@0,0/pci10de,375@f/pci108e,286@0/disk@0,0
       1. c2t201500A0B856312Cd31 <DEFAULT cyl 17 alt 2 hd 64 sec 32>
          /pci@7c,0/pci10de,377@f/pci1077,143@0/fp@0,0/disk@w201500a0b856312c,1f
       2. c2t202400A0B856312Cd31 <DEFAULT cyl 17 alt 2 hd 64 sec 32>
          /pci@7c,0/pci10de,377@f/pci1077,143@0/fp@0,0/disk@w202400a0b856312c,1f
       3. c3t201400A0B856312Cd31 <DEFAULT cyl 17 alt 2 hd 64 sec 32>
...
Specify disk (enter its number): 0
```

We choose disk 0 here. Then we'll enter the partition management tool:

```bash {linenos=table,hl_lines=[1]}
format> partition

PARTITION MENU:
        0      - change `0' partition
        1      - change `1' partition
        2      - change `2' partition
        3      - change `3' partition
        4      - change `4' partition
        5      - change `5' partition
        6      - change `6' partition
        7      - change `7' partition
        select - select a predefined table
        modify - modify a predefined partition table
        name   - name the current table
        print  - display the current table
        label  - write partition map and label to the disk
       !<cmd> - execute <cmd>, then return
        quit
```

Next, we'll display the content to see the current partitions:

```bash {linenos=table,hl_lines=[1,11]}
partition> p
Current partition table (original):
Total disk cylinders available: 53499 + 2 (reserved cylinders)

Part      Tag    Flag     Cylinders         Size            Blocks
  0       root    wm     524 -  3134       20.00GB    (2611/0/0)   41945715
  1       swap    wu       1 -   523        4.01GB    (523/0/0)     8401995
  2     backup    wm       0 - 53498      409.82GB    (53499/0/0) 859461435
  3        var    wm    3135 -  4440       10.00GB    (1306/0/0)   20980890
  4 unassigned    wm    4441 -  4571        1.00GB    (131/0/0)     2104515
  5 unassigned    wm       0                0         (0/0/0)             0
  6 unassigned    wm       0                0         (0/0/0)             0
  7       home    wm    4572 -  7182       20.00GB    (2611/0/0)   41945715
  8       boot    wu       0 -     0        7.84MB    (1/0/0)         16065
  9 unassigned    wm       0                0         (0/0/0)             0
```

We can see that the last cylinder is 7182 out of 53498. We'll continue after this cylinder. Let's take a random slice (slice 5 for example) and create a new swap partition on it:

```bash {linenos=table,hl_lines=[1,"4-7"]}
partition> 5
Part      Tag    Flag     Cylinders         Size            Blocks
  5 unassigned    wm       0                0         (0/0/0)             0
Enter partition id tag[unassigned]: swap
Enter partition permission flags[wm]: wu
Enter new starting cyl[1]: 7183
Enter partition size[0b, 0c, 7183e, 0.00mb, 0.00gb]: 70gb
```

Here I've created a partition starting from the last used cylinder (7182) + 1 (7183), a swap partition with the corresponding flag (wu), with a size of 70GB.
Then I display the new partition table:

```bash {linenos=table,hl_lines=[1,11]}
partition> p
Current partition table (unnamed):
Total disk cylinders available: 53499 + 2 (reserved cylinders)

Part      Tag    Flag     Cylinders         Size            Blocks
  0       root    wm     524 -  3134       20.00GB    (2611/0/0)   41945715
  1       swap    wu       1 -   523        4.01GB    (523/0/0)     8401995
  2     backup    wm       0 - 53498      409.82GB    (53499/0/0) 859461435
  3        var    wm    3135 -  4440       10.00GB    (1306/0/0)   20980890
  4 unassigned    wm    4441 -  4571        1.00GB    (131/0/0)     2104515
  5       swap    wu    7183 - 16320       70.00GB    (9138/0/0)  146801970
  6 unassigned    wm       0                0         (0/0/0)             0
  7       home    wm    4572 -  7182       20.00GB    (2611/0/0)   41945715
  8       boot    wu       0 -     0        7.84MB    (1/0/0)         16065
  9 unassigned    wm       0                0         (0/0/0)             0
```

Now I can see my new swap partition. Let's write this new data to the disk:

```bash {linenos=table,hl_lines=[1]}
partition> label
Ready to label disk, continue? y
```

We'll exit and declare this new partition as a swap partition:

```bash
swap -a /dev/dsk/c1t0d0s5
```

Now, when I display the list of active partitions, I can see the new one:

```bash
> swap -l
swapfile             dev  swaplo blocks   free
/dev/dsk/c1t0d0s1   30,65      8 8401984 8401984
/dev/dsk/c1t0d0s5   30,69      8 146801960 146801960
```

I just need to add a line in the vfstab for persistence:

```bash {linenos=table,hl_lines=[9]}
#device device  mount   FS  fsck    mount   mount
#to mount   to  fsck        point       type    pass    at boot options
#
fd  -   /dev/fd fd  -   no  -
/proc   -   /proc   proc    -   no  -
/dev/dsk/c1t0d0s1   -   -   swap    -   no  -
/dev/dsk/c1t0d0s0   /dev/rdsk/c1t0d0s0  /   ufs 1   no  -
/dev/dsk/c1t0d0s3   /dev/rdsk/c1t0d0s3  /var    ufs 1   no  -
/dev/dsk/c1t0d0s5   -   -   swap    -   no  -
/dev/dsk/c1t0d0s7   /dev/rdsk/c1t0d0s7  /export/home    ufs 2   yes -
#/dev/dsk/c1t0d0s4  /dev/rdsk/c1t0d0s4  /globaldevices  ufs 2   yes -
/devices    -   /devices    devfs   -   no  -
sharefs -   /etc/dfs/sharetab   sharefs -   no  -
ctfs    -   /system/contract    ctfs    -   no  -
objfs   -   /system/object  objfs   -   no  -
swap    -   /tmp    tmpfs   -   yes -
/dev/did/dsk/d4s4 /dev/did/rdsk/d4s4 /global/.devices/node@2 ufs 2 no global
```

## Swap on ZFS

You simply need to increase the size of the ZFS associated with the swap.

### Adding a Swap

Let's check how many swaps are allocated:

```bash
> swap -l
swapfile             dev  swaplo blocks   free
/dev/dsk/c1t0d0s1   30,65      8 8401984 8401984
```

Now we add a ZFS:

```bash
zfs create -V 30G rpool/swap1
```

Here we've created a 30GB swap. Then we declare this new partition as swap:

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

If you get this kind of message:

```
/dev/zvol/dsk/rpool/swap is in use for live upgrade -. Please see ludelete(1M).
```

You'll need to use the following command to activate it:

```bash
/sbin/swapadd
```

### Enlarging a Swap

When the machine is running and the swap space is being used, you can increase the size of the swap so that the system can use it. This will require deactivation and reactivation for the new space to be recognized. To do this, we'll enlarge the ZFS:

```bash
zfs set volsize=72G rpool/swap
zfs set refreservation=72G rpool/swap
```

Now we'll deactivate the swap:

```bash
swap -d /dev/zvol/dsk/rpool/swap
```

We need to remove or comment out the swap entry in /etc/vfstab corresponding to the swap, as it will be automatically created in the next step:

```bash
#/dev/zvol/dsk/rpool/swap    -   -   swap    -   no  -
```

Then reactivate it for the new size to be recognized:

```bash
swap -a /dev/zvol/dsk/rpool/swap
```

We can verify the swap size:

```bash
> swap -l
swapfile             dev  swaplo blocs   libres
/dev/zvol/dsk/rpool/swap 181,1       8 150994936 150994936
```

## Swap File

For a swap file, it's the quickest method to implement but also the least elegant. Find a place on your disk where you have space and create an empty file of the desired size:

```bash
mkfile 70g /swap1
```

Then we activate this new swap:

```bash
swap -a /swap1
```

Now, when I display the list of active partitions, I can see the new one:

```bash
> swap -l
swapfile             dev  swaplo blocks   free
/dev/dsk/c1t0d0s1   30,65      8 8401984 8401984
/dev/dsk/c1t0d0s5   30,69      8 146801960 146801960
```
