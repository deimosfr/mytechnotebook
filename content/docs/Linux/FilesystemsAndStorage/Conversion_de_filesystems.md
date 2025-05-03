---
weight: 999
url: "/Conversion_de_filesystems/" 
title: "Filesystem Conversion"
description: "Guide on how to convert filesystems without reformatting, focusing on converting from ext3 to ext4."
categories: ["Ubuntu", "Linux"]
date: "2009-05-08T06:39:00+02:00"
lastmod: "2009-05-08T06:39:00+02:00"
tags: ["filesystem", "ext3", "ext4", "conversion", "linux"]
toc: true
---

## Introduction

Sometimes you need to convert your filesystem to a more recent one when possible, without having to reformat everything.

## EXT3 -> EXT4

It's obvious that the partition to be converted must be unmounted beforehand. Here, I'll use an LVM partition in ext3 to convert it to ext4:

```bash
tune2fs -O extents,uninit_bg,dir_index /dev/mapper/lvm-home
fsck -pf /dev/mapper/lvm-home
```

If you're modifying the boot partition, you need to boot from a live CD for example, then run this command to revalidate grub:

```bash
grub-install /dev/sda
```

## Resources
- http://doc.ubuntu-fr.org/ext4
