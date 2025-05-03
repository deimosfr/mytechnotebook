---
weight: 999
url: "/OpenElec_\\:_Solution_multimedia_pour_Raspberry_Pi/"
title: "OpenElec: Multimedia Solution for Raspberry Pi"
description: "Guide on installing and configuring OpenElec multimedia center on Raspberry Pi with support for remote controls and hardware video decoding."
categories: ["Linux", "Storage"]
date: "2013-02-12T07:46:00+02:00"
lastmod: "2013-02-12T07:46:00+02:00"
tags: ["OpenElec", "Raspberry Pi", "XBMC", "Multimedia"]
toc: true
---

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 3.0 (2.99.2) |
| **Operating System** | Linux |
| **Website** | [OpenElec Website](https://openelec.tv/) |
| **Last Update** | 12/02/2013 |
| **Others** | Raspberry Pi B |
{{< /table >}}

## Introduction

OpenELEC is an embedded GNU/Linux distribution aimed at allowing the use of a media center (HTPC — Home Theatre PC) in the same way as any other device connected to your television, such as a DVD player or external digital TV receiver. Turn on your device and OpenELEC is ready to use in less than 10 seconds — as fast as some DVD players. All you need is a simple remote control to access its functions.

The Raspberry Pi is a single-board computer with an ARM processor designed by game inventor David Braben as part of his Raspberry Pi Foundation.
The computer is the size of a credit card, it allows the execution of several variants of the free GNU/Linux operating system and compatible software. It is supplied bare (motherboard only, without case, power supply, keyboard, mouse or screen) with the aim of reducing costs and allowing the use of recovered hardware.
This computer is intended to encourage learning computer programming. However, it is sufficiently open (USB and network ports) and powerful (ARM 700 MHz, 256 MB of RAM for the original model, 512 MB on the latest versions) to allow a wide range of uses. Its BMC Videocore 4 graphics circuit in particular can decode full HD Blu-ray streams (1080p 30 frames per second), emulate old consoles, and run relatively recent video games.

We will see here how to set up OpenElec on Raspberry Pi and add some features.

## Installation

To install OpenElec on Raspberry Pi, you first need to download [the latest version](https://sources.openelec.tv/tmp/image/), then decompress it:

```bash
wget http://openelec.tv/get-openelec/download/finish/10-raspberry-pi-builds/30-openelec-testing-raspberry-pi-arm
tar -xf 30-openelec-testing-raspberry-pi-arm
```

Now, insert your SD card into your computer's reader, check the name of the associated device:

```bash
dmesg
```

Then build the system using the provided script:

```bash
cd OpenELEC-RPi.arm-2.99.2
./create_sdcard /dev/sdb
```

Now you can install the SD card in the Raspberry Pi, boot it and connect to it via SSH:

```
root/openelec
```

## Configuration

For the configuration part, there are several things that need to be done. First, you should know that if you look at the displayed total memory, it's 380 MB by default because the video card shares its memory with RAM.

### Logitech Harmony One

At the time of writing, the remote control works perfectly with an [IR605Q Dongle](https://www.mediahd.fr/recepteur-infrarouge/14-recepteur-infrarouge-ir605q.html) and a stable version of OpenElec. With development versions, this is less reliable, so pay attention to the version you choose!

For setup, using the "Logitech Harmony Remote Software" tool, configure a new "Device" as follows:

- Device: Computer -> Media PC
  - Manufacturer: Microsoft
    - Model: MCE-1039

Then access your OpenElec via SSH:

- Login: root
- Password: openelec

And create this file:

```xml
<keymap>

<!-- ************* GLOBAL ***************************** -->
<!-- Define R as reload skin, S as screenshot, anywhere -->
<!-- and define some special remote keys -->
<!-- Red toggles fullscreen view -->
<!-- Green Toggles the watched status of an item -->
<!-- And yellow will toggle the library views between 'show everything' and 'show watched only' -->

<global>
    <keyboard>
       <r>XBMC.ReloadSkin()</r>
       <s>Screenshot</s>
    </keyboard>
    <remote>
 <Red>FullScreen</Red>
 <Green>ToggleWatched</Green>
 <Yellow>SendClick(25,14)</Yellow>
    </remote>
</global>

<!-- ************************************************ -->

<!-- On the home screen, 1 cleans the library, 2 triggers an update -->

<home>
     <remote>
       <one>XBMC.CleanLibrary(video)</one>
       <two>XBMC.UpdateLibrary(video)</two>
     </remote>
</home>

<!-- Set up Audio Delay Easy Keys -->

<FullscreenVideo>
    <remote>
 <Yellow>AudioDelayMinus</Yellow>
 <Blue>AudioDelayPlus</Blue>
    </remote>
</FullscreenVideo>

<!-- Set up zooming in picture slideshows -->

<Slideshow>
    <remote>
 <Yellow>ZoomOut</Yellow>
 <Blue>ZoomIn</Blue>
    </remote>
</Slideshow>

<!-- make the info window close when info is pressed again -->

  <movieinformation>
    <remote>
      <info>Close</info>
    </remote>
  </movieinformation>
  <musicinformation>
    <remote>
      <info>Close</info>
    </remote>
  </musicinformation>

</keymap>
```

Then restart your OpenElec and the remote control will work :-)

### Hardware MPEG-2/VC-1 video decoding

It's more efficient to decode MPEG-2 and VC-1 format videos using hardware rather than software since there is a dedicated chip for this in the Raspberry Pi. The problem is that you need to purchase licenses to be allowed to do this. Fortunately, they are not expensive. Go to [the Raspberry Pi website](https://www.raspberrypi.com/) and purchase them.

You will be asked for your Raspberry Pi's serial number. To retrieve it:

```bash {linenos=table,hl_lines=[13]}
> cat /proc/cpuinfo
Processor	: ARMv6-compatible processor rev 7 (v6l)
BogoMIPS	: 697.95
Features	: swp half thumb fastmult vfp edsp java tls
CPU implementer	: 0x41
CPU architecture: 7
CPU variant	: 0x0
CPU part	: 0xb76
CPU revision	: 7

Hardware	: BCM2708
Revision	: 000e
Serial		: 0000000000000000
```

Then you just need to insert the codes in this file (`/flash/config.txt`):

```bash {linenos=table,hl_lines=[7,8]}
[...]
################################################################################
# License keys to enable GPU hardware decoding for various codecs
# to obtain keys visit the shop at http://www.raspberrypi.com
################################################################################

decode_MPG2=0x00000000
decode_WVC1=0x00000000
[...]
```

Restart and you're done.

## Backing Up Your Configuration

You may want to backup your entire configuration. This is very simple, everything is in file format, so you just have to copy it wherever you want. The configuration folder is located at:

```
/storage/.xbmc
```

When you want to restore, you just need to copy the same folder to the same location.

## References

1. http://linuxfr.org/news/openelec-2-0-annonce
2. http://fr.wikipedia.org/wiki/Raspberry_Pi
