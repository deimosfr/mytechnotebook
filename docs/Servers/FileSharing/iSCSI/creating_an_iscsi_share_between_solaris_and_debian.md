
---
weight: 999
url: "/Créer_un_partage_iSCSI_entre_Solaris_et_Debian/"
title: "Creating an iSCSI Share Between Solaris and Debian"
description: "How to set up an iSCSI share between Solaris (server) and Debian (client)"
categories: ["Debian", "Linux"]
date: "2009-11-28T15:28:00+02:00"
lastmod: "2009-11-28T15:28:00+02:00"
tags: ["Servers", "Linux", "Debian", "Solaris", "iSCSI", "Storage"]
toc: true
---

## 1 Introduction

Technology is beautiful, isn't it! Here's a guide to create an iSCSI share between Solaris and Debian. Be careful with performance as this remains iSCSI - I don't recommend it for production environments unless you know what you're doing.

In this setup, Solaris will be our server and Debian our client.

## 2 Installation

### 2.1 Solaris

Let's verify that we have the required packages:

```bash
$ pkginfo | grep iscsi
system      SUNWiscsir                       Sun iSCSI Device Driver (root)
system      SUNWiscsitgtr                    Sun iSCSI Target (Root)
system      SUNWiscsitgtu                    Sun iSCSI Target (Usr)
system      SUNWiscsiu                       Sun iSCSI Management Utilities (usr)
```

### 2.2 Debian

Let's install Open iSCSI:

```bash
apt-get install open-iscsi
```

## 3 Configuration

### 3.1 Solaris

Activate the iSCSI target service:

```bash
svcadm enable iscsitgt
```

Now, let's configure iSCSI discovery:

```bash
iscsiadm modify discovery –sendtargets enable
iscsiadm add discovery-address <ip_solaris>
```

Create ZFS volumes to export:

```bash
zfs create tank/xen
zfs set shareiscsi=on tank/xen
```

Volume1 will inherit properties from tank/xen and will be shared automatically:

```bash
zfs create -s -V 10g tank/xen/volume1
```

Then verify that everything is working:

```bash
iscsiadm list target
Target: iqn.1986-03.com.sun:02:6bc5ce3d-eb83-4055-fe67-d1fd9a7eb7b7
        Alias: tank/xen/volume1
        TPGT: 1
        ISID: 4000002a0000
        Connections: 1
```

### 3.2 Debian

Log into the target:

```bash
iscsiadm -m discovery -t st -p <ip_solaris>
<ip_solaris>:3260,1 iqn.1986-03.com.sun:02:6bc5ce3d-eb83-4055-fe67-d1fd9a7eb7b7
debian-test# iscsiadm  -m node -l -T "iqn.1986-03.com.sun:02:6bc5ce3d-eb83-4055-fe67-d1fd9a7eb7b7?
Logging in to [iface: default, target: iqn.1986-03.com.sun:02:6bc5ce3d-eb83-4055-fe67-d1fd9a7eb7b7, portal: <ip_solaris>,3260]
Login to [iface: default, target: iqn.1986-03.com.sun:02:6bc5ce3d-eb83-4055-fe67-d1fd9a7eb7b7, portal: <ip_solaris>,3260]: successful
```

Verify that everything is working:

```bash
[58182.163989] sd 3:0:0:0: [sdc] 20971520 512-byte hardware sectors (10737 MB)
[58182.167839] sd 3:0:0:0: [sdc] Write Protect is off
[58182.167867] sd 3:0:0:0: [sdc] Mode Sense: 67 00 00 08
[58182.171977] sd 3:0:0:0: [sdc] Write cache: disabled, read cache: enabled, doesn't support DPO or FUA
[58182.175981] sd 3:0:0:0: [sdc] 20971520 512-byte hardware sectors (10737 MB)
[58182.175989] sd 3:0:0:0: [sdc] Write Protect is off
[58182.175989] sd 3:0:0:0: [sdc] Mode Sense: 67 00 00 08
[58182.183974] sd 3:0:0:0: [sdc] Write cache: disabled, read cache: enabled, doesn't support DPO or FUA
[58182.183989]  sdc: sdc1
[58182.187989] sd 3:0:0:0: [sdc] Attached SCSI disk
```

## 4 Resources

http://www.rottenbytes.info/?p=50
