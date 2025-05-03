---
weight: 999
url: "/Grimwepa\\_\\:_le_hack_wifi_facile/"
title: "Grimwepa: Easy WiFi Hacking"
description: "A guide to using Grimwepa for wireless network security testing in Ubuntu"
categories: ["Linux", "Ubuntu", "Security"]
date: "2010-03-07T21:55:00+02:00"
lastmod: "2010-03-07T21:55:00+02:00"
tags: ["WiFi", "Security", "Hacking", "Aircrack-ng", "Grimwepa"]
toc: true
---

## Introduction

This method is really designed for beginners and allows for easy cracking of wireless networks without any networking knowledge. It's not the kind of method I usually prefer since it enables 16-year-olds to think they're NASA-level hackers with these tools... but anyway.

This method is very practical when you don't have much time. For the OS, I obviously recommend BackTrack, but Ubuntu can also work. For this tutorial, I'll use Ubuntu.

## Installation

### aircrack-ng

Let's install aircrack-ng to get all the necessary binaries:

```bash
aptitude install aircrack-ng openjdk-6-jre
```

### Grimwepa

It's recommended to install grimwepa using this method:

```bash
wget http://grimwepa.googlecode.com/files/grimstall.sh
chmod 755 grimstall.sh
sudo ./grimstall.sh install
```

## Configuration

For configuration, we just need to activate monitoring mode on our wireless interface. I'm using a DLINK DWL-G122 with a RALINK chipset that allows me to perform injections, etc. To activate this mode:

```bash
sudo airmon-ng start wlan1

Found 7 processes that could cause trouble.
If airodump-ng, aireplay-ng or airtun-ng stops working after
a short period of time, you may want to kill (some of) them!

PID	Name
1200	NetworkManager
1202	avahi-daemon
1204	avahi-daemon
1478	wpa_supplicant
3183	dhclient
16461	dhclient
17269	dhclient
Process with PID 16461 (dhclient) is running on interface wlan0
Process with PID 17269 (dhclient) is running on interface wlan1

Interface	Chipset		Driver

wlan0		Intel 4965/5xxx	iwlagn - [phy0]
wlan1		Ralink 2573 USB	rt73usb - [phy1]
				(monitor mode enabled on mon0)
```

Monitor mode is now active on mon0 :-)

## Utilization

Now, let's launch grimwepa:

```bash
sudo java -jar grimwepa_X.X.jar
```

Select the mon0 interface, then click on "Refresh Targets", you should see it scanning:

![WEPA0](/images/wepa0.avif)

Stop after about 3 scans, that's sufficient. Choose a network with WEP encryption (faster because it's older and therefore easier to crack). Then select "Fragmentation" as the attack method and choose an available client. Then click on "Start Attack":

![WEPA-1](/images/wepa-1.avif)

A window should open that will listen to what's happening on this network:

![Airodump](/images/airodump.avif)

Once there is enough data (which can take some time depending on traffic), an airmon-ng window will start to launch injections. From this point, it will go relatively quickly. The WEP key cracking will follow. The key will then be displayed in the status.

## Resources
- [https://code.google.com/p/grimwepa/](https://code.google.com/p/grimwepa/)
