---
weight: 999
url: "/MariaDB_\\:_Migration_depuis_MySQL/"
title: "MariaDB: Migration from MySQL"
description: "A guide on how to migrate from MySQL to MariaDB, including installation and configuration steps."
categories: ["Debian", "Storage", "Database"]
date: "2014-04-01T06:22:00+02:00"
lastmod: "2014-04-01T06:22:00+02:00"
tags: ["MariaDB", "MySQL", "Database", "Migration", "Debian"]
toc: true
---

![MariaDB](/images/mariadb-logo.avif)

{{< table "table-hover table-striped" >}}
|||
|---|---|
| **Software version** | 10 |
| **Operating System** | Debian 7 |
| **Website** | [MariaDB Website](https://mariadb.org/) |
| **Last Update** | 01/04/2014 |
{{< /table >}}

## Introduction

MariaDB is a community-developed fork of the MySQL relational database management system, the impetus being the community maintenance of its free status under the GNU GPL. As a fork of a leading open source software system, it is notable for being led by its original developers and triggered by concerns over direction by an acquiring commercial company Oracle. Contributors are required to share their copyright with Monty Program AB.

The intent is also to maintain high compatibility with MySQL, ensuring a "drop-in" replacement capability with library binary equivalency and exacting matching with MySQL APIs and commands. It includes the XtraDB storage engine as a replacement for InnoDB, as well as a new storage engine, Aria, that intends to be both a transactional and non-transactional engine perhaps even included in future versions of MySQL.

Its lead developer is Michael "Monty" Widenius, the founder of MySQL and Monty Program AB. He had previously sold his company, MySQL AB, to Sun Microsystems for 1 billion USD. MariaDB is named after Monty's younger daughter, Maria.

For a migration from MySQL to MariaDB, it's recommended to keep the same version (Example: MySQL 5.1 & MariaDB 5.1). Then you can upgrade to upper versions of MariaDB.

## Installation

To install MariaDB, it's unfortunately not embedded in Debian, so we'll add a repository. First of all, install a python tool to get aptkey:

```bash
aptitude install python-software-properties
```

Then let's add this repository (https://downloads.mariadb.org/mariadb/repositories/):

```bash
apt-key adv --recv-keys --keyserver keyserver.ubuntu.com 0xcbcb082a1bb943db
add-apt-repository 'deb http://mirrors.linsrv.net/mariadb/repo/10.0/debian wheezy main'
```

We're now going to change apt pinning to prioritize MariaDB's repository (`/etc/apt/preferences.d/mariadb`):

```bash
Package: *
Pin: release o=MariaDB
Pin-Priority: 1000
```

You can delete MySQL if it was already installed:

```bash
aptitude remove mysql-server mysql-client
```

Then we install MariaDB:

```bash
aptitude update
aptitude install mariadb-server
```

That's it, nothing else to do! Your MySQL databases are still accessible like if it was MySQL, but you're running MariaDB :-)

## References

- http://mariadb.org/
- http://www.skysql.com/
