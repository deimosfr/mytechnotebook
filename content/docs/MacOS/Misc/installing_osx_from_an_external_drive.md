---
weight: 999
url: "/Installer_OSX_depuis_un_disque_externe/"
title: "Installing OSX from an External Drive"
description: "How to install OSX when your Mac doesn't have a working DVD drive by using an external drive."
categories: ["Mac OS X"]
date: "2007-10-27T13:46:00+02:00"
lastmod: "2007-10-27T13:46:00+02:00"
tags: ["Mac OS X", "installation", "recovery"]
toc: true
---

## Introduction

Maybe your machine doesn't have a DVD drive or it's simply not working. How can you install OSX again? The solution is here...

## Instructions

- Create a disk image using Disk Utility (on Mac A)
- Transfer the .dmg file to the target computer (Mac B) via FireWire or Ethernet
- You'll need an external FireWire hard drive, name it "Mac OS X Install DVD" with Disk Utility
- Select the newly created partition and click on the Restore tab
- Drag and drop the Tiger disk image into **Source** and the new partition into **Destination**. Then click on Restore
- Once the restoration is complete, go to System Preferences / Startup Disk, choose the new partition as the startup disk, then restart

![External OSX Installation](/images/dmg-01-tm.avif)
