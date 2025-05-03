---
weight: 999
url: "/Xtrabackup_\\:_Optimiser_ses_backups_MySQL/"
title: "XtraBackup: Optimizing Your MySQL Backups"
description: "A guide on how to optimize MySQL backups using Percona XtraBackup, including installation, full backups, incremental backups, and restoration procedures."
categories: ["Backup", "Linux", "Debian"]
date: "2013-12-27T08:43:00+02:00"
lastmod: "2013-12-27T08:43:00+02:00"
tags: ["XtraBackup", "MySQL", "Backup", "Database", "Percona"]
toc: true
---

![XtraBackup](/images/percona-xtrabackup-logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 2.1.2-611.wheezy |
| **Operating System** | Debian 7 |
| **Website** | [XtraBackup Website](https://www.percona.com/doc/percona-xtrabackup/) |
| **Last Update** | 27/12/2013 |
{{< /table >}}

## Introduction

[XtraBackup](https://www.percona.com/doc/percona-xtrabackup/) is an open-source solution for optimizing your MySQL backups. It's much faster than mysqldump since it works directly with files rather than SQL queries.

In this article, we'll see how to use it. It's important to know that XtraBackup is particularly effective with the InnoDB engine.

## Installation

Installing it on Debian is very simple. First, we'll add the repository:

```bash
apt-key adv --keyserver keys.gnupg.net --recv-keys 1C4CBDCDCD2EFD2A
```

Create this file to add the repository:

```bash
# /etc/apt/sources.list.d/percona.list
deb http://repo.percona.com/apt VERSION main
deb-src http://repo.percona.com/apt VERSION main
```

Update and install:

```bash
aptitude update
aptitude install xtrabackup
```

## Usage

It can be useful to store the credentials in a file if you want to back up locally:

```bash
# /etc/mysql/conf.d/xtrabackup.cnf
[xtrabackup]
user=<user>
password=<password>
```

You can also use a file like ~/.my.cnf.

### Backup (full)

To make a full backup:

```bash
innobackupex --user=xxxxx --password=xxxx --databases=database_name directory_to_store_backup
```

The backup will be stored in a folder named after its creation timestamp (example: 2011-09-19_09-32-19).

It's possible to add the '--apply-log' option, which will save more information to allow a simple restore by copying the folder to `/var/lib/mysql`:

```bash
innobackupex --apply-log --user=xxxxx --password=xxxx --databases=database_name directory_to_store_backup
```

{{< alert context="info" text="With the --apply-log option, it's impossible to make incremental backups with XtraBackup" />}}

### Backup (incremental)

A full backup is required to be able to make an incremental backup. Once you have one, use this command to generate an incremental backup:

```bash
innobackupex --incremental directory_to_store_backup --incremental-basedir=directory_containing_the_full --user=root --password=xxxxx
```

### Restoration

#### From a Full backup

To restore from a full backup, the operation is done in three parts:

1. Use the apply-log argument:

```bash
innobackupex --apply-log directory_of_the_full_backup
```

2. Move the existing database to avoid any unwanted residue:

```bash
mv /var/lib/mysql/database_directory{,.old}
```

3. Run this command to restore (copy-back argument):

```bash
innobackupex --ibbackup=xtrabackup --copy-back directory_of_the_full_backup
```

#### From an incremental backup

To perform a restoration from an incremental backup, you again need several steps:

1. Prepare the full backup. Here, use-memory is optional; it just speeds up the process:

```bash
innobackupex --apply-log --redo-only directory_of_the_full_backup [--use-memory=1G] --user=root --password=xxxxx
```

2. Apply the incremental backups, in order:

```bash
innobackupex --apply-log directory_of_the_full_backup --incremental-dir=directory_of_the_incremental_backup [--use-memory=1G] --user=root --password=xxxxx
```

3. Prepare the final backup:

```bash
innobackupex-1.5.1 --apply-log directory_of_the_full_backup [--use-memory=1G] --user=root --password=xxxxx
```

4. Restore the final backup:

```bash
mv /var/lib/mysql/database_directory{,.old}
innobackupex --ibbackup=xtrabackup --copy-back directory_of_the_full_backup
```

### Backup and restore from the slave

Here is a solution when you do not have enough space on the local master to store backups. Simply use Netcat to grab the backup directly from the slave:

```bash
mkdir /tmp/backups/
nc -l -p 8080 | tar xfi - -C /tmp/backups/
```

On the master:

```bash
innobackupex --stream=tar /tmp/ --slave-info | nc sql-slave 8080
```

Apply the logs on the slave:

```bash
innobackupex --apply-log --ibbackup=xtrabackup /tmp/backups/
```
