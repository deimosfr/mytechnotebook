---
weight: 999
url: "/Générer_un_fichier_configure_pour_pré-make/"
title: "Generate a Configure File for Pre-make"
description: "How to generate a configure file before running make when retrieving sources from SVN or CVS repositories."
categories: ["Linux"]
date: "2007-05-10T05:15:00+02:00"
lastmod: "2007-05-10T05:15:00+02:00"
tags: ["Development", "Linux", "Command Line", "Build Tools"]
toc: true
---

## Introduction

As you probably already know, the three magic commands are:

- `./configure`
- `make`
- `make install`

This is nice and cool, but if you're retrieving sources from SVN or CVS repositories, you might not always have the configure file. This is why you only need two simple commands to generate one:

```bash
/usr/bin/autoheader ; /usr/bin/autoconf
```
