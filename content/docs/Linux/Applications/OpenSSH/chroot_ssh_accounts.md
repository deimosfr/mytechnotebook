---
weight: 999
url: "/Chrooter_des_comptes_SSH/"
title: "Chroot SSH Accounts"
description: "How to configure SSH accounts with chroot environment for enhanced security" 
categories: ["Linux", "Security", "SSH"]
date: "2008-04-01T10:31:00+02:00"
lastmod: "2008-04-01T10:31:00+02:00"
tags: ["SSH", "security", "chroot", "SFTP"]
toc: true
---

## Introduction

The OpenSSH version (4.8p1 for the GNU/Linux port) features a new configuration option: `ChrootDirectory`.
This has been made possible by a new SFTP subsystem statically linked to sshd.

This makes it easy to replace a basic FTP service without the hassle of configuring encryption and/or bothering with FTP passive and active modes when operating through a NAT router. This is also simpler than packages such as rssh, scponly or other patches because it does not require setting up and maintaining (i.e. security updates) a chroot environment.

## Installation

To enable it, you obviously need the new version 4.8p1. I personally use the cvs version and the debian/ directory of the sid package to build a well integrated Debian package 4.8p1~cvs-1.

## Configuration

In `/etc/ssh/sshd_config`:

You need to configure OpenSSH to use its internal SFTP subsystem.

```bash
Subsystem sftp internal-sftp
```

Then, I configured chroot()ing in a match rule.

```bash
Match group sftponly
        ChrootDirectory /home/%u
        X11Forwarding no
        AllowTcpForwarding no
        ForceCommand internal-sftp
```

The directory in which to chroot() must be owned by root. After the call to chroot(), sshd changes directory to the home directory relative to the new root directory. That is why I use / as home directory.

```bash
chown root.root /home/user
usermod -d / user
adduser user sftponly
```

This seems to work as expected:

```bash
$ sftp user@host
Connecting to host...
user@host's password:
sftp> ls
build               cowbuildinall       incoming            johnbuilderclean
sftp> pwd
Remote working directory: /
sftp> cd ..
sftp> ls
build               cowbuildinall       incoming            johnbuilderclean
```

The only thing I miss is file transfers logging, but I did not investigate this at all. More on this whenever I find some time to do so.
