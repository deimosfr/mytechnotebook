---
weight: 999
url: "/trick-samba-share-size-display/"
title: "Trick Samba Share Size Display"
description: "How to modify the displayed available space on a Samba share to overcome space limitations with nested mount points"
categories: ["Linux", "Samba", "File Sharing"]
date: "2012-02-16T10:17:00+02:00"
lastmod: "2012-02-16T10:17:00+02:00"
tags: ["Samba", "Storage", "Linux", "File Sharing"]
toc: true
---

## Introduction

A colleague of mine found himself in a rather delicate situation. Let me explain the scenario:

- 2 mount points in /mnt, with one nested inside the other
- 1 share on the primary mount point

When the primary mount point is full, you can't copy anything anymore, even if the second nested mount point still has free space. For those who still don't understand:

- /mnt/: 30 MB remaining
- /mnt/disk1: 10 GB remaining
- share: /mnt/

The share tells me that it can't copy more than 30 MB, even into /share/disk1.

## Solution

Here is a solution that allows you to bypass the fact that Windows will analyze the remaining size of the shared folder before copying what you want. In the Samba configuration file, adjust your share like this:

```bash {linenos=table,hl_lines=[8]}
...
[Share]
   comment = Share file space
   path = /mnt/shares/Share
   read only = no
   public = yes
   guest ok = yes
   dfree command = /etc/samba/dfree
   dfree cache time = 3600
   vfs objects = recycle
   create mask = 0775
   directory mask = 0775
   #force user = nobody
   #force group = Team
   recycle:exclude = *.tmp *.temp *.o *.obj ~$*
   recycle:exclude = *.tmp *.temp *.o *.obj ~$*
   recycle:keeptree = True
   recycle:touch = True
   recycle:versions = True
   recycle:noversions = .doc|.xls|.ppt
   recycle:repository = .recycle
   recycle:maxsize = 0
   admin users = @admins
   inherit permissions = Yes
   #case sensitive = no
   #preserve case = yes
...
```

dfree is the argument needed for the Samba daemon to determine the size to display for a given share at startup:

```bash
#!/usr/bin/env bash
df -Pk $1 | tail -1 | awk '{print $2" "$4}'
```

Then apply the proper permissions:

```bash
chmod 700 /etc/samba/dfree
```

## References

[https://www.samba.org/samba/docs/man/manpages-3/smb.conf.5.html](https://www.samba.org/samba/docs/man/manpages-3/smb.conf.5.html)
