---
weight: 999
url: "/AutoFsck_\\:_Changer_les_checks_filesystem_sur_Ubuntu/"
title: "AutoFsck: Changing Filesystem Checks on Ubuntu"
description: "Learn how to modify filesystem check behavior in Ubuntu using AutoFsck to make them run during shutdown instead of boot time."
categories: ["Linux", "Ubuntu", "System"]
date: "2007-08-25T21:57:00+02:00"
lastmod: "2007-08-25T21:57:00+02:00"
tags: ["Ubuntu", "Filesystem", "fsck", "Maintenance"]
toc: true
---

If you've used Ubuntu Linux for longer than a month, you've no doubt realized that every 30 times you boot up you are forced to run a filesystem check. This filesystem check is necessary in order to keep your filesystem healthy. Some people advise turning the check off completely, but that is generally not a recommended solution. Another solution is to increase the number of maximum mounts from 30 to some larger number like 100. That way it's about 3 times less annoying. But this solution is also not recommended. Enter AutoFsck.

AutoFsck is a set of scripts that replaces the file system check script that comes shipped with Ubuntu. The difference is that AutoFsck doesn't ruin your day if you are so unfortunate to encounter the 30th mount. The most important difference is that AutoFsck does its dirty work when you shut your computer down, not during boot when you need your computer the most!

The 30th time you mount your filesystem, AutoFsck will wait until you shut down your computer. It will then ask you if it is convenient for you to check your filesystem. If it is convenient for you, then AutoFsck will restart your computer, automatically execute the filesystem check, and then immediately power down your system. If it is not convenient for you to check your filesystem at that moment, then AutoFsck will wait until the next time you shut down your computer to ask you again. Being prompted for a file system check during shutdown is infinitely more convenient than being forced to sit through a 15 minute check during boot up.

[https://wiki.ubuntu.com/AutoFsck](https://wiki.ubuntu.com/AutoFsck)
