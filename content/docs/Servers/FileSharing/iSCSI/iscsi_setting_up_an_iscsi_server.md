---
weight: 999
url: "/ISCSI_\\:_Mise_en_place_d'un_serveur_iSCSI/"
title: "iSCSI: Setting up an iSCSI Server"
description: "This guide explains how to install and configure an iSCSI server and client on Red Hat systems."
categories: ["Linux", "Storage", "Networking"]
date: "2013-06-04T14:04:00+02:00"
lastmod: "2013-06-04T14:04:00+02:00"
tags: ["iSCSI", "Storage", "Red Hat", "Linux"]
toc: true
---

![ISCSI](/images/iscsi_logo.avif)

## Introduction

[iSCSI](https://fr.wikipedia.org/wiki/ISCSI) (internet SCSI) is an application layer protocol that enables the transport of SCSI commands over a TCP/IP network.

This documentation was created on Red Hat 5 and is compatible with Red Hat 6.

## Server

### Installation

To install an iSCSI server:

```bash
yum install scsi-target-utils
```

### Creating partitions

Create your partitions, then have them detected with a command like:

```bash
partx -a /dev/sda
```

or

```bash
partprobe /dev/sda
```

Preferably use partx.

**Note: For the following steps, I strongly recommend using UUIDs instead of device paths (/dev/xxx). In this documentation, device paths are used for simplicity.**

### Configuration

Let's edit the server configuration file and uncomment the "target" section:

```bash {linenos=table,hl_lines=[1,4,12,13]}
<target iqn.2012-02.fr.deimos.www:iscsi>
        # List of files to export as LUNs
        #backing-store /usr/storage/disk_1.img
        backing-store /dev/sda1
        # Authentication:
        # if no "incominguser" is specified, it is not used
        #incominguser backup secretpass12

        # Access control:
        # defaults to ALL if no "initiator-address" is specified
        #initiator-address 192.168.1.2
        initiator-address 192.168.0.1
</target>
```

- Here's how to name the iSCSI target, which must be unique: **ign.<date>.<reverse_dns>.<strings>[:<substring>]**
  - date: year + month (yyyy-mm)
  - reverse_dns: reversed DNS (fr.deimos.www)
  - strings: name to identify this device (myiscsi)
  - :<substring>: optional, allows adding a name
- backing-store: the disk device or disk image to use
- initiator-address: client addresses authorized to mount this device

Start the service and make it persistent:

```bash
chkconfig tgtd on
service tgtd start
```

Check the configuration like this:

```bash {linenos=table,hl_lines=[1]}
> tgt-admin -s
Target 1: iqn.2012-02.fr.deimos.www:iscsi
    System information:
        Driver: iscsi
        State: ready
    I_T nexus information:
    LUN information:
        LUN: 0
            Type: controller
            SCSI ID: deadbeaf1:0
            SCSI SN: beaf10
            Size: 0 MB
            Online: Yes
            Removable media: No
            Backing store: No backing store
        LUN: 1
            Type: disk
            SCSI ID: deadbeaf1:1
            SCSI SN: beaf11
            Size: 5369 MB
            Online: Yes
            Removable media: No
            Backing store: /dev/sda1
    Account information:
    ACL information:
        192.168.0.1
```

Here's how to get information about iSCSI devices:

- `/sys/class/scsi_host`: all detected iSCSI adapters
- `/sys/block`: lists the peripherals

## Client

### Installation

Clients are called "Initiators" and the target is the recipient (disk array/server).
To install the client:

```bash
yum install iscsi-initiator-utils
```

Then we'll start the service:

```bash
chkconfig iscsi on
service iscsi start
```

### Usage

#### Mounting

First, let's perform a "discovery" to see what devices are available to us:

```bash {linenos=table,hl_lines=[1]}
> iscsiadm --mode discoverydb --type sendtargets --portal <server> --discover
iqn.2012-02.fr.deimos.www:iscsi
```

- `<server>`: enter the **IP address** of the server (not DNS!!!)

Then we'll log in to the device so that it will be mounted on each reboot:

```bash
iscsiadm --mode node --targetname <iqn.2012-02.fr.deimos.www:iscsi> --portal 192.168.1.1:3260 --login
```

- iqn.2012-02.fr.deimos.www:iscsi: the IQN of the server to use
- `<server>`: enter the **IP address** of the server (not DNS!!!)

Once logged in, you can retrieve device information from the logs:

```bash
tail -20 /var/log/messages | grep "/dev"
```

If you need more information, use verbose mode:

```bash
iscsiadm -m node -P 1
```

Next, format the partition in the desired format. Add a line to fstab with the "\_netdev" option, otherwise the machine won't be able to reboot because of the rc.sysinit script. This specifies that the device doesn't use networking:

```bash
/dev/sda1 /mnt/iscsi  ext3   defaults,auto,_netdev 0 0
```

#### Unmounting

To temporarily unmount an iSCSI device (until the next reboot):

```bash
iscsiadm -m node -T <iqn.2012-02.fr.deimos.www:iscsi> -p <server> -u
```

And if you want to delete it permanently:

```bash
iscsiadm -m node -T <iqn.2012-02.fr.deimos.www:iscsi> -p <server> -o delete
```

## Resources

For additional resources on this topic, you might want to consult similar documentation on iSCSI setup for other Linux distributions.
