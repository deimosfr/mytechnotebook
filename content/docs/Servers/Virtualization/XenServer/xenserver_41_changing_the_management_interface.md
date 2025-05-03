---
weight: 999
url: "/XenServer_4.1_\\:_Changer_l'interface_de_Management/"
title: "XenServer 4.1: Changing the Management Interface"
description: "How to change the management interface on XenServer 4.1 when you've selected the wrong one during installation."
categories: ["Linux", "Virtualization"]
date: "2009-02-04T18:58:00+02:00"
lastmod: "2009-02-04T18:58:00+02:00"
tags: ["XenServer", "Networking", "Virtualization", "Servers"]
toc: true
---

The other day I installed a XenServer 4.1 server. Unfortunately, I selected the wrong management interface during installation, choosing eth0 instead of eth1.

The question was: how to change the management interface on XenServer 4.1?

## Retrieving the UUID of the interface

First, we need to retrieve the UUID of the interface (they call it "PIF"):

```bash
[root@xenserver-backup2 ~]# xe pif-list 
uuid ( RO)                 : f34c8861-94d7-c3da-0437-3a068b273db5
                device ( RO): eth3
    currently-attached ( RO): true
                  VLAN ( RO): -1
          network-uuid ( RO): c7c6e869-f611-11ad-871e-c130bbb9d13f


uuid ( RO)                 : ffcba213-30dd-09d2-ee7e-0708507516f0
                device ( RO): eth2
    currently-attached ( RO): true
                  VLAN ( RO): -1
          network-uuid ( RO): 55efa2d9-768b-52a0-4595-cebf771ce602


uuid ( RO)                 : 7e5946e3-3c11-ebec-2636-713a69ce60ae
                device ( RO): eth1
    currently-attached ( RO): true
                  VLAN ( RO): -1
          network-uuid ( RO): 7755aae6-a565-6e0a-a6b0-49cef6abee5e


uuid ( RO)                 : b726e0f2-7b58-3489-dca7-0269630f8d4e
                device ( RO): eth0
    currently-attached ( RO): false
                  VLAN ( RO): -1
          network-uuid ( RO): ae8ca5b4-9fa1-912a-7cca-55f4060a782c
```

## Configure IP options

Then we assign IP options with a command like this:

```bash
xe pif-reconfigure-ip uuid=7e5946e3-3c11-ebec-2636-713a69ce60ae IP=192.168.0.185 netmask=255.255.255.0 gateway=192.168.0.245 DNS=192.168.0.216 mode=static
```

(The mode can be "none", "static" or "dhcp")

## Assign the management function

Finally, we assign the management function to this interface:

```bash
xe host-management-reconfigure pif-uuid=7e5946e3-3c11-ebec-2636-713a69ce60ae
```
