---
weight: 999
url: "/Patcher_sa_Solaris/"
title: "Patching Solaris"
description: "A guide on how to detect, retrieve and install patches for Solaris systems."
categories: 
  - Linux
  - Servers
date: "2006-11-28T18:14:00+02:00"
lastmod: "2006-11-28T18:14:00+02:00"
tags: 
  - Solaris
  - Servers
  - Network
  - cd ~
toc: true
---

## Introduction

I won't go into detailed explanation of what a patch is, but you should know that by applying patches (fixes), you can eliminate bugs and security vulnerabilities.

## Detecting Patches

```bash
showrev -p
```

```
Patch: 106793-01 Obsoletes:  Requires:  Incompatibles: Packages: SUNWhea . . .
```

```bash
patchadd -p
```

```
Patch: 106793-01 Obsoletes: Requires: Incompatibles: Packages: SUNWhea 
. . .
```

These versions should be identical. If they're not, it means a patch has been applied.

You can also see the different system patches here:

```bash
ls /var/sadm/patch
```

```
107558-05  107594-04  107630-01  107663-01  107683-01 107696-01
107817-01  107582-01  107612-06  107640-03
```

## Retrieving Patches

For France, here is the address for Solaris patches: http://sunsolve.sun.fr

Then, on Sun's FTP site, get the latest version:

```bash
cd /var/tmp
```

```bash
ftp sunsolve.sun.com
Connected to sunsolve.sun.com.
(output omitted)
Name (sunsolve:usera): anonymous
331 Guest login ok, send your complete e-mail address as password.
Password: yourpassword
(output omitted)
ftp> bin
200 Type set to I.
ftp> cd /patchroot/reports
ftp> get public_patch_report
(output omitted)
ftp> cd /patchroot/clusters
ftp> get 10_SunAlert_Patch_Cluster.README
(output omitted)
ftp> cd /patchroot/current_unsigned
ftp> mget 112605*
mget 112605-01.zip? y
(output omitted)
mget 112605.readme? y
ftp> bye
```

Decompress the patch:

```bash
/usr/bin/unzip 105050-01.zip
```

## Implementing Patches

Here are the existing commands:

- patchadd - Install a patch
- patchrm - Remove a patch
- smpatch - utility to download and install a patch

### patchadd

Let's install the patch:

```bash
cd /var/tmp
```

```bash
patchadd 105050-01
Checking installed patches...
Verifying sufficient filesystem capacity (dry run method)
Installing patch packages...
Patch number 105050-01 has been successfully installed.
See /var/sadm/patch/105050-01/log for details.
Patch packages installed:
  SUNWhea
```

### pathrm

To remove a patch:

```bash
patchrm 105050-01
```

```bash
Checking installed packages and patches...
Backing out patch 105050-01...
Patch 105050-01 has been backed out.
```

Don't forget to restart the machine after applying a patch.

### smpatch

There's also another utility that allows you to automatically download and install (or remove) a patch:

```bash
smpatch get -L patchpro.patch.source patchpro.download.directory
```

```bash
https://updateserver.sun.com/solaris/
/var/sadm/spool
```
