---
weight: 999
url: "/Compiz_\\:_Mise_en_place_d'un_bureau_3D/"
title: "Compiz: Setting Up a 3D Desktop"
description: "Guide to setting up Compiz for a 3D desktop environment on Linux systems"
categories: ["Linux", "Desktop"]
date: "2007-10-13T07:52:00+02:00"
lastmod: "2007-10-13T07:52:00+02:00"
tags: ["compiz", "3D", "desktop", "linux", "nvidia", "gnome", "kde"]
toc: true
---

## Introduction

For those who are not familiar, Compiz is one of the simplest ways to have a 3D desktop. Unlike compiz-fusion, it is less advanced, but already offers many features. Here, we'll see how to do a quick and simple deployment.

## Installation

For the installation, we will install all the compiz packages:

```bash
apt-get install compiz compiz-core compiz-gnome compiz-gtk compiz-plugins libdecoration0
```

Note: Replace **gnome** with kde in the package name if you are using KDE.

## Configuration

### Drivers

Make sure your NVIDIA drivers are properly installed (I was only able to test with NVIDIA). Install your kernel headers, gcc, and then install the drivers.
Don't forget to stop X before installing the NVIDIA drivers:

```bash
/etc/init.d/gdm stop
```

Note: Or kdm for KDE

To test that your drivers are properly installed with 3D acceleration, you can run glxgears. If it's choppy, your drivers were not installed correctly:

```bash
glxgears
```

### xorg.conf

Create a backup of your current file:

```bash
cp /etc/X11/xorg.conf /etc/X11/xorg.conf.bak
```

Then add this to the file /etc/X11/xorg.conf:

```
Section "Device"
    Identifier     "nVidia Corporation G70 [GeForce 7600 GT]"
    Driver         "nvidia"
    Option "XAANoOffscreenPixmaps" "true"
    Option "AllowGLXWithComposite" "true"
    Option "TripleBuffer" "true"
EndSection
```

Then add the following at the end of the file:

```
Section "Extensions"
          Option  "Composite" "Enable"
          Option "RenderAccel" "true"
          Option "AllowGLXWithComposite" "true"
EndSection
```

Now restart your graphical session:

```bash
/etc/init.d/gdm start
```

### GNOME

For GNOME, we'll just indicate that we want to use the 3D cube and some nice effects:

```bash
gconftool --set /apps/compiz/general/allscreens/options/active_plugins --type list --list-type string '[gconf,png,svg,decoration,wobbly,fade,minimize,cube,rotate,zoom,scale,move,place,switcher,screenshot,resize]'
```

## Launch

Now that everything is set up, you can launch compiz with this command:

```bash
compiz --replace
```

## References

[NVIDIA Drivers](https://www.nvidia.fr/object/linux_fr.html)
[Compiz Website](https://compiz.org/)
[Compiz Fusion for those who want more advanced effects](https://www.compiz-fusion.org/)
