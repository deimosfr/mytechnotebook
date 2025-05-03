---
weight: 999
url: "/Mise_en_place_dun_Watchdog/"
title: "Setting up a Watchdog"
description: "Instructions for setting up and configuring a Watchdog on a system to ensure automatic restart when the system becomes unresponsive."
categories: ["Linux"]
date: "2008-12-29T11:31:00+02:00"
lastmod: "2008-12-29T11:31:00+02:00"
tags: ["Servers", "Linux"]
toc: true
---

## Introduction

A [watchdog](https://fr.wikipedia.org/wiki/Watchdog) is an electronic circuit or software used in digital electronics to ensure that an automated system or computer does not remain stuck at a particular step of its processing. It's a protection measure generally intended to restart the system if a defined action is not executed within a specified time.

In industrial computing, the watchdog is often implemented through an electronic device, generally a monostable flip-flop. It works on the principle that each step of the processing must execute within a maximum time. At each step, the system arms a timer before execution. If the flip-flop returns to its stable state before the task is completed, the watchdog triggers. It implements a backup system that can either trigger an alarm, restart the automaton, or activate a redundant system. Watchdogs are often integrated into microcontrollers and motherboards dedicated to real-time processing.

When implemented in software, it generally consists of a counter that is regularly reset to zero. If the counter exceeds a given value (timeout), the system performs a reset (restart). The watchdog often consists of a register that is updated via a regular interrupt. It can also consist of an interrupt routine that must perform certain maintenance tasks before returning control to the main program. If a routine enters an infinite loop, the watchdog counter will no longer be reset to zero and a reset is ordered. The watchdog also allows for a restart if no instruction is provided for this purpose. It is then sufficient to write a value exceeding the counter's capacity directly into the register: the watchdog will initiate the reset.

## Configuration

Configuration of the watchdog:

- Edit the `/etc/rc.conf` file and modify the following line:

```bash
watchdogd_flags=NO # for normal use: ""
```

to

```bash
watchdogd_flags="" # for normal use: ""
```

- Reboot, a watchdogd process should be present, check the parameters:

```bash
sysctl -a
```

- Modify the sysctl parameters:

```bash {linenos=table}
kern.watchdog.period=30
kern.watchdog.auto=1
```
