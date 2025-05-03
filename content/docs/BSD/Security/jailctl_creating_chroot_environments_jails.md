---
weight: 999
url: "/Jailctl_\\:_Création_de_chroot_(jails)/"
title: "Jailctl: Creating Chroot Environments (Jails)"
description: "A guide on how to create, configure, and manage FreeBSD jails using Jailctl"
categories: ["FreeBSD", "Linux", "Backup"]
date: "2007-04-15T21:50:00+02:00"
lastmod: "2007-04-15T21:50:00+02:00"
tags: ["FreeBSD", "Chroot", "Jails", "System Administration", "Backup", "Security"]
toc: true
---

## Introduction

Jailctl is a shell tool for creating/launching/stopping/updating/backing up/restoring/destroying jails.
By jail, we mean here a "virtual server" and not simply a method for isolating a service.

## Installation

Here's the command to install the package:

```bash
pkg_add -vr jailctl
```

or

```bash
cd /usr/ports/sysutils/jailctl ; make install clean
```

## Configuration

- You need a config file: `/usr/local/etc/jails.conf`
- You also need: a directory where the jails will be stored (`/data` in this example)
- A runme.sh script provided with jailctl that lives by default in `/usr/local/jails/addons/`
- A file dellist4.txt that contains a list of files to delete in the jails because they are not needed (for example commands like mount)
- A file dellist5.txt that contains more files to delete in case jailctl runs on a 5.x or a 6.x (jailctl is indeed compatible with all versions from 4.x to 6.x)
- And finally an etc/ directory with configuration files to install in the new jails (by default, login.conf and make.conf).

All this lives in `/usr/local/jails/addons/` which we will need to move to `/data/addons/` in our example.

### Changes in login.conf

It is recommended to modify the following line:

```bash
      :setenv=PS1=[\\u@\\h] \\w\\\\$ ,MAIL=/var/mail/$,BLOCKSIZE=K,FTP_PASSIVE_MODE=YES, \
       PACKAGEROOT=ftp\c//ftp.no.freebsd.org,CLICOLOR=1,EDITOR=/usr/local/bin/nano:\
```

to use a closer mirror:

```bash
      :setenv=PS1=[\\u@\\h] \\w\\\\$ ,MAIL=/var/mail/$,BLOCKSIZE=K,FTP_PASSIVE_MODE=YES, \
       PACKAGEROOT=ftp\c//ftp.fr.freebsd.org,CLICOLOR=1,EDITOR=/usr/local/bin/nano:\
```

You can also customize the default editor and other settings if needed.

### Changes in jails.conf

Here is the main configuration. Since the file is very well documented internally, here are just the mandatory elements to get started quickly.

Interface on which to add the jail IPs:

```bash
IF="em0"
```

Where the jails will be stored:

```bash
JAIL_HOME="/data/"
```

Where the jail backups will be stored (by default in the same place):

```bash
BACKUPDIR=$JAIL_HOME
```

What not to back up:

```bash
BACKUP_EXCLUDE="--exclude ./usr/ports/* --exclude ./tmp/* --exclude ./var/tmp/* --exclude ./usr/src/*"
```

The jails themselves:

```bash
JAILS=""
JAILS="$JAILS chii.domaine.com:192.168.1.43"
JAILS="$JAILS motosuwa.domaine.com:192.168.1.44"
JAILS="$JAILS sumomo.domaine.com:192.168.1.46"
JAILS="$JAILS yuzuki.domaine.com:192.168.1.47:/usr/local/jails"
JAILS="$JAILS yoshiyuki.domaine.com:192.168.1.48:/data2/yoshi/"
JAILS="$JAILS example.domaine.com:192.168.1.49"
```

Note a special feature recently added to jailctl. You can customize the directory where a specific jail will be stored.

