---
weight: 999
url: "/Créer_un_repository_Debian/"
title: "Creating a Debian Repository"
description: "This guide explains how to create your own Debian package repository for hosting custom packages."
categories: ["Ubuntu", "Debian", "Linux"]
date: "2007-12-11T20:46:00+02:00"
lastmod: "2007-12-11T20:46:00+02:00"
tags: ["Servers", "Linux", "Debian"]
toc: true
---

## Introduction

I am a co-designer of the MySecureShell project. The problem is that getting your project accepted in the official Debian/Ubuntu repositories requires significant effort. Therefore, while waiting, we decided to create our own repository. Here, I will describe the steps to create your own repository.

## Preparation

This article assumes that the packages to be made available on the repository are already generated. Assuming the package is called "mysecureshell", you should have the following files:

- mysecureshell.orig.tar.gz
- mysecureshell.diff.gz
- mysecureshell.dsc
- mysecureshell.changes
- mysecureshell.deb

## Repository Generation

### Generating the Structure

First, you need to generate the repository tree structure with the following commands:

```bash
mkdir -p /var/www/mss/debian/dists/testing/main/{binary-i386,source}
```

You need to copy your package files to your repository:

```bash
cp mysecureshell_1.0.dsc mysecureshell_1.0_i386.deb /var/www/mss/debian/dists/testing/main/binary-i386/
cp mysecureshell.diff.gz mysecureshell.dsc mysecureshell.orig.tar.gz mysecureshell.orig.changes /var/www/mss/debian/dists/testing/main/source
```

### Generating Repository Files

Then you need to generate the two files Packages.gz and Sources.gz required for the repository:

```bash
cd /var/www/mss/debian/dists/testing/main
dpkg-scanpackages binary-i386 /dev/null dists/testing/main/ | gzip -f9 > binary-i386/Packages.gz
dpkg-scansources source /dev/null dists/testing/main/ | gzip -f9 > source/Sources.gz
```

### Generating Description Files

These two files must be regenerated every time you need to put a new version of your package on the repository.

Finally, you need to create two description files for your repository. The first file should be placed in the binary-i386 directory, called Release, and should contain:

```
Archive: testing
Version: 1.0
Component: main
Origin: MySecureShell
Label: mysecureshell
Architecture: i386
```

The second file should be placed in the source directory, also called Release, and should contain:

```
Archive: testing
Version: 1.0
Component: main
Origin: MySecureShell
Label: mysecureshell
Architecture: source
```

Your Debian repository is now ready! Now you just need to deploy it on your HTTP server (I'll let you do that ;-)).

## Usage

Creating a Debian package repository is fine, but you need to know how to use it.

Users who want to use your repository need to add one of the following two lines to their `/etc/apt/sources.list` file:

```bash
deb http://mysecureshell.free.fr/debian testing main
deb-src http://mysecureshell.free.fr/debian testing main
```

Then the procedure is the same as usual for the package management system to know all the packages available on your repositories:

```bash
apt-get update
```

Finally, installing a package from the repository is done with the usual command for all Debian users:

```bash
apt-get install mysecureshell
```

## References

http://www.debian.org/doc/manuals/repository-howto/repository-howto.fr.html  
http://www.debianaddict.org/article31.html
