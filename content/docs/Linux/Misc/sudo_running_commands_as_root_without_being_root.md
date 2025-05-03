---
weight: 999
url: "/Sudo_\\:_Exécuter_des_commandes_en_root_sans_l'être/"
title: "Sudo: Running commands as root without being root"
description: "Learn how to use sudo to execute commands with root privileges without logging in as root, including installation, configuration and usage examples."
categories: ["Linux"]
date: "2009-09-20T18:01:00+02:00"
lastmod: "2009-09-20T18:01:00+02:00"
tags: ["sudo", "security", "linux", "permissions", "administration"]
toc: true
---

## Introduction

Sudo is frequently used and very practical because it allows occasional execution of commands as root without being logged in as root. It has several security options for usage.

## Installation

It's super simple as usual:

```bash
apt-get install sudo
```

## Configuration

Edit the `/etc/sudoers` file and adapt according to your needs...

### Give all rights to a person

Warning: this operation is equivalent to giving all root rights to a person. They will be able to change the root password, see everything, delete everything. If you wish to apply this type of rights, add this:

```bash
username    ALL=(ALL) ALL
```

Replace username with the name of the user in question.

### Allow only one application to be run as root

If, for example, a user needs to be root to perform a recurring task, which is (most of the time) a script that will run in the background, there should not be a password request, otherwise it cannot execute. This is why you should do it like this:

```bash
username ALL=NOPASSWD: /my/script.sh
```

Put the username at the beginning, then the script or command that this user will have the right to run. You can even put arguments to force the user to use this command only with certain arguments.

### Multiple authorizations

To combine authorizations, simply put a comma:

```bash
username ALL=(ALL) ALL, NOPASSWD: /my/script.sh
```

## Usage

Once a person has the rights, just use sudo, followed by the command:

```bash
username@machine $ sudo /my/script.sh
```

## Get the list of available commands

Here I'm checking the properties for user pmavro:

```bash
$ sudo -l
[sudo] password for pmavro: 
User pmavro may run the following commands on this host:
    (ALL) ALL
```
