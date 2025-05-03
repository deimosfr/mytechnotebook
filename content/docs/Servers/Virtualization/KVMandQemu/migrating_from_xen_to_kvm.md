---
weight: 999
url: "/Migration_Xen_vers_KVM/"
title: "Migrating from Xen to KVM"
description: "A guide on how to migrate virtual machines from Xen to KVM virtualization platform."
categories: ["Debian", "Linux"]
date: "2008-04-18T22:50:00+02:00"
lastmod: "2008-04-18T22:50:00+02:00"
tags: ["virtualization", "Xen", "KVM", "migration", "Debian"]
toc: true
---

## Introduction

You may have, like me, surrendered to the joys of user-friendly virtualization to manage some of your projects, for testing, or just to show off... In short, you used the tools from the "xen-tools" package on your Debian machine, particularly `/usr/bin/xen-create-image` to create your DomU instances.

## Problem

Over time, maintaining your various modules with the xenified kernel of your Dom0 has become tedious, so you decided to switch to KVM, "it's so simple" (except, of course, if you're using a precompiled Debian kernel).

The problem is that your "foo.img" image is not directly bootable, since by default, Debian used a kernel placed on the Dom0 to boot the DomUs.

## Solutions

Here are a few tricks to handle the problem...

### First solution

If you've kept those kernels (meaning you didn't remove everything while screaming "never again this !@#$ Xen"), then you just need to launch KVM with the appropriate parameter, but I won't go into detail as it's trivial. (For example: `qemu -hda image.img -kernel /boot/vmlinuz-2.6.22-3-686`)

### Second Solution

If you've removed the kernels, then using a host kernel will be difficult, unless you recreate a ramdisk, and then it will be less easy to move the image to run it elsewhere with qemu... So in this case, the nice and relatively quick solution is:

- Download a Debian installer ISO (just as an example)
- Boot from it with: `qemu -hda image.img -cdrom ../Downloads/debian-40r1-i386-businesscard.iso -boot d`
- Follow the installation process until you have network access, then switch to a console
- Mount /dev/hda in /mnt (after a possible fdisk)
- Mount /proc, /sys and /dev in /mnt (`mount -o loop`)
- Chroot into /mnt
- Install a suitable kernel
- Install LILO (since GRUB seems to have trouble with this type of system as it stands; you can install it later)
- Edit the LILO configuration and execute LILO

Now you can boot your image directly with KVM. You can also take the opportunity to convert the raw image to qcow2, change the network configuration, etc.
