---
weight: 999
url: "/Recompiler_son_du_noyau_(Kernel)/"
title: "Recompile your kernel"
description: "A guide on how to recompile the Linux kernel on various distributions including CentOS, Debian, Fedora, Mandriva, Suse, and Ubuntu"
categories: ["Debian", "CentOS", "Linux"]
date: "2016-11-10T23:54:00+02:00"
lastmod: "2016-11-10T23:54:00+02:00"
tags: ["kernel", "compilation", "development", "linux", "debian", "ubuntu", "centos"]
toc: true
---

## Linux

### CentOS

[CentOS Kernel 2.6](/pdf/centos-kernel.pdf)

### Debian

Here is a list of packages needed to recompile the kernel:

```bash
aptitude install bzip2 libncurses5-dev fakeroot kernel-package
```

Let's go to `/usr/src`:

```bash
cd /usr/src
```

Then, go to www.kernel.org and download the latest version in "**Full**" of the latest kernel. Here, it's version 4.8.4, then we extract it:

```bash
wget https://cdn.kernel.org/pub/linux/kernel/v4.x/linux-4.8.4.tar.xz
tar -xJf linux-4.8.4.tar.xz
```

Now let's create a symbolic link:

```bash
ln -s linux-4.8.4 linux
```

The kernel is ready to be configured. Let's launch the configuration tool:

```bash
cd linux
make menuconfig
```

Or copy the configuration of your existing kernel:

```bash
cp /boot/config-$(uname -r) .config
```

To avoid a "x509_certificate" error during compilation, we will disable the kernel signature (we don't have the signature key, it belongs to Debian):

```bash
sed -i 's/^CONFIG_SYSTEM_TRUSTED_KEY/# CONFIG_SYSTEM_TRUSTED_KEY/g' .config
sed -i 's/^CONFIG_MODULE_SIG_KEY/# CONFIG_MODULE_SIG_KEY/g' .config
```

It's possible to use all of your cores to speed up compilation time:

```bash
export CONCURRENCY_LEVEL=`grep -c "^processor" /proc/cpuinfo`
```

#### New method

Once configured, all that's left is to launch the compilation:

```bash
make clean
make deb-pkg LOCALVERSION=-custom KDEB_PKGVERSION=$(make kernelversion)-1 -j $CONCURRENCY_LEVEL
```

You can change the name of LOCALVERSION to a name that suits you better and increment KDEB_PKGVERSION each time you compile.

#### Old method

Once configured, all that's left is to launch the compilation:

```bash
make-kpkg clean
make-kpkg --initrd --revision=1.0 kernel_image
```

or without initrd:

```bash
make-kpkg --revision=1.0 kernel_image
```

The revision tag is used to put a version number on your kernel. That way, if during the next boot you get a kernel panic, restart with the old one and start over. When recompiling it, increment the version by 1 (ex: --revision=2.0).

Your kernel is now finished, let's install it:

```bash
dpkg -i ../linux-image-4.8.4_1.0_amd64.deb
```

Now restart your machine and boot on your new kernel :-)

Here are other documentations:  
[Debian Kernel 2.4](/pdf/debian-kernel-2.4.pdf)  
[Debian Kernel 2.6](/pdf/debian-kernel-2.6.pdf)

### Fedora Core

[Fedora Kernel](/pdf/fedora-kernel.pdf)

### Mandriva

[Mandriva Kernel](/pdf/mandriva-kernel.pdf)

### Suse

[Suse Kernel](/pdf/suse-kernel.pdf)

### Ubuntu

[Ubuntu Kernel](/pdf/ubuntu-kernel.pdf)
