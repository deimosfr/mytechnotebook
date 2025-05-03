---
weight: 999
url: "/Mise_en_place_de_Xen/"
title: "Setting up Xen"
description: "A comprehensive guide for setting up and configuring Xen virtualization on Linux systems, covering installation, network configuration, virtual machine creation and troubleshooting."
categories:
  - "Linux"
  - "Virtualization"
  - "Debian"
date: "2009-12-13T16:24:00+02:00"
lastmod: "2009-12-13T16:24:00+02:00"
tags:
  - "xen"
  - "virtualization"
  - "linux"
  - "hypervisor"
  - "debian"
toc: true
---

## Introduction

Xen allows multiple operating systems (and their applications) to run in isolation on the same physical machine. Guest operating systems share the host machine's resources.

Xen is a "paravirtualizer" or "hypervisor" for virtual machines. Guest operating systems are "aware" of the underlying Xen system and need to be "ported" (adapted) to work on Xen. Linux, NetBSD, FreeBSD, and Plan 9 can already run on Xen.

Xen 3 can also run unmodified systems like Windows on processors that support VT technology.

With Intel Vanderpool and AMD Pacifica technologies, this porting will soon no longer be necessary, and all operating systems will be supported.

The x86, x64, IA-64, PowerPC, and SPARC architectures are supported. Multiprocessor (SMP) and partially Hyper-Threading are supported.

Some might ask, why not use XenExpress or a paid version to get additional functionality? Apart from having support, there's no real reason except for having a nice graphical interface to manage your VMs.

In my opinion, unless you're managing a park of 100 physical machines, the paid version is not really necessary. Here are the differences from the Citrix DataSheet (01/26/2008):

![Xendiff](/images/xendiff.avif)

The version we'll install below has no restrictions and is free :-). On the other hand, you'll spend more time on configuration than with a GUI, that's for sure! It's up to you to see what you really need.

## Installation

### 32 bits

It's very easy to install Xen on Debian:

```bash
apt-get install linux-image-2.6-xen-686 xen-hypervisor-3.0.3-1-i386-pae xen-tools xen-linux-system-2.6.18-5-xen-686 linux-headers-2.6-xen-686 libc6-xen
```

And if you also need to install Windows, then add this:

```bash
apt-get install xen-ioemu-3.0.3-1
```

### 64 bits

It's very easy to install Xen on Debian:

```bash
apt-get install linux-image-2.6-xen-amd64 linux-image-xen-amd64 xen-linux-system-2.6.18-5-xen-amd64 linux-headers-2.6-xen-amd64 xen-hypervisor-3.2-1-amd64 xen-tools xenstore-utils xenwatch xen-shell
```

And if you also need to install Windows, then add this:

```bash
apt-get install xen-ioemu-3.0.3-1
```

## Configuration

### Kernel

To configure the kernel, we'll use certain directives to ensure that dom0_mem never takes more than 512 MB of memory. This is to leave all available space for our domUs:

```bash
title           Xen 3.0.3-1-i386-pae / Debian GNU/Linux, kernel 2.6.18-5-xen-686
root            (hd0,1)
kernel          /boot/xen-3.0.3-1-i386-pae.gz dom0_mem=512m
module          /boot/vmlinuz-2.6.18-5-xen-686 root=/dev/sda2 ro console=tty0 max_loop=64
module          /boot/initrd.img-2.6.18-5-xen-686
```

The "max_loop=64" directive ensures that we won't run out of loopback devices, which are heavily used by Xen. This is a classic error and can be recognized by this type of message:

```bash
Error: Device 2049 (vbd) could not be connected. Backend device not found.
```

### Modules

We'll also load the loop module and set it to 64 as above

```bash
# /etc/modules
loop max_loop=64
```

### Network

#### Configuration of the interface

##### Physical Interface

Let's set up a bridge interface, `/etc/network/interfaces`:

