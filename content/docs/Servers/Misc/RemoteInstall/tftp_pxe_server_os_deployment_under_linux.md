---
weight: 999
url: "/TFTP_\\:_PXE_Serveur,_déploiement_d'OS_sous_Linux/"
title: "TFTP: PXE Server, OS Deployment under Linux"
description: "Guide to set up a PXE server for OS deployment using TFTP under Linux"
categories: ["Debian", "Linux", "Ubuntu"]
date: "2010-05-20T05:14:00+02:00"
lastmod: "2010-05-20T05:14:00+02:00"
tags:
  [
    "TFTP",
    "PXE",
    "Network",
    "Servers",
    "Debian",
    "OpenBSD",
    "Red Hat",
    "Memtest86+",
  ]
toc: true
---

## Introduction

PXE boot (Pre-boot eXecution Environment) allows a workstation to boot from the network an operating system that is stored on a server.

It also allows automatic and remote installation of servers with various operating systems.

To enable PXE, you first need to configure it in the BIOS. The option is frequently found in a menu related to the network card.

PXE booting is performed in several steps:

- Search for an IP address on a DHCP server as well as the file to boot
- Download the boot file from a Trivial FTP server
- Execute the boot file

It should be noted that the size of the boot file does not allow for directly booting a Linux kernel, for example, but requires that the boot software download and execute it itself.

## Prerequisites

The prerequisites are quite simple; you just need a [DHCP server](./DHCP3_:_Installation_et_configuration_d'un_serveur_DHCP.html) that is able to boot on PXE. We will see here the configuration of this DHCP server so that it accepts network booting.

## Installation

To install the PXE server:

```bash
apt-get install tftpd-hpa syslinux
```

## Configuration

### tftpd-hpa

We will edit the file `/etc/default/tftpd-hpa` to replace the value of RUN_DAEMON:

```bash
RUN_DAEMON="yes"
```

### Inetd

We disable the tsize of tftp-hpa which limits the size of files to be downloaded. For this, add a line in `/etc/inetd.conf` and check that another one is commented out:

```bash
tftp            dgram   udp     wait    root  /usr/sbin/in.tftpd /usr/sbin/in.tftpd -s /var/lib/tftpboot
# tftp            dgram   udp     wait    nobody    /usr/sbin/tcpd /usr/sbin/in.tftpd -r blksize /tftpboot
```

Once done, we will restart inetd and tftpd:

```bash
/etc/init.d/inetd restart
/etc/init.d/tftpd-hpa start
```

To verify that everything is working:

```bash
$ netstat -uap | grep tftp
udp        0      0 *:tftp                  *:*                                30265/in.tftpd
```

If the line above appears, everything went well :-)

### Iptables

Here's the nice line to add to iptables to allow tftp:

```bash
iptables -A INPUT -s 10.1.1.0/255.255.255.0 -p udp -j ACCEPT
```

### DHCP under Linux

If your DHCP is under Linux, edit the `/etc/dhcp3/dhcpd.conf` file and add these lines in your subnet:

```bash
subnet 192.168.0.0 netmask 255.255.255.0 {
...
filename "pxelinux.0";
next-server 192.168.1.254;
...
}
```

Next-server is to specify the IP address of the PXE server.  
Then restart your DHCP server:

```bash
/etc/init.d/dhcp3 restart
```

### DHCP under Windows

If your DHCP is under Windows, in your DHCP configuration (general or not), add the address of the TFTP server.

### Boot loader

Now, we must prepare and organize our TFTP server:

```bash
cd /var/lib/tftpboot
mkdir pxelinux.cfg os-installer
touch boot.txt
cp /usr/lib/syslinux/{pxelinux.0,menu.c32} .
```

We have inserted pxelinux.0, which is essential for booting our OSes, and menu.c32, which provides a basic but practical menu when we have our OSes installed.

Let's configure the global configuration of the server. Create and edit the file pxelinux.cfg/default to insert this:

```bash
PROMPT 1
DISPLAY boot.txt

F1 boot-screens/f1.txt
F2 boot-screens/f2.txt
F3 boot-screens/f3.txt
F4 boot-screens/f4.txt
F5 boot-screens/f5.txt
F6 boot-screens/f6.txt
F7 boot-screens/f7.txt
F8 boot-screens/f8.txt
F9 boot-screens/f9.txt
F0 boot-screens/f10.txt

# On définit ce qui sera lancer par defaut lors du boot, à savoir le menu graphique choisi
DEFAULT menu.c32
NOESCAPE 1
# On choisi un titre pour l'écran d'arrivé
MENU TITLE -=[ TFTP Server - OS Installer ]=-
# Il y a un boot automatique au bout de 20 secondes
TIMEOUT  200

# Le boot automatique s'effectue sur le disque dur en locale
LABEL Local Hard Drive Boot
        localboot 0  --
```

The basic configuration is now ready. We only need to add operating systems.

## Setting up Operating Systems

Let's see how to set up different types of operating systems. Before continuing, go to this folder:

```bash
cd os-installer
```

### Debian

Let's create what we need, that is, a folder, and then insert the kernel. We'll do both the 32-bit and 64-bit versions:

```bash
mkdir -p debian-installer/{x64,x86}
```

For the 32-bit version:

```bash
cd debian-installer/x86
wget http://ftp.debian.org/dists/stable/main/installer-i386/current/images/netboot/debian-installer/i386/initrd.gz \
http://ftp.debian.org/dists/stable/main/installer-i386/current/images/netboot/debian-installer/i386/linux
```

For the 64-bit version:

