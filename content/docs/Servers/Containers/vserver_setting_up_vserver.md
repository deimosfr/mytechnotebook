---
weight: 999
url: "/VServer_\\:_Mise_en_place_de_VServer/"
title: "VServer: Setting Up VServer"
description: "A guide on setting up and managing VServer for server virtualization on Linux systems, including creation, management, networking, and troubleshooting."
categories: ["Debian", "Networking", "Virtualization"]
date: "2011-05-24T16:43:00+02:00"
lastmod: "2011-05-24T16:43:00+02:00"
tags:
  [
    "VServer",
    "Virtualization",
    "Networking",
    "Linux",
    "NFS",
    "RAM",
    "DNS",
    "System Administration",
  ]
toc: true
---

## Introduction

Linux-VServer is a security context isolator combined with segmented routing, chroot, extended quotas, and other standard tools.

Initially launched by Jacques Gélinas as the CTX patch, Linux-VServer consists of a patch for the Linux kernel that allows multiple applications to run in different security contexts on the same host machine. Linux-VServer is also equipped with a set of tools to install/manage these contexts.

This project allows one or more operating environments (operating systems without the kernel) to run on a distribution, meaning you can run one or more distributions on a single distribution.

Linux-VServer is a much more advanced virtualization solution than simple chroot.

Not to be confused with the Linux Virtual Server Project.

## Experience Feedback

VServer is an excellent product for applications that don't require many simultaneous network connections. In terms of CPU/RAM, I've never encountered any particular issues. The real peculiarities are on the network side. For example, a Nagios will struggle with more than 800 checks, or a MySQL replication might not function correctly, etc.

VServer remains a major asset but should be used for specific needs without heavy network load.

## Installation

Install the minimum requirements:

```bash
apt-get install linux-image-vserver-686 util-vserver vserver-debiantools ssh
```

Now reboot on the new kernel.

## Creating a VServer

Create the folder where your VServer will be installed and enter it:

```bash
mkdir -p /home/deimos/vserver/vservertest && cd /home/deimos/vserver/
```

Create a VServer like this:

```bash
newvserver --vsroot $(pwd) --hostname vservertest --domain mydomain.local --ip 192.168.0.70/24 --dist etch --mirror http://ftp.fr.debian.org/debian --interface eth0
```

Note: use dummy0 if you have configured a dummy interface!

- vsroot - by default in `/var/lib/vservers`, this is where the vserver will be installed
- Hostname - the name of the vserver
- Domain - the domain of the vserver (like that of the machine)
- IP Address - the IP address of the vserver
- CIDR Range - the CIDR (subnet mask)
- Dist - the distribution to use
- Debian Mirror - the Debian mirror
- Interface - the network interface to use (eth0 by default on most systems)

## VServer Management

To start a VServer and connect to it:

```bash
vserver vservertest start; vserver vservertest enter
```

Here are the basic options:

- Start a vserver: **vserver _vserver_name_ start**
- Stop a vserver: **vserver _vserver_name_ stop**
- Restart a vserver: **vserver _vserver_name_ restart**
- Enter a vserver shell: **vserver _vserver_name_ enter**

### Memory Limits

A vserver kernel keeps track many resources used by each guest (context). Some of these relate to memory usage by the guest. You can place limits on these resources to prevent guests from using all the host memory and making the host unusable.

Two resources are particularly important in this regard:

- The **Resident Set Size** (`rss`) is the amount of pages currently present in RAM.
- The **Address Space** (`as`) is the total amount of memory (pages) mapped in each process in the context.

Both are measured in **pages**, which are 4 kB each on Intel machines (i386). So a value of 200000 means a limit of 800,000 kB, a little less than 800 MB.

Each resource has a **soft** and a **hard limit**.

