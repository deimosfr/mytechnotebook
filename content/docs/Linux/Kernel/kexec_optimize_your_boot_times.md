---
weight: 999
url: "/Kexec_\\:_optimisez_vos_temps_de_boot/"
title: "Kexec: Optimize Your Boot Times"
description: "A guide on how to use kexec to optimize boot times by skipping hardware initialization, especially useful for high availability systems."
categories: ["Linux", "Red Hat", "Debian"]
date: "2012-03-15T13:41:00+02:00"
lastmod: "2012-03-15T13:41:00+02:00"
tags: ["Linux", "Boot", "System Administration", "Performance", "Red Hat", "Debian"]
toc: true
---

## Introduction

Kexec is a tool that allows rebooting a machine without going through the hardware layer. It will stop all services, shut down the init processes (sysV) to reach the bootloader. Then it will start normally, bypassing the hardware reboot phase.

This technique can be very useful on High Availability systems where downtime is critical.

## Installation

### Debian

On Debian, you need to install this package:

```bash
aptitude install kexec-tools
```

### Red Hat

On Red Hat, you need to have this package installed:

```bash
yum install kexec-tools
```

## Configuration

### Debian

There is a configuration file that allows you to make some modifications such as loading a newer kernel if one exists. This technique is not particularly recommended as some applications may not handle it well, though this is relatively rare. This option is not enabled by default, however, you can activate it if needed:

```bash {linenos=table,hl_lines=[15]}
# Defaults for kexec initscript
# sourced by /etc/init.d/kexec and /etc/init.d/kexec-load

# Load a kexec kernel (true/false)
LOAD_KEXEC=true

# Kernel and initrd image
KERNEL_IMAGE="/vmlinuz"
INITRD="/initrd.img"

# If empty, use current /proc/cmdline
APPEND=""

# Load the default kernel from grub config (true/false)
USE_GRUB_CONFIG=false
```

### Red Hat

Like Debian, Red Hat has its own configuration for kexec:

```bash
# Kernel Version string for the -kdump kernel, such as 2.6.13-1544.FC5kdump
# If no version is specified, then the init script will try to find a
# kdump kernel with the same version number as the running kernel.
KDUMP_KERNELVER=""

# The kdump commandline is the command line that needs to be passed off to
# the kdump kernel.  This will likely match the contents of the grub kernel
# line.  For example:
#   KDUMP_COMMANDLINE="ro root=LABEL=/"
# If a command line is not specified, the default will be taken from
# /proc/cmdline
KDUMP_COMMANDLINE=""

# This variable lets us append arguments to the current kdump commandline
# As taken from either KDUMP_COMMANDLINE above, or from /proc/cmdline
KDUMP_COMMANDLINE_APPEND="irqpoll nr_cpus=1 reset_devices cgroup_disable=memory"

# Any additional /sbin/mkdumprd arguments required.
MKDUMPRD_ARGS=""

# Any additional kexec arguments required.  In most situations, this should
# be left empty
#
# Example:
#   KEXEC_ARGS="--elf32-core-headers"
KEXEC_ARGS=""

#Where to find the boot image
KDUMP_BOOTDIR="/boot"

#What is the image type used for kdump
KDUMP_IMG="vmlinuz"

#What is the images extension.  Relocatable kernels don't have one
KDUMP_IMG_EXT=""
```

The default options do not need to be modified.

## Utilization

### Debian

When kexec is installed, the reboot command calls kexec and therefore natively reboots via kexec. Here are the useful commands:

* reboot: executes a fast reboot of the machine via kexec
* coldreboot: performs a standard reboot including hardware reboot

### Red Hat

For Red Hat, run this command:

```bash
kexec -l /boot/vmlinuz-`uname -r` --initrd=/boot/initramfs-`uname -r`.img --command-line="`sed 's/ rhgb\| quiet//g' /proc/cmdline`"
```

Now, to launch a fast reboot, use the standard reboot command.

I made a small script to do all of this:

```bash
#!/bin/sh
# Made by Pierre Mavro (Deimosfr) | 15/03/2012

echo "Do you want to perform a fast reboot (without hardware reboot) ? (y/n)"
read ans
if [ $ans = 'y' ] ; then
   echo "Fast reboot in progress..."
   kexec -l /boot/vmlinuz-`uname -r` --initrd=/boot/initramfs-`uname -r`.img --command-line="`sed 's/ rhgb\| quiet//g' /proc/cmdline`"
   reboot
   exit 0
fi
echo "Fast reboot cancelled"
```

## Resources
- http://wiki.debian.org/BootProcessSpeedup#Using_kexec_for_warm_reboots
- http://archive09.linux.com/feature/150202.html
- http://fedoraproject.org/wiki/Kernel/kexec
- http://fedoraproject.org/wiki/Archive:FC6KdumpKexecHowTo
