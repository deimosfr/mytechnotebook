---
weight: 999
url: "/Trouver_toutes_les_dépendances_liées_à_un_package/"
title: "Find all dependencies related to a package"
description: "How to find all dependencies related to a package in Debian-based systems using apt-rdepends"
categories: ["Linux", "Debian"]
date: "2008-09-24T15:00:00+02:00"
lastmod: "2008-09-24T15:00:00+02:00"
tags: ["apt", "dependencies", "package management", "Debian", "Ubuntu"]
toc: true
---

## Introduction

Sometimes it's very practical to find all dependencies of a package. Here's a useful command for this purpose.

## Installation

Installation:

```bash
apt-get install apt-rdepends
```

Example search:

```bash
$ apt-rdepends libapache2-mod-php5
Reading package lists... Done
Building dependency tree... Done
libapache2-mod-php5
  Depends: apache2-mpm-itk
  Depends: apache2-mpm-prefork (>> 2.0.52)
  Depends: apache2.2-common
  Depends: libbz2-1.0
  Depends: libc6 (>= 2.3.6-6)
  Depends: libcomerr2 (>= 1.33-3)
  Depends: libdb4.4
  Depends: libkrb53 (>= 1.4.2)
  Depends: libmagic1
  Depends: libpcre3 (>= 4.5)
  Depends: libssl0.9.8 (>= 0.9.8c-1)
  Depends: libxml2 (>= 2.6.27)
  Depends: mime-support (>= 2.03-1)
  Depends: php5-common (= 5.2.0-8+etch1)
  Depends: ucf
  Depends: zlib1g (>= 1:1.2.1)
apache2-mpm-itk
  Depends: apache2.2-common (= 2.2.3-4+etch5)
  Depends: libapr1
  Depends: libaprutil1
  Depends: libc6 (>= 2.3.6-6)
  Depends: libcap1
  Depends: libdb4.4
  Depends: libexpat1 (>= 1.95.8)
  Depends: libldap2 (>= 2.1.17-1)
  Depends: libpcre3 (>= 4.5)
  Depends: libpq4 (>= 8.1.4)
  Depends: libsqlite3-0 (>= 3.3.8)
  Depends: libuuid1
apache2.2-common
  Depends: apache2-utils
  Depends: libmagic1
  Depends: lsb-base
  Depends: mime-support
  Depends: net-tools
  Depends: procps
apache2-utils
  Depends: libapr1
  Depends: libaprutil1
  Depends: libc6 (>= 2.3.6-6)
  Depends: libdb4.4
  Depends: libexpat1 (>= 1.95.8)
  Depends: libldap2 (>= 2.1.17-1)
  Depends: libpcre3 (>= 4.5)
  Depends: libpq4 (>= 8.1.4)
  Depends: libsqlite3-0 (>= 3.3.8)
  Depends: libssl0.9.8 (>= 0.9.8c-1)
  Depends: libuuid1
libapr1
  Depends: libc6 (>= 2.3.6-6)
  Depends: libuuid1
libc6
  Depends: tzdata
tzdata
libuuid1
  Depends: libc6 (>= 2.3.6-6)
libaprutil1
  Depends: libapr1
  Depends: libc6 (>= 2.3.6-6)
  Depends: libdb4.4
  Depends: libexpat1 (>= 1.95.8)
  Depends: libldap2 (>= 2.1.17-1)
  Depends: libpq4 (>= 8.1.4)
  Depends: libsqlite3-0 (>= 3.3.7)
  Depends: libuuid1
libdb4.4
  Depends: libc6 (>= 2.3.6-6)
libexpat1
  Depends: libc6 (>= 2.3.6-6)
libldap2
  Depends: libc6 (>= 2.3.6-6)
  Depends: libgnutls13 (>= 1.4.0-0)
  Depends: libsasl2-2
libgnutls13
  Depends: libc6 (>= 2.3.6-6)
  Depends: libgcrypt11 (>= 1.2.2)
  Depends: libgpg-error0 (>= 1.4)
  Depends: liblzo1
  Depends: libopencdk8 (>= 0.5.8)
  Depends: libtasn1-3 (>= 0.3.4)
  Depends: zlib1g (>= 1:1.2.1)
libgcrypt11
  Depends: libc6 (>= 2.3.6-6)
  Depends: libgpg-error0 (>= 1.2)
libgpg-error0
  Depends: libc6 (>= 2.3.6-6)
liblzo1
  Depends: libc6 (>= 2.3.5-1)
libopencdk8
  Depends: libc6 (>= 2.3.6-6)
  Depends: libgcrypt11 (>= 1.2.2)
  Depends: libgpg-error0 (>= 1.4)
  Depends: zlib1g (>= 1:1.2.1)
zlib1g
  Depends: libc6 (>= 2.3.6-6)
libtasn1-3
  Depends: libc6 (>= 2.3.6-6)
libsasl2-2
  Depends: libc6 (>= 2.3.6-6)
  Depends: libdb4.2
libdb4.2
  Depends: libc6 (>= 2.3.6-6)
libpq4
  Depends: libc6 (>= 2.3.6-6)
  Depends: libcomerr2 (>= 1.33-3)
  Depends: libkrb53 (>= 1.4.2)
  Depends: libssl0.9.8 (>= 0.9.8c-1)
libcomerr2
  Depends: libc6 (>= 2.3.6-6)
libkrb53
  Depends: libc6 (>= 2.3.6-6)
  Depends: libcomerr2 (>= 1.33-3)
libssl0.9.8
  Depends: debconf (>= 0.5)
  Depends: debconf-2.0
  Depends: libc6 (>= 2.3.6-6)
  Depends: zlib1g (>= 1:1.2.1)
debconf
  Depends: debconf-english
  Depends: debconf-i18n
  PreDepends: perl-base (>= 5.6.1-4)
debconf-english
  Depends: debconf
debconf-i18n
  Depends: debconf
  Depends: liblocale-gettext-perl
  Depends: libtext-charwidth-perl
  Depends: libtext-iconv-perl
  Depends: libtext-wrapi18n-perl
liblocale-gettext-perl
  Depends: libc6 (>= 2.3.2.ds1-21)
  PreDepends: perl-base (>= 5.8.7-3)
  PreDepends: perlapi-5.8.7
perl-base
  PreDepends: libc6 (>= 2.3.6-6)
perlapi-5.8.7
libtext-charwidth-perl
  Depends: libc6 (>= 2.3.6-6)
  Depends: perl-base (>= 5.8.8-6)
  Depends: perlapi-5.8.8
perlapi-5.8.8
libtext-iconv-perl
  Depends: libc6 (>= 2.3.6-6)
  Depends: perl-base (>= 5.8.8-6)
  Depends: perlapi-5.8.8
libtext-wrapi18n-perl
  Depends: libtext-charwidth-perl
debconf-2.0
libsqlite3-0
  Depends: libc6 (>= 2.3.6-6)
libpcre3
  Depends: libc6 (>= 2.3.6-6)
libmagic1
  Depends: libc6 (>= 2.3.6-6)
  Depends: zlib1g (>= 1:1.2.1)
lsb-base
  Depends: ncurses-bin
  Depends: sed
ncurses-bin
  PreDepends: libc6 (>= 2.3.6-6)
  PreDepends: libncurses5 (>= 5.4-5)
libncurses5
  Depends: libc6 (>= 2.3.6-6)
sed
  PreDepends: libc6 (>= 2.3.6-6)
mime-support
net-tools
  Depends: libc6 (>= 2.3.2.ds1-21)
procps
  Depends: libc6 (>= 2.3.6-6)
  Depends: libncurses5 (>= 5.4-5)
  Depends: lsb-base (>= 3.0-10)
libcap1
  Depends: libc6 (>= 2.3.2.ds1-4)
apache2-mpm-prefork
  Depends: apache2.2-common (= 2.2.3-4)
  Depends: libapr1
  Depends: libaprutil1
  Depends: libc6 (>= 2.3.6-6)
  Depends: libdb4.4
  Depends: libexpat1 (>= 1.95.8)
  Depends: libldap2 (>= 2.1.17-1)
  Depends: libpcre3 (>= 4.5)
  Depends: libpq4 (>= 8.1.4)
  Depends: libsqlite3-0 (>= 3.3.8)
  Depends: libuuid1
libbz2-1.0
  Depends: libc6 (>= 2.3.6-6)
libxml2
  Depends: libc6 (>= 2.3.6-6)
  Depends: zlib1g (>= 1:1.2.1)
php5-common
  Depends: sed (>= 4.1.1-1)
ucf
  Depends: coreutils (>= 5.91)
  Depends: debconf (>= 1.2.0)
  Depends: debconf-2.0
coreutils
  PreDepends: libacl1 (>= 2.2.11-1)
  PreDepends: libc6 (>= 2.3.6-6)
  PreDepends: libselinux1 (>= 1.32)
libacl1
  Depends: libattr1 (>= 2.4.4-1)
  Depends: libc6 (>= 2.3.6-6)
libattr1
  Depends: libc6 (>= 2.3.5-1)
libselinux1
  Depends: libc6 (>= 2.3.6-6)
  Depends: libsepol1 (>= 1.14)
libsepol1
  Depends: libc6 (>= 2.3.6-6)
```