Until now, jails were always stored in $JAIL_HOME/name.of.domain.com. Now, you can either specify another general directory (if you don't put a / at the end, for example here yuzuki will be in `/usr/local/jails/yuzuki.domaine.com/`) or a full directory for a given jail (if you put a / at the end, yoshiyuki will be in `/data2/yoshi/`).

An rc.conf will be placed in the jail, containing:

```bash
RC_CONF='sendmail_enable="NO" sshd_enable="YES" portmap_enable="NO" \
network_interfaces="" tcp_keepalive="NO" inetd_enable="NO"'
```

Finally, you need to provide a DNS for the jail's resolv.conf:

```bash
NAMESERVER="195.95.225.104"
```

It is essential that this DNS is reachable during the "create" of the jail, as packages will be installed by runme.sh at the end of the creation.

Finally, if desired, you can specify scripts that will be executed before/after certain jailctl commands (the scripts will receive as argument $1 the name of the jail and as $2 its jail ID as long as you're at least on a 5.x):

```bash
BEFORESTATUS_HOOKS="/usr/bin/true"
AFTERSTATUS_HOOKS="/usr/bin/true"
BEFORESTART_HOOKS="/usr/bin/true"
AFTERSTART_HOOKS="/usr/bin/true"
BEFORESTOP_HOOKS="/usr/bin/true"
AFTERSTOP_HOOKS="/usr/bin/true"
```

## Advice

WARNING: It is STRONGLY advised against creating jails with a different environment than the host machine, for example a host in -STABLE and jails in -RELEASE, or vice versa. With jailctl, this essentially means that you should not do a cvsup between compiling the host and installing the jails.

Take the time to read jails.conf as well as runme.sh before doing anything to customize them.

## Practical Implementation

### Jail status

```bash
# jailctl status
Jail status (*=running, !=not configured):
*chii.domaine.com (192.168.1.43)
*motosuwa.domaine.com (192.168.1.44)
*sumomo.domaine.com (192.168.1.46)
*yuzuki.domaine.com (192.168.1.47)
 yoshiyuki.domaine.com (192.168.1.48)
*example.domaine.com (192.168.1.49)
```

In this example, all jails are installed and yoshiyuki is not running. A jail not yet created would be marked with an exclamation mark.

### Creating a jail

```bash
# jailctl create example.domaine.com
Creating jail example.domaine.com...
>>> Making hierarchy
>>> Installing everything
Setting root password in jail
Changing local password for root
New Password:
Retype New Password:
chsh: user information updated
use.perl: not found
ifconfig em0 inet 192.168.1.49 netmask 0xffffffff alias
jail /data/example.domaine.com example.domaine.com 192.168.1.49 /bin/sh /runme.sh
ifconfig em0 inet 192.168.1.49 netmask 0xffffffff -alias
```

The use.perl is still there for compatibility reasons. This is not an error. The only information needed for the installation is the root password of the jail, if jailctl is not run in batch mode.

### Starting a jail

```bash
# jailctl start example.domaine.com
Starting jail example.domaine.com...
stty: stdin isn't a terminal
yellow-sub# ln: /dev/log: Operation not permitted
```

The errors are normal and simply due to jail peculiarities.

### Stopping a jail

```bash
# jailctl stop example.domaine.com
Stopping jail example.domaine.com...
Sending TERM signal to jail processes...
Stopping cron.
Shutting down daemon processes:.
Shutting down local daemons:.
Terminated
.
```

### Backing up a jail (not running)

```bash
# jailctl backup example.domaine.com
Doing cold backup of jail example.domaine.com...
```

A nice tar.gz appears in the directory where the jails are located.

Restoring a jail (necessarily not running)

```bash
# jailctl restore example.domaine.com
No valid jail specified!
 
Usage:
jailctl <command> <jail> [<path>]
<command> = start|stop|status|create|delete|upgrade|backup|restore
<jail> = hostname|all
<path> = Backup destination / restore source
```

But what's happening?

To restore a jail, it must not exist.

Deleting a jail

```bash
# jailctl delete example.domaine.com
Deleting jail example.domaine.com...
```

Restoring a jail, second attempt

```bash
# jailctl restore example.domaine.com
Restoring jail example.domaine.com from backup
```

Backing up a jail (running)

This is the "premium option".

```bash
# jailctl backup example.domaine.com
Doing warm backup of jail example.domaine.com...
tar: ././var/run/log: tar format cannot archive socket: Inappropriate file type or format
tar: ././var/run/logpriv: tar format cannot archive socket: Inappropriate file type or format
```

The errors are normal.
