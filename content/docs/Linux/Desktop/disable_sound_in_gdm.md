---
weight: 999
url: "/Désactiver_le_son_sous_GDM/"
title: "Disable Sound in GDM"
description: "How to disable the login sound in GDM on Ubuntu systems for a quieter computing experience."
categories: ["Ubuntu", "Linux"]
date: "2009-10-31T21:36:00+02:00"
lastmod: "2009-10-31T21:36:00+02:00"
tags: ["Servers", "Ubuntu", "Sound", "GDM", "Configuration"]
toc: true
---

## Introduction

If like me, you want to have Ubuntu as silent as possible, you need to disable the startup sound in GDM.

## Command

```bash
sudo -u gdm gconftool-2 --set /desktop/gnome/sound/event_sounds --type bool false
```

## Resources
- [https://www.webupd8.org/2009/10/turn-off-login-sound-in-ubuntu-karmic.html](https://www.webupd8.org/2009/10/turn-off-login-sound-in-ubuntu-karmic.html)
