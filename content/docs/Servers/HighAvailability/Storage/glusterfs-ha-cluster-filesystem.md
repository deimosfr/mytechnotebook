---
weight: 999
url: "/glusterfs-ha-cluster-filesystem/"
title: "GlusterFS: High Availability Cluster Filesystem"
description: "How to setup GlusterFS to create a high availability clustered filesystem with automatic file replication across multiple servers."
categories: ["Linux", "File System", "Cluster", "High Availability"]
date: "2011-04-11T08:18:00+02:00"
lastmod: "2011-04-11T08:18:00+02:00"
tags: ["glusterfs", "cluster", "ha", "filesystem", "linux", "openvz"]
toc: true
---

## Introduction

[GlusterFS](https://fr.wikipedia.org/wiki/GlusterFS) is an open source distributed parallel file system capable of scaling to several petabytes.
GlusterFS is a cluster/network file system. GlusterFS comes with two components, a server and a client.
The storage server (or each server in a cluster) runs glusterfsd and clients use the mount command or glusterfs client to mount the file systems served, using FUSE.

The goal here is to run 2 servers that will perform complete replication of part of a filesystem.

Be careful not to run this type of architecture on the Internet as performance will be catastrophic. Indeed, when a node wants to read access a file, it must contact all other nodes to see if there are any discrepancies. Only then does it authorize reading, which can take a long time depending on the architectures.

## Installation

To install on Debian...easy move:

```bash
aptitude install glusterfs-server glusterfs-examples
```

## Configuration

### hosts

As with any respectable cluster, we must correctly configure the hosts table to avoid troubles in case of DNS loss. Add the following hosts:

(`/etc/hosts`)

```bash
192.168.110.2 rafiki.deimos.fr  rafiki
192.168.20.6 ed.deimos.fr  ed
```

### Generating Configurations

We'll simplify things here by generating a configuration for a RAID 1:

```bash
cd /etc/glusterfs
rm -f *
/usr/bin/glusterfs-volgen --name www --raid 1 rafiki:/var/www-orig ed:/var/www-orig
```

Then on each server, rename the file corresponding to the server you're on to glusterfsd.vol and the tcp file to glusterfs.vol:

```bash
mv rafiki-www-export.vol glusterfsd.vol
mv www-tcp.vol glusterfs.vol
```

Don't forget to do the same on the other server and you can restart your glusterfs server.

### Server

On the server side, we'll apply this configuration:

(`/etc/glusterfs/glusterfsd.vol`)

```bash
### file: server-volume.vol.sample
 
#####################################
###  GlusterFS Server Volume File  ##
#####################################
 
#### CONFIG FILE RULES:
### "#" is comment character.
### - Config file is case sensitive
### - Options within a volume block can be in any order.
### - Spaces or tabs are used as delimitter within a line. 
### - Multiple values to options will be : delimitted.
### - Each option should end within a line.
### - Missing or commented fields will assume default values.
### - Blank/commented lines are allowed.
### - Sub-volumes should already be defined above before referring.
 
volume posix1
  type storage/posix
  option directory /var/www
end-volume
 
volume locks1
    type features/locks
    subvolumes posix1
end-volume
 
volume brick1
    type performance/io-threads
    option thread-count 8
    subvolumes locks1
end-volume
 
volume server-tcp
    type protocol/server
    option transport-type tcp
    option auth.addr.brick1.allow *
    option transport.socket.listen-port 6996
    option transport.socket.nodelay on
    subvolumes brick1
end-volume
```

### Client

For the client part, we tell it that we want to do "raid1". Here is the configuration to apply on the "ed" node:

(`/etc/glusterfs/glusterfs.vol`)

```bash
### file: client-volume.vol.sample
 
#####################################
###  GlusterFS Client Volume File  ##
#####################################
 
#### CONFIG FILE RULES:
### "#" is comment character.
### - Config file is case sensitive
### - Options within a volume block can be in any order.
### - Spaces or tabs are used as delimitter within a line. 
### - Each option should end within a line.
### - Missing or commented fields will assume default values.
### - Blank/commented lines are allowed.
### - Sub-volumes should already be defined above before referring.
 
# RAID 1
# TRANSPORT-TYPE tcp
volume ed-1
    type protocol/client
    option transport-type tcp
    option remote-host ed
    option transport.socket.nodelay on
    option transport.remote-port 6996
    option remote-subvolume brick1
end-volume
 
volume rafiki-1
    type protocol/client
    option transport-type tcp
    option remote-host rafiki
    option transport.socket.nodelay on
    option transport.remote-port 6996
    option remote-subvolume brick1
end-volume
 
volume mirror-0
    type cluster/replicate
    subvolumes rafiki-1 ed-1
end-volume
 
volume readahead
    type performance/read-ahead
    option page-count 4
    subvolumes mirror-0
end-volume
 
volume iocache
    type performance/io-cache
    option cache-size `echo $(( $(grep 'MemTotal' /proc/meminfo | sed 's/[^0-9]//g') / 5120 ))`MB
    option cache-timeout 1
    subvolumes readahead
end-volume
 
volume quickread
    type performance/quick-read
    option cache-timeout 1
    option max-file-size 64kB
    subvolumes iocache
end-volume
 
volume writebehind
    type performance/write-behind
    option cache-size 4MB
    subvolumes quickread
end-volume
 
volume statprefetch
    type performance/stat-prefetch
    subvolumes writebehind
end-volume
```

## Execution

### Server

Restart glusterfs after adapting to your needs.

### Client

Simply mount the glusterfs partition:

```bash
glusterfs /var/www
```

You now have access to your glusterfs mount point in /var/www.

## FAQ

### Force Client Synchronization

If you want to force data synchronization for a client, it's simple. Just go to the directory where the glusterfs share is located (here /mnt/glusterfs), then perform a directory traversal like this:

```bash
ls -lRa
```

This will read everything and therefore copy everything.

### www-posix: Extended attribute not supported

If you look in your logs and see something like this:

(`/var/log/glusterfs/glusterfsd.vol.log`)

```bash
...
+------------------------------------------------------------------------------+
[2010-10-17 00:40:30] W [afr.c:2743:init] www-replicate: Volume is dangling.
[2010-10-17 00:40:30] C [posix.c:4936:init] www-posix: Extended attribute not supported, exiting.
[2010-10-17 00:40:30] E [xlator.c:844:xlator_init_rec] www-posix: Initialization of volume 'www-posix' failed, review your volfile again
[2010-10-17 00:40:30] E [glusterfsd.c:591:_xlator_graph_init] glusterfs: initializing translator failed
[2010-10-17 00:40:30] E [glusterfsd.c:1395:main] glusterfs: translator initialization failed. exiting
```

It means you have permission problems. In my case, this happened in an OpenVZ container. To solve the problem, here's the solution to apply on the host machine (not in the VE) (**warning, this requires stopping, applying configurations, then restarting the VE**):

If you want to do glusterfs in a VE, you may encounter permission problems:

```
fuse: failed to open /dev/fuse: Permission denied
```

To work around them, we'll create the fuse device from the host on the VE in question and add admin rights to it (not great in terms of security, but no choice):

```bash
vzctl set $my_veid --devices c:10:229:rw --save
vzctl exec $my_veid mknod /dev/fuse c 10 229
vzctl set $my_veid --capability sys_admin:on --save
```

Note: Don't forget to load the fuse module on your host machine:

(`/etc/modules`)

```bash
...
fuse
```

## Resources
- [High-Availability Storage Cluster With GlusterFS](/pdf/high-availability_storage_cluster_with_glusterfs_on_ubuntu.pdf)
- [High-availability storage with GlusterFS on Debian Lenny](https://www.howtoforge.org/high-availability-storage-with-glusterfs-on-debian-lenny-automatic-file-replication-across-two-storage-servers)
- [OpenVZ Forum Discussion on GlusterFS](https://forum.openvz.org/index.php?t=msg&goto=35230&)
- [Additional GlusterFS Configuration Example](https://olemskoi.ru/node/3788)
