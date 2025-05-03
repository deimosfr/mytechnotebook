---
weight: 999
url: "/Compiler_avec_gcc_sur_plusieurs_architectures_(ex\\:_PPC_et_Intel)/"
title: "Compiling with GCC on Multiple Architectures (e.g., PPC and Intel)"
description: "How to compile cross-architecture binaries using GCC, specifically for PPC and Intel architectures"
categories: ["Development", "Mac OS"]
date: "2007-07-12T07:20:00+02:00"
lastmod: "2007-07-12T07:20:00+02:00"
tags: ["gcc", "compilation", "cross-compilation", "PPC", "Intel", "Mac OS X", "development"]
toc: true
---

## Introduction

This type of compilation happens in two phases. The first phase is to compile separately for Intel and PPC architectures. The second phase is to create a binary that combines both architecture binaries.

## Creating Binaries for Different Architectures

On Intel, I can compile with gcc like this:

```bash
gcc -arch ppc -isysroot /Developer/SDKs/MacOSX10.4u.sdk prog.c
```

For the Makefile, edit it and look for the following options to specify the architecture:

```
-arch and -isysroot
```

## Assembly

Once you have your two architecture binaries, you need to merge them:

```bash
lipo -create ppc/prog i386/prog -output prog
```

Note: It seems that at the gcc compilation level, you can simply use the following options to avoid having to do the assembly step:

```bash
-arch ppc -arch i386
```
