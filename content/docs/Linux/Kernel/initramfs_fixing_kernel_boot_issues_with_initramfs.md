---
weight: 999
url: "/Initramfs_\\:_corriger_les_petits_problèmes_de_boot_kernel_grâce_à_initramfs/"
title: "Initramfs: Fixing Kernel Boot Issues with Initramfs"
description: "This guide explains how to use initramfs to solve kernel boot issues when disk detection order changes."
categories: ["Linux", "Debian", "System Administration"]
date: "2007-09-18T22:19:00+02:00"
lastmod: "2007-09-18T22:19:00+02:00"
tags: ["Kernel", "Boot", "Initramfs", "Troubleshooting", "RAID", "Storage"]
toc: true
---

## Problem

I just installed a server with many SATA disks. The machine has 4.5TB of storage spread across 2 Areca ARC-1280ML controllers. The Debian/etch installation went without issues using kernel 2.6.18-5-686. After installing the system on 2 disks connected to the motherboard (ICH5R controller using the ata_piix driver), RAID5 volumes are created on the Areca cards (arcmsr driver). However, boot stops at an initramfs command prompt, unable to find the root partition:

```bash
Begin: Waiting for root file system... ...
Done.
        Check root= bootarg cat /proc/cmdline
        or missing modules, devices: cat /proc/modules ls /dev
ALERT! /dev/sda1 does not exist. Dropping to a shell!


Busybox v1.1.3 (Debian 1:1.1.1-4) Built-in shell (ash)
Enter 'help' for a list of built-in commands.

/bin/sh: can't access tty; job control turned off
(initramfs)
```

## Explanation

What's happening? Simply put, the new volumes are detected by the kernel before the disk on which the system is installed. As a result, the system is no longer on `/dev/sda` but on `/dev/sdc`. And the funniest part is that it's sometimes on `/dev/sdb` because the Areca controllers take time to initialize.

How do we fix this issue? By working with the initialization RAM partition, namely initramfs.

## Solution

It's extremely simple. We'll ask the RAM boot partition to load the SATA modules in the order we want. In our case, the ata_piix driver before arcmsr. Debian tools make this very easy, just add the modules you want loaded during startup to the `/etc/initramfs-tools/modules` file. Modules should be listed one per line in the desired loading order. In our case, we just need to specify the module that handles the boot disk.

```bash
# cat /etc/initramfs-tools/modules
[...]
ata_piix
```

Now we need to update the RAM image to apply these changes:

```bash
# update-initramfs -v -k 2.6.18-5-686 -t -u
Keeping /boot/initrd.img-2.6.18-5-686.dpkg-bak
update-initramfs: Generating /boot/initrd.img-2.6.18-5-686
Adding module /lib/modules/2.6.18-5-686/kernel/drivers/scsi/scsi_mod.ko
Adding module /lib/modules/2.6.18-5-686/kernel/drivers/scsi/scsi_transport_spi.ko
Adding module /lib/modules/2.6.18-5-686/kernel/drivers/scsi/aic7xxx/aic7xxx.ko
[...]
Adding binary /sbin/mdrun
Building cpio /boot/initrd.img-2.6.18-5-686 initramfs
Backup /boot/initrd.img-2.6.18-5-686.bak
```

After a reboot, everything is back to normal. The best part is that when you need to update your kernel, the new kernel will automatically rebuild the initialization RAM image.
