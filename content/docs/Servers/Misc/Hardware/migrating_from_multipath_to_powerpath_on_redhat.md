---
weight: 999
url: "/Migrer_de_Multipath_à_Powerpath_sur_RedHat/"
title: "Migrating from Multipath to Powerpath on RedHat"
description: "A guide on how to migrate from Multipath to Powerpath on RedHat Linux systems, including uninstallation of multipath, installation of Powerpath, verification, and LVM configuration."
categories: 
  - Linux
date: "2008-06-03T12:10:00+02:00"
lastmod: "2008-06-03T12:10:00+02:00"
tags: 
  - RedHat
  - Storage
  - EMC
  - Powerpath
  - LVM
toc: true
---

## Introduction

Powerpath is the multipathing solution for EMC arrays. The multipath package is so buggy on RedHat that you shouldn't install it in production environments. This migration has been released on RedHat 4.6.EL.

Reminder: Multipathing brings redundancy functionalities with 2 links on a disk array without having I/O errors.

## Uninstalling multipath

Multipath package name (if installed) can be found this way:

```bash
rpm -qa | grep multipath
```

Next, you just need to use the package name and add it to the rpm command:

```bash
rpm -e device-mapper-multipath-0.4.5-27.el4_6.3
```

**Now reboot!**

## Installation

Now we'll install the package. Take it from the CD or anywhere else and install it:

```bash
rpm -ivh EMCpower.Linux*.rpm
```

Now you'll need to launch the license key command:

```bash
emcpreg -install
```

Now update the initial ramdisk:

```bash
mkinitrd -f /boot/initrd-`uname -r`.img `uname -r`
```

Then reboot again!

## Powerpath verification

You may need to start Powerpath service:

```bash
/etc/init.d/PowerPath start
```

And also may reboot the server again.
Now enter powermt command to verify if your LUN can be shown:

```bash
$ powermt display dev=all
Pseudo name=emcpowerb ? All the paths below are handled through this pseudo-device
CLARiiON ID=APM00023500472
Logical device ID=600601F0310A00006F80C3D32D69D711? Unique ID of LUN (LUN Properties)
state=alive; policy=CLAROpt; priority=0; queued-IOs=0? properties/status of the paths


==============================================================================
---------------- Host ---------------   - Stor -   -- I/O Path -  -- Stats ---
### HW Path                 I/O Paths    Interf.   Mode    State  Q-IOs Errors
==============================================================================
  2 QLogic Fibre Channel 2300 sdc        SP A1     active  alive      0      0
  2 QLogic Fibre Channel 2300 sde        SP B0     active  alive      0      0
  3 QLogic Fibre Channel 2300 sdg        SP B1     active  alive      0      0
  3 QLogic Fibre Channel 2300 sdi        SP A0     active  alive      0      0
```

You can also check that you have emc devices:

```bash
$ ls /dev/emcpower*
emcpower emcpowera
```

## Configuring LVM

Now all is ok, but you need to recover your previous partitions from LVM. So edit the LVM config file and add these arguments:

```bash {linenos=table,hl_lines=[3]}
...
filter = [ "r/sd*/", "a/.*/"," "a|/dev/sdb[1-9]|", "a|/dev/mapper/.*$|", "r|.*|" ]
...
```

Now you can reboot or rebuild the LVM cache:

```bash
vgscan -v
lvmdiskscan
```

Now you should see all your disks :-). Look also at the `/dev/mapper` and you'll see them too.
