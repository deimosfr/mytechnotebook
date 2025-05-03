---
weight: 999
url: "/MFS_\\:_Utiliser_un_filesystem_en_RAM/"
title: "MFS: Using a RAM Filesystem"
description: "A guide on how to use MFS (Memory FileSystem) to improve performance by creating partitions in RAM"
categories: ["Linux", "BSD", "Performance"]
date: "2009-01-05T02:37:00+02:00"
lastmod: "2009-01-05T02:37:00+02:00"
tags: ["MFS", "RAM", "Filesystem", "Performance", "OpenBSD"]
toc: true
---

## Introduction

MFS allows you to place a partition in RAM. The advantage is speed. The disadvantage is that you lose all modifications made to it after each reboot. With a simple rsync setup, this can be resolved, which I'll show you how to do.

## Configuration

### Partition /tmp

The tmp directory is interesting to move into RAM since the data there is temporary anyway and doesn't need to be stored on the filesystem. Edit the fstab file:

```bash
...
swap /tmp mfs rw,nodev,nosuid,-s=32768 0 0
```

Nothing more needs to be done :-)

### Partition /var

After installing the necessary packages (such as net-snmp, pftop, pfstat, screen...), configuring the crontab, and configuring chrooted services like bind or an Apache reverse proxy, we can begin setting up the /var partition in MFS by copying its content to the partition reserved for this purpose.

Note that from now on, all modifications to /var should be made from /mfs/var:

```bash
find /var | cpio -dumpv /mfs/
```

Now let's edit the fstab:

```bash
swap /var mfs rw,-P=/mfs/var,nodev,nosuid,-s=64000 0 0
```

### Syncronisation of modifications

We'll use rsync to update the data. Since we don't need to have everything in real time, a weekly update can be sufficient. Install rsync:

```bash
pkg_add -iv rsync
```

Add this line to the root crontab:

```bash
3  0  *  *  */1  /usr/local/bin/rsync -az --delete /var/ /mfs/var/
```

And finally for machine shutdown:

```bash
/usr/local/bin/rsync -vaz --delete /var/ /mfs/var/
```

## References

- http://www.openbsd.org/cgi-bin/man.cgi?query=mfs
- http://wiki.gcu.info/doku.php?id=openbsd:install_soekris
- https://alpage.org/wiki/doc/openbsd/root_ro
- http://blog.spoofed.org/2007/12/openbsd-on-soekris-cheaters-guide.html
