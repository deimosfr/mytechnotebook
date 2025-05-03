---
weight: 999
url: "/FAQ_OpenSSH/"
title: "OpenSSH FAQ"
description: "Frequently Asked Questions about OpenSSH. Solutions to common problems and configuration tips."
categories: ["Network", "Security", "SSH"]
date: "2013-05-08T19:19:00+02:00"
lastmod: "2013-05-08T19:19:00+02:00"
tags: ["openssh", "ssh", "security", "authentication"]
toc: true
---

## Introduction

OpenSSH is not always simple, which is why a small documentation is useful.

## FAQ

### fatal: Timeout before authentication for @ip

Your DNS on your SSH server might not be up to date. Check them.

### Some clients take a long time to connect

On the SSH server, it is very likely that the server is trying to resolve names, which is not always possible or practical. The solution is to disable this (`/etc/ssh/sshd_config`):

```bash
...
UseDNS no
...
```

You just need to restart the SSH server.

### Unspecified GSS failure. Minor code may provide more information

Add this line to your configuration (`/etc/ssh/sshd_config`):

```bash
GSSAPIAuthentication no
```

### Limit users authorized in SSH

```bash
[...]
AllowUsers deimos
[...]
```
