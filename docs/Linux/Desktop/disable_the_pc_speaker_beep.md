---
title: "Disable the PC Speaker Beep"
slug: disable-the-pc-speaker-beep/
description: "How to disable the annoying PC speaker beep in Linux systems both temporarily and permanently."
categories: ["Linux"]
date: "2009-04-26T14:52:00+02:00"
lastmod: "2009-04-26T14:52:00+02:00"
tags: ["Linux", "Configuration", "System"]
---

## Introduction

It's very annoying to hear BEEPs for every impossible completions or errors when you're working in a shell!

## Solutions

### Temporary Solution

Simply, from a graphical interface, run this command:

```bash
xset -b
```

That's the solution to stop the biiiiiiiiiiiiiip :-)

### Permanent Solution

If you want a long-term solution! Add this line to the `/etc/modprobe.d/blacklist` file:

```bash
blacklist pcspkr
```

This will disable all PC speaker sounds (startup, shutdown...). Pure bliss!
