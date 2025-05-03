---
weight: 999
url: "/apt-aptitude-les-commandes-utiles/"
title: "Apt & Aptitude: Useful Commands"
description: "A comparison of Apt and Aptitude package management tools in Debian-based Linux distributions, highlighting the advantages of Aptitude"
categories: ["Linux", "Debian", "System Administration"]
date: "2013-05-07T11:24:00+02:00"
lastmod: "2013-05-07T11:24:00+02:00"
tags: ["Apt", "Aptitude", "Debian", "Ubuntu", "Package Management"]
toc: true
---

## Introduction

Being an Ubuntu/Debian user (yes, I use and advocate both), I have fallen in love with the Advanced Packaging Tool, also known as apt. Before Ubuntu, I played in the world of RPM hell, with distros such as Red Hat itself, Mandrake (as it was called back then), and even SuSE. I would find some piece of software, try to install it, only to find that it would choke, saying that it relied on some certain dependencies. I would install the dependencies, only to find conflicting versions with newer software. Hell indeed. So when I discovered the Debian way of installing software, I wondered why no one had mentioned it to me before. It was heaven. This is the way to software, I thought.

## Apt

So, as any new user to the world of apt learns, apt-get is the way to install software in your system. After working on a Debian-based system that uses apt, such as Ubuntu, you also learn the various tools:

* apt-get: Installing and removing packages from your system, as well as updating package lists and upgrading the software itself.
* apt-cache: Search for packages in the package list maintained by apt on the local system
* dpkg- Used for various administrative tasks to your system, such as reconfiguring Xorg.

Those are probably the first few tools that you learn while on a Debian-based distro, if you plan on getting down and dirty at any length. But the buck doesn't stop there. You need to memorize, and learn other tools, if you are to further administrate your system. These include:

* apt-listbugs: See what bugs are listed on a software package before you install it.
* apt-listchanges: Same thing as apt-listbugs, but for non-bug changes.
* apt-rdepends: Tool for viewing dependency trees on packages.
* deborphan- Look for orphaned dependencies on the system left from removing parent packages.
* debfoster- Helps deborphan identify what package dependencies you no longer need on your system.
* dselect- Curses interface for viewing, selecting and searching for packages on your system.
* apt-show-versions -b: show which package comes from which Debian version

There's even more: apt-cdrom, apt-config, apt-extracttemplates, apt-ftparchive, apt-key, apt-mark and apt-sortpkgs.

If any of you have noticed, that is 16 different tools that you need to become familiar with, if you are to start learning about your Debian-based distro. I don't know about you, but doesn't that seem a bit bass-ackwards? I mean, when I'm using OpenSSH, for example, other than scp, all of the functionality of OpenSSH is filed under one tool: ssh. So, wouldn't you think that all the functionality of apt would be under one tool, namely just 'apt'?

Further more, apt-get has a big problem that hasn't really been addressed until only just recently. The problem is in removing packages. You see, apt-get does a great job of indentifying what dependencies need to be installed when you want a certain package, but it fails miserably when you want to remove that package. If dependencies were required, 'apt-get remove' will remove your packages, but leave orphaned dependencies on your system. Psychocats.net has a great writeup on this very phenomenon, by simply installing and removing the package kword. The solution? Aptitude.

Now, before I continue, I want to say that yes, I am aware of 'apt-get autoremove' finally being able to handle orphaned dependencies. This is a step in the right direction, for sure. However, apt-get, with its many other tools, is an okay way of doing things, if you like to learn 16 different tools. Aptitude, as I will show you, is one tool for them all.

## Aptitude

Aptitude is the superior way to install, remove, upgrade, and otherwise administer packages on you system with apt. For one, since it's inception, aptitude has been solving orphaned dependencies. Second, it has a curses interface that blows the doors off of dselect. Finally, and most importantly, it takes advantage of one tool, doing many many functions. Let's take a look:

* aptitude: Running it with no arguments brings up a beautiful interface to search, navigate, install, update and otherwise administer packages.
* aptitude install: Installing software for your system, installing needed dependencies as well.
* aptitude remove: Removing packages as well as orphaned dependencies.
* aptitude purge: Removing packages and orphaned dependencies as well as any configuration files left behind.
* aptitude search: Search for packages in the local apt package lists.
* aptitude update: Update the local packages lists.
* aptitude upgrade: Upgrade any installed packages that have been updated.
* aptitude clean: Delete any downloaded files necessary for installing the software on your system.
* aptitude dist-upgrade: Upgrade packages, even if it means uninstalling certain packages.
* aptitude show: Show details about a package name.
* aptitude autoclean: Delete only out-of-date packages, but keep current ones.
* aptitude hold: Fix a package at it's current version, and don't update it

Are we starting to see a pattern here? One command with different readable options (no unnecessary flags). And that's just the tip of the ice berg. It gets better. For example, when searching for a package using aptitude, the results are sorted alphabetically (gee, imagine that) and justified in column width format. Heck, it will even tell you which one you have installed on your system already, instead of haphazardly listing the packages in some random, unreadable format, like apt-cache.

I've already mentioned it, but aptitude run with no options will pull up a curses application for you to navigate your apt system. If any of you have used it, you know that it is far superior to dselect- talk about a shoddy application. Aptitude makes searching for packages, updating them, removing them, getting details and other necessary tools, easy. Spend 20 minutes inside the console, and you begin to feel like this is an application done right. Spend 20 minutes in dselect, and you'll begin to get massive headaches, and feel lost inside Pan's Labyrinth.

Aptitude is just superior to apt-get in every way, shape, and form. Better dependency handling. Better curses application. Better options. ONE tool. Better stdout formatting. The list goes on and on. I see constantly, on forums, IRC and email, the use of apt-get. We need to better educate our brethren and sisters about the proper use of tools, and show them the enlightened way of aptitude. I've been using aptitude since I first learned about it, ad will continue to do so the remainder of my Debian/Ubuntu days.
