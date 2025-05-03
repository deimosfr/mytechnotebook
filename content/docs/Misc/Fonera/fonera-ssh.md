---
weight: 999
url: "/activer-le-ssh-sur-sa-fonera-plus/"
title: "Enabling SSH on the Fonera+"
description: "A guide to enable SSH access on a Fonera+ router by flashing it with a custom firmware."
categories: ["Networking", "Linux", "Hardware"]
date: "2009-05-06T20:40:00+02:00"
lastmod: "2009-05-06T20:40:00+02:00"
tags: ["Fonera", "SSH", "Router", "Firmware", "OpenWRT"]
toc: true
---

## Introduction

The Fonera is a small device based on OpenWRT that allows you to distribute WiFi (HotSpot). To learn more, I invite you to visit the official website. In short, it is possible to play a bit more with the basic functions, which is why I tackled SSH access on this device.

**Important: Any firmware modification can take 10 minutes to update. So don't reboot before that!**

## Prerequisites

* A web server

Apache will do fine... (on Mac: /Library/WebServer/Documents)

* Remember to turn off your firewall during the process
* Perl and the perl-Net-Telnet dependency:

Open cpan and do:

```bash
install Net::Telnet
```

* Install fping:

```bash
apt-get install fping
```

or on Mac:

```bash
sudo port install fping
```

