---
weight: 999
url: "/Gnome_\\:_Tableaux_de_bord_verrouillés/"
title: "Gnome: Locked Dashboards"
description: "How to lock Gnome dashboards to prevent panels and applet icons from moving around accidentally."
categories: ["Linux"]
date: "2009-08-07T16:14:00+02:00"
lastmod: "2009-08-07T16:14:00+02:00"
tags: ["Gnome", "Desktop", "Configuration"]
toc: true
---

## Introduction

In Gnome, dashboard panels and applet icons are so movable that they tend to move around and wander everywhere.

While this option is certainly very flexible, it quickly becomes annoying in daily use. The solution: lock the dashboards!

Gnome includes many configuration tools that allow each user to perfect and control their work environment. One of the most interesting tools is the Configuration Editor, gconf-editor. Warning: gconf-editor is for advanced users. If you don't understand what you're doing, don't touch everything! You could break things!

## Implementation

Launch the gconf-editor command from a terminal.

In the left panel, open the "apps" directory, then the "panel" directory, and finally "global". Then, in the right part of the window, look for the "locked_down" line and check the box next to it. That's it.

## Resources
- [https://ubunteros.tuxfamily.org/spip.php?article210](https://ubunteros.tuxfamily.org/spip.php?article210)
