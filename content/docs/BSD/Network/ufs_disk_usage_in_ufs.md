---
weight: 999
url: "/UFS_\\:_utilisation_des_disques_en_UFS/"
title: "UFS: Disk usage in UFS"
description: "A guide on how to create and manage UFS slices in FreeBSD, including partition creation, formatting, and mounting."
categories: ["FreeBSD", "Linux", "Backup"]
date: "2012-06-14T12:30:00+02:00"
lastmod: "2012-06-14T12:30:00+02:00"
tags: ["FreeBSD", "UFS", "Disk Management", "Partitioning"]
toc: true
---

![FreeBSD](/images/poweredbyfreebsd.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Operating System** | FreeBSD 9 |
| **Website** | [FreeBSD Website](https://www.freebsd.org) |
| **Last Update** | 14/06/2012 |
{{< /table >}}

## Introduction

Today with BTRFS and ZFS, it's super simple to create and delete partitions, etc., but with UFS (the old way), it's less obvious, although not very complicated. I'll talk about some common practices here :-)

## Creating a slice

First, let's display the content of our disk:

```bash
> gpart show ada0
=>       34  976773101  ada0  GPT  (465G)
         34        128     1  freebsd-boot  (64k)
        162   41943040     2  freebsd-ufs  (20G)
   41943202  924843904        - free -  (441G)
  966787106    8388608     3  freebsd-swap  (4.0G)
  975175714    1597421        - free -  (780M)
```

Here we can see I have 3 slices. The idea is to add a slice (for backup needs, so we'll call it backups) on this GPT partition table. I specify this because we're going to use the gpt command and not fdisk! Let's add 50G:

```bash
> gpart add -t freebsd-ufs -l backups -s 50G ada0
ada0p4 added
```

And if we display it now:

```bash {linenos=table,hl_lines=[5]}
> gpart show ada0
=>       34  976773101  ada0  GPT  (465G)
         34        128     1  freebsd-boot  (64k)
        162   41943040     2  freebsd-ufs  (20G)
   41943202  104857600     4  freebsd-ufs  (50G)
  146800802  819986304        - free -  (391G)
  966787106    8388608     3  freebsd-swap  (4.0G)
  975175714    1597421        - free -  (780M)
```

Oh yeah! Now we just need to format it as UFS:

```bash
> newfs -U /dev/gpt/backups
/dev/gpt/backups: 51200.0MB (104857600 sectors) block size 32768, fragment size 4096
[...]
```

And mount it:

```bash
mount /dev/gpt/backups /mnt/
```

If we want persistence, we need to add a line in fstab:

```bash {linenos=table,hl_lines=[4]}
# Device        Mountpoint      FStype  Options Dump    Pass#
/dev/ada0p3     none            swap    sw      0       0
/dev/ada0p2     /               ufs     rw      1       1
/dev/ada0p4     /mnt            ufs     rw      2       2
```

## References

http://www.wonkity.com/~wblock/docs/html/disksetup.html
