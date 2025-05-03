---
weight: 999
url: "/Debian_\\:_réinstallation_d'un_serveur_presque_à_l'identique/"
title: "Debian: Reinstalling a Server Almost Identically"
description: "How to reinstall a Debian server with the exact same packages as a reference server."
categories: ["Debian", "Linux"]
date: "2013-05-06T15:36:00+02:00"
lastmod: "2013-05-06T15:36:00+02:00"
tags: ["Debian", "Packages", "Server", "Installation"]
toc: true
---

## Principle

From a reference server, I want to quickly install another server with exactly the same packages.

## Usage

We'll extract the list of installed packages from the reference server:

```bash
dpkg --get-selections > my_packages.txt
```

Then reinstall everything from this package list:

```bash
dpkg --clear-selections
dpkg --set-selections < my_packages.txt
apt-get update
apt-get dselect-upgrade
apt-get dist-upgrade
apt-get upgrade
```
