---
weight: 999
url: "/SWAP_\\:_Création_de_swap_dynamique/"
title: "SWAP: Creating Dynamic Swap"
description: "Guide on how to create and manage dynamic swap space in Linux through partition-based and image-based methods."
categories: ["Linux"]
date: "2009-09-19T21:47:00+02:00"
lastmod: "2009-09-19T21:47:00+02:00"
tags: ["swap", "Linux", "system management", "memory", "partition"]
toc: true
---

## Introduction

Creating a swap partition is quite straightforward but extremely useful when you're running out of RAM.

## Creating Swap

### On a Partition

Use fdisk on the target partition:

```bash
fdisk /dev/hdc1
```

Then create a primary partition with the key combination **"n p 1"**. Next, change the ID to indicate it's a swap partition with **"t 82"**.

Save everything with the "w" key. To finish setting up our partition, we need to initialize it:

```bash
mkswap /dev/hdc1
```

Now let's use it:

```bash
swapon /dev/hdc1
```

I recommend modifying your `/etc/fstab` file so that the swap is used on the next boot:

```bash
...
/dev/hdc1               swap                    swap    defaults        0 0
...
```

### On a Disk Image

If you can't create a swap partition, you have the option to create a disk image and assign it as additional swap:

```bash
dd if=/dev/zero of=/swapfile bs=1M count=64
chmod 600 /swapfile
mkswap /swapfile
swapon /swapfile
```

Obviously, this size is small and serves as an example. You should adjust the dd command according to the available space on one of your partitions and adapt it to your swap needs.

## Resources
- [Documentation on Creating dynamic swap space](/pdf/creating_dynamic_swap_space.pdf)