- If a guest exceeds the `rss` hard limit, the kernel will invoke the Out-of-Memory (OOM) killer to kill some process in the guest.
- The `rss` soft limit is shown inside the guest as the maximum available memory. If a guest exceeds the `rss` soft limit, it will get an extra "bonus" for the OOM killer (proportional to the oversize).
- If a guest exceeds the `as` hard limit, memory allocation attempts will return an error, but no process is killed.
- The `as` soft limit is not utilized until now. In the future it may be used to penalize guests over that limit or it could be used to force swapping on them and such...

Bertl explained the difference between **rss** and **as** with the following example. If two processes share 100 MB of memory, then only 100 MB worth of virtual memory pages can be used at most, so the RSS use of the guest increases by 100 MB. However, two processes are using it, so the AS use increases by 200 MB.

This makes me think that limiting AS is less useful than limiting RSS, since it doesn't directly reflect real, limited resources (RAM and swap) on the host, that deprive other virtual machines of those resources. Bertl says that AS limits can be used to give guests a "gentle" warning that they are running out of memory, but I don't know how much more gentle it is, or how to set it accurately.

For example, 100 processes each mapping a 100 MB file would consume a total of 10 GB of address space (AS), but no more than 100 MB of resources on the host. But if you set the AS limit to 10 GB, then it will not stop one process from allocating 4 GB of RAM, which could kill the host or result in that process being killed by the OOM killer.

You can set the hard limit on a particular context, effective immediately, with this command:

```bash
/usr/sbin/vlimit -c <xid> --<resource> <value>
```

`<xid>` is the context ID of the guest, which you can determine with the `/usr/sbin/vserver-stat` command.

For example, if you want to change the **rss** hard limit for the vserver with `<xid>` 49000, and limit it to 10,000 pages (40 MB), you could use this command:

```bash
/usr/sbin/vlimit -c 49000 --rss 10000
```

You can change the soft limit instead by adding the -S parameter.

Changes made with the vlimit command are effective only until the vserver is stopped. To make permanent changes, write the value to this file:

```bash
/etc/vservers/<name>/rlimits/<resource>.hard
```

To set a soft limit, use the same file name without the `.hard` extension. The `rlimits` directory is not created by default, so you may need to create it yourself.

Changes to these files take effect only when the vserver is started. To make immediate and permanent changes to a running vserver, you need to run vlimit **and** update the rlimits file.

The safest setting, to prevent any guest from interfering with any other, is to set the total of all RSS hard limits (across all running guests) to be less than the total virtual memory (RAM and swap) on the host. It should be sufficiently less to leave room for processes running on the host, and some disk cache, perhaps 100 MB.

However, this is very conservative, since it assumes the worst case where all guests are using the maximum amount of memory at one time. In practice, you can usually get away with contended resources, i.e. allowing guests to use more than this value.

### Network

To configure vservers to have distinct IP addresses, modify the `/etc/network/interfaces` file of the main host:

```bash
# The loopback network interface
auto lo bond0 dummy0
iface lo inet loopback

# Bonding interface
iface bond0 inet static
        address 192.168.0.98
        netmask 255.255.255.0
        gateway 192.168.0.248
        network 192.168.0.0
        up ifenslave bond0 eth0 eth1
        down ifenslave -d bond0 eth0 eth1

# Vserver interface
iface dummy0 inet static
        address 192.168.0.10
        netmask 255.255.255.0
```

Then restart the network:

```bash
/etc/init.d/network restart
```

Now create a vserver and enter **dummy0** as the interface:

```bash
newvserver --vsroot $(pwd) --hostname hostname --domain mydomain.local --ip 192.168.0.12/24 --dist lenny --mirror http://ftp.fr.debian.org/debian --interface dummy0
```

Then add a **name** file containing the last number of the IP you assigned to:

```bash
echo 12 > /etc/vservers/hostname/interfaces/0/name
```

Start your vserver and if you do an ifconfig, you should have something like this:

