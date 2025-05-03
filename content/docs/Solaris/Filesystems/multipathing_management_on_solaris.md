---
weight: 999
url: "/Multipathing_management_on_Solaris/"
title: "Multipathing Management on Solaris"
description: "A guide to manage multipathing on Solaris systems, including configuration, verification of fibers, and volume management."
categories: ["Linux", "Solaris"]
date: "2009-11-06T16:40:00+02:00"
lastmod: "2009-11-06T16:40:00+02:00"
tags: ["Solaris", "Storage", "Multipathing", "SAN", "Fibre Channel"]
toc: true
---

## Introduction

Multipathing allows for connection to multiple links. For example, a disk array connected via fiber to machines can have 2 fibers per machine.

To manage this type of configuration, you need to use multipathing.

## Configuration

To check if your multipath is enabled, it's simple. Look at your devices, they should look like this:

```
/dev/dsk/c3t2000002037CD9F72d0s0
```

instead of this:

```
/dev/dsk/c1t1d0s0
```

If this is not the case, then perform the actions that follow.

### Kernel

We need to enable multipathing at the kernel level. To do this, replace the following value (`/kernel/drv/fp.conf`):

```bash
mpxio-disable="yes";
```

to

```bash
mpxio-disable="no";
```

Let's make sure the change is applied at the next restart:

```bash
touch /reconfigure
```

Then restart the server.

## Management

### Verification of fibers

You need to check the status of the fibers before continuing. To check the status of HBAs for example:

```bash
$ fcinfo hba-port
HBA Port WWN: 2100001b3281b4e8
        OS Device Name: /dev/cfg/c2
        Manufacturer: QLogic Corp.
        Model: 375-3356-02
        Firmware Version: 4.04.01
        FCode/BIOS Version:  BIOS: 1.24; fcode: 1.24; EFI: 1.8;
        Serial Number: 0402H00-0850613916
        Driver Name: qlc
        Driver Version: 20080617-2.29
        Type: N-port
        State: online
        Supported Speeds: 1Gb 2Gb 4Gb 
        Current Speed: 4Gb 
        Node WWN: 2000001b3281b4e8
HBA Port WWN: 2101001b32a1b4e8
        OS Device Name: /dev/cfg/c3
        Manufacturer: QLogic Corp.
        Model: 375-3356-02
        Firmware Version: 4.04.01
        FCode/BIOS Version:  BIOS: 1.24; fcode: 1.24; EFI: 1.8;
        Serial Number: 0402H00-0850613916
        Driver Name: qlc
        Driver Version: 20080617-2.29
        Type: N-port
        State: online
        Supported Speeds: 1Gb 2Gb 4Gb 
        Current Speed: 4Gb 
        Node WWN: 2001001b32a1b4e8
```

### Multipathing

* Get the logical units of the system:

```bash
$ mpathadm list lu
        /scsi_vhci/disk@g600a0b8000492c63000005fd49e8305b
                Total Path Count: 4
                Operational Path Count: 4
```

* To get more details:

```bash
$ mpathadm show lu /scsi_vhci/disk@g600a0b8000492c63000005fd49e8305b
Logical Unit:  /scsi_vhci/disk@g600a0b8000492c63000005fd49e8305b
        mpath-support:  libmpscsi_vhci.so
        Vendor:  SUN     
        Product:  LCSM100_F       
        Revision:  0735
        Name Type:  unknown type
        Name:  600a0b8000492c63000005fd49e8305b
        Asymmetric:  yes
        Current Load Balance:  round-robin
        Logical Unit Group ID:  NA
        Auto Failback:  on
        Auto Probing:  NA

        Paths:  
                Initiator Port Name:  2101001b32a12ae9
                Target Port Name:  203400a0b8492c63
                Override Path:  NA
                Path State:  OK
                Disabled:  no

                Initiator Port Name:  2101001b32a12ae9
                Target Port Name:  203500a0b8492c63
                Override Path:  NA
                Path State:  OK
                Disabled:  no

                Initiator Port Name:  2100001b32812ae9
                Target Port Name:  202500a0b8492c63
                Override Path:  NA
                Path State:  OK
                Disabled:  no

                Initiator Port Name:  2100001b32812ae9
                Target Port Name:  202400a0b8492c63
                Override Path:  NA
                Path State:  OK
                Disabled:  no

        Target Port Groups:  
                ID:  1
                Explicit Failover:  yes
                Access State:  active
                Target Ports:
                        Name:  203400a0b8492c63
                        Relative ID:  0

                        Name:  202400a0b8492c63
                        Relative ID:  0

                ID:  2
                Explicit Failover:  yes
                Access State:  standby
                Target Ports:
                        Name:  203500a0b8492c63
                        Relative ID:  0

                        Name:  202500a0b8492c63
                        Relative ID:  0
```

### Application of your Volumes

Quick info for your freshly created LUNs. You may not see them immediately on your systems. Why? Because you might have other previous tasks - check in CAM if you're not queued for the creation of these LUNs in the jobs list. If they're created but you still don't see anything, you'll need to refresh everything rather than use the cache. Run the devfsadm command once so that new volumes are seen:

```bash
devfsadm
```

If it still doesn't work, run this command:

```bash
cfgadm -al
```

### Flushing LUNs that no longer exist

If you have LUNs that have been removed and are not visible to the servers, you may encounter some issues with Sun Cluster if you don't reboot. That's why you should use the following commands after each deletion, on all nodes:

```bash
cfgadm -c configure cN
cfgadm -c configure cN+1
cfgadm -c unconfigure -o unusable_SCSI_LUN cN
cfgadm -c unconfigure -o unusable_SCSI_LUN cN+1
devfsadm -C -v
```

Then use the following command on only one node:

```bash
scgdevs
```

## References

http://docs.sun.com/source/819-0139/ch_3_admin_multi_devices.html  
http://www.princeton.edu/~unix/Solaris/troubleshoot/mpathadm.html
