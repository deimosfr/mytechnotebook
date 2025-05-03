---
weight: 999
url: "/Créer_un_Apple_Time_Machine_réseaux_sous_FreeBSD/"
title: "Creating an Apple Time Machine Network on FreeBSD"
icon: "apple"
icontype: "simple"
description: "How to set up a FreeBSD server as a Time Machine network backup destination for Mac OS X."
categories: ["FreeBSD", "Backup", "Mac OS X"]
date: "2009-01-28T02:29:00+02:00"
lastmod: "2009-01-28T02:29:00+02:00"
tags: ["Time Machine", "Apple", "Backup", "FreeBSD", "Network", "Mac OS X"]
toc: true
---

## Introduction

Here's a quick guide on how to set up Time Machine on Mac OS X to back up to a networked machine running FreeBSD.

## On the FreeBSD server

- Build & Install net/netatalk from ports.
- Edit `/usr/local/etc/AppleVolumes.default`
- Append: "/your_time_machine_path TimeMachine allow:your_user_name cnidscheme:cdb options:usedots" and replace your path and your username in the proper places.
- Optionally, remove the "~" already present in that file if you don't want to share users home directories.
- Add "netatalk_enable="YES"" and "afpd_enable="YES"" to `/etc/rc.conf`.
- `/usr/local/etc/rc.d/netatalk start` (nothing will be printed).

## Mac OS X machine client (running Leopard, of course)

- Mount your remote volume. Command+K on the Finder and then type: "afp://<machine IP address or local hostname if you have a local DNS server>". You can't type the machine name because we're not using multicast DNS.
- Build a sparse bundle image using "Disk utility" (HFS+ case-sensitive formatted). Usually, the size should be something that gives you enough room for expansion. If you want to backup your whole MacBook/iMac/etc. disk, you can set the sparse bundle image size the same as the disk your are backing up.
- The name of this image is important. It should be "Your_Computer_Name_MACAddress.sparsebundle". Check your computer name from the "Sharing" section of "System Preferences" and the MAC address comes from the interface you'll be using to do the backup. I really recommend using your Wired interface. Check the MAC address via ifconfig(1) or via the "Network" section of "System Preferences". E.g., if you're John Doe, have a MacBook and your MAC address is 00:01:02:03:04:05, your file should be named "John Doe's MacBook_000102030405.sparsebundle".
- On the Terminal, type "defaults write com.apple.systempreferences TMShowUnsupportedNetworkVolumes 1". This is the crucial thing.
- Go to "System Preferences", "Time Machine" and enable it. The networked volume will now show up on the list.
- Before selecting the volume on which you'll dump the backup, copy the sparse bundle file you've created to your networked volume called "TimeMachine".
- Select the networked Volume from the Time Machine volumes list.
- Initiate the backup!

Enjoy!

As Remko points out in the comments, the MAC address is not restrictive. So, if you want to backup via wired interface and after that via wireless, Time Machine will work using both interfaces. I suppose that Time Machine inspects all MAC addresses in your machine and then searches a sparse bundle in the networked volume that matches.

## Resources

- [https://blogs.freebsdish.org/rpaulo/2008/10/04/apple-time-machine-freebsd-in-14-steps/](https://blogs.freebsdish.org/rpaulo/2008/10/04/apple-time-machine-freebsd-in-14-steps/)
