---
weight: 999
url: "/Modifier_son_kernel_générique/"
title: "Modifying Your Generic Kernel"
description: "How to modify and switch between kernel types in BSD, including changing from single core to multi-core and booting with alternate kernels."
categories: ["Linux", "BSD"]
date: "2007-11-13T09:15:00+02:00"
lastmod: "2007-11-13T09:15:00+02:00"
tags: ["kernel", "BSD", "multicore", "system administration"]
toc: true
---

## Introduction

For those coming from the Linux world, they'll find that this is super simple on BSD. In my case, I have a multicore server, and after installing BSD, it only detects one core. It's a shame to run with just one core when you have several. I'll explain here the procedure to switch to multicore, but this works with other kernel modifications as well!

## Single Core to Multicore

During installation, you need to select the "bsd.mp" kernel. When the machine reboots after installation, it boots on the single-core kernel. We will therefore move the current kernel and replace it with the multicore one. To do this, it's very simple:

```bash
mv /bsd /bsd.mono
```

Here we moved the old kernel to bsd.mono (for single processor) and now we'll rename the multiprocessor kernel to the default kernel name (**because /bsd is the default kernel**):

```bash
mv /bsd.mp /bsd
```

Now all you have to do is reboot. That wasn't complicated, was it?

## Booting on Another Kernel

If you don't want to boot with the default kernel, you can do this before the kernel loads:

```bash
boot> b /bsd.mono
```

or

```bash
boot> boot hd0a:/bsd.mono
```
