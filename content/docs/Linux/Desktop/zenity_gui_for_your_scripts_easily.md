---
weight: 999
url: "/Zenity_\\:_Du_GUI_pour_vos_scripts_simplement/"
title: "Zenity: GUI for Your Scripts Easily"
description: "Learn how to use Zenity, a tool developed by the Gnome team that allows you to easily display GTK dialog boxes from shell scripts."
categories: 
  - Linux
date: "2009-12-14T08:31:00+02:00"
lastmod: "2009-12-14T08:31:00+02:00"
tags:
  - GUI
  - Shell
  - Scripts
  - Zenity
toc: true
---

## Introduction

I recently discovered with delight [Zenity](https://library.gnome.org/users/zenity/index.html.fr), developed by the Gnome team, which allows you to easily display GTK dialog boxes from shell scripts.

There is also dialog for ASCII GUI options.

## Example of Restarting a Service

Here is an example to restart the openvpn service (`openvpn-restart.sh`):

```bash
zenity --question --text="Start OpenVPN ?" --ok-label="Yes" --cancel-label="No"
if [ $? = 0 ] ; then
    sudo /etc/init.d/openvpn restart
fi
```
