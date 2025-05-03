---
weight: 999
url: "/Hdiutil_\\:_Créer_un_ISO_DVD_depuis_un_dossier_VIDEO_TS/"
title: "Hdiutil: Creating a DVD ISO from a VIDEO TS folder"
description: "How to create a DVD ISO from a VIDEO_TS folder using hdiutil on Mac OS X"
categories: ["Mac OS X"]
date: "2009-11-28T16:06:00+02:00"
lastmod: "2009-11-28T16:06:00+02:00"
tags: ["DVD", "ISO", "Mac OS X", "hdiutil"]
toc: true
---

## Introduction

You may sometimes need to create a DVD from a VIDEO_TS folder. This is the solution.

## Usage

```bash
$ hdiutil makehybrid -udf -udf-volume-name DVD_NAME -o MY_DVD.iso /path/
```

`/path/` is the root folder of the DVD, not the VIDEO_TS folder.
