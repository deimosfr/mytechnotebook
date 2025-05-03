---
weight: 999
url: "/KVM_\\:_Mise_en_place_de_KVM/"
title: "KVM: Setting Up KVM"
description: "A comprehensive guide on setting up and configuring KVM virtualization on Linux systems, including performance optimization, networking, VM management, and advanced features."
categories: ["Virtualization", "Linux", "System Administration"]
date: "2015-03-03T09:41:00+02:00"
lastmod: "2015-03-03T09:41:00+02:00"
tags:
  ["KVM", "Virtualization", "Linux", "Network", "LVM", "QEMU", "Performance"]
toc: true
---

![KVM](/images/kvm-logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 0.12.5 |
| **Operating System** | Debian 6 |
| **Website** | [KVM Website](https://www.linux-kvm.org/) |
| **Last Update** | 03/03/2015 |
{{< /table >}}

## Introduction

KVM (for Kernel-based Virtual Machine) is a full virtualization solution for Linux on x86 hardware containing virtualization extensions (Intel VT or AMD-V). It consists of a loadable kernel module, kvm.ko, that provides the core virtualization infrastructure and a processor specific module, kvm-intel.ko or kvm-amd.ko. KVM also requires a modified QEMU although work is underway to get the required changes upstream.

Using KVM, one can run multiple virtual machines running unmodified Linux or Windows images. Each virtual machine has private virtualized hardware: a network card, disk, graphics adapter, etc.

The kernel component of KVM is included in mainline Linux, as of 2.6.20.

KVM is open source software.

## Install

To install and run KVM on Debian, follow these steps:

Run these commands as root:

```bash
aptitude update
aptitude install kvm qemu bridge-utils libvirt-bin virtinst
```

- qemu is necessary, this is the base
- kvm is for full acceleration (need processor with vmx or svm technology)
- bridge-utils are tools for bridging VMs network
- libvirt is a managing solution for your VMs

Depending on if you are using an AMD or Intel processor, run one of these commands:

```bash
modprobe kvm-amd
```

or

```bash
modprobe kvm-intel
```

Then add it to `/etc/modules`:

```bash
kvm-intel
```

## Configuration

### System performances

To get as performances as possible, we'll need to set some specific options.

#### Disks

First, we'll set the disks algorithm to deadline in grub:

```bash {linenos=table,hl_lines=[7]}
# If you change this file, run 'update-grub' afterwards to update
# /boot/grub/grub.cfg.

GRUB_DEFAULT=0
GRUB_TIMEOUT=5
GRUB_DISTRIBUTOR=`lsb_release -i -s 2> /dev/null || echo Debian`
GRUB_CMDLINE_LINUX_DEFAULT="quiet elevator=deadline"
GRUB_CMDLINE_LINUX=""

# Uncomment to enable BadRAM filtering, modify to suit your needs
# This works with Linux (no patch required) and with any kernel that obtains
# the memory map information from GRUB (GNU Mach, kernel of FreeBSD ...)
#GRUB_BADRAM="0x01234567,0xfefefefe,0x89abcdef,0xefefefef"

# Uncomment to disable graphical terminal (grub-pc only)
#GRUB_TERMINAL=console

# The resolution used on graphical terminal
# note that you can use only modes which your graphic card supports via VBE
# you can see them in real GRUB with the command `vbeinfo'
#GRUB_GFXMODE=640x480

# Uncomment if you don't want GRUB to pass "root=UUID=xxx" parameter to Linux
#GRUB_DISABLE_LINUX_UUID=true

# Uncomment to disable generation of recovery mode menu entries
#GRUB_DISABLE_LINUX_RECOVERY="true"

# Uncomment to get a beep at grub start
#GRUB_INIT_TUNE="480 440 1"
```

And update grub:

```bash
update-grub
```

Now, to reduces data copies and bus traffic, when you're using LVM partitions, disable the cache and use virtio drivers which are the fastest:

```bash
virt-install ... --file=/dev/vg-name/lv-name,cache=none,if=virtio ...
```

#### Memory

Then, we will enable KSM. Kernel Samepage Merging (KSM) is a feature of the Linux kernel introduced in the 2.6.32 kernel. KSM allows for an application to register with the kernel to have its pages merged with other processes that also register to have their pages merged. For KVM, the KSM mechanism allows for guest virtual machines to share pages with each other. In an environment where many of the guest operating systems are similar, this can result in significant memory savings.

To enable it, add this line:

```bash {linenos=table,hl_lines=[15]}
#!/bin/sh -e
#
# rc.local
#
# This script is executed at the end of each multiuser runlevel.
# Make sure that the script will "exit 0" on success or any other
# value on error.
#
# In order to enable or disable this script just change the execution
# bits.
#
# By default this script does nothing.

# KSM
echo 1 > /sys/kernel/mm/ksm/run
exit 0
```

You can see at anytime the status of KSM by:

```bash
for i in /sys/kernel/mm/ksm/*; do echo -n "$i: "; cat $i; done
```

And in addition, we will disable swapiness to avoid having too much memory consumption. Add those lines in sysctl:

```bash
# Swapiness
vm.swappiness = 1
```

#### Network

For security and performances issues, you should disable ipv6 on bridged interfaces by adding those 3 lines:

```bash
net.bridge.bridge-nf-call-ip6tables = 0
net.bridge.bridge-nf-call-iptables = 0
net.bridge.bridge-nf-call-arptables = 0
```

When you will create network interfaces, uses tap with virtio drivers:

```bash
virt-install ... --network tap,bridge=br0,model=virtio ...
```

#### Virtio

If you want to always enable VirtIO, to get maximum performances, load those modules:

```bash
virtio_blk
virtio_pci
virtio_net
```

### Add user to group

- The installation of kvm created a new system group named kvm in /etc/group. You need to add the user accounts that will run kvm to this group (replace username with the user account name to add):

```bash
adduser username kvm
```

- For those who would like to use libvirt (recommanded), add your user to this group too:

```bash
adduser username libvirt
```

### Storage

There are 2 types of solutions:

- Disk Image: Easier but slower
- Need LVM acknoledges but is faster and simpler to backup (LVM Snapshot)

#### Create Disks image

Create a virtual disk image (10 gigabytes in the example, but it is a sparse file and will only take as much space as is actually used, which is 0 at first, as can be seen with the du command: du vdisk.qcow, while ls -l vdisk.qcow shows the sparse file size):

```bash
qemu-img create -f qcow2 disk0.qcow2 10G
```

You can also use QED, which is faster than qcow2:

```bash
qemu-img create -f qed disk0.qed 10G
```

#### Create LVM LV

To create a Logical Volume, simply:

```bash
lvcreate -L 4G -n vm-name vg-name
```

- -L: size of the logical volume
- -n: Name of the VM
- vg-name: replace by your Volume Group Name

### Network Interfaces

#### Bridged configuration

Now modify your `/etc/network/interfaces` file to add bridged configuration:

```bash
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
EBTABLES_LOAD_ON_START="yes"
EBTABLES_SAVE_ON_STOP="yes"
EBTABLES_SAVE_ON_RESTART="yes"
```

And enable VLAN tagging on bridged interfaces:

```bash
ebtables -t broute -A BROUTING -i eth0 -p 802_1Q -j DROP
```

```bash {linenos=table,hl_lines=["8-15","34-41"]}
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

#### Nat configuration

Nat is the default configuration. But you may need to do some adjustements. Add the forwarding to sysctl:

```bash
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

```xml {linenos=table,hl_lines=[5,6]}
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

## Create a VM

### New method (with libvirt)

- If you want to create a VM with disk image and bridged configuration:

```bash
virt-install --ram=1024 --name=lenny --file=/mnt/vms/lenny/disk0.qcow2 --cdrom=/mnt/isos/debian-6.0.4-amd64-CD-1.iso --hvm --vnc --noautoconsole --accelerate --network=bridge:br0,model=virtio
```

- If you are using LVM and bridged configuration:

```bash
virt-install --os-variant=debiansqueeze --ram=256 --name=vmname --disk path=/dev/mapper/vgname-lvname,device=disk,cache=none,bus=virtio --cdrom=/mnt/isos/debian-6.0.4-amd64-CD-1.iso --hvm --vnc --noautoconsole --accelerate --network=bridge:br0,model=virtio
```

Now a configuration file for this VM has been created in /etc/libvirt/qemu.

- If you are using LVM and nated configuration:

```bash
virt-install --os-variant=debiansqueeze --ram=256 --name=vmname --disk path=/dev/mapper/vg--kvm-vmname,device=disk,cache=none,bus=virtio --cdrom=/mnt/isos/debian-6.0.5-amd64-CD-1.iso --hvm --vnc --noautoconsole --accelerate --network=network:default,model=virtio
```

### Old method (without libvirt)

To make a clean installation of a Guest KVM, you can create a script for each VM you want to create.  
here is a script exemple to launch a KVM using VNC for display instead of X11 display:

```bash
#!/bin/sh
clear
# Var Definition

# Name of your KVM
HOSTNAME="Client 1"

# Path to your virtual hard drive image
HDD="-hda /mnt/vms/lenny/disk0.qcow2"

# Path to your CD-Rom
# note: you can use an iso
CDROM="-cdrom /dev/cdrom"

# Boot Sequence
# note: "c": HDD, "d": CD-Rom, "a": Floppy
BOOT="-boot c"

# TAP Device Creation
# note: don't forget to change "ifname" if you are using mutliples KVM's!
#        This is for "bridged mode". if you wan't to use the "user mode", remove the "TAP" variable
TAP="-net tap,vlan=0,ifname=tap0,script=/etc/kvm/kvm-ifup"

# Virtual Network Card parameters
# note: default model is a "ne2k_pci" (rtl8029) and works on Windows XP and Vista
#        "rtl8139" has better performances and is detected as a 100Mb Adapter
#        "pcnet" or "ai82551" are better for BSD's
NIC="-net nic,model=rtl8139,vlan=0"

# Amount of memory (in Megabyte)
MEM="-m 384"

# Miscelaneous options
# note: "-k fr": if using VNC, it corrects the keyboard problem
#        "-usbdevice" tablet: correct the problem of mouse desynchronisation
#        "-no-acpi": if you are installing a Windows based guest or a BSD
MISC="-k fr -localtime -no-acpi -usbdevice tablet"

# VNC Mode
# note: "-vnc <ip>:<display>": allows clients to connect to specified display only ON (not "from") the specified IP Address
VNC="-vnc 192.168.0.80:1"

# Starting the KVM
# Cheapy Design by Hostin!! ;-p
echo -e "\n\n################################"
echo "##### Starting KVM with... #####"
echo "################################"
echo -e "\n Hard Disk: \"$HDD\" "
echo " CD-Rom: \"$CDROM\" "
echo " Boot Sequence: \"$BOOT\" "
echo " TAP Device: \"$TAP\" "
echo " Virtual Network Card: \"$NIC\" "
echo " Memory size: \"$MEM\" "
echo " Miscelaneous: \"$MISC\" "
echo " VNC Mode: \"$VNC\" "
echo -e "\n################################"
echo -e "\n\n######################################################################"
echo "Kernel-based Virtual Machine: $HOSTNAME - Running"
echo "######################################################################"
echo -e "\n\nLoading kvm-intel kernel module..."
modprobe kvm-intel
exec kvm $HDD $CDROM $BOOT $TAP $NIC $MEM $MISC $VNC
```

## Manage VM

### Manual method

After installation is complete, run it with:

```bash
kvm -hda vdisk.img -m 384
```

Here is a good solution for \*BSD guests:

```bash
/usr/local/bin/qemu-system-x86_64 /data/virt/netbsd.img \
	-net nic,macaddr=00:56:01:02:03:04,model=i82551 \
	-net tap,ifname=tap0,script=/etc/qemu-ifup \
	-m 256 \
	-no-acpi \
	-localtime \
	-daemonize
```

### Virt Manager (GUI)

You can easily install locally or remotly this GUI to manage your VMs:

```bash
apt-get install virt-manager virt-viewer
```

![Virt Manager Screenshot](/images/virt-manager-screenshot.avif)

Just connect remotly or locally and double click to launch a VM (and use it as vnc).

### Virt Manager (Command line)

Use virsh command to manage your VMs. Here is a list of useless examples.

#### Start, stop, list

- List VMs:

```bash
 $ virsh list --all
 Id Name                 State
----------------------------------
  2 Mails        running
  4 Backups      running
  6 Web          running
```

```bash
virsh start vm_name
```

- Shutdown gracefully a VM:

```bash
virsh shutdown vm_name
```

- Force to shutdown the VM:

```bash
virsh destroy vm_name
```

#### Suspend, restore

- Suspend a VM:

```bash
virsh suspend vm_name
```

```bash
virsh resume vm_name
```

#### Delete

- Will delete configuration file (xml in /etc/libvirt/qemu/):

```bash
virsh undefine vm_name
```

#### Backups, restore

- Save a VM:

```bash
virsh save vm_name vm_name.dump
```

- Restore a VM:

```bash
virsh restore vm_name.dump
```

#### Autostart

If you want to set a VM to start automatically on boot:

```bash
virsh autostart <vm_name>
```

## Snapshots

There is a long list on how to manage snapshots ([I suggest that link](https://kashyapc.fedorapeople.org/virt/lc-2012/snapshots-handout.html)). There are 2 kinds of snapshot:

- Internal: snapshot are stored inside the image file
- External: snapshot are stored in an external image file

### Create a snapshot

To create the snapshot:

```bash
virsh snapshot-create-as --domain <vm_name> <snapshot_name> <snapshot_description> --disk-only --diskspec vda,snapshot=external,file=/mnt/vms/<snapshot_name>.qed --atomic
```

- <vm_name>: set the vm domain name
- <snapshot_name>: the name of the snapshot (not as a file, but as it will be displayed on virsh)
- <snapshot_description>: a description of that snapshot
- vda: select the name of the VM device to backup
- <snapshot_name>: set the **full path** of the snapshot file (where it should be stored)

Once that command launched, the base VM disk (not the snapshot) becomes read only and the snapshot is read/write. You could copy the base for the backup if you want.

You can check your disk currently used for write:

```bash
> virsh domblklist pma_qedtest

Target     Source
------------------------------------------------
vda        /mnt/vms/vms/pma_qedtest-snap.qed
```

### External snapshot

With the external snapshot, there are multiple way to create snapshots (blockcommit, blockpull...). They all got their pros and cons.

#### blockpull

This is the simplest and oldest way to merge the base to the snapshot:

```bash
virsh blockpull --domain test --path /mnt/vms/<snapshot_name>.img
```

When finished (look at iotop status), you can remove the base image and keep the snapshot.

#### blockcommit

The blockcommit, is my favorite way to create backups. The actual problem is on Debian 7, this is not present as virsh require a version upper or equal to 0.10.2 and it's only available on Debian unstable for the moment. Anyway, if you've got this version, here is how to do it.

Now we've got something like that:

```
[base(1)]--->[snapshot(2)]
```

If I now want to merge snapshot to base and got only one disk file:

```
[base(1)]<---[snapshot(2)]
[base(2)]
```

you need to follow it:

```bash
virsh blockcommit --domain <vm_name> vda --base /mnt/vms/<base_name>.qed --top /mnt/vms/<snapshot_name>.qed --wait --verbose
```

- --base: the actual read only base disk
- --top: the snapshot disk file to merge in base

## Others

### Reload configuration

If you have made modifications on the xml file and wish reload it:

```bash
virsh define /etc/libvirt/qemu/my_vm.xml
```

### Serial Port Connection

If you want to connect to serial port, [you need to have configured your guest to enable it]({{< ref "docs/Linux/Misc/activating-serial-port-on-linux.md">}}), then connect with virsh:

```bash
virsh connect <hostname>
```

### Add a disk to an existing VM

This is very simple (my VM is called 'ed'):

```bash
virsh attach-disk vmname /dev/mapper/vg-lv sdb
```

- vmname: set the name of your VM
- /dev/mapper/vg-lv: set the disk
- sdb: the device name shown inside the VM

On Debian 6, there is a little bug. So the configuration needs to be reviewed. Replace the current configuration by the running one:

```bash
vmname=ed
virsh dumpxml $vmname > /etc/libvirt/qemu/$vmname.xml
virsh shutdown $vmname
```

Then edit the xml of your VM and change driver name from 'phy' to 'qemu':

```xml {linenos=table,hl_lines=[3]}
[...]
    <disk type='block' device='disk'>
      <driver name='phy' type='raw'/>
[...]
```

And load it:

```bash
virsh define /etc/libvirt/qemu/$vmname.xml
```

Now you can launch your VM

### Bind CPU/Cores to a VM

If you want to bind some CPU/Cores to a VM/Container, there is a solution called CPU Pining :-). First, look at the available cores on your server:

```bash
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
core id 0: processors 0 and 4
core id 1: processors 1 and 5
core id 2: processors 2 and 6
core id 3: processors 3 and 7
```

Now, if I want on a VM a dedicated CPU with it's additional thread, I would prefer do 2 virtual CPU (vpcu) and bind the good core on it. So first, look at the current configuration:

```bash
> virsh vcpuinfo vmname
VCPU:           0
CPU:            6
State:          running
CPU time:       7,5s
CPU Affinity:   yyyyyyyy
```

You can see there is only 1 vcpu. And all the cores of the CPU are used (count the number of 'y' in CPU Affinity, here 8). If we want the best performances, we need to add as many vcpu as we want of cores on a VM, you will see the advantage later... So let's add some cores:

```bash
virsh setvcpus <vmname> <number_of_vcpus>
```

So here for example, we set 4 vcpus. That mean the VM will see 4 cores! Now, we're going to bind processor 0 and 4 on both vcpu! Why? Because if an application doesn't know how to multithread, it will use all the cores! And if applications knows how to use multi cores, they will use it like that. So in any case, you will have good performances :-).

```bash
virsh vcpupin vmname 0 2,6,3,7
virsh vcpupin vmname 1 2,6,3,7
virsh vcpupin vmname 2 2,6,3,7
virsh vcpupin vmname 3 2,6,3,7
```

So now I added 4 virtuals CPU (0 and 1) and added 2 cores (2 and 3) with their associated thread (6 and 7).

In Debian 6 version, it will be done on the fly, but won't be set definitely in the configuration. That's why you'll need to add those parameters (cpuset) in the XML of your VM:

```xml
[...]
  <vcpu cpuset='2,6,3,7'>4</vcpu>
[...]
```

Do not forget to apply the new configuration:

```bash
virsh define /etc/libvirt/qemu/vmname.xml
```

We can now check the new configuration:

```bash
> virsh vcpuinfo vmname
VCPU:           0
CPU:            2
State:          running
CPU time:       8,1s
CPU Affinity:   --yy--yy

VCPU:           1
CPU:            2
State:          running
CPU time:       2,6s
CPU Affinity:   --yy--yy

VCPU:           2
CPU:            3
State:          running
CPU time:       2,7s
CPU Affinity:   --yy--yy

VCPU:           3
CPU:            7
State:          running
CPU time:       6,6s
CPU Affinity:   --yy--yy
```

If you want to know more [how cpusets works, follow that link]({{< ref "docs/Linux/Kernel/process_latency_and_kernel_timing.md#cpuset" >}}).

## Others

### Convert a disk based VM on a LVM parition

You may have a couple of VM based on disk image like qcow2 and my want to convert them into LVM partition. Fortunatly, there is a solution! First convert into your qcow into raw format:

```bash
kvm-img convert disk0.qcow2 -O raw disk0.raw
```

and then put the raw bits into the LVM volume:

```bash
dd if=disk0.raw of=/dev/vg-name/lv-name bs=1M
```

Now edit your xml file and make those changes:

```xml
...
   <disk type='file' device='disk'>
     <source file='/dev/vg0/vm10'/>
     <target dev='hda' bus='ide'/>
   </disk>
...
```

Change to:

```xml
...
   <disk type='block' device='disk'>
     <source dev='vg-name/lv-name'/>
     <target dev='hda' bus='ide'/>
   </disk>
...
```

Now reload your xml file of VM:

```bash
virsh define /etc/libvirt/qemu/vm.xml
```

And now you can start the VM :-)

### Convert an LVM parition to a disk image

To convert an LVM to QED for example, launch that command and adapt it:

```bash
qemu-img convert -O qed /dev/vg_name/lv_name/ /var/lib/libvirt/images/image_name.qed
```

Then edit the VM libvirt configuration file like this:

```xml
    <disk type='file' device='disk'>
      <driver name='qemu' type='qed'/>
      <source file='/var/lib/libvirt/images/image_name.qed'/>
```

Now reload your xml file of VM:

```bash
virsh define /etc/libvirt/qemu/vm.xml
```

And now you can start the VM :-)

### Transfert a LVM disk based VM

If you need to transfer from one server to another a VM based on LVM, there is an easy way solution. You need to first stop the Virtual Machine to have consistency datas, then you can transfer them:

```bash
dd if=/dev/vgname/lvname bs=1M | ssh root@new-server 'dd of=/dev/vgname/lvname bs=1M'
```

Do not forget to transfer xml file configuration of the VM and adapt LVM disks name if needed. Then "virsh define" the new xml file.

### Graphically access to VMs without Virt Manager

If you want to access thought your VMs without installing any manager, you can. First you have to be sure when you created your VM, you entered the **--vnc** option **or** when you launch it, you use this option.

If if it's not hte case and you're using libvirt, please add it to your wished VM:

```xml
...
    <graphics type='vnc' port='-1' autoport='yes' keymap='fr'/>
  </devices>
</domain>
```

Now this is done, you need to change the default listening address of VNC on libvirt. By default, it's listening on 127.0.0.1. This is the most secure choice. However, you may have a secured LAN and wished to open it to anybody. Open so the qemu.conf and modify it to bind on you secure server IP address:

```bash
vnc_listen = "192.168.0.1"
```

If you need as well to activate secure VNC connections, please activate TLS in the same config file.

Then restart or reload libvirt-bin.

### Suspend guests VMs on host shutdown

If your desktop hosts several VMs, it could be interesting to auto suspend them when you restart your computer for example. There is a service for that to make it easy. Simply edit libvirt-guests file configuration:

```bash {linenos=table,hl_lines=[3,10,14,24,28,35]}
# URIs to check for running guests
# example: URIS='default xen:/// vbox+tcp://host/system lxc:///'
URIS=qemu:///system
# action taken on host boot
# - start   all guests which were running on shutdown are started on boot
#           regardless on their autostart settings
# - ignore  libvirt-guests init script won't start any guest on boot, however,
#           guests marked as autostart will still be automatically started by
#           libvirtd
ON_BOOT=ignore

# Number of seconds to wait between each guest start. Set to 0 to allow
# parallel startup.
START_DELAY=0

# action taken on host shutdown
# - suspend   all running guests are suspended using virsh managedsave
# - shutdown  all running guests are asked to shutdown. Please be careful with
#             this settings since there is no way to distinguish between a
#             guest which is stuck or ignores shutdown requests and a guest
#             which just needs a long time to shutdown. When setting
#             ON_SHUTDOWN=shutdown, you must also set SHUTDOWN_TIMEOUT to a
#             value suitable for your guests.
ON_SHUTDOWN=suspend

# If set to non-zero, shutdown will suspend guests concurrently. Number of
# guests on shutdown at any time will not exceed number set in this variable.
PARALLEL_SHUTDOWN=3

# Number of seconds we're willing to wait for a guest to shut down. If parallel
# shutdown is enabled, this timeout applies as a timeout for shutting down all
# guests on a single URI defined in the variable URIS. If this is 0, then there
# is no time out (use with caution, as guests might not respond to a shutdown
# request). The default value is 300 seconds (5 minutes).
SHUTDOWN_TIMEOUT=600

# If non-zero, try to bypass the file system cache when saving and
# restoring guests, even though this may give slower operation for
# some file systems.
#BYPASS_CACHE=0
```

That's all :-)

## FAQ

Read the manual page for more information:

```bash
man kvm
```

### warning: could not open /dev/net/tun: no virtual network emulation

This happen when you want to charge the tun device and you don't have permissions. Simply run your kvm command with sudo.

### Solaris reboot all the time on grub menu

- Run through the installer as usual
- On completion and reboot, the VM will perpetually reboot. "Stop" the VM.
- Start it up again, and immediately open a vnc console and select the Safe Boot from the options screen
- When prompted if you want to try and recover the boot block, say yes
- You should now have a Bourne terminal with your existing filesystem mounted on /a
- Run /a/usr/bin/bash (my preferred shell)
- export TERM=xterm
- vi /a/boot/grub/menu.1st (editing the bootloader on your mounted filesystem), to add "kernel/unix" to the kernel options for the non-safe-mode boot. Ex:

```bash
...
kernel$ /platform/i86pc/multiboot -B $ZFS-BOOTFS kernel/unix
...
```

- Save the file and restart the VM - that's it!

### error: Timed out during operation: cannot acquire state change lock

If you got this kind of error while starting a VM:

```
error: Failed to start domain <vm>
error: Timed out during operation: cannot acquire state change lock
```

it's due to a bug and could be resolved like this:

```bash
/etc/init.d/libvirt-bin stop
rm -Rf /var/run/libvirt
/etc/init.d/libvirt-bin start
```

## Ressources

- https://help.ubuntu.com/community/KVM
- [Documentation for Speeding up QEMU with KVM and KQEMU](/pdf/speed_up_qemu.pdf)
- [Documentation on using KVM on Ubuntu](/pdf/using_kvm_on_ubuntu.pdf)
- [Virtualization With KVM](/pdf/virtualization_with_kvm.pdf)
- [KVM Guest Management With Virt-Manager](/pdf/kvm_guest_management_with_virt-manager.pdf)
- http://www.linux-kvm.org/page/Using_VirtIO_NIC
- http://blog.loftninjas.org/2008/10/22/kvm-virtio-network-performance/
- http://www.linux-kvm.org/page/Tuning_KVM
- http://blog.bodhizazen.net/linux/improve-kvm-performance/
- http://blog.allanglesit.com/2011/05/linux-kvm-vlan-tagging-for-guest-connectivity/
- https://wiki.archlinux.org/index.php/KVM#Enabling_KSM
- http://fr.gentoo-wiki.com/wiki/Libvirt
- http://docs.redhat.com/docs/en-US/Red_Hat_Enterprise_Linux/6/html/Virtualization_Host_Configuration_and_Guest_Installation_Guide/chap-Virtualization_Host_Configuration_and_Guest_Installation_Guide-Network_Configuration.html#sect-Virtualization_Host_Configuration_and_Guest_Installation_Guide-Network_Configuration-Network_address_translation_NAT_with_libvirt
- http://docs.fedoraproject.org/en-US/Fedora/13/html/Virtualization_Guide/ch25s06.html
- http://berrange.com/posts/2010/02/12/controlling-guest-cpu-numa-affinity-in-libvirt-with-qemu-kvm-xen/
- http://wiki.kartbuilding.net/index.php/KVM_Setup_on_Debian_Squeeze
- http://kashyapc.fedorapeople.org/virt/lc-2012/snapshots-handout.html
- http://fedoraproject.org/wiki/Features/Virt_Live_Snapshots
- https://www.redhat.com/archives/libvirt-users/2012-September/msg00063.html
