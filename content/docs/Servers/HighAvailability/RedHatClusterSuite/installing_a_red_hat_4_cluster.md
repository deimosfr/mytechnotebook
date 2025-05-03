---
weight: 999
url: "/Installation_d\\'un_cluster_Red_Hat_4/"
title: "Installing a Red Hat 4 Cluster"
description: "Step-by-step guide to installing and configuring a Red Hat Cluster Suite on Red Hat Enterprise Linux 4."
categories: ["Linux", "Red Hat"]
date: "2012-02-26T13:24:00+02:00"
lastmod: "2012-02-26T13:24:00+02:00"
tags: ["Cluster", "Red Hat", "High Availability"]
toc: true
---

**To install Red Hat Enterprise with Cluster Suite, take CDs or DVD and boot from it.**

We are going to see step by step the interesting components.

## Language

```
Set the English language and not French to have more verbose output during future problems.
```

## Keyboard

```
Keep the French keyboard
```

## Partitioning (do the same partitioning on all the cluster nodes)

{{< table "table-hover table-striped" >}}
| Location | Type | Size | LVM Name |
|---------|------|------|----------|
| /boot | ext3 | 128 MB | **NO LVM!!!** |
| LVM | Group | Full Size (space left) | VolGroup00 |
| / | ext3 | 4096 MB | racine |
| /var | ext3 | 4096 MB | var |
| /usr | ext3 | 10240 MB | usr |
| /tmp | ext3 | 2048 MB | tmp |
| swap | | 4096 MB | swap |
| /home | ext3 | Space left | home |
{{< /table >}}

## List of packages installation

```
base
base-x
kde-desktop or gnome-desktop
editors (only vim)
mail-server (- spamassassin)
development-tools
compat-arch-development
admin-tools (+ sysstem-boot-config)
system-tools
compat-arch-support
```

After rebooting, install these packages from the cluster suite (cdrom/RedHat/RPMS). Warning: pay attention for SMP; if not SMP servers, do not install it, install other relative packages.

```
ccs
cman
cman-kernel-smp
dlm
dlm-kernel-smp
fence
gulm
ipvsadm
magma
magma-plugins
perl-Net-Telnet
rgmanager
system-config-cluster
```

Or run this command directly:

```bash
rpm -ivh --aid ccs-1.0.7-0.x86_64.rpm cman-1.0.11-0.x86_64.rpm cman-kernel-smp-2.6.9-45.2.x86_64.rpm  dlm-kernel-smp-2.6.9-42.10.x86_64.rpm fence-1.32.25-1.x86_64.rpm gulm-1.0.7-0.x86_64.rpm ipvsadm-1.24-6.x86_64.rpm magma-1.0.6-0.x86_64.rpm magma-plugins-1.0.9-0.x86_64.rpm perl-Net-Telnet-3.03-3.noarch.rpm rgmanager-1.9.53-0.x86_64.rpm system-config-cluster-1.0.27-1.0.noarch.rpm dlm-1.0.1-1.x86_64.rpm
```

A little update might not be so bad (optional):

```bash
up2date -u
```

Modify the `/etc/hosts` like this on each node:

```bash
127.0.0.1               localhost.localdomain localhost
192.168.0.242           secondary.cluster.net secondary
192.168.0.241           primary.cluster.net primary
```

Here is an example of the cluster.conf file in `/etc/cluster`:

```xml
<?xml version="1.0"?>
<cluster config_version="13" name="alpha_cluster">
	<fence_daemon post_fail_delay="0" post_join_delay="3"/>
	<clusternodes>
		<clusternode name="primary.cluster.net" votes="1">
			<fence>
				<method name="1">
					<device name="fence1" nodename="primary.cluster.net"/>
				</method>
			</fence>
		</clusternode>
		<clusternode name="secondary.cluster.net" votes="1">
			<fence>
				<method name="1">
					<device name="fence2" nodename="secondary.cluster.net"/>
				</method>
			</fence>
		</clusternode>
	</clusternodes>
	<cman expected_votes="1" two_node="1"/>
	<fencedevices>
		<fencedevice agent="fence_manual" name="fence1"/>
		<fencedevice agent="fence_manual" name="fence2"/>
	</fencedevices>
	<rm>
		<failoverdomains>
			<failoverdomain name="myfailoverdomain" ordered="0" restricted="1">
				<failoverdomainnode name="primary.cluster.net" priority="1"/>
				<failoverdomainnode name="secondary.cluster.net" priority="1"/>
			</failoverdomain>
		</failoverdomains>
		<resources>
			<ip address="192.168.0.238" monitor_link="1"/>
			<fs device="/dev/sdb1" force_fsck="0" force_unmount="0" fsid="44028" fstype="ext3" mountpoint="/mnt/cluster" name="disk2" options="" self_fence="0"/>
		</resources>
		<service autostart="1" domain="cluster" name="clusterdev">
			<ip ref="192.168.0.238"/>
			<fs ref="disk2"/>
			<script file="/etc/init.d/ntpd" name="ntp"/>
		</service>
	</rm>
</cluster>
```

And of course, don't forget to run services:

```bash
service ccsd restart
service cman restart
service fenced restart
service rgmanager restart
```

If you want to stop the cluster, do it in the reverse order:

```bash
service rgmanager stop
service fenced stop
service cman stop
service ccsd stop
```

**It's very important to stop in this order!**
