---
weight: 999
url: "/Ganeti_\\:_Management_d'un_Cluster_Xen/"
title: "Ganeti: Xen Cluster Management"
description: "Ganeti is a virtual server management tool built on top of Xen virtual machine monitor for easy cluster management of virtual servers."
categories: ["Linux", "Virtualization", "Clusters"]
date: "2009-11-27T21:15:00+02:00"
lastmod: "2009-11-27T21:15:00+02:00"
tags: ["Ganeti", "Xen", "Cluster", "Virtualization"]
toc: true
---

## Introduction

[Ganeti](https://code.google.com/p/ganeti/) is a virtual server management software tool built on top of Xen virtual machine monitor and other Open Source software.

However, Ganeti requires pre-installed virtualization software on your servers in order to function. Once installed, the tool will take over the management part of the virtual instances (Xen DomU), e.g. disk creation management, operating system installation for these instances (in co-operation with OS-specific install scripts), and startup, shutdown, failover between physical systems. It has been designed to facilitate cluster management of virtual servers and to provide fast and simple recovery after physical failures using commodity hardware.

## Documentation

[Documentation on Xen Cluster Management With Ganeti](/pdf/xen_cluster_management_with_ganeti.pdf)  
[Xen Cluster Management With Ganeti On Debian Lenny](/pdf/xen_cluster_management_with_ganeti_on_debian_lenny.pdf)