```bash
auto eth0
iface eth0 inet static
        address 192.168.0.90
        netmask 255.255.255.0
        gateway 192.168.0.248

auto eth1
iface eth1 inet manual

auto xenbr0
iface xenbr0 inet static
        address 192.168.130.254
        netmask 255.255.255.0
        bridge_ports eth1
        bridge_maxwait 0
```

We bridge on eth1 because we only have one public IP on this machine. It is therefore out of the question to bridge on the public interface eth0. You'll need to create a dummy0 interface (simulated) if you only have a single physical interface.

**It is also strongly recommended to create a dummy interface to avoid any POSTROUTING problems with Iptables and network slowdowns:**

```bash
iptables -t nat -A PREROUTING -i eth0 -p tcp --dport 80 -j DNAT --to 192.168.0.3:80
```

To see bridges:

```bash
$ brctl show
bridge name	bridge id		STP enabled	interfaces
xenbr0		8000.ee5ad8739af7	no		vif0.0
							peth0
							tap0
							vif6.0
							vif31.0
```

##### Dummy Interface

Modify the `/etc/network/interfaces` file:

```bash
auto eth0
iface eth0 inet static
        address 192.168.0.90
        netmask 255.255.255.0
        gateway 192.168.0.248

iface dummy0 inet manual

auto xenbr0
iface xenbr0 inet static
        address 192.168.130.254
        netmask 255.255.255.0
        bridge_ports dummy0
        bridge_maxwait 0
```

You can check your bridges with this command:

```bash
$ brctl show
xenbr0		8000.1a87802de454	no		dummy0
							vif1.0
							tap0
```

#### Bridge Mode

The "Bridge" mode is set up by the script `/etc/xen/scripts/network-bridge`.
Here's how it works:

- Creation of the new Bridge "_xenbr0_"
- Stopping the "real" network card "_eth0_"
- Copying the MAC and IP addresses from "_eth0_" to a virtual network interface "_veth0_"
- Renaming "_eth0_" to "_peth0_"
- Renaming "_veth0_" to "_eth0_"
- Attaching "_peth0_" and "_vif0.0_" to the bridge "_xenbr0_"
- Starting the bridge interfaces "_xenbr0_", "_peth0_", "_eth0_" and "_vif0.0_"

To enable bridge mode, edit the file **/etc/xen/xend-config.sxp** and uncomment:

```bash
(network-script network-bridge)
(vif-script vif-bridge)
```

Edit this file `/etc/xen-tools/xen-tools.conf` and put this:

```bash
lvm         = my_lvm_vg # Will use an LVM volgroup to create partitions on the fly
debootstrap = 1

size   = 5Gb      # Disk image size.
memory = 128Mb    # Memory size
swap   = 128Mb    # Swap size
fs     = ext3     # use the EXT3 filesystem for the disk image.
dist   = etch     # Default distribution to install.
image  = full

gateway   = 192.168.0.11
netmask   = 255.255.255.0

passwd = 1

kernel = /boot/vmlinuz-`uname -r`
initrd = /boot/initrd.img-`uname -r`

serial_device = hvc0
```

### Linux

To create an image with debootstrap (Debian for example), use the file **/etc/xen-tools/xen-tools.conf**. When you've edited it as desired, you can create the appropriate image. Here are some examples:

```bash
xen-create-image --hostname=vm03.example.com --ip=192.168.0.103 --netmask=255.255.255.0 --gateway=192.168.0.1 --dir=/vserver/images --dist=sarge --debootstrap
```

```bash
xen-create-image --debootstrap --hostname xen-etch --dhcp --dist=etch
```

```bash
xen-create-image --hostname=xen4 --size=3Gb --swap=128Mb --memory=512Mb --ip=172.30.4.155
```

### Windows

For Windows, don't expect amazing performance in terms of network and disk, because until the PV drivers are released from the commercial versions (even Xen Express), performance will remain poor (e.g., network: 1.5 MB/s max).

Nevertheless, if this is sufficient for you, insert this to create a 4 GB file:

```bash
dd if=/dev/zero of=/var/xen/images/WinXP.img bs=1M count=4096
```

Then configure what's needed below to make it work:

