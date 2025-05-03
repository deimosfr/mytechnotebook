---
weight: 999
url: "/Modifier_le_firmware_par_le_FrancoFON/"
title: "Modify the firmware with FrancoFON"
description: "Guide to modify your Fonera firmware using FrancoFON to add enhanced functionality including improved firewall management, port redirection, and more advanced networking options."
categories: ["Linux", "Network"]
date: "2008-01-01T17:31:00+02:00"
lastmod: "2008-01-01T17:31:00+02:00"
tags: ["Firmware", "Fonera", "Network", "SSH", "Router", "Customization"]
toc: true
---

## Introduction

I'd like to modify the original firmware to have more functionality. That's why I opted for this firmware. Unfortunately, at the time of writing this article, it's not possible to run this firmware on a Fonera+. Here's what it adds compared to the base firmware:

- Improved firewall management
- Simplified port redirection
- DHCP address management for the private network
- Blacklisting of certain sites for the public network
- Adding/removing DNS
- Antenna power adjustment
- Remote reboot
- Management of dyndns accounts
- Hosts file management
- Live view of connections to private and public networks
- Comparison of the Fonera version with the latest available version
- Hidden SSID management for the private network
- Language management (English, French, Romanian for now)
- Ability to update firmware from a server other than the FrancoFON server
- Ability to configure the Fonera in pppoe with the @ and / characters required by some ISPs, up to 64 characters
- Opening/closing SSH
- Opening/closing console access via ethernet
- Script to know the number of Foneras using the FrancoFON firmware
- Ability to block clients on public and private signals by MAC address with scheduler
- Ability to connect the Fonera to a WiFi AP to use its internet connection (autonomous "Fonera mode", also called "Ponte2")

## Prerequisites

You must have SSH installed on your Fonera. Use this documentation:
[Enable SSH on your Fonera+]({{< ref "docs/Misc/Fonera/fonera-ssh.md" >}})

## Installation of the new firmware

Connect to your Fonera via WiFi, then via SSH:

```bash
ssh root@192.168.10.1
```

The default root password is "admin" or the Fonera management password you defined on the website.

Then, let's start the installation:

```bash
cd /tmp
wget http://download.francofon.fr/update.sh
chmod +x update.sh
sh ./update.sh
```

The Fonera will then update:

```bash
*****************************************************
*    Script D'installation du firmware FrancoFON    *
*****************************************************
 Auteurs : FrancoFON.FR
 Version : 1.2
 Contact : webmaster@francofon.fr
*****************************************************
Le fichier contenant le firmware FrancoFON va être téléchargé
Le fichier a correctement ete telecharge
La decompression va commencer. Le processus peut être long.
Desarchivage en cours......
```

and the Fonera will reboot.

## FAQ

### I have too many problems, I want to restore factory settings

- Unplug the network cable and power from the Fonera.
- Locate the reset button by turning the Fonera over.
- Press the "Reset" button and hold it.
- Reconnect the power without releasing the button.
- Hold the button for at least 45 seconds.
- Wait until the WLAN light turns on.
- Press the reset button again for 10 seconds then release.
- Wait until the WLAN light turns on.
- Check via the console that the Fonera has returned to its factory settings.

## Resources
- http://www.francofon.fr/modules/mediawiki/index.php/Le_firmware_FrancoFON
- http://www.francofon.fr/modules/mediawiki/index.php/La_Fonera/Reset
