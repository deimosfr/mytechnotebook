---
weight: 999
url: "/Infinality_fonts_for_retina_display/"
title: "Infinality fonts for retina display"
description: "Installing and configuring Infinality fonts for high DPI screens and retina displays on Linux"
categories: ["Linux", "Desktop", "Customization"]
date: "2014-02-17T08:15:00+02:00"
lastmod: "2014-02-17T08:15:00+02:00"
tags: ["Infinality", "Fonts", "Retina", "Linux", "Debian"]
toc: true
---

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Operating System** | Debian 8 |
| **Website** | [Infinality Website](https://www.infinality.net/blog/) |
| **Last Update** | 17/02/2014 |
{{< /table >}}

## Introduction

Big resolution with small screens is more and more common and we're facing a problem on fonts visualization. We have a lot of non adapted fonts for retina display, even if most of them are not so bad. The solution is to install additional fonts to replace the default one. One of the best one I've seen is [Infinality fonts](https://www.infinality.net/blog/)!

## Installation

### Manual

Clone the git repository:

```bash
cd /tmp
git clone https://github.com/chenxiaolong/Debian-Packages.git
cd Debian-Packages/
```

Install the build dependencies. Run the following command and install the packages it lists using apt-get/synaptic/etc.:

```bash
aptitude install docbook-to-man libx11-dev x11proto-core-dev libz-dev quilt
cd freetype-infinality/
dpkg-checkbuilddeps
cd ../fontconfig-infinality/
dpkg-checkbuilddeps
```

Build the packages:

```bash
cd ../freetype-infinality/
./build.sh
cd ../fontconfig-infinality/
./build.sh
```

Install the deb files:

```bash
cd ..
sudo dpkg -i freetype-infinality/*.deb fontconfig-infinality/*.deb
```

### With packages

You can grab them directly from this site[^1]:

```bash
cd /tmp
wget https://dl.dropboxusercontent.com/u/106654446/infinality_jessie/fontconfig-infinality_1-2_all.deb
wget https://dl.dropboxusercontent.com/u/106654446/infinality_jessie/freetype-infinality_2.4.9-3_all.deb
wget https://dl.dropboxusercontent.com/u/106654446/infinality_jessie/libfreetype-infinality6_2.4.9-3_amd64.deb
dpkg -i *.deb
```

Or via the ppa http://ppa.launchpad.net/no1wantdthisname/ppa/ubuntu/[^2]

## Select your configuration

Select the configuration

```bash
> sudo /etc/fonts/infinality/infctl.sh setstyle
Select a style:
1) debug       3) linux	      5) osx2	     7) win98
2) infinality  4) osx	      6) win7	     8) winxp
#? 3
```

Restart your session and you're done :-)

## References

[^1]: http://forums.debian.net/viewtopic.php?f=16&t=88545
[^2]: http://www.webupd8.org/2013/06/better-font-rendering-in-linux-with.html
