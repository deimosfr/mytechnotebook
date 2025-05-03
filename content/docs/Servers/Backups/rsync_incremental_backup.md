---
weight: 999
url: "/Rsync_\\:_Sauvegarde_incrémentale/"
title: "Rsync: Incremental Backup"
description: "Guide on how to use Rsync for incremental backups, including manual examples and script resources"
categories: ["Linux", "Backup"]
date: "2013-09-06T08:43:00+02:00"
lastmod: "2013-09-06T08:43:00+02:00"
tags: ["rsync", "backup", "incremental", "ssh"]
toc: true
---

## Scripts and Programs

There are scripts and programs available to facilitate the use of this solution:

- [Easy Automated Snapshot-Style Backups with Linux and Rsync](https://www.mikerubel.org/computers/rsync_snapshots/)
- [rsnapshot](https://www.rsnapshot.org/)

## A Manual Example

You want to copy your disk/backup a partition.  
Assuming you have created a partition on another disk and mounted it (in this example `/mnt/usbharddrivemain`):

- Duplicate a partition:

```bash
rsync --progress --stats -avxzl --exclude "/mnt/usbharddrivemain/" --exclude "/mnt/usbharddriveboot/" --exclude "/usr/portage/" --exclude "/proc/" --exclude "/root/.ccache/" --exclude "/var/log/" --exclude "/sys" --exclude "/dev" --exclude "tmp/" /* /mnt/usbharddrivemain
```

- Duplicate a partition and delete files that are no longer current (that have been deleted from the source partition):

```bash
rsync --progress --stats --delete -avxzl --exclude "/mnt/usbharddrivemain/" --exclude "/mnt/usbharddriveboot/" --exclude "/usr/portage/" --exclude "/proc/" --exclude "/root/.ccache/" --exclude "/var/log/" --exclude "/sys" --exclude "/dev" --exclude "tmp/" /* /mnt/usbharddrivemain
```

- To backup `/boot` (another partition):

```bash
rsync --progress --stats -avxzl /boot /mnt/usbharddriveboot
rsync --progress -avxzl --stats --delete /boot /mnt/usbharddriveboot
```

To restore, you can boot from the second hard drive or use a Live CD. Repeat the previous commands by changing the source and destination, for example: `/mnt/usbharddrivemain /mnt/driveToRestoreTo`.

- You may encounter problems with certain files, such as LDAP databases. To work around this issue, use the 'Sparse' option:

```bash
rsync -avxzl --sparse /var/lib/ldap/mdb/db/data.mdb /mnt/usbharddriveboot
```

- To synchronize via SSH:

```bash
rsync -e ssh -av --delete "/rsync" remote:backupdir
```

## Resources
- [Incremental Backup Server](/pdf/serveur_de_sauvegardes_incrémentales.pdf)
- [Backing up your data with Rsync](/pdf/rsync.pdf)
- [Information about Rsync memory consumption](https://www.samba.org/rsync/FAQ.html#4)
