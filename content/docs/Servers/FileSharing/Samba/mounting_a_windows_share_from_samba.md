---
weight: 999
url: "/Monter_un_partage_Windows_depuis_Samba/"
title: "Mounting a Windows Share from Samba"
description: "How to mount Windows shares from Linux using Samba and CIFS methods, with configuration options for different Windows versions."
categories: ["Linux", "Servers", "Network"]
date: "2009-04-19T16:12:00+02:00"
lastmod: "2009-04-19T16:12:00+02:00"
tags: ["Samba", "CIFS", "Windows", "Network", "Mounting", "File Sharing"]
toc: true
---

## Introduction

Mounting Windows shares from a Linux machine with Samba is straightforward with older Windows versions. However, for versions integrating NT technology (starting with Windows 2000), you might encounter problems like:

```
params.c:Parameter() - Ignoring badly formed line in configuration file: +########## Domains ###########
cli_negprot: SMB signing is mandatory and we have disabled it.
30432: protocol negotiation failed
SMB connection failed
```

Or sometimes:

```
smb_add_request: request [f65480c0, mid=24160] timed out!
smb_add_request: request [f6548ec0, mid=24161] timed out!
smb_add_request: request [f6548ec0, mid=24162] timed out!
smb_add_request: request [f6548ec0, mid=24163] timed out!
smb_add_request: request [f6548ec0, mid=24164] timed out!
```

This can be problematic when you need to perform this operation quickly. That's why I'm sharing the solution.

## Installation

Start by installing the minimum requirements (ensure your kernel is >= 2.2.x):

```bash
apt-get install smbfs
```

### Older than Windows 2000

For these versions, check that you have what you need at the kernel level (insert your kernel version):

```bash
grep CONFIG_SMB_FS < /boot/config-2.*
CONFIG_SMB_FS=m
```

### Windows 2000 or newer

For versions newer than NT4, we'll use [CIFS](https://en.wikipedia.org/wiki/CIFS):

```bash
grep CONFIG_CIFS < /boot/config-2.*
CONFIG_CIFS=m
```

You can decide whether to keep SMB_FS or CIFS as modules or recompile them directly into the kernel.

## Configuration

Now we're ready to connect to Windows shares.

### Older than Windows 2000

If using SMB, here's the command:

```bash
mount -t smbfs -o username=USER,password=PASS //SERVER/SHARE /DESTINATION/
```

Or in `/etc/fstab`:

```
//nasServer/SHARE /home/NAS/LM smbfs user,noauto,rw,username=user,password=pass 0 0
```

If you encounter this type of message:

```
params.c:Parameter() - Ignoring badly formed line in configuration file: +########## Domains ###########
cli_negprot: SMB signing is mandatory and we have disabled it.
28570: protocol negotiation failed
SMB connection failed
```

You have two options:
- Switch to CIFS (see below)
- Modify Windows security options (not recommended)

#### Modifying Windows Security Options

As mentioned above, this is strongly discouraged. But if you have no other choice, here's the solution:

On your domain controller:

- Open Administrative Tools
- Open domain controller security settings
- Then Security Settings
- Local Policies
- Security Options
- Edit the properties of "Microsoft network server: Digitally sign communications (always)"
- Disable this option

Then apply the policy following this article: [Windows: Refresh Security Policies (GPO)]()

### Windows 2000 or newer

Here we'll use CIFS to avoid problems with Windows security. Use this command:

```bash
mount -t cifs -o username=USER,password=PASS //SERVER/SHARE /DESTINATION/
```

Or for newer versions:

```bash
mount -t cifs //SERVER/SHARE /DESTINATION/ -o username=USER,password=PASS
```

Or in `/etc/fstab`:

```
192.168.10.2:/SHARE /home/NAS/LM  cifs  user,noauto,rw,username=user,password=pass,gid=1000,uid=1000 0 0
```

Your share is now mounted.
