---
weight: 999
url: "/Tmpfs_\\:_monter_un_filesystem_en_RAM_sous_Solaris/"
title: "Tmpfs: Mounting a RAM filesystem on Solaris"
description: "A guide on how to set up and manage tmpfs on Solaris systems to create temporary filesystems in RAM"
categories: ["Solaris", "Linux"]
date: "2012-02-19T09:07:00+02:00"
lastmod: "2012-02-19T09:07:00+02:00"
tags: ["Solaris", "Filesystem", "RAM", "Memory", "System Administration", "Performance"]
toc: true
---

## Introduction

TmpFS (Temporary File System) is the generic name given to any temporary Unix file system. Any file created in such a filesystem disappears when the system is shut down.

The default implementation of tmpfs in Linux 2.6.x kernels is based on ramfs, which uses the caching mechanism to optimize memory management.
It is also available on Solaris 10.

However, tmpfs additionally offers a memory size limit that is set at mount time and can be modified on-the-fly with the "remount" option for security purposes. Tmpfs also allows the system to use swap space when necessary, which provides an additional guarantee.

Unlike a RAM Disk, it allocates memory dynamically to avoid excessive usage and offers better performance due to its extreme simplicity.

## Usage

### Prerequisites

We will create a mount point at `/media/montmpfs`.

First, create the directory:

```bash
mkdir -p /media/montmpfs
```

Then change the permissions on this directory so everyone can read/write/execute:

```bash
chmod 777 /media/montmpfs
```

### Mounting

Once the prerequisites are completed, we can mount the partition:

```bash
mount -F tmpfs -o size=2048m  swap /media/montmpfs
```

* -o size=2048m: specify the desired size. If none is defined, it will be the size of RAM + swap.
* /media/montmpfs: mount point for tmpfs

To have this mount persist across reboots, add it to the vfstab file (`/etc/vfstab`):

```bash
...
swap    -   /media/montmpfs    tmpfs   -   yes size=2048m
...
```

### Expanding the Partition

You didn't plan enough space for your partition? We can expand it - it's not simple, but it's feasible.

**WARNING: Performing this operation in production can be very risky**

#### Retrieving Information

First, let's check the available space:

```bash
> df -h /media/montmpfs
Filesystem             size   used  avail capacity  Mounted on
swap                   2,0G     272K   2,0G     1%    /media/montmpfs
```

Let's retrieve the memory address used for our tmpfs:

```bash {linenos=table,hl_lines=[1,2]}
> echo "::fsinfo" | mdb -k | egrep "VFSP|/media/montmpfs"
            VFSP FS              MOUNT
ffffffffbd1e10c0 tmpfs           /media/montmpfs
```

Here's the allocation address: **ffffffffbd1e10c0**.

Now, let's retrieve the address of the tm_anonmax variable so we can change its value later:

```bash {linenos=table,hl_lines=[1]}
> echo "ffffffffbd1e10c0::print vfs_t vfs_data | ::print -ta struct tmount tm_anonmax" | mdb -k
ffffffffbfb42068 ulong_t tm_anonmax = 0x80000
```

Here tm_anonmax (number of pages) is equal to 0x80000 (512kb) at memory address ffffffffbfb42068.

#### Applying the New Value

We will now change its value. Let's say we want to increase it to 3GB. First, we need to retrieve the default block size allocated for swap:

```bash {linenos=table,hl_lines=[2]}
> pagesize 
4096
```

So I have 4kb blocks.

If I want to change to 3GB, I need to calculate the value in hexadecimal:

```
Desired size => desired size in kb / block size in kb = size in Kb = size in hexadecimal
3G => 3145728Kb / 4kb = 786432Kb = 0xC0000
```

Now let's use these values and apply them to the current memory address to expand it:

```bash {linenos=table,hl_lines=[1]}
> echo "ffffffffbfb42068/Z 0xC0000" | mdb -kw
0x3000f488d00:                  0x80000                 =       0xC0000
```

Let's verify that our changes were applied correctly:

```bash {linenos=table,hl_lines=[1]}
> echo "ffffffffbfb42068/J" | mdb -k
0x3000f488d00:             0xC0000
```

or

```bash {linenos=table,hl_lines=[1]}
> echo "ffffffffbfb42068::print vfs_t vfs_data | ::print struct tmount tm_anonmax" | mdb -k
tm_anonmax = 0xC0000
```

And check the result:

```bash {linenos=table,hl_lines=[3]}
> df -h /media/montmpfs
Filesystem             size   used  avail capacity  Mounted on
swap                   3,0G     272K   3,0G     1%    /media/montmpfs
```

## Resources
- http://docs.oracle.com/cd/E19963-01/html/821-1459/fscreate-99040.html
- http://ilapstech.blogspot.com/2009/11/grow-tmp-filesystem-tmpfs-on-line-under.html
- [Solaris Tmpfs documentation](/pdf/solaris_tmpfs.pdf)
