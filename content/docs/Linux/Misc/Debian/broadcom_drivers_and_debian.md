---
weight: 999
url: "/Drivers_Broadcom_et_Debian/"
title: "Broadcom Drivers and Debian"
description: "Installation guide for Broadcom drivers on Debian systems, providing a simple solution for hardware compatibility issues."
categories: ["Debian", "Linux"]
date: "2009-02-16T18:55:00+02:00"
lastmod: "2009-02-16T18:55:00+02:00"
tags: ["Servers", "Linux", "Drivers", "Hardware"]
toc: true
---

## Introduction

Troublesome drivers! Thanks to manufacturers who provide drivers with problematic licenses. In short, with Debian version 5, network cards with Broadcom chipsets stopped working properly.

No problem, here's the solution.

## Installation

Simply install this package:

```bash
apt-get install firmware-bnx2
```

A quick reboot and you're all set!

## Resources
- http://forum.hardware.fr/hfr/OSAlternatifs/reseaux-securite/activer-broadcom-debian-sujet_67918_1.htm