```bash
dummy0    Lien encap:Ethernet  HWaddr 22:22:xx:xx:CF:xx
          UP BROADCAST RUNNING NOARP  MTU:1500  Metric:1
          RX packets:0 errors:0 dropped:0 overruns:0 frame:0
          TX packets:3 errors:0 dropped:0 overruns:0 carrier:0
          collisions:0 lg file transmission:0
          RX bytes:0 (0.0 b)  TX bytes:210 (210.0 b)

dummy0:12 Lien encap:Ethernet  HWaddr 22:22:xx:xx:CF:xx
          inet adr:192.168.0.12  Bcast:192.168.0.255  Masque:255.255.255.0
          UP BROADCAST RUNNING NOARP  MTU:1500  Metric:1
```

#### Binding

**Important: If you want to install software that runs on a port already in use by another vserver, try to bind to the address of the corresponding server. Ex: with ssh, edit the /etc/ssh/sshd_config file and specify the bind with the IP of the current vserver.**

Here is a table of Applications with configuration files to modify for Binding:

{{< table "table-hover table-striped" >}}
| Application | Config File | Line |
|------------|-------------|------|
| OpenSSH | /etc/ssh/sshd_config | ListenAddress 192.168.0.12 |
| Postgresql | /etc/postgresql/8.1/main/postgresql.conf | listen_addresses = '192.168.0.12' |
| Apache2 | /etc/apache2/ports.conf | Listen 192.168.0.12:80 |
| Munin | /etc/munin/munin-node.conf | host 192.168.0.12 |
| Tomcat5.5 | /etc/tomcat5.5/server.xml | \<Server address="192.168.0.12" port="8005" shutdown="SHUTDOWN"> |
| MySQL | /etc/mysql/my.cnf | bind-address = 192.168.0.12 |
| Postfix | /etc/postfix/main.cf | inet_interfaces = 192.168.0.12 |
| Lighttpd | /etc/lighttpd/lighttpd.conf | server.bind = "192.168.0.12" |
| Nagios NRPE | /etc/nagios/nrpe.cfg | server_address=192.168.0.12 |
{{< /table >}}

#### NFS Mounting

For security reasons, we cannot mount NFS shares by default. To bypass this security, simply create a file called "bcapabilities" in the configuration of the vserver in question:

```bash
echo "CAP_SYS_ADMIN" > /etc/vservers/my_vserver/bcapabilities
```

Then restart your vserver and you will have the ability to mount NFS shares :-)

You may encounter connection problems because you haven't specified the main VServer machine:

```bash
mount: block device backuppc:/mnt/backups/disk_array0/intranet is write-protected, mounting read-only
mount: cannot mount block device backuppc:/mnt/backups/disk_array0/intranet read-only
```

Edit the `/etc/export` file on the server:

```bash
/mnt/backups/disk_array0/intranet     Main_Server_IP(sync,insecure,rw,no_root_squash)   Vserver_IP(sync,insecure,rw,no_root_squash)
```

### VPS

Here's a kind of **ps** but for VServer:

```bash
# vps -ef
root    8102     0 MAIN       8100  0 12:40 pts/1    00:00:00 -bash
root    8210 49159 ns0        5542  0 12:49 ?        00:00:00 sshd: root@pts/2
root    8212 49159 ns0        8210  0 12:49 pts/2    00:00:00 -bash
root    8271     1 ALL_PROC   8102  0 12:57 pts/1    00:00:00 vps -ef
root    8272     1 ALL_PROC   8271  0 12:57 pts/1    00:00:00 ps -ef
```

### Automatically Starting a Vserver at Machine Boot

To automatically start a vserver when the machine boots, simply create a mark file on the desired machine and insert default:

```bash
echo "default" > /etc/vservers/my_vserver/apps/init/mark
```

Here the vserver called "my_vserver" will start automatically when the machine boots.

### Mounting a Folder in Multiple Locations at the Same Time

Modify the file `/etc/vservers/my_server/fstab`. Here's an example:

