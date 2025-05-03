---
weight: 999
url: "/Compilations_foireuses/"
title: "Failed Compilations"
description: "Solutions for common Linux compilation failures and required libraries"
categories: ["Linux", "Development"]
date: "2006-08-05T19:44:00+02:00"
lastmod: "2006-08-05T19:44:00+02:00"
tags: ["compilation", "gcc", "libraries", "development", "troubleshooting"]
toc: true
---

Another failed compilation? As usual, you're missing a library. I've listed here the minimum packages you should have to avoid problems.

If you get a message like this during compilation:

```bash
checking for C compiler default output... configure: error: C compiler cannot create executables
```

This is typically the kind of error you might encounter. Install the following packages:

```bash
gcc libc6 libc6-dev make autoconf
```

Also consider:

```bash
glibc2 glibc2-dev
```

For Debian users like me:

```bash
apt-get install libc6 libc6-dev make autoconf glibc2 glibc2-dev
```
