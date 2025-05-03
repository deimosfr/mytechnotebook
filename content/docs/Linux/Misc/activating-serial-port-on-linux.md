---
weight: 999
url: "/Activer_le_port_série_sur_Linux/"
title: "Activating the Serial Port on Linux"
description: "Learn how to configure and activate the serial port on Linux systems for remote access and management."
categories: ["Linux", "System Administration"]
date: "2012-05-14T17:09:00+02:00"
lastmod: "2012-05-14T17:09:00+02:00"
tags: ["Linux", "Serial Port", "Console", "Debian", "Grub", "Configuration"]
toc: true
---

## Introduction

It's often convenient to access a machine via the serial port. Unfortunately, we usually realize this too late. Here's how to activate it.

For usage information, please refer to the documentation on Minicom.

## Configuration

### Debian

Modify the following lines to ensure that the console part loads correctly during the Grub boot:

```bash {linenos=table,hl_lines=[8,16,17],anchorlinenos=true}
# If you change this file, run 'update-grub' afterwards to update
# /boot/grub/grub.cfg.

GRUB_DEFAULT=0
GRUB_TIMEOUT=5
GRUB_DISTRIBUTOR=`lsb_release -i -s 2> /dev/null || echo Debian`
GRUB_CMDLINE_LINUX_DEFAULT="quiet"
GRUB_CMDLINE_LINUX="console=tty0 console=ttyS0,9600n8"

# Uncomment to enable BadRAM filtering, modify to suit your needs
# This works with Linux (no patch required) and with any kernel that obtains
# the memory map information from GRUB (GNU Mach, kernel of FreeBSD ...)
#GRUB_BADRAM="0x01234567,0xfefefefe,0x89abcdef,0xefefefef"

# Uncomment to disable graphical terminal (grub-pc only)
GRUB_TERMINAL=console
GRUB_SERIAL_COMMAND="serial --speed=9600 --unit=0 --word=8 --parity=no --stop=1"

# The resolution used on graphical terminal
# note that you can use only modes which your graphic card supports via VBE
# you can see them in real GRUB with the command `vbeinfo'
#GRUB_GFXMODE=640x480

# Uncomment if you don't want GRUB to pass "root=UUID=xxx" parameter to Linux
#GRUB_DISABLE_LINUX_UUID=true

# Uncomment to disable generation of recovery mode menu entries
#GRUB_DISABLE_LINUX_RECOVERY="true"

# Uncomment to get a beep at grub start
#GRUB_INIT_TUNE="480 440 1"
```

Then update Grub:

```bash
update-grub
```

Uncomment the following line (`/etc/inittab`):

```bash {linenos=table,hl_lines=[4]}
[...]
# Example how to put a getty on a serial line (for a terminal)
#
T0:23:respawn:/sbin/getty -L ttyS0 9600 vt100
#T1:23:respawn:/sbin/getty -L ttyS1 9600 vt100
[...]
```

Verify that ttyS0 is also present (`/etc/securetty`):

```bash {linenos=table,hl_lines=[3]}
[...]
# UART serial ports
ttyS0
[...]
```

Reboot your machine! Your system will then be accessible via the serial port.

## References

- [https://www.cyberciti.biz/faq/howto-setup-serial-console-on-debian-linux/](https://www.cyberciti.biz/faq/howto-setup-serial-console-on-debian-linux/)
- [https://codepoets.co.uk/2011/getting-a-kvm-serial-console-with-grub2/](https://codepoets.co.uk/2011/getting-a-kvm-serial-console-with-grub2/)
