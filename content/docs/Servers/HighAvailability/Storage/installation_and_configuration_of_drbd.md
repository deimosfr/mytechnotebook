---
weight: 999
url: "/Installation_et_configuration_de_DRBD/"
title: "Installation and Configuration of DRBD"
description: "How to set up and configure DRBD for disk replication across nodes in a network for high availability"
categories: ["Linux", "Debian"]
date: "2013-09-07T09:56:00+02:00"
lastmod: "2013-09-07T09:56:00+02:00"
tags: ["DRBD", "High Availability", "Clustering", "Storage", "Replication"]
toc: true
---

## Introduction

**DRBD** is a system that allows you to create software RAID1 over a local network.
This enables high availability and resource sharing on a cluster without a disk array.

Here we will install **DRBD8**, with the goal of implementing a cluster filesystem (see documentation on OCFS2) which is not supported on DRBD7.
We'll use the DRBD8 packages from Debian repositories. We'll work on a 2-node cluster.

## Installation

First, install the following packages:

```bash
aptitude install drbd8-utils
```

Then we'll load the module and make it persistent (for future reboots):

```bash
modprobe drbd
echo "drbd" >> /etc/modules
```

## Configuration

### drbd.conf

The drbd.conf file is pretty good by default as it allows you to write an extensible configuration:

```bash
# You can find an example in  /usr/share/doc/drbd.../drbd.conf.example

include "drbd.d/global_common.conf";
include "drbd.d/*.res";
```

I didn't modify it.

### global_common.conf

This file is the default file, which can contain host configurations, but also allows you to have a global configuration for your different DRBD configurations (common section):

```bash
# Global configuration
global {
    # Do not report statistics usage to LinBit
    usage-count no;
}

# All resources inherit the options set in this section
common {
    # C (Synchronous replication protocol)
    protocol C;

    startup {
        # Wait for connection timeout (in seconds)
        wfc-timeout 1 ;
        # Wait for connection timeout, if this node was a degraded cluster (in seconds)
        degr-wfc-timeout 1 ;
    }

    net {
        # Maximum number of requests to be allocated by DRBD
        max-buffers 8192;
        # The highest number of data blocks between two write barriers
        max-epoch-size 8192;
        # The size of the TCP socket send buffer
        sndbuf-size 512k;
        # How often the I/O subsystem's controller is forced to process pending I/O requests
        unplug-watermark 8192;
        # The HMAC algorithm to enable peer authentication at all
        cram-hmac-alg sha1;
        # The shared secret used in peer authentication
        shared-secret "xxx";
        # Split brains
        # Split brain, resource is not in the Primary role on any host
        after-sb-0pri disconnect;
        # Split brain, resource is in the Primary role on one host
        after-sb-1pri disconnect;
        # Split brain, resource is in the Primary role on both host
        after-sb-2pri disconnect;
        # Helps to solve the cases when the outcome of the resync decision is incompatible with the current role assignment
        rr-conflict disconnect;
    }

    handlers {
        # If the node is primary, degraded and if the local copy of the data is inconsistent
        pri-on-incon-degr "echo Current node is primary, degraded and the local copy of the data is inconsistent | wall ";
    }

    disk {
        # The node downgrades the disk status to inconsistent on io errors
        on-io-error pass_on;
        # Disable protecting data if power failure (done by hardware)
        no-disk-barrier;
        # Disable the backing device to support disk flushes
        no-disk-flushes;
        # Do not let write requests drain before write requests of a new reordering domain are issued
        no-disk-drain;
        # Disables the use of disk flushes and barrier BIOs when accessing the meta data device
        no-md-flushes;
    }

    syncer {
        # The maximum bandwidth a resource uses for background re-synchronization
        rate 500M;
        # Control how big the hot area (= active set) can get
        al-extents 3833;
    } 
}
```

I've commented all my changes.

### r0.res

Now we'll create a file to add our resource 0:

```bash
resource r0 {
    # Node 1
    on srv1 {
        device       /dev/drbd0;
        # Disk containing the drbd partition
        disk         /dev/mapper/datas-drbd;
        # IP address of this host
        address      192.168.100.1:7788;
        # Store metadata on the same device
        meta-disk    internal;
    }
    # Node 2
    on srv2 {
        device      /dev/drbd0;
        disk        /dev/mapper/lvm-drbd;
        address     192.168.20.4:7788;
        meta-disk   internal;
    }
}
```

