---
weight: 999
url: "/Obtenir_les_informations_hardware_d'une_machine/"
title: "Getting Hardware Information from a Machine"
description: "A guide on how to obtain hardware information from different operating systems including Linux, BSD, and Solaris."
categories: ["Linux", "Storage", "BSD", "Solaris"]
date: "2009-08-26T13:30:00+02:00"
lastmod: "2009-08-26T13:30:00+02:00"
tags: ["hardware", "dmidecode", "dmesg", "hdparm", "prtdiag", "prtconf", "Linux", "BSD", "Solaris"]
toc: true
---

## Introduction

On machines (especially remote ones), it's often useful to obtain information such as available RAM, etc.

In short, we might need other information, and I'm not sure where to put all of this data.

## Linux

### dmesg

This command shows what is written during kernel boot and all kernel calls (module loading, etc.). On some distributions, it's even written in the logs in `/var/log/*`. Here's an example output:

```bash
dmesg 
[    0.000000] Linux version 2.6.22-14-generic (buildd@king) (gcc version 4.1.3 20070929 (prerelease) (Ubuntu 4.1.2-16ubuntu2)) #1 SMP Tue Feb 12 02:46:46 UTC 2008 (Ubuntu 2.6.22-14.52-generic)
[    0.000000] Command line: root=/dev/mapper/lvm-racine ro quiet splash
[    0.000000] BIOS-provided physical RAM map:
[    0.000000]  BIOS-e820: 0000000000000000 - 000000000009fc00 (usable)
[    0.000000]  BIOS-e820: 00000000000f0000 - 0000000000100000 (reserved)
[    0.000000]  BIOS-e820: 0000000000100000 - 000000007fe0ac00 (usable)
...
```

### dmidecode

This command allows you to view RAM usage, free space, and other useful information:

```bash
# dmidecode 2.9
SMBIOS 2.3 present.
75 structures occupying 2618 bytes.
Table at 0x000F0450.

Handle 0xDA00, DMI type 218, 65 bytes
OEM-specific Type
        Header and Data:
                DA 41 00 DA B2 00 17 0B 0E 38 00 00 80 00 80 01
                00 02 80 02 80 01 00 00 A0 00 A0 01 00 58 00 58
                00 01 00 59 00 59 00 01 00 75 01 75 01 01 00 76
                01 76 01 01 00 05 80 05 80 01 00 FF FF 00 00 00
                00

Handle 0xDA01, DMI type 218, 35 bytes
OEM-specific Type
        Header and Data:
                DA 23 01 DA B2 00 17 0B 0E 38 00 10 F5 10 F5 00
                00 11 F5 11 F5 00 00 12 F5 12 F5 00 00 FF FF 00
...
```

### Hard Disks

For hard disk information, use the hdparm command to get all the details:

```bash
$ hdparm -i /dev/sda

/dev/sda:

 Model=ST3500320AS                             , FwRev=SD15    , SerialNo=            9QM89WJF
 Config={ HardSect NotMFM HdSw>15uSec Fixed DTR>10Mbs RotSpdTol>.5% }
 RawCHS=16383/16/63, TrkSize=0, SectSize=0, ECCbytes=4
 BuffType=unknown, BuffSize=0kB, MaxMultSect=16, MultSect=?8?
 CurCHS=16383/16/63, CurSects=16514064, LBA=yes, LBAsects=976773168
 IORDY=on/off, tPIO={min:120,w/IORDY:120}, tDMA={min:120,rec:120}
 PIO modes:  pio0 pio1 pio2 pio3 pio4 
 DMA modes:  mdma0 mdma1 mdma2 
 UDMA modes: udma0 udma1 udma2 udma3 udma4 udma5 *udma6 
 AdvancedPM=no WriteCache=enabled
 Drive conforms to: unknown:  ATA/ATAPI-4,5,6,7

 * signifies the current active mode
```

## BSD

### dmesg

Just like in Linux, the dmesg command exists, but the logs are not in the same place. You need to look in the `/var/run/dmesg.boot` file:

```bash
$ cat /var/run/dmesg.boot
OpenBSD 4.2 (GENERIC.MP) #1378: Tue Aug 28 10:48:58 MDT 2007
    deraadt@amd64.openbsd.org:/usr/src/sys/arch/amd64/compile/GENERIC.MP
real mem = 2146729984 (2047MB)
avail mem = 2073427968 (1977MB)
mainbus0 at root
bios0 at mainbus0: SMBIOS rev. 2.4 @ 0x7ffbc000 (62 entries)
bios0: vendor Dell Inc. version "1.5.1" date 08/10/2007
bios0: Dell Inc. PowerEdge 1950
acpi at mainbus0 not configured
ipmi0 at mainbus0: version 2.0 interface KCS iobase 0xca8/8 spacing 4
mainbus0: Intel MP Specification (Version 1.4)
cpu0 at mainbus0: apid 0 (boot processor)
cpu0: Intel(R) Xeon(R) CPU 5130 @ 2.00GHz, 1995.33 MHz
cpu0: FPU,VME,DE,PSE,TSC,MSR,PAE,MCE,CX8,APIC,SEP,MTRR,PGE,MCA,CMOV,PAT,PSE36,CFLUSH,DS,ACPI,MMX,FXSR,SSE,SSE2,SS,HTT,TM,SBF,SSE3,MWAIT,DS-CPL,VMX,TM2,CX16,xTPR,NXE,LONG
cpu0: 4MB 64b/line 16-way L2 cache
cpu0: apic clock running at 332MHz
cpu1 at mainbus0: apid 1 (application processor)
cpu1: Intel(R) Xeon(R) CPU 5130 @ 2.00GHz, 1995.02 MHz
cpu1: FPU,VME,DE,PSE,TSC,MSR,PAE,MCE,CX8,APIC,SEP,MTRR,PGE,MCA,CMOV,PAT,PSE36,CFLUSH,DS,ACPI,MMX,FXSR,SSE,SSE2,SS,HTT,TM,SBF,SSE3,MWAIT,DS-CPL,VMX,TM2,CX16,xTPR,NXE,LONG
cpu1: 4MB 64b/line 16-way L2 cache
mpbios: bus 0 is type PCI
mpbios: bus 1 is type PCI
...
```

