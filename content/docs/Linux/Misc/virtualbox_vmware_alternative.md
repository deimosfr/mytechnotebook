---
weight: 999
url: "/VirtualBox_\\:_Alternative_à_Vmware/"
title: "VirtualBox: VMware Alternative"
description: "A guide on installing and configuring VirtualBox as an alternative to VMware on Linux systems"
categories: ["Linux", "Debian", "Ubuntu"]
date: "2009-01-28T03:17:00+02:00"
lastmod: "2009-01-28T03:17:00+02:00"
tags: ["Virtualization", "VirtualBox", "VMware", "Networking", "Solaris"]
toc: true
---

## Introduction

VirtualBox is a virtual machine created by InnoTek for Windows, 32 and 64-bit GNU/Linux, and Mac OS X hosts supporting Windows (including Vista), Linux 2.x, OS/2 Warp, OpenBSD, and FreeBSD as guest systems. After several years of development, VirtualBox was released under the GNU GPL license in January 2007.

## Installation

### Debian

Once you've downloaded the version, installation is quite simple:

```bash
dpkg -i virtualbox*.deb
```

Then you'll need the bridge utilities:

```bash
apt-get install bridge-utils
```

### Solaris

Once you've downloaded the version, installation is quite simple:

```bash
gtar -xzvf virtualbox*.tgz
pkgadd -G -d VirtualBoxKern*.pkg
pkgadd -d VirtualBox-*.pkg
```

## Configuration

### Networks

Note that since version 2.1.0, it is no longer necessary to manage the network configuration as explained below, as everything is pre-configured.

#### Network Card Configuration

Here's the necessary network configuration:

```bash
# This file describes the network interfaces available on your system
# and how to activate them. For more information, see interfaces(5).

# The loopback network interface
auto lo eth0 br0
iface lo inet loopback

# The primary network interface
iface eth0 inet static

iface br0 inet dhcp
bridge_ports eth0
```

Now let's restart the network:

```bash
/etc/init.d/networking restart
```

#### Bridged Interfaces Configuration

If you want to use your VMs in bridge mode (as if they were separate computers on the network), you need to bridge your network cards. Add interfaces planned for this purpose:

```bash
VBoxAddIF vbox0 pmavro br0
VBoxAddIF vbox1 pmavro br0
VBoxAddIF vbox2 pmavro br0
VBoxAddIF vbox3 pmavro br0
```

Here, I've added 4 to have some margin.

### Adding the User to the vboxusers Group

Next, **add your current user to the vboxusers group**.

And that's it! Log out and log back in, and you're all set!

## FAQ

### I Changed Kernel and VirtualBox No Longer Starts VMs

Simply run this command:

```bash
/etc/init.d/vboxdrv setup
```

### Failed to open/create the internal network...

This small issue on Solaris can be resolved as follows:

```bash
rem_drv vboxflt
add_drv vboxflt
```

### verr_vm_driver_not_installed

To solve this problem:

```bash
ln -s /devices/pseudo/vboxdrv\@0\:vboxdrv /dev/vboxdrv
```

### Ubuntu: Virtualbox Ubuntu unable to boot please use a kernel appropriate

This is due to the default kernel of the server version which is 686. Change it to 386 and everything will work like magic.

At the end of the installation, you can chroot into your new system and install the kernel: linux-image-386.

## Resources
- [VirtualBox Documentation](/pdf/virtualboxfc6centos.pdf)
- [Controlling VirtualBox from the Command Line](https://www.linux.com/feature/151029)
- [Advanced Networking Linux configuration for VirtualBox](https://www.virtualbox.org/wiki/Advanced_Networking_Linux)
