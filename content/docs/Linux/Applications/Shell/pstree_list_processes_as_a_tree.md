---
weight: 999
url: "/Pstree_\\:_lister_ses_process_sous_forme_d'arbre/"
title: "Pstree: List Processes as a Tree"
description: "Learn how to use the pstree command to visualize running processes in a tree format, showing parent-child relationships between processes."
categories: ["Linux"]
date: "2008-08-31T08:55:00+02:00"
lastmod: "2008-08-31T08:55:00+02:00"
tags:
  ["Command Line", "Process Management", "System Administration", "Linux Tools"]
toc: true
---

## Introduction

If you're using a system which has a lot of users, and you'd like to see who has started a particular script, daemon, or binary, then the pstree utility is very helpful. It draws a tree of all currently running processes - allowing you to see which processes are related.

In much the same way as the tree command isn't likely to be generally useful, this command might seem a little pointless if you're on a single-user machine, and you essentially start everything yourself. But even so it can be helpful to see where processes have come from.

## Examples

The most basic usage would look something like this:

```bash
$  pstree
init-+-apache2---10*[apache2]
     |-atd
     |-clamd
     |-cron
     |-events/0
     |-exim4---exim4
     |-freshclam
     |-getty
     |-gpg-agent
     |-khelper
     |-ksoftirqd/0
     |-kthread-+-aio/0
     |         |-ata/0
     |         |-ata_aux
     |         |-kblockd/0
     |         |-khubd
     |         |-kjournald
     |         |-kmirrord
     |         |-kseriod
     |         |-kswapd0
     |         |-2*[pdflush]
     |         |-xenbus
     |         `-xenwatch
     |-memcached
     |-migration/0
     |-monit---{monit}
     |-munin-node
     |-mysqld_safe-+-logger
     |             `-mysqld---16*[{mysqld}]
     |-pdnsd---4*[{pdnsd}]
     |-python
     |-qpsmtpd-forkser
     |-roundup-server---roundup-server
     |-screen---bash---irssi
     |-ssh-agent
     |-sshd-+-sshd---sshd---bash
     |      `-sshd---sshd---bash---pstree
     |-syslog-ng
     `-watchdog/0
```

Here we can see several kernel processes running, (aio, ata, kseriod, etc.), several system daemons (syslog-ng, qpsmtpd, etc), as well as the ssh processes open for my current user.

There aren't many ways of customizing the output of the display, although you can modify several things such as the display of command line arguments with:

```bash
pstree -a
...
...
  |-rinetd
  |-rpc.statd
  |-screen -S foo
  |   `-bash
  |-screen -S bar
  |   `-bash
  |-ssh-agent
  |-sshd
  |-syslog-ng -p /var/run/syslog-ng.pid
  |-udevd --daemon
  |-vino-session --sm-client-id default5
  |-(watchdog/0)
  |-(watchdog/1)
  |-wnck-applet--oaf-activate-iid=OAFIID:GNOME_Wncklet_Factory
  |-xenconsoled
  |   `-{xenconsoled}
  `-xenstored --pid-file /var/run/xenstore.pid
```

For all the available options please read the manpage via "man pstree".

If you don't have the pstree command available then you may find it in [the psmisc package](https://packages.debian.org/psmisc), and it may be installed by apt-get, or aptitude in the usual manner.

## References

http://www.debian-administration.org/articles/607
