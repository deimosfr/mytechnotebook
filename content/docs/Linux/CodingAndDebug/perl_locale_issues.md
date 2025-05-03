---
weight: 999
url: "/Problèmes_de_locales_avec_Perl/"
title: "Perl Locale Issues"
description: "How to solve locale issues when using Perl on Debian and Ubuntu systems."
categories: ["Debian", "Linux", "Ubuntu"]
date: "2008-11-19T07:14:00+02:00"
lastmod: "2008-11-19T07:14:00+02:00"
tags: ["Perl", "Locale", "Configuration", "Linux", "Ubuntu", "Debian"]
toc: true
---

## Introduction

Perl is great. However, error messages are not something we enjoy. If your environment variables are misconfigured or nothing is defined at the system level, you may encounter problems when launching Perl.

## Problem

Here is what you might encounter when launching Perl:

```bash
perl: warning: Setting locale failed.
perl: warning: Please check that your locale settings:
	LANGUAGE = (unset),
	LC_ALL = (unset),
	LANG = "fr_FR@euro"
    are supported and installed on your system.
perl: warning: Falling back to the standard locale ("C").
locale: Cannot set LC_CTYPE to default locale: No such file or directory
locale: Cannot set LC_MESSAGES to default locale: No such file or directory
locale: Cannot set LC_ALL to default locale: No such file or directory
```

## Solution

### Debian

- To solve the problem, make sure that "fr_FR@euro" is included in your system locales:

```bash
dpkg-reconfigure locales
```

- If it still doesn't work, here's what you can do:

```bash
dpkg-reconfigure console-data
```

When it asks for the default locale, set it to "none".

- Then, for your shell, you need to set some minimum configuration:

```bash
export LANGUAGE=fr_FR@euro
export LC_ALL=fr_FR@euro
export LANG=fr_FR@euro
```

### Ubuntu

#### Solution 1

Execute this command:

```bash
sudo locale-gen fr_FR@euro
```

#### Solution 2

Add this to the file `/var/lib/locales/supported.d/local`:

```bash
fr_FR.UTF-8 @euro
```

For the rest, follow the Debian method described above.
