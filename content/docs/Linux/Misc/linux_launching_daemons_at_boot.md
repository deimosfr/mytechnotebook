---
weight: 999
url: "/Linux_\\:_Lancement_des_démons_au_boot/"
title: "Linux: Launching Daemons at Boot"
description: "Guide explaining how to configure daemons to automatically start at boot time on different Linux distributions such as Debian, Red Hat, and Gentoo."
categories: 
  - Debian
  - Linux
  - Red Hat
date: "2009-05-05T07:11:00+02:00"
lastmod: "2009-05-05T07:11:00+02:00"
tags:
  - Linux
  - Servers
  - System Administration
  - Boot
toc: true
---

## Introduction

To launch daemons at boot time, you need to act differently depending on the distribution. But before starting, the script must already exist in `/etc/init.d`.

## Debian

Here's the command:

```bash
update-rc.d
```

You can test with the `-n` argument:

```bash
update-rc.d -n name_of_service defaults
```

Once you're ready, execute:

```bash
update-rc.d name_of_service defaults
```

And if you want to remove it:

```bash
update-rc.d -f name_of_service remove
```

## Red Hat

Run this command:

```bash
ntsysv
```

A menu will open. It's up to you to make your choice.

## Gentoo

For Gentoo, here's the procedure:

```bash
rc-update add name_of_service default
```
