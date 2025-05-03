---
weight: 999
url: "/Restaurer_les_permissions_d\\'une_Red_Hat/"
title: "Restore permissions on Red Hat"
description: "How to restore file permissions on a Red Hat system after a permissions mistake"
categories: ["Linux", "Red Hat"]
date: "2012-02-04T17:50:00+02:00"
lastmod: "2012-02-04T17:50:00+02:00"
tags: ["permissions", "rpm", "recovery", "system administration"]
toc: true
---

## Introduction

A colleague of mine made a serious mistake (running `chown -Rf mysql` on `/`!). This caused a huge mess, and we had to find a solution to restore the correct permissions.

Fortunately, Red Hat anticipated these kinds of errors and included the `--setperms` and `--setugids` options in the `rpm` command to repair permissions on installed packages. Basically, this gives you a way to repair your machine.

So if you also made a mistake like this, know that there is a solution on Red Hat.

## Usage

Here are the two magic commands:

```bash
for u in $(rpm -qa); do rpm --setugids $u; done
for p in $(rpm -qa); do rpm --setperms $p; done
```

You'll need to do a bit of verification afterward because this only repairs the permissions of files and directories contained in installed packages. Your personal files will not have their permissions restored with this method.

## Resources
- http://www.adminlinux.org/2009/07/how-to-restore-default-system.html
