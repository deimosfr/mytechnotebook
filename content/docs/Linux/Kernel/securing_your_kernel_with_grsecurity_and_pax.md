---
weight: 999
url: "/Sécurisation_de_son_noyau_avec_Grsecurity_et_PaX/"
title: "Securing Your Kernel with Grsecurity and PaX"
description: "How to secure a Linux kernel with Grsecurity and PaX patches to enhance security against buffer overflows and other vulnerabilities."
categories: ["Linux", "Security"]
date: "2009-02-03T07:18:00+02:00"
lastmod: "2009-02-03T07:18:00+02:00"
tags: ["grsecurity", "pax", "kernel", "linux", "security", "hardening"]
toc: true
---

## Grsecurity & PaX

### Introduction

Are buffer overflows and script kiddies something you're concerned about? To counter these types of problems, there exists a kernel patch. There are many solutions to protect against "classic" attacks which are complex for professionals (but not impossible). Here is THE kernel patch that helps us protect our systems. There's even a learning phase to help with system configuration, resulting in a highly secured solution.

http://www.grsecurity.net  
http://pax.grsecurity.net  
http://www.kernel.org

### Installation of Required Packages

Here is the list of packages that need to be installed **before applying the patch** and recompiling the kernel.

```bash
apt-get install paxctl paxtest chpax gradm2
```

### Patching and Recompiling the Kernel

Now we need to download the kernel and the Grsecurity patch. We'll decompress the kernel and then the patch:

```bash
wget http://www.kernel.org/pub/linux/kernel/v2.6/linux-2.6.*.tar.bz2
tar -xjvf linux-2.6.*.tar.bz2
wget http://www.grsecurity.net/grsecurity-*.patch.gz
gzip -d grsecurity-*.patch.gz
```

What remains is to patch the kernel, make any modifications as you see fit (or following the [Quickstart](https://www.grsecurity.net/quickstart.pdf) documentation which is the official installation guide, or find more explanations [on the site](https://grsecurity.net/confighelp.php)), then create and install the kernel:

```bash
ln -s linux-2.6.* linux
cd linux
patch -p1 < ../grsecurity-*.patch
make menuconfig
make-kpkg clean
make-kpkg --revision=1.0 kernel_image
dpkg -i ../linux-2.6.*.deb
```

### System Configuration

Good, we are ready to reboot the machine. We don't need to recompile certain software as mentioned in the documentation, such as gradm, since they exist as packages (thanks to Debian). So it's easier to install and allows automatic updates when upgrading the distribution.

Here are some commands that will be useful now:

```bash
# Start the system: 
gradm -E
# Stop the system:
gradm -D
# Authenticate with administrator rights:
gradm –a admin
# De-authenticate from administrator rights, otherwise exit shell or enter:
gradm –u admin
```

**Now make sure that all your services are working, that everything has started correctly! If not, consult the [FAQ](#faq)**

Now that everything is working properly, we will switch to learning mode. This mode will learn everything you do on your system to normalize actions. Once you think you've done everything that you consider "normal," we will disable this learning mode to go into production.
**However, do not perform any administrative tasks such as starting/stopping daemons, adding/removing users/software**

To start this learning mode, here is the command:

```bash
gradm –F –L /etc/grsec/learning.log
```

Once learning is complete (**leave it for at least 2 days**), we can stop it:

```bash
gradm –F –L /etc/grsec/learning.log –O etc/grsec/acl
```

That's it, all you need to do now is keep an eye on the Gradm logs to see the types of protections your kernel is implementing.

## FAQ

### How to Debug?

To debug and view memory properties at the application level:

```bash
strace /application
```

And to read memory at the library level:

```bash
readelf -e /usr/lib/library
```

### error while loading shared libraries Permission denied

Yes, not everything works. You may encounter this type of error. To check the status of an application and how it will be handled, run the chpax command:

```bash
chpax --help               
chpax 0.7 .::. Manage PaX flags for binaries
Usage: chpax OPTIONS FILE1 FILE2 FILEN ...
  -P    enforce paging based non-executable pages
  -p    do not enforce paging based non-executable pages
  -E    emulate trampolines
  -e    do not emulate trampolines
  -M    restrict mprotect()
  -m    do not restrict mprotect()
  -R    randomize mmap() base [ELF only]
  -r    do not randomize mmap() base [ELF only]
  -X    randomize ET_EXEC base [ELF only]
  -x    do not randomize ET_EXEC base [ELF only]
  -S    enforce segmentation based non-executable pages
  -s    do not enforce segmentation based non-executable pages
  -v    view current flag mask 
  -z    zero flag mask (next flags still apply)

The flags only have effect when running the patched Linux kernel.
```

Next, try to see what is wrong with the application:

```bash
chpax -v /usr/sbin/openvpn 

----[ chpax 0.7 : Current flags for /usr/sbin/openvpn (PemRxS) ]---- 

 * Paging based PAGE_EXEC       : enabled (overridden) 
 * Trampolines                  : not emulated 
 * mprotect()                   : restricted 
 * mmap() base                  : randomized 
 * ET_EXEC base                 : not randomized 
 * Segmentation based PAGE_EXEC : enabled
```

Here, we can see that there is a restriction in OpenVPN at the mprotect level. We can work around this by removing the protection. But this means we're allowing a vulnerability here. I'll disable it like this:

```bash
chpax -m /usr/sbin/openvpn
```

To finally have:

```bash
* mprotect()                   : not restricted
```

At this point, my application will work without any issues :-)

### Bind Refuses to Start: permission denied

This happens due to a kernel module that you either haven't compiled or have enabled as a module. It needs to be implemented in the kernel.

```bash
CONFIG_SECURITY_CAPABILITIES=y
```

## Resources
- [Le fonctionnement de PaX](/pdf/le_fonctionnement_de_pax.pdf)
- [Hardening The Linux Kernel With Grsecurity](/pdf/hardening_the_linux_kernel_with_grsecurity.pdf)
