---
weight: 999
url: "/fuse-refus-de-monter-les-disques-a-cause-de-dev-fuse/"
title: "FUSE: Unable to Mount Disks Due to /dev/fuse"
description: "How to solve the FUSE disk mounting issue when /dev/fuse is missing."
categories: ["Linux"]
date: "2006-08-05T12:58:00+02:00" 
lastmod: "2006-08-05T12:58:00+02:00"
tags: ["FUSE", "Encryption", "Device Management"]
toc: true
---

I once encountered this problem. I didn't look into it for long, but it's quite annoying, and the worst part is that I don't have an explanation for it.

After a reboot, I got an error message saying that FUSE could not mount my encrypted disks because "/dev/fuse" did not allow it.

Indeed, the device file didn't exist, which is why I had to run this command to recreate it:

```bash
sudo mknod /dev/fuse -m 0666 c 10 229
```
