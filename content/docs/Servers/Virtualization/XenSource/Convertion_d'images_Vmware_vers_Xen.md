---
weight: 999
url: "/Convertion_dimages_Vmware_vers_Xen/"
title: "Converting VMware Images to Xen"
description: "Guide on how to convert VMware disk images (vmdk) to Xen format and run them in HVM mode with full virtualization."
categories: ["Linux"]
date: "2009-12-11T21:50:00+02:00"
lastmod: "2009-12-11T21:50:00+02:00"
tags: ["Virtualization", "VMware", "Xen", "Conversion", "HVM"]
toc: true
---

## Introduction

At work, I received some new toys. The particularity of these machines is that they are equipped with "Dual-Core Intel Xeon 5140 Processor (2.33 GHz, 1333 FSB)" CPUs, and the special feature of these CPUs is that in small print at the bottom of their spec sheet, you can read "NOTE: Intel Xeon Processor 5100/5000 sequence are 64-bit, Dual-Core, 4MB L2 Cache, and support Intel VT technology." And that's awesome.

I had been waiting for Santa to come for a while to fulfill a long-thought-out fantasy: shifting a complete lab platform from VMware to Xen. To make it short: not only does it work, but as expected, it performs amazingly well.

## Conversion

Here is a summary of what you need to know to migrate a VMware vmdk image to a Xen image:

Using qemu-img, a tool integrated into qemu which you will obviously install beforehand, we convert the .vmdk to a raw image usable by Xen:

```bash
qemu-img convert test.vmdk test.img
```

## Launching

Here is a working configuration that takes into account the HVM (Hardware-assisted Virtual Machine) mode, aka full virtualization, without modifying the guest. This feature requires a CPU that supports VT technology:

```bash
kernel = '/usr/lib/xen/boot/hvmloader'
builder = 'hvm'

memory  = '256'
name = 'test'

device_model = '/usr/lib/xen/bin/qemu-dm'

# We declare two network cards of type pcnet, those that VMware emulates
nic=2
vif = [ 'type=ioemu, mac=00:50:56:01:09:01, bridge=xenbr0, model=pcnet', \
 'type=ioemu, mac=00:56:3e:00:00:02, bridge=xenbr0, model=pcnet' ]

sdl=0
# The output will be visible on a vnc server on display 1
vnc=1
vnclisten='192.168.10.20'
vncunused=0
vncdisplay=1
vncpasswd=''

disk = [ 'file:/home/xen/test/test.img,hda,w' ]
```

After classically creating your domain:

```bash
xm create test
```

Just connect to the vnc console from a vncviewer client:

```bash
vncviewer 192.168.10.20:5901
```

## Conclusion

You can now show off in front of your salespeople.

## Resources
- [Converting a disk image to LVM](https://wiki.deimos.fr/images/c/cc/Xen-_How_to_Convert_An_Image-Based_Guest_To_An_LVM-Based_Guest.pdf)
