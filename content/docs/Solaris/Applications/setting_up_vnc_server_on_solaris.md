---
weight: 999
url: "/Mise_en_place_de_vncserver_sur_Solaris/"
title: "Setting up VNC server on Solaris"
description: "This guide explains how to install and configure VNC Server on Solaris systems to enable remote desktop access."
categories: ["Unix", "Solaris", "Network"]
date: "2008-12-03T15:53:00+02:00"
lastmod: "2008-12-03T15:53:00+02:00"
tags: ["VNC", "Solaris", "Remote Desktop", "Servers", "Network"]
toc: true
---

## Introduction

Solaris, my love! You who never have packages and are never very convenient to use! In short, my difficult love! And yes, VNC server is doable! Recompiling manually is too difficult and time-consuming, so I opted for packages. But even that isn't all roses!

## Installation

### Method 1

Add this line to your profile:

```bash
PATH=$PATH:/usr/openwin/bin:/usr/X11/bin
```

### Method 2

For prerequisites, you need pkg-get! Then let's get started:

```bash
pkg-get -i vncserver
```

Let's refine things by editing `/etc/init.d/boot.server` and adding the magic line:

```bash
test -d /tmp/.X11-unix && /usr/bin/chmod 1777 /tmp/.X11-unix 
```

This line allows all users to access vncserver from boot.

That's it! Now you just need to launch it.

## Configuration

We'll start vncserver once, then kill it to get the configuration. Also enter the password when prompted:

```bash
$ vncserver
You will require a password to access your desktops.

Password:
Verify:

New 'unknown:1 (root)' desktop is unknown:1

Creating default startup script //.vnc/xstartup
Starting applications specified in //.vnc/xstartup
Log file is //.vnc/unknown:1.log
```

```bash
vncserver -kill :1
```

Now we can edit the `~/.vnc/xstartup` file for all users:

```bash
#!/bin/sh

[ -r $HOME/.Xresources ] && xrdb $HOME/.Xresources
xsetroot -solid grey
#vncconfig -iconic &
#xterm -geometry 80x24+10+10 -ls -title "$VNCDESKTOP Desktop" &
gnome-session &
```

## Launching

Let's start vncserver and we're good to go!

```bash
$ vncserver
```
