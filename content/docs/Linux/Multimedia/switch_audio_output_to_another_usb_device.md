---
weight: 999
url: "/Swith_audio_output_to_another_USB_device/"
title: "Switch audio output to another USB device"
description: "Instructions for setting up USB audio devices on Linux, including both manual and automatic approaches with udev rules"
categories: ["Linux", "Debian"]
date: "2013-05-23T15:00:00+02:00"
lastmod: "2013-05-23T15:00:00+02:00"
tags: ["Linux", "Debian", "USB", "Audio", "udev"]
toc: true
---

![Linux](/images/poweredbylinux.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Operating System** | Debian 7 |
| **Last Update** | 23/05/2013 |
{{< /table >}}

## Introduction

I've recently bought an Enermax keyboard (Caesar) and have input and output audio directly on the keyboard! That's great and better when it works. Unfortunately, as I have a minimal OS, it doesn't work out of the box. So here is how I could achieve it.

I will present 2 ways to make it work and I personally use the manual way.

## Manual way

First of all, you need to see if you could find your USB Audio device (normally there is no reason why you shouldn't see it):

```bash
> lsusb
[...]
Bus 003 Device 004: ID 0d8c:0105 C-Media Electronics, Inc. CM108 Audio Controller
[...]
```

Then check this module is loaded:

```bash
> lsmod | grep snd_usb_audio
snd_usb_audio          84836  2
```

If it's not the case, add it to your module list and load it:

```bash
echo "snd_usb_audio" >> /etc/modules
modprobe snd_usb_audio
```

Now find the corresponding card ID (you could also have it in AlsaMixer):

```bash {linenos=table,hl_lines=[3]}
> aplay -l
[...]
card 1: Device [USB Multimedia Audio Device], périphérique 0: USB Audio [USB Audio]
  Subdevices: 0/1
  Subdevice #0: subdevice #0
[...]
```

My card number is 1! Nice, now let's add to our personal configuration (replace 1 at the end of the lines by your card number):

```bash
defaults.ctl.card 1
defaults.pcm.card 1
```

Now reboot to make it work!

## Automatic way

For the automatic rules, we're going to play with udev. Simply create a rule for it:

```bash
# Set USB device as default sound card when plugged in
KERNEL=="pcmC[D0-9cp]*", ACTION=="add", PROGRAM="/bin/sh -c 'K=%k; K=$${K#pcmC}; K=$${K%%D*}; echo defaults.ctl.card $$K > /etc/asound.conf; echo defaults.pcm.card $$K >>/etc/asound.conf'"

# Restore default sound card when USB device unplugged
KERNEL=="pcmC[D0-9cp]*", ACTION=="remove", PROGRAM="/bin/sh -c 'echo defaults.ctl.card 0 > /etc/asound.conf; echo defaults.pcm.card 0 >>/etc/asound.conf'"
```

Udev will do the work to make it automatically work :-)

## References

[https://archlinux.me/w0ng/2012/07/06/alsa-switch-audio-usb-headset/](https://archlinux.me/w0ng/2012/07/06/alsa-switch-audio-usb-headset/)
