---
weight: 999
url: "/Corbeille_réseau/"
title: "Network Recycle Bin"
description: "How to set up a network recycle bin with Samba 3 on Red Hat Enterprise Linux"
categories: ["Red Hat", "Linux"]
date: "2006-10-22T20:28:00+02:00"
lastmod: "2006-10-22T20:28:00+02:00"
tags: ["Servers", "Linux", "Network", "Samba"]
toc: true
---

There is a network recycle bin for each share. To set this up, here are the files to configure with the necessary options. My tests were done on a Red Hat Enterprise 4 with Samba 3. This configuration only works with Samba 3. For earlier versions, it's a .recycle file with a different content. But let's proceed with the configuration:

Edit the `smb.conf` file and add these lines:

```bash
vfs objects = recycle
recycle:exclude = *.tmp *.temp *.o *.obj ~$*
recycle:keeptree = True
recycle:touch = True
recycle:versions = True
recycle:noversions = .doc|.xls|.ppt
recycle:repository = .recycle
recycle:maxsize = 0
```

A small script in the crontab to remove items older than 1 week and you're good to go :-)

```bash
#!/bin/sh #

# This is the name of the Dust bin
recyclename=".recycle"

for dustshare in "/home/data/$recyclename" "/home/sales/$recyclename" "/home/share/$recyclename" ; do
     test -d $dustshare || mkdir $dustshare && chown nobody:nobody $dustshare && chmod 700 $dustshare
     find $dustshare -mtime +168 -exec rm -f {} \;
done
```

Of course, don't forget to restart Samba!

```bash
/etc/init.d/smbd restart
```
