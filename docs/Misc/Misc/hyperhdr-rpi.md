---
title: 'HyperHDR with Raspberry Pi Zero 2W'
slug: hyperhdr/
description: 'Get your TV true colors with Raspberry Pi Zero 2W and HyperHDR'
categories: ['Hyperhdr', 'Raspberry Pi']
tags: ['Hyperhdr', 'Raspberry Pi']
toc: true
date: '2025-05-04T15:08:37+02:00'
---

I really love the [HyperHDR](https://github.com/awawa-dev/HyperHDR) project. I use it to control my LED strips on my Sony TV. I've even [designed a 3D case](https://www.printables.com/model/1072533-ambilight-case-for-hyperhdrhyperion) to put my Raspberry Pi Zero 2W and all the requirements in one place.

## Prerequisites

### Operating system

I use [Raspberry Pi OS](https://www.raspberrypi.org/software/operating-systems/). The latest version is [Raspberry Pi OS (64-bit) Lite](https://www.raspberrypi.org/software/operating-systems/).

### Kernel fix

We'll start to install a [Kernel fix](https://github.com/awawa-dev/P010_for_V4L2) with DKMS. Run this as root:

```bash
wget https://raw.githubusercontent.com/awawa-dev/P010_for_V4L2/refs/heads/master/dkms-installer.sh
chmod +x ./dkms-installer.sh
./dkms-installer.sh
```

Then answer `1` for RPi:

```bash
 Please enter any one of the following values and press Enter
1 for Raspberry Pi OS
2 for other Debian/Ubuntu x64 system
Enter value:
1
```

Then reboot.

## Installation

Create a hyperhdr user and connect with it:

```bash
adduser hyperhdr
su hyperhdr
```

HyperHDR install is [straightforward](https://awawa-dev.github.io/wiki/Installation#debian-and-ubuntu-apt-repository):

```bash
sudo apt-get update
sudo apt-get remove hyperhdr
sudo apt-get install -y curl
curl -fsSL https://awawa-dev.github.io/hyperhdr.public.apt.gpg.key | sudo dd of=/usr/share/keyrings/hyperhdr.public.apt.gpg.key
sudo chmod go+r /usr/share/keyrings/hyperhdr.public.apt.gpg.key
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/hyperhdr.public.apt.gpg.key] https://awawa-dev.github.io $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hyperhdr.list > /dev/null
sudo apt-get update
sudo apt-get install hyperhdr -y
```

At the end, you should see something like this:

```bash
 +-----------------------------------------------------------------------+
 |                 HyperHDR has been installed/updated!                  |
 +-----------------------------------------------------------------------+
 |  For configuration, visit with your browser: 192.168.0.2:8090         |
 |  If already used by another service try: 192.168.0.2:8091             |
 |  Start the service: sudo systemctl start hyperhdr@hyperhdr            |
 |  Stop the service: sudo systemctl stop hyperhdr@hyperhdr              |
 |  Troubleshooting? Run HyperHDR manually: /usr/bin/hyperhdr            |
 +-----------------------------------------------------------------------+
 |  HyperHDR is installed as a service and starts automatically          |
 +-----------------------------------------------------------------------+
 |  Webpage: https://hyperhdr.eu                                         |
 |  GitHub: https://github.com/awawa-dev/HyperHDR                        |
 +-----------------------------------------------------------------------+
```

You can now access the web interface at `http://<your-rpi-ip>:8090`.

## OS Configuration

### Video device permissions

Now allow hyperhdr user to access video devices:

```bash
sudo usermod -a -G video hyperhdr
```

### Enable service

Enable the service to enable reboot persistence:

```bash 
sudo systemctl enable hyperhdr@hyperhdr
```

### Read only file system

!!! info
    It will protect your Raspberry Pi system from corruption if you don't have a dedicated & safe powering off mechanism for the Rpi. And will extend the life of your SD card since all writes will be stored in the memory.

Once you've a working HyperHDR and won't update it for some time, you can enable the read only file system:

* Run `sudo raspi-config`
* Select `Performance Options`
* Select `P2 Overlay File System Enable/disable read-only file system`
* Select `Yes`

Then reboot.

## HyperHDR configuration

### LUT

First we have to download LUT files based on your hardware config:

![LUT download](../../static/images/hyperhdr_lut_download.avif)

### Video capturing

To configure video capturing, leaving everything to automatic should be enough. I just advise you to check the `Auto resume` box:

![Video capturing](../../static/images/hyperhdr_video_capturing.avif)

### LED strip configuration

#### Led hardware

Using an ESP32 with [WLED](https://kno.wled.ge/) is a very good thing to manage effitently your LED strip with and without HyperHDR.

So you have to configure the endpoint:

![LED hardware](../../static/images/hyperhdr_led_hardware.avif)

!!! note
    You have to properly configure WLED first, ensuring the voltage, number of LEDs, and the correct protocol are set. Don't go to the LED layout configuration without having a working WLED.

#### LED layout

Finally you have to configure the LED layout to match your LED strip configuration:

![LED layout](../../static/images/hyperhdr_led_layout.avif)

### LED visualisation (check)

You can check the capture and LEDs with the LED visualisation:

![LED visualisation](../../static/images/hyperhdr_led_viz.avif)

## Troubleshooting

### Test your video capturing colors

You can test your video capturing colors by running video like [color wheel on YouTube](https://www.youtube.com/watch?v=xAwB9lQnxAY).

### Swapped colors

I encountered a problem with my Raspberry Pi Zero 2W where the colors are swapped. It looks to be a [kernel issue](https://github.com/awawa-dev/HyperHDR/discussions/848) and the only solution I've found is to revert to a working kernel:

```bash
sudo rpi-update d86b5843d68b9972a5430a6d3da1b271cfc83521
```

Then reboot.

!!! note
    I've tried many times tweaking with v4l2-ctl but I couldn't get something as good as downgrading the kernel.
