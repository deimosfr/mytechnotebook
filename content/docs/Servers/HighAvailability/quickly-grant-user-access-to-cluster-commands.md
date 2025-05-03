---
weight: 999
url: "/Autoriser_rapidement_un_utilisateur_à_avoir_accès_aux_commandes_cluster/"
title: "Quickly Grant User Access to Cluster Commands"
description: "Learn how to quickly configure user permissions to allow non-root users to execute cluster management commands."
categories: ["Linux", "Security", "Cluster"]
date: "2007-05-31T10:12:00+02:00"
lastmod: "2007-05-31T10:12:00+02:00"
tags: ["sudo", "Cluster", "Permissions", "Linux"]
toc: true
---

## Introduction

We often need users to have access to specific commands without being root, and for cluster management, if you have dedicated administrators, it's quite useful. Here's a simple way to give them the necessary permissions...

## Configuration

To give a user permissions to simply use cluster commands, here are the files to modify:

* /etc/sudoers:

```
# Cmnd alias specification
Cmnd_Alias     CLUSTAT          = /usr/sbin/clustat
Cmnd_Alias     CLUSVCADM        = /usr/sbin/clusvcadm
Cmnd_Alias     MOUNT            = /bin/mount
Cmnd_Alias     UMOUNT           = /bin/umount

# Defaults specification

# User privilege specification
root            ALL=(ALL) ALL
my_user          ALL=NOPASSWD:CLUSTAT,NOPASSWD:CLUSVCADM,NOPASSWD:MOUNT,NOPASSWD:UMOUNT
```

* ~/.bashrc (for the user)

```bash
# User specific aliases and functions
alias clustat='sudo /usr/sbin/clustat'
alias clusvcadm='sudo /usr/sbin/clusvcadm'
alias mount='sudo /bin/mount'
alias umount='sudo /bin/umount'
```

## Usage

With my user account, I can simply run the commands and they will be executed as root:

```bash
$ clustat
```
