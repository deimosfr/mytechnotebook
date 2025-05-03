---
weight: 999
url: "/GFS2_\\:_Le_Filesystem_Cluster_de_Red_Hat/"
title: "GFS2: Red Hat Cluster Filesystem"
description: "Learn how to install and configure GFS2, the Global File System cluster filesystem developed by Red Hat for Linux."
categories: ["Linux", "File Systems", "Clusters"]
date: "2012-03-06T12:54:00+02:00"
lastmod: "2012-03-06T12:54:00+02:00"
tags: ["GFS2", "Red Hat", "Cluster", "Filesystem", "LVM"]
toc: true
---

## Introduction

[Global File System (GFS)](https://fr.wikipedia.org/wiki/Global_File_System) is a shared file system designed for Linux or IRIX clusters. GFS and GFS2 are different from distributed file systems like AFS, Coda, or InterMezzo because they allow all nodes to have direct concurrent access to the same block storage device. Additionally, GFS and GFS2 can also be used as a local file system.

GFS does not have a disconnected mode; there are no clients or servers. All nodes in a GFS cluster are equal. Using GFS in a computer cluster requires hardware that allows access to shared data and a lock manager to control access to the data. The lock manager is a separate module, therefore GFS and GFS2 can use a Distributed Lock Manager (DLM) for cluster configurations and the "nolock" lock manager for local file systems. Older versions of GFS also support GULM, a server-based lock manager that implements redundancy through failover.

GFS2 tends to experience performance degradation when there are too many files in a directory. However, there are no issues with large files or many directories with few files in them.

GFS and GFS2 are free software, distributed under the GPL license.

This documentation was created on Red Hat Enterprise Linux 6.2.

## Prerequisites

- You need to [know how to use LVM]({{< ref "docs/Linux/FilesystemsAndStorage/lvm_working_with_logical_volume_management.md" >}}).
- You must use [LVM Cluster]({{< ref "docs/Servers/HighAvailability/RedHatClusterSuite/installation_and_configuration_of_red_hat_cluster_suite.md#clvmd">}}) to be able to use GFS2 (and therefore the cluster part as well).
- If you want to use quotas, it's better to have a good understanding of how standard filesystem quotas work by following [this documentation]({{< ref "docs/Linux/FilesystemsAndStorage/setting_up_quotas_on_linux.md">}}).

## Installation

To install GFS2:

```bash
yum -y install gfs2-utils
```

## Configuration

To format a GFS2 partition:

```bash
> mkfs.gfs2 -p lock_dlm -t cluter_name:gfs_name -j 3 /dev/my_vg/my_gfs
This will destroy any data on /dev/my_vg/my_gfs.

Are you sure you want to proceed? [y/n] y

Device:                    /dev/my_vg/my_gfs
Blocksize:                 4096
Device Size                0.50 GB (131072 blocks)
Filesystem Size:           0.50 GB (131070 blocks)
Journals:                  1
Resource Groups:           2
Locking Protocol:          "lock_dlm"
Lock Table:                "cluter:gfs_web"
UUID:                      27DBECB7-E200-E8D9-D484-179AB59B5595
```

- `-p lock_dlm`: allows locking the partition during formatting
- `-t cluter_name:gfs_name`: the cluster name and a name for the GFS
- `-j 3`: indicates the number of journals to create (1 per node, if there aren't enough, it won't be mountable on more nodes)

Then to mount a GFS2 partition:

```bash
mount -t gfs2 /dev/my_vg/my_gfs /mnt
```

In case of emergency, you need to access it to recover data:

```bash
mount -t gfs2 -o lockproto=lock_nolock /dev/my_vg/my_gfs /mnt
```

If the partition needs to be mounted on another cluster:

```bash
mount -t gfs2 -o locktable=cluster_name:gfs_name /dev/my_vg/my_gfs /mnt
```

When you add a node that needs to access a GFS2 partition, you need to add a journal, otherwise the mount will be refused:

```bash
gfs2_jadd -j 1 /mnt
```

I used the mount point here, but you can also use the device in `/dev`.

### Quotas

You can use quotas, but they work differently from standard quotas:

```bash
mount -t gfs2 -o quota=on /dev/my_vg/my_gfs /mnt
```

For hard quota limitation:

```bash
gfs2_quota limit -u user -l size /mountpoint
```

- size: size in Megabytes

For soft quota limitation:

```bash
gfs2_quota warn -u user -l size /mountpoint
```

To list a user's quotas:

```bash
gfs2_quota get -u user
```

And to list quotas for a mount point:

```bash
gfs2_quota list -f /mountpoint
```

### Expanding a Partition

To expand a partition, you'll need to create a new PV, extend the VG, and extend the LV:

```bash
pvcreate /dev/sdb
vgextend my_vg /dev/sdb
lvextend -L +10G /dev/my_vg/my_gfs
```

```bash
gfs2_grow -v /mountpoint
```

### Repairing a Partition

```bash
gfs2_fsck /dev/my_vg/my_gfs
```

### Modifying the Superblock

Sometimes it's necessary to change the superblock, for example, if you have a "lock_nolock" option that was created during the mkfs of your partition and you now want to change it.

You must unmount your partition before proceeding further. To change the lock manager:

```bash
gfs2_tool sb <dev> proto [lock_dlm,lock_nolock]
```

To lock the table name:

```bash
gfs2_tool sb <dev> table cluster:my_gfs
```

To list information about the superblock:

```bash
gfs2_tool sb <dev> all
```
