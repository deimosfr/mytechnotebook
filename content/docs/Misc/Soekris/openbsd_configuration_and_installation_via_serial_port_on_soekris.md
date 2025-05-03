---
weight: 999
url: "/Configuration_et_installation_via_port_série_d'OpenBSD_sur_Soekris/"
title: "OpenBSD Configuration and Installation via Serial Port on Soekris"
description: "Complete guide to install OpenBSD on a Soekris device using the serial port and optimizing for CompactFlash storage"
categories: ["BSD", "OpenBSD", "Soekris", "Embedded"]
date: "2009-10-25T20:14:00+02:00"
lastmod: "2009-10-25T20:14:00+02:00"
tags: ["OpenBSD", "Soekris", "serial", "embedded", "PXE", "CompactFlash"]
toc: true
---

## Introduction

You've just bought a Soekris board without thinking too much about it... don't worry, it works like a charm :-).

I'm going to explain how I was able to put OpenBSD on my Soekris with a CompactFlash drive.

## Materials Used

For my installation, I needed:

- A DB9 Female/Female serial cable (also called NULL-MODEM)
- A USB to serial port cable (I no longer have serial ports on my machines)
- A network cable connected to interface 1 (the first one) of the Soekris for the PXE boot to work
- A Compact Flash card (Kingston Elite Pro 16GB 133X) (SanDisk is preferred for compatibility reasons)
- And finally, the Soekris 5501-70

## Prerequisites

Before starting, certain things need to be set up. I won't explain how all the services below work, so I'll simply provide my configurations. However, I invite you to check out:

