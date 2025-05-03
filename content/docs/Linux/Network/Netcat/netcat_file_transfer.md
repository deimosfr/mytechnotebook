---
weight: 999
url: "/Netcat_\\:_Transfert_de_fichiers/"
title: "Netcat: File Transfer"
description: "A guide on how to use Netcat for file transfers between systems, including HTTP sharing options."
categories:
  - "Linux"
  - "Network"
date: "2012-05-14T16:30:00+02:00"
lastmod: "2012-05-14T16:30:00+02:00"
tags:
  - "Netcat"
  - "Network"
  - "File Transfer"
  - "Linux"
  - "Servers"
toc: true
---

## Introduction

Netcat (nc) is a "...simple Unix utility which reads and writes data across network connections, using TCP or UDP protocol. It is designed to be a reliable back-end tool that can be used directly or easily driven by other programs and scripts. At the same time, it is a feature-rich network debugging and exploration tool, since it can create almost any kind of connection you would need and has several interesting built-in capabilities."

## Action

Basically it's another small, cool Unix tool that allows you to do tons of cool stuff. I found this example out there that lets you transfer files via tar from one box to another. As with anything to do with nc, it's dead simple, and logical. On the target box, start nc to listen on a port, and tar up anything it 'hears' like this:

```bash
nc -l <PORT> | tar -xf -
```

Then, on the source system, have tar pipe out to netcat, that is pointed to the target host/ip:

```bash
tar -cf - <DIRECTORY> | nc <HOST> <PORT>
```

Damn, how cool. There's plenty more info out there, and the more you look the more you'll realize what you can do with nc. Tons of great info at the above Wikipedia link, and I also found a great overview at Vulwatch.org. Have fun!

## Sharing file through http 80 port

```bash
nc -w 5 -v -l 80 < file.ext
```

From the other machine open a web navigator and go to ip from the machine who launch netcat, http://ip-address/

If you have some web server listening at 80 port then you would need stop them or select another port before launch net cat ;-)

## Resources
- [Netcat: Creating a listening port](./Netcat_\:_Créer_un_port_d'écoute.html)
- [Netcat: Remote partition backup](./netcat_\:_sauvegarde_de_partions_à_distance.html)
- [Netcat Documentation](/pdf/netcat.pdf)
