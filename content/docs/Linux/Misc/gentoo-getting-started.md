---
weight: 999
url: "/Gentoo_\\:_Bien_commencer_avec_Gentoo/"
title: "Gentoo: Getting Started"
description: "Learn essential commands and tips for getting started with Gentoo Linux after a fresh installation."
categories: ["Linux", "Distributions", "Gentoo"]
date: "2007-01-18T08:34:00+02:00"
lastmod: "2007-01-18T08:34:00+02:00"
tags: ["Gentoo", "Linux", "Package Management", "Portage"]
toc: true
---

## Introduction

Gentoo is a Linux distribution known as a source-based distribution. It was designed to be modular, portable, and optimized for the user's hardware. As such, all programs must be compiled from source code. However, many software packages available in precompiled form for different architectures can also be used. This is managed through Gentoo's Portage system.

Its particularity is the complete (or partial) compilation of a GNU/Linux system from sources, similar to Linux From Scratch but automated.

Its package management tools are inspired by BSD ports. This process allows for complete optimization and customization of the system but takes some time to compile all the necessary software.

This type of installation allows you to make the best use of your machine's architecture. Indeed, the source code will be compiled taking into account the possible optimizations of the processor's instruction set. Most distributions are compiled with the i386 instruction set and not for a more recent processor, in order to maintain functionality on as many machines as possible. More recent processors then operate minimally without using the manufacturer's optimizations.

In addition, this type of installation makes it easy to manage dependencies, even during a major update of the entire distribution. When installing each program, the development libraries that accompany it are automatically installed, and other programs that use these libraries will be automatically recompiled with the new version of these libraries during the update. The result is a high-performance, coherent, and stable system.

Since Gentoo is a bit special in some ways, I'll note here the essential points that helped me make a clean installation.

## After an Installation

### Authorizing a User to Connect as Root

For security reasons, users can switch to root with su only if they belong to the wheel group. To add a username to the wheel group, type the following command as root:

```bash
gpasswd -a username wheel
```

### Updating the Package List

Use this command:

```bash
emerge --sync
```

### Installing Software

To install software:

```bash
emerge screen
```

Here I'm installing screen.

### Updating Gentoo

#### Software Updates

Here's a command to see what's left to update:

```bash
emerge -Dvp world
```

Now let's update Gentoo:

```bash
emerge world
```

This command allows you to recompile the entire system:

```bash
emerge -e world
```

#### Configuration Updates

If you have a message like:

```
* IMPORTANT: 33 config files in /etc need updating.
```

You can find out which configuration files want to be replaced:

```bash
find /etc -iname '._cfg????_*'
```

For more information, consult the command:

```bash
emerge --help config
```

### Searching for a Package

To search for a package or description:

```bash
emerge --searchdesc searchword
```

**Note:** Currently on Gentoo, it is preferable to install "esearch" or "eix" to make searches. Then, it's used like this:

```bash
eix searchword
```
```bash
esearch searchword
```

### Installing a Specific Package

For example, I'd like to install munin. However, I'm having issues because I'm on the stable version ("x86") and I want to install the version that's in unstable ("**~**x86") because it doesn't exist in stable. When I do:

```bash
emerge munin
```

I get this:

```
Calculating dependencies
!!! All ebuilds that could satisfy "munin" have been masked.
!!! One
