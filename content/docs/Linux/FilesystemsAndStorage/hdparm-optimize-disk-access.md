---
weight: 999
url: "/Hdparm_\\:_Optimiser_les_accès_disques/"
title: "Hdparm: Optimizing Disk Access"
description: "How to optimize disk and optical drive performance on Linux systems using hdparm utility."
categories: ["Linux", "System Administration", "Performance"]
date: "2006-10-18T09:14:00+02:00"
lastmod: "2006-10-18T09:14:00+02:00"
tags: ["hdparm", "DMA", "Linux", "Performance Tuning", "Disk I/O"]
toc: true
---

## Introduction

When you copy large files from one hard drive to another or copy the contents of a CD-ROM to a hard drive, you may have noticed a significant slowdown in your system. During the transfer, music might begin to crackle, for example, or DVD-ROM playback is particularly slow.

* Considerations

Throughout this article, we consider that your hard drive is the master on the first IDE interface, meaning it is connected to the entry point `/dev/hda`.

## Prerequisites

### Kernel configuration

Your kernel must support DMA. You must have the following options when compiling your kernel:

{{< table "table-hover table-striped" >}}
| Section | Kernel Option | As module or built-in? |
|---------|---------------|------------------------|
| ATA/IDE/MFM/RLL support | IDE, ATA and ATAPI Block devices / Generic PCI IDE chipset support | Built-in |
| ATA/IDE/MFM/RLL support | IDE, ATA and ATAPI Block devices / Generic PCI bus-master DMA support | Built-in |
| ATA/IDE/MFM/RLL support | IDE, ATA and ATAPI Block devices / Use PCI DMA by default when available | Built-in |
{{< /table >}}

You must replace XXXXXXXX with the chipset reference of your motherboard. Refer to your hardware user manual to find this reference.

### The hdparm package

The tool we will use to test and optimize the hard disk transfer rate is called hdparm, which corresponds to the package of the same name. We install this package:

```bash
apt-get install hdparm
```

Note that certain options must be activated in the kernel to enable the DMA channel of your IDE devices. All kernels available for Debian GNU/Linux have these options enabled, but if it's a kernel you've compiled yourself, it's better to check that the following options are present:

## Improving the transfer rate of your hard drives

To check the transfer rate of your hard drive, simply type the following command:

```bash
hdparm -tT /dev/hda
```

Without optimization, you should get something similar to this:

```
/dev/hda:
Timing buffer-cache reads: 128 MB in 1.06 seconds = 120.75 MB/sec
Timing buffered disk reads: 64 MB in 35.70 seconds = 1.79 MB/sec
```

The speed of a hard drive is generally between 10 and 30 MB/s for the second test. You can see that here the hard drive is horribly slow. We will therefore fix this problem by activating the DMA controller and 32-bit transfer for your hard drive. The DMA controller (acronym for Direct Memory Access) is a process that allows access to RAM without going through the processor. You can activate this option without any worries, using hdparm.

To activate this optimization:

```bash
hdparm -c1 -d1 /dev/hda
```

Which produces the following result:

```
/dev/hda:
setting 32-bit I/O support flag to 1
setting using_dma to 1 (on)
I/O support = 1 (32-bit)
using_dma = 1 (on)
```

In the above command:
* -c1 corresponds to activating 32-bit transfer
* -d1 corresponds to activating the DMA channel

You can test your hard drive again to verify that the optimization does produce a performance gain. The rate is on average multiplied by 15. However, this value can vary depending on your hardware!

## Improving the transfer rate of your CD-ROM or DVD-ROM drive

The optimization of a CD-ROM or DVD-ROM drive can be done regardless of the drivers you use to manage your CD-ROM drives. Thus, for CD-ROM to CD-ROM copies, SCSI emulation is absolutely essential. So, whether you have SCSI emulation or not, you should type the following command:

```bash
hdparm -c1 -d1 /dev/hdc
```

You should then get the following result:

```
/dev/hda: setting 32-bit I/O support flag to 1
setting using_dma to 1 (on)
I/O support = 1 (32-bit)
using_dma = 1 (on)
```

Normally your transfer rate should have been multiplied by 2 and you can notice that DVD-ROM playback is much smoother. In addition, to reduce the noise made by the CD-ROM or DVD-ROM drive, you can choose its reading speed with this command (where 40 corresponds to the chosen speed, i.e., 40X):

```bash
hdparm -E 40 /dev/hdc
```

## Making your optimizations permanent

The optimizations you just made are certainly interesting, but at the next restart, you'll have to do everything again. To overcome this problem, we will write them in the configuration file of the hdparm program.

You need to edit the `/etc/hdparm.conf` file. This file contains in the first part all the options you can use. You then need to define for each of your disks the list of options you want to activate.

The following block activates DMA and 32-bit access for the `/dev/hda` disk.

```
/dev/hda {
 quiet
 dma = on
 io32_support = 1
}
```

* quiet parameter

The quiet parameter makes the modification of the hard disk properties silent. Without this parameter, you will get information in the console about the status of modifications made to the hard disk.

If you have a CD-ROM drive, you can be inspired by the block below:

```
/dev/hdc {
 quiet
 dma = on
 io32_support = 1
 cd_speed = 40
}
```

To activate these changes immediately, you can execute the command:

```bash
/etc/init.d/hdparm start
```

## Appendix: parameters for the /etc/hdparm.conf file

Here are the first lines of the `/etc/hdparm.conf` file that describe the different possible options for the blocks you can define for each of your disks.

```
# -q be quiet
#quiet
# -a sector count for filesystem read-ahead
#read_ahead_sect = 12
# -A disable/enable the IDE drive's read-lookahead feature
#lookahead = on
# -b bus state
#bus = on
# -c enable (E)IDE 32-bit I/O support - can be any of 0,1,3
#io32_support = 1
# -d disable/enable the "using_dma" flag for this drive
#dma = on
# -D enable/disable the on-drive defect management
#defect_mana = off
# -E cdrom speed
#cd_speed = 40
# -m sector count for multiple sector I/O
#mult_sect_io = 32
# -P maximum sector count for the drive's internal prefetch mechanism
#prefetch_sect = 12
# -r read-only flag for device
#read_only = off
# -S standby (spindown) timeout for the drive
#spindown_time = 24
# -u interrupt-unmask flag for the drive
#interrupt_unmask = on
# -W Disable/enable the IDE drive's write-caching feature
#write_cache = off
# -X IDE transfer mode for newer (E)IDE/ATA2 drives
#transfer_mode = 34
# -y force to immediately enter the standby mode
#standby
# -Y force to immediately enter the sleep mode
#sleep
# -Z Disable the power-saving function of certain Seagate drives
#disable_seagate
# -M Set the acoustic management properties of a drive
#acoustic_management
```
