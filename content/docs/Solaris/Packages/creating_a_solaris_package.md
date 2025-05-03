---
weight: 999
url: "/Créer_un_package_Solaris/"
title: "Creating a Solaris Package"
description: "Learn how to create your own Solaris packages for efficient software distribution and management."
categories: ["Linux"]
date: "2009-11-14T16:04:00+02:00"
lastmod: "2009-11-14T16:04:00+02:00"
tags: ["Servers", "Linux", "Solaris"]
toc: true
---

## 1 Introduction

Creating your own Solaris package can be very practical. Why? It saves time and it's clean :-)

We will need several files, but not all of them are mandatory. I'll try to detail as much as possible to make everything clear.

It's also possible to secure your package by certifying it via SSL, GPG, etc.

## 2 Creating Necessary Files

Let's create a package folder and work in it. Here's the composition of a package:

```bash
package
|-- NAMEofpackage
|   |-- etc
|   |   `-- deimos
|   `-- usr
|       |-- bin
|       |   `-- deimos_cluster
|       `-- perl5
|           |-- 5.8.4
|           |   |-- lib
|           |   |   `-- i86pc-solaris-64int
|           |   |       `-- perllocal.pod
|           |   `-- man
|           |       `-- man3
|           |           `-- Unix::Syslog.3
|           `-- site_perl
|               `-- 5.8.4
|                   `-- i86pc-solaris-64int
|                       |-- Unix
|                       |   `-- Syslog.pm
|                       `-- auto
|                           `-- Unix
|                               `-- Syslog
|                                   |-- Syslog.bs
|                                   `-- Syslog.so
|-- copyright
|-- depend
|-- pkginfo
`-- prototype
```

This includes:

* NAMEofpackage: A package name with 3 to 4 unique capital letters, which will contain all the files necessary for installation
* copyright and other files: Files that will interact during the package usage phase or simply contain information

Let me detail the last point. Here are the files you can use:

* pkginfo: Contains information about the package **(mandatory)**
* prototype: Contains the list of files and directories to copy. Also contains the package configuration files **(mandatory)**
* compver: Defines older versions of packages compatible with this one
* depend: Defines other packets required to install this packet
* space: Defines the disk space required for this package
* copyright: The copyright
* request: Prerequisites for installation
* checkinstall: Script performing tests to verify that everything is fine before installation
* preinstall: Script to prepare for installation
* postinstall: Script to finish installation
* preremove: Script to prepare for package removal
* postremove: Script to finish the removal

**Create a folder containing the desired name of your package and the desired tree structure as shown in the example above.**

### 2.1 copyright

For copyright, put your own. Something simple:

```
Copyright (c) 2009 Deimos
All Rights Reserved
 
This product is protected by copyright and distributed under
licenses restricting copying, distribution, and decompilation.
```

### 2.2 depend

If you have dependencies for your package, you must list them here. For example, in this case, I need to have Perl and associated libraries. I'll gather the list of packet names I need and put them in my depend file:

```bash
$ pkginfo | grep perl584 | grep -v man > depend
system      SUNWperl584core                  Perl 5.8.4 (core)
system      SUNWperl584usr                   Perl 5.8.4 (non-core)
```

Then, we need to replace 'system' with 'P' to get this in our file:

```
P      SUNWperl584core                  Perl 5.8.4 (core)
P      SUNWperl584usr                   Perl 5.8.4 (non-core)
```

### 2.3 pkginfo

Now, the famous file containing the package description:

```bash
ARCH=i386
CATEGORY=application
NAME=Deimos Clustering Solution
PKG=NAMEofpackage
VERSION=1.0
DESC=Clustering Solution for Sun Cluster (SunPlex)
VENDOR=Deimos - http://www.deimos.fr
EMAIL=xxx@mycompany.com
BASEDIR=/
```

### 2.4 prototype

Finally, this file is the most complex. It must contain all existing file names (f), folders (d), and also information (i).

Use only what you need (pkginfo is mandatory):

```bash
echo "i pkginfo=~/package/pkginfo" > prototype
echo "i copyright=~/package/copyright" >> prototype
echo "i depend=~/package/depend" >> prototype
```

Then I'll go to my NAMEofpackage folder and add what's left:

```bash
cd NAMEofpackage
pkgproto . >> ../prototype
```

The content of my file should ultimately look something like this:

```bash
i pkginfo=~/package/pkginfo
i copyright=~/package/copyright
i depend=~/package/depend
d none etc 0755 root sys
d none etc/deimos 0755 root root
d none usr 0755 root root
d none usr/perl5 0755 root root
d none usr/perl5/site_perl 0755 root root
d none usr/perl5/site_perl/5.8.4 0755 root root
d none usr/perl5/site_perl/5.8.4/i86pc-solaris-64int 0755 root root
d none usr/perl5/site_perl/5.8.4/i86pc-solaris-64int/auto 0755 root root
d none usr/perl5/site_perl/5.8.4/i86pc-solaris-64int/auto/Unix 0755 root root
d none usr/perl5/site_perl/5.8.4/i86pc-solaris-64int/auto/Unix/Syslog 0755 root root
f none usr/perl5/site_perl/5.8.4/i86pc-solaris-64int/auto/Unix/Syslog/Syslog.bs 0644 root root
f none usr/perl5/site_perl/5.8.4/i86pc-solaris-64int/auto/Unix/Syslog/Syslog.so 0755 root root
f none usr/perl5/site_perl/5.8.4/i86pc-solaris-64int/auto/Unix/Syslog/.packlist 0644 root root
d none usr/perl5/site_perl/5.8.4/i86pc-solaris-64int/Unix 0755 root root
f none usr/perl5/site_perl/5.8.4/i86pc-solaris-64int/Unix/Syslog.pm 0644 root root
d none usr/perl5/5.8.4 0755 root root
d none usr/perl5/5.8.4/man 0755 root root
d none usr/perl5/5.8.4/man/man3 0755 root root
f none usr/perl5/5.8.4/man/man3/Unix::Syslog.3 0644 root root
d none usr/perl5/5.8.4/lib 0755 root root
d none usr/perl5/5.8.4/lib/i86pc-solaris-64int 0755 root root
f none usr/perl5/5.8.4/lib/i86pc-solaris-64int/perllocal.pod 0644 root root
d none usr/bin 0755 root root
f none usr/bin/deimos_cluster 0740 root root
```

## 3 Creating the Package

Now we'll create the package. We'll make a package_done folder to create the package in:

```bash
cd ~/package
mkdir package_done
pkgmk -o -b . -d package_done -f prototype
```

At this point, we have our Solaris package in folder form and not with a pkg extension (which is classier, more beautiful, and glows in the dark). That's why here's the final step:

```bash
cd package_done
pkgtrans -s . NAMEofpackage.pkg NAMEofpackage
```

And now I have my beautiful package in .pkg format :-)

## 4 Resources

[https://docs.sun.com/app/docs/doc/817-0406/ch1designpkg-51728?a=view](https://docs.sun.com/app/docs/doc/817-0406/ch1designpkg-51728?a=view)
