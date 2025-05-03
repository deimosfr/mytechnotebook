---
weight: 999
url: "/Xen_convertion_P2V/"
title: "Xen Conversion P2V"
description: "Guide on how to convert a physical server to a Xen virtual machine through P2V conversion process."
categories: ["Linux", "Debian"]
date: "2008-03-04T17:32:00+02:00"
lastmod: "2008-03-04T17:32:00+02:00"
tags: ["Virtualization", "Xen", "P2V", "Migration", "Server Management"]
toc: true
---

## 1 Introduction

How to remedy the physical deterioration of a server (oldpc) by migrating the OS into a brand new Xen HVM domU.

## 2 Preparing space for oldpc on the xen server

Create a disk image (20 GB size):

```bash
dd if=/dev/zero of=/vieuxpc/vieuxpc.hda.img bs=1024k count=20000
```

Create a loop device corresponding to the entire disk image:

```bash
losetup /dev/loop0 /vieuxpc/vieuxpc.hda.img
```

Partition the image:

```bash
fdisk /dev/loop0
```

List the created partitions:

```bash
fdisk -ul  /dev/loop0 
Disk /dev/loop0: 20.9 GB, 20971520000 bytes
255 heads, 63 sectors/track, 2549 cylinders, total 40960000 sectors
Units = sectors of 1 * 512 = 512 bytes

      Device Boot      Start         End      Blocks   Id  System
/dev/loop0p1              63    17591174     8795556   83  Linux
/dev/loop0p2        17591175    19567169      987997+  82  Linux swap / Solaris
```

Attach the partitions as loop devices:

```bash
losetup -o $((63*512)) /dev/loop1 /vieuxpc/vieuxpc.hda.img
losetup -o $((17591175*512)) /dev/loop2 /vieuxpc/vieuxpc.hda.img
```

Format: (note that 8795556 is the number of blocks as displayed by fdisk)

```bash
mkfs.ext3 /dev/loop1 8795556
mkswap /dev/loop2 987997
```

## 3 Copy the data from oldpc

Mount the root partition of oldpc and copy the data:

```bash
mount /dev/loop0 /target
mount /dev/sdc2 /source
cd /source
tar cf - . | ( cd /target/ && tar xf -)
```

Edit `/source/boot/grub/menu.lst` to match your partition scheme:

```bash
root            (hd0,0)
kernel          /boot/vmlinuz-2.6.18-5-486 root=/dev/hda1 ro
initrd          /boot/initrd.img-2.6.18-5-486
```

Unmount everything:

```bash
umount /target
umount /source
losetup -d /dev/loop{0,1,2}
```

## 4 Install the bootloader

Install GRUB on the disk using a GRUB rescue floppy:

```bash
cat /usr/lib/grub/stage[12] > floppy.img
```

Note: Debian users can install grub-disk, a bootable GRUB image preconfigured to boot many operating systems.

```bash
qemu -fda floppy.img -hda  /vieuxpc/vieuxpc.hda.img -boot a 

GNU GRUB  version 0.95  (639K lower / 31744K upper memory)

[ Minimal BASH-like line editing is supported.  For the first word, TAB
  lists possible command completions.  Anywhere else TAB lists the possible
  completions of a device/filename. ]

grub> root (hd0,0)
Filesystem type is ext2fs, partition type 0x83

grub> setup (hd0)
Checking if "/boot/grub/stage1" exists... yes
Checking if "/boot/grub/stage2" exists... yes
Checking if "/boot/grub/e2fs_stage1_5" exists... no
Running "install /boot/grub/stage1 (hd0) /boot/grub/stage2 p /boot/grub/menu.lst "... succeeded
Done.

grub>
```

## 5 Test and put into production

Test everything with qemu:

```bash
qemu -hda  /vieuxpc/vieuxpc.hda.img
```

Create the Xen-hvm configuration file (see http://wiki.gcu.info/doku.php?id=unix:xen_hvm for example):

```bash
name="vieuxpc"
kernel = "/usr/lib/xen/boot/hvmloader"
builder = "hvm"
vif=['type=ioemu, mac=00:16:3E:00:03:05, bridge=xenbr0']
disk = [ 'file:/vieuxpc/vieuxpc.hda.img,ioemu:hda,w' ]
device_model = "/usr/lib/xen/bin/qemu-dm"
dhcp="dhcp"
memory="256"
on_poweroff = 'destroy'
on_reboot   = 'restart'
on_crash    = 'restart'
vnc=1
```

On first boot, reconfigure the network interface, since the physical interface obviously no longer exists. Instead, you have an emulated ne2000 network card.
