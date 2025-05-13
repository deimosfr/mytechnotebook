---
title: "Mounting a Filesystem in Multiple Places Simultaneously"
slug: mounting-a-filesystem-in-multiple-places-simultaneously/
description: "How to mount the same filesystem in multiple locations on Linux and BSD systems."
categories: ["Linux"]
date: "2008-10-28T21:04:00+02:00"
lastmod: "2008-10-28T21:04:00+02:00"
tags: ["Commandes", "Commands", "Filesystem", "Mount", "Linux", "BSD"]
---

## Introduction

Some might say: "There's no need for this, just use symbolic links with ln -s". But I disagree - it's really not the same thing. This approach allows you to have a global view of a particular directory.

For example, if I want to mount `/var/jails` in `/jails`, it's possible and here's the result once done:

```bash
Filesystem     Size    Used   Avail Capacity  Mounted on
/dev/ad4s1a    989M    129M    782M    14%    /
devfs          1.0K    1.0K      0B   100%    /dev
/dev/ad4s1d    224G    1.4G    204G     1%    /usr
/dev/ad4s1e    224G    4.0M    206G     0%    /var
/var/jails     224G    4.0M    206G     0%    /jails
```

## Commands

=== "Linux"

    ```bash
    mount --bind /folder1 /folder2
    ```
=== "BSD"

    ```bash
    mount -t nullfs /folder1 /folder2
    ```

