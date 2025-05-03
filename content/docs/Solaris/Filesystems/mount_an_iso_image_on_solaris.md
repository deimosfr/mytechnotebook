---
weight: 999
url: "/Monter_une_image_ISO_sous_Solaris/"
title: "Mount an ISO image on Solaris"
description: "How to mount an ISO image on Solaris using lofiadm and mount commands."
categories: ["Unix", "Solaris"]
date: "2010-06-15T07:24:00+02:00"
lastmod: "2010-06-15T07:24:00+02:00"
tags: ["ISO", "Mount", "Solaris", "lofiadm", "hsfs"]
toc: true
---

## Introduction

Many software packages can be downloaded in the form of an ISO image. Rather than burning the image to a CD-ROM to access its contents, it is easy to mount the image directly into the filesystem using the lofiadm and mount commands.

## Usage

Given an ISO image in `/export/temp/software.iso`, a loopback file device (`/dev/lofi/1`) is created with the following command:

```bash
lofiadm -a /export/temp/software.iso /dev/lofi/1
```

The lofi device creates a block device version of a file. This block device can be mounted to `/mnt` with the following command:

```bash
mount -F hsfs -o ro /dev/lofi/1 /mnt
```

These commands can be combined into a single command:

```bash
mount -F hsfs -o ro `lofiadm -a /export/temp/software.iso` /mnt
```

## Resources
- [https://www.tech-recipes.com/rx/218/mount-an-iso-image-on-a-solaris-filesystem-with-lofiadm/](https://www.tech-recipes.com/rx/218/mount-an-iso-image-on-a-solaris-filesystem-with-lofiadm/)
