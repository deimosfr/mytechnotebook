---
weight: 999
url: "/Yum_\\:_utilisation_des_packages_sous_RedHat/"
title: "Yum: Package Management in Red Hat"
description: "A guide on using Yum package manager in Red Hat-based distributions, including installation, removal, updates, and other package management operations."
categories: ["Red Hat", "Linux"]
date: "2012-03-02T19:08:00+02:00"
lastmod: "2012-03-02T19:08:00+02:00"
tags: ["package management", "rpm", "red hat", "fedora", "yum"]
toc: true
---

## Introduction

Yum, which stands for Yellow dog Updater Modified, is a package manager for Linux distributions like Fedora and Red Hat Enterprise Linux, created by Yellow Dog Linux.

It allows you to manage the installation and updating of software installed on a distribution. It's a layer on top of RPM that handles downloads and dependencies, similar to Debian's APT or Mandriva's Urpmi.

## Usage

- Installing a package:

```bash
yum install <package>
```

- Reinstalling a package:

```bash
yum reinstall <package>
```

- Installing a local RPM:

```bash
yum localinstall <package.rpm>
```

- Removing a package:

```bash
yum remove <package>
```

- Updating packages or a specific package:

```bash
yum update <package>
```

- Getting information about a package:

```bash
yum info <package>
```

- Installing a package group:

```bash
yum groupinstall <group>
```

- Viewing available packages (installed or not):

```bash
yum list
```

or

```bash
yum list htt*
```

- Viewing available package groups:

```bash
yum grouplist
```

- Viewing repositories:

```bash
yum repolist
```

- Finding which package a file belongs to (equivalent of apt-file):

```bash
yum provides <file/command>
```

or

```bash
yum whatprovides <file/command>
```

- Ignoring missing GPG key:

```bash
yum ... --nogpg
```

- Listing all installed packages:

```bash
yum list installed
```

- Forcing the protection of an rpm on a specific redhat version:

```bash
yum protectbase <package>
```

- Checking compatibility:

```bash
yum verify <package>
```

- Downloading packages only

You will need the yum-downloadonly package first to have this option in yum:

```bash
yum install yum-downloadonly
```

Then to download the package:

```bash
yum install -y postfix --downloadonly
```

- Viewing the contents of a package:

```bash
rpm -qla <package>
```

or

```bash
repoquery -qla <package>
```
