---
weight: 999
url: "/Migration_\\:_Migrer_des_comptes_linux_vers_BSD/"
title: "Migration: Migrating Linux accounts to BSD"
description: "Guide for migrating user accounts from Linux systems to BSD systems while maintaining account details and passwords."
categories: ["Linux", "BSD", "System Administration"]
date: "2007-11-23T11:16:00+02:00"
lastmod: "2007-11-23T11:16:00+02:00"
tags: ["Linux", "BSD", "migration", "user accounts", "system administration"]
toc: true
---

## Introduction

Here is a solution that allows you to easily migrate Linux accounts to BSD. **The only constraint is that you must not have two identical logins or identifiers after this migration.**

## Linux

Commands to execute (as root) on your Linux machine for exporting:

- Gathering data from /etc/passwd and /etc/shadow files:

```bash
pwunconv
```

- Transforming the /etc/passwd file to be usable by BSD systems (grep -v 'root\|daemon' excludes root and daemon users):

```bash
cat /etc/passwd | grep -v '^root\|^daemon' | awk -F: '{printf("%s:%s:%s:%s::0:0:%s:%s:%s\n", $1,$2,$3,$4,$5,$6,$7);}' > ~/linux_passwd
```

- Separate the data from /etc/passwd and /etc/shadow again:

```bash
pwconv
```

## BSD

Retrieve the generated file on your BSD system, then for importing:

- Add the content of the linux_passwd file to the end of /etc/master.passwd:

```bash
cat linux_passwd >> /etc/master.passwd
```

- Regenerate the /etc/pwd.db and /etc/spwd.db files:

```bash
pwd_mkdb -p /etc/master.passwd
```
