---
weight: 999
url: "/Installation_et_configuration_de_CVS/"
title: "CVS Installation and Configuration"
description: "How to install and configure CVS on Debian and Red Hat Linux distributions."
categories: ["Linux", "Backup", "Debian"]
date: "2008-12-26T18:57:00+02:00"
lastmod: "2008-12-26T18:57:00+02:00"
tags: ["CVS", "Version Control", "Linux", "Debian", "Red Hat"]
toc: true
---

## Installation

### Debian

```bash
apt-get install cvs xinetd
```

### Red Hat

```bash
up2date cvs
```

## Creating Admin Group and Directories

```bash
adduser cvsadmin
mkdir -p /home/cvsadmin/repository
mkdir /home/cvsadmin/.lock
chown -Rf :cvsadmin /home/cvsadmin
chmod -Rf 775 /home/cvsadmin
chmod -Rf 777 /home/cvsadmin/.lock
ln -s /home/cvsadmin/repository /usr/local/cvsroot
```

## Repository Initialization

```bash
export CVSROOT=/usr/local/cvsroot
cvs -d $CVSROOT init
chown -Rf root:cvsadmin /home/cvsadmin/repository/* && chmod -Rf 777 /home/cvsadmin/repository/*
```

## Activating Lock Files

To ensure that lock files are stored in a directory writable by all, you need to modify the `/usr/local/cvsroot-backup/CVSROOT/config` file as follows:

```bash
# Put CVS lock files in this directory rather than directly in the repository.
LockDir=/home/cvsadmin/.lock
```

## Xinetd Configuration

Add a file **"/etc/xinetd.d/cvspserver"** for the CVS PServer in the directory:

```text
service cvspserver
{
        port = 2401
        socket_type = stream
        protocol = tcp
        user = root
        wait = no
        disable = no
        type = UNLISTED

        # Debian
        server = /usr/bin/cvs
        # RedHat
        # server = /usr/local/bin/cvs

        server_args = -f --allow-root /usr/local/cvsroot pserver
}
```

## User Management

Each user must have their own home directory and be the owner of it.
By default, all users will have the ability to commit and write in each module (directories at the root of the repository) but not to create new modules.
To allow a user to write inside the repository (create/delete modules, etc.), simply add them to the "cvsadmin" group by modifying the `/etc/group` file as follows:

```bash
cvsadmin:x:509:username
```
