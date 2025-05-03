---
weight: 999
url: "/Le_système_de_Packages_OpenBSD/"
title: "OpenBSD Package System"
description: "A guide to understanding and using the OpenBSD package and ports systems for software management."
categories: ["Security", "Linux"]
date: "2007-11-27T18:07:00+02:00"
lastmod: "2007-11-27T18:07:00+02:00"
tags: ["BSD", "OpenBSD", "Security", "Package Management", "System Administration"]
toc: true
---

## Introduction

Transitioning to OpenBSD isn't very straightforward when coming from the Linux world (PS: if you're coming from Windows, it's better to stay there or switch to Linux first).

In OpenBSD, there are two package systems:

- The first one is provided by the base OpenBSD system (pkg), containing packages with almost no security vulnerabilities (just a reminder: only 2 vulnerabilities discovered in 10 years). These packages are pre-compiled.
- The second contains many more software applications (port), but according to OpenBSD, they can compromise system stability and security. That said, it's still preferable to use these rather than compiling yourself from various sources since these packages have still been validated by the OpenBSD team. These packages are compiled on demand.

## PGK

### Preparation

First, you need to [choose an FTP server](https://www.openbsd.org/fr/ftp.html) to specify the repository. I chose a French server:

```bash
export PKG_PATH=ftp://ftp.arcane-networks.fr/pub/OpenBSD/`uname -r`/packages/`machine -a`/
```

To make this permanent, add this line to your `.profile` file.

### Installing a package

To install packages, there are several arguments available. I use `-i` for interactive mode (available since OpenBSD 3.9) and the `-v` option for verbose output. To install postfix, I do:

```bash
$ pkg_add -iv postfix
Ambiguous: postfix could be postfix-2.3.2 postfix-2.3.2-ldap postfix-2.3.2-mysql postfix-2.3.2-sasl2 postfix-2.4.20060727 postfix-2.4.20060727-ldap postfix-2.4.20060727-mysql postfix-2.4.20060727-sasl2
Choose one package
         0: <None>
         1: postfix-2.3.2
         2: postfix-2.3.2-ldap
         3: postfix-2.3.2-mysql
         4: postfix-2.3.2-sasl2
         5: postfix-2.4.20060727
         6: postfix-2.4.20060727-ldap
         7: postfix-2.4.20060727-mysql
         8: postfix-2.4.20060727-sasl2
Your choice:
```

I can then choose the version number I want to install. Otherwise, I can directly specify what I want:

```bash
pkg_add -v postfix-2.4.20060727
```

### Getting information

To get information about my installed packages, I can use:

```bash
$ pkg_info
screen-4.0.3        multi-screen window manager
tcl-8.4.7p1         Tool Command Language
vim-7.0.42-no_x11   vi clone, many additional features
zsh-4.2.6p0         Z shell, Bourne shell-compatible
```

### Updating packages

To update packages, just use the `pkg_add` command with an argument and the name of the package to update:

```bash
pkg_add -u screen
```

If you add the `-c` option, it will overwrite the current configuration files with the default ones.

### Removing a package

Again, it's not very complicated, use the `pkg_delete` command followed by the package to remove:

```bash
pkg_delete screen
```

### Upgrade

When there's a version upgrade, remember to modify your repository line, then do this to update all packages:

```bash
pkg_add -ui -F update -F updatedepends
```

## Port

### Getting the ports tree

Let's get the tar.gz and decompress it:

```bash
cd /tmp
ftp ftp://ftp.openbsd.org/pub/OpenBSD/`uname -r`/ports.tar.gz
cd /usr
tar xzf /tmp/ports.tar.gz
```

### Configuring the ports system

If you've created a new user, add them to the sudoers list, and give them the rights:

```bash
chgrp -R wsrc /usr/ports
find /usr/ports -type d -exec chmod g+w {} \;
```

Next, let's add a few lines to the mk.conf file:

```bash
echo "USE_SYSTRACE=Yes" >> /etc/mk.conf
echo "WRKOBJDIR=/usr/obj/ports" >> /etc/mk.conf
echo "DISTDIR=/usr/distfiles" >> /etc/mk.conf
echo "PACKAGE_REPOSITORY=/usr/packages" >> /etc/mk.conf
```

### Searching for packages

To perform a search:

```bash
cd /usr/ports
make search key=nmap
```

```bash
Port:   nmap-4.11
Path:   net/nmap
Info:   scan ports and fingerprint stack of network hosts
Maint:  Okan Demirmen <okan@demirmen.com>
Index:  net security
L-deps: dnet::net/libdnet gdk_pixbuf-2.0.0.0,gdk-x11-2.0.0.0,gtk-x11-2.0.0.0::x11/gtk+2 iconv.>=4::converters/libiconv intl.>=3:gettext->=0.10.38:devel/gettext pcre::devel/pcre
B-deps: :devel/gmake gettext->=0.14.5:devel/gettext pkgconfig-*:devel/pkgconfig
R-deps: gettext->=0.10.38:devel/gettext
Archs:  any
 
Port:   nmap-4.11-no_x11
Path:   net/nmap,no_x11
Info:   scan ports and fingerprint stack of network hosts
Maint:  Okan Demirmen <okan@demirmen.com>
Index:  net security
L-deps: dnet::net/libdnet pcre::devel/pcre
B-deps: :devel/gmake
R-deps: 
Archs:  any
```

Here I don't have a graphical interface, so I don't need X11, the no_x11 option is interesting to me.

### Installing a package

To install a package, just go to the right section and start the compilation:

```bash
cd /usr/ports/net/nmap/
```

Now let's see the available options:

```bash
$ make show=FLAVORS
no_x11
```

Once again, we can see that I can compile nmap without the graphical interface. I pass my argument before compiling, then run make install:

```bash
env FLAVOR="no_x11" make install
```

### Cleaning up after compilation

You probably want to clean up the default working directory of the port after building and installing the package.

```bash
$ make clean
===>  Cleaning for rsnapshot-1.2.9
```

You can also clean the working directories of all the port's dependencies with this make target:

```bash
$ make clean=depends
===>  Cleaning for rsync-2.6.8
===>  Cleaning for rsnapshot-1.2.9
```

If you want to remove all the port's distribution sources, you can use:

```bash
$ make clean=dist
===>  Cleaning for rsnapshot-1.2.9
===>  Dist cleaning for rsnapshot-1.2.9
```

In case you've compiled multiple flavors of the same port, you can clean the working directories of all these flavors at once using:

```bash
$ make clean=flavors
```
