---
weight: 999
url: "/OpenSSH\\:_Multiplexage_des_connexions_SSH/"
title: "OpenSSH: SSH Connection Multiplexing"
description: "How to set up SSH connection multiplexing to accelerate login times for multiple connections to the same host."
categories: 
  - "Linux"
  - "Network"
date: "2007-05-14T19:37:00+02:00"
lastmod: "2007-05-14T19:37:00+02:00"
tags:
  - "ssh"
  - "openssh"
  - "network"
  - "security"
  - "performance"
toc: true
---

Since version 4.0, OpenSSH allows multiplexing several connections into one, which speeds up the connection time for subsequent logins.

This tip requires OpenSSH version 4.2 or higher to work.

Just add this to your `~/.ssh/config` file:

```bash
Host *
ControlMaster auto
ControlPath ~/.ssh/master-%r@%h:%p
```

All new connections to a host where you are already connected will go through this existing connection. In addition to speeding up connection time, this has the advantage of not prompting for passwords on subsequent connections.
