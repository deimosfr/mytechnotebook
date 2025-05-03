---
weight: 999
url: "/OCS_Inventory_\\:_Mise_en_place_d'un_inventaire_de_parc_automatique/"
title: "OCS Inventory: Setting up Automatic Network Inventory"
description: "How to set up OCS Inventory for automatic network inventory management with client installation on Debian and Windows platforms"
categories: ["Linux", "Debian", "Windows"]
date: "2008-12-30T10:08:00+02:00"
lastmod: "2008-12-30T10:08:00+02:00"
tags: ["Inventory", "OCS", "Network", "System Administration", "Servers"]
toc: true
---

## Introduction

Open Computer and Software Inventory Next Generation is an application designed to help system or network administrators keep track of the network machines' configuration and the software installed on them.

OCS Inventory is also able to detect any active device on the network, such as switches, routers, printers, and other unexpected hardware. For each one, it stores MAC and IP addresses and allows you to classify them.

If the administration server runs on Linux, and if nmap and smblookup are available, you also have the option to scan an IP or a subnet for detailed information about uninventoried hosts.

Last but not least, OCS Inventory NG integrates package deployment features on client machines. From the administration console, you can upload packages (software installations, commands, or just files to store on client computers) that will be downloaded via HTTP/HTTPS and executed by agents on clients.

Here's a setup document:  
[OCS Inventory Setup](/pdf/mise_en_place_d'ocs_inventory.pdf)

## Client Installation

### Debian

For Debian, here's what you need to do to make it work. Extract the archive:

```bash
tar -xzvf OCSNG_LINUX_AGENT_1.01_with_require.tar.gz
```

Then install the perl modules via cpan:

```bash
cpan
install XML::Simple
install Compress::Zlib
install Net::IP
install LWP
quit
```

Then the others via apt:

```bash
apt-get install libmd5-perl libnet-ssleay-perl
sh setup.sh
```

### Windows

For your domain, here is a small script that can be executed at user logon:

```bash
@echo off
"\\mon_serveur_de_domaine\netlogon\192.168.0.16.exe" /DEBUG /NP /INSTALL
c:
cd "%ProgramFiles%\OCS Inventory Agent"
OCSInventory /SERVER:192.168.0.16 /TAG:"%username%"
```

Here 192.168.0.16 corresponds to my OCS server.

## Resources
- [OCS Inventory Setup with GLPI](/pdf/mise_en_place_d'ocs_inventory.pdf)
