---
weight: 999
url: "/Pkg-get_\\:_Mise_en_place_d'un_système_de_repository_pour_Solaris/"
title: "Pkg-get: Setting up a Repository System for Solaris"
description: "This guide explains how to set up a repository system using pkg-get for Solaris to easily install packages."
categories: ["Solaris", "Linux"]
date: "2013-01-23T10:54:00+02:00"
lastmod: "2013-01-23T10:54:00+02:00"
tags: ["Solaris", "Repository", "Package Management", "pkgutil", "pkg-get"]
toc: true
---

## Introduction

Package installation in Solaris is quite basic by default, and recompiling sources is not always simple or fast. This is why I suggest using pkg-get, which is a very practical utility that allows you to install your desired packages for Solaris (SPARC or x86) in just a few seconds or minutes.

## Installation

### New Method

The new method makes things much simpler:

```bash
pkgadd -d http://get.opencsw.org/now
```

Then add `/opt/csw/bin` to your path:

```bash
export PATH=$PATH:/opt/csw/bin
```

And that's it :-)

### Old Method

First, you need wget! I suggest downloading it from these URLs:

- x86: http://www.blastwave.org/wget-i386.bin
- SPARC: http://www.blastwave.org/wget-sparc.bin

Then place it in the `/usr/bin/` directory of your Solaris. Rename it to wget to simplify the task:

```bash
mv /usr/bin/wget-i386.bin /usr/bin/wget
chmod 755 /usr/bin/wget
```

Next, add `/opt/csw/bin` to your path:

```bash
export PATH=$PATH:/opt/csw/bin
```

This is obviously a temporary solution for your PATH, but I recommend adding it to the `/etc/profile` file:

```bash
echo "export PATH=$PATH:/opt/csw/bin" >> /etc/profile
```

You should also have gzip installed (normally it's included by default, but I prefer to specify it...). Now that wget is installed, we just need to download pkg_get:

```bash
wget http://www.blastwave.org/pkg_get-3.8-all-CSW.pkg
```

And let's proceed with the installation:

```bash
pkgadd -d pkg_get-3.8-all-CSW.pkg
```

In response to the questions, answer **all** and **yes** every time :-).

## Configuration

Let's edit the `/opt/csw/etc/pkg-get.conf` file to select the most appropriate mirror from [this list](https://www.blastwave.org/mirrors.php).

## Usage

Now that everything is set up, we can use it. First, update the list of available packages using pkg-get or pkgutil (depending on the version of opencsw you have):

```bash
pkg-get -U
```

or

```bash
pkgutil -U
```

To search (for example, for vim), use this command:

```bash
pkg-get -a 
```

And then to install it:

```bash
pkg-get install vim
```

or

```bash
pkgutil -i vim
```

It's not complicated and saves a tremendous amount of time :-)
