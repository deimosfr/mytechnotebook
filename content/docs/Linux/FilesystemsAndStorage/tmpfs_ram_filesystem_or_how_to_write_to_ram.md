---
weight: 999
url: "/Tmpfs_\\:_un_filesystem_en_ram_ou_comment_écrire_en_ram/"
title: "Tmpfs: RAM filesystem or how to write to RAM"
description: "A guide on how to create and use a temporary filesystem in RAM for fast access and temporary storage"
categories: ["Linux"]
date: "2012-02-18T18:47:00+02:00"
lastmod: "2012-02-18T18:47:00+02:00"
tags: ["Linux", "Filesystem", "RAM", "Performance"]
toc: true
---

## Introduction

TmpFS (Temporary File System) is the generic name given to any temporary Unix filesystem. Any file created in such a filesystem disappears when the system shuts down.

The default implementation of tmpfs in Linux 2.6.x kernels is based on ramfs which uses the caching mechanism to optimize memory management.  
It is also available on Solaris 10.

However, for security reasons, tmpfs additionally offers a memory size limit set at mount time that can be changed on-the-fly with the "remount" option. Tmpfs also allows the system to use swap when necessary, which provides an additional guarantee.

Unlike a RAM Disk, it dynamically allocates memory so as not to use it excessively, and offers better performance thanks to its extreme simplicity.

## Usage

### Prerequisites

We will create a mount point on `/media/montmpfs`.

First, we need to create the directory:

```bash
mkdir -p /media/montmpfs
```

Then, if necessary, change the permissions on this directory so that everyone can read/write/execute:

```bash
chmod 777 /media/montmpfs
```

### Mounting

Finally, a tmpfs is mounted like all mount points in Linux, with the mount command.

```bash
mount -t tmpfs -o size=256M tmpfs /media/montmpfs
```

The options are:

- `-t`: to specify the file type
- `-o`: for options, including size (if not specified, the default size is equal to half of the RAM)
- then the device, here tmpfs or none (personally I use tmpfs, because with the df command it's listed as tmpfs as the Filesystem)
- then the mount point.

To mount it automatically at startup, you need to edit the `/etc/fstab` file.

Example of a line to add:

```bash
tmpfs /tmp tmpfs defaults,size=1g 0 0
```

## Resources
- [Storing Files Directories In Memory With tmpfs](/pdf/storing_files_directories_in_memory_with_tmpfs.pdf)
- http://www.generation-linux.fr/index.php?post/2009/05/04/tmpfs-%3A-utiliser-sa-ram-comme-repertoire-de-stockage
