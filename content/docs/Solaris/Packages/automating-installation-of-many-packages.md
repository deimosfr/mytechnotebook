---
weight: 999
url: "/Automatiser_l'installation_de_beaucoup_de_packages/"
title: "Automating Installation of Many Packages"
description: "Learn how to automate the installation of numerous packages on Solaris systems using response and admin files."
categories: ["Solaris", "System Administration", "Automation"]
date: "2012-02-15T18:07:00+02:00"
lastmod: "2012-02-15T18:07:00+02:00"
tags: ["Solaris", "pkgadd", "Automation", "Packages", "Installation"]
toc: true
---

## Introduction

Pkgadd isn't particularly user-friendly; it has plenty of options but they're a bit too hidden for my taste. For work, I needed to deploy approximately 1000 packages on several machines. However, with each package installation, I had to enter 'y' and press enter, which quickly becomes tedious.

To solve this problem, I looked at the man page and found the `-n` argument, which is useful since it lets you bypass interactive mode. The drawback is that it doesn't work perfectly. So I searched a bit on the internet and found the response file. This is a file where we pass options. These options are tested during a package installation, and here we force them by default to avoid being bothered.

_Note: there's also the pkgask command that allows you to create a response file, but it seems to have some limitations compared to this method._

**IMPORTANT: This doesn't work for .pkg files! [Convert them beforehand]({{< ref "docs/Solaris/Packages/solaris_package_management.md#pkgtrans">}}).**

## Configuration

### Response File

A response file can be generated using the **pkgask** command:

```bash
pkgask -r /answer -d . <package>
```

This will generate a response file called "answer". This file contains options specific to the package and cannot be used by any other package.

When trying to generate a response file, the command might return this:

```
pkgask: ERROR: package does not contain an interactive request script

Processing of request script failed.
No changes were made to the system.
```

This indicates that the package doesn't contain any customization information, and in this case, using an admin file is necessary.

### Admin File

Not all packages allow the use of a response file, which is why we can also use an admin file.
The admin file works the same way as the response file but uses standard options for a package installation rather than custom options specific to the package.

A template admin file can be found in the directory **/var/sadm/install/admin/default**.
Copy it, then modify it to meet your needs.

Here's what an admin file looks like:

```
# /answer
mail=
instance=overwrite
partial=nocheck
runlevel=nocheck
idepend=nocheck
rdepend=nocheck
space=nocheck
setuid=nocheck
conflict=nocheck
action=nocheck
networktimeout=60
networkretries=3
authentication=quit
keystore=/var/sadm/security
proxy=
basedir=default
```

**WARNING:** If you want to use this method to install packages in the "Finish" script of a jumpstart, don't forget to modify the parameter **"basedir=default"** to **"basedir=/a/"** which corresponds to the mount point of the partition during installation.

### Installation Script

You can simply navigate to the folder in question and execute pkgadd:

- With a response file

```bash
pkgadd -n -r /answer -d . *
```

- With an admin file

```bash
pkgadd -n -a /answer -d . *
```

**You must specify the full paths when using the response file!**

Alternatively, you can use this script (I prefer it, but it's up to you):

- With a response file

```bash
for i in *; do
   test -d /export/home/packages_migration/$i && pkgadd -n -r /answer -d . $i
done
```

- With an admin file

```bash
for i in *; do
   test -d /export/home/packages_migration/$i && pkgadd -n -a /answer -d . $i
done
```

Run this and you'll be set!

## Resources
- [https://forums13.itrc.hp.com/service/forums/questionanswer.do?admit=109447627+1253258346531+28353475&threadId=223010](https://forums13.itrc.hp.com/service/forums/questionanswer.do?admit=109447627+1253258346531+28353475&threadId=223010)
