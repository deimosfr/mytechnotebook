---
weight: 999
url: "/LXC_\\:_Install_and_configure_the_Linux_Containers/"
title: "LXC: Install and configure the Linux Containers"
description: "A comprehensive guide on installing, configuring and using Linux Containers (LXC) on Debian systems."
categories: ["Virtualization", "Linux", "Containers"]
date: "2015-02-22T00:00:00+02:00"
lastmod: "2015-02-22T00:00:00+02:00"
tags: ["LXC", "Linux", "Containers", "Virtualization", "Debian"]
toc: true
---

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 0.8 |
| **Operating System** | Debian 7 |
| **Website** | [LXC Website](https://lxc.sourceforge.net/) |
| **Last Update** | 22/02/2015 |
{{< /table >}}

## Introduction

LXC is the userspace control package for Linux Containers, a lightweight virtual system mechanism sometimes described as "chroot on steroids".

LXC builds up from chroot to implement complete virtual systems, adding resource management and isolation mechanisms to Linux's existing process management infrastructure.

Linux Containers (lxc) implement:

- Resource management via "process control groups" (implemented via the cgroup filesystem)
- Resource isolation via new flags to the clone(2) system call (capable of create several types of new namespaces for things like PIDs and network routing)
- Several additional isolation mechanisms (such as the "-o newinstance" flag to the devpts filesystem).

The LXC package combines these Linux kernel mechanisms to provide a userspace container object, a lightweight virtual system with full resource isolation and resource control for an application or a system.

Linux Containers take a completely different approach than system virtualization technologies such as KVM and Xen, which started by booting separate virtual systems on emulated hardware and then attempted to lower their overhead via paravirtualization and related mechanisms. Instead of retrofitting efficiency onto full isolation, LXC started out with an efficient mechanism (existing Linux process management) and added isolation, resulting in a system virtualization mechanism as scalable and portable as chroot, capable of simultaneously supporting thousands of emulated systems on a single server while also providing lightweight virtualization options to routers and smart phones.

The first objective of this project is to make the life easier for the kernel developers involved in the containers project and especially to continue working on the Checkpoint/Restart new features. The lxc is small enough to easily manage a container with simple command lines and complete enough to be used for other purposes.[^1]

## Installation

To install LXC, we do not need too much packages. As I will want to manage my LXC containers with libvirt, I need to install it as well:

```bash
aptitude install lxc bridge-utils debootstrap git git-core
```

At the time where I write this sentence, there is an issue with LVM container creation (here is [a first Debian bug](https://bugs.debian.org/cgi-bin/bugreport.cgi?bug=680469) and [a second one](https://bugs.debian.org/cgi-bin/bugreport.cgi?bug=716839)) on Debian Wheezy and it doesn't seams to be resolve soon.

Here is a workaround to avoid errors during LVM containers initialization:

```bash
cd /tmp
git clone https://github.com/simonvanderveldt/lxc-debian-wheezy-template.git
cp lxc-debian-wheezy-template/lxc-debian-0.9.0-upstream /usr/share/lxc/templates/lxc-debian-squeeze
cp lxc-debian-wheezy-template/lxc-debian-0.9.0-upstream /usr/share/lxc/templates/lxc-debian-wheezy
sed -i 's/-squeeze/-wheezy/' /usr/share/lxc/templates/lxc-debian-wheezy
chmod 755 /usr/share/lxc/templates/lxc-debian-*
```

On Jessie, you'll also have to install this:

```bash
aptitude install cgroupfs-mount
```

### Kernel

It's recommended to get a recent kernel as LXC grow very fast, get better performances, stabilities and new features. To get a newer kernel, we're going to use a kernel from the testing repo:

```bash
# /etc/apt/preferences.d/kernel
Package: *
Pin: release a=stable
Pin-priority: 900

Package: *
Pin: release a=testing
Pin-priority: 100

Package: linux-image-*
Pin: release a=testing
Pin-priority: 1001

Package: linux-headers-*
Pin: release a=testing
Pin-priority: 1001

Package: linux-kbuild-*
Pin: release a=testing
Pin-priority: 1001
```

Add then this testing content:

```bash
# /etc/apt/sources.list.d/testing.list
# Testing
deb http://ftp.fr.debian.org/debian/ testing main non-free contrib
deb-src http://ftp.fr.debian.org/debian/ testing main non-free contrib
```

Then you can install the latest kernel image:

```bash
aptitude update
aptitude install linux-image-amd64
```

If it's not enough, you'll need to install the package with specific kernel version number corresponding to latest (ex. linux-image-3.11-2-amd64) and reboot on this new kernel.

## Configuration

### Cgroups

LXC is based on cgroups. Those are used to limit CPU, RAM etc...[You can check here for more informations]({{< ref "docs/Linux/Kernel/process_latency_and_kernel_timing.md" >}}).

We need to enable Cgroups. Add this line in fstab:

```bash
# /etc/fstab
[...]
cgroup  /sys/fs/cgroup  cgroup  defaults                        0       0
```

As we want to manage memory and swap on containers, as it's not available by default, add cgroup argument to grub to activate those functionality:

- cgroup RAM feature: "cgroup_enable=memory"
- cgroup SWAP feature: "swapaccount=1"

```bash {linenos=table,hl_lines=[10],anchorlinenos=true}
# /etc/default/grub
# If you change this file, run 'update-grub' afterwards to update
# /boot/grub/grub.cfg.
# For full documentation of the options in this file, see:
#   info -f grub -n 'Simple configuration'

GRUB_DEFAULT=0
GRUB_TIMEOUT=5
GRUB_DISTRIBUTOR=`lsb_release -i -s 2> /dev/null || echo Debian`
GRUB_CMDLINE_LINUX_DEFAULT="quiet cgroup_enable=memory swapaccount=1"
GRUB_CMDLINE_LINUX=""
[...]
```

Then regenerate grub config:

```bash
update-grub
```

{{< alert context="info" text="You'll need to reboot to make the cgroup memory feature active" />}}

Then mount it:

```bash
mount /sys/fs/cgroup
```

### Check LXC configuration

```bash
> lxc-checkconfig
Kernel config /proc/config.gz not found, looking in other places...
Found kernel config file /boot/config-3.2.0-4-amd64
--- Namespaces ---
Namespaces: enabled
Utsname namespace: enabled
Ipc namespace: enabled
Pid namespace: enabled
User namespace: enabled
Network namespace: enabled
Multiple /dev/pts instances: enabled

--- Control groups ---
Cgroup: enabled
Cgroup clone_children flag: enabled
Cgroup device: enabled
Cgroup sched: enabled
Cgroup cpu account: enabled
Cgroup memory controller: enabled
Cgroup cpuset: enabled

--- Misc ---
Veth pair device: enabled
Macvlan: enabled
Vlan: enabled
File capabilities: enabled

Note : Before booting a new kernel, you can check its configuration
usage : CONFIG=/path/to/config /usr/bin/lxc-checkconfig
```

{{< alert context="danger" text="All should be enabled to ensure it will work as expected!" />}}

### Network

#### No specific configuration (same than host)

If you don't configure your network configuration after container initialization, you'll have the exact same configuration on your guests (containers) than your host. That mean all network interfaces are available on the guests and they will have full access to the host.

{{< alert context="danger" text="This is not the recommended solution for production usages" />}}

- The pro of that "no" configuration, is to have network working out of the box for the guests (perfect for quick tests)
- Another con, is to have the access to process on host. **I mean that a SSH server running on host will have it's port available on the guest too**. So you cannot have a SSH server running on guests without changing port (or you'll have a network binding conflict).

You can easily check this configuration in opening a port on the host (here 80):

```bash
nc -lp 80
```

now on a guest, you can see it listening:

```bash
> netstat -aunt | grep 80
tcp        0      0 0.0.0.0:80              0.0.0.0:*               LISTEN
```

#### Nat configuration

How does it works for NAT configuration?

1. You need to **choose** which kind of configuration you want to use: **libvirt or dnsmaq**
2. Iptables: to help to access nated containers from outside and help containers to get internet
3. Configure the network container

##### With Libvirt

You'll need to install libvirt first:

```bash
aptitude install libvirt-bin
```

Nat is the default configuration. But you may need to do some adjustements. Add the forwarding to sysctl:

```bash
# /etc/sysctl.conf
net.ipv4.ip_forward = 1
```

Check that the connecting is active:

```bash
> virsh net-list --all
Name                 State      Autostart
-----------------------------------------
default              inactive   no
```

If it's not the case, then set the default network configuration:

```bash
virsh net-define /etc/libvirt/qemu/networks/default.xml
virsh net-autostart default
virsh net-start default
```

Then you should see it enable:

```bash
> virsh net-list --all
Name                 State      Autostart
-----------------------------------------
default              active     yes
```

Edit the configuration to add your range of IP:

```xml {linenos=table,hl_lines=[5,7],anchorlinenos=true}
<network>
  <name>default</name>
  <bridge name="virbr0" />
  <forward/>
  <ip address="192.168.122.1" netmask="255.255.255.0">
    <dhcp>
      <range start="192.168.122.100" end="192.168.122.200" />
    </dhcp>
  </ip>
</network>
```

Now it's done, restart libvirt.

##### With dnsmasq

Libvirt is not necessary as for the moment it doesn't manage LXC containers very well. So you can manage your own dnsmasq server to give DNS and DHCP to your containers. First of all, install it:

```bash
aptitude install dnsmasq dnsmasq-utils
```

Then configure it:

```bash
# /etc/dnsmasq.conf
# Bind it to the LXC interface
interface=lxcbr0
bind-interfaces
# Want DHCP client FQDN
dhcp-fqdn
# Domain name and ip range with lease time
domain=deimos.fr,192.168.122.0/24
dhcp-range=192.168.122.100,192.168.122.200,1h
# DHCP options
dhcp-option=40,deimos.fr
log-dhcp
```

Then restart the dnsmasq service. And now configure the lxcbr0 interface:

```bash
# /etc/network/interfaces
# This bridge will is used to NAT LXC containers' traffic
auto lxcbr0
iface lxcbr0 inet static
    pre-up brctl addbr lxcbr0
    bridge_fd 0
    bridge_maxwait 0
    address 192.168.122.1
    netmask 255.255.255.0
    post-up iptables -A FORWARD -i lxcbr0 -s 192.168.122.1/24 -j ACCEPT
    post-up iptables -A POSTROUTING -t nat -s 192.168.122.1/24 -j MASQUERADE
    # add checksum so that dhclient does not complain.
    # udp packets staying on the same host never have a checksum filled else
    post-up iptables -A POSTROUTING -t mangle -p udp --dport bootpc -s 192.168.122.1/24 -j CHECKSUM --checksum-fill
```

##### Iptables

You may need to configure iptables for example if you're on a dedicated box where the provider doesn't allow bridge configuration. Here is a working iptables configuration to permit incoming connexions to Nated guests:

```bash
#!/bin/bash
# Made by Pierre Mavro / Deimosfr
# This script will Nat you KVM/containers hosts
# and help you to get access from outside

#-------------------------------------------------------------------------
# Essentials
#-------------------------------------------------------------------------

IPTABLES='/sbin/iptables'
modprobe nf_conntrack_ftp

#-------------------------------------------------------------------------
# Physical and virtual interfaces definitions
#-------------------------------------------------------------------------

# Interfaces
wan1_if="eth0"
wan2_if="eth0:0"
kvm_if="virbr0"

#-------------------------------------------------------------------------
# Networks definitions
#-------------------------------------------------------------------------

# Networks
wan1_ip="x.x.x.x"
wan2_ip="x.x.x.x"
vms_net="192.168.122.0/24"

# Dedibox internals IPs
web_ip="192.168.122.10"
mail_ip="192.168.122.20"

#-------------------------------------------------------------------------
# Global Rules input / output / forward
#-------------------------------------------------------------------------

# Flushing tables
$IPTABLES -F
$IPTABLES -X
$IPTABLES -t nat -F

# Define default policy
$IPTABLES -P INPUT DROP
$IPTABLES -P OUTPUT ACCEPT
$IPTABLES -P FORWARD ACCEPT

## Loopback accepte
${IPTABLES} -A FORWARD -i lo -o lo -j ACCEPT
${IPTABLES} -A INPUT -i lo -j ACCEPT
${IPTABLES} -A OUTPUT -o lo -j ACCEPT

# Allow KVM DHCP/dnsmasq
${IPTABLES} -A INPUT -i $kvm_if -p udp --dport 67 -j ACCEPT
${IPTABLES} -A INPUT -i $kvm_if -p udp --dport 69 -j ACCEPT

$IPTABLES -A INPUT -j ACCEPT -d $vms_net
$IPTABLES -A INPUT -j ACCEPT -m state --state ESTABLISHED,RELATED

#-------------------------------------------------------------------------
# Allow masquerading for KVM VMs
#-------------------------------------------------------------------------

# Activating masquerade to get Internet from KVM VMs
$IPTABLES -t nat -A POSTROUTING -o $wan1_if -s $vms_net -j MASQUERADE

#-------------------------------------------------------------------------
# Allow ports on KVM host
#-------------------------------------------------------------------------

# Allow ICMP
$IPTABLES -A INPUT -j ACCEPT -p icmp

# SSH access
$IPTABLES -A INPUT -j ACCEPT -p tcp --dport 22

# HTTPS access
$IPTABLES -A INPUT -j ACCEPT -p tcp --dport 443

#-------------------------------------------------------------------------
# Redirections for incoming connections (wan1)
#-------------------------------------------------------------------------

# HTTP access
$IPTABLES -t nat -A PREROUTING -p tcp --dport 80 -d $wan1_ip -j DNAT --to-destination $web_ip:80

# HTTP access
$IPTABLES -t nat -A PREROUTING -p tcp --dport 443 -d $wan1_ip -j DNAT --to-destination $web_ip:443

# Mail for mailsrv
$IPTABLES -t nat -A PREROUTING -p tcp --dport 25 -d $wan1_ip -j DNAT --to-destination $mail_ip:25

#-------------------------------------------------------------------------
# Reload fail2ban
#-------------------------------------------------------------------------
/etc/init.d/fail2ban reload
```

##### Nat on containers

###### DHCP

On each containers you want to use NAT configuration, you need to add those lines for DHCP configuration[^2]:

```bash
# /var/lib/lxc/mycontainer/config
## Network
lxc.network.type = veth
lxc.network.flags = up

# Network host side
lxc.network.link = virbr0
lxc.network.veth.pair = veth0
# lxc.network.hwaddr = 00:FF:AA:00:00:01

# Network container side
lxc.network.name = eth0
lxc.network.ipv4 = 0.0.0.0/24
```

Then in the LXC container (mount the LV if you did LVM) configure the network like this:

```bash
# /var/lib/lxc/mycontainer/rootfs/etc/network/interfaces
auto eth0
iface eth0 inet dhcp
```

###### Static IP

{{< alert context="info" text="This is only applicable for Libvirt" />}}

You can also configure manual static IP if you want by changing 'lxc.network.ipv4'. Another elegant method is to ask DHCP to fix it:

```xml {linenos=table,hl_lines=[8],anchorlinenos=true}
<network>
  <name>default</name>
  <bridge name="virbr0" />
  <forward/>
  <ip address="192.168.122.1" netmask="255.255.255.0">
    <dhcp>
      <range start="192.168.122.2" end="192.168.122.254" />
      <host mac="00:FF:AA:00:00:01" name="fixed.deimos.fr" ip="192.168.122.101" />
    </dhcp>
  </ip>
</network>
```

{{% alert context="warning" %}}
Do not forget to fix lxc.network.hwaddr parameter. Here is a way to generate mac address:

```bash
openssl rand -hex 6 | sed 's/\\(..\\\)/\\1:/g; s/.$//'
```

{{% /alert %}}

#### Private container interface

You can create a private interface for your containers. Containers will be able to communicate together though this dedicated interface. Here are the steps to create one between 2 hosts.

On the host server, install UML utilities:

```bash
aptitude install uml-utilities
```

Then edit the network configuration file and add a bridge:

```bash
# /etc/network/interfaces
auto privbr0
iface privbr0 inet static
    pre-up /usr/sbin/tunctl -t tap0
    pre-up /sbin/ifup tap0
    post-down /sbin/ifdown tap0
    bridge_ports tap0
    bridge_fd 0
```

You can restart your network or launch it manually if you can't restart now:

```bash
tunctl -t tap0
brctl addbr privbr0
brctl addif privbr0 tap0
```

Then edit both containers that will have this dedicated interface and replace or add those lines:

```bash
# /var/lib/lxc/mycontainer/config
# Private interface
lxc.network.type = veth
lxc.network.flags = up
lxc.network.link = privbr0
lxc.network.ipv4 = 10.0.0.1
```

Now start the container and you'll have the 10.0.0.X dedicated network.

#### Bridged configuration

Now modify your /etc/network/interfaces file to add bridged configuration:

```bash
# /etc/network/interfaces
auto lo br0
iface lo inet loopback

# The primary network interface
iface br0 inet static
        address 192.168.0.80
        netmask 255.255.255.0
        gateway 192.168.0.252
        broadcast 192.168.0.255
        network 192.168.0.0
        bridge_ports eth0
        bridge_fd 9
        bridge_hello 2
        bridge_maxage 12
        bridge_stp off
```

br0 is replacing eth0 for bridging.

Here is another configuration with 2 network cards and 2 bridges:

```bash
# /etc/network/interfaces
# This file describes the network interfaces available on your system
# and how to activate them. For more information, see interfaces(5).

# The loopback network interface
auto lo br0 br1
iface lo inet loopback

# DMZ
iface br0 inet static
	address 192.168.10.1
	netmask 255.255.255.0
	gateway 192.168.10.254
	network 192.168.10.0
	broadcast 192.168.10.255
        bridge_ports eth0
        bridge_fd 9
        bridge_hello 2
        bridge_maxage 12
        bridge_stp off

# Internal
iface br1 inet static
	address 192.168.0.1
	netmask 255.255.255.0
	gateway 192.168.0.254
	network 192.168.0.0
	broadcast 192.168.0.255
        bridge_ports eth1
        bridge_fd 9
        bridge_hello 2
        bridge_maxage 12
        bridge_stp off
```

#### VLAN Bridged configuration

And a last one with Vlans bridged ([look at this documentation to enable it before]({{< ref "docs/Servers/Network/setting_up_vlan.md" >}})):

We will need to use etables (iptables for bridged interfaces). Install this:

```bash
aptitude install ebtables
```

Check you etables configuration:

```bash
# /etc/default/ebtables
EBTABLES_LOAD_ON_START="yes"
EBTABLES_SAVE_ON_STOP="yes"
EBTABLES_SAVE_ON_RESTART="yes"
```

And enable VLAN tagging on bridged interfaces:

```bash
ebtables -t broute -A BROUTING -i eth0 -p 802_1Q -j DROP
```

```bash {linenos=table,hl_lines=["10-16","35-42"],anchorlinenos=true}
# /etc/network/interfaces
# This file describes the network interfaces available on your system
# and how to activate them. For more information, see interfaces(5).

# The loopback network interface
auto lo
iface lo inet loopback

# The primary network interface
allow-hotplug eth0
auto eth0
iface eth0 inet manual

auto eth0.110
iface eth0.110 inet manual
        vlan_raw_device eth0

# The bridged interface
auto vmbr0
iface vmbr0 inet static
        address 192.168.100.1
        netmask 255.255.255.0
        network 192.168.100.0
        broadcast 192.168.100.255
        gateway 192.168.100.254
        # dns-* options are implemented by the resolvconf package, if installed
        dns-nameservers 192.168.100.254
        dns-search deimos.fr
        bridge_ports eth0
        bridge_fd 9
        bridge_hello 2
        bridge_maxage 12
        bridge_stp off

auto vmbr0.110
iface vmbr0.110 inet static
        address 192.168.110.1
        netmask 255.255.255.0
        bridge_ports eth0.190
        bridge_stp off
        bridge_maxwait 0
        bridge_fd 0
```

### Security

It's recommended to use Grsecurity kernel (may be not compatible with the [testing kernel](#kernel))or Apparmor.

With Grsecurity, here are the parameters[^3]:

```bash
# /etc/sysctl.d/grsecurity.conf
kernel.grsecurity.chroot_deny_mount = 0
kernel.grsecurity.chroot_deny_chroot = 0
kernel.grsecurity.chroot_deny_chmod = 0
kernel.grsecurity.chroot_deny_mknod = 0
kernel.grsecurity.chroot_caps = 0
kernel.grsecurity.chroot_findtask = 0
```

## Basic Usage

### Create a container

#### Classic method

To create a container with a wizard:

```bash
lxc-create -n mycontainer -t debian
```

or

```bash
lxc-create -n mycontainer -t debian-wheezy
```

or

```bash
lxc-create -n mycontainer -t debian-wheezy -f /etc/lxc/lxc-nat.conf
```

- n: the name of the container
- t: the template of the container. You can find the list in this folder: `/usr/share/lxc/templates`
- f: the configuration template

It will deploy through debootstrap a new container.

If you want to plug yourself to a container through the console, your first need to create devices:

```bash
chroot /var/lib/lxc/mycontainer/rootfs
mknod -m 666 /dev/tty1 c 4 1
mknod -m 666 /dev/tty2 c 4 2
mknod -m 666 /dev/tty3 c 4 3
mknod -m 666 /dev/tty4 c 4 4
mknod -m 666 /dev/tty5 c 4 5
mknod -m 666 /dev/tty6 c 4 6
```

Then you'll be able to connect:

```bash
lxc-console -n mycontainer
```

#### LVM method

If you're using LVM to store your containers (strongly recommended), you can ask to LXC to auto create the logical volume and mkfs it for you:

```bash
lxcname=mycontainerlvm
lxc-create -t debian-wheezy -n $lxcname -B lvm --vgname lxc --lvname $lxcname --fssize 4G --fstype ext3
```

- t: specify the wished template
- B: we want to use LVM as backend (BTRFS is also supported)
- vgname: set the volume group (VG) name where logical volume (LV) should be created
- lvname: set the wished LV name for that container
- fssize: set the size of the LV
- fstype: set the filesystem for this container (full list is available in /proc/filesystems)

#### BTRFS method

If your host has a btrfs /var, the LXC administration tools will detect this and automatically exploit it by cloning containers using btrfs snapshots.[^4]

#### Templating configuration

You can template configuration if you want to simplify your deployments. It could be useful if you need to do specific lxc configuration. To do it, simply create a file (name it as you want) and add your lxc configuration (here the network configuration):

```bash
# /etc/lxc/lxc-nat.conf
## Network
lxc.network.type = veth
lxc.network.flags = up

# Network host side
lxc.network.link = virbr0
lxc.network.veth.pair = veth-$name

# Network container side
lxc.network.name = eth0
lxc.network.ipv4 = 0.0.0.0/24
```

Then you could call it when you'll create a container with -f argument. You can create as many configuration as you want and place them were you want. I did it in /etc/lxc as I felt it well.

### List containers

You can list containers:

```bash
> lxc-list
RUNNING
  mycontainer

FROZEN

STOPPED
  mycontainer2
```

If you want to list running containers:

```bash
lxc-ls --active
```

### Start a container

To start a container as a deamon:

```bash
lxc-start -n mycontainer -d
```

- d: Run the container as a daemon. As the container has no more tty, if an error occurs nothing will be displayed, the log file can be used to check the error.

### Plug to console

You can connect to the console:

```bash
lxc-console -n mycontainer
```

If you've problems to connect to the container, [do this](#cant-connect-to-console).

### Stop a container

To stop a container:

```bash
lxc-shutdown -n mycontainer
```

or

```bash
lxc-halt -n mycontainer
```

{{< alert context="info" text="You can't lxc-halt/lxc-shutdown on a container based on LVM in the current Debian version(Wheezy)" />}}

### Force shutdown

If you need to force a container to halt:

```bash
lxc-stop -n mycontainer
```

### Autostart on boot

If you need to get LXC containers to autostart on boot, you'll need to create symlink:

```bash
ln -s /var/lib/lxc/mycontainer/config /etc/lxc/auto/mycontainer
```

### Delete a container

You can delete a container like this:

```bash
lxc-destroy -n mycontainer
```

{{< alert context="danger" text="This will remove all your data as well. Do a backup before doing destroy!" />}}

### Monitoring

If you want to know the state of a container:

```bash
lxc-info -n mycontainer
state:   RUNNING
pid:      3034
```

Available state are:

- ABORTING
- RUNNING
- STARTING
- STOPPED
- STOPPING

### Freeze

If you want to freeze (suspend like) a container:

```bash
lxc-freeze -n mycontainer
```

### Unfreeze/Restore

If you want to unfreeze (resume/restore like) a container:

```bash
lxc-unfreeze -n mycontainer
```

### Monitor changes

You can monitor changes of a container with lxc-monitor:

```bash
> lxc-monitor -n mycontainer
'mycontainer' changed state to [STOPPING]
'mycontainer' changed state to [STOPPED]
'mycontainer' changed state to [STARTING]
'mycontainer' changed state to [RUNNING]
```

You can see all container changing states.

### Trigger changes

You can also use 'lxc-wait command with '-s' parameter to wait a specific state and execute something afterward:

```bash
lxc-wait -n mycontainer -s STOPPED && echo "Container stopped" | mail -s 'You need to restart it' xxx@mycompany.com
```

### Launch command in a running container

You can launch a command in a running container without being inside it:

```bash
lxc-attach -n mycontainer -- /etc/init.d/cron restart
```

This restart the cron service in "mycontianer" container.

## Convert/Migrate a VM/Host to a LXC container

If you already have a running machine on KVM/VirtualBox or anything else and want to convert to an LXC container, it's easy. I've wrote a script (strongly inspired from the lxc-create) that helps me to initiate the missing elements. You can copy it in /usr/bin folder (`lxc-convert`).

```bash
#!/bin/bash

#
# lxc: linux Container library

# Authors:
# Pierre MAVRO <xxx@company.com>

# This library is free software; you can redistribute it and/or
# modify it under the terms of the GNU Lesser General Public
# License as published by the Free Software Foundation; either
# version 2.1 of the License, or (at your option) any later version.

# This library is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
 # MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
# Lesser General Public License for more details.

# You should have received a copy of the GNU Lesser General Public
# License along with this library; if not, write to the Free Software
# Foundation, Inc., 59 Temple Place, Suite 330, Boston, MA 02111-1307 USA

configure_debian()
{
    rootfs=$1
    hostname=$2

    # Remove unneeded folders
    rm -Rf $rootfs/{dev,proc,sys,run}

    # Decompress dev devices
    tar -xzf /usr/share/debootstrap/devices.tar.gz -C $rootfs
    rootfs_dev=$rootfs/dev

    # Create missing dev devices
    mkdir -m 755 $rootfs_dev/pts
    mkdir -m 1777 $rootfs_dev/shm
    mknod -m 666 $rootfs_dev/tty0 c 4 0
    mknod -m 600 $rootfs_dev/initctl p

    # Create folders
    mkdir -p $rootfs/{proc,sys,run}

    # Do not use fstab
    mv $rootfs/etc/fstab{,.old}
    touch $rootfs/etc/fstab

    # squeeze only has /dev/tty and /dev/tty0 by default,
    # therefore creating missing device nodes for tty1-4.
    for tty in $(seq 1 4); do
    if [ ! -e $rootfs/dev/tty$tty ]; then
        mknod $rootfs/dev/tty$tty c 4 $tty
    fi
    done

    # configure the inittab
    cat <<EOF > $rootfs/etc/inittab
id:2:initdefault:
si::sysinit:/etc/init.d/rcS
l0:0:wait:/etc/init.d/rc 0
l1:1:wait:/etc/init.d/rc 1
l2:2:wait:/etc/init.d/rc 2
l3:3:wait:/etc/init.d/rc 3
l4:4:wait:/etc/init.d/rc 4
l5:5:wait:/etc/init.d/rc 5
l6:6:wait:/etc/init.d/rc 6
# Normally not reached, but fallthrough in case of emergency.
z6:6:respawn:/sbin/sulogin
1:2345:respawn:/sbin/getty 38400 console
c1:12345:respawn:/sbin/getty 38400 tty1 linux
c2:12345:respawn:/sbin/getty 38400 tty2 linux
c3:12345:respawn:/sbin/getty 38400 tty3 linux
c4:12345:respawn:/sbin/getty 38400 tty4 linux
EOF

    # add daemontools-run entry
    if [ -e $rootfs/var/lib/dpkg/info/daemontools.list ]; then
        cat <<EOF >> $rootfs/etc/inittab
#-- daemontools-run begin
SV:123456:respawn:/usr/bin/svscanboot
#-- daemontools-run end
EOF
    fi

    # Remove grub and kernel
    chroot $rootfs apt-get --yes -o  Dpkg::Options::="--force-confdef" -o Dpkg::Options::="--force-confold" remove grub grub2 grub-pc grub-common linux-image-amd64

    # remove pointless services in a container
    chroot $rootfs "LANG=C /usr/sbin/update-rc.d -f checkroot.sh remove" # S
    chroot $rootfs "LANG=C /usr/sbin/update-rc.d checkroot.sh stop 09 S ."

    chroot $rootfs "LANG=C /usr/sbin/update-rc.d -f umountfs remove" # 0 6
    chroot $rootfs "LANG=C /usr/sbin/update-rc.d umountfs start 09 0 6 ."

    chroot $rootfs "LANG=C /usr/sbin/update-rc.d -f umountroot remove" # 0 6
    chroot $rootfs "LANG=C /usr/sbin/update-rc.d umountroot start 10 0 6 ."

    # The following initscripts don't provide an empty start or stop block.
    # To prevent them being enabled on upgrades, we leave a start link on
    # runlevel 3.
    chroot $rootfs "LANG=C /usr/sbin/update-rc.d -f hwclock.sh remove" # S 0 6
    chroot $rootfs "LANG=C /usr/sbin/update-rc.d hwclock.sh start 10 3 ."

    chroot $rootfs "LANG=C /usr/sbin/update-rc.d -f hwclockfirst.sh remove" # S
    chroot $rootfs "LANG=C /usr/sbin/update-rc.d hwclockfirst start 08 3 ."

    chroot $rootfs "LANG=C /usr/sbin/update-rc.d -f module-init-tools remove" # S
    chroot $rootfs "LANG=C /usr/sbin/update-rc.d module-init-tools start 10 3 ."

    return 0
}

copy_configuration()
{
    path=$1
    rootfs=$2
    name=$3

    cat <<EOF > $path/config
# $path/config

## Container
lxc.utsname                             = $hostname
lxc.tty                                 = 4
lxc.pts                                 = 1024
#lxc.console                            = /var/log/lxc/$name.console

## Capabilities
#lxc.cap.drop                           = mac_admin
#lxc.cap.drop                           = mac_override
lxc.cap.drop                            = sys_admin
#lxc.cap.drop                           = sys_module

## Devices
# Allow all devices
#lxc.cgroup.devices.allow = a
# Deny all devices
lxc.cgroup.devices.deny = a
# /dev/null and zero
lxc.cgroup.devices.allow = c 1:3 rwm
lxc.cgroup.devices.allow = c 1:5 rwm
# /dev/consoles
lxc.cgroup.devices.allow = c 5:1 rwm
# /dev/tty
lxc.cgroup.devices.allow = c 5:0 rwm
lxc.cgroup.devices.allow = c 4:0 rwm
lxc.cgroup.devices.allow = c 4:1 rwm
# /dev/{,u}random
lxc.cgroup.devices.allow = c 1:9 rwm
# /dev/random
lxc.cgroup.devices.allow = c 1:8 rwm
# /dev/pts/*
lxc.cgroup.devices.allow = c 136:* rwm
# /dev/ptmx
lxc.cgroup.devices.allow = c 5:2 rwm
# /dev/rtc
lxc.cgroup.devices.allow = c 254:0 rwm
# /dev/fuse
lxc.cgroup.devices.allow = c 10:229 rwm

## Limits
#lxc.cgroup.cpu.shares                  = 1024
#lxc.cgroup.cpuset.cpus                 = 0
#lxc.cgroup.memory.limit_in_bytes       = 256M
#lxc.cgroup.memory.memsw.limit_in_bytes = 1G
#lxc.cgroup.blkio.weight                = 500


## Filesystem
lxc.mount.entry                         = proc $rootfs/proc proc nodev,noexec,nosuid 0 0
lxc.mount.entry                         = sysfs $rootfs/sys sysfs defaults,ro 0 0
lxc.rootfs                              = $rootfs
# LVM
#lxc.rootfs                             = /dev/vg/lvname
EOF

    # Adding shared data directory if existing
    if [ -d /srv/share/$hostname ]; then
    echo "lxc.mount.entry                         = /srv/share/$hostname $rootfs/srv/$hostname none defaults,bind 0 0" >> $path/config
    else
    echo "#lxc.mount.entry                        = /srv/share/$hostname $rootfs/srv/$hostname none defaults,bind 0 0" >> $path/config
    fi

    gen_mac=`openssl rand -hex 6 | sed 's/\(..\)/\1:/g; s/.$//'`
    cat >> $path/config << EOF

#lxc.mount.entry                        = /srv/$hostname $rootfs/srv/$hostname none defaults,bind 0 0

## Network
lxc.network.type                        = veth
lxc.network.flags                       = up
#lxc.network.hwaddr                      = $gen_mac
lxc.network.link                        = lxcbr0
lxc.network.name                        = eth0
lxc.network.veth.pair                   = veth-$hostname
EOF

    if [ $? -ne 0 ]; then
    echo "Failed to add configuration"
    return 1
    fi

    return 0
}


usage()
{
    cat <<EOF
$1 -h|--help -p|--path=<path> -n|--name=name
EOF
    return 0
}

options=$(getopt -o hp:n:c -l help,path:,name:,clean -- "$@")
if [ $? -ne 0 ]; then
        usage $(basename $0)
    exit 1
fi
eval set -- "$options"

while true
do
    case "$1" in
        -h|--help)      usage $0 && exit 0;;
        -p|--path)      path=$2; shift 2;;
        -n|--name)      name=$2; shift 2;;
        --)             shift 1; break ;;
        *)              break ;;
    esac
done

if [ ! -z "$clean" -a -z "$path" ]; then
    clean || exit 1
    exit 0
fi

if [ -z "$path" ]; then
    echo "'path' parameter is required"
    exit 1
fi

if [ "$(id -u)" != "0" ]; then
    echo "This script should be run as 'root'"
    exit 1
fi

rootfs=$path/rootfs

configure_debian $rootfs $name
if [ $? -ne 0 ]; then
    echo "failed to configure debian for a container"
    exit 1
fi

copy_configuration $path $rootfs
if [ $? -ne 0 ]; then
    echo "failed write configuration file"
    exit 1
fi
```

To use it, it's easy. First of all mount or copy all your datas in the rootfs folder, be sure to have enough space, then launch the lxc-convert script like in this example :

```bash
migrated_container=migrated_container
mkdir -p /var/lib/lxc/$migrated_container/rootfs
rsync -e ssh -a --exclude '/dev' --exclude '/proc' --exclude '/sys' <old_vm>:/ /var/lib/lxc/$migrated_container/rootfs
lxc-convert -p /var/lib/lxc/$migrated_container -n $migrated_container
```

Adapt the remote host to your distant SSH host or rsync without SSH if it's possible. During the transfer, you need to exclude some folders to avoid errors (/proc, /sys, /dev). They will be recreated during the lxc-convert.

Then you'll be able to start it :-)

## Container configuration

Once you've initialized your container, there are a lot of interesting options. Here are some for a classical configuration (`/var/lib/lxc/mycontainer/config`):

```
## Container
# Container name
lxc.utsname                             = mycontainer
# Path where default container is based
lxc.rootfs                              = /var/lib/lxc/mycontainer/rootfs
# Set architecture type
lxc.arch                                = x86_64
#lxc.console                            = /var/log/lxc/mycontainer.console
# Number of tty/pts available for that container
lxc.tty                                 = 6
lxc.pts                                 = 1024

## Capabilities
lxc.cap.drop                            = mac_admin
lxc.cap.drop                            = mac_override
lxc.cap.drop                            = sys_admin
lxc.cap.drop                            = sys_module
## Devices
# Allow all devices
#lxc.cgroup.devices.allow               = a
# Deny all devices
lxc.cgroup.devices.deny                 = a
# Allow to mknod all devices (but not using them)
lxc.cgroup.devices.allow                = c *:* m
lxc.cgroup.devices.allow                = b *:* m

# /dev/console
lxc.cgroup.devices.allow                = c 5:1 rwm
# /dev/fuse
lxc.cgroup.devices.allow                = c 10:229 rwm
# /dev/null
lxc.cgroup.devices.allow                = c 1:3 rwm
# /dev/ptmx
lxc.cgroup.devices.allow                = c 5:2 rwm
# /dev/pts/*
lxc.cgroup.devices.allow                = c 136:* rwm
# /dev/random
lxc.cgroup.devices.allow                = c 1:8 rwm
# /dev/rtc
lxc.cgroup.devices.allow                = c 254:0 rwm
# /dev/tty
lxc.cgroup.devices.allow                = c 5:0 rwm
# /dev/urandom
lxc.cgroup.devices.allow                = c 1:9 rwm
# /dev/zero
lxc.cgroup.devices.allow                = c 1:5 rwm

## Limits
#lxc.cgroup.cpu.shares                  = 1024
#lxc.cgroup.cpuset.cpus                 = 0
#lxc.cgroup.memory.limit_in_bytes       = 256M
#lxc.cgroup.memory.memsw.limit_in_bytes = 1G

## Filesystem
# fstab for the containers with advanced features like bindind mount.
# Mount bind between host and containers (mount --bind equivalent)
lxc.mount.entry                         = proc /var/lib/lxc/mycontainer/rootfs/proc proc nodev,noexec,nosuid 0 0
lxc.mount.entry                         = sysfs /var/lib/lxc/mycontainer/rootfs/sys sysfs defaults,ro 0 0
#lxc.mount.entry                        = /srv/mycontainer /var/lib/lxc/mycontainer/rootfs/srv/mycontainer none defaults,bind 0 0

## Network
lxc.network.type                        = veth
lxc.network.flags                       = up
lxc.network.hwaddr                      = 11:22:33:44:55:66
lxc.network.link                        = br0
lxc.network.mtu                         = 1500
lxc.network.name                        = eth0
lxc.network.veth.pair                   = veth-$name
```

For an LVM configuration:

```
lxc.tty = 4
lxc.pts = 1024
lxc.utsname = mycontainer

## Capabilities
lxc.cap.drop                            = sys_admin

# When using LXC with apparmor, uncomment the next line to run unconfined:
#lxc.aa_profile = unconfined

lxc.cgroup.devices.deny = a
# /dev/null and zero
lxc.cgroup.devices.allow = c 1:3 rwm
lxc.cgroup.devices.allow = c 1:5 rwm
# consoles
lxc.cgroup.devices.allow = c 5:1 rwm
lxc.cgroup.devices.allow = c 5:0 rwm
lxc.cgroup.devices.allow = c 4:0 rwm
lxc.cgroup.devices.allow = c 4:1 rwm
# /dev/{,u}random
lxc.cgroup.devices.allow = c 1:9 rwm
lxc.cgroup.devices.allow = c 1:8 rwm
lxc.cgroup.devices.allow = c 136:* rwm
lxc.cgroup.devices.allow = c 5:2 rwm
# rtc
lxc.cgroup.devices.allow = c 254:0 rwm

# mounts point
lxc.mount.entry = proc proc proc nodev,noexec,nosuid 0 0
lxc.mount.entry = sysfs sys sysfs defaults  0 0
lxc.rootfs = /dev/vg/mycontainer
```

### Architectures

You can set the container architecture on a container. For example, you can use an x86 container on a x64 kernel (`/var/lib/lxc/mycontainer/config`):

```
lxc.arch=x86
```

### Capabilities

You can specify the capability to be dropped in the container. A single line defining several capabilities with a space separation is allowed. The format is the lower case of the capability definition without the "CAP\_" prefix, eg. CAP_SYS_MODULE should be specified as sys_module. You can see the complete list of linux capabilities with explanations by reading the man page :

```
man 7 capabilities
```

### Devices

You can manage (allow/deny) accessible devices directly from your containers. By default, everything is disabled (`/var/lib/lxc/mycontainer/config`):

```
## Devices
# Allow all devices
#lxc.cgroup.devices.allow               = a
# Deny all devices
lxc.cgroup.devices.deny                 = a
```

- a : means all devices

You can then allow some of them easily[^5] (`/var/lib/lxc/mycontainer/config`):

```
lxc.cgroup.devices.allow = c 5:1 rwm # dev/console
lxc.cgroup.devices.allow = c 5:0 rwm # dev/tty
lxc.cgroup.devices.allow = c 4:0 rwm # dev/tty0
```

To get the complete list of allowed devices (`lxc-cgroup`):

```
lxc-cgroup -n mycontainer devices.list
```

To get a better understanding of this, here are explanations[^6]:

```
lxc.cgroup.devices.allow = <type> <major>:<minor> <perm>
```

- <type> : b (block), c (char), etc ...
- <major> : major number
- <minor> : minor number (wildcard is accepted)
- <perms> : r (read), w (write), m (mapping)

### Container limits (cgroups)

If you've never played with Cgroups, [look at my documentation]({{< ref "docs/Linux/Kernel/process_latency_and_kernel_timing.md" >}}). With LXC, here are the available ways to setup cgroups to your containers :

- You can change cgroups values with lxc-cgroup command (on the fly):

```bash
lxc-cgroup -n <container_name> <cgroup-name> <value>
```

- You can directly play with /proc (on the fly):

```bash
echo <value> > /sys/fs/cgroup/lxc/<container_name>/<cgroup-name>
```

- And set it directly in the config file (persistent way) (`/var/lib/lxc/mycontainer/config`):

```
lxc.cgroup.<cgroup-name> = <value>
```

Cgroups can be changed on the fly.

{{< alert context="danger" text="You should warn when you reduce some of them, especially the memory (be sure that you do not reduce more than used)." />}}

If you want to see available cgroups for a container :

```bash
> ls -1 /sys/fs/cgroup/lxc/<container_name>
blkio.io_merged
blkio.io_queued
blkio.io_service_bytes
blkio.io_serviced
blkio.io_service_time
blkio.io_wait_time
blkio.reset_stats
blkio.sectors
blkio.time
blkio.weight
blkio.weight_device
cgroup.clone_children
cgroup.event_control
cgroup.procs
cpuacct.stat
cpuacct.usage
cpuacct.usage_percpu
cpuset.cpu_exclusive
cpuset.cpus
cpuset.mem_exclusive
cpuset.mem_hardwall
cpuset.memory_migrate
cpuset.memory_pressure
cpuset.memory_spread_page
cpuset.memory_spread_slab
cpuset.mems
cpuset.sched_load_balance
cpuset.sched_relax_domain_level
cpu.shares
devices.allow
devices.deny
devices.list
freezer.state
memory.failcnt
memory.force_empty
memory.limit_in_bytes
memory.max_usage_in_bytes
memory.memsw.failcnt
memory.memsw.limit_in_bytes
memory.memsw.max_usage_in_bytes
memory.memsw.usage_in_bytes
memory.move_charge_at_immigrate
memory.numa_stat
memory.oom_control
memory.soft_limit_in_bytes
memory.stat
memory.swappiness
memory.usage_in_bytes
memory.use_hierarchy
net_cls.classid
notify_on_release
tasks
```

#### CPU

##### CPU Pining

If you want to bind some CPU/Cores to a VM/Container, there is a solution called CPU Pining :). First, look at the available cores on your server :

```
> grep -e processor -e core /proc/cpuinfo | sed 's/processor/\nprocessor/'

processor	: 0
core id		: 0
cpu cores	: 4

processor	: 1
core id		: 1
cpu cores	: 4

processor	: 2
core id		: 2
cpu cores	: 4

processor	: 3
core id		: 3
cpu cores	: 4

processor	: 4
core id		: 0
cpu cores	: 4

processor	: 5
core id		: 1
cpu cores	: 4

processor	: 6
core id		: 2
cpu cores	: 4

processor	: 7
core id		: 3
cpu cores	: 4
```

You can see there are 7 cores (called processor). In fact there are 4 cores with 2 thread each on this CPU. That's why there are 4 cores id and 8 detected cores.

So here is the list of the cores with their attached core:

```
core id 0 : processors 0 and 4
core id 1 : processors 1 and 5
core id 2 : processors 2 and 6
core id 3 : processors 3 and 7
```

Now if I want to dedicate a set of cores to a container (`/var/lib/lxc/mycontainer/config`):

```
lxc.cgroup.cpuset.cpus = 0-3
```

This will add the 4 firsts cores to the container. Then if I only want CPU 1 and 3 (`/var/lib/lxc/mycontainer/config`):

```
lxc.cgroup.cpuset.cpus = 1,3
```

###### Check CPU assignments

You can check how many CPU Pining are working for a container. On the host file, launch "htop" for example and in the container launch stress :

```bash
aptitude install stress
stress --cpu 2 --timeout 10s
```

This will stress 2 CPU at 100% for 10 seconds. You'll see your htop CPU bars at 100%. If I change 2 by 3 and only binded 2 CPUs, only 2 will be at 100% :-)

##### Scheduler

This is the other method to assign CPU to a container. You need to add weight to VMs so that the scheduler can decide which container should use CPU time form the CPU clock. For instance, if a container is set to 512 and another to 1024, the last one will have twice more CPU time than the first container. To edit this property (`/var/lib/lxc/mycontainer/config`):

```
lxc.cgroup.cpu.shares = 512
```

If you need more documentation, look at the kernel page[^7].

#### Memory

You can limit the memory in a container like this (`/var/lib/lxc/mycontainer/config`):

```
lxc.cgroup.memory.limit_in_bytes = 128M
```

If you've got error when trying to limit memory, [check the FAQ](#cant-limit-container-memory).

```
lxc.cgroup.memory.limit_in_bytes = 128M
```

##### Check memory from the host

You can check memory from the host like that[^8] :

- Current memory usage:

```bash
cat /sys/fs/cgroup/lxc/mycontainer/memory.usage_in_bytes
```

- Current memory + swap usage:

```bash
cat /sys/fs/cgroup/lxc/mycontainer/memory.memsw.usage_in_bytes
```

- Maximum memory usage:

```bash
cat /sys/fs/cgroup/lxc/mycontainer/memory.max_usage_in_bytes
```

- Maximum memory + swap usage :

```bash
cat /sys/fs/cgroup/lxc/mycontainer/memory.memsw.max_usage_in_bytes
```

Here is an easier solution to read informations :

```bash
awk '{ printf "%sK\n", $1/ 1024 }' /sys/fs/cgroup/lxc/mycontainer/memory.usage_in_bytes
awk '{ printf "%sM\n", $1/ 1024 / 1024 }' /sys/fs/cgroup/lxc/mycontainer/memory.usage_in_bytes
```

##### Check memory in the container

The actual problem is you can't check how many memory you've set and is available for your container. For the moment /proc/meminfo is not correctly updated[^9]. If you need to validate the available memory on a container, you have to write fake data into the allocated memory area to trigger the memory checks of the kernel/visualization tool.

Memory overcommit is a Linux kernel feature that lets applications allocate more memory than is actually available. The idea behind this feature is that some applications allocate large amounts of memory just in case, but never actually use it. Thus, memory overcommit allows you to run more applications than actually fit in your memory, provided the applications don’t actually use the memory they have allocated. If they do, then the kernel (via OOM killer) terminates the application.

Here is the code[^10] (`memory_allocation.c`):

```c
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

main() {

   int i;
   for(i=0;i<9000;i++) {
        int *ptr = malloc(i*1024*1024);
        if (ptr == NULL) {
          printf ("Soft memory allocation failed for %i MB\n",i );
          break;
        } else {
            free(ptr);
            ptr = NULL;
        }
   }


   for(i=0;i<9000;i++) {

        int *ptr = malloc(i*1024*1024);
        if (ptr == NULL) {
          printf ("Memory allocation failed for %i MB\n",i );
          break;
        } else {
            memset(ptr, 0, i*1024*1024);
            printf("Wrote %i MB to memory\n", i);
            //usleep(1000);
            free(ptr);
            ptr = NULL;
        }

   }
}
```

Then compil it (with gcc) :

```bash
aptitude install gcc
gcc memory_allocation.c -o memory_allocation
```

You can now run the test :

```bash
> ./memory_allocation
Soft memory allocation failed for 654 MB
Wrote 0 MB to memory
Wrote 1 MB to memory
Wrote 2 MB to memory
Wrote 3 MB to memory
[...]
Wrote 185 MB to memory
Wrote 186 MB to memory
Wrote 187 MB to memory
Wrote 188 MB to memory
Killed
```

#### SWAP

You can limit the swap in a container like this (`/var/lib/lxc/mycontainer/config`):

```
lxc.cgroup.memory.memsw.limit_in_bytes = 192M
```

{{< alert context="danger" text="This limit is not only SWAP but Memory + SWAP" />}}

That mean that **"lxc.cgroup.memory.memsw.limit_in_bytes" should be at least equal to "lxc.cgroup.memory.limit_in_bytes".**

If you've got error when trying to limit swap, check the [FAQ](#cant-limit-container-memory).

#### Disks

By default, LXC doesn't provide any disks limitation. Anyway, there are enough solution today to make that kind of limitations:

- [LVM]({{< ref "docs/Linux/FilesystemsAndStorage/lvm_working_with_logical_volume_management.md" >}}): create one LV per container
- [BTRFS]({{< ref "docs/Linux/FilesystemsAndStorage/btrfs-using-the-ext4-replacement.md">}}): using integrated BTRFS quotas
- [ZFS]({{< ref "docs/Solaris/Filesystems/zfs_the_filesystem_par_excellence.md">}}): if you're using ZFS on Linux, you can use integrated zfs/zpool quotas
- [Quotas]({{< ref "docs/Linux/FilesystemsAndStorage/setting_up_quotas_on_linux.md" >}}): using classical Linux quotas (not the recommended solution)
- Disk image: you can use QCOW/QCOW2/RAW/QED images

##### Mount

{{< alert context="warning" text="You should take care if you want to create a mount entry in a subdirectory of /mnt." />}}

It won't work so easily. The reason this happens is that by default 'mnt' is the directory used as pivotdir, where the old_root is placed during pivot_root(). After that, everything under pivotdir is unmounted.

A workaround is to specify an alternate 'lxc.pivotdir' in the container configuration file.[^11]

###### Block Device

You can mount block devices in adding in your container configuration lines like this (adapt with your needs) (`/var/lib/lxc/mycontainer/config`):

```
lxc.mount.entry = /dev/sdb1 /var/lib/lxc/mycontainer/rootfs/mnt ext4 rw 0 2
```

###### Bind mount

You also can mount bind mountpoints like that (adapt with your needs) (`/var/lib/lxc/mycontainer/config`):

```
lxc.mount.entry = /path/in/host/mount_point /var/lib/lxc/mycontainer/rootfs/mount_moint none bind 0 0
```

##### Disk priority

You can set disk priority like that (default is 500) (`/var/lib/lxc/mycontainer/config`):

```
lxc.cgroup.blkio.weight = 500
```

Higher the value is, more the priority will be important. You can get more informations here. Maximum value is 1000 and lowest is 10.

{{< alert context="info" text="You need to have CFQ scheduler to make it work properly" />}}

##### Disk bandwidth

Another solution is to limit bandwidth usage, but the Wheezy kernel doesn't have the "CONFIG_BLK_DEV_THROTTLING" activated. You need to take a testing/unstable kernel instead or recompile a new one with this option activated. To do this follow the kernel procedure.

Then, you'll be able to limit bandwidth like that (`/var/lib/lxc/mycontainer/config`):

```
# Limit to 1Mb/s
lxc.cgroup.blkio.throttle.read_bps_device = 100
```

#### Network

You can limit network bandwidth using native kernel QOS directly on cgroups. For example, we have 2 containers : A and B. To get a good understanding, look at this schema:

![](/images/cgroup_network.avif)[^12]

Now you've understand how it could looks like. Now if I want to limit a container to 30Mb and the other one to 40Mb, here is how I should achieve it. Assign IDs on containers that should have quality of service :

```bash
echo 0x1001 > /sys/fs/cgroup/lxc/<containerA>/net_cls.classid
echo 0x1002 > /sys/fs/cgroup/lxc/<containerB>/net_cls.classid
```

- 0x1001 : corresponding to 10:1
- 0x1002 : corresponding to 10:2

Then select the desired QOS algorithm (HTB) :

```bash
tc qdisc add dev eth0 root handle 10: htb
```

Choose the desired bandwidth on containers IDs :

```bash
tc class add dev eth0 parent 10: classid 10:1 htb rate 40mbit
tc class add dev eth0 parent 10: classid 10:2 htb rate 30mbit
```

Enable filtering :

```bash
tc filter add dev eth0 parent 10: protocol ip prio 10 handle 1: cgroup
```

### Resources statistics

Unfortunately, you can't have informations directly on the containers, however you can have informations from the host. Here is a little script to do it (`/usr/bin/lxc-resources-stats`):

```bash
#!/bin/bash

cd /sys/fs/cgroup/lxc/
for i in * ; do
    if [ -d $i ] ; then
        echo "===== $i ====="
        echo "CPU, cap:               " $(cat /sys/fs/cgroup/lxc/$i/cpuset.cpus)
        echo "CPU, shares:            " $(cat /sys/fs/cgroup/lxc/$i/cpu.shares)
        awk '{ printf "RAM, limit usage:        %sM\n", $1/ 1024/1024 }' /sys/fs/cgroup/lxc/$i/memory.limit_in_bytes
        awk '{ printf "RAM+SWAP, limit usage:   %sM\n", $1/ 1024/1024 }' /sys/fs/cgroup/lxc/$i/memory.memsw.limit_in_bytes
        awk '{ printf "RAM, current usage:      %sM\n", $1/ 1024/1024 }' /sys/fs/cgroup/lxc/$i/memory.usage_in_bytes
        awk '{ printf "RAM+SWAP, current usage: %sM\n", $1/ 1024/1024 }' /sys/fs/cgroup/lxc/$i/memory.memsw.usage_in_bytes
        awk '{ printf "RAM, max usage:          %sM\n", $1/ 1024/1024 }' /sys/fs/cgroup/lxc/$i/memory.max_usage_in_bytes
        awk '{ printf "RAM+SWAP, max usage:     %sM\n", $1/ 1024/1024 }' /sys/fs/cgroup/lxc/$i/memory.memsw.max_usage_in_bytes
        echo "DISK I/O weight:        " $(cat /sys/fs/cgroup/lxc/$i/blkio.weight)
        echo ""
    fi
done
```

Here is the result:

```bash
> lxc-resources-stats
===== mycontainer =====
CPU, cap:                3-4,6-7
CPU, shares:             1024
RAM, limit usage:        2048M
RAM+SWAP, limit usage:   3072M
RAM, current usage:      1577.33M
RAM+SWAP, current usage: 1582.79M
RAM, max usage:          2048M
RAM+SWAP, max usage:     2060.5M
DISK I/O weight:         500
```

## FAQ

### How could I know if I'm in a container or not?

There's an easy way to know that :

```bash
> cat /proc/$$/cgroup
1:perf_event,blkio,net_cls,freezer,devices,memory,cpuacct,cpu,cpuset:/lxc/mycontainer
```

You can see in the cpuset, the container name where I am ("mycontainer" here).

### Can't connect to console

If you want to plug yourself to a container through the console, your first need to create devices :

```bash
chroot /var/lib/lxc/mycontainer/rootfs
mknod -m 666 /dev/tty1 c 4 1
mknod -m 666 /dev/tty2 c 4 2
mknod -m 666 /dev/tty3 c 4 3
mknod -m 666 /dev/tty4 c 4 4
mknod -m 666 /dev/tty5 c 4 5
mknod -m 666 /dev/tty6 c 4 6
```

Then you'll be able to connect:

```
lxc-console -n mycontainer
```

### Can't create a LXC LVM container

If you get this kind of error during LVM :

```
Copying local cache to /var/lib/lxc/mycontainerlvm/rootfs.../usr/share/lxc/templates/lxc-debian: line 101: /var/lib/lxc/mycontainerlvm/rootfs/etc/apt/sources.list.d/debian.list: No such file or directory
/usr/share/lxc/templates/lxc-debian: line 107: /var/lib/lxc/mycontainerlvm/rootfs/etc/apt/sources.list.d/debian.list: No such file or directory
/usr/share/lxc/templates/lxc-debian: line 111: /var/lib/lxc/mycontainerlvm/rootfs/etc/apt/sources.list.d/debian.list: No such file or directory
/usr/share/lxc/templates/lxc-debian: line 115: /var/lib/lxc/mycontainerlvm/rootfs/etc/apt/sources.list.d/debian.list: No such file or directory
/usr/share/lxc/templates/lxc-debian: line 183: /var/lib/lxc/mycontainerlvm/rootfs/etc/fstab: No such file or directory
mount: mount point /var/lib/lxc/mycontainerlvm/rootfs/dev/pts does not exist
mount: mount point /var/lib/lxc/mycontainerlvm/rootfs/proc does not exist
mount: mount point /var/lib/lxc/mycontainerlvm/rootfs/sys does not exist
mount: mount point /var/lib/lxc/mycontainerlvm/rootfs/var/cache/apt/archives does not exist
/usr/share/lxc/templates/lxc-debian: line 49: /var/lib/lxc/mycontainerlvm/rootfs/etc/dpkg/dpkg.cfg.d/lxc-debconf: No such file or directory
/usr/share/lxc/templates/lxc-debian: line 55: /var/lib/lxc/mycontainerlvm/rootfs/usr/sbin/policy-rc.d: No such file or directory
chmod: cannot access `/var/lib/lxc/mycontainerlvm/rootfs/usr/sbin/policy-rc.d': No such file or directory
chroot: failed to run command `/usr/bin/env': No such file or directory
chroot: failed to run command `/usr/bin/env': No such file or directory
chroot: failed to run command `/usr/bin/env': No such file or directory
umount: /var/lib/lxc/mycontainerlvm/rootfs/var/cache/apt/archives: not found
chroot: failed to run command `/usr/bin/env': No such file or directory
chroot: failed to run command `/usr/bin/env': No such file or directory
chroot: failed to run command `/usr/bin/env': No such file or directory
chroot: failed to run command `/usr/bin/env': No such file or directory
chroot: failed to run command `/usr/bin/env': No such file or directory
chroot: failed to run command `/usr/bin/env': No such file or directory
umount: /var/lib/lxc/mycontainerlvm/rootfs/dev/pts: not found
umount: /var/lib/lxc/mycontainerlvm/rootfs/proc: not found
umount: /var/lib/lxc/mycontainerlvm/rootfs/sys: not found
'debian' template installed
Unmounting LVM
'mycontainerlvm' created
```

this is because of a Debian bug that the maintainer doesn't want to fix :-(. [Here is a workaround](#installation).

### Can't limit container memory or swap

If you can't limit container memory and have this kind of issue:

```bash
> lxc-cgroup -n mycontainer memory.limit_in_bytes "128M"
lxc-cgroup: cgroup is not mounted
lxc-cgroup: failed to assign '128M' value to 'memory.limit_in_bytes' for 'mycontainer'
```

This is because cgroup memory capability is not loaded from your kernel. You can check it like that :

```bash {linenos=table,hl_lines=[6],anchorlinenos=true}
> cat /proc/cgroups
#subsys_name	hierarchy	num_cgroups	enabled
cpuset	1	4	1
cpu	1	4	1
cpuacct	1	4	1
memory	0	1	0
devices	1	4	1
freezer	1	4	1
net_cls	1	4	1
blkio	1	4	1
perf_event	1	4	1
```

As we want to manage memory and swap on containers, as it's not available by default, add cgroup argument to grub to activate those functionality:

- cgroup RAM feature : `cgroup_enable=memory`
- cgroup SWAP feature : `swapaccount=1`

With grub (`/etc/default/grub`):

```ini {linenos=table,hl_lines=[9],anchorlinenos=true}
# If you change this file, run 'update-grub' afterwards to update
# /boot/grub/grub.cfg.
# For full documentation of the options in this file, see:
#   info -f grub -n 'Simple configuration'

GRUB_DEFAULT=0
GRUB_TIMEOUT=5
GRUB_DISTRIBUTOR=`lsb_release -i -s 2> /dev/null || echo Debian`
GRUB_CMDLINE_LINUX_DEFAULT="quiet cgroup_enable=memory swapaccount=1"
GRUB_CMDLINE_LINUX=""
[...]
```

Then regenerate grub config:

```bash
update-grub
```

**Now reboot to make changes available.**

After reboot you can check that memory is activated:

```bash {linenos=table,hl_lines=[6],anchorlinenos=true}
> cat /proc/cgroups
#subsys_name	hierarchy	num_cgroups	enabled
cpuset	1	6	1
cpu	1	6	1
cpuacct	1	6	1
memory	1	6	1
devices	1	6	1
freezer	1	6	1
net_cls	1	6	1
blkio	1	6	1
perf_event	1	6	1
```

Another way to check is the mount command:

```bash
> mount | grep cgroup
cgroup on /sys/fs/cgroup type cgroup (rw,relatime,perf_event,blkio,net_cls,freezer,devices,memory,cpuacct,cpu,cpuset,clone_children)
```

You can see that memory is available here[^13].

### I can't start my container, how could I debug?

You can debug a container on boot by this way:

```bash
lxc-start -n mycontainer -l debug -o debug.out
```

Now you can look at debug.out and see what's wrong.

### /usr/sbin/grub-probe: error: cannot find a device for / (is /dev mounted?)

I got dpkg error issue while I wanted to upgrade an LXC containers running on Debian because grub couldn't find /.

To resolve that issue, I needed to remove definitively grub and grub-pc. Then the system accepted to remove the kernel.

### telinit: /run/initctl: No such file or directory

If you got his kind of error when you want to properly shutdown your LXC container, you need to create a device in your container:

```bash
mknod -m 600 /var/lib/lxc/<container_name>/rootfs/run/initctl p
```

And then add this in the container configuration file (`/var/lib/lxc/<container_name>/config`):

```
lxc.cap.drop = sys_admin
```

You can now shutdown it properly without any issue :-)

