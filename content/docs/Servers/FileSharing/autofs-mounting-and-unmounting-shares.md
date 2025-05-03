---
url: "/AutoFS_\\:_montage_et_démontage_de_partages/"
title: "AutoFS: Mounting and Unmounting Shares"
description: "Learn how to set up AutoFS to automatically mount and unmount various network filesystems including NFS, CIFS, SSH, and FTP."
categories: ["Linux", "Storage", "Network"]
date: "2013-09-26T07:53:00+02:00"
lastmod: "2013-09-26T07:53:00+02:00"
tags: ["AutoFS", "NFS", "CIFS", "Samba", "SSH", "FTP", "Linux"]
toc: true
---

## Introduction

You already installed Linux on your networked desktop PC and now you want to work with files stored on some other PCs in your network. This is where autofs comes into play. This tutorial shows how to configure autofs to use CIFS to access Windows or Samba shares from Linux Desktop PCs. It also includes a tailored configuration file.

If autofs Version 4.0 or newer is already installed, you should find these files:

```
/etc/auto.master
/etc/auto.smb
```

on your system. Otherwise start the package manager of your distribution (e.g. YaST on SuSE, synaptic on Debian or Ubuntu, ...) and install it. When you are at it, also install the Samba client package (look for smbclient), because we will also need this.

## Installation

Just do:

```bash
apt-get install autofs
```

## Configuration

### NFS

#### Installation

For NFS, it's not very hard, just be sure it works fine when you mount it manually. You also need to install portmap:

```bash
apt-get install portmap
```

#### Configuration

Now edit for example the auto.master file, and add a nfs file line (`/etc/auto.master`):

```
/mnt    /etc/auto.nfs --timeout=600
```

This means it will be mounted in `/mnt` and will be automatically unmounted after 10 min. Now edit your new file (`/etc/auto.nfs`):

```
* -fstype=nfs,rw 192.168.0.187:/home
```

With some virtualized systems like OpenVZ, you'll need to add 'nolock' option in addition to fstype options (`/etc/auto.nfs`):

```
* -fstype=nfs,rw,nolock 192.168.0.187:/home
```

### CIFS

#### Installation

If autofs is already installed, it is probably still not configured and not working. Assuming your Linux Distribution contains a Linux 2.6.x kernel I recommend to use the common internet file system (cifs) module to access files on the network. You also need smbfs to be installed:

```bash
apt-get install smbfs cifs-utils
```

#### Configuration

Please store the following file as (`/etc/auto.master`):

```
/mnt     /etc/auto.cifs
```

on your computer. You need root (or sudo) to have the permissions to do this (`/etc/auto.cifs`):

```
#
share  -fstype=cifs,rw,noperm,credentials=/etc/auto.cred ://server/share
```

This file must be executable to work (chmod 755)!

Just create a /etc/auto.cred file where credentials will be added. Then add those 2 lines (`/etc/auto.cred`):

```
username=le_login
password=le_password
```

Change the permissions:

```bash
chmod 755 /etc/auto.cifs
chmod a-x /etc/auto.cifs
chmod 600 /etc/auto.cred
```

Restart service:

```bash
/etc/init.d/autofs restart
```

That's all.

### SSH

#### Installation

First we need to install sshfs:

```bash
apt-get install sshfs fuse
```

#### Configuration

##### Key exchange

You first need to do [key exchange]({{< ref "docs/Linux/Network/OpenSSH/openssh_ssh_key_exchange.md" >}}) from the root user on remote host.

##### Autofs configuration

Now edit the master file and add this line (`/etc/auto.master`):c/

```
/mnt/remote_server   /etc/auto.sshfs  --timeout=360,--ghost
```

Now add this and adapt it for your needs (`/etc/auto.sshfs`):

```bash
myfolder -fstype=fuse,port=22,rw,allow_other :sshfs\#login@remote_host\:/the/share
```

- myfolder: folder will be accessible in /mnt/remote_server/myfolder
- login: type your ssh login
- /the/share: set the remote share

Do not forget to create the folder on your client machine.

Now reload autofs:

```bash
/etc/init.d/autofs reload
```

Now access to /mnt/remote_server/myfolder it will be automatically mounted.

### FTP

You'll need to install this package:

```bash
aptitude install curlftpfs
```

Then add this line in /etc/auto.master:

```
# /etc/auto.master
/mnt/autofs    /etc/auto.ftp --timeout=600
```

Now we need to create this file and add a line like that (replace with your information):

```
# /etc/auto.ftp
backups_ftp      -fstype=fuse,allow_other,user=user:password   :curlftpfs\server
```

### Create a mountpoint per home user

To create a home mountpoint per connected user, you have to configure (NFS for example) like this:

```
# /etc/auto.nfs
* -fstype=nfs,rw 192.168.0.187:/home/&
```

&: corresponds to each user

## Verify

Use the command:

```bash
ls -als /cifs/FILESERVERNAME/SHARENAME
```

or

```bash
mount.cifs //server/share /test -o username=login,password=pass
```

to check if it works. If not, consult the system logfiles (usually /var/log/messages or /var/log/syslog) for messages.

## FAQ

### My NFS connection is very slow

I can see this kind of things in the syslogs:

```
# /var/log/syslog
...
Jun  6 16:26:40 debusertest kernel: portmap: server localhost not responding, timed out
Jun  6 16:26:40 debusertest kernel: RPC: failed to contact portmap (errno -5).
```

You just need to install the portmap package :-)

### I encounter NFS delay problem

If you have this kind of problem:

```
Jul  8 15:50:39 deb-devtest2 automount[8415]: mount(nfs): mkdir_path /mnt/share/drive failed: No such file or directory
Jul  8 15:50:39 deb-devtest2 automount[8415]: failed to mount /mnt/share/drive
```

Grow your auto.nfs file, option timeo:

```
# /etc/auto.nfs
drive -fstype=nfs,rw,rsize=8192,wsize=8192,timeo=60,intr server:/700G/share/drive
```

### My CIFS doesn't want to mount

You may have a problem while automounting a cifs share getting these errors in the logs:

```
Dec 27 11:36:07 pmavro-laptop automount[27824]: lookup(program): lookup for backups failed
Dec 27 11:36:07 pmavro-laptop automount[27824]: failed to mount /mnt/backups/backups
```

As recommended in this documentation, auto.cifs needs execute rights. But depending on the version of autofs you use, you may need to remove those rights:

```bash
chmod a-x /etc/auto.cifs
```

Then restart autofs daemon.

### Why autofs volumes are hidden when not mounted?

It is possible to show wished volumes that are not mounted, this is easier to know which ones are available. To do it, simply add '--ghost' on a line on auto.master. Example:

```
# /etc/auto.master
[...]
/mnt/mount1      /etc/auto.sshfs --timeout=720,--ghost
/mnt/mount2      /etc/auto.cifs --timeout=720,--ghost
```
