---
weight: 999
url: "/Numlockx_\\:_Activer_le_pavé_numérique_au_boot/"
title: "Numlockx: Enabling the numeric keypad at boot"
description: "How to automatically enable the numeric keypad at boot time using numlockx on Linux systems."
categories: ["Linux"]
date: "2008-07-09T06:40:00+02:00"
lastmod: "2008-07-09T06:40:00+02:00"
tags: ["Linux", "Configuration", "Boot", "Desktop"]
toc: true
---

## Introduction

The numeric keypad is not necessarily activated at boot time, which can quickly become annoying. Here's a method to enable it automatically on Gnome.

## Installation

```bash
apt-get install numlockx
```

## Configuration

We'll use numlockx with GDM. Insert this line into `/etc/gdm/Init/Default`:

```bash
test -x /usr/bin/numlockx && /usr/bin/numlockx on
```

Now you just need to restart your GDM and it works! :-)
