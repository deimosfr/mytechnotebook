---
weight: 999
url: "/Ulimit_\\:_Utiliser_les_limites_systèmes/"
title: "Ulimit: Using System Limits"
description: "An overview of how to use ulimit to manage system-wide resource limits for users and processes on Linux and Solaris systems."
categories: ["Linux", "Security"]
date: "2012-07-03T10:08:00+02:00"
lastmod: "2012-07-03T10:08:00+02:00"
tags: ["Linux", "Security", "System Administration", "Solaris", "Resource Management"]
toc: true
---

## Introduction

The ulimit programs allow to limit system-wide resource use using a normal configuration file - `/etc/security/limits.conf`. This can help a lot in system administration, e.g. when a user starts too many processes and therefore makes the system unresponsive for other users.

## Usage

### Linux

```bash
ulimit -a
```

```
core file size          (blocks, -c) 0
data seg size           (kbytes, -d) unlimited
scheduling priority             (-e) 0
file size               (blocks, -f) unlimited
pending signals                 (-i) 7671
max locked memory       (kbytes, -l) 64
max memory size         (kbytes, -m) 811664
open files                      (-n) 1024
pipe size            (512 bytes, -p) 8
POSIX message queues     (bytes, -q) 819200
real-time priority              (-r) 0
stack size              (kbytes, -s) 8192
cpu time               (seconds, -t) unlimited
max user processes              (-u) 7671
virtual memory          (kbytes, -v) 1175120
file locks                      (-x) unlimited
```

All these settings can be manipulated. A good example is this forkbomb that forks as many processes as possible and can crash systems where no user limits are set

Now this is not good - any user with shell access to your box could take it down. But if that user can only start 20 processes the damage will be minimal. So let's set a process limit of MAX 20 process for a particular users in the system, this can be done by inserting the simple one line in limit.conf file.

Following will prevent a "fork bomb" (`/etc/security/limits.conf`):

```bash
deimos hard nproc 20
@group1 hard nproc 50
```

Above will prevent user "deimos" to create more than 20 process and anyone in the group1 from having more than 50 processes.

There are many more setting and limits that you can set on a particular user or to a entire group like..

Using below configuration will prevent any users in the system to logins not more than 3 places at same time (`/etc/security/limits.conf`):

```bash
hard maxlogins 3
```

Limit on size of core file (`/etc/security/limits.conf`):

```bash
hard core 0
```

### Solaris

To get all information:

```bash
ulimit -a
```

To display a process' current file descriptor limit, run:

```bash
/usr/proc/bin/pfiles pid
```

Remove the grep to see all files linked to a process.

To change the files descriptor for example:

```bash
ulimit -n 1024
```
