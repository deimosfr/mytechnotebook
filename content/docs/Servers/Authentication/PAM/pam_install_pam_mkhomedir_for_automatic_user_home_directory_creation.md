---
weight: 999
url: "/PAM_\\:_Installer_pam_mkhomedir_pour_la_création_automatique_des_home_utilisateurs/"
title: "PAM: Install pam_mkhomedir for Automatic User Home Directory Creation"
description: "How to install and configure pam_mkhomedir on Solaris to automatically create home directories for users at login time"
categories: ["Linux", "Security"]
date: "2009-10-06T08:44:00+02:00"
lastmod: "2009-10-06T08:44:00+02:00"
tags: ["PAM", "Solaris", "Security", "Authentication"]
toc: true
---

## Introduction

You may be in a situation like mine where you have your Solaris connected to an LDAP server or similar, and you're facing the problem of automatic home directory creation for users who connect. After 2-3 hours of struggles, I finally managed to compile and configure everything.

Here's my documentation to save time for others who might face the same issue.

## Prerequisites

First, you'll need a small package to avoid errors like this:

```bash
ld: fatal: file values-Xa.o: open failed: No such file or directory
```

To avoid this problem, install the SUNWarc package:

```bash
pkgadd -d . /cdrom/cdrom0/Solaris_10/Product/SUNWarc
```

You'll also need to install gcc:

```bash
pkg-get install gcc4core gcc4g++
```

Then, we'll add these to our path:

```bash
export PATH=$PATH:/opt/csw/gcc4/bin
export CC=gcc
```

## Compilation

Now we can proceed with the compilation.

```bash
wget http://www.kernel.org/pub/linux/libs/pam/pre/library/Linux-PAM-0.81.tar.bz2
./configure
cp _pam_aconf.h libpam/include/security
cd modules/pammodutil
gcc -c -O2 -D_REENTRANT -DPAM_DYNAMIC -Wall -fPIC -I../../libpam/include -I../../libpamc/include -Iinclude modutil_cleanup.c
gcc -c -O2 -D_REENTRANT -DPAM_DYNAMIC -Wall -fPIC -I../../libpam/include -I../../libpamc/include -Iinclude modutil_ioloop.c
gcc -c -O2 -D_REENTRANT -DPAM_DYNAMIC -Wall -fPIC -I../../libpam/include -I../../libpamc/include -Iinclude modutil_getpwnam.c -D_POSIX_PTHREAD_SEMANTICS
cd ../pam_mkhomedir
gcc -c -O2 -D_REENTRANT -DPAM_DYNAMIC -Wall -fPIC -I../../libpam/include -I../../libpamc/include -I../pammodutil/include pam_mkhomedir.c
/usr/ccs/bin/ld -o pam_mkhomedir.so -B dynamic -G -lc pam_mkhomedir.o ../pammodutil/modutil_*.o
```

## Installation

Installation is quite simple. Still in the source directory, execute:

```bash
cp pam_mkhomedir.so /usr/lib/security/pam_mkhomedir.so.1
cd /usr/lib/security
ln -s pam_mkhomedir.so.1 pam_mkhomedir.so
chown root:bin pam_mkhomedir.so.1
```

That's it! :-)

## Configuration

For configuration, we'll simply edit pam.conf and add this line:

```bash {linenos=table,hl_lines=[4]}
...
# Default definition for Session management
# Used when service name is not explicitly mentioned for session management
#
other   session required        pam_unix_session.so.1
other   session required        pam_mkhomedir.so.1 skel=/etc/skel/ umask=0022
...
```

You can now log in and you'll get this message on the first login of a user:

```
Creating directory '/home/pmavro'.
```

## Resources
- http://www.webservertalk.com/message1674632.html
- http://osdir.com/ml/linux.pam/2006-12/msg00018.html
