---
weight: 999
url: "/Installation_et_Configuration_de_Red_Hat_Cluster_Suite/"
title: "Installation and Configuration of Red Hat Cluster Suite"
description: "How to install and configure Red Hat Cluster Suite to create a high availability infrastructure with failover capabilities."
categories: ["RHEL", "Monitoring", "Linux"]
date: "2012-05-14T08:53:00+02:00"
lastmod: "2012-05-14T08:53:00+02:00"
tags: ["RHEL", "Cluster", "High Availability", "HA", "Red Hat"]
toc: true
---

![Red Hat Cluster Suite](/images/red_hat_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 5/6 |
| **Operating System** | RHEL 5/6 |
| **Website** | [Red Hat Website](https://www.redhat.com) |
| **Last Update** | 14/05/2012 |
{{< /table >}}

## Introduction

Red Hat offers a cluster based on Open Source solutions. This document explains how to set one up.

Here is the type of infrastructure you should have in place to install and configure a cluster:

![Cluster architecture RHEL](/images/cluster_archi_rhel.avif)

Red Hat cluster suite works with a maximum of 16 nodes in a cluster.

Here is the list of services we will cover with their descriptions:

{{< table "table-hover table-striped" >}}
| Service | Description |
|---------|-------------|
| ccsd | Cluster configuration service |
| aisexec | Low-level framework for cluster management (OpenAIS) for Red Hat 5 |
| corosync | Low-level framework for cluster management (Corosync) for Red Hat 6 |
| cman | Cluster manager |
| fenced | Fencing service that enables remote machine reboot |
| DLM | Lock system |
| clvmd | Cluster version of LVM |
| rgmanager | Resources Groups Manager |
| GFS2 | Cluster filesystem |
| ricci | Remote cluster management service |
| lucci | Web frontend connecting to ricci for remote cluster management |
{{< /table >}}

On Red Hat 5 and Red Hat 6, as you can see, different frameworks are used, but this is transparent to us.

## Prerequisites

- You need dedicated interfaces for the cluster.
- The switches must be capable of multicasting on private interfaces.
- Use bonding on your network cards.
- Add as much redundancy as possible.

### ACPI

First, turn off ACPI to avoid problems:

```bash
chkconfig acpid off
service acpid stop
```

Or in grub, add the acpi parameter set to off:

```bash {linenos=table,hl_lines=[3]}
title Red Hat Enterprise Linux (2.6.32-220.el6.x86_64)
	root (hd0,0)
	kernel /vmlinuz-2.6.32-220.el6.x86_64 ro root=/dev/mapper/myvg-rootvol rd_NO_LUKS  KEYBOARDTYPE=pc KEYTABLE=fr LANG=en_US.UTF-8 rd_NO_MD quiet SYSFONT=latarcyrheb-sun16 rhgb crashkernel=auto rd_LVM_LV=myvg/rootvol rd_NO_DM acpi=off
	initrd /initramfs-2.6.32-220.el6.x86_64.img
```

Choose the method that seems simplest to you.

### Hostname

**It is mandatory** that the hostname is correctly configured on all nodes:

```bash
> hostname
node1.deimos.fr
```

If not, modify the line in the following file:

```bash
HOSTNAME=node1.deimos.fr
```

Similarly, check that all nodes are reachable by their DNS, or list them in the hosts file:

```bash
127.0.0.1   localhost localhost.localdomain localhost4 localhost4.localdomain4
::1         localhost localhost.localdomain localhost6 localhost6.localdomain6

# Private interfaces
10.102.2.72     node1
10.102.2.73     node2
10.102.2.74     node3

# Public interfaces
192.168.0.1     server1
192.168.0.2     server2
192.168.0.3     server3
```

The private interfaces (node1, node2 and node3) will be used **for the heartbeat part and therefore dedicated to the cluster**. Public interfaces should be used for the rest (direct connections, VIP...). So **when configuring your cluster, use the private interfaces** for its creation.

### Bonding

**IMPORTANT: The bonding of private interfaces (cluster heartbeat) can only be configured in mode 1! Only this mode is accepted.**

If you are on RHEL 6, create the following line in /etc/modprobe.d/bonding.conf:

```bash
alias bond0 bonding
```

If you are on RHEL <6, you'll need to insert the same content in the /etc/modprobes.conf file.

In /etc/sysconfig/network-script/, create the file ifcfg-bond0 to set up the bonding configuration:

```bash
DEVICE=bond0
IPADDR=192.168.0.104
NETMASK=255.255.255.0
GATEWAY=192.168.0.1
ONBOOT=yes
BOOTPROTO=static
BONDING_OPTS="mode 1"
```

Then we'll configure our network interfaces to tell them they will be used for bonding:

```bash
DEVICE=eth0
MASTER=bond0
ONBOOT=yes
SLAVE=yes
BOOTPROTO=static
```

Do the same for interface 1:

```bash
DEVICE=eth1
MASTER=bond0
ONBOOT=yes
SLAVE=yes
BOOTPROTO=static
```

### Multipathing

Generally, you will use a disk array so that data can be accessed from any machine. For this, you will need to use multipathing. [Follow this documentation to set it up]({{< ref "docs/Servers/Misc/Hardware/multipath_configuring_multiple_paths_for_external_disk_access.md">}}).

### Firewall

If you use a firewall for the interconnection part of your nodes, here is the list of ports:

{{< table "table-hover table-striped" >}}
| Service name | Port/Protocol |
|--------------|---------------|
| cman | 5404/udp, 5405/udp |
| ricci | 11111/tcp |
| gnbd | 14567/tcp |
| modclusterd | 16851/tcp |
| dlm | 21064/tcp |
| ccsd | 50006/tcp, 50007/udp, 50008/tcp, 50009/tcp |
{{< /table >}}

_Note: It is strongly **advised against having a software firewall** on this part and strongly recommended to have a **dedicated switch**!_

## Installation

### Cluster

On all your nodes, install cman and all your dependencies will be installed at once:

```bash
yum install cman ricci openais cman ccs rgmanager lvm2-cluster gfs2-utils
```

Only install what you need of course, you may not need gfs for example.

Then, set a password for ricci so we can connect to it later:

```bash
passwd ricci
```

Then we're going to make the services necessary for the proper functioning of the cluster persistent:

```bash
chkconfig cman on
chkconfig ricci on
chkconfig clvmd on
chkconfig gfs2 on
chkconfig rgmanager on
service cman start
service ricci start
service clvmd start
service gfs2 start
service rgmanager start
```

### Luci

On the client side, preferably use a dedicated machine, you will need to install luci:

```bash
yum install luci
chkconfig luci on
```

The idea is to install this interface on a remote machine for the administration of our cluster.

### Clvmd

Clvmd is the service that allows you to have LVM in a cluster and thus to have identical volume names on all your nodes.

A quorum must absolutely be present if we want Clvm. To activate cluster LVM, here is the command to run on all nodes:

```bash
lvmconf --enable-cluster
```

Otherwise, edit the lvm configuration file to modify this line:

```bash
...
#locking_type = 1
locking_type = 3
...
```

You may need to look at the "preferred_names" parameter if you're using multipathing. Be careful with WWNs, if you use them, you must configure this correctly:

```bash
...
    # If several entries in the scanned directories correspond to the
    # same block device and the tools need to display a name for device,
    # all the pathnames are matched against each item in the following
    # list of regular expressions in turn and the first match is used.
    # preferred_names = [ ]

    # Try to avoid using undescriptive /dev/dm-N names, if present.
    preferred_names = [ "^/dev/mpath/", "^/dev/mapper/mpath", "^/dev/[hs]d" ]
...
```

Then restart and make the service persistent:

```bash
chkconfig clvmd o
service clvmd restart
```

Then you can create LVM volumes as usual. If you are not very familiar with LVM, [check this documentation]({{< ref "docs/Linux/FilesystemsAndStorage/lvm_working_with_logical_volume_management.md">}}).

You can check the status this way:

```bash
> service clvmd status
clvmd (pid  16802) is running...
Clustered Volume Groups: (none)
Active clustered Logical Volumes: (none)
```

Let's quickly create a volume on a node and you can check that it appears everywhere. In this example, multipath is not used. So we're going to create a partition and give it an LVM tag. On all nodes, refresh the partition table on the disk where the new partition was created:

```bash
partprobe /dev/sdb
```

Once done, check on all your nodes:

```bash {linenos=table,hl_lines=[3]}
> fdisk -cul /dev/sdb
Périphérique Amorce  Début        Fin      Blocs     Id  Système
/dev/sdb1            2048     8386181     4192067   8e  Linux LVM
```

Now we're going to create the volume:

```bash
> pvcreate /dev/sdb1
  Writing physical volume data to disk "/dev/sdb1"
  Physical volume "/dev/sdb1" successfully created
> vgcreate shared_vg /dev/sdb1
  Clustered volume group "shared_vg" successfully created
> lvcreate -n shared_lv1 -L 256M shared_vg
> mkfs.ext4 /dev/shared_vg/shared_lv1
```

Our volume is now available on all nodes. **Caution: this does not mean you will be able to write to it simultaneously**. To do this kind of operation, you will need to use [GFS2](#gfs2).

### GFS2

If you need to have a filesystem shared across all nodes, you will need to [install and use GFS2]({{< ref "docs/Servers/HighAvailability/Storage/gfs2-red-hat-cluster-filesystem.md">}}).

## Configuration

### Quorum

The Quorum is an important element in the design of a High Availability cluster. It is mandatory on a single-node cluster and optional beyond. It must be an element common to the nodes (a partition of a disk array for example) and it is used to place locks to guarantee the general proper state of the cluster.

**To avoid split brain situations, the weight must be at least half of the nodes + 1.**

(expected_votes / 2) + 1

By default 1 node has a weight of 1. The "Expected votes" is the number of weights needed for the cluster to function normally.
We will play with this kind of thing when there are application dependencies on other applications not on the same machines. To avoid problems, we will increase the weight of the nodes on the machines.

#### Expected votes

To see the number of votes required for the cluster to function properly:

```bash
> cman_tool status | grep Expected
Expected votes: 2
```

To change the expected votes value for a maintenance operation for example:

```bash
cman_tool -e 2
```

To see the weight on each node:

```bash
ccs_tool lsnode
```

To modify the weight on a given node:

```bash
cman_tool votes -v <votes>
```

#### Qdisk

The Quorum disk (qdisk) is only necessary on a 2-node cluster. Beyond that, it's not necessary and is not recommended!
Qdisk is used by cman and ccsd and must be a partition/LUN, and in no case a logical volume.

##### Configuration

For a good configuration, I invite you to look at these links:

- http://magazine.redhat.com/2007/12/19/enhancing-cluster-quorum-with-qdisk/
- https://access.redhat.com/knowledge/node/2881

You need to properly configure the heartbeat, heuristics and configure the quorum with the right options.

##### Creating the qdisk

To create the quorum, you must have a single partition that is visible on all nodes:

```bash
> mkqdisk -c /dev/sda2 -l label
mkqdisk v0.6.0
Writing new quorum disk label 'qdisk' to /dev/sda2.
WARNING: About to destroy all data on /dev/sda2; proceed [N/y]? y
Initializing status block for node 1...
Initializing status block for node 2...
Initializing status block for node 3...
Initializing status block for node 4...
Initializing status block for node 5...
Initializing status block for node 6...
Initializing status block for node 7...
Initializing status block for node 8...
Initializing status block for node 9...
Initializing status block for node 10...
Initializing status block for node 11...
Initializing status block for node 12...
Initializing status block for node 13...
Initializing status block for node 14...
Initializing status block for node 15...
Initializing status block for node 16...
```

- -c: the partition corresponding to the quorum
- -l: the label created during fdisk

To list active quorums:

```bash
> mkqdisk -L
mkqdisk v0.6.0
/dev/disk/by-id/scsi-S_beaf11-part2:
/dev/disk/by-path/ip-172.17.101.254:3260-iscsi-iqn.2009-10.com.example.cluster1:iscsi-lun-1-part2:
/dev/sda2:
        Magic:                eb7a62c2
        Label:                qdisk
        Created:              Wed Feb 29 16:27:13 2012
        Host:                 node1.deimos.fr
        Kernel Sector Size:   512
        Recorded Sector Size: 512
```

We're going to start the service and make it persistent:

```bash
chkconfig qdiskd on
service qdiskd start
```

Then you'll need to restart the cluster for the quorum to be taken into account.

#### cman

We can configure CMAN when we have a 2-node cluster so that it does not use the quorum:

```xml
<cman expected_votes="1" two_node="1" />
```

So there will be a race to the fencing, to see who will reboot the other one first. All this to avoid splitbrains.

### Fencing

Fencing is mandatory, because a problematic node can corrupt data on mounted partitions. It is therefore preferable that the other nodes fence (yes, that's the verb ;)) the node that may cause problems.
For this there will be a fenced daemon managed by cman. The fenced agents are stored in /sbin/fence\_\*.

We can test our fencing configuration:

```bash
fence_node node2.deimos.fr
```

It therefore reads the configuration with ccs and uses the right agent to fence the node in question.

There are 2 fencing methods:

- STONITH (Shoot The Other Node In The Head): to do the equivalent of a power cut
- Fabric: at the level of a switch or equipment like a disk array

#### Stonith

Here's an example with an HP ILO, to fence a machine:

```bash
/sbin/fence_ilo -a <ip> -l <login> -p <password> -v -o reboot
```

#### Fabric

For SCSI, there is a "scsi_reserve" service that allows to generate a unique key and create records for each machine. Any node can therefore delete the registration that has been made to prevent a problematic machine from continuing to write to SCSI devices.

We can see if it is possible to do it or not (here no because I use iscsi which does not support it):

```bash {linenos=table,hl_lines=[8]}
fence_scsi_test -c

Testing devices in cluster volumes...
DEBUG: pv_name = /dev/sda1

Attempted to register with devices:
-------------------------------------
        /dev/sda1       Failure
-------------------------------------
Number of devices tested: 1
Number of devices passed: 0
Number of devices failed: 1
```

#### Fencing on virtual machines

If your fencing is not done via RACs (Remote Access Card), and you are using Xen or KVM for example, you will need to resort to software fencing by directly calling the hypervisor to shoot down a node. Let's see how to proceed.

We install cman on the host machine if your cluster nodes are on a virtual machine:

```bash
yum install cman
```

And we copy the content of one of my nodes /etc/cluster/\* to my luci machine:

```bash
scp node1.deimos.fr:/etc/cluster/fence_xvm.key /etc/cluster/
```

If this key doesn't exist, you can recreate it like this:

```bash
dd if=/dev/urandom of=/etc/cluster/fence-xvm.key bs=4k count=1
```

Then copy it as explained above.

Then we have "Shared fence devices" which appears in the interface. Now we click on add and choose virtual:

Then we configure the "Main Fencing Method" for each node (in domain).

Failover domains allows you to put a group of machines to a given application so that an application can only start on one of its groups of machines.

If you use Xen and want to be able to reboot a VM, you will need to add this:

```bash
/sbin/fence_xvmd -L -I <cluster_interface>
```

- -L: listens on the network and nodes, but is not a member of the clusters
- -I <cluster_interface>: put the name of the network interface corresponding to the cluster

### Failover domains

Allows you to create groups of nodes on which services will be allowed to switch.

### Example configuration file

```xml
<?xml version="1.0"?>
<cluster alias="cluster1" config_version="16" name="cluster1">
        <fence_daemon clean_start="0" post_fail_delay="0" post_join_delay="3"/>
        <clusternodes>
                <clusternode name="node2.deimos.fr" nodeid="1" votes="1">
                        <fence>
                                <method name="1">
                                        <device domain="xm destroy node1.cluster2.example.com" name="xenfence1"/>
                                </method>
                        </fence>
                </clusternode>
                <clusternode name="node1.deimos.fr" nodeid="2" votes="1">
                        <fence>
                                <method name="1">
                                        <device domain="xm destroy node1.deimos.fr" name="xenfence1"/>
                                </method>
                        </fence>
                </clusternode>
        </clusternodes>
        <cman expected_votes="1" two_node="1"/>
        <fencedevices>
                <fencedevice agent="fence_xvm" name="xenfence1"/>
        </fencedevices>
        <rm>
                <failoverdomains>
                        <failoverdomain name="prefer_node1" nofailback="0" ordered="1" restricted="1">
                                <failoverdomainnode name="node2.deimos.fr" priority="2"/>
                                <failoverdomainnode name="node1.deimos.fr" priority="1"/>
                        </failoverdomain>
                </failoverdomains>
                <resources>
                        <ip address="172.16.50.16" monitor_link="1"/>
                        <apache config_file="conf/httpd.conf" name="httpd" server_root="/etc/httpd" shutdown_wait="5"/>
                        <clusterfs device="/dev/mapper/gfslv-webfs" force_unmount="0" fsid="45506" fstype="gfs2" mountpoint="/var/www/html" name="gfs_www" self_fence="0"/>
                </resources>
                <service autostart="0" domain="prefer_node1" exclusive="0" name="webby" recovery="relocate">
                        <ip ref="172.16.50.16"/>
                        <clusterfs ref="gfs_www">
                                <apache ref="httpd"/>
                        </clusterfs>
                </service>
        </rm>
</cluster>
```

### Manually updating the configuration

You can manually change the cluster configuration:

```
CLUSTER
   |__CMAN or GULM (RHEL4 only)
   |             |__LOCKSERVER (Only for GULM; 1,3,4 or 5 only)
   |
   |__FENCE_XVMD (RHEL5 option, used when hosting a virtual cluster)
   |
   |__CLUSTERNODES
   |         |_____CLUSTERNODE+
   |                     |______FENCE
   |                               |__METHOD+
   |                                      |___DEVICE+
   |
   |__FENCEDEVICES
   |        |______FENCEDEVICE+
   |
   |__RM (Resource Manager Block)
   |   |__FAILOVERDOMAINS
   |   |          |_______FAILOVERDOMAIN*
   |   |                         |________FAILOVERDOMAINNODE*
   |   |__RESOURCES
   |   |      |_____(See Resource List Below)
   |   |__SERVICE*
   |
   |__FENCE_DAEMON
```

Then once you're satisfied, change the version value + 1:

```xml
<cluster alias="cluster1" config_version="16" name="cluster1">
```

Then update this configuration on all nodes:

- On RHEL 5, go to the node in question, then:

```bash
ccs_tool update /etc/cluster/cluster.conf
```

- On RHEL 6:

```bash
ccs -h <host> -p <password> --sync --activate
```

- host: the node on which the action should be done
- password: the ricci password
- --sync: synchronize the configuration file on all nodes
- --activate: activate the new configuration

### cman

View the status of the cluster:

```bash
> cman_tool status
Version: 6.2.0
Config Version: 14
Cluster Name: cluster1
Cluster Id: 26777
Cluster Member: Yes
Cluster Generation: 12
Membership state: Cluster-Member
Nodes: 3
Expected votes: 3
Total votes: 3
Quorum: 2
Active subsystems: 9
Flags: Dirty
Ports Bound: 0 11 177
Node name: node1.deimos.fr
Node ID: 2
Multicast addresses: 239.192.104.2
Node addresses: 10.0.0.1
```

View the status of nodes:

```bash
cman_tool nodes
```

View the status of services:

```bash
cman_tool services
```

Register a node in the cluster:

```bash
cman_tool join
```

To remove a node from the cluster (no services should be running on this node):

```bash
cman_tool leave
```

We can test the fencing (reboot) on a machine:

```bash
cman_tool kill -n node2.deimos.fr
```

### Managing cluster services

Here is the order of starting the services, so if you have a problem, check these services:

```bash
service cman start
service qdiskd start
service clvmd start
service gfs start
service rgmanager start
service ricci start
```

And vice versa for shutdown.

Here's a little script that will only turn on the useful services in case of non-cluster boot or problems:

```bash
#!/bin/sh
for i in cman qdiskd clvmd gfs rgmanager ricci ; do
   if [ $(chkconfig --list | grep $i | grep -c on) -ge 1 ]; then
      service $i start
   fi
done
```

Then run it.

### Preventing a node from joining the cluster at reboot

It is possible to prevent a node from rejoining a cluster like this:

```bash
for i in rgmanagers gfs clvmd qdiskq cman ; do
   chkconfig --level 2345 $i off
done
```

### Cluster shutdown

Be careful if you are using GFS, to completely shut down your cluster, you need to lower the "quorum expected" value to avoid problems when you shut down your nodes one by one. Otherwise you will not be able to unmount your GFS at all!

- Solution 1 (recommended):

For example, on a 3-node cluster, your quorum is at 3 for example, lower your quorum:

```bash
cman_tool expected 2
```

Turn off one of your nodes, lower the quorum again:

```bash
cman_tool expected 1
```

Then you can unmount everything and shut down your nodes.

- Solution 2 (avoid):

You can run the following command which will take each node out of the cluster and therefore we can properly shut them down, but there is a risk of splitbrain:

```bash
cman_tool leave remove
```

## Usage

### Luci

#### Initialization

If you are on Red Hat 5, we need to create the first user for Luci:

```bash
> luci_admin init
Initializing the luci server


Creating the 'admin' user

Enter password:
Confirm password:

Please wait...
The admin password has been successfully set.
Generating SSL certificates...
The luci server has been successfully initialized


You must restart the luci server for changes to take effect.

Run "service luci restart" to do so
```

Then start the service:

```bash
service luci restart
```

Then log in http://luci:8084

- Use admin with the password you just entered via the luci_admin command if you are on RHEL 5
- If you are on RHEL 6, just use the root login and password

### Creating a cluster

To create a cluster, go to one of the future nodes and run the following command:

```bash
ccs -h <node1> --createcluster <cluster>
```

- node1: the name of the node on which the configuration should be created (**use the name of the private interface**)
- cluster: the name of the cluster (no more than 15 characters)
- -f: you can specify the configuration file if you don't want the cluster one to be overwritten

### Adding nodes

To add nodes to a cluster:

```bash
ccs -h node1 --addnode node1
ccs -h node1 --addnode node2
ccs -h node1 --addnode node3
```

These lines will add nodes 1 to 3 to the cluster configuration on node1.

To check the list of nodes:

```bash
> ccs -h node1 --lsnodes
node1: nodeid=1
node2: nodeid=2
node3: nodeid=3
```

### Fencing

We need to add a mandatory fencing method to be able to shoot down a node in case of a problem, so that it can reintegrate the cluster as quickly as possible and does not corrupt data on the disks.

To list the available fencing methods:

```bash
> ccs -h node1 --lsfenceopts
fence_rps10 - RPS10 Serial Switch
fence_vixel - No description available
fence_egenera - No description available
fence_xcat - No description available
fence_na - Node Assassin
fence_apc - Fence agent for APC over telnet/ssh
fence_apc_snmp - Fence agent for APC over SNMP
fence_bladecenter - Fence agent for IBM BladeCenter
fence_bladecenter_snmp - Fence agent for IBM BladeCenter over SNMP
fence_cisco_mds - Fence agent for Cisco MDS
fence_cisco_ucs - Fence agent for Cisco UCS
fence_drac5 - Fence agent for Dell DRAC CMC/5
fence_eps - Fence agent for ePowerSwitch
fence_ibmblade - Fence agent for IBM BladeCenter over SNMP
fence_ifmib - Fence agent for IF MIB
fence_ilo - Fence agent for HP iLO
fence_ilo_mp - Fence agent for HP iLO MP
fence_intelmodular - Fence agent for Intel Modular
fence_ipmilan - Fence agent for IPMI over LAN
fence_kdump - Fence agent for use with kdump
fence_rhevm - Fence agent for RHEV-M REST API
fence_rsa - Fence agent for IBM RSA
fence_sanbox2 - Fence agent for QLogic SANBox2 FC switches
fence_scsi - fence agent for SCSI-3 persistent reservations
fence_virsh - Fence agent for virsh
fence_virt - Fence agent for virtual machines
fence_vmware - Fence agent for VMWare
fence_vmware_soap - Fence agent for VMWare over SOAP API
fence_wti - Fence agent for WTI
fence_xvm - Fence agent for virtual machines
```

To list the possible options on a fencing method:

```bash
> ccs -h node1 --lsfenceopts fence_vmware_soap
 fence_vmware_soap - Fence agent for VMWare over SOAP API
  Required Options:
  Optional Options:
    option: No description available
    action: Fencing Action
    ipaddr: IP Address or Hostname
    login: Login Name
    passwd: Login password or passphrase
    passwd_script: Script to retrieve password
    ssl: SSL connection
    port: Physical plug number or name of virtual machine
    uuid: The UUID of the virtual machine to fence.
    ipport: TCP port to use for connection with device
    verbose: Verbose mode
    debug: Write debug information to given file
    version: Display version information and exit
    help: Display help and exit
    separator: Separator for CSV created by operation list
    power_timeout: Test X seconds for status change after ON/OFF
    shell_timeout: Wait X seconds for cmd prompt after issuing command
    login_timeout: Wait X seconds for cmd prompt after login
    power_wait: Wait X seconds after issuing ON/OFF
    delay: Wait X seconds before fencing is started
    retry_on: Count of attempts to retry power on
```

Here's how to add a fencing method for Vmware:

```bash
ccs -h node1 --addfencedev <vmware_fence> agent=fence_vmware_soap ipaddr=<vcenter.deimos.fr> login=<login> password=<password>
```

- vmware_fence: name you want to give to fencing

For Vmware, we only need one. But if you have RACs, you will need to specify information for each of them.

If you want to remove one:

```bash
ccs -h node1 --rmfencedev <vmware_fence>
```

#### Adding fencing for nodes

Above, we have declared the fencing methods. Now we need to add them to each of our nodes:

```bash
ccs -h node1 --addmethod <vmware_fence> node1
ccs -h node1 --addfenceinst port=node1 ssl=on uuid=412921c1-c259-f9a4-0ee2-37cc047eb4ed
ccs -h node1 --addmethod <vmware_fence> node2
ccs -h node1 --addfenceinst port=node2 ssl=on uuid=422921c1-c259-f9a4-0ee2-37cc047eb4ed
ccs -h node1 --addmethod <vmware_fence> node3
ccs -h node1 --addfenceinst port=node3 ssl=on uuid=432921c1-c259-f9a4-0ee2-37cc047eb4ed
```

_Note: Here's a PowerShell command with PowerCLI to retrieve the UUIDs of VMs:_

```bash
Get-VM <vm_name> | %{(Get-View $_.Id).config.uuid}
```

### Creating a service

--> creation of a typical service

Once our service is created, all that's left is to do the clustat:

```bash
> clustat
Cluster Status for cluster1 @ Tue Feb 29 13:50:09 2012
Member Status: Quorate

 Member Name                                                     ID   Status
 ------ ----                                                     ---- ------
 node2.deimos.fr                                                 1 Online, rgmanager
 node1.deimos.fr                                                 2 Online, Local, rgmanager

 Service Name                                                     Owner (Last)                                                     State
 ------- ----                                                     ----- ------                                                     -----
 service:apache                                                   node2.deimos.fr                                                  started
```

#### NFS

For clustering with NFS, it's a bit special, you need:

- A file system
- NFS export: monnfs
- NFS client:
  - Name: monnfs
  - Target: the IPs allowed to connect
  - Options: ro/rw

As dependencies, here's what it should look like:

- FS
  - NFS Export
    - NFS Client

### rgmanager

The cluster scripts (resource types) available are in OCF (Open Cluster Framework) format and the scripts are available in /usr/share/cluster/.

When you create a service (Resource Group), you can create multiple resources and associate them or create dependencies on each other to define a startup order.
You can retrieve the predefined information related to resource types in /usr/share/cluster/service.sh

### Managing services/RG

Enable a service:

```bash
clusvcadm -e <service> -m <node>
```

Disable a service:

```bash
clusvcadm -d <service>
```

Relocate a service to another node:

```bash
clusvcadm -r <service> -m <node>
```

### Resources Recovery

By default it tries to restart the service on the same node, but there is also:

- relocate: it will try to switch it to another node
- disable: takes no actions in case of problems

#### Check Status

The status must be performed on each resource, which is by default at 30 seconds. You should absolutely not go below 5s. You can modify all these default values in /usr/share/cluster/\*.

The lines to modify look like:

```xml
        <action name="status" interval="30s" timeout="0"/>
        <action name="monitor" interval="30s" timeout="0"/>
```

You can change the check interval and the timeout if you want one.

#### Custom scripts

You can develop your own scripts, and they must be able to respond to start, stop, restart and status. Additionally, return codes must be handled correctly.

## Monitoring

You can monitor in several ways:

- Via the clustat command
- Via SNMP

### Clustat

This command shows the state of the cluster:

```bash
> clustat
Cluster Status for cluster1 @ Thu Mar  1 09:45:21 2012
Member Status: Quorate

 Member Name                                                     ID   Status
 ------ ----                                                     ---- ------
 node2.deimos.fr                                                 1 Online, rgmanager
 node1.deimos.fr                                                 2 Online, Local, rgmanager
 node3.deimos.fr                                                 3 Online, rgmanager

 Service Name                                                     Owner (Last)                                                     State
 ------- ----                                                     ----- ------                                                     -----
 service:websrv                                                   (none)
```

And you can do an XML export:

```xml
> clustat -x

<?xml version="1.0"?>
<clustat version="4.1.1">
  <cluster name="cluster1" id="26777" generation="188"/>
  <quorum quorate="1" groupmember="1"/>
  <nodes>
    <node name="node2.deimos.fr" state="1" local="0" estranged="0" rgmanager="1" rgmanager_master="0" qdisk="0" nodeid="0x00000001"/>
    <node name="node1.deimos.fr" state="1" local="1" estranged="0" rgmanager="1" rgmanager_master="0" qdisk="0" nodeid="0x00000002"/>
    <node name="node3.deimos.fr" state="1" local="0" estranged="0" rgmanager="1" rgmanager_master="0" qdisk="0" nodeid="0x00000003"/>
  </nodes>
  <groups>
    <group name="service:webby" state="119" state_str="disabled" flags="0" flags_str="" owner="none" last_owner="none" restarts="0" last_transition="0" last_transition_str="Thu Jan  1 01:00:00 1970"/>
  </groups>
</clustat>
```

For states, here is the information:

- Started: the service is started
- Pending: the service is starting
- Disabled: the service is disabled
- Stopped: the service is temporarily stopped
- Failed: the service could not start properly

### SNMP

#### Installation

You can query the cluster via SNMP. You will need to install this on the nodes:

```bash
yum install cluster-snmp net-snmp
```

#### Configuration

You can read the file /usr/share/doc/cluster-snmp-0.12.1/README.snmpd which contains the 2 configuration lines for SNMP:

```bash
...
######################################
## cluster monitoring configuration ##
######################################

dlmod RedHatCluster     /usr/lib/cluster-snmp/libClusterMonitorSnmp.so
view    systemview      included        REDHAT-CLUSTER-MIB:RedHatCluster
...
```

And add them to your SNMP configuration file:

```bash
dlmod RedHatCluster     /usr/lib/cluster-snmp/libClusterMonitorSnmp.so
view    systemview      included        REDHAT-CLUSTER-MIB:RedHatCluster
```

Restart your SNMP service. Now we will be able to query the cluster (requires net-snmp-utils for the client):

```bash
> snmpwalk -v1 -c public localhost REDHAT-CLUSTER-MIB:RedHatCluster
REDHAT-CLUSTER-MIB::rhcMIBVersion.0 = INTEGER: 1
REDHAT-CLUSTER-MIB::rhcClusterName.0 = STRING: "cluster1"
REDHAT-CLUSTER-MIB::rhcClusterStatusCode.0 = INTEGER: 4
REDHAT-CLUSTER-MIB::rhcClusterStatusDesc.0 = STRING: "Some services not running"
REDHAT-CLUSTER-MIB::rhcClusterVotesNeededForQuorum.0 = INTEGER: 2
REDHAT-CLUSTER-MIB::rhcClusterVotes.0 = INTEGER: 3
REDHAT-CLUSTER-MIB::rhcClusterQuorate.0 = INTEGER: 1
REDHAT-CLUSTER-MIB::rhcClusterNodesNum.0 = INTEGER: 3
REDHAT-CLUSTER-MIB::rhcClusterNodesNames.0 = STRING: "node1.deimos.fr, node2.deimos.fr, node3.deimos.fr"
REDHAT-CLUSTER-MIB::rhcClusterAvailNodesNum.0 = INTEGER: 3
REDHAT-CLUSTER-MIB::rhcClusterAvailNodesNames.0 = STRING: "node1.deimos.fr, node2.deimos.fr, node3.deimos.fr"
REDHAT-CLUSTER-MIB::rhcClusterUnavailNodesNum.0 = INTEGER: 0
REDHAT-CLUSTER-MIB::rhcClusterUnavailNodesNames.0 = ""
REDHAT-CLUSTER-MIB::rhcClusterServicesNum.0 = INTEGER: 1
REDHAT-CLUSTER-MIB::rhcClusterServicesNames.0 = STRING: "webby"
REDHAT-CLUSTER-MIB::rhcClusterRunningServicesNum.0 = INTEGER: 0
REDHAT-CLUSTER-MIB::rhcClusterRunningServicesNames.0 = ""
REDHAT-CLUSTER-MIB::rhcClusterStoppedServicesNum.0 = INTEGER: 1
REDHAT-CLUSTER-MIB::rhcClusterStoppedServicesNames.0 = STRING: "webby"
REDHAT-CLUSTER-MIB::rhcClusterFailedServicesNum.0 = INTEGER: 0
REDHAT-CLUSTER-MIB::rhcClusterFailedServicesNames.0 = ""
REDHAT-CLUSTER-MIB::rhcNodeName."node1.deimos.fr" = STRING: "node1.deimos.fr"
REDHAT-CLUSTER-MIB::rhcNodeName."node2.deimos.fr" = STRING: "node2.deimos.fr"
REDHAT-CLUSTER-MIB::rhcNodeName."node3.deimos.fr" = STRING: "node3.deimos.fr"
REDHAT-CLUSTER-MIB::rhcNodeStatusCode."node1.deimos.fr" = INTEGER: 0
REDHAT-CLUSTER-MIB::rhcNodeStatusCode."node2.deimos.fr" = INTEGER: 0
REDHAT-CLUSTER-MIB::rhcNodeStatusCode."node3.deimos.fr" = INTEGER: 0
REDHAT-CLUSTER-MIB::rhcNodeStatusDesc."node1.deimos.fr" = STRING: "Participating in cluster"
REDHAT-CLUSTER-MIB::rhcNodeStatusDesc."node2.deimos.fr" = STRING: "Participating in cluster"
REDHAT-CLUSTER-MIB::rhcNodeStatusDesc."node3.deimos.fr" = STRING: "Participating in cluster"
REDHAT-CLUSTER-MIB::rhcNodeRunningServicesNum."node1.deimos.fr" = INTEGER: 0
REDHAT-CLUSTER-MIB::rhcNodeRunningServicesNum."node2.deimos.fr" = INTEGER: 0
REDHAT-CLUSTER-MIB::rhcNodeRunningServicesNum."node3.deimos.fr" = INTEGER: 0
REDHAT-CLUSTER-MIB::rhcNodeRunningServicesNames."node1.deimos.fr" = ""
REDHAT-CLUSTER-MIB::rhcNodeRunningServicesNames."node2.deimos.fr" = ""
REDHAT-CLUSTER-MIB::rhcNodeRunningServicesNames."node3.deimos.fr" = ""
REDHAT-CLUSTER-MIB::rhcServiceName."webby" = STRING: "webby"
REDHAT-CLUSTER-MIB::rhcServiceStatusCode."webby" = INTEGER: 1
REDHAT-CLUSTER-MIB::rhcServiceStatusDesc."webby" = STRING: "stopped"
REDHAT-CLUSTER-MIB::rhcServiceStartMode."webby" = STRING: "manual"
REDHAT-CLUSTER-MIB::rhcServiceRunningOnNode."webby" = ""
End of MIB
```

To get the SNMP cluster hierarchy:

```bash
> snmptranslate -Os -Tp REDHAT-CLUSTER-MIB:RedHatCluster
+--RedHatCluster(8)
   |
   +--rhcMIBInfo(1)
   |  |
   |  +-- -R-- Integer32 rhcMIBVersion(1)
   |
   +--rhcCluster(2)
   |  |
   |  +-- -R-- String    rhcClusterName(1)
   |  +-- -R-- Integer32 rhcClusterStatusCode(2)
   |  +-- -R-- String    rhcClusterStatusDesc(3)
   |  +-- -R-- Integer32 rhcClusterVotesNeededForQuorum(4)
   |  +-- -R-- Integer32 rhcClusterVotes(5)
   |  +-- -R-- Integer32 rhcClusterQuorate(6)
   |  +-- -R-- Integer32 rhcClusterNodesNum(7)
   |  +-- -R-- String    rhcClusterNodesNames(8)
   |  +-- -R-- Integer32 rhcClusterAvailNodesNum(9)
   |  +-- -R-- String    rhcClusterAvailNodesNames(10)
   |  +-- -R-- Integer32 rhcClusterUnavailNodesNum(11)
   |  +-- -R-- String    rhcClusterUnavailNodesNames(12)
   |  +-- -R-- Integer32 rhcClusterServicesNum(13)
   |  +-- -R-- String    rhcClusterServicesNames(14)
   |  +-- -R-- Integer32 rhcClusterRunningServicesNum(15)
   |  +-- -R-- String    rhcClusterRunningServicesNames(16)
   |  +-- -R-- Integer32 rhcClusterStoppedServicesNum(17)
...
```

## FAQ

### Logs

The logs are in /var/log/messages, and the verbosity level can be adjusted in /etc/cluster/cluster.conf.

You can then test your cluster log level like this:

```bash
clulog -s warning "test"
```

### Common errors

Here is a list of the most common errors:

- Poorly written custom scripts that don't return the correct values
- Is the service being checked too often?
- Are resources starting in the right order?

### Starting cman... Can't determine address family of nodename

If you have this kind of message when starting cman:

```bash {linenos=table,hl_lines=[1]}
> service cman start
Starting cluster:
   Checking if cluster has been disabled at boot...        [  OK  ]
   Checking Network Manager...                             [  OK  ]
   Global setup...                                         [  OK  ]
   Loading kernel modules...                               [  OK  ]
   Mounting configfs...                                    [  OK  ]
   Starting cman... Can't determine address family of nodename
Unable to get the configuration
Can't determine address family of nodename
cman_tool: corosync daemon didn't start Check cluster logs for details
```

Check that you have [hostnames configured correctly](#hostname) and that all nodes can communicate with each other.

### Error locking on node

You may encounter this kind of problem if your physical volume is not detected on all nodes.

```bash {linenos=table,hl_lines=[2,3]}
> lvcreate -n shared_lv1 -L 256M shared_vg
  Error locking on node node3: Volume group for uuid not found: EfDxVCE2XWh2Se7ohirVFpXyjXSEJqZxigQMWiiUXNjRXeWd2SzLHwZv3bFropf1
  Error locking on node node1: Volume group for uuid not found: EfDxVCE2XWh2Se7ohirVFpXyjXSEJqZxigQMWiiUXNjRXeWd2SzLHwZv3bFropf1
  Failed to activate new LV.
```

To fix the problem, do a partprobe on the disks containing the new partition.

### Cluster is not quorate. Refusing connection

You need to check that all cluster services have started correctly. We can check the status of the cluster like this:

```bash {linenos=table,hl_lines=[3]}
> clustat
Cluster Status for cluster1 @ Wed Feb 29 14:16:31 2012
Member Status: Quorate
```

And here it is quorate, so no worries :-)

## Resources

- http://sources.redhat.com/cluster/doc/cluster_schema.html
- https://alteeve.com/w/RHCS_v2_cluster.conf
- http://magazine.redhat.com/2007/12/19/enhancing-cluster-quorum-with-qdisk/
- https://access.redhat.com/knowledge/node/2881
- [RHEL Cluster Archi Visio](/others/cluster_archi_rhel.vsd)
