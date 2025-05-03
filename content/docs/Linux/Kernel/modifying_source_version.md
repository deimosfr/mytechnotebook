---
weight: 999
url: "/Modifier_la_version_des_sources/"
title: "Modifying Source Version"
description: "How to modify kernel source version to match your running kernel for compiling software"
categories: ["Linux"]
date: "2007-06-20T15:05:00+02:00"
lastmod: "2007-06-20T15:05:00+02:00"
tags: ["Kernel", "Development", "Compilation"]
toc: true
---

- You have recompiled your own little kernel, GREAT!
- Your kernel is called "2.6.21-a_bibi"
- You want to recompile a software that relies on the source/headers of your kernel and "BAM!" it doesn't work: "The sources you are using do not match your kernel! You must be kidding!!!" and so on.

**I say NO ladies and gentlemen!**

All you need to do is modify your sources to fool your kernel!

Here is the content of the file **version.h** located in **\<source_path\>/include/linux/version.h**:

```bash
#define UTS_RELEASE "2.6.21-a_bibi"
#define LINUX_VERSION_CODE 132628
#define KERNEL_VERSION(a,b,c) (((a) << 16) + ((b) << 8) + (c))
```

You need to modify the UTS_RELEASE line to make it match with the result of the `uname -r` command on your machine.

Additionally, you need to do the same in the **utsrelease.h** file located in **\<source_path\>/include/linux/version.h**:

```bash
#define UTS_RELEASE "2.6.21-a_bibi"
```

And then you say wowww!
