---
weight: 998
url: "/Activer_le_port_série_sous_FreeBSD/"
title: "Activating the Serial Port on FreeBSD"
description: "Learn how to configure and activate the serial port on FreeBSD systems for remote access and console management."
categories: ["BSD", "System Administration"]
date: "2012-05-19T23:07:00+02:00"
lastmod: "2012-05-19T23:07:00+02:00"
tags: ["FreeBSD", "Serial Port", "Console", "Configuration", "Soekris"]
toc: true
---

## Introduction

By default, the serial port is not activated. If you have a Soekris device, for example, it's essential to activate it. Let's see how to do this.

## Configuration

To activate the standard output on the serial port (only stdout):

```bash
echo "-h" > /boot.config
```

Alternatively, you can choose to activate the serial port only if no keyboard is connected to the machine:

```bash
echo "-P" > /boot.config
```

Next, you'll need to configure the stdin part (keyboard) by enabling this line by changing it to "on":

```bash {linenos=false,hl_lines=[5]}
# /etc/ttys
[...]
# Serial terminals
# The 'dialup' keyword identifies dialin lines to login, fingerd etc.
ttyu0   "/usr/libexec/getty std.9600"   dialup  on secure
[...]
```

Reboot and then all you need to do is connect to it.

## References

- [https://www.freebsd.org/doc/en_US.ISO8859-1/books/handbook/serialconsole-setup.html](https://www.freebsd.org/doc/en_US.ISO8859-1/books/handbook/serialconsole-setup.html)
