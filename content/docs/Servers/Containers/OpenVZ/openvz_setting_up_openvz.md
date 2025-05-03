---
weight: 999
url: "/OpenVZ_\\:_Mise_en_place_d'OpenVZ/"
title: "OpenVZ: Setting Up OpenVZ"
description: "A comprehensive guide to installing, configuring and managing OpenVZ virtualization technology on Linux."
categories: ["Linux", "Debian", "Backup", "Virtualization"]
date: "2014-04-07T15:52:00+02:00"
lastmod: "2014-04-07T15:52:00+02:00"
tags: ["OpenVZ"]
toc: true
---

## Introduction

[OpenVZ](https://en.wikipedia.org/wiki/Openvz) is an operating system-level virtualization technology based on the Linux kernel. OpenVZ allows a physical server to run multiple isolated instances of operating systems, known as Virtual Private Servers (VPS) or Virtual Environments (VE).

Compared to virtual machines like VMware and paravirtualization technologies like Xen, OpenVZ offers less flexibility in the choice of operating system: both guest and host operating systems must be Linux (although different Linux distributions can be used in different VEs). However, OpenVZ's OS-level virtualization offers better performance, better scalability, higher density, better dynamic resource management, and easier administration than its alternatives. According to the OpenVZ website, this virtualization method introduces a very low performance penalty: only 1 to 3% loss compared to a physical computer.

OpenVZ is the basis of Virtuozzo, a proprietary product provided by SWsoft, Inc. OpenVZ is distributed under the GNU General Public License version 2.

OpenVZ includes the Linux kernel and a set of user commands.

## Installation

Installing OpenVZ on Debian is straightforward:

```bash
aptitude install vzquota vzctl linux-image-2.6-openvz-amd64 linux-image-openvz-amd64 linux-headers-2.6-openvz-amd64 debootstrap
```

We also need to edit the sysctl file to add these settings:

```bash
# On Hardware Node we generally need
# packet forwarding enabled and proxy arp disabled

net.ipv4.conf.default.forwarding=1
net.ipv4.conf.default.proxy_arp = 0
net.ipv4.ip_forward=1

# Enables source route verification
net.ipv4.conf.all.rp_filter = 1

# Enables the magic-sysrq key
kernel.sysrq = 1

# TCP Explict Congestion Notification
#net.ipv4.tcp_ecn = 0

# we do not want all our interfaces to send redirects
net.ipv4.conf.default.send_redirects = 1
net.ipv4.conf.all.send_redirects = 0
```

And activate what we need:

```bash
sysctl -p
```

## Configuration

### Preparation of Environments

You can create your environments in two ways:

- Using debootstrap
- Using templates

Both solutions are valid. The template method is faster as the template is already on your machine.

Personally, I used templates for a long time. More recently, I had to migrate my server to Squeeze (which was in testing at the time) and no templates existed. So I opted for the debootstrap version.

### Debootstrap

Creating a VE with debootstrap is very simple. Here's the list of available distributions:

```bash
> ls /usr/share/debootstrap/scripts/
breezy	edgy  etch-m68k  gutsy	hoary	      intrepid	karmic	lucid	  potato  sarge.buildd	    sid      stable   unstable	warty.buildd  woody.buildd
dapper	etch  feisty	 hardy	hoary.buildd  jaunty	lenny	maverick  sarge   sarge.fakechroot  squeeze  testing  warty	woody
```

For my part, I chose Squeeze:

```bash
debootstrap --arch amd64 squeeze /mnt/containers/private/101
```

- --arch: Choose the desired architecture, here amd64.
- squeeze: The name of the desired distribution
- /mnt/containers/private/101: The location where the VE should be placed (I'm using the VE ID here)

### Template

Instead of creating our template, we can download pre-existing ones from: http://wiki.openvz.org/Download/template/precreated and place them in the right location. I modified the file `/etc/vz/vz.conf` to change the default paths. It now points to `/mnt/containers`.

```bash
cd /mnt/containers/template/cache
wget http://download.openvz.org/template/precreated/contrib/debian-5.0-amd64-minimal.tar.gz
```

## Networking

### NAT Mode

If you need to have multiple VEs available from a single IP and you can't have other directly accessible IPs, you'll need to use NAT mode. This is especially the case if you have a Dedibox (Illiad) or a Kimsufi server (OVH) where only one WAN IP is available (by default) for a given server and it's impossible to use [bridge mode](#with-a-bridged-interface).

Let's start by configuring the network card:

```bash
# This file describes the network interfaces available on your system
# and how to activate them. For more information, see interfaces(5).

# The loopback network interface
auto lo
iface lo inet loopback

# The primary network interface
auto eth0
iface eth0 inet static
        address 88.xx.xx.xx
        netmask 255.255.255.0
        network 88.xx.xx.0
        broadcast 88.xx.xx.255
        gateway 88.xx.xx.xx
```

So far, everything is standard. However, after creating a VE, we'll need masquerading to provide internet access from the VEs and PREROUTING for port redirection. We'll use IPtables to provide external access from VEs and port redirections from the outside to the VEs:

```bash {linenos=table}
#!/bin/bash

#-------------------------------------------------------------------------
# Essentials
#-------------------------------------------------------------------------

IPTABLES='/sbin/iptables';
modprobe nf_conntrack_ftp

#-------------------------------------------------------------------------
# Physical and virtual interfaces definitions
#-------------------------------------------------------------------------

# Interfaces
wan_if="eth0";
vpn_if="tap0";

#-------------------------------------------------------------------------
# Networks definitions
#-------------------------------------------------------------------------

# Networks
wan_ip="x.x.x.x";
lan_net="192.168.90.0/24";
vpn_net="192.168.20.0/24";

# IPs
ed_ip="192.168.90.1";
banzai_ip="192.168.90.2";

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

$IPTABLES -A INPUT -j ACCEPT -d $lan_net;
$IPTABLES -A INPUT -j ACCEPT -m state --state ESTABLISHED,RELATED

#-------------------------------------------------------------------------
# Allow masquerading for VE
#-------------------------------------------------------------------------

# Activating masquerade to get Internet from VE
$IPTABLES -t nat -A POSTROUTING -o $wan_if -s $lan_net -j MASQUERADE

# Activating masquerade to get VPN access from VE
$IPTABLES -t nat -A POSTROUTING -o tap0 -j MASQUERADE

#-------------------------------------------------------------------------
# Allow ports on CT
#-------------------------------------------------------------------------

# Allow ICMP
$IPTABLES -A INPUT -j ACCEPT -p icmp

# SSH access
$IPTABLES -A INPUT -j ACCEPT -p tcp --dport 22

#-------------------------------------------------------------------------
# Redirections for incoming connections (wan)
#-------------------------------------------------------------------------

# HTTP access
$IPTABLES -t nat -A PREROUTING -p tcp --dport 80 -d $wan_ip -j DNAT --to-destination $ed_ip:80

# HTTPS access
$IPTABLES -t nat -A PREROUTING -p tcp --dport 443 -d $wan_ip -j DNAT --to-destination $ed_ip:443
```

Here I have ports 80 and 443 redirected to one machine (Ed).

### With a Bridged Interface

The major advantage of the venet interface is that it allows the administrator of the physical server to decide on the network configuration of the virtual machine. This is particularly appreciated for a hosting provider who wants to sell a virtual machine hosting service, as they can freely let their customers manage their OpenVZ virtual server while being certain that they cannot disrupt the network, allocate more IP addresses than planned, modify routing tables, etc. If you're using OpenVZ to provide hosting services to your customers, the veth interface is not for you. Since I'm not in this situation, the potential security issues posed by the veth interface for a hosting provider don't concern me.

I use both venet and veth on the servers I manage: veth when I need to have access to ethernet layers from virtual servers, and venet in all other cases. But, unlike the simplicity offered by venet interfaces, to be able to use a veth interface in a virtual server, you need to perform some additional configuration operations to prepare the physical server to receive virtual servers using veth. That's what I'll describe in the following sections.

First, let's install the necessary tools:

```bash
apt-get install bridge-utils
```

Then create a typical configuration:

```bash
auto lo
iface lo inet loopback

auto eth0
iface eth0 inet manual

auto vmbr0
iface vmbr0 inet static
    bridge_ports eth0
    address 192.168.1.216
    netmask 255.255.255.0
    network 192.168.1.0
    broadcast 192.168.1.255
    gateway 192.168.1.254
```

Now reboot the server.

For OpenVZ to dynamically add or remove a veth interface to the bridge when starting a virtual server, we need to create some files on the physical server. This operation needs to be done once for all on the physical server, regardless of the number of virtual servers that will be created later. It's also possible to do it on a server that, in the end, hosts only virtual servers using venet interfaces - it doesn't disturb them at all.

First, create the `/etc/vz/vznet.conf` file with the following content:

```bash
#!/bin/bash
EXTERNAL_SCRIPT="/usr/sbin/vznetaddbr"
```

Then create the file `/usr/sbin/vznetaddbr` with the following content:

```bash
#!/bin/sh
#
# Add virtual network interfaces (veth's) in a container to a bridge on CT0

CONFIGFILE=/etc/vz/conf/$VEID.conf
. $CONFIGFILE

NETIFLIST=$(printf %s "$NETIF" |tr ';' '\n')

if [ -z "$NETIFLIST" ]; then
   echo >&2 "According to $CONFIGFILE, CT$VEID has no veth interface configured."
   exit 1
fi

for iface in $NETIFLIST; do
    bridge=
    host_ifname=

    for str in $(printf %s "$iface" |tr ',' '\n'); do
	case "$str" in
	    bridge=*|host_ifname=*)
		eval "${str%%=*}=\${str#*=}" ;;
	esac
    done

    [ "$host_ifname" = "$3" ] ||
	continue

    [ -n "$bridge" ] ||
	bridge=vmbr0

    echo "Adding interface $host_ifname to bridge $bridge on CT0 for CT$VEID"
    ip link set dev "$host_ifname" up
    echo 1 >"/proc/sys/net/ipv4/conf/$host_ifname/proxy_arp"
    echo 1 >"/proc/sys/net/ipv4/conf/$host_ifname/forwarding"
    brctl addif "$bridge" "$host_ifname"

    break
done

exit 0
```

For information, I modified the line containing "bridge=" with my interface vmbr0.

You also need to set the execute permission on this file by typing:

```bash
chmod 0500 /usr/sbin/vznetaddbr
```

Now let's configure the VE (see also how to manage a VE below):

```bash
vzctl set $my_veid --netif_add eth0 --save
```

Finally edit your machine's configuration and add the bridge interface for it:

```bash
CONFIG_CUSTOMIZED="yes"
VZHOSTBR="vmbr0"
```

The end of the configuration file should look like this:

```bash
...
OSTEMPLATE="debian-5.0-amd64-minimal"
ORIGIN_SAMPLE="vps.basic"
HOSTNAME="vz.deimos.fr"
CONFIG_CUSTOMIZED="yes"
VZHOSTBR="vmbr0"
IP_ADDRESS=""
NAMESERVER="192.168.100.3"
CAPABILITY="SYS_TIME:on "
NETIF="ifname=eth0,mac=00:18:51:96:D4:8D,host_ifname=veth101.0,host_mac=00:18:51:B8:B8:CF"
```

Now you just need to configure the eth0 interface as you normally would in your VE.

### With VLANs

You may need to create VLANs in your VEs. This works very well with a bridged interface. To do this, on the host machine, you must have a configured VLAN ([use this documentation for setup]({{< ref "docs/Servers/Network/setting_up_vlan.md" >}})). For those who still want an example:

```bash
# This file describes the network interfaces available on your system
# and how to activate them. For more information, see interfaces(5).

# The loopback network interface
auto lo
iface lo inet loopback

# The primary network interface
allow-hotplug eth0
auto eth0
iface eth0 inet manual

# The bridged interface
auto vmbr0
iface vmbr0 inet static
        address 192.168.100.1
        netmask 255.255.255.0
        gateway 192.168.100.254
        broadcast 192.168.100.255
        network 192.168.100.0
        bridge_ports eth0
        bridge_fd 9
        bridge_hello 2
        bridge_maxage 12
        bridge_stp off

# The DMZ Vlan 110
auto vmbr0.110
iface vmbr0.110 inet static
	address 192.168.110.1
	netmask 255.255.255.0
	broadcast 192.168.110.255
	vlan_raw_device vmbr0
```

This example is made with a bridged interface because I have [KVM]({{< ref "docs/Servers/Virtualization/KVMandQemu/kvm_setting_up_kvm.md" >}}) running on it, but there's nothing forcing it to be bridged.

Then, when you create your VE, you don't have to do anything special when creating the network interface for your VE. Launch the creation of your VE and don't forget to install the "vlan" package to be able to create VLAN access within your VE. Here's another example to give you an idea of the VE network configuration:

```bash
...
CONFIG_CUSTOMIZED="yes"
VZHOSTBR="vmbr0"
IP_ADDRESS=""
NETIF="ifname=eth0,mac=00:18:50:FE:EF:0B,host_ifname=veth101.0,host_mac=00:18:50:07:B8:F4"
```

For the VE configuration, it's almost identical to the host machine configuration - you need to create a VLAN interface on the main interface (again, there's no need to have the main interface configured, just the VLAN is enough). For those who are still skeptical, here's an example of configuration in a VE:

```bash
# This configuration file is auto-generated.
# WARNING: Do not edit this file, your changes will be lost.
# Please create/edit /etc/network/interfaces.head and /etc/network/interfaces.tail instead,
# their contents will be inserted at the beginning and at the end
# of this file, respectively.
#
# NOTE: it is NOT guaranteed that the contents of /etc/network/interfaces.tail
# will be at the very end of this file.

# Auto generated lo interface
auto lo
iface lo inet loopback

# VE interface
auto eth0
iface eth0 inet manual

# VLAN 110 interface
auto eth0.110
iface eth0.110 inet static
	address 192.168.110.2
	netmask 255.255.255.0
	gateway 192.168.110.254
	broadcast 192.168.110.255
	vlan_raw_device eth0
```

### With Bonding

You may need to create bridged bonding in your VEs. To do this, on the host machine, you must have a configured bonding ([use this documentation for setup]({{< ref "docs/Linux/Network/network_creating_bonding.md" >}})). For those who still want an example:

```bash
# This file describes the network interfaces available on your system
# and how to activate them. For more information, see interfaces(5).

# The loopback network interface
auto lo eth0 eth1
iface lo inet loopback

iface eth0 inet manual
iface eth1 inet manual

auto bond0
iface bond0 inet manual
    slaves eth0 eth1
    bond_mode active-backup
    bond_miimon 100
    bond_downdelay 200
    bond_updelay 200

auto vmbr0
iface vmbr0 inet static
    address 192.168.0.227
    netmask 255.255.255.0
    network 192.168.0.0
    gateway 192.168.0.245
    bridge_ports bond0
```

## NFS

Here I will only cover the client side in a VE, not a server in a VE.

### Server

First on the server side, [set up your NFS server by following this documentation]({{< ref "docs/Servers/FileSharing/nfs_setting_up_an_nfs_server.md" >}}).

Then we'll install this, which will allow us to use not only the NFS v3 protocol, but also to abstract the kernel layer (useful in case of a crash):

```bash
aptitude install nfs-user-server
```

### Client

- On the host machine (HN), add these lines to sysctl.conf:

```bash
...
# OpenVZ NFS Client
sunrpc.ve_allow_rpc = 1
fs.nfs.ve_allow_nfs = 1
kernel.ve_allow_kthreads = 1
```

Then apply all this on the fly:

```bash
sysctl -p
```

Now let's activate NFS on the VEs that interest us:

```bash
vzctl set $my_veid --features "nfs:on" --save
```

I don't know if this step is essential, but just in case, I'll add it anyway:

```bash
aptitude install nfs-user-server
```

Then you can mount your NFS mount points normally. However, you may encounter some permission issues. That's why I recommend adding the 'no_root_squash' option in the exports file on the server:

```bash
/mnt/backups/backups 192.168.0.127(rw,no_root_squash)
```

And on the client, add the nolock option to mount the NFS:

```bash
mount -t nfs -o nolock @IP:/my/share my_mount_point
```

## Mount Bind

Mount binds can sometimes be very useful. Here's how to do them:

### Mount script

Create a script in `/etc/vz/conf/vps.mount` for all VEs or `/etc/vz/conf/CTID.mount` for a specific VE (replace CTID with the VE number):

```bash
#!/bin/bash
source /etc/vz/vz.conf
source ${VE_CONFFILE}
mount --bind /mnt/disk ${VE_ROOT}/mnt/disk
```

And adapt the last line to your needs.

### Unmount script

And finally the same for the unmounting script, so `/etc/vz/conf/vps.umount` or `/etc/vz/conf/CTID.umount`:

```bash
#!/bin/bash
source /etc/vz/vz.conf
source ${VE_CONFFILE}
umount ${VE_ROOT}/mnt/disk
exit 0
```

Finally apply the necessary permissions to make it executable:

```bash
chmod u+x /etc/vz/conf/CTID.mount /etc/vz/conf/CTID.umount
```

## Management

### Create a VE

Choose the method you want depending on the VM creation method you prefer (template or debootstrap).

#### Template

To create a container, use this command:

```bash
my_veid=101
vzctl create $my_veid --ostemplate debian-5.0-amd64-minimal --config vps.basic
vzctl set $my_veid --onboot yes --save
vzctl set $my_veid --hostname nagios.mycompany.com --save
vzctl set $my_veid --ipadd 192.168.0.130 --save
vzctl set $my_veid --nameserver 192.168.0.27 --save
```

If you want to configure your interfaces in bridge mode, don't forget [this part](#with-a-bridged-interface).

#### Debootstrap

At the time of writing, Squeeze is not yet the stable version, but has just frozen to become stable (today itself). So there are a few small things that may vary, like this that we need to create:

```bash
cp /etc/vz/conf/ve-basic.conf-sample /etc/vz/conf/ve-vps.basic.conf-sample
```

We're copying the default parameters of a VE to a new name that it will be able to take by default. Let's configure it:

```bash
> my_veid=101
> vzctl set $my_veid --applyconfig vps.basic --save
WARNING: /etc/vz/conf/101.conf not found: No such file or directory
Saved parameters for CT 101
```

Then we add some additional lines needed in our configuration file:

```bash
echo "OSTEMPLATE=debian" >> /etc/vz/conf/$my_veid.conf
```

Then we configure the network parameters:

```bash
vzctl set $my_veid --ipadd 192.168.0.130 --save
vzctl set $my_veid --nameserver 192.168.0.27 --save
```

Next we'll remove udev, which can prevent the VM from booting:

```bash
rm /mnt/containers/private/101/etc/rcS.d/S02udev
```

Then we'll start the VE and enter it:

```bash
vzctl start $my_veid
vzctl enter $my_veid
```

Configure your source.list file:

```bash
deb http://ftp.fr.debian.org/debian/ squeeze main non-free contrib
deb-src http://ftp.fr.debian.org/debian/ squeeze main non-free contrib

deb http://security.debian.org/ squeeze/updates main contrib non-free
deb-src http://security.debian.org/ squeeze/updates main contrib non-free
```

Then run these commands to remove the superfluous and remove the gettys (VEs don't use them):

```bash
sed -i -e '/getty/d' /etc/inittab
rm -Rf /lib/udev/
```

We solve the mtab problem:

```bash
rm -f /etc/mtab
ln -s /proc/mounts /etc/mtab
```

### Start a VE

Now you can start your VE:

```bash
vzctl start 101
```

### Change the root password of a VE

If you want to change the root password:

```bash
vzctl exec 101 passwd
```

### List VEs

To list your VEs, use the vzlist command:

```bash
$ vzlist
     VEID      NPROC STATUS  IP_ADDR         HOSTNAME
      101          8 running 192.168.0.130   nagios.mycompany.com
```

### Stop a VE

To stop your VE:

```bash
vzctl stop 101
Stopping container ...
Container was stopped
Container is unmounted
```

### Restart a VE

To restart a VE:

```bash
$ vzctl restart 101
Restarting VE
Stopping VE ...
VE was stopped
VE is unmounted
Starting VE ...
VE is mounted
Adding IP address(es): 192.168.0.130
Setting CPU units: 1000
Configure meminfo: 500000
Set hostname: nagios.mycompany.com
File resolv.conf was modified
VE start in progress...
```

### Destroy a VE

To permanently delete a VE:

```bash
$ vzctl destroy 101
Destroying container private area: /vz/private/101
Container private area was destroyed
```

### Enter a VE

```bash
$ vzctl enter 101
entered into VE 101
```

### VE Limits

By default there are limits that VEs cannot exceed. We'll see how to manage them here. If you still want to go further, know that everything is explained on the official website http://wiki.openvz.org/User_beancounters.

#### List the limits

To list the limits, we'll execute this command:

```bash
$ vzctl exec 101 cat /proc/user_beancounters
Version: 2.5
      uid  resource                     held              maxheld              barrier                limit              failcnt
     101:  kmemsize                  1301153              3963508             14372700             14790164                    0
           lockedpages                     0                    8                  256                  256                    0
           privvmpages                  8987                49706                65536                69632                    0
           shmpages                      640                  656                21504                21504                    0
           dummy                           0                    0                    0                    0                    0
           numproc                         9                   21                  240                  240                    0
           physpages                    1786                27164                    0  9223372036854775807                    0
           vmguarpages                     0                    0                33792  9223372036854775807                    0
           oomguarpages                 1786                27164                26112  9223372036854775807                    0
           numtcpsock                      3                    6                  360                  360                    0
           numflock                        1                    7                  188                  206                    0
           numpty                          1                    2                   16                   16                    0
           numsiginfo                      0                    2                  256                  256                    0
           tcpsndbuf                  106952               106952              1720320              2703360                    0
           tcprcvbuf                   49152              1512000              1720320              2703360                    0
           othersockbuf                    0                21008              1126080              2097152                    0
           dgramrcvbuf                     0                 5648               262144               262144                    0
           numothersock                   25                   29                  360                  360                    0
           dcachesize                 115497               191165              3409920              3624960                    0
           numfile                       234                  500                 9312                 9312                    0
           dummy                           0                    0                    0                    0                    0
           dummy                           0                    0                    0                    0                    0
           dummy                           0                    0                    0                    0                    0
           numiptent                      10                   10                  128                  128                    0
```

The last line should always be 0, otherwise you've clearly reached the limits.

Here's some important information:

- held: Currently the value of the resource
- maxheld: the maximum value the resource has reached
- barrier: the soft limit (warning). It means that the held value has already been up to this value
- limit: the hard limit. The held value will never exceed this value.

To increase them, you have two solutions:

- You can do it [on the fly](#apply-limits-on-the-fly).
- Go to `/etc/vz/conf/101.conf` and increase the values that are causing problems. Unfortunately, this requires a restart of the VE for these parameters to be taken into account.

It's possible that the failcnt isn't reset to 0... don't worry, it will happen later.

#### Maximize the limits

If you have too many restrictions compared to your needs or you feel that the VE you're going to create will be heavy, maximize this from the start:

```bash
vzsplit -n 2 -f max-limits
```

- 2: change this number to the number of VMs you want to run on your machine and it will take care of calculating the limits for your VE in the best way

A file **/etc/vz/conf/ve-max-limits.conf-sample** is created. You can edit it at any time and make the modifications you want. If you want to apply it to a VE:

```bash
vzctl set 101 --applyconfig max-limits --save
```

VE 101 now has a new configuration via the 'max-limits' config.

#### Apply limits on the fly

For example, if you want to change the RAM size on the fly:

```bash
vzctl set 101 --privvmpages 786432:1048576 --save --setmod restart
```

- 786432: corresponding to the soft limit (barrier)
- 1048576: corresponding to the hard limit (limit)

Thanks to 'setmod restart', I can now apply limits on the fly.

#### Disk quotas

You may have set up disk quotas and you're reaching 100% of your VE disk. No problem, we'll increase the quota on the fly:

```bash
vzctl set 101 --diskspace 14G:15G --save
```

If you don't want a quota, add this line to the VE configuration:

```bash
DISK_QUOTA=no
```

#### Limit the CPU

If you want to limit the CPU, there are several ways to do it. Either you play with the scheduler, or you apply a percentage. For basic uses, the second solution will suit you:

```bash
> vzctl set $my_veid --cpulimit 80 --save --setmod restart
Setting CPU limit: 80
Saved parameters for CT 102
```

For example, here I limited my VE 102 to 80% of CPU. For more info: http://wiki.openvz.org/Resource_shortage

#### Disable the most restrictive limits

If you have a lot of load on a VE, you'll regularly need to increase certain parameters, which can quickly become tiring. To be at ease with some of them, you can simply disable the limits by entering the maximum value:

```bash
vzctl set $my_veid --kmemsize unlimited --save --setmod restart
vzctl set $my_veid --lockedpages unlimited --save --setmod restart
vzctl set $my_veid --privvmpages unlimited --save --setmod restart
vzctl set $my_veid --shmpages unlimited --save --setmod restart
vzctl set $my_veid --numproc unlimited --save --setmod restart
vzctl set $my_veid --numtcpsock unlimited --save --setmod restart
vzctl set $my_veid --numflock unlimited --save --setmod restart
vzctl set $my_veid --numpty unlimited --save --setmod restart
vzctl set $my_veid --numsiginfo unlimited --save --setmod restart
vzctl set $my_veid --tcpsndbuf unlimited --save --setmod restart
vzctl set $my_veid --tcprcvbuf unlimited --save --setmod restart
vzctl set $my_veid --othersockbuf unlimited --save --setmod restart
vzctl set $my_veid --dgramrcvbuf unlimited --save --setmod restart
vzctl set $my_veid --numothersock unlimited --save --setmod restart
vzctl set $my_veid --dcachesize unlimited --save --setmod restart
vzctl set $my_veid --numfile unlimited --save --setmod restart
vzctl set $my_veid --numiptent unlimited --save --setmod restart
```

#### Validation of limits

If you want to validate your limits to make sure everything is ok:

```bash
> vzcfgvalidate /etc/vz/conf/101.conf
Validation completed: success
```

#### Automatic validation of limits

If, for example, you want to validate your machine's resource usage in the best way and also follow OpenVZ best practices, you can use the vzcfgvalidate tool:

```bash
vzcfgvalidate -i /etc/vz/conf/101.conf
```

This will launch an interactive mode informing you of what's wrong. It's up to you to validate whether or not you want to make the changes. If you don't have the necessary analysis capabilities for the problems or simply don't want to bother with this, replace the '-i' option with '-r', which will automatically adjust to the correct values.

For those who don't want to bother at all with limits, here's a small script that will automatically adjust all limits of all your VEs. Just add it to crontab:

```bash
#!/bin/sh
# Auto Ajust VE
# Made by Pierre Mavro / Deimos
# Contact: xxx@mycompany.com

# Configure VE directory config path
VE_CONF="/etc/vz/conf"

echo "Auto ajustment start:"
for i in `ls $VE_CONF/*.conf`; do
	if [ $i != "$VE_CONF/0.conf" ]; then
                echo "Ajusting $i..."
		vzcfgvalidate -r $i
	fi
done
echo "Auto ajustment done"
```

#### Automatic adjustment of limits

Whether for my work or for home, the constant modification of resource values is a bit of a headache for me. So I created a tool to manage all this on its own called [vzautobean (OpenVZ Auto Management User Beancounters)](https://www.deimos.fr/gitweb/?p=vzautobean.git;a=summary).

To install it is very simple, we'll fetch it and place it in the folder with the other OpenVZ binaries:

```bash
git clone git://git.deimos.fr/git/vzautobean.git
sudo mv vzautobean/vzautobean /usr/sbin/
```

Then either you launch it and it will run on all your VEs in use, or you can specify one or more specific VEs:

```bash
vzautobean -e 102 -e 103
```

Here it will do VE 102 and 103. By default it increases all barriers by 10% and limits by 20% relative to maxheld.

If you wish, you can change these values by adding 2 arguments:

```bash
vzautobean -b 30 -l 50
```

Here the barrier should be set to 30% above maxheld and the limit to 50% above. And this for all VEs that are launched!

Now all you have to do is put this in cron if you want this program to run regularly.

### Backup and restore

#### Installation

You may need to backup an entire environment, and for that there's the vzdump command:

```bash
aptitude install vzdump
```

If you're on a Debian older than 6, you can find the command here:

```bash
wget http://download.proxmox.com/debian/dists/lenny/pve/binary-amd64/vzdump_1.2-11_all.deb
dpkg -i vzdump_1.2-11_all.deb
```

#### Backup

To backup, it's very simple:

```bash
vzdump $my_veid
```

#### Restore

If you want to restore:

```bash
vzrestore backup_101.tar $my_veid
```

The machine will be restored on the VE corresponding to $my_veid.

## Services

Some services may have some problems working properly (always the same ones). That's why I'll give solutions for the ones I've encountered.

### NTP

For the NTP service, simply configure the VE like this:

```bash
vzctl set 101 --capability sys_time:on --save
```

### GlusterFS

If you want to use glusterfs in a VE, you may encounter permission problems:

```
fuse: failed to open /dev/fuse: Permission denied
```

To work around this, we'll create the fuse device from the host on the VE in question and add admin rights to it (a bit poor in terms of security, but no choice):

```bash
vzctl set $my_veid --devices c:10:229:rw --save
vzctl exec $my_veid mknod /dev/fuse c 10 229
vzctl set $my_veid --capability sys_admin:on --save
```

### Encfs

If you want to use encfs in a VE, you may encounter permission problems:

```
EncFS Password:
fuse: device not found, try 'modprobe fuse' first
fuse failed.  Common problems:
 - fuse kernel module not installed (modprobe fuse)
 - invalid options -- see usage message
```

Be aware that you need to load the fuse module at the VZ level for VEs to inherit it. Add this to your VZ to avoid having to load the module at each boot:

```bash
...
# Load Fuse
fuse
...
```

Then load it dynamically to access it afterwards:

```bash
modprobe fuse
```

To work around these issues, we'll create the fuse device from the host on the VE in question and add admin rights to it (a bit poor in terms of security, but no choice):

```bash
vzctl set $my_veid --devices c:10:229:rw --save
vzctl exec $my_veid mknod /dev/fuse c 10 229
vzctl set $my_veid --capability sys_admin:on --save
```

It's possible that the second line doesn't work when the VE is off. Run it once it's turned on, then mount your encfs partition.

### OpenVPN

If you want to run an OpenVPN server in a VE, add these kinds of permissions and create the necessary devices:

```bash
vzctl set $my_veid --devices c:10:200:rw --save
vzctl set $my_veid --capability net_admin:on --save
vzctl exec $my_veid mkdir -p /dev/net
vzctl exec $my_veid mknod /dev/net/tun c 10 200
vzctl exec $my_veid chmod 600 /dev/net/tun
vzctl set $my_veid --devnodes net/tun:rw --save
```

## FAQ

### bridge vmbr0 does not exist!

Damn!!! Oh I had a hard time with this sh\*t! Indeed by default Debian assigns the name of the interface to vmbr0 for bridging (certainly a new nomenclature for Debian 6). I had found this problem when installing my server and then like an idiot I hadn't noted it down! When you do an update, it may stop working (the bridge of your VMs) because the file `/usr/sbin/vznetaddbr` may be replaced. When launching a VE, this gives you something like:

```
Adding interface veth101.0 to bridge vmbr0 on CT0 for CT101
bridge vmbr0 does not exist!
```

To fix the problem quickly, modify this file and replace the name of this interface with the correct one (here br0). For my part, I changed the name of the interface to avoid being bothered in the future.

### I have freezes when booting/stopping my VZs

If you have freezes at boot or stop of VZs and you're in bridge mode, you need to change the MAC addresses like this:

```bash
sed -i -e "s/00:18/FE:FF/g" /etc/vz/conf/*.conf
sed -i -e "s/00:AB:BA/FE:FF:BA/g" /etc/vz/conf/*.conf
```

## Resources

- [Some Tips On OpenVZ Deployment](/pdf/some_tips_on_openvz_deployment.pdf)
- [Splitting Resources Evenly Between OpenVZ VMs With vzsplit](/pdf/splitting_resources_evenly_between_openvz_vms_with_vzsplit.pdf)
- [Installing And Using OpenVZ On Debian Lenny AMD64](/pdf/installing_and_using_openvz_on_debian_lenny_amd64.pdf)
- http://www.libresys.fr/2008/10/14/les-differentes-formes-de-configuration-du-reseau-avec-openvz/
- http://www.famille-fontes.net/comments.php?y=10&m=02&entry=entry100208-094400
- http://wiki.openvz.org/Backup_of_a_running_container_with_vzdump
- http://www.kafe-in.net/Blog/OpenVZ-Les-ressources-syst-me-dans-un-VE
