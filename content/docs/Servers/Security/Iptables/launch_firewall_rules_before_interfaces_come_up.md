---
weight: 999
url: "/Lancer_les_règles_de_Firewalling_avant_que_les_interfaces_deviennent_up/"
title: "Launch Firewall Rules Before Interfaces Come Up"
description: "This guide explains how to configure your firewall rules to load before network interfaces come up, ensuring your system is always protected by a firewall."
categories: ["Networking", "Linux", "Debian"]
date: "2008-09-24T11:42:00+02:00"
lastmod: "2008-09-24T11:42:00+02:00"
tags: ["Firewall", "Linux", "Networking", "Security", "iptables"]
toc: true
---

## Introduction

There used to be a script to do it automatically via init.d files, but now the suggested method is to use ifup.d networking scripts, which are executed on state changes of the network interfaces. So I submit here my simple script, which does the trick for me nicely.

## Configuration

Drop this script into `/etc/network/if-pre-up.d` in a file called iptables:

```bash
#!/bin/sh

# Load iptables rules before interfaces are brought online
# This ensures that we are always protected by the firewall
#
# Note: if bad rules are inadvertently (or purposely) saved it could block
# access to the server except via the serial tty interface.
#

RESTORE=/sbin/iptables-restore
STAT=/usr/bin/stat
IPSTATE=/etc/iptables.conf

test -x $RESTORE || exit 0
test -x $STAT || exit 0

# Check permissions and ownership (rw------- for root)
if test `$STAT --format="%a"` $IPSTATE -ne "600"; then
  echo "Permissions for $IPSTATE must be 600 (rw-------)"
  exit 0
fi

# Since only the owner can read/write to the file, we can trust that it is
# secure. We need not worry about group permissions since they should be
# zeroed per our previous check; but we must make sure root owns it.
if test `$STAT --format="%u"` $IPSTATE -ne "0"; then
  echo "The superuser must have ownership for $IPSTATE (uid 0)"
  exit 0
fi

# Now we are ready to restore the tables
$RESTORE < $IPSTATE
```

Then make sure you make the script executable:

```bash
chmod +x iptables
chown root:root iptables
```

It loads the settings from $IPSTATE - by default, `/etc/iptables.conf`. You have to save the rules manually; this ensures that you make sure your rules are working properly (i.e. doesn't block you from logging in remotely, for example) before you decide to save them.

You do this running the command: "iptables-save > /etc/iptables.conf" (or whatever file you have chosen to use as your $IPSTATE file)

## Resources
- [https://www.debian-administration.org/articles/615](https://www.debian-administration.org/articles/615)
