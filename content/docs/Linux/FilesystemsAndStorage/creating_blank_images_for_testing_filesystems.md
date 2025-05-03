---
weight: 999
url: "/Créer_des_images_vierges_pour_tester_des_filesystems/"
title: "Creating Blank Images for Testing Filesystems"
description: "How to create blank disk images to safely test filesystems without risking your existing partitions"
categories: ["Linux"]
date: "2007-06-20T08:16:00+02:00"
lastmod: "2007-06-20T08:16:00+02:00"
tags: ["Filesystem", "Testing", "Linux", "Disk Images"]
toc: true
---

If, like me, you want to test a new filesystem or hardware speed without risking damage to one of your partitions, here's a small tip that allows you to create a blank disk image to work with:

```bash
dd if=/dev/zero of=./mon_image.img bs=1M count=128
```

or

```bash
dd if=/dev/zero of=./10M.img bs=10m count=1
```

The second line is more for BSD systems.

Here, I'll have a 128 MB image. Change the last number if you want a different size.

Now, all that's left is to format the partition and start experimenting:

```bash
mkfs.ext3 mon_image.img
```