```bash
/mnt/backups/disk_array0 /mnt/backups/disk_array0 none acl,bind 0 0
```

### Changing the TMP Size

Modify the file `/etc/vservers/my_server/fstab`, by default the `/tmp` is 16MB, simply modify the **Size** parameter and restart the Vserver afterwards:

```bash
none    /tmp            tmpfs   size=2g,mode=1777       0 0
```

## Information on Running VMs

To get information, you need to run a command to get the number of a VM that interests us:

```bash
# vserver-stat
CTX   PROC    VSZ    RSS  userTIME   sysTIME    UPTIME NAME
0       71   1.1G 101.1M   2h33m03   7m39s14   4d01h03 root server
49159    2  14.8M   1.4M   0m00s00   0m00s00  29m09s20 deb-vserv1
49163    2  14.8M   1.4M   0m00s00   0m00s00   0m00s69 deb-vserv2
```

The number in the CTX column corresponds to the VM number.

### RAM

If we take 49159 as an example, we can find out the current RAM:

```bash
# cat /proc/virtual/49159/limit
PROC:            2               4              -1           0
VM:           3808           11410          100000           0
VML:             0               0              -1           0
RSS:           282            1163              -1           0
ANON:           81              81              -1           0
FILES:          37              37              -1           0
OFD:            22              22              -1           0
LOCKS:           1               1              -1           0
SOCK:            1               1              -1           0
MSGQ:            0               0              -1           0
SHM:             0               0              -1           0
```

Here, we have 3808KB used, 11410KB is the maximum memory that has already been used, and the limit is set to 100000KB.

### Load Average

For load average, it's also useful since we can monitor other things. Still on 49159:

```bash
cat cvirt
BiasUptime:     347673.95
SysName:        Linux
NodeName:       deb-vserv1
Release:        2.6.18-4-xen-vserver-amd64
Version:        #1 SMP Wed Apr 18 20:24:37 UTC 2007
Machine:        x86_64
DomainName:     (none)
nr_threads:     2
nr_running:     0
nr_unintr:      0
nr_onhold:      0
load_updates:   418
loadavg:        0.00 0.00 0.00
total_forks:    38
```

Here we have quite a bit of info on the OS as well.

## FAQ

### I'm encountering various network problems

#### Setting network flags (nflags)

- Add the flags to a file named `nflags`:

```bash
echo HIDE_NETIF >> nflags
```

- The default nflags are:

```bash
HIDE_NETIF
```

#### Setting network capabilities (ncaps)

- Add the capabilities to a file named `ncapabilities`:

```bash
echo ^12 >> ncapabilities
```

- There are no default ncaps.

### DNS resolution doesn't work

It's possible that you have a cache Bind server on your host machine that prevents resolution via other DNS servers. In this case, simply authorize the range of your virtual machines to resolve via this bind server.

## Resources

- [Virtualization of servers using Linux vserver](/pdf/virtualisation-de-serveur-grace-a-linux-vserver.pdf)
- [Documentation on Memory Allocation](/pdf/vserver_memory_allocation.pdf)
- [Documentation Multi Vservers](/pdf/multi_vservers.pdf)
- [Xen and vserver: monitoring VMs on a PHP page](./xen_et_vserver_:_monitoring_des_vm_sur_une_page_php)
- [Documentation on Setting up a Web Server with Apache, LVS and Heartbeat 2](/pdf/setup_web_server_cluster.pdf)
- [Capabilities and Flags](https://linux-vserver.org/util-vserver:Capabilities_and_Flags)
- [Going further with Linux vserver](/pdf/un_peu_plus_loin_avec_linux_vserver.pdf)
- [https://linux-vserver.org/Networking_vserver_guests](https://linux-vserver.org/Networking_vserver_guests)
- [https://fr.wikibooks.org/wiki/Linux_VServer](https://fr.wikibooks.org/wiki/Linux_VServer)
