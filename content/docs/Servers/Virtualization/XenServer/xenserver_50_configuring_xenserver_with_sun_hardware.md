---
weight: 999
url: "/XenServer_5.0_\\:_Configuration_d'un_XenServer_avec_du_matériel_SUN/"
title: "XenServer 5.0: Configuring XenServer with SUN Hardware"
description: "Guide for configuring XenServer 5.0 with SUN hardware, including multipathing setup for StorageTek storage arrays and fiber channel connectivity."
categories: 
  - Storage
  - Linux
date: "2009-04-19T09:28:00+02:00"
lastmod: "2009-04-19T09:28:00+02:00"
tags: 
  - Servers
  - Storage
  - Virtualization
  - XenServer
toc: true
---

## Introduction

To perform installation on multiple XenServer, I needed to set up multiple hardware devices connected via Fiber channel. I completed two installations, both with this kind of hardware:

- SUN X4150
- SUN X4600
- SUN StorageTek ST2540
- SUN StorageTek ST6140

## Hardware Connectivity

To setup connectivity, I used fiber channel. To make it work, you need to have Fiber Channel Switch and if possible, redundancy (so, 2 switches).

Having this kind of configuration creates multiple multipath links. This is very beneficial for fault tolerance and load balancing.

## Lun Configuration

To have all paths active for multipathing, you need to present your LUN in a specific format to your XenServers. Simply connect to the disk array with CAM (Common Array Manager), select the LUN (or go to Administration to set all the LUNs) and set the host type: **Windows 2000/Server 2003 non-clustered (with Veritas DMP)**.

## Multipathing configuration

By default, no preconfigured device is set for multipath configuration for SUN hardware. That's why we need to change the default configuration to add configuration for the StorageTek.

### SUN StorageTek ST2540

Add these lines to the `/etc/multipath.conf` configuration file:

```bash
        device {
                vendor "SUN"
                product "LCSM100_F"
                getuid_callout "/sbin/scsi_id -g -u -s /block/%n"
                prio_callout "/sbin/mpath_prio_rdac /dev/%n"
                features "0"
                hardware_handler "1 rdac"
                path_grouping_policy group_by_prio
                failback immediate
                rr_weight uniform
                no_path_retry queue
                rr_min_io 1000
                path_checker rdac
        }
```

### SUN StorageTek ST6140 or ST6540

Add these lines to the `/etc/multipath.conf` configuration file:

```bash
        device {
                vendor "SUN"
                product "CSM200_R"
                getuid_callout "/sbin/scsi_id -g -u -s /block/%n"
                prio_callout "/sbin/mpath_prio_rdac /dev/%n"
                features "0"
                hardware_handler "1 rdac"
                path_grouping_policy group_by_prio
                failback immediate
                rr_weight uniform
                no_path_retry queue
                rr_min_io 1000
                path_checker rdac
        }
```

## Refresh configuration

Now you need to restart the multipathd daemon or reboot your servers for changes to take effect.

## FAQ

### I applied a Xen Update Hotfix and my SAN is not recognized anymore

Reapply the configuration to your multipath. Each update will erase it.

## References

http://forums.citrix.com/thread.jspa?threadID=241199&tstart=0  
http://support.citrix.com/article/ctx118791