```bash
kernel          = "/usr/lib/xen-3.0.3-1/boot/hvmloader"
builder         = 'hvm'
name            = "WindowsXP"
disk            = [ 'file:/mnt/disk1/Winxp/WinXP.img,ioemu:hda,w', 'file:/mnt/disk1/xp-sp2.iso,hdb:cdrom,r' ]

device_model    = '/usr/lib/xen-3.0.3-1/bin/qemu-dm'

# cdrom is no longer used since XEN 3.02
cdrom           = '/dev/hdb'

memory          = 512
boot            = 'dca'
device_model    = 'qemu-dm'
nic=2
vif = [ 'type=ioemu, mac=00:50:56:01:09:01, bridge=xenbr0, model=pcnet' ]
sdl             = 1
vncviewer       = 0
localtime       = 1
ne2000          = 0
vcpus           = 1
serial          = 'pty'

# Correct mouse problems
usbdevice       = 'tablet'

on_poweroff = 'destroy'
on_reboot   = 'restart'
on_crash    = 'restart'
```

Change the vcpu option depending on the number of cores available on your machine.

All that's left is to start the install:

```bash
xm create /etc/xen/WindowsXP.cfg
```

### BSD

```bash
kernel = '/usr/lib/xen-3.0.3-1/boot/hvmloader'
builder = 'hvm'

memory = '512'
name = 'OpenBSD'

device_model    = '/usr/lib/xen-3.0.3-1/bin/qemu-dm'

# Number of network cards
nic = 1
# ne2k_pci is not the fastest (RTL8139), but works with xBSD
# Otherwise pcnet also works
vif = [ 'type=ioemu, bridge=xenbr0, model=pcnet' ]

sdl = 0

# The output will be visible on a vnc server on display 1
vnc = 1
vnclisten = '192.168.130.20'
vncunused = 0
vncdisplay = 1
vncpasswd = ''

# For physical disk: phy:/dev/mapper/example
# For a file: file:/mnt/iso/examples.iso
disk = [ 'file:/root/test.img,ioemu:hda,w', 'file:/root/iso/cd42.iso,hdc:cdrom,r' ]

boot = 'cd'

on_poweroff = 'destroy'
on_reboot = 'restart'
on_crash = 'restart'

acpi = 0
apic = 0
```

For disk images, here are the supported formats:

```bash
vvfat vpc bochs dmg cloop vmdk qcow cow raw parallels qcow2 dmg raw host_device
```

### Additional Options

Here are some additional options that can be very useful...

#### VNC Server at VM Boot

To launch a VNC server when booting a VM, add these lines to your VM config file (here: `/etc/xen/WindowsXP.cfg`):

```bash
vnc             = 1
vncviewer       = 1
vncdisplay      = 1
stdvga          = 0
sdl             = 0
```

Then edit `/etc/xen/xend-config.sxp` and add this:

```bash
# The interface for VNC servers to listen on. Defaults
# to 127.0.0.1  To restore old 'listen everywhere' behaviour
# set this to 0.0.0.0
(vnc-listen '0.0.0.0')
```

You may need a package for this to work:

```bash
apt-get install libsdl1.2debian-all
```

#### Different Boot Devices

- To load a hard disk from an image **for Windows**:

```bash
disk            = [ 'file:/mnt/disk1/Winxp/WindowsXP.img,ioemu:hda,w' ]
```

- Load a CD from a drive **for Windows**:

```bash
disk            = [ 'phy:/dev/cdrom,ioemu:hda,r' ]
```

- Load an ISO image:

```bash
disk            = [ 'file:/mnt/disk1/xp-sp2.iso,hda:cdrom,r' ]
```

hda: **must match your cdrom's udev!**

- Load a CD from a drive:

```bash
disk            = [ 'phy:/dev/cdrom,hda:cdrom,r' ]
```

hda: **must match your cdrom's udev!**

#### Limitations

CPU allocation: more/less:  
Hot increase of used memory:

```bash
xm mem-max xen4 612
xm mem-set xen4 612
```

Don't forget to modify the virtual server configuration file that also contains the memory size.