- [DHCP3: Installation and Configuration of a DHCP Server](/DHCP3_:_Installation_et_configuration_d'un_serveur_DHCP/)
- [TFTP: PXE Server, OS Deployment under Linux](/TFTP_:_PXE_Serveur,_déploiement_d'OS_sous_Linux/)
- <!-- Link to "Connecting via Serial Port on Debian" removed - file not found -->

_My DHCP and TFTP server has the IP address 192.168.10.107 here._

### DHCP Server

Here's my DHCP configuration:

```bash
ddns-update-style ad-hoc;
option domain-name-servers 212.27.40.241, 212.27.40.240;
option routers 192.168.10.138;
log-facility local7;
subnet 192.168.10.0 netmask 255.255.255.0 {
 range 192.168.10.70 192.168.10.80;
 filename "pxeboot";
 next-server 192.168.10.107;
 option root-path "/var/lib/tftpboot";
}
```

Don't forget to restart the DHCP server after making changes.

### TFTP Server

For TFTP, go to /var/lib/tftpboot and do this:

```bash
cd /var/lib/tftpboot
wget http://ftp.arcane-networks.fr/pub/OpenBSD/snapshots/i386/bsd.rd
wget http://ftp.arcane-networks.fr/pub/OpenBSD/snapshots/i386/pxeboot
mkdir etc
touch etc/boot.conf
chmod -Rf 777 .
```

Then edit the boot.conf file and insert these lines:

```bash
set tty com0
stty com0 19200
boot bsd.rd
```

Don't forget to restart your TFTP server.

### Soekris BIOS

#### Connection via Minicom

Launch minicom:

```bash
$ minicom -s
```

Then go to:

- Serial port configuration
- Baud rate/Parity/Bits
- Set it to 19200 8N1 (key combination c+a+q)

Save everything and exit to validate the configuration. Now you're connected.

#### Boot

Let's configure the Soekris BIOS. Use your com port to connect to it and press Ctrl+P at boot to enter the BIOS:

```
POST: 012345689bcefghips1234ajklnopqr,,,tvwxy

comBIOS ver. 1.33  20070103  Copyright (C) 2000-2007 Soekris Engineering.

net5501

0512 Mbyte Memory                        CPU Geode LX 500 Mhz

Pri Mas  ELITE PRO CF CARD 16GB          LBA Xlt 1024-255-63  15761 Mbyte

Slot   Vend Dev  ClassRev Cmd  Stat CL LT HT  Base1    Base2   Int
-------------------------------------------------------------------
0:01:2 1022 2082 10100000 0006 0220 08 00 00 A0000000 00000000 10
0:06:0 1106 3053 02000096 0117 0210 08 40 00 0000E101 A0004000 11
0:07:0 1106 3053 02000096 0117 0210 08 40 00 0000E201 A0004100 05
0:08:0 1106 3053 02000096 0117 0210 08 40 00 0000E301 A0004200 09
0:09:0 1106 3053 02000096 0117 0210 08 40 00 0000E401 A0004300 12
0:14:0 104C AC23 06040002 0107 0210 08 40 01 00000000 00000000
0:20:0 1022 2090 06010003 0009 02A0 08 40 80 00006001 00006101
0:20:2 1022 209A 01018001 0005 02A0 08 00 00 00000000 00000000
0:21:0 1022 2094 0C031002 0006 0230 08 00 80 A0005000 00000000 15
0:21:1 1022 2095 0C032002 0006 0230 08 00 00 A0006000 00000000 15
1:00:0 100B 0020 02000000 0107 0290 00 40 00 0000D001 A4000000 10
1:01:0 100B 0020 02000000 0107 0290 00 40 00 0000D101 A4001000 07
1:02:0 100B 0020 02000000 0107 0290 00 40 00 0000D201 A4002000 10
1:03:0 100B 0020 02000000 0107 0290 00 40 00 0000D301 A4003000 07

 5 Seconds to automatic boot.   Press Ctrl-P for entering Monitor.

comBIOS Monitor.   Press ? for help.
```

Once inside, set the date and time:

```bash
time HH:MM:SS
date YYYY/MM/DD
```

Now for some brief explanations about the boot process. Enter the BIOS again and use the show command to see the available options:

```bash
> show

ConSpeed = 19200
ConLock = Enabled
ConMute = Disabled
BIOSentry = Enabled
PCIROMS = Enabled
PXEBoot = Enabled
FLASH = Primary
BootDelay = 5
FastBoot = Disabled
BootPartition = Disabled
BootDrive = 80 81 F0 FF
ShowPCI = Enabled
Reset = Hard
CpuSpeed = Default
```

The BootDrive devices are indicated as follows:

- 80: hard disk (IDE or SATA)
- 81: compact flash
- F0: PXE

We'll boot from PXE to launch the OpenBSD installation:

```bash
boot F0
```

## OpenBSD Installation

Perform your installation as you normally would, except it would be good not to set up a swap partition to avoid stressing the Compact Flash.

Then, toward the end of the installation, don't forget to specify that you also want to use the com port to connect:

```
Change the default console to com0? [no] yes
Available speeds are: 9600 19200 38400 57600 115200.
Which one should com0 use? (or 'done') [9600] 19200
Saving configuration files...done.
```

Then reboot, and once again we'll touch the BIOS one last time to specify the boot order:

```bash
set BootDrive=81 80 F0 FF
reboot
```

There are subtleties described below due to the short lifespan of compact flash. We'll do everything possible to preserve it as much as possible.

### Remove Access Information

Let's remove access information on the compact flash with the noatime option in fstab:

```bash
/dev/wd0a / ffs rw,noatime 1 1
/dev/wd0d /var ffs rw,nodev,nosuid,noatime 1 2
```

### Add MFS Filesystems

Add MFS filesystems to your data that changes regularly:  
[MFS: Using a Filesystem in RAM]({{< ref "docs/BSD/Filesystems/mfs_using_a_ram_filesystem.md" >}})

### Watchdog

I strongly recommend enabling the watchdog:  
[Setting up a Watchdog](/Mise_en_place_d'un_Watchdog/)

## Updating OpenBSD

For updates, simply specify during installation that you want to upgrade the system. Then, after rebooting, download the etcXX.tgz file of the OpenBSD version you're upgrading to in /tmp and run this command:

```bash
sysmerge -s etcXX.tgz
```

You'll then see menus for merging your configurations.

## FAQ

### Why doesn't the OS boot on first attempt? I have to do a cold boot for it to work.

I had messages like these:

```
open(hd0a:/etc/boot.conf): Unknown error: code 102
boot>
booting hd0a:/bsd: open hd0a:/bsd: Unknown error: code 102
 failed(102). will try /bsd
boot>
booting hd0a:/bsd: open hd0a:/bsd: Invalid argument
 failed(22). will try /bsd
Turning timeout off.
boot>
```

Forced to reset at the Monitor (BIOS) level, and therefore do a Cold boot for the OS to boot, what a joke...

To solve the problem, there are 2 solutions:

- <!-- Link to "Update the Soekris BIOS" removed - file not found -->
- Change your Flash card to a more compatible one (SanDisk for example, as indicated on the Soekris site)

## References

http://ludique.u-bourgogne.fr/~leclercq/wiki/index.php/Soekris  
[Installing OpenBSD on Soekris via QEMU](https://blog.arofarn.info/openbsd/installation-sur-soekris-net4801-50-grace-a-qemu/)  
http://www.lininfo.org/spip.php?article11  
http://wiki.gcu.info/doku.php?id=openbsd:install_soekris
