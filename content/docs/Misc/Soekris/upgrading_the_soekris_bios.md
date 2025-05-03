---
weight: 999
url: "/Upgrader_le_BIOS_de_la_Soekris/"
title: "Upgrading the Soekris BIOS"
description: "How to upgrade the BIOS firmware on a Soekris device"
categories: ["Linux", "Network"]
date: "2010-12-01T21:25:00+02:00"
lastmod: "2010-12-01T21:25:00+02:00"
tags: ["Soekris", "BIOS", "Firmware", "Embedded Systems"]
toc: true
---

## Introduction

I strongly recommend upgrading the BIOS whenever an update is available: [https://www.soekris.com/downloads.htm](https://www.soekris.com/downloads.htm)

Since it's not necessarily obvious how to do it, I'll explain the process.

## Installation

A small package will be needed:

```bash
apt-get install lrzsz
```

## Upgrading

- Connect to your BIOS monitor using minicom and run the download command:

```bash
> download
```

Or try this if it doesn't work:

```bash
> download -
```

- We'll send the file using sx:

```bash
$ sx -X ~/b5501_133.bin > /dev/ttyUSB0 < /dev/ttyUSB0
Sending b5501_133c.bin, 784 blocks: Give your local XMODEM receive command now.
Xmodem sectors/kbytes sent:  96/12kRetry 0: NAK on sector
Xmodem sectors/kbytes sent:  97/12kRetry 0: NAK on sector
Xmodem sectors/kbytes sent: 131/16kRetry 0: NAK on sector
Xmodem sectors/kbytes sent: 132/16kRetry 0: NAK on sector
Xmodem sectors/kbytes sent: 528/66kRetry 0: NAK on sector
Xmodem sectors/kbytes sent: 529/66kRetry 0: NAK on sector
Xmodem sectors/kbytes sent: 578/72kRetry 0: NAK on sector
Xmodem sectors/kbytes sent: 579/72kRetry 0: NAK on sector
Bytes Sent: 100352   BPS:1519

Transfer complete
```

Once completed, flash and reboot:

```bash
> flashupdate
?Updating BIOS Flash ,,,,,,,,,,,,,,,,,,,,,,,,,,,,..,,,,.... Done.

> reboot
```

## Reference

- [https://soekris.kd85.com/flashupdate_4801](https://soekris.kd85.com/flashupdate_4801)
- [https://mguesdon.oxymium.net/blog/?postid=128](https://mguesdon.oxymium.net/blog/?postid=128)
- [https://www.mail-archive.com/soekris-tech@lists.soekris.com/msg04537.html](https://www.mail-archive.com/soekris-tech@lists.soekris.com/msg04537.html)
- [https://wiki.soekris.info/Updating_Bios#Minicom](https://wiki.soekris.info/Updating_Bios#Minicom)
