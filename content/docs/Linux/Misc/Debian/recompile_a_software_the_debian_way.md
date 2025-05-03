---
weight: 999
url: "/Recompiler_un_soft_à_la_sauce_Debian/"
title: "Recompile a Software the Debian Way"
description: "Tutorial on how to recompile Debian packages while keeping the Debian package system benefits."
categories: ["Debian", "Linux"]
date: "2010-03-28T20:27:00+02:00"
lastmod: "2010-03-28T20:27:00+02:00"
tags: ["Debian", "Linux", "Package", "Compilation"]
toc: true
---

## Introduction

Sometimes you need to recompile some software. But when using Debian packages, it's not always easy because you want to keep the Debian conveniences. The solution is apt-get source :-)

## Installation

We will need these packages:

```bash
aptitude install dpkg-dev
```

## Example

For example, I want to install nginx but recompile it with an additional option. I download the sources as follows:

```bash
apt-get source nginx
```

Then I have a directory containing the sources and I reconfigure my sources with ./configure.

After that, all that's left is to recreate the package:

```bash
dpkg-buildpackage -us -uc
```

You may be missing some packages that will be indicated afterwards (usually autotools-dev). Install them and restart the above line.

Then install the package like this:

```bash
dpkg -i nginx.deb
```
