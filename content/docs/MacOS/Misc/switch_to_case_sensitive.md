---
weight: 999
url: "/switch_mac_os_case_sensitive/"
title: "Switching to Case-Sensitive File System on Mac OS X"
description: "Instructions on how to clear the DNS cache on Mac OS X systems"
categories: ["Mac OS X", "Network"]
tags: ["Mac OS X", "Filesystem"]
toc: true
---

## Introduction

On Mac OS X, the default file system is case-insensitive. However, you can switch to a case-sensitive file system if needed. This is particularly useful for certain applications or development environments that require case sensitivity.

{{< alert context="warning" text="Switching to a case-sensitive file system can cause issues with applications that expect a case-insensitive file system. This is why moving to a case-sensitive file system is not recommended." />}}

Instead, it is preferable to create a dedicated volume (not a partition) with a case-sensitive file system. This way, you can keep your main volume case-insensitive while having a separate volume for applications that require case sensitivity. The other advantage is that you won't encounter any issue with TimeMachine backups and restoring your system (you can't restore a case-sensitive volume to a case-insensitive one without a long and painful procedure).

## Create a new volume

To create a new volume with a case-sensitive file system, follow these steps:

1. Open **Disk Utility** (found in Applications > Utilities).
2. Select your main disk (usually named "Macintosh HD").
3. Click on the + sign next to Volume.

![Mac OS create volume](/images/macos_casesensitive_vol_create.avif)

Choose a name for the new volume (e.g., "workspace") and set it to **APFS (Case-sensitive)** or **APFS (Case-sensitive, Encrypted)** (it's recommended to use the Encrypted for security reasons):

![Mac OS volume naming](/images/macos_casesensitive_vol_naming.avif)

## Use the new volume

Once the new volume is created, you should see it and find relecant information. You can use it for applications or development environments that require a case-sensitive file system.

![Mac OS created volume](/images/macos_casesensitive_vol_created.avif)

From the terminal, create a symbolic link to the new volume:

```bash
ln -s /Volumes/workspace ~/workspace
```

This will create a symbolic link in your home directory that points to the new volume. You can now use `~/workspace` as a path to access the case-sensitive volume. It's very useful for development purposes, especially if you are using a version control system like Git or Mercurial that may have issues with case-insensitive file systems.
