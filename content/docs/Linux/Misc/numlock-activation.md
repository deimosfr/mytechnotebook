---
weight: 999
url: "/activer-le-pave-numerique/"
title: "Activating the Numeric Keypad"
description: "How to enable the numeric keypad by default in console mode and in X Window System."
categories: ["Linux", "System Configuration"]
date: "2008-01-22T08:49:00+02:00"
lastmod: "2008-01-22T08:49:00+02:00"
tags: ["Linux", "NumLock", "Console", "X Window"]
toc: true
---

To have the numlock activated in console and under X:

* In console mode:

```bash
echo "LEDS=+num" >> /etc/console-tools/config
```

* In X Window System:

```bash
apt-get install numlockx
```
