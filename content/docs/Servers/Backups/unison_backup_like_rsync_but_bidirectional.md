---
weight: 999
url: "/Unison_\\:_Sauvegarde_comme_rsync_mais_bidirectionnelle/"
title: "Unison: Backup like rsync but bidirectional"
description: "How to use Unison for bidirectional file synchronization between systems"
categories: ["Linux", "Backup"]
date: "2011-05-13T17:46:00+02:00"
lastmod: "2011-05-13T17:46:00+02:00"
tags: ["Backup", "Synchronization", "File Management"]
toc: true
---

## Introduction

[Unison](https://en.wikipedia.org/wiki/Unison_(file_synchronizer)) is a popular file synchronization software that also offers functionality for creating and managing directory backups. The synchronization is bidirectional (meaning that modifications in one directory are reflected in the other and vice-versa), making it useful for keeping directories in sync between two different machines.

Unison is free software released under the GPL license. It works on a wide range of operating systems (Windows, Linux, Mac OS X), allowing file synchronization between different operating systems.

## Installation

On Debian, installation is very simple once again:

```bash
aptitude install unison
```

## Configuration

I recommend checking the manual for complete information on how it works, but I'll provide a configuration that I use to replicate my website. Rather than using a long command line, I prefer using a configuration file that contains all the elements I want to back up and how to handle conflicts. In `~/.unison`, you can create `*.prf` files. Here's my configuration:

```bash
root = /var/www
root = ssh://192.168.90.1//var/www

ignore = Name w3tc
ignore = Name piwik/tmp
ignore = Name captcha-temp
batch = true
auto = true
silent = true
log = true
logfile = /tmp/unison.log
```

The options used are:

- root: The two sources and destinations to replicate. I use one local host and another over SSH.
- ignore: Allows using restrictions
- batch: We're in automatic mode and want to avoid being asked questions
- auto: Indicates that we'll use unison in an automated way
- silent: Silent prevents any output
- log: We enable logging
- logfile: We indicate where we want the logs (by default ~/unison.log)

## Usage

To use the configuration file we just created, it's very simple:

```bash
unison www.prf
```

You can also use options one after another (in nearly the same format, see the manual) if you don't want to use a configuration file.
