---
weight: 999
url: "/Gentoo_\\:_Utilisation_des_portages/"
title: "Gentoo: Using Portage"
description: "A guide to effectively use Gentoo's Portage system for package management, including installation, updates, and best practices."
categories: ["Linux", "Gentoo", "Package Management"]
date: "2007-02-24T14:40:00+02:00"
lastmod: "2007-02-24T14:40:00+02:00"
tags: ["Gentoo", "Linux", "Portage", "Emerge", "Package Management"]
toc: true
---

## Introduction

The Gentoo Portage system is an excellent resource (in terms of software) when used correctly. However, incorrect use can lead to a bloated system with unidentifiable packages and files that cannot be updated. This guide will help users properly manage Gentoo portages to have a better system.

## Installing packages

First, what do we mean by "Emerge"? The description of the **emerge** command (available in English [here](https://gentoo-wiki.com/MAN_emerge)), indicates that **emerge** is the command-line program that serves as the interface to the Portage system. This command allows packages to be installed on the system.
This installation includes (these steps are fully automated):

- finding dependencies for the package in question
- installing and/or updating dependencies if necessary
- installing the package in question

After "emerging" a package, the package is integrated into the system, which can use it directly.

## The Portage tree

Portage has a large database of packages it can install. This database is called **the Portage tree**.
This tree is stored on the hard drive, usually in the `/usr/portage/` directory.
To benefit from the latest packages, it is necessary to keep the tree up-to-date by synchronizing it with the official Gentoo tree, which is updated hundreds of times per day and stored on dedicated servers.
This operation is done very simply on a Gentoo system connected to the Internet by typing:

```bash
# emerge --sync
```

It is unnecessary and not recommended to update the Portage tree more than once a day.

## Choosing the right mirrors

To update its tree or obtain the sources of packages to install, Portage must download files from one of its servers.
There are many mirrors containing these files, so it's better to choose one that is fastest for your geographic area.
This choice can be made automatically using the **mirrorselect** program.

To install mirrorselect:

```bash
emerge mirrorselect
```

Once the installation is complete, to choose the 4 best mirrors for downloading sources:

```bash
mirrorselect -D -s4 -t5
```

You can optionally specify the -D option for a more accurate evaluation of mirror performance, but the operation will take longer.
The `/etc/make.conf` file is automatically updated to take into account the chosen mirrors.

By running **mirrorselect -ir**, you can also choose the geographic area of servers to contact to synchronize the Portage tree. However, French mirrors are often not very performant, so it's best to leave the default option (SYNC="rsync://rsync.gentoo.org/gentoo-portage" in the `/etc/make.conf` file).

## ACCEPT_KEYWORDS

The "emerge" function makes it very easy to manage the installation of stable and unstable packages (newer versions of the same software but insufficiently tested). A simple (even simplistic) method to install an unstable package (let's take vlc as an example) for an x86 architecture would be:

```bash
ACCEPT_KEYWORDS="~x86" emerge vlc
```

Unfortunately, this simple method only allows installing the unstable package (vlc) temporarily. During the next system update, the emerge -u world command will try to replace the unstable version of the package with its stable version.
A cleaner method is to indicate in the `/etc/portage/package.keywords` file that we want to use the unstable version:

```bash
echo media-video/vlc >> /etc/portage/package.keywords
```

Thus, each time we want to emerge vlc (during an update for example), the ACCEPT_KEYWORDS="~x86" command will be implied.

## Masked packages

Various situations may lead us to mask (or unmask) certain packages. Similarly to above, you need to use the `/etc/portage/package.mask` file to do this.

```bash
echo x11-base/xfree >> /etc/portage/package.mask
```

And to unmask a masked package, you obviously do:

```bash
echo media-video/realone >> /etc/portage/package.unmask
```

## USE variable

The best way to affect compilation options for packages that interest us is to assign each package (as they are installed) a line in the `/etc/portage/package.use` file as follows:

```bash
echo net-p2p/bittorrent -X >> /etc/portage/package.use
```

## Package maintenance

When you want to update your system with the 'emerge -u world' command, it is possible that Portage might suggest replacing an application with an older version.

## About the 'world' file

First, a brief explanation of what the 'world' file is. The 'world' file lists all packages that the user wants to be able to automatically update with Portage (using the emerge -u world command).
For example, by running:

```bash
emerge gnome
```

Portage records the corresponding package (here gnome-base/gnome) in the 'world' file - **only** gnome, not its dependencies.
This is why, when you run **emerge -u** world, only packages in the world file are updated, not the dependencies. To update dependencies as well, use the --deep option, or -D for short.

If you lose your 'world' file, Portage will no longer know which packages you want to update, which is generally considered problematic.

To **try** to recover your 'world' file, run:

```bash
regenworld
```

The world file resides in **`/var/lib/portage/world`**. You can check its content but don't edit it manually, otherwise Portage could behave erratically.

Also take a look at [this thread](https://forums.gentoo.org/viewtopic.php?t=136627) which gives instructions for recovering your 'world' file if regenworld doesn't work.

## Updating installed software

Using this method, you'll ensure that the 'emerge -u world' command does its job correctly, and you'll have a perfectly configured machine.

```bash
emerge -uDavt world
```

This is the best way to update your Gentoo system.
-u for upgrade (update to the latest packages)
-D for deep (update dependencies)
-a for ask (displays the list of packages to update and asks before compiling)
-v for verbose (displays maximum information, e.g., USE flags used for each package)
-t for tree (displays packages as a dependency tree).

I recommend trying [ecatmur's Cruft script](https://forums.gentoo.org/viewtopic.php?t=152618) or [hepta_sean's more recent 'findcruft' script](https://forums.gentoo.org/viewtopic.php?t=254197) to keep your system clean and tidy!
