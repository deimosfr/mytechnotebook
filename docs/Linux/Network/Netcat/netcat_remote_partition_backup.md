---
title: "Netcat: Remote Partition Backup"
slug: netcat-remote-partition-backup/
description: "Guide on how to use Netcat for remote partition backups, including commands for sending compressed images over the network."
categories: ["Linux", "Network"]
date: "2007-08-02T14:34:00+02:00"
lastmod: "2007-08-02T14:34:00+02:00"
tags: ["Netcat", "Backup", "Network", "Linux", "dd"]
---

## Remote Backups

For a simple but efficient way to perform remote backups, you can use the dd and netcat commands.

Netcat is available in two flavors ;-):

```bash
emerge gnu-netcat
```

or

```bash
emerge netcat
```

To create an image of your entire hda1 partition, start netcat in passive (listening) mode on the remote machine:

```bash
netcat -l -p 10000 > image.gz
```

On your machine, run dd to read the partition, gzip to compress it, and netcat to transfer the image to the other machine:

```bash
dd if=/dev/hda1 | gzip | netcat -w 5 remote_ip 10000
```

Refer to [How to clone a Linux box using netcat](https://www.ebruni.it/en/docs/clone_linux/index.htm) for more information.

## Resources
- [Netcat: Creating a Listening Port](../../Network/Netcat/netcat_creating_a_listening_port.md)
- [Netcat: File Transfer](../../Network/Netcat/netcat_file_transfer.md)
- [Netcat Documentation](../../../static/pdf/netcat.pdf)
