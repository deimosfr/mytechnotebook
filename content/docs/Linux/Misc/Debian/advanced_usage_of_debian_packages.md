---
weight: 999
url: "/Utilisation_avancée_des_packages_Debian/"
title: "Advanced Usage of Debian Packages"
description: "A comprehensive guide to advanced Debian package management using dpkg, aptitude, and apt-get tools with practical examples."
categories: ["Debian", "Linux", "Ubuntu"]
date: "2009-12-11T21:54:00+02:00"
lastmod: "2009-12-11T21:54:00+02:00"
tags: ["dpkg", "apt-get", "aptitude", "package management", "linux", "debian", "ubuntu"]
toc: true
---

## Introduction

Quite possibly the most distinguishing feature of Debian-based Linux distributions (such as Ubuntu, Mepis, Knoppix, etc) is their package system – APT. Also known as the Advanced Package Tool, APT was first introduced in Debian 2.1 in 1999. APT is not so much a specific program as it is a collection of separate, related packages.

With APT, Linux gained the ability to install and manage software packages in a much simpler and more efficient way than was previously possible. Before its introduction, most software had to be installed either by manually compiling the source code, or using individual packages with no automatic dependency handling (such as RPM files). This could mean hours of "dependency hell" even to install a fairly trivial program.

In this article, we are going to highlight some of APT's best features, and share a few of the lesser known features of APT and its cousin dpkg.

## Dpkg

The base of Debian's package system is dpkg. It performs all the low level functions of software installation. If you were so inclined, you could use dpkg alone to manage your software. It can install, remove, and provide information on your system's software collection. Here are some of my favorite features.

## Basic installation of local file

Some software authors create Debian packages of their programs, but do not provide a repository for APT to fetch from. In this case they just provide a downloadable .deb file. This is very similar to RPM packages, or even Windows .msi files. It contains all the files and configuration information necessary to install the program. To install a program from a .deb file, you simply need:

```bash
dpkg -i MyNewProgram.deb
```

The -i, as you may guess, tells dpkg to install this piece of software.

## Listing a package's contents

You may find yourself, after installing a program, unable to figure out how to run that program. Sometimes, you need to know where to find the config files for your new game. Dpkg provides an easy way to list all the files that belong to a particular package.

```bash
dpkg -L MyNewProgram
```

Note that case matters. -L and -l are entirely different options.

Often, a package has so many files it can be difficult to sift through the list to find the one(s) you're looking for. If that's the case, we can use grep to filter the results. The following command does the same as above, but only shows results that have "bin" in the path, such as /usr/bin.

```bash
dpkg -L MyNewProgram 
```

I won't even begin to go into the awesome power that is grep, but in its simplest form it can be used, like above, to quickly and easily filter a program's output.

## Discover what package a file belongs to

Occasionally, you find yourself in a situation that's the reverse of the section above. You have a file, but you don't know what package it belongs to. Once again, dpkg has you covered.

```bash
dpkg -S mysteryfile.cfg
```

This will tell you which package created/owns that file.

## Listing what you've got installed

Let's say you're about to reinstall your system, and you want to know exactly what you've already got installed. You could open up an app like Synaptic and set a filter to show everything marked as "installed", or you could do it quickly and easily from the command line with dpkg.

```bash
dpkg -l
```

or

```bash
dpkg --get-selections
```

That will give you one big long list of everything you've got installed. Advanced users could use these commands to create a text file with all their packages listed, which could be fed into APT later to reinstall everything at once!

## Reconfiguring a package

When a .deb package is installed, it goes through a few stages. One of those is the configuration stage, where developers can put a series of actions that take place once all the files have been installed to a proper location. This includes things like start/stopping services, or creating logs, or other such things. Sometimes you need to repeat those steps, without going through the whole reinstallation process. For that, you use:

```bash
dpkg-reconfigure (packagename)
```

This will redo all the post-install steps needed for that package without forcing you to reinstall. Believe me, this one comes in handy.

## Aptitude/Apt-get

There's some debate and confusion regarding these two tools. Many Linux users have a hard time telling when/why to use one over the other, as they do roughly the same thing.

Short answer: use Aptitude.

Long answer: Both can be used to manage all software installations removals, and both will do a good job. The Debian team officially recommends using Aptitude. It's not that it's a LOT better than apt-get, but that it's a little better, in lots of ways. You can use either one and it will meet your package management needs, but don't mix and match on the same system. Pick one and stick with it.

## Finding the right package

I often find myself in need of software to do a certain thing, but I don't know the name of any programs to do it. For example, I may need a FLAC player, but don't know offhand what player will work…

Aptitude:

```bash
aptitude search flac
```

Classic APT:

```bash
apt-cache search flac
```

You'll get a list of available packages that have "flac" in the name or description.

## Preventing a package from updating

On some occasions, I have a version of a package that I want to keep even though there may be upgrades. When it comes to my kernel, for example, I prefer to update manually.

Aptitude:

```bash
aptitude hold (packagename)
```

dpkg:

```bash
echo "(packagename) hold" 
```

## Upgrading

Both Aptitude and classic APT provide two methods of upgrading your system: upgrade and dist-upgrade. This is another thing that causes some confusion. An upgrade is an upgrade, right? Well not exactly.

A regular upgrade will read your list of packages, check online for newer versions, and install them as needed. It will NOT, however, perform any upgrades that would require new packages installed, or existing ones removed. This is what dist-upgrade is for. It will get every newer version it finds, even if it involves installing something new (such as a dependency) or removing an existing package (if it's obsolete or is no longer needed).

Aptitude:

```bash
aptitude safe-upgrade
aptitude full-upgrade
```

Classic APT:

```bash
apt-get upgrade
apt-get dist-upgrade
```

## Learn about a package

Finally, sometimes you just need to know a little about a package. What version is it? Who maintains it? Is it already installed? All these things and more you can find with:

Aptitude:

```bash
aptitude show (packagename)
```

APT:

```bash
apt-cache showpkg (packagename)
```

All of the programs mentioned here are capable of far more than I've shown. The tips here should go a long way in helping you use this amazing package system to its full potential.

Also, for those with the patience to read all the way to the end,

```bash
apt-get moo
```
