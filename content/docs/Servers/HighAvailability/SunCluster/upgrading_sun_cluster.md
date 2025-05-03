---
weight: 999
url: "/Upgrader_SUN_Cluster/"
title: "Upgrading SUN Cluster"
description: "Guide on how to upgrade a SUN Cluster with detailed step-by-step instructions"
categories: ["Solaris", "Servers"]
date: "2010-02-13T14:37:00+02:00"
lastmod: "2010-02-13T14:37:00+02:00"
tags: ["SUN", "Cluster", "Upgrade", "Solaris"]
toc: true
---

## 1. Introduction

Upgrading a cluster is not an easy task. I'll explain here the steps to follow.

## 2. Preparing environment

First of all, you need to be sure you don't have any running services (You should do those step for all your servers).

Then you need to boot in a non cluster mode:

```bash
reboot -- -x
```

You also need to download the latest version of Sun Cluster and unzip it in `/export/home/patchs/` for example.

## 3. Upgrade

### 3.1 Normal way

Now start installer (here: `/export/home/patchs/Solaris_x86`):

```bash
/export/home/patchs/Solaris_x86/installer
```

During the wizard, choose "All Shared Components" to upgrade your SUN Cluster version.

Now you need to upgrade all your component, choose to upgrade the whole cluster structure:

```bash
> /export/home/patchs/Solaris_x86/Product/sun_cluster/Solaris_10/Tools/scinstall

  *** Main Menu ***

    Please select from one of the following (*) options:

        1) Create a new cluster or add a cluster node
        2) Configure a cluster to be JumpStarted from this install server
      * 3) Manage a dual-partition upgrade
      * 4) Upgrade this cluster node
      * 5) Print release information for this cluster node

      * ?) Help with menu options
      * q) Quit

    Option:  4
```

So you normally need to select option 1.

After that, you need to update all agents (using option 2).

Then reboot this node and apply this to all the nodes of your cluster.

### 3.2 Fast way

If you need a faster way, you can do like this (you need to exactly know what you have installed):

```bash
cd /export/home/patchs/Solaris_x86/Product/sun_cluster/Solaris_10/Tools
./scinstall -u update
./scinstall -u update -d /export/home/patchs/Solaris_x86/Product/sun_cluster_agents -s tomcat,smb,PostgreSQL,mys,dhcp,container,9ias,oracle,iws,dns,apache
reboot
```

## 4. Verification

You can now verify by looking at `/etc/cluster/release` file:

```bash
> cat /etc/cluster/release
                    Sun Cluster 3.2u3 for Solaris 10 i386
          Copyright 2008 Sun Microsystems, Inc. All Rights Reserved.
```
