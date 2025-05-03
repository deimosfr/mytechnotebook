---
weight: 999
url: "/Créer_une_image_ISO/"
title: "Creating an ISO Image"
description: "Methods for creating ISO images in Linux using simple commands like cat and dd."
categories: ["Linux"]
date: "2007-08-14T06:19:00+02:00"
lastmod: "2007-08-14T06:19:00+02:00"
tags: ["Servers", "Mac OS X", "Divers", "Windows", "Solaris"]
toc: true
---

## Introduction

Creating an ISO image in Windows can be somewhat tedious in some cases (e.g., bootable CDs, etc.). In Linux, one might think it would be complicated, but it's actually very simple.

## Method using cat

The easiest method:

```bash
cat /dev/hda > ~/image.iso
```

hda: corresponds to your CD device (using `/dev/cdrom` should also work)

## Method using dd

Here's another method that's not much more complicated:

```bash
dd if=/dev/cdrom of=winxp.iso
```

And if you want compression afterward:

```bash
dd if=/dev/ad0 bs=8192 | gzip > my_image_file.dd.gz
```