## Solaris

### prtdiag

Here's how to obtain hardware information on Solaris:

```bash
System Configuration: Sun Microsystems     Sun Fire X4140
BIOS Configuration: American Megatrends Inc. 080014  10/15/2008
BMC Configuration: IPMI 2.0 (KCS: Keyboard Controller Style)

==== Processor Sockets ====================================

Version                          Location Tag
-------------------------------- --------------------------
Quad-Core AMD Opteron(tm) Processor 2384 CPU 1
Quad-Core AMD Opteron(tm) Processor 2384 CPU 2

==== Memory Device Sockets ================================

Type    Status Set Device Locator      Bank Locator
------- ------ --- ------------------- --------------------
unknown empty  0   DIMM0               BANK0
unknown empty  0   DIMM1               BANK1
DDR2    in use 0   DIMM2               BANK2
DDR2    in use 0   DIMM3               BANK3
DDR2    in use 0   DIMM4               BANK4
DDR2    in use 0   DIMM5               BANK5
DDR2    in use 0   DIMM6               BANK6
DDR2    in use 0   DIMM7               BANK7
unknown empty  0   DIMM8               BANK8
unknown empty  0   DIMM9               BANK9
DDR2    in use 0   DIMM10              BANK10
DDR2    in use 0   DIMM11              BANK11
DDR2    in use 0   DIMM12              BANK12
DDR2    in use 0   DIMM13              BANK13
DDR2    in use 0   DIMM14              BANK14
DDR2    in use 0   DIMM15              BANK15

==== On-Board Devices =====================================
 Gigabit Ethernet #1
 Gigabit Ethernet #2
 Gigabit Ethernet #3
 Gigabit Ethernet #4
 AST2000 VGA

==== Upgradeable Slots ====================================

ID  Status    Type             Description
--- --------- ---------------- ----------------------------
0   in use    PCI Express      PCIExp SLOT0
1   available PCI Express      PCIExp SLOT1
2   in use    PCI Express      PCIExp SLOT2
3   available PCI Express      PCIExp SLOT3
4   available PCI Express      PCIExp SLOT4
5   available PCI Express      PCIExp SLOT5
```

### prtconf

```bash
System Configuration:  Sun Microsystems  i86pc
Memory size: 73728 Megabytes
System Peripherals (Software Nodes):

i86pc
    scsi_vhci, instance #0
        disk, instance #5
        disk, instance #6
        disk, instance #7
        disk, instance #8
    isa, instance #0
        asy, instance #0
        motherboard (driver not attached)
    pci, instance #0
        pci10de,cb84 (driver not attached)
        pci10de,cb84 (driver not attached)
        pci10de,cb84 (driver not attached)
        pci10de,cb84, instance #0
            device, instance #0
                keyboard, instance #0
                mouse, instance #1
        pci10de,cb84, instance #0
            storage, instance #0
                disk, instance #2
            hub, instance #0
        pci10de,cb84, instance #0
        pci10de,cb84, instance #1
        pci10de,cb84, instance #2
        pci10de,370, instance #0
            display, instance #0
        pci10de,cb84, instance #0
        pci10de,cb84, instance #1
        pci10de,377 (driver not attached)
        pci10de,376 (driver not attached)
        pci10de,375, instance #2
            pci108e,286, instance #0
                disk, instance #3
        pci1022,1200 (driver not attached)
        pci1022,1201 (driver not attached)
        pci1022,1202 (driver not attached)
        pci1022,1203 (driver not attached)
        pci1022,1204 (driver not attached)
        pci1022,1200 (driver not attached)
        pci1022,1201 (driver not attached)
        pci1022,1202 (driver not attached)
        pci1022,1203 (driver not attached)
        pci1022,1204 (driver not attached)
    pci, instance #1
        pci10de,cb84 (driver not attached)
        pci10de,cb84 (driver not attached)
        pci10de,368 (driver not attached)
        pci10de,cb84, instance #3
        pci10de,cb84, instance #4
        pci10de,cb84, instance #5
        pci10de,cb84, instance #2
        pci10de,cb84, instance #3
        pci10de,378 (driver not attached)
        pci10de,376 (driver not attached)
        pci10de,377, instance #5
            pci1077,143, instance #0
                fp, instance #0
                    disk, instance #9
                    disk, instance #10
            pci1077,143, instance #1
                fp, instance #1
                    disk, instance #4
    iscsi, instance #0
    pseudo, instance #0
    options, instance #0
    agpgart, instance #0
    xsvc, instance #0
    objmgr, instance #0
    used-resources (driver not attached)
    cpus, instance #0
        cpu, instance #0
        cpu, instance #1
        cpu, instance #2
        cpu, instance #3
        cpu, instance #4
        cpu, instance #5
        cpu, instance #6
        cpu, instance #7
```

## Resources
- [Hardware Configuration Mastery](/pdf/maitriser_sa_configuration_materielle.pdf)
- [Dmidecode - Finding Out Hardware Details Without Opening The Computer Case](/pdf/dmidecode-_finding_out_hardware_details_without_opening_the_computer_case.pdf)
