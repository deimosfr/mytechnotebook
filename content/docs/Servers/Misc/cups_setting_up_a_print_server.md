---
weight: 999
url: "/Cups_\\:_mise_en_place_d'un_serveur_d'impression/"
title: "CUPS: Setting Up a Print Server"
description: "Guide on how to set up a CUPS print server on Linux systems, including installation, configuration, and administration."
categories: ["Linux"]
date: "2010-04-01T12:58:00+02:00"
lastmod: "2010-04-01T12:58:00+02:00"
tags: ["Linux", "Servers", "Printing", "CUPS", "System Administration"]
toc: true
---

## Introduction

The [Common Unix Printing System (CUPS)](https://fr.wikipedia.org/wiki/Cups) is a modular printing system for Unix and Unix-like operating systems. Any computer using CUPS can act as a print server; it can accept documents sent from other machines (client computers), process them, and send them to the appropriate printer.

## Installation

The installation is very simple:

```bash
aptitude install cupsys cupsys-client cupsys-bsd cupsys-driver-gimpprint samba-client
```

I've included a module for creating PDFs and a package containing drivers.

## Configuration

### Parallel Ports

This technology is now almost obsolete, so it's unlikely that you have a printer operating on a parallel port. If like me you don't need it, disable this option:

```bash
# Cups configure options

# LOAD_LP_MODULE: enable/disable to load "lp" parallel printer driver module
LOAD_LP_MODULE=no
```

### cupsd.conf

Now, let's edit the basic configuration and change it so that we can access it. Here are the lines to replace (in this case my server is 192.168.0.1 in a 192.168.0.0 network):

```bash
...
# Only listen for connections from the local machine.
Listen 192.168.0.1:631
...
# Restrict access to the server...
<Location />
    Order Deny,Allow
    Deny From All
    Allow From 192.168.0.*
    Allow From @LOCAL
</Location>
...
# Restrict access to the admin pages...
<Location /admin>
    AuthType Basic
    AuthClass System
    Order Deny,Allow
    Deny From All
    Allow From 192.168.0.*
    Allow From @LOCAL
</Location>
```

Restart the CUPS service afterward.

You can now access the administration interface via: http://192.168.0.1:631 or https://192.168.0.1:631

### Administering Printers

To authorize a user account to administer printers, simply add it to the lpadmin group:

```bash
adduser username lpadmin
```

## References

[CUPS Documentation](/pdf/cups.pdf)
