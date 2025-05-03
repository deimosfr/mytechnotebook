---
weight: 999
url: "/Activer_Desactiver_ecran_veille_ligne_commande/"
title: "Enabling and Disabling Screen Savers from Command Line"
description: "How to control screen power management using the command line in Linux for automatic display shutoff and power saving"
categories: ["Linux", "Desktop"]
date: "2008-07-31T10:21:00+02:00"
lastmod: "2008-07-31T10:21:00+02:00"
tags:
  ["x11", "xset", "dpms", "screensaver", "power management", "cron", "linux"]
toc: true
---

## Introduction

For my job, we got nice 42-inch LCD screens for monitoring. However, these are Toshiba models that don't even include a timer for turning the displays on and off. Instead of manually turning them on and off every morning, I preferred to use screen saver activation and deactivation features to manage them automatically.

## Usage

Here's a test that should turn off your screen, wait, and then turn it back on:

```bash
xset dpms force off; sleep 5; xset dpms force on; xrandr -s 1; xrandr -s 0
```

To automate this process, you'll want to create a crontab entry. You'll need to add the following to your crontab:

```
DISPLAY=:0.0
0 20 * * 1-5 xset dpms force off
0 8 * * 1-5 xset dpms force on
```

This will turn off the displays at 8:00 PM and turn them back on at 8:00 AM, Monday through Friday.
