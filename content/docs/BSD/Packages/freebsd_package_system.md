---
weight: 999
url: "/Le_système_de_Packages_FreeBSD/"
title: "FreeBSD Package System"
description: "A guide to understanding and using FreeBSD's package management system including both the new and old methods, portage system, and system updates."
categories: ["FreeBSD", "Linux"]
date: "2014-07-18T19:39:00+02:00"
lastmod: "2014-07-18T19:39:00+02:00"
tags:
  [
    "FreeBSD",
    "Package Management",
    "System Administration",
    "Portage",
    "Security Updates",
  ]
toc: true
---

## Introduction

FreeBSD is one of the most widely used BSD distributions for servers as it is very up-to-date and offers a very comprehensive port system (more than 16,000 ports available). Additionally, it integrates advanced functions at the source level.

This allows you, for example, to modify parameters to best adapt to your needs, similar to Gentoo.

## Precompiled Packages

### New Method

The latest method involves using pkgng. To set it up, you need to convert the current database:

```bash
pkg2ng
```

Then modify/add this line:

```bash
/etc/make.conf
```

And modify the base repository if you cannot access it (as at the time of writing, a security incident has forced the removal of binary packages from the official site):

```bash
packagesite: http://mirror.exonetric.net/pub/pkgng/${ABI}/latest
#packagesite: http://pkgbeta.FreeBSD.org/freebsd:9:x86:32/latest
```

All that remains is to update the repositories:

```bash
pkg update
```

And to install software:

```bash
pkg install <software>
```

### Old Method

#### Adding Software

If I want to install lsof with FreeBSD packages:

```bash
pkg_add -r lsof
```

#### Finding Installed Software

pkg_version is a utility that summarizes the versions of all pre-compiled software installed:

```bash
pkg_version
```

If you want a description of the software installed on your machine:

```bash
pkg_info
```

{{< table "table-hover table-striped" >}}
| Symbols | Meanings |
|---------|----------|
| = | The installed pre-compiled software version is equivalent to that found in the local ports catalog. |
| < | The installed version is older than the one available in the ports catalog. |
| > | The installed version is newer than the one found in the local ports catalog. (The local ports catalog is probably outdated) |
| ? | The pre-compiled software cannot be found in the ports catalog index. (This can happen when, for example, installed software is removed from the ports catalog or renamed.) |
| \* | There are multiple versions of this pre-compiled software. |
{{< /table >}}

#### Deleting a Package

To delete an installed package:

```bash
pkg_delete
```

## The Portage System

### Searching for a Package

For example, if you are looking for lsof:

```
# cd /usr/ports
# make search name=lsof
Port:   lsof-4.56.4
Path:   /usr/ports/sysutils/lsof
Info:   Lists information about open files (similar to fstat(1))
Maint:  obrien@FreeBSD.org
Index:  sysutils
B-deps:
R-deps:
```

You need to use \*make search key=**string\***.

## Updating Security Patches

To update security patches:

```bash
freebsd-update fetch
freebsd-update install
```

If you want to rollback:

```bash
freebsd-update rollback
```

You can also be notified when updates are available by adding this line to the crontab:

```bash
@daily root freebsd-update cron
```

## Updating Your FreeBSD

To upgrade from one version to another, run this command indicating the version you want to upgrade to (here 10.0):

```bash
freebsd-update -r 10.0-RELEASE upgrade
```

This command will download updates and merge them. Then reboot, apply the updates:

```bash
freebsd-update install
```

Reboot and run this command again.
