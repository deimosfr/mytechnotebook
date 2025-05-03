---
weight: 999
url: "/Se_connecter_par_port_serie_sur_sa_debian/"
title: "Connect to Debian via Serial Port"
description: "Guide to set up and use a serial port connection to a Linux machine, including minicom configuration and Grub setup."
categories: ["Debian", "Linux"]
date: "2011-05-20T13:41:00+02:00"
lastmod: "2011-05-20T13:41:00+02:00"
tags: ["Serial port", "Minicom", "Console", "Terminal", "Debian", "Grub"]
toc: true
---

## Introduction

This documentation helps you set up a connection to a Linux machine via a serial port. It also explains how to configure minicom.

## Installation

On Debian:

```bash
apt-get install minicom
```

## Configuration

### Linux

First, we need to determine your COM port configuration:

```bash
$ dmesg | grep tty
[    0.004000] console [tty0] enabled
[    1.629013] tty ttye7: hash matches
[    1.629021] tty ttyba: hash matches
[   12.843809] usb 5-2: pl2303 converter now attached to ttyUSB0
[  129.343389] type=1503 audit(1230854350.490:5): operation="inode_permission" requested_mask="w::" denied_mask="w::" fsuid=0 name="/dev/ttyUSB0" pid=6609 profile="/usr/sbin/cupsd"
```

We can see that ttyUSB0 is detected here.

### Mac

On Mac, remember to install your drivers if you don't find `/dev/cu.usbserial`. As you may have guessed, this is the device to use.

Launch minicom with your regular user:

```bash
minicom -s
```

To launch the help, press Ctrl+A and Z. In our case, we want to configure the communication device and speed. Press the "O" key and follow the instructions.

## Setting up the COM port in Grub

[Documentation on setting up the serial port to connect to Debian](/pdf/linuxserialconsole.pdf)
