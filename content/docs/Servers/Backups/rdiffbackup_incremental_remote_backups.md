---
weight: 999
url: "/Rdiff-backup_\\:_Sauvegardes_distantes_incrémentielles/"
title: "Rdiff-backup: Incremental Remote Backups"
description: "Guide for setting up and using rdiff-backup for incremental remote backups over SSH."
categories:
  - "Linux"
  - "Backup"
date: "2007-01-28T21:31:00+02:00"
lastmod: "2007-01-28T21:31:00+02:00"
tags:
  - "backup"
  - "rdiff-backup"
  - "ssh"
  - "incremental backup"
toc: true
---

rdiff-backup is a simple and efficient utility. It can be used to make a remote backup of a directory via SSH. The differences are then saved in archives. [rdiff-backup Main page](https://www.nongnu.org/rdiff-backup/)

To install it:

```bash
emerge rdiff-backup
```

To make a backup on a remote computer, rdiff-backup must be installed on both machines.

For example, if the remote machine is called `remotehost.remotedomain`, to make a backup you simply do:

```bash
rdiff-backup ~/mydir remoteuser@remotehost.remotedomain::mydir-backup
```

This creates a new directory named `mydir-backup` in the `HOME` directory of the user `remoteuser` on `remotehost.remotedomain`. If the directory already exists, it is updated according to the content of `mydir`, and the differences are also stored.

For more complete examples: [Examples on the rdiff-backup site](https://www.nongnu.org/rdiff-backup/examples.html).

To make authentication automatic, it is recommended to use key-based authentication with SSH.
