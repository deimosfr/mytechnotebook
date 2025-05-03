---
weight: 999
url: "/Cisco_VPN_Client_\\:_Installation_sur_une_Red_Hat_6/"
title: "Cisco VPN Client: Installation on Red Hat 6"
description: "A guide to installing the Cisco VPN Client on Red Hat 6, including necessary patches and configurations"
categories: ["Network", "Linux", "Security"]
date: "2012-04-06T00:15:00+02:00"
lastmod: "2012-04-06T00:15:00+02:00"
tags: ["Cisco", "VPN", "Red Hat", "Enterprise Linux"]
toc: true
---

![Cisco](/images/poweredbycisco.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 4.8.02 |
| **Operating System** | RHEL 6.2 |
| **Website** | [Cisco Website](https://www.cisco.com/) |
| **Last Update** | 06/04/2012 |
{{< /table >}}

## Introduction

Installing the Cisco VPN Client on Linux can be quite challenging! That's why I decided to create this documentation where you just need to copy and paste the commands.

## Prerequisites

You'll need all the tools required for compilation (gcc, g++...):

```bash
yum install gcc glibc
```

## Installation

Let's download a version 4:

```bash
wget http://projects.tuxx-home.at/ciscovpn/clients/linux/4.8.02/vpnclient-linux-x86_64-4.8.02.0030-k9.tar.gz
```

Let's extract everything:

```bash
tar -xzvf vpnclient-linux-x86_64-4.8.02.0030-k9.tar.gz
cd vpnclient
```

Apply the following patch for 64-bit compatibility:

```c
diff -urN vpnclient.orig/Makefile vpnclient/Makefile
--- vpnclient.orig/Makefile	2008-06-23 17:59:12.000000000 +0100
+++ vpnclient/Makefile	2008-07-09 23:16:54.000000000 +0100
@@ -12,7 +12,9 @@
 SOURCE_OBJS := linuxcniapi.o frag.o IPSecDrvOS_linux.o interceptor.o linuxkernelapi.o
 
 ifeq ($(SUBARCH),x86_64)
-CFLAGS += -mcmodel=kernel -mno-red-zone
+# Must NOT fiddle with CFLAGS
+# CFLAGS += -mcmodel=kernel -mno-red-zone
+EXTRA_CFLAGS += -mcmodel=kernel -mno-red-zone
 NO_SOURCE_OBJS := libdriver64.so
 else
 NO_SOURCE_OBJS := libdriver.so
diff -urN vpnclient.orig/frag.c vpnclient/frag.c
--- vpnclient.orig/frag.c	2008-06-23 17:59:12.000000000 +0100
+++ vpnclient/frag.c	2008-07-09 23:16:54.000000000 +0100
@@ -22,7 +22,9 @@
 #include "frag.h"
 
 #if LINUX_VERSION_CODE >= KERNEL_VERSION(2,6,22)
-#define SKB_IPHDR(skb) ((struct iphdr*)skb->network_header)
+/* 2.6.22 added an inline function for 32-/64-bit usage here, so use it.
+ */
+#define SKB_IPHDR(skb) ((struct iphdr*)skb_network_header)
 #else
 #define SKB_IPHDR(skb) skb->nh.iph
 #endif
diff -urN vpnclient.orig/interceptor.c vpnclient/interceptor.c
--- vpnclient.orig/interceptor.c	2008-06-23 17:59:12.000000000 +0100
+++ vpnclient/interceptor.c	2008-07-09 23:34:51.000000000 +0100
@@ -637,19 +637,30 @@
 
     reset_inject_status(&pBinding->recv_stat);
 #if LINUX_VERSION_CODE >= KERNEL_VERSION(2,6,22)
-    if (skb->mac_header)
+/* 2.6.22 added an inline function for 32-/64-bit usage here, so use it.
+ */
+    if (skb_mac_header_was_set(skb))
 #else
     if (skb->mac.raw)
 #endif
     {
 #if LINUX_VERSION_CODE >= KERNEL_VERSION(2,6,22)
-        hard_header_len = skb->data - skb->mac_header;
+/* 2.6.22 added an inline function for 32-/64-bit usage here, so use it.
+ */
+        hard_header_len = skb->data - skb_mac_header(skb);
 #else
         hard_header_len = skb->data - skb->mac.raw;
 #endif
         if ((hard_header_len < 0) || (hard_header_len > skb_headroom(skb)))
         {
             printk(KERN_DEBUG "bad hh len %d\n", hard_header_len);
+
+            printk(KERN_DEBUG "bad hh len %d, mac: %d, data: %p, head: %p\n",
+                hard_header_len,
+                skb->mac_header,    /* actualy ptr in 32-bit */
+                skb->data,
+                skb->head);
+
             hard_header_len = 0;
         }
     }
@@ -664,7 +675,9 @@
     {
     case ETH_HLEN:
 #if LINUX_VERSION_CODE >= KERNEL_VERSION(2,6,22)
-        CniNewFragment(ETH_HLEN, skb->mac_header, &MacHdr, CNI_USE_BUFFER);
+/* 2.6.22 added an inline function for 32-/64-bit usage here, so use it.
+ */
+        CniNewFragment(ETH_HLEN, skb_mac_header(skb), &MacHdr, CNI_USE_BUFFER);
 #else
         CniNewFragment(ETH_HLEN, skb->mac.raw, &MacHdr, CNI_USE_BUFFER);
 #endif
@@ -782,7 +795,9 @@
 #endif //LINUX_VERSION_CODE >= KERNEL_VERSION(2,4,0)
     reset_inject_status(&pBinding->send_stat);
 #if LINUX_VERSION_CODE >= KERNEL_VERSION(2,6,22)
-    hard_header_len = skb->network_header - skb->data;
+/* 2.6.22 added an inline function for 32-/64-bit usage here, so use it.
+ */
+    hard_header_len = skb_network_header(skb) - skb->data;
 #else
     hard_header_len = skb->nh.raw - skb->data;
 #endif
diff -urN vpnclient.orig/linuxcniapi.c vpnclient/linuxcniapi.c
--- vpnclient.orig/linuxcniapi.c	2008-06-23 17:59:12.000000000 +0100
+++ vpnclient/linuxcniapi.c	2008-07-09 23:16:54.000000000 +0100
@@ -338,8 +338,12 @@
     skb->ip_summed = CHECKSUM_UNNECESSARY;
 
 #if LINUX_VERSION_CODE >= KERNEL_VERSION(2,6,22)
-    skb->network_header = (sk_buff_data_t) skb->data;
-    skb->mac_header = (sk_buff_data_t)pMac;
+/* 2.6.22 added an inline function for 32-/64-bit usage here, so use it.
+ * We have to use (pMac - skb->data) to get an offset.
+ * We need to cast ptrs to byte ptrs and take the difference.
+ */
+    skb_reset_network_header(skb);
+    skb_set_mac_header(skb, (int)((void *)pMac - (void *)skb->data));
 #else
     skb->nh.iph = (struct iphdr *) skb->data;
     skb->mac.raw = pMac;
@@ -478,8 +482,12 @@
     skb->dev = pBinding->pDevice;
 
 #if LINUX_VERSION_CODE >= KERNEL_VERSION(2,6,22)
-    skb->mac_header = (sk_buff_data_t)pMac;
-    skb->network_header = (sk_buff_data_t)pIP;
+/* 2.6.22 added an inline function for 32-/64-bit usage here, so use it.
+ * We have to use (pIP/pMac - skb->data) to get an offset.
+ * We need to cast ptrs to byte ptrs and take the difference.
+ */
+    skb_set_mac_header(skb, (int)((void *)pMac - (void *)skb->data));
+    skb_set_network_header(skb, (int)((void *)pIP - (void *)skb->data));
 #else
     skb->mac.raw = pMac;
     skb->nh.raw = pIP;
@@ -487,8 +495,13 @@
 
     /*ip header length is in 32bit words */
 #if LINUX_VERSION_CODE >= KERNEL_VERSION(2,6,22)
-    skb->transport_header = (sk_buff_data_t)
-      (pIP + (((struct iphdr*)(skb->network_header))->ihl * 4));
+/* 2.6.22 added an inline function for 32-/64-bit usage here, so use it.
+ * We have to use (pIP - skb->data) to get an offset.
+ * We need to cast ptrs to byte ptrs and take the difference.
+ */
+    skb_set_transport_header(skb,
+        ((int)((void *)pIP - (void *)skb->data) +
+           ((((struct iphdr*)(skb_network_header(skb))))->ihl * 4)));
 #else
     skb->h.raw = pIP + (skb->nh.iph->ihl * 4);
 #endif
diff -urN vpnclient.orig/linuxkernelapi.c vpnclient/linuxkernelapi.c
--- vpnclient.orig/linuxkernelapi.c	2008-06-23 17:59:12.000000000 +0100
+++ vpnclient/linuxkernelapi.c	2008-07-09 23:16:54.000000000 +0100
@@ -9,7 +9,10 @@
     void*rc = kmalloc(size, GFP_ATOMIC);
     if(NULL == rc)
     {
-        printk("<1> os_malloc size %d failed\n",size);
+/* Allow for 32- or 64-bit size        
+ *        printk("<1> os_malloc size %d failed\n",size);
+ */
+        printk("<1> os_malloc size %ld failed\n", (long)size);
     }
 
     return rc;
```

And also this patch for kernels above 2.6.31:

```c
diff -uBbr vpnclient.orig/interceptor.c vpnclient/interceptor.c
--- vpnclient.orig/interceptor.c	2009-10-07 20:22:56.000000000 +0200
+++ vpnclient/interceptor.c	2009-10-07 20:28:48.000000000 +0200
@@ -120,6 +120,14 @@
     .notifier_call = handle_netdev_event,
 };
 
+#if LINUX_VERSION_CODE >= KERNEL_VERSION(2,6,31)
+static const struct net_device_ops vpn_netdev_ops = {
+ .ndo_start_xmit = interceptor_tx,
+ .ndo_get_stats = interceptor_stats,
+ .ndo_do_ioctl = interceptor_ioctl,
+};
+#endif
+
 #if LINUX_VERSION_CODE >= KERNEL_VERSION(2,6,22)
 static int
 #else
@@ -128,10 +136,13 @@
 interceptor_init(struct net_device *dev)
 {
     ether_setup(dev);
-
+    #if LINUX_VERSION_CODE >= KERNEL_VERSION(2,6,31)
+    dev->netdev_ops = &vpn_netdev_ops;
+    #else
     dev->hard_start_xmit = interceptor_tx;
     dev->get_stats = interceptor_stats;
     dev->do_ioctl = interceptor_ioctl;
+    #endif
 
     dev->mtu = ETH_DATA_LEN-MTU_REDUCTION;
     kernel_memcpy(dev->dev_addr, interceptor_eth_addr,ETH_ALEN);
@@ -268,8 +279,13 @@
     Bindings[i].original_mtu = dev->mtu;
 
     /*replace the original send function with our send function */
+    #if LINUX_VERSION_CODE >= KERNEL_VERSION(2,6,31)
+    Bindings[i].InjectSend = dev->netdev_ops->ndo_start_xmit;
+    dev->netdev_ops->ndo_start_xmit = replacement_dev_xmit;
+    #else
     Bindings[i].InjectSend = dev->hard_start_xmit;
     dev->hard_start_xmit = replacement_dev_xmit;
+    #endif
 
     /*copy in the ip packet handler function and packet type struct */
     Bindings[i].InjectReceive = original_ip_handler.orig_handler_func;
@@ -291,7 +307,11 @@
     if (b)
     {   
         rc = 0;
+        #if LINUX_VERSION_CODE >= KERNEL_VERSION(2,6,31)
+        dev->netdev_ops->ndo_start_xmit = b->InjectSend;
+        #else
         dev->hard_start_xmit = b->InjectSend;
+        #endif
         kernel_memset(b, 0, sizeof(BINDING));
     }
     else
```

Now let's apply the patches:

```bash
patch < vpnclient-linux-4.8.02-64bit.patch 
patch < vpnclient-linux-2.6.31-final.diff
```

The last two modifications for the route:

```bash
sed -i 's/^CFLAGS/EXTRA_CFLAGS/' Makefile
sed -i 's/const\ struct\ net_device_ops\ \*netdev_ops;/struct\ net_device_ops\ \*netdev_ops;/' `find /usr/src -name netdevice.h`
```

And finally, let's run the installation:

```bash
./vpn_install
```

This should produce output like:

```
Cisco Systems VPN Client Version 4.8.02 (0030) Linux Installer
Copyright (C) 1998-2006 Cisco Systems, Inc. All Rights Reserved.

By installing this product you agree that you have read the
license.txt file (The VPN Client license) and will comply with
its terms. 


Directory where binaries will be installed [/usr/local/bin]

Automatically start the VPN service at boot time [yes]

In order to build the VPN kernel module, you must have the
kernel headers for the version of the kernel you are running.


Directory containing linux kernel source code [/lib/modules/2.6.32-131.0.15.el6.x86_64/build]

* Binaries will be installed in "/usr/local/bin".
* Modules will be installed in "/lib/modules/2.6.32-131.0.15.el6.x86_64/CiscoVPN".
* The VPN service will be started AUTOMATICALLY at boot time.
* Kernel source from "/lib/modules/2.6.32-131.0.15.el6.x86_64/build" will be used to build the module.

Is the above correct [y]

Making module
make -C /lib/modules/2.6.32-131.0.15.el6.x86_64/build SUBDIRS=/root/test/vpnclient modules
make[1]: Entering directory `/usr/src/kernels/2.6.32-131.0.15.el6.x86_64'
  CC [M]  /root/test/vpnclient/linuxcniapi.o
  CC [M]  /root/test/vpnclient/frag.o
  CC [M]  /root/test/vpnclient/interceptor.o
/root/test/vpnclient/interceptor.c: In function 'interceptor_init':
/root/test/vpnclient/interceptor.c:140: warning: assignment discards qualifiers from pointer target type
  CC [M]  /root/test/vpnclient/linuxkernelapi.o
  LD [M]  /root/test/vpnclient/cisco_ipsec.o
  Building modules, stage 2.
  MODPOST 1 modules
WARNING: could not find /root/test/vpnclient/.libdriver64.so.cmd for /root/test/vpnclient/libdriver64.so
  CC      /root/test/vpnclient/cisco_ipsec.mod.o
  LD [M]  /root/test/vpnclient/cisco_ipsec.ko.unsigned
  NO SIGN [M] /root/test/vpnclient/cisco_ipsec.ko
make[1]: Leaving directory `/usr/src/kernels/2.6.32-131.0.15.el6.x86_64'
Create module directory "/lib/modules/2.6.32-131.0.15.el6.x86_64/CiscoVPN".
Copying module to directory "/lib/modules/2.6.32-131.0.15.el6.x86_64/CiscoVPN".
Already have group 'bin'

Creating start/stop script "/etc/init.d/vpnclient_init".
    /etc/init.d/vpnclient_init
Enabling start/stop script for run level 3,4 and 5.
Creating global config /etc/opt/cisco-vpnclient

Installing license.txt (VPN Client license) in "/opt/cisco-vpnclient/":
    /opt/cisco-vpnclient/license.txt

Installing bundled user profiles in "/etc/opt/cisco-vpnclient/Profiles/":
* New Profiles    : sample 

Copying binaries to directory "/opt/cisco-vpnclient/bin".
Adding symlinks to "/usr/local/bin".
    /opt/cisco-vpnclient/bin/vpnclient
    /opt/cisco-vpnclient/bin/cisco_cert_mgr
    /opt/cisco-vpnclient/bin/ipseclog
Copying setuid binaries to directory "/opt/cisco-vpnclient/bin".
    /opt/cisco-vpnclient/bin/cvpnd
Copying libraries to directory "/opt/cisco-vpnclient/lib".
    /opt/cisco-vpnclient/lib/libvpnapi.so
Copying header files to directory "/opt/cisco-vpnclient/include".
    /opt/cisco-vpnclient/include/vpnapi.h

Setting permissions.
    /opt/cisco-vpnclient/bin/cvpnd (setuid root)
    /opt/cisco-vpnclient (group bin readable)
    /etc/opt/cisco-vpnclient (group bin readable)
    /etc/opt/cisco-vpnclient/Profiles (group bin readable)
    /etc/opt/cisco-vpnclient/Certificates (group bin readable)
* You may wish to change these permissions to restrict access to root.
* You must run "/etc/init.d/vpnclient_init start" before using the client.
* This script will be run AUTOMATICALLY every time you reboot your computer.
```

## Resources
- http://micro.stanford.edu/wiki/How_to_install_and_configure_the_Cisco_VPN_client_on_a_Linux_computer
- http://www.lamnk.com/blog/vpn/how-to-install-cisco-vpn-client-on-ubuntu-jaunty-jackalope-and-karmic-koala-64-bit/
- http://forum.tuxx-home.at/viewtopic.php?f=15&t=957
