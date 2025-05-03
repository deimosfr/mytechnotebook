---
weight: 999
url: "/redemarrer_x/"
title: "Restart X"
description: "How to restart the X Window System using keyboard shortcuts and command line methods."
categories: ["Linux", "Ubuntu"]
date: "2009-04-15T21:00:00+02:00"
lastmod: "2009-04-15T21:00:00+02:00"
tags: ["X11", "X Server", "Shortcuts", "Ubuntu", "System Management"]
toc: true
---

## Introduction

To restart X, there aren't 36 solutions, but there is one in particular that I really like.

## Solution

It's a keyboard shortcut:

```
Ctrl + Alt + Backspace
```

This avoids having to type, for example:

```bash
/etc/init.d/gdm restart
```

That's it, X restarts :-). **Warning, this should only be done while in X, otherwise the PC will reboot :-(**.

### Problem starting from Ubuntu 9.04

Here's what you can read:

```
Ctrl-Alt-Backspace is now disabled, to reduce issues experienced by users who accidentally trigger the key combo. Users who do want this function can enable it in their xorg.conf, or via the command
```

How annoying! To reactivate it:

```bash
apt-get install dontzap
```

Then:

```bash
dontzap --disable
```
