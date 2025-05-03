---
weight: 999
url: "/Récupérer_son_OpenBSD_après_une_mauvaise_manip/"
title: "Recovering your OpenBSD after a bad manipulation"
description: "Guide on how to recover an OpenBSD system after accidentally deleting the /dev directory"
categories: ["Linux", "BSD", "System Recovery"]
date: "2008-01-20T07:58:00+02:00"
lastmod: "2008-01-20T07:58:00+02:00"
tags: ["OpenBSD", "Recovery", "System Administration", "Troubleshooting"]
toc: true
---

## 1. Concrete case

While manipulating my external partitions, I accidentally executed:

```bash
# rm -fr /dev
```

At boot time I get:

```bash
"/dev/console not found"
```

How to fix this problem?

## 2. Solutions

- A simple method with minimal manipulation is to boot from a ramdisk (your /bsd.rd file that you should have, or an installation CD) and perform an Upgrade.

  This won't overwrite your configuration files, just the files provided in the base system excluding etc42.tgz and xetc42.tgz.

- If you want to do it manually, boot from the ramdisk, mount your root partition, recreate the dev directory, copy the MAKEDEV script inside (script that you'll find in `/usr/src/etc/etc.<arch>/MAKEDEV`, or in base42.tgz) and execute:

```bash
cd dev && sh MAKEDEV all
```

- Alternatively, boot from the CD, then choose shell (install, upgrade or shell). Execute `/dev/MAKEDEV all` (Retrieve the MAKEDEV script for your version)
