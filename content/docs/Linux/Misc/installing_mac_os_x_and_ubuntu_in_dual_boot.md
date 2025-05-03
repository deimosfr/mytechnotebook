---
weight: 999
url: "/Installer_Mac_OS_X_et_ubuntu_en_dual_boot/"
title: "Installing Mac OS X and Ubuntu in Dual Boot"
description: "A guide on how to install Mac OS X and Ubuntu in dual boot configuration on the same hard drive."
categories: ["Linux", "Ubuntu", "Mac"]
date: "2008-10-05T09:07:00+02:00"
lastmod: "2008-10-05T09:07:00+02:00"
tags: ["Dual Boot", "Ubuntu", "Mac OS X", "Refit", "Installation"]
toc: true
---

## Introduction

If you want to have Ubuntu alongside your Mac OS X on the same hard drive and be able to choose which OS to launch at boot time, then this article is for you.

## Setup

Here are the steps:

- Install Mac OS X
- Install [Refit](https://refit.sourceforge.net/) and enable it:

```bash
sudo /efi/refit/enable-always.sh
```

- Once completed, launch Boot Camp and decide how much space you want to allocate
- Reboot from the Ubuntu CD
- Install [Ubuntu](https://www.ubuntu-fr.org/)
- When Ubuntu reboots, go back to Mac OS X (as it's impossible to boot Ubuntu for now)
- Reinstall [Refit](https://refit.sourceforge.net/) over the existing installation, then enable it again:

```bash
sudo /efi/refit/enable-always.sh
```

- Now you can reboot to boot into Ubuntu :-)