```bash
cd debian-installer/x64
wget http://ftp.debian.org/dists/stable/main/installer-amd64/current/images/netboot/debian-installer/amd64/initrd.gz \
http://ftp.debian.org/dists/stable/main/installer-amd64/current/images/netboot/debian-installer/amd64/linux
```

Then add these lines (depending on the architecture you have chosen) in the file `/var/lib/tftpboot/pxelinux.cfg/default`:

```bash
LABEL x64 - Debian
        kernel os-installer/debian-installer/x86/linux
        append vga=791 priority=low initrd=os-installer/debian-installer/x86/initrd.gz --

LABEL x86 - Debian
        kernel os-installer/debian-installer/x64/linux
        append vga=791 priority=low initrd=os-installer/debian-installer/x64/initrd.gz --
```

- vga=791: loads 1024\*768 resolution
- priority=low: loads Debian expert mode

Note: To automate installations, [follow this link: Automate a Debian installation]({{< relref "docs/Linux/Misc/Debian/automate-debian-installation.md" >}}).

### Memtest86+

At the time of writing, the latest version is 1.70. So I'll use this for my example:

```bash
mkdir -p memtest86
```

Let's download this version (we'll take the bootable binary):

```bash
cd memtest86
wget http://www.memtest.org/download/1.70/memtest86+-1.70.bin.gz
gzip -d http://www.memtest.org/download/1.70/memtest86+-1.70.bin.gz
```

Then a small subtlety, we need to rename and remove the .bin for it to work:

```bash
mv memtest86+-1.70{,.bin}
```

Then add these lines (depending on the architecture you have chosen) in the file `/var/lib/tftpboot/pxelinux.cfg/default`:

```bash
LABEL Memtest86+ (RAM Testing Program)
        kernel os-installer/memtest/memtest86+-1.70
```

### OpenBSD

Again, we'll do what's necessary to be able to launch OpenBSD in 32-bit and 64-bit versions:

```bash
mkdir -p openbsd-installer/{x64,x86}
```

For the 32-bit version:

```bash
cd openbsd-installer/x86
wget http://ftp.arcane-networks.fr/pub/OpenBSD/4.1/i386/floppy41.fs
```

For the 64-bit version:

```bash
cd openbsd-installer/x64
wget http://ftp.arcane-networks.fr/pub/OpenBSD/4.1/amd64/floppy41.fs
```

We're using the floppy versions here and not the CD versions because we'll be using a new module called memdisk that can load an ISO but only smaller than the size of a floppy disk. So copy this module:

```bash
cp /usr/lib/syslinux/memdisk /var/lib/tftpboot/
```

Then add these lines (depending on the architecture you have chosen) in the file `/var/lib/tftpboot/pxelinux.cfg/default`:

```bash
LABEL x64 - OpenBSD 4.1
        kernel memdisk
        append initrd=x64/openbsd-installer/floppy41.fs  --

LABEL x86 - OpenBSD 4.1
        kernel memdisk
        append initrd=x86/openbsd-installer/floppy41.fs  --
```

### Red Hat

Red Hat is a bit special because we'll need to create a DVD, then copy it to insert the kernel. We'll do the 32-bit and 64-bit versions:

```bash
mkdir -p redhat-installer/{x64,x86}
```

[Create the DVD]({{< ref "docs/Linux/Packages/RedHat/redhat-dvd-repository.md" >}}), then copy it to the proper directory according to your version (32 or 64 bits).

For the 32-bit version:

```bash
cd redhat-installer/x86
cp -Rf votre_dvd/* votre_dvd/.* .
```

For the 64-bit version:

```bash
cd redhat-installer/x64
cp -Rf votre_dvd/* votre_dvd/.* .
```

Then add these lines (depending on the architecture you have chosen) in the file `/var/lib/tftpboot/pxelinux.cfg/default`:

```bash
LABEL x64 - Red Hat
        kernel os-installer/redhat-installer/x86/linux
        append vga=791 priority=low initrd=os-installer/redhat-installer/x86/initrd.gz --

LABEL x86 - Red Hat
        kernel os-installer/redhat-installer/x64/linux
        append vga=791 priority=low initrd=os-installer/redhat-installer/x64/initrd.gz --
```

- vga=791: loads 1024\*768 resolution

## Password Protection

The SYSLINUX archive contains an executable called sha1pass (it's a Perl script) that generates passwords in the correct format. To use it under Debian, you need the appropriate Perl module:

```bash
apt-get install libdigest-sha1-perl
```

Then run the command with the password as a parameter and it will give us the string to paste into the configuration file. For example, to protect Ghost:

```bash
LABEL ghost
        MENU LABEL Ghost
        MENU PASSWD $4$jfoBirJg$rSDbzznCZtmJAES9RH/lC92/3Rs$
        kernel memdisk
        append initrd=ghost/ghost288.IMA
```

## Resources

- [TFTP Documentation on Ubuntu](/pdf/pxe_ubuntu.pdf)
- [TFTP Documentation on Debian](/pdf/booting_on_pxe.pdf)
- [Start without a Disk with PXE, Grub and NFS](/pdf/démarrez_sans_disque_avec_pxe,_grub_et_nfs.pdf)
- [Configuration and Installation via Serial Port of OpenBSD on Soekris](Configuration_et_installation_via_port_série_d'OpenBSD_sur_Soekris.html)
- [Boot BSDs with PXElinux](https://www.thegibson.org/blog/archives/10)
- [Setting up a PXE and Solving its Problems](https://en.opensuse.org/SuSE_install_with_PXE_boot)
