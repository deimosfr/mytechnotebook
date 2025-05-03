---
weight: 999
url: "/Changer_les_locales_de_Solaris/"
title: "Changing Solaris Locales"
description: "How to change locale settings on Solaris systems to fix locale-related error messages and properly configure internationalization."
categories: ["Solaris"]
date: "2010-02-08T16:22:00+02:00"
lastmod: "2010-02-08T16:22:00+02:00"
tags: ["solaris", "locale", "internationalization", "system administration"]
toc: true
---

## Introduction

You may encounter messages like "couldn't set locale correctly" which can quickly become annoying to see in the display.

## Problem Explanation

This occurs because the locales installed on the machine do not match those in your shell's environment variables.

To see what you have in your shell:

```bash
> env 
```

And to see what's available on the system, it's just as simple:

```bash
> ls /usr/lib/locale
C
```

## Solution

For my part, I live in France, so I need the locales for my country. I'm going to install the Western European locales. For this, you'll need the Solaris DVD in the drive:

```bash
cd /cdrom/cdrom0/Solaris_10/Product/
pkgadd -d . SUNWweuos
```

Now it's good, there will be no more error messages.

To change Solaris locales at the system level, edit the `/etc/default/init` file and adapt according to your needs:

```bash
TZ=Europe/Paris
CMASK=022
LC_COLLATE=fr_FR.ISO8859-15
LC_CTYPE=fr_FR.ISO8859-15
LC_MESSAGES=fr
LC_MONETARY=fr_FR.ISO8859-15
LC_NUMERIC=fr_FR.ISO8859-15
LC_TIME=fr_FR.ISO8859-15
LC_COLLATE=fr_FR.ISO8859-15
LC_CTYPE=fr_FR.ISO8859-15
LC_MESSAGES=fr
LC_MONETARY=fr_FR.ISO8859-15
LC_NUMERIC=fr_FR.ISO8859-15
LC_TIME=fr_FR.ISO8859-15
```

## Resources
- [https://developers.sun.com/dev/gadc/faq/locale.html](https://developers.sun.com/dev/gadc/faq/locale.html)
