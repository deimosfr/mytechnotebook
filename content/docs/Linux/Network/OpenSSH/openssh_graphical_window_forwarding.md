---
weight: 999
url: "/OpenSSH_\\:_Export_de_fenêtre_graphiques/"
title: "OpenSSH: Graphical Window Forwarding"
description: "How to export X windows using OpenSSH tunneling for accessing graphical applications on remote machines securely"
categories: ["Linux", "Network"]
date: "2009-01-12T09:37:00+02:00"
lastmod: "2009-01-12T09:37:00+02:00"
tags: ["SSH", "X11", "Security", "Tunneling", "VNC"]
toc: true
---

## Introduction

OpenSSH is capable of exporting X windows from another machine (creating an SSH tunnel). For example, you can connect to a server that has X and you only have SSH access to the remote machine.

In your SSH configuration file (`/etc/ssh/sshd_config`), set this to yes:

```bash
X11Forwarding yes
```

## Launching the Session

Here's an example session. Here we'll export VNC which is running on the direct machine:

```bash
/usr/bin/ssh -gL5901:127.0.0.1:5901 -C xxx@mycompany.com
```

This will export the remote port 5901 to the local machine's port 5901.

If we need to go through an SSH gateway first:

```bash
/usr/bin/ssh gateway_machine -L 5901:machine_on_my_network:5901
```

It's also possible to automate an SSH tunnel by adding a line like this in the SSH configuration file (`~/.ssh/config`):

```bash
LocalForward <local port> <target machine>:<target port>
```

Here's an example:

```bash
Host mycompany.com
    User username                       // To use a different username than the current one
    LocalForward 993 localhost:993      // To access my own IMAPS server
    LocalForward 119 news.free.fr:119   // To access the free news server
```

## Connecting to the Remote Session

Launch vncviewer and connect to **"localhost:1"**. You will then see the remote server screen.

## Conclusion

SSH is capable of forwarding any window and any port. For security reasons, it's preferable to open as few ports as possible. Just open SSH to pass these types of services.

## Resources
- [Documentation on Best Practices on SSH](/pdf/ssh-_best_practices.pdf)
- [Principles and Usage of SSH](/pdf/principes_et_utilisation_de_ssh.pdf)
