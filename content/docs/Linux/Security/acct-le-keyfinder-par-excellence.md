---
weight: 999
url: "/Acct_Le_keyfinder_par_excellence/"
title: "Acct: The Ultimate Keyfinder"
description: "Learn how to track user commands and processes with the acct tools for system auditing and monitoring."
categories: ["Linux", "Security", "Administration"]
date: "2013-05-06T13:59:00+02:00"
lastmod: "2013-05-06T13:59:00+02:00"
tags:
  ["acct", "monitoring", "commands", "logging", "security", "auditing", "linux"]
toc: true
---

## Introduction

In a production environment, it can be useful to know what each person is doing. This is particularly helpful when a mistake happens and nobody admits to it (yes, it happens). Novice hackers (aka script kiddies) who call themselves hackers because they've put a keylogger on a machine might also be interested in this. However, the purpose is obviously not the same.

Two commands are useful:

- sa: obtains statistics on process launches
- lastcomm: obtains a list of commands launched by users

## Installation

The installation is done as follows:

```bash
apt-get install acct
```

## Configuration

- All log files will be written to this file:

```bash
/var/log/account/pacct
```

- If you want to change the file, execute this action:

```bash
accton FileName
```

- For activation, edit the file /etc/default/acct:

```bash
# Activate acct
ACCT_ENABLE="1"

# Amount of days that the logs are kept.
ACCT_LOGGING="30"
```

## Usage

### lastcomm

- To list the commands used:

```bash
lastcomm
```

{{< alert context="warning" text="Beware, you can also see what the shell executes on startup" />}}

- List commands recently launched by a user:

```bash
lastcomm user
```

- Search in history for who launched a given command and when:

```bash
lastcomm apachectl
```

- Find out which commands were launched directly from the physical terminal of the machine:

```bash
lastcomm --tty tty1
```

### sa

- List commands that ran the longest:

```bash
sa --sort-real-time | head
```

- List commands that consume the most I/O:

```bash
sa -d | head
```

- List all commands with the user who launched them:

```bash
sa -u
```

- Consumption by user:

```bash
sa -m
```

The output contains:

- Number of calls
- re: time spent
- cp: amount of CPU consumed (in seconds)
- avio: average number of I/O operations (very useful for diagnosing which process is using the disk)
- Memory consumed per second (k, this value is not very intuitive)

## References

[https://tldp.org/HOWTO/Process-Accounting/pa.html](https://tldp.org/HOWTO/Process-Accounting/pa.html)