#### Migration

Memory migration from the source server:

```bash
xm migrate -l xen2migrate xenhost_destination
```

## Launching a Virtual Machine

To start a machine, nothing could be simpler:

```bash
xm create -c /etc/xen/my_xen.cfg
```

**The -c option is used to get control immediately after execution.**  
If you don't use -c, you can get a console like this:

```bash
xm console the_name_of_my_xen_machine
```

If you want to exit the console:

```bash
Ctrl + AltGr + ]
```

Then you can check the status:

```bash
xm list
```

## FAQ

### 4gb seg fixup, process syslogd (pid 15584), cs:ip 73:b7ed76d0

The guest kernel, which is a master vserver, gives a number of insults. The solution is indicated here:  
http://bugs.donarmstrong.com/cgi-bin/bugreport.cgi?bug=405223

It involves installing libc6-xen and doing mv /lib/tls /lib/tls.disabled on the guest vservers.

### Error: Device 2049 (vbd) could not be connected. Hotplug scripts not working

This message is obtained when launching a DomU (xm create toto.cfg). One cause could be the non-existence of one of the DomU partitions. If you use LVM, check that the volumes exist.

### Error: Device 0 (vif) could not be connected. Backend device not found

This is probably because your Xend network is not configured. To fix this, replace the line in your `/etc/xen/xend-config.sxp` file:

```bash
(network-script network-dummy)
```

with

```bash
(network-script 'network-bridge netdev=eth0')
```

Don't forget to restart xend:

```bash
/etc/init.d/xend restart
```

### re0: watchdog timeout

This is all we got when we chose the rtl8139 NIC driver for a NetBSD or OpenBSD domU.

Finally, it's a response from [Manuel Bouyer on port-xen](https://mail-index4.netbsd.org/port-xen/2006/10/25/0001.html) that gives the solution:  
Disable **re\***! It's then **rtk\*** that takes over, and there, no timeout, latency, or other malfunction, "it just works".

## Resources

- http://www.cl.cam.ac.uk/research/srg/netos/xen/readmes/user/
- [Xen Documentation](/pdf/xen.pdf)
- [Paravirtualization with Xen](/pdf/paravirtualization_avec_xen.pdf)
- [Xen Setting up a Perfect Server](/pdf/xen_serveur_parfait.pdf)
- [Xentools Documentation (Xen-Shell and Argo)](/pdf/xen-tools_xen-shell_and_argo.pdf)
- [XenExpress Documentation](/pdf/xenexpress.pdf)
- [Xen for Debian Documentation](/pdf/xen4debian.pdf)
- [How To Make Your Xen-PAE Kernel Work With More Than 4GB RAM](/pdf/how_to_make_your_xen-pae_kernel_work_with_more_than_4gb_ram.pdf)
- [Documentation on Heartbeat2 Xen cluster with drbd8 and OCFS2](/pdf/heartbeat2_xen_cluster_with_drbd8_and_ocfs2.pdf)
- [XEN On An Ubuntu - High Performance](/pdf/xen_on_an_ubuntu_high_performance.pdf)
- [NetBSD Xen Guide](/pdf/netbsd_xen_guide.pdf)
- [XEN and disk space optimization](/pdf/XEN_et_l'optimisation_d'espace_disque.pdf)
- [Consolidation issues and achievement of service level objectives (SLO) with Xen](/pdf/problematique_de_consolidation_slo_sur_xen.pdf)
- [How To Run Fully-Virtualized Guests (HVM) With Xen 3 2 On Debian Lenny](/pdf/how_to_run_fully-virtualized_guests_hvm.pdf)
- [Xen Live Migration Of An LVM-Based Virtual Machine With iSCSI On Debian Lenny](/pdf/xen_live_migration_of_an_lvm-based_virtual_machine_with_iscsi_on_debian_lenny.pdf)
- [Creating A Fully Encrypted Para-Virtualised Xen Guest System Using Debian Lenny](/pdf/creating_a_fully_encrypted_para-virtualised_xen_guest_system_using_debian_lenny.pdf)
