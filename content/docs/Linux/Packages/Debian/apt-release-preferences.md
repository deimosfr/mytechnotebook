---
weight: 999
url: "/apt-ajouter-des-preferences-de-release-sur-certains-packages/"
title: "APT: Adding Release Preferences for Specific Packages"
description: "How to configure APT to use packages from different Debian versions by setting package release preferences."
categories: ["Linux", "Debian", "Package Management"]
date: "2013-12-29T17:44:00+02:00"
lastmod: "2013-12-29T17:44:00+02:00"
tags: ["APT", "Debian", "Package Management", "Linux", "Pinning"]
toc: true
---

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Operating System** | Debian 6 |
| **Last Update** | 29/12/2013 |
{{< /table >}}

## Introduction

You may have encountered a situation where you're running the stable version of your system, but you'd like to use a more up-to-date package from testing, for example. This is possible! :-)

## Configuration

### source.list

I've chosen backuppc as an example. Let's say I'm using the stable version, and I'd like to permanently install the testing version. I need to edit my `/etc/apt/source.list` file so it contains what's necessary to download from both stable and testing repositories:

```bash
# Stable
deb http://ftp.fr.debian.org/debian/ stable main non-free contrib
deb-src http://ftp.fr.debian.org/debian/ stable main non-free contrib

# Testing
deb http://ftp.fr.debian.org/debian/ testing main non-free contrib
deb-src http://ftp.fr.debian.org/debian/ testing main non-free contrib

deb http://security.debian.org/ stable/updates main contrib non-free
deb-src http://security.debian.org/ stable/updates main contrib non-free
```

Then we update:

```bash
apt-get update
```

## Preferences

Now, we'll create a `/etc/apt/preferences` file and fill it as follows:

```bash
Package: *
Pin: release a=stable
Pin-priority: 900

Package: *
Pin: release a=testing
Pin-priority: 100
```

Let me try to be clear:

- Pin: 'Package \*' with Pin must be present to indicate each Debian version you want to use (here stable and testing)
- Pin-priority: priorities for Debian versions range from 1 to 1000. The highest value takes precedence over others. So here stable (900) is stronger than testing (500).

That was for the essential part. Now to add backuppc, we'll add these lines:

```bash
Package: backuppc
Pin: release a=testing
Pin-priority: 1001
```

Again, a brief explanation:

- Package: we indicate the name of the package to update in a certain version
- Pin: I'm indicating that I want to switch to testing for the desired package
- Pin-priority: this must be a number above 1000 to override previous restrictions.

And now to verify, if I do:

```bash
apt-get install backuppc
```

I'm offered the testing version, without being offered the rest of the system in testing :-). Here's some useful information:

{{< table "table-hover table-striped" >}}
| Pin Priority | Effect on the package |
|-------------|----------------------|
| 1001 | Install the package even if it's a downgrade |
| 990 | Default for target version archive |
| 500 | Default for normal archive |
| 100 | Default for non-automatic archive but with automatic upgrades |
| 100 | Used for installed package |
| 1 | Default for non-automatic archive |
| -1 | Never install the package even if recommended |
{{< /table >}}

### Installing a specific package without using the preferences file

It's possible to install a specific package without explicitly configuring it. For example, for Nginx, here's how to get the list of available packages:

```bash
apt-cache policy nginx
nginx:
  Installé : (aucun)
  Candidat : 1.2.1-2.2+wheezy2
 Table de version :
     1.4.4-2 0
        100 http://ftp.fr.debian.org/debian/ unstable/main amd64 Packages
     1.2.1-2.2+wheezy2 0
        500 http://ftp.fr.debian.org/debian/ wheezy/main amd64 Packages
```

I can see the unstable and wheezy versions. If I want to install the unstable version, I'll use the '-t' option:

```bash
aptitude install -t unstable nginx
```

### Blocking updates for a specific package

I never clearly remember how to block a Debian package update. Yet it's quite simple if you follow the documentation!

I continue to use and appreciate swiftfox but I have an issue with the latest version (2.0.0.9-1) which doesn't work (problem loading libXcomposite.so.1). I don't have the patience to search for a solution, so I decided to stay with the previous version I have installed: 2.0.0.6-1. For this, it's very simple, just add the following lines to /etc/apt/preferences:

```bash
Package: swiftfox-athlon64
Pin: version 2.0.0.6-1
Pin-priority: 1001
```

The priority 1001 means that the package will never be updated, which is exactly what I want! We can verify this has been taken into account in two ways:

- By trying to update (apt-get upgrade). We shouldn't see the swiftfox package.
- By using apt-cache policy swiftfox-athlon64

Here's the output of this last command:

```bash
 > apt-cache policy swiftfox-athlon64
 swiftfox-athlon64:
 Installé : 2.0.0.6-1
 Candidat : 2.0.0.6-1
 Étiquette de paquet : 2.0.0.6-1
 Table de version :
     2.0.0.9-1 1001
        500 http://getswiftfox.com unstable/non-free Packages
 *** 2.0.0.6-1 1001
        100 /var/lib/dpkg/status
```

## FAQ

### APT: apt-get error: "E: Dynamic MMap ran out of room"

This message may appear during an "apt-get update" when apt no longer has enough space for its cache.
To fix this problem, simply create a file `/etc/apt/apt.conf` and add the following line:

```bash
APT::Cache-Limit "10000000"
```

Run a quick:

```bash
apt-get update
```

And you're done!

## Resources
- [https://www.debian.org/doc/manuals/debian-reference/ch02.en.html](https://www.debian.org/doc/manuals/debian-reference/ch02.en.html)
