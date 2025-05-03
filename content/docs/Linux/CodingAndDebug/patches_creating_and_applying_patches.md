---
weight: 999
url: "/Patchs_\\:_Création_et_applications_de_patchs/"
title: "Patches: Creating and Applying Patches"
description: "A guide on how to create and apply patches in Linux, including the use of diff and patch commands for file modification tracking and distribution."
categories: ["Linux"]
date: "2013-05-07T12:39:00+02:00"
lastmod: "2013-05-07T12:39:00+02:00"
tags: ["Patching", "Linux", "Development", "diff", "patch", "Command Line"]
toc: true
---

## Introduction

When translating or modifying open-source software, it's important to find a simple and effective way to redistribute your work. To do this, you need to distribute it in the lightest possible manner. The solution is to distribute only the difference between the original version and yours. In short, you need to create a patch.

## Patching

This term refers to the action of applying a patch to one or more files. By patching a file, you automatically make the necessary modifications to update it. The patch contains the list of differences between the old and new versions of the file(s). Result: this difference file is much lighter and can be applied in a single operation.

The command used to apply a patch in Linux is naturally called patch. This utility, created by Larry Wall (creator of Perl), does more than just "look at" and change files. In fact, patch can react based on context and apply a patch to already modified files.

## Creating a Patch

To create a difference file, you need to use the diff command, which analyzes and determines the "contextual differences" between two files. It can be used recursively to analyze files between two distinct directories. Here's a concrete example to understand how the diff command works:

Let's say we have a C source file hello.c.old containing the following text:

```c
 //Demonstration file for diff
 //Creator :Diamonds Edition 1998
 #include<stdio.h>
 void main(void) {
  printf("Hello World !");
 }
```

We modify it to create a French version under the name hello.c, which gives us:

```c
//Fichier de démonstration pour diff
//Créateur :Copyright 2000 Diamond Editions/Linux magazine France 1998
#include<stdio.h>
void main(void) {
 printf("Bonjour le Monde !");
 }
```

Then, to create a difference file, we use the command:

```bash
diff -c hello.c.old hello.c > hello.diff
```

The -c indicates a contextual comparison, resulting in a hello.diff file that can then be distributed to anyone who has the English version of hello.c.

Often, the source code of software is not limited to a single file. The program is distributed across multiple source files that will be compiled, then linked to form the executable program. Translation (and therefore modification) spans multiple files and sometimes multiple subdirectories.

The method for creating the patch is to modify the sources in their original directory, then install the original sources in a different directory. Let's say the /h.old directory contains the English version of hello.c and the /h directory will contain the French version. Place yourself at the root of the disk and type:

```bash
diff -cr h.old h >hello.diff
```

The resulting hello.diff file will contain the contextual differences between the files in the /h.old directory (the English original) and /h (the French version). The r parameter passed in addition to -c indicates recursive operation.

To create a clean and more readable diff:

```bash
diff -b -Nur h.old h >hello.diff
```

If you want to see a colored diff, use colordiff instead.

### Statistics

To get statistics on a diff, you can use the diffstat command:

```bash
$ diff test1 test2 | diffstat
unknown |    2 +-
1 file changed, 1 insertion(+), 1 deletion(-)
```

## Applying a Patch

To apply the patch from the first example, we'll use:

```bash
patch < hello.diff
```

in the directory where the English hello.c file is located. The patch will be applied and the hello.c file will be modified to become our French version.

In the case of a patch applying to multiple files in a directory, copy the .diff file to the root of the disk, then type:

```bash
patch -p0 <hello.diff
```

The -p0 parameter allows patch to work recursively. The file in the /h directory will be modified to become our French version.

If it asks you to specify the files, it's because it's trying to compare directories that don't exist. For example, when I want to upgrade MediaWiki and take the patch, the paths are hardcoded in the patch:

* mediawiki-1.15.1/blabla
* mediawiki-1.15.2/blabla

If I want to ignore the mediawiki-1.15... directory, I just need to increment the -p argument. It will be:

```bash
patch -p1 < myfile.patch
```
