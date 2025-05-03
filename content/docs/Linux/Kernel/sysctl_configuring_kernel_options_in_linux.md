---
weight: 999
url: "/Sysctl_\\:_configurer_les_options_kernel_sous_Linux/"
title: "Sysctl: Configuring Kernel Options in Linux"
description: "A guide to using sysctl to configure and manage Linux kernel parameters through /proc and /sys interfaces."
categories: 
  - "Linux"
  - "Debian"
date: "2012-02-12T14:26:00+02:00"
lastmod: "2012-02-12T14:26:00+02:00"
tags:
  - "Linux"
  - "Kernel"
  - "Network"
  - "System"
  - "Documentation"
toc: true
---

## Introduction

[Sysctl](https://fr.wikipedia.org/wiki/Sysctl) is an interface that allows you to examine and dynamically modify parameters of BSD and Linux operating systems. The implementation is very different between these two systems.

In Linux, the sysctl interface mechanism is also exported as part of procfs in the sys directory. This difference means that checking the value of certain parameters requires opening a file in the virtual file system, reading and interpreting its content, and then closing it. The sysctl system call exists in Linux, but is not encapsulated by a glibc function and its use is discouraged.

## Usage

First, you should know that there are two main paths: `/proc` and `/sys`:

- `/proc`: concerns memory, CPU, network, etc.
- `/sys`: concerns devices, disks, etc.

Now let's get to the point. To read a sysctl parameter, we'll use the cat command:

```bash
cat /proc/sys/net/ipv4/ip_forward
```

### /proc

`/sys` corresponds to what we call sysctls.

For persistence, after "/proc/sys/", you simply need to replace the "/" with "." to know the name of the parameter:

```
/proc/sys/net/ipv4/ip_forward
```

corresponds to 

```
net.ipv4.ip_forward
```

#### Applying values on the fly

##### Method 1

Here the default value is 0. If I want to enable IP forwarding on the fly, here's how to do it:

```bash
echo 1 > /proc/sys/net/ipv4/ip_forward
```

This will activate the new value but in a non-persistent way (value reset after reboot).

##### Method 2

To apply on the fly:

```bash
sysctl -w net.ipv4.ip_forward=1
```

#### Applying values persistently

Simply add this to the `/etc/sysctl.conf` file:

```bash
net.ipv4.ip_forward = 1
```

If you haven't yet activated your new parameter, you can activate everything in the `/etc/sysctl.conf` file like this:

```bash
sysctl -p
```

### /sys

For `/sys`, it's unfortunately less straightforward than for `/proc` since only the `/etc/rc.local` file exists for persistence.

## Documentation

### Debian

To get the kernel documentation, install this:

```bash
aptitude linux-doc-*
```

The documentation is located in `/usr/share/doc/linux-doc*/Documentation`

### RedHat

To get the kernel documentation, install this:

```bash
aptitude kernel-doc
```

The documentation is located in `/usr/share/doc/kernel-doc*/Documentation`
