---
weight: 999
url: "/Connaitre_son_architecture/"
title: "How to Determine Your System Architecture"
description: "A simple command to check your system's architecture and CPU capabilities"
categories: ["Solaris", "System Administration"]
date: "2007-08-08T12:30:00+02:00"
lastmod: "2007-08-08T12:30:00+02:00"
tags: ["architecture", "CPU", "64-bit", "32-bit", "system"]
toc: true
---

Here's a simple but very useful trick. The `isainfo` command allows you to determine the architecture of a machine:

```bash
$ isainfo -v
64-bit amd64 applications
        pause sse2 sse fxsr amd_3dnowx amd_3dnow amd_mmx mmx cmov amd_sysc cx8 
        tsc fpu 
32-bit i386 applications
        pause sse2 sse fxsr amd_3dnowx amd_3dnow amd_mmx mmx cmov amd_sysc cx8 
        tsc fpu
```

I can see here that my system is 64-bit and also has compatibility for 32-bit applications.
