---
weight: 999
url: "/Watchdog_\\:_détection_de_problèmes_hardware/"
title: "Watchdog: Hardware Problem Detection"
description: "This article explains how to use the Linux watchdog system to monitor hardware issues and automatically reboot the system when problems are detected."
categories: ["Linux"]
date: "2006-12-17T22:29:00+02:00"
lastmod: "2006-12-17T22:29:00+02:00"
tags: ["Linux", "Servers", "Hardware", "System Monitoring", "System Administration"]
toc: true
---

## Introduction

In computer hardware, a watchdog is an electronic or software mechanism designed to ensure that an automated system hasn't become stuck at a particular processing step. It's a protection mechanism designed to restart the system if a defined action is not executed within a given time period.

When implemented in software, it typically consists of a counter that is regularly reset to zero. If the counter exceeds a given value (timeout), then a system reset is triggered. The watchdog often consists of a register that is updated via a regular interrupt. It can also consist of an interrupt routine that must perform certain maintenance tasks before returning control to the main program. If a routine enters an infinite loop, the watchdog counter will no longer be reset to zero, and a reset is ordered. The watchdog also allows a restart if no instruction is provided for this purpose. You simply need to write a value exceeding the counter's capacity directly into the register. The watchdog will then initiate the reset.

In industrial computing, the watchdog is often implemented as an electronic device, generally a monostable flip-flop. It is based on the principle that each processing step must execute within a maximum time. It is therefore possible to arm a timer before its execution. When the flip-flop returns to its stable state before the task is complete, the watchdog is triggered. It implements a backup system that can either trigger an alarm, restart the device, or activate a redundant system.

Watchdogs are often integrated into microcontrollers and motherboards dedicated to real-time operations.

## Installation

Watchdog is simple to install and set up:

```bash
apt-get install watchdog
```

There's also a small kernel component to implement:

```bash
CONFIG_WATCHDOG=y
CONFIG_SOFT_WATCHDOG=y
```

## Configuration

Configuration is done in the `/etc/watchdog.conf` file. Let's look at some different mechanisms.

### Networks

Take here for example the IP addresses *192.168.0.138* and *192.168.0.1*. This means that we're going to continuously ping these IP addresses, and if one of them doesn't respond, it means we have a failure on our machine and therefore need to reboot. This method is quite dangerous in production, so be sure of what you're doing.

```bash
ping                   = 192.168.0.138
ping                   = 192.168.0.1
interface              = eth0
file                   = /var/log/messages
```

The pings are sent from the network card *eth0* and are logged in */var/log/messages*.

### System Load

If you believe your machine contains no bugs and that if the memory load is too high, there's a problem and you need to reboot, then here's an option that will appeal to you:

```bash
max-load-1             = 24
max-load-5             = 18
max-load-15            = 12
```

Modify the values according to your needs.

### Temperature

If you monitor your machine and want to restart in case of overheating, use this:

```bash
temperature-device     = /dev/hda
max-temperature        = 50
```

You should normally have already configured sensors beforehand (e.g., hdparm & lm-sensors).

### Default Options

You must also set the default options and adapt them to your needs:

```bash
# Defaults compiled into the binary
admin                   = root
interval                = 10
logtick                = 1

# This greatly decreases the chance that watchdog won't be scheduled before
# your machine is really loaded
realtime                = yes
priority                = 1

# Check if syslogd is still running by enabling the following line
pidfile         = /var/run/syslogd.pid
```

Once finished, apply the changes by restarting the service:

```bash
/etc/init.d/watchdog restart
```
