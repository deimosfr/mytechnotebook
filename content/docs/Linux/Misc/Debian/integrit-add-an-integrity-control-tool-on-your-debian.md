---
weight: 999
url: "/Integrit_\\:_Add_an_integrity_control_tool_on_your_Debian/"
title: "Integrit: Add an integrity control tool on your Debian"
description: "How to install and configure Integrit, a simple yet secure alternative to tripwire for file integrity monitoring on Debian systems."
categories: ["Linux", "Security", "System Administration"]
date: "2013-05-07T14:44:00+02:00"
lastmod: "2013-05-07T14:44:00+02:00"
tags: ["Integrit", "Debian", "Security", "File Integrity", "Monitoring"]
toc: true
---

![Integrit](/images/poweredbylinux.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 4.1 |
| **Operating System** | Debian 7 |
| **Website** | [Integrit Website](https://sourceforge.net/projects/integrit/) |
| **Last Update** | 07/05/2013 |
{{< /table >}}

## Introduction

Integrit[^1] is a simple yet secure alternative to products like tripwire. It has a small memory footprint, uses up-to-date cryptographic algorithms, and has features that make sense (like including the MD5 checksum of newly generated databases in the report

## Installation

To install Integrit:

```bash
aptitude install integrit
```

## Configuration

```bash {linenos=table,hl_lines=[10,13]}
# /etc/integrit/integrit.debian.conf
# Configuration of the example daily cron job /etc/cron.daily/integrit

# Set the configuration file(s) for integrit.  /etc/cron.daily/integrit
# will run ``integrit -uc -C <file>'' for each file specified in CONFIGS.
# An empty CONFIGS variable disables /etc/cron.daily/integrit.  Multiple
# file names are separated with spaces, e.g.:
# CONFIGS="/etc/integrit/usr.conf /etc/integrit/lib.conf"
# CONFIGS="/etc/integrit/integrit.conf"
CONFIGS="/etc/integrit/integrit.conf"

# Set the mail address reports are sent to
EMAIL_RCPT="xxx@mycompany.com"

# Set the subject line for the report mails
EMAIL_SUBJ="[integrit] `hostname -f`: report on changes in the filesystems"

# If ALWAYS_EMAIL is set to ``true'', a report is mailed on every run.
# Normally a report is only generated when integrit(1) exits non-zero.
ALWAYS_EMAIL=false
```

You need to adapt the vars vars listed bellow:

- CONFIGS: set your main configuration or multiples if you have so
- EMAIL_RCPT: your email address (the recipient)
- EMAIL_SUBJ: the email subject if this one doesn't suit you
- ALWAYS_EMAIL: set it to false if you want to receive emails only when a change occur

Now we're going to edit the main configuration of Integrit:

```bash {linenos=table,hl_lines=["19-21","41-44"]}
# /etc/integrit/integrit.conf
# /etc/integrit.conf : configuration file for integrit
#
# See integrit(1) and /usr/share/doc/integrit/examples/
# for more information.
#
# *** WARNING ***
#
# This is a simple default configuration file for Debian systems.
# It contains only comments, therefore integrit will not run with
# it. To make integrit functional, you must edit this file according
# to your needs.
#
# Please read README.Debian before running integrit.
#
# *** WARNING ***

#
root=/
known=/var/lib/integrit/known.cdb
current=/var/lib/integrit/current.cdb
#
# # Here's a table of letters and the corresponding checks / options:
# # Uppercase turns the check off, lowercase turns it on.
# #
# # 	  s	checksum
# # 	  i	inode
# # 	  p	permissions
# # 	  l	number of links
# # 	  u	uid
# # 	  g	gid
# # 	  z	file size (redundant if checksums are on)
# # 	  a	access time
# # 	  m	modification time
# # 	  c	ctime (time UN*X file info last changed)
# # 	  r	reset access time (use with care)
#
# # ignore directories that are expected to change
#
# !/cdrom
!/dev
!/lost+found
!/proc
!sys
# !/etc
# !/floppy
# !/home
# !/mnt
# !/root
# !/tmp
# !/var
#
# # ignore inode, change time and modification time
# # for ephemeral module files.
#
# /lib/modules/2.4.3/modules.dep IMC
# /lib/modules/2.4.3/modules.generic_string IMC
# /lib/modules/2.4.3/modules.isapnpmap IMC
# /lib/modules/2.4.3/modules.parportmap IMC
# /lib/modules/2.4.3/modules.pcimap IMC
# /lib/modules/2.4.3/modules.usbmap IMC
#
# # to cut down on runtime and db size:
#
# =/usr/include
# =/usr/X11R6/include
#
# =/usr/doc
# =/usr/info
# =/usr/share
#
# =/usr/X11R6/man
# =/usr/X11R6/lib/X11/fonts
#
# # ignore user-dependant directories
#
# !/usr/local
# !/usr/src
```

To give you a quick understand of this configuration file:

- !: do not scan this folder/file
- =: do not search recursively if it's a folder
- $: tells not not inherit from the parent folder regarding the checking method
- /etc **MC**: this example ask to not check mtime + ctime verification on /etc

Now we're going to initialize the known database:

```bash
integrit -C /etc/integrit/integrit.conf -u
integrit: ---- integrit, version 4.1 -----------------
integrit:                      output : human-readable
integrit:                   conf file : /etc/integrit/integrit.conf
integrit:                    known db : /var/lib/integrit/known.cdb
integrit:                  current db : /var/lib/integrit/current.cdb
integrit:                        root : /
integrit:                    do check : no
integrit:                   do update : yes
```

Move the current known database to known:

```bash
mv /var/lib/integrit/current.cdb /var/lib/integrit/known.cdb
```

Next and to finish, you can update manually (or let cron do) the database:

```bash
integrit -C /etc/integrit/integrit.conf -c
```

{{< alert context="danger" text="This is strongly recommanded that you put the known database on a read only share" />}}

## References

[^1]: [https://sourceforge.net/projects/integrit/](https://sourceforge.net/projects/integrit/)
