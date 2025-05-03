---
weight: 999
url: "/Installation_from_Scratch_d'un_Cisco_Pix/"
title: "Installation from Scratch of a Cisco Pix"
description: "A guide on how to install and configure a Cisco Pix firewall from scratch, including initial setup and telnet activation."
categories: ["Network"]
date: "2007-05-23T14:32:00+02:00"
lastmod: "2007-05-23T14:32:00+02:00"
tags: ["Cisco", "Firewall", "Network", "Security"]
toc: true
---

## Introduction

The installation of a Cisco device is not very complicated, but if you're not familiar with it, it's not always obvious. That's why I made this small guide, since it's not something we do every day.

## Installation

Connected via serial port? Let's get started:

```
Pre-configure PIX Firewall now through interactive prompts [yes]?
Enable password [<use current password>]: mot_de_passe
Clock (UTC):
 Year [2007]:
 Month [May]:
 Day [23]:
 Time [02:57:33]: 12:01:00
Inside IP address: 192.168.0.77
Inside network mask: 255.255.255.0
Host name: hk-pix-bak
Domain name: mon_domaine
IP address of host running PIX Device Manager: 192.168.0.104

The following configuration will be used:
Enable password: mot_de_passe
Clock (UTC): 12:01:00 May 23 2007
Inside IP address: 192.168.0.77
Inside network mask: 255.255.255.0
Host name: hk-pix-bak
Domain name: mon_domaine
IP address of host running PIX Device Manager: 192.168.0.104

Use this configuration and write to flash? y
```

To summarize:
- Inside IP address: this is the address of the Cisco device
- IP address of host running PIX Device Manager: this is the address of the machine that will have the right to connect via HTTPS to configure the PIX

## Activation of Telnet

Here, we want to enable telnet on the **inside** interface for the network **192.168.0.0**:

```
telnet 192.168.0.0 255.255.255.0 inside
```

Now, we can verify that it works:

```
# show run | grep telnet
telnet 192.168.0.0 255.255.255.0 inside
telnet timeout 5
```
