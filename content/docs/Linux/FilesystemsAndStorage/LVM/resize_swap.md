---
weight: 999
url: "/Resizer_sa_swap/"
title: "Resize Swap"
description: "How to resize swap space in Linux by extending LVM volumes, creating new swap volumes, or using swap files"
categories:
  - Linux
date: "2007-03-05T10:40:00+02:00"
lastmod: "2007-03-05T10:40:00+02:00"
tags:
  - Linux
  - LVM
  - Swap
  - System Administration
toc: true
---

## Introduction

Sometimes it is necessary to add more swap space after installation. For example, you may upgrade the amount of RAM in your system from 128 MB to 256 MB, but there is only 256 MB of swap space. It might be advantageous to increase the amount of swap space to 512 MB if you perform memory-intense operations or run applications that require a large amount of memory.

You have three options: create a new swap partition, create a new swap file, or extend swap on an existing LVM2 logical volume. It is recommended that you extend an existing logical volume.

## Extending Swap on an LVM2 Logical Volume

To extend an LVM2 swap logical volume (assuming `/dev/VolGroup00/LogVol01` is the volume you want to extend):

Disable swapping for the associated logical volume:

```bash
# swapoff -v /dev/VolGroup00/LogVol01
```

Resize the LVM2 logical volume by 256 MB:

```bash
# lvm lvresize /dev/VolGroup00/LogVol01 -L +256M
```

Format the new swap space:

```bash
# mkswap /dev/VolGroup00/LogVol01
```

Enable the extended logical volume:

```bash
# swapon -va
```

Test that the logical volume has been extended properly:

```bash
# cat /proc/swaps
# free -m
```

## Creating an LVM2 Logical Volume for Swap

To add a swap volume group (assuming `/dev/VolGroup00/LogVol02` is the swap volume you want to add):

Create the LVM2 logical volume of size 256 MB:

```bash
# lvm lvcreate VolGroup00 -n LogVol02 -L 256M
```

Format the new swap space:

```bash
# mkswap /dev/VolGroup00/LogVol02
```

Add the following entry to the `/etc/fstab` file:

```bash
/dev/VolGroup00/LogVol02   swap     swap    defaults     0 0
```

Enable the extended logical volume:

```bash
# swapon -va
```

Test that the logical volume has been extended properly:

```bash
# cat /proc/swaps
# free
```

## Creating a Swap File

To add a swap file:

Determine the size of the new swap file in megabytes and multiply by 1024 to determine the number of blocks. For example, the block size of a 64 MB swap file is 65536.

At a shell prompt as root, type the following command with count being equal to the desired block size:

```bash
# dd if=/dev/zero of=/swapfile bs=1024 count=65536
```

Setup the swap file with the command:

```bash
# mkswap /swapfile
```

To enable the swap file immediately but not automatically at boot time:

```bash
# swapon /swapfile
```

To enable it at boot time, edit `/etc/fstab` to include the following entry:

```bash
# /swapfile          swap            swap    defaults        0 0
```

The next time the system boots, it enables the new swap file.

After adding the new swap file and enabling it, verify it is enabled by viewing the output of the command `cat /proc/swaps` or `free`.
