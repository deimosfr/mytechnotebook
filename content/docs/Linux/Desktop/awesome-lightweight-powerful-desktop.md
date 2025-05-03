---
weight: 999
url: "/Awesome_\\:_un_bureau_léger_et_puissant/"
title: "Awesome: A Lightweight and Powerful Desktop"
description: "Learn how to install and configure Awesome, a lightweight and powerful tiling window manager for Linux."
categories: ["Linux", "Desktop", "Window Manager"]
date: "2013-03-28T08:21:00+02:00"
lastmod: "2013-03-28T08:21:00+02:00"
tags: ["Awesome", "Linux", "Window Manager", "Tiling", "Desktop Environment"]
toc: true
---

![Awesome](/images/awesome_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 3.4.11 |
| **Operating System** | Debian 7 |
| **Website** | [Website](https://awesome.naquadah.org) |
| **Last Update** | 28/03/2013 |
{{< /table >}}

## Introduction

Awesome is a free window manager that runs on top of the X Window system on UNIX-type machines. Its goal is to remain very lightweight and offer several window layouts: maximized, floating, but also automatically placed in the form of tiles (a mode called tiling), similar to Ion.

## Installation

```bash
aptitude install awesome
```

## Configuration

The configuration part is the most challenging aspect of Awesome. You can do really nice things provided you spend enough time on it.

### Checking your configuration

First, let's make sure we have our configuration file ~/.config/awesome/rc.lua. Every time you modify it, it's wise to test your configuration before restarting awesome:

```bash
> awesome -k ~/.config/awesome/rc.lua
? Configuration file syntax OK.
```

### Testing your configuration

It's very convenient to test your rc.lua configuration without having to break your own desktop[^1]. A solution exists, which consists of launching a virtual desktop in your desktop. For this, we'll use Xephyr:

```bash
aptitude install xserver-xephyr
```

Then, we'll initialize this virtual environment:

```bash
Xephyr :1 -ac -br -noreset -screen 1152x720 &
```

and run our new awesome configuration (rc.lua.new):

```bash
DISPLAY=:1.0 awesome -c ~/.config/awesome/rc.lua.new
```

## References

[^1]: https://wiki.archlinux.org/index.php/Awesome#Debugging_rc.lua
