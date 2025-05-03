---
weight: 999
url: "/PAM_mount\\:_Monter_des_partages_réseaux_au_login/"
title: "PAM mount: Automatically Mount Network Shares at Login"
description: "Guide on how to use PAM mount to automatically mount network shares when users log in to a system."
categories: ["Linux", "Security"]
date: "2008-07-16T16:23:00+02:00"
lastmod: "2008-07-16T16:23:00+02:00"
tags: ["PAM", "Network", "Mount", "NFS", "Security"]
toc: true
---

## Introduction

You may need sometimes to automatically mount network shares. This can be done with the [pam_mount](https://pam-mount.sourceforge.net/) module.

## Installation

For installation:

```bash
apt-get install pam_mount
```

## Configuration

### pam_mount.conf

Edit the `/etc/security/pam_mount.conf` file and configure what you need. Here we would like users to have nfs home share to be mounted from the server at logon. Add this line at the end of the file:

```bash
volume * nfs my_nfs_server /home/& ~ - - -
```

You may also need to have other mount points for other things like smb, cifs or fuse:

```bash
volume user smbfs krueger public /home/user/krueger - - -
volume * fuse - "sshfs#&@fileserver:" /home/& - - -
```

### Application

You must choose an application to mount automatically your share. For example I choose SSH. When a user logs into SSH it must mount the NFS share, so edit this file:

```bash {linenos=table,hl_lines=[3,7]}
...
auth       required     pam_env.so envfile=/etc/default/locale
auth       required    pam_mount.so
@include   common-auth
...
@include   common-account
session    required pam_mount.so
@include   common-session
...
```

Now when you'll login, it will automatically mount your home.
