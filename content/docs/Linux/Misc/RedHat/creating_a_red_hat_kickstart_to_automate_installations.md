---
weight: 999
url: "/Créer_un_Kickstart_Red_Hat_pour_automatiser_les_installation/"
title: "Creating a Red Hat Kickstart to Automate Installations"
description: "Learn how to automate Red Hat installations using the Kickstart method to deploy multiple machines efficiently."
categories: ["RHEL", "Linux", "Red Hat"]
date: "2012-06-07T07:34:00+02:00"
lastmod: "2012-06-07T07:34:00+02:00"
tags: ["Servers", "Linux", "RHEL", "Automation", "Installation"]
toc: true
---

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Operating System** | RHEL 6.2 |
| **Website** | [Debian Website](https://www.debian.org) |
| **Last Update** | 07/05/2013 |
{{< /table >}}

## Introduction

Automating Red Hat installations quickly becomes essential if you have many machines to deploy. This is why the Kickstart method exists, allowing you to boot from virtually any media.

## Kickstart File

The kickstart file is where all the installation and configuration options are defined. I won't describe all available options as they are well documented on the Red Hat website, but rather describe some interesting methods to achieve certain things that aren't natively possible.

Here's an example of a very standard kickstart file to give you an idea:

```bash
# Kickstart file
# Minimal installation for production usage
# Made by Deimos

# Install instead of upgrade
install
# Use text mode install
text
# Use CDROM installation media
cdrom
# Installation logging level
logging --level=info

# System keyboard
keyboard us
# System language
lang en_US

# Disable firewall
firewall --disabled
# SELinux configuration
selinux --disabled

# Root password
rootpw --iscrypted $1$6qaKy76d$R8ToaKxZD4Q89pJWrpE/y.
# System authorization information
auth  --useshadow  --passalgo=sha512

# Do not configure the X Window System
skipx
# System timezone
timezone --isUtc Europe/Paris
# Network information
network --bootproto=dhcp --device=eth0 --onboot=on

# System bootloader configuration
bootloader --location=mbr
# Clear the Master Boot Record
zerombr

# Partition clearing information
clearpart --all --initlabel
# Disk partitioning information
part /boot --fstype="ext4" --size=512
part pv.os --size=65536 --grow
part swap --size=1000 --maxsize=2000
volgroup vgos pv.os
logvol / --vgname=vgos --name=root --size=8096
logvol /var --vgname=vgos --name=var --size=8096
logvol /home --vgname=vgos --name=home --size=32768

# Packages installation
%packages
@base
%end

# Reboot after installation
reboot
```

### Swap = RAM

If you want to have the swap size equal to the RAM size, you'll need to create a pre-script that stores the swap creation line in a temporary file. Then it will be called by the kickstart:

```bash
[...]
# Partition clearing information
clearpart --all --initlabel
# Disk partitioning information
part /boot --fstype="ext4" --size=512
part pv.os --size=65536 --grow
volgroup vgos pv.os
%include /tmp/swappart
logvol / --vgname=vgos --name=root --size=8096
logvol /var --vgname=vgos --name=var --size=8096
logvol /home --vgname=vgos --name=home --size=32768

%pre
#!/bin/sh
act_mem=$(awk '/MemTotal/{print int($2/1024)}' /proc/meminfo)
echo "logvol swap --fstype swap --name=swap --vgname=vgos --size=$act_mem" > /tmp/swappart
```

### Package Groups

To install package groups in the packages section, you need to precede the package group name with @. For example, here's how to get the list of available groups:

```bash
> yum -v grouplist
Installed Groups:
   Base (base)
   Compatibility libraries (compat-libraries)
   E-mail server (mail-server)
[...]
Available Groups:
[...]
   High Availability (ha)
   High Availability Management (ha-management)
[...]
Available Language Groups:
   Afrikaans Support (afrikaans-support) [af]
[...]
```

If, for example, I want to add the High Availability group, I need to add this line:

```bash
# Packages installation
%packages
@base
@ha
%end
```

## Creation

### DVD

If you want to boot from a DVD, insert an original DVD in the drive, then dump it somewhere on your machine:

```bash
mkdir /mnt/rhdvd ~/iso
mount /dev/cdrom /mnt/rhdvd
cp -Rf /mnt/rhdvd/* ~/iso
```

Next, we'll edit the isolinux.cfg file to insert the kickstart parameters at the end:

```bash {linenos=table,hl_lines=[22,26],anchorlinenos=true}
default vesamenu.c32
#prompt 1
timeout 600

display boot.msg

menu background splash.jpg
menu title Welcome to Red Hat Enterprise Linux 6.2!
menu color border 0 #ffffffff #00000000
menu color sel 7 #ffffffff #ff000000
menu color title 0 #ffffffff #00000000
menu color tabmsg 0 #ffffffff #00000000
menu color unsel 0 #ffffffff #00000000
menu color hotsel 0 #ff000000 #ffffffff
menu color hotkey 7 #ffffffff #ff000000
menu color scrollbar 0 #ffffffff #00000000

label linux
  menu label ^Install or upgrade an existing system
  menu default
  kernel vmlinuz
  append initrd=initrd.img ks=cdrom:/ks.cfg
label vesa
  menu label Install system with ^basic video driver
  kernel vmlinuz
  append initrd=initrd.img xdriver=vesa nomodeset ks=cdrom:/ks.cfg
label rescue
  menu label ^Rescue installed system
  kernel vmlinuz
  append initrd=initrd.img rescue
label local
  menu label Boot from ^local drive
  localboot 0xffff
label memtest86
  menu label ^Memory test
  kernel memtest
  append -
```

Then add your kickstart file (ks.cfg) to the root of your iso (~/iso).

Now, let's generate the new ISO:

```bash
mkisofs -o ~/rhks.iso -b isolinux/isolinux.bin -c isolinux/boot.cat -no-emul-boot -boot-load-size 4 -boot-info-table -J -R -V "RedHatKS" .
```

All that's left is to load/burn your ISO and let things happen :-)

## References

http://ihazem.wordpress.com/2012/04/06/creating-dynamic-swap-space-during-linux-kickstart/  
http://docs.redhat.com/docs/en-US/Red_Hat_Network_Satellite/5.3/html/Deployment_Guide/satops-kickstart.html  
http://techspotting.org/creating-a-kickstart-cd-dvd-fedora-redhat-centos/  
http://fedoraproject.org/wiki/Anaconda/Kickstart  
http://ihazem.wordpress.com/2012/04/06/creating-dynamic-swap-space-during-linux-kickstart/
