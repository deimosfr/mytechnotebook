---
weight: 999
url: "/XenServer_5.0_\\:_Astuces/"
title: "XenServer 5.0: Tips & Tricks"
description: "A collection of tips and tricks for XenServer 5.0, including solutions for common issues with XenCenter and boot options."
categories: 
  - Linux
date: "2009-05-08T15:25:00+02:00"
lastmod: "2009-05-08T15:25:00+02:00"
tags:
  - Servers
  - Virtualization
  - XenServer
toc: true
---

## Introduction

Xen can be very temperamental, especially the Citrix version with the incomplete XenCenter. That's why, here are some tips to help you navigate through difficult situations.

## Adding Boot Options in XenCenter "Startup Options"

To add this option, which is missing when creating a VM using a template.

You need to access the server console and check the UUID of the concerned VM:

```bash
xe vm-list 
```

Retrieve the UUID preceding the VM: uuid=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxxx

Then run the following command:

```bash
xe vm-param-set HVM-boot-policy=BIOS\ order uuid=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxxx
```

To verify that the command worked correctly:

```bash
xe vm-param-list uuid=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxxx
```

Check that the HVM-boot-policy value is indeed set to 'BIOS order'

It's quite possible that the VM will no longer boot normally with this option activated, or that you simply want to return to the original mode. To deactivate it, use the following line:

```bash
xe vm-param-set HVM-boot-policy= uuid=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxxx
```
