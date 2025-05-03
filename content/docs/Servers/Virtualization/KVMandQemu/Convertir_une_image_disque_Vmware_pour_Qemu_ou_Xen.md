---
weight: 999
url: "/Convertir_une_image_disque_Vmware_pour_Qemu_ou_Xen/"
title: "Converting a VMware Disk Image for Qemu or Xen"
description: "This guide explains how to convert VMware disk images to be compatible with Qemu or Xen virtualization platforms."
categories: ["Ubuntu", "Linux", "Virtualization"]
date: "2008-04-25T08:54:00+02:00"
lastmod: "2008-04-25T08:54:00+02:00"
tags: ["VMware", "Qemu", "Xen", "Virtualization", "Image Conversion"]
toc: true
---

## Introduction

I've switched to the new Ubuntu, great! Except for one thing: VMware doesn't work yet, you have to wait for a patch, etc... but I don't have time to wait!

VMware is nice and pretty, but it's becoming annoying. So I decided to use KVM and QTqemu for the graphical interface. And here I am, ready to convert my VMware images.

## Conversion

For conversion, we'll first use VMware for preparation and then qemu for the actual conversion:

```bash
vmware-vdiskmanager -r vmware_image.vmdk -t 0 temporary_image.vmdk
qemu-img convert -f vmdk temporary_image.vmdk -O raw xen_compatible.img
```

## FAQ

### Great, but my Debian disk won't boot

This happens because if you configured your VMware disk as SCSI, you'll be using sdx device names. Your Debian busybox will launch, and you'll need to mount the partition containing /boot and edit the menu.lst file to change the kernel root. Here's some help:

```bash
mkdir /test
mount -t ext3 /dev/hda1 /test
vi /test/boot/grub/menu.lst
```

Change this:

```bash
kernel		/vmlinuz-2.6.18-6-amd64 root=/dev/sda1 ro
```

to:

```bash
kernel		/vmlinuz-2.6.18-6-amd64 root=/dev/hda1 ro
```

And finally, reboot :-)
