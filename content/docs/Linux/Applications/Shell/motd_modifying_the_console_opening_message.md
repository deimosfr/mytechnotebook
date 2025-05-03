---
weight: 999
url: "/Motd_\\:_Modification_du_message_d'ouverture_de_console/"
title: "MOTD: Modifying the Console Opening Message"
description: "Learn how to customize your console opening message (MOTD) using various utilities like cowsay, boxes, linuxlogo, and figlet."
categories: ["Linux"]
date: "2007-08-24T21:40:00+02:00"
lastmod: "2007-08-24T21:40:00+02:00"
tags: ["Linux", "Servers", "Console", "Customization"]
toc: true
---

Don't like your console opening message? I'll show you how to change it. The file is located at:

`/etc/motd`

There are plenty of small utilities that can create drawings and other interesting displays:

## Cowsay

Allows you to put text in a small drawing (by default a cow, but other drawings are available):

```bash
sudo apt-get install cowsay
```

Example:

```bash
echo Serveur Toto|cowsay -f eyes
ls /usr/share/cowsay/cows
```

## Boxes

Allows you to place text in small drawings (mainly for adding comments in code lines):

```bash
sudo apt-get install boxes
```

Example:

```bash
echo Acces on this server is stricly restricted | boxes -d peek
```

## LinuxLogo

Lets you use colorized Linux banners:

```bash
sudo apt-get install linuxlogo
```

Example:

```bash
linux_logo
```

## FIGlet

Allows you to write text in ASCII art form:

```bash
sudo apt-get install figlet
```

Example:

```bash
figlet -f small Access Restricted
showfigfonts|more
```

## Disabling MOTD in SSH

If you no longer want any messages when booting via SSH, this command is enough:

```bash
touch ~/.hushlogin
```
