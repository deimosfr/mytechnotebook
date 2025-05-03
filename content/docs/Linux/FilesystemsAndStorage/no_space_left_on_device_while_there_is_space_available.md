---
weight: 999
url: "/No_space_left_on_device_alors_qu\\'il_y_a_de_la_place/"
title: "No space left on device while there is space available"
description: "How to troubleshoot 'No space left on device' errors when disk space seems available by checking inode usage."
categories: ["Linux"]
date: "2008-10-13T20:43:00+02:00"
lastmod: "2008-10-13T20:43:00+02:00"
tags: ["Servers", "Linux", "Troubleshooting", "File System"]
toc: true
---

## The problem

Sometimes the system can no longer write to a partition and reports "No Space Left On Device" even though "df -h" shows that there is enough space available.

For example, I can no longer write at all to my "/var" partition. Let's check the space with df -h:

```bash
Sys. de fich.         Tail. Occ. Disp. %Occ. Monté sur
/dev/sda3             9,7G  2,4G  6,9G  26% /
tmpfs                 498M     0  498M   0% /lib/init/rw
tmpfs                 498M     0  498M   0% /dev/shm
/dev/sda1             122M  9,9M  106M   9% /boot
/dev/sda6             132G  3,6G  121G   3% /mnt/datas
/dev/sda5             4,9G  4,0G  629M  87% /var
```

## The solution

BUT WHY??? After some reflection and research... OF COURSE! The inodes!!!

Indeed, this can be due to a lack of inodes, which can be caused by too many small files in a directory. To check if the inodes are fine, we'll use **"df -i"**:

```bash
Sys. de fich.         Inodes   IUtil.  ILib. %IUti. Monté sur
/dev/sda3            1281696   85671 1196025    7% /
tmpfs                 127353       6  127347    1% /lib/init/rw
tmpfs                 127353       1  127352    1% /dev/shm
/dev/sda1              32256      29   32227    1% /boot
/dev/sda6            17481728    7595 17474133    1% /mnt/datas
/dev/sda5             640000  640000       0  100% /var
```

Well, look at that! 100% inodes on **/var**

Now we just need to find the guilty directory and deal with it. (In my case, it was the **/var/amavis/virusmail** directory which was full of small compressed files.)
