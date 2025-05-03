---
weight: 999
url: "/Modifier_les_paramètres_de_boot/"
title: "Modifying Boot Parameters"
description: "How to modify boot parameters in BSD systems by configuring boot.conf."
categories: ["Linux"]
date: "2009-05-17T09:25:00+02:00"
lastmod: "2009-05-17T09:25:00+02:00"
tags: ["BSD", "boot", "configuration"]
toc: true
---

## Introduction

You may want to modify certain boot parameters (Grub equivalent). Here is some information about that.

## boot.conf

Edit the file `/etc/boot.conf`:

```bash
stty com0 19200
set tty com0
set timeout 1
```

Personally, I just added the "set timeout" line to 1, the other parameters were already there. For those who don't understand these lines, I'll explain a bit:

- `stty com0 19200`: Sets the baud rate of the com0 port to 19200
- `set tty com0`: Defines the tty on com0 to initialize it
- `set timeout 1`: Sets the waiting time for the boot system (in seconds, here 1 sec)