### Some containers are loosing their IP addresse at boot

If you're experiencing issues with booting containers which are loosing their static IP at boot[^14] there is a solution. The first thing to do to recover is:

```bash
ifdown eth0 && ifup eth0
```

But is is a temporary solution. You in fact need to add in your LXC configuration file, the IP address with CIDR of your container (`/var/lib/lxc/<container_name>/config`):

```
lxc.network.ipv4 = 192.168.0.50/24
lxc.network.ipv4.gateway = auto
```

The automatic gateway setting is will in fact address to the container, the IP of the interface on which the container is attached. Then you have to modify your container network configuration and change static configuration to manual of eth0 interface. You should have something like this:

```
allow-hotplug eth0
iface eth0 inet manual
```

You're now ok, on next reboot the IP will be properly configured automatically by LXC and it will work anytime.

### OpenVPN

To make openvpn working, you need to allow tun devices. In the LXC configuration, simply add this (`container.conf`):

```
lxc.cgroup.devices.allow = c 10:200 rwm
```

And in the container, create it:

```bash
mkdir /dev/net
mknod /dev/net/tun c 10 200
chmod 0666 /dev/net/tun
```

### LXC inception or Docker in LXC

To get Docker in LXC or LXC in LXC working, you need to have some packages installed inside the LXC container:

