---
title: "Hdiutil: Creating a DVD ISO from a VIDEO TS folder"
slug: hdiutil-creating-a-dvd-iso-from-a-video-ts-folder/
description: "How to create a DVD ISO from a VIDEO_TS folder using hdiutil on Mac OS X"
categories: ["Mac OS X"]
date: "2009-11-28T16:06:00+02:00"
lastmod: "2009-11-28T16:06:00+02:00"
tags: ["DVD", "ISO", "Mac OS X", "hdiutil"]
---

## Introduction

You may sometimes need to create a DVD from a VIDEO_TS folder. This is the solution.

## Usage

```bash
$ hdiutil makehybrid -udf -udf-volume-name DVD_NAME -o MY_DVD.iso /path/
```

`/path/` is the root folder of the DVD, not the VIDEO_TS folder.
