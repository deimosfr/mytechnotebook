---
weight: 999
url: "/Debian_\\:_Fini_les_erreurs_de_dépendances_quand_vous_voulez_configurer_des_sources/"
title: "Debian: No More Dependency Errors When Configuring From Source"
description: "Learn how to use auto-apt in Debian to automatically resolve missing dependencies when compiling from source."
categories: ["Linux", "Debian"]
date: "2006-11-28T07:48:00+02:00"
lastmod: "2006-11-28T07:48:00+02:00"
tags: ["Dependencies", "Debian", "Source Code", "Compilation", "auto-apt"]
toc: true
---

Imagine you want to install the latest version of xyz from source on your Debian system. You run ./configure and... BLAM!!! You get a ton of errors because of libraries that aren't installed.

Well, this won't happen anymore, thanks to a very user-friendly utility called auto-apt.

To install it:

```bash
apt-get install auto-apt
```

Then, enjoy the pleasure of ./configure without errors. To use it, just type:

```bash
auto-apt run ./configure
```

instead. That's all. If apt detects an error due to a missing file, it will kindly offer to install it for you.

Just remember to update its databases regularly with:

```bash
auto-apt update
auto-apt updatedb
auto-apt update-local
```

Note: auto-apt works with any command that might need missing files: auto-apt run command.