```bash
apt-get install cgroup-bin libcgroup1 cgroupfs-mount
```

In the container configuration, you also need to have that line (`config`):

```
lxc.mount.auto = cgroup
```

Then it's ok :-)

### LXC control device mapper

In Docker, you may want to use devicemapper driver. To get it working, you need to let your LXC container to control devicemappers. To do so, just add those 1 lines in your container configuration:

```
lxc.cgroup.devices.allow = c 10:236 rwm
lxc.cgroup.devices.allow = b 252:* rwm
```

## References

- http://www.pointroot.org/index.php/2013/05/12/installation-du-systeme-de-virtualisation-lxc-linux-containers-sur-debian-wheezy/
- http://box.matto.nl/lxconlaptop.html
- https://help.ubuntu.com/lts/serverguide/lxc.html
- http://debian-handbook.info/browse/stable/sect.virtualization.html
- http://www.fitzdsl.net/2012/12/installation-dun-conteneur-lxc-sur-dedibox/
- http://freedomboxblog.nl/installing-lxc-dhcp-and-dns-on-my-freedombox/
- http://containerops.org/2013/11/19/lxc-networking/

[^1]: http://lxc.sourceforge.net/
[^2]: http://pi.lastr.us/doku.php/virtualizacion:lxc:digitalocean-wheezy
[^3]: http://philpep.org/blog/lxc-sur-debian-squeeze
[^4]: https://help.ubuntu.com/lts/serverguide/lxc.html
[^5]: http://lwn.net/Articles/273208/
[^6]: http://wiki.rot13.org/rot13/index.cgi?action=display_html;page_name=lxc
[^7]: https://www.kernel.org/doc/Documentation/scheduler/sched-design-CFS.txt
[^8]: http://www.mattfischer.com/blog/?p=399
[^9]: http://webcache.googleusercontent.com/search?q=cache:vWmLMNBRKIYJ:comments.gmane.org/gmane.linux.kernel.containers/23094+&cd=6&hl=fr&ct=clnk&gl=fr&client=firefox-a
[^10]: http://www.jotschi.de/Uncategorized/2010/11/11/memory-allocation-test.html
[^11]: https://bugs.launchpad.net/ubuntu/+source/lxc/+bug/986385
[^12]: http://vger.kernel.org/netconf2009_slides/Network%20Control%20Group%20Whitepaper.odt
[^13]: http://vin0x64.fr/2012/01/debian-limite-de-memoire-sur-conteneur-lxc/
[^14]: http://serverfault.com/questions/571714/setting-up-bridged-lxc-containers-with-static-ips/586577#586577