* This perl file [redboot.pl](https://download.francofon.fr/fonera_plus_ssh/redboot.pl) to be placed in the root of your web server
* The [firmware](https://download.francofon.fr/fonera_plus_ssh/firmware_francofon.bin) which is also to be placed in the root of your web server
* A direct network connection with the Fonera (via its WAN port, the black one)
* Configure your IP address to 192.168.1.254

## Flashing with the New Firmware

Launch the redboot.pl script that you downloaded like this:

```bash
perl redboot.pl 192.168.1.1
```

Once the connection is established, specify the Fonera's IP:

```bash
ip_address -l 192.168.1.1/24 -h 192.168.1.254
```

then type:

```bash
fis delete image
load -r -b 0x80100000 /firmware_francofon.bin -m HTTP -h 192.168.1.254
fis create -b 0x80100000 -l 0x00237040 -f 0xA8040000  -e 0x80040400  -r 0x80040400 image
```

You should see something like this:

```bash
192.168.1.1 is unreachable
192.168.1.1 is alive
-> == Executing boot script in 1.910 seconds - enter ^C to abort
<- ^C
Trying 192.168.1.1...
Connected to 192.168.1.1.
Escape character is '^]'.
RedBoot> ip_address -l 192.168.1.1/24 -h 192.168.1.254
IP: 192.168.1.1/255.255.255.0, Gateway: 0.0.0.0
Default server: 192.168.1.254
```

```bash
RedBoot> fis delete image
Delete image 'image' - continue (y/n)? y
... Erase from 0xa8040000-0xa8277040: ....................................
... Erase from 0xa87e0000-0xa87f0000: .
... Program from 0x80ff0000-0x81000000 at 0xa87e0000: .
```

```bash
RedBoot> load -r -b 0x80100000 /firmware_francofon.bin -m HTTP -h 192.168.1.254
Raw file loaded 0x80100000-0x8033703f, assumed entry at 0x80100000
```

```bash
RedBoot> fis create -b 0x80100000 -l 0x00237040 -f 0xA8040000  -e 0x80040400  -r 0x80040400 image
... Erase from 0xa8040000-0xa8277040: ....................................
... Program from 0x80100000-0x80337040 at 0xa8040000: ....................................
... Erase from 0xa87e0000-0xa87f0000: .
... Program from 0x80ff0000-0x81000000 at 0xa87e0000: .
```

Wait until the flashing is complete. A quick reboot:

```bash
RedBoot> reset
```

You should now have SSH access to your Fonera :-):

```bash
ssh -l root 192.168.10.1
The authenticity of host '192.168.10.1 (192.168.10.1)' can't be establish
RSA key fingerprint is 5c:d3:42:ed:52:6d:c0:c6:fb:ec:84:57:18:24:d7:be.
Are you sure you want to continue connecting (yes/no)? yes
Warning: Permanently added '192.168.10.1' (RSA) to the list of known host
root@192.168.10.1's password:
BusyBox v1.4.1 (2007-09-03 10:39:50 UTC) Built-in shell (ash)
Enter 'help' for a list of built-in commands.
 ______                                           __
/\  ___\                                         /\ \
\ \ \__/  __     ___      __   _ __    __        \_\ \___
 \ \  _\/ __`\ /' _ `\  /'__`\/\`'__\/'__`\     /\___  __\
  \ \ \/\ \L\ \/\ \/\ \/\  __/\ \ \//\ \L\.\_   \/__/\ \_/
   \ \_\ \____/\ \_\ \_\ \____\\ \_\\ \__/.\_\      \ \_\
    \/_/\/___/  \/_/\/_/\/____/ \/_/ \/__/\/_/       \/_/
--------------  Fonera 1.5 Firmware (v1.1.1.1) -----------------
            * Based on OpenWrt - http://openwrt.org
            * Powered by FON - http://www.fon.com
     -----------------------------------------------------
root@OpenWrt:~# Wow!!! Your Fonera+ is now FREE!
```

## FAQ

### The redboot script cannot find my Fonera

Make sure the POWER light is on. If not, unplug the Fonera for 10 seconds and then plug it back in!

### FON_ATTENTION_PLEASE_CONNECT

When I scan the WiFi at home, I discover an unencrypted SSID called "FON_ATTENTION_PLEASE_CONNECT". This is a failed flashing. So you need to download the [official firmware](https://www.fon.com/en/download#) from the official site to restore the original firmware.

Connect to this signal and then type [https://192.168.1.1](https://192.168.1.1). Now, upload the firmware and you will see something like this:

```bash
Firmware upgrade and hotfix installation

Ooooops! Looks like the La Fonera is not working properly. You need to reinstall the fon software in it. Please, provide a valid full firmware in the box below (you can find them at fon's download page) or contact fon support at support@fon.com
FON upgrade file 	
Upgrading... 

Wait for the router to reflash itself. This can take up to some minutes. DO NOT DICONNECT THE LA FONERA in 10min

This is a FON reflash v2 archive
Verified OK
Upgrade name: reflash_all
grep: /etc/hotfix: No such file or directory
Upgrading FON firmware and rebooting...
This may take up to 10 minutes. Please be patient.
The power light will be alternating green and orange. When the process is finished the light will stay orange while rebooting
Flashing image...

The upgrade process was successful

Press here to return to the index page
```

### Updating personal config from FON

Just after flashing your Fonera+ will reset with factory default settings. You can verify this by going into your HTTP console on 192.168.10.1

To update your config, log on to www.fon.com, and access userzone. Select your router, and update WLAN private and public SSID names. If you don't want to change the name, please just change one letter, click on "update" button, and change again to the right name. For the private WLAN: change the WEP/WPA key encryption using the same method.

Fon.com servers will send the new config to your Fonera+. Wait a few minutes and check in your local HTTP console. You don't need to reboot.

### Registered or not?

If your Fonera+ has been registered before the SSH-unlock, check on your local HTTP console status if all is OK. If the logo displayed is "your Fonera+ has not been registered", it is important to change this parameter to give access to users on your public WLAN.

To do this, open SSH console:

```bash
echo 1 > /etc/config/registered
```

Reboot your Fonera+, connect again to your HTTP local console, and verify the change to the logo: "Your Fonera is registered OK"

### Bandwidth, QoS, transfer rate

Once your Fonera+ is running, configured and registered, check your transfer rate on all ports!

Default settings in original FON firmware 1.1.1r1 are 1024kb/s for download and 128kb for upload (WAN port).

Adjust these settings to your ISP speed line, in this example 2048kbs for D/L and 256 kbs for U/L.

```bash
uci set qos.wan.upload=256
uci set qos.wan.download=2048
uci commit
```

Reboot and perform a new speed transfer test on WLAN and LAN.

Disable FON QoS service (not recommended):

```bash
uci set qos.wan.enabled=0
uci commit
```

## Resources
- [https://www.jopa.fr/index.php/2008/03/24/jouer-avec-la-fonera-2eme-partie-hacker-la-fonera/](https://www.jopa.fr/index.php/2008/03/24/jouer-avec-la-fonera-2eme-partie-hacker-la-fonera/)
- [https://www.cs.helsinki.fi/u/sklvarjo/lafon.htm](https://www.cs.helsinki.fi/u/sklvarjo/lafon.htm)
- [https://www.dd-wrt.com/dd-wrtv3/dd-wrt/downloads.html](https://www.dd-wrt.com/dd-wrtv3/dd-wrt/downloads.html)
