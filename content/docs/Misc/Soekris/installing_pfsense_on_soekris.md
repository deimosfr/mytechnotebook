---
weight: 999
url: "/Installer_pfSense_sur_Soekris/"
title: "Installing pfSense on Soekris"
description: "A guide on how to install pfSense on a Soekris net5501 hardware platform including solutions for common issues."
categories: ["Linux", "Networking"]
date: "2010-12-03T23:29:00+02:00"
lastmod: "2010-12-03T23:29:00+02:00"
tags: ["pfSense", "Soekris", "Firewall", "Installation"]
toc: true
---

## Introduction

I spent too much time trying to find how to install pfSense on a Soekris net5501. Why? Because no PXE versions exist or are easily installable, and also due to configuration and connection issues. To help others save time and for my own reference, I've decided to write this article.

## Installation

Let's say I have a hard drive to install pfSense on. I connected it to my Ubuntu Desktop laptop through a USB external 2.5" box and used KVM/QEMU. Here are the required packages:

```bash
aptitude install kvm virt-manager
```

Do not forget to add your user to the 'kvm' group.

Now use virt-manager to create a VM with:

* CDROM: the pfSense ISO
* HDD: the direct USB box

Then boot the VM and perform the installation.

## Configuration

### 1st boot

For the first configuration part, you should install the hard drive in the Soekris box and boot it. I'm using minicom and I've set these serial parameters as: **9600 7E1**

Then you'll likely encounter this issue:

```
Trying to mount root from ufs:/dev/ad0s1a
Trying to mount root from ufs:/dev/ad0s1a
Trying to mount root from ufs:/dev/ad0s1a

Manual root filesystem specification:
  <fstype>:<device>  Mount <device> using filesystem <fstype>
                     eg. ufs:da0s1a
  ?                  List valid disk boot devices
  <empty line>       Abort manual input
```

This is because the device during installation didn't get the same device name on the Soekris. There is a way to fix this. Try something like this:

```
mountroot> ?
List of GEOM managed disk devices:
  ufsid/4cf95e52836e2e4f ad1s1c ad1s1b ad1s1a ad1s1 ad1

Manual root filesystem specification:
  <fstype>:<device>  Mount <device> using filesystem <fstype>
                     eg. ufs:da0s1a
  ?                  List valid disk boot devices
  <empty line>       Abort manual input
```

Now we're going to boot on the correct root slice:

```bash
ufs:ad1s1a
```

### Web interface

Now the web interface is available through port 0 on IP 192.168.1.1. The default credentials are:

* Login: admin
* Password: pfsense

Proceed with your desired configuration.

### Remote connection

Now you're able to set default connections. I mean correct default serial parameters and no exotic ones, as well as the SSH server. Go to System -> Advanced to enable it.

### Fstab

To avoid manually mounting the root filesystem at next boot, we'll change the fstab file through our new SSH connection. Change the device values to the correct ones. For me, I had to change from ad0 to ad1 (`/etc/fstab`):

```
# Device                Mountpoint      FStype  Options         Dump    Pass#
/dev/ad1s1a             /               ufs     rw              1       1
/dev/ad1s1b             none            swap    sw              0       0
```

### Soekris

Now reboot and set the connection parameters to 9600 to get the hardware and OS at the same configuration level.

That's all :-)