## Synchronization

We need to launch the first sync now.

On both nodes, run this command:

```bash
drbdadm create-md r0
```

Still on both nodes, run this command to activate the resource:

```bash
drbdadm up r0
```

### Node 1

We'll ask the first node to do the first block-by-block replication:

```bash
drbdadm -- --overwrite-data-of-peer primary r0
```

Then we'll have to wait for the sync to finish before continuing:

```bash
> cat /proc/drbd
version: 8.3.7 (api:88/proto:86-91)
srcversion: EE47D8BF18AC166BE219757 
 0: cs:SyncSource ro:Primary/Secondary ds:UpToDate/Inconsistent C r----
   ns:912248 nr:0 dw:0 dr:920640 al:0 bm:55 lo:1 pe:388 ua:2048 ap:0 ep:1 wo:b oos:3283604
        [===>................] sync'ed: 21.9% (3283604/4194304)K
        finish: 1:08:24 speed: 580 (452) K/sec
```

The display of /proc/drbd allows you to see the replication status. At the end, you should have something like this:

```bash
> cat /proc/drbd
version: 8.3.7 (api:88/proto:86-91)
srcversion: EE47D8BF18AC166BE219757
 0: cs:Connected ro:Secondary/Primary ds:UpToDate/UpToDate C r----
    ns:0 nr:4194304 dw:4194304 dr:0 al:0 bm:256 lo:0 pe:0 ua:0 ap:0 ep:1 wo:b oos:0
```

### Node 2

If you want to do dual master, this option must be active in the configuration:

```bash
resource <resource>
  startup {
    become-primary-on both;
  }
  net {
    protocol C;
    allow-two-primaries yes;
  }
}
```

Now we can activate the other node as primary:

```bash
drbdadm primary r0
```

Once the synchronization is complete, DRBD is installed and properly configured.
You now need to format the device `/dev/drbd0` with a filesystem, such as ext3 for active/passive or OCFS2 for example if you want active/active (there are others like GFS2).

```bash
mkfs.ext3 /dev/drbd0
```
or
```bash
mkfs.ocfs2 /dev/drbd0
```

Then mount the volume in a folder to access the data:

```bash
mount /dev/drbd0 /mnt/data
```

Only a primary node can mount and access the data on the DRBD volume.
When DRBD works with HeartBeat in CRM mode, if the primary node goes down, the cluster is able to switch the secondary node to primary.
When the old primary is "UP" again, it will synchronize and become a secondary in turn.

## Usage

### Become master

To set all volumes as primary:

```bash
drbdadm primary all
```

{{< alert context="info" text="Replace <b>all</b> with the name of your volume if you only want to operate on one." />}}

### Become slave

To set a volume as slave:

```bash
drbdadm secondary all
```

### Manual synchronization

To start a manual synchronization (will invalidate all your data):

```bash
drbdadm invalidate all
```

To do the same but on other nodes:

```bash
drbdadm invalidate_remote all
```

## FAQ

### My sync doesn't work, I have: Secondary/Unknown

If you have this type of message:

```bash
> cat /proc/drbd
version: 8.3.7 (api:88/proto:86-91)
srcversion: EE47D8BF18AC166BE219757 
 0: cs:StandAlone ro:Secondary/Unknown ds:Inconsistent/DUnknown   r----
    ns:0 nr:0 dw:0 dr:0 al:0 bm:0 lo:0 pe:0 ua:0 ap:0 ep:1 wo:b oos:4194304
```

You need to check if the machines are properly configured for your resources and also if they can telnet to each other (firewalling etc...)

### What to do in case of split brain?

If you find yourself in this situation:

```bash
> cat /proc/drbd
primary/unknown
```
or
```bash
secondary/unknown
```

1. Unmount the drbd volumes
2. On the primary:

```bash
drbdadm connect all
```

3. On the secondary (this will destroy all data and reimport from the master)

```bash
drbdadm -- --discard-my-data connect all
```

## Resources
- [Documentation on Heartbeat2 Xen cluster with drbd8 and OCFS2](/pdf/heartbeat2_xen_cluster_with_drbd8_and_ocfs2.pdf)
- [DRBD 8 3 Third Node Replication](/pdf/drbd_8_3_third_node_replication_with_debian_etch.pdf)
- [DRBD advanced usages](https://fr.slideshare.net/deimosfr/drbd-25911753)
