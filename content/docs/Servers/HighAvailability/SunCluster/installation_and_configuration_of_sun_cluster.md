---
weight: 999
url: "/Installation_et_configuration_du_SUN_Cluster/"
title: "Installation and Configuration of SUN Cluster"
description: "A comprehensive guide covering the installation and configuration of SUN Cluster, including requirements, network setup, maintenance and troubleshooting procedures."
categories: ["Security", "Database", "Operating Systems"]
date: "2011-11-03T00:24:00+02:00"
lastmod: "2011-11-03T00:24:00+02:00"
tags: ["Solaris", "Cluster", "High Availability", "Sun", "Storage"]
toc: true
---

## Introduction

Solaris Cluster (sometimes Sun Cluster or SunCluster) is a high-availability cluster software product for the Solaris Operating System, created by Sun Microsystems.

It is used to improve the availability of software services such as databases, file sharing on a network, electronic commerce websites, or other applications. Sun Cluster operates by having redundant computers or nodes where one or more computers continue to provide service if another fails. Nodes may be located in the same data center or on different continents.

This Documentation has been released with:

- Sun Solaris update 7
- Sun Cluster 3.2u2

## Requirements

All of the following items are required before installing Sun Cluster. Follow all these steps before installation.

### Hardware

To make a real cluster, here is the required hardware list:

- 2 nodes
  - sun-node1
  - sun-node2
- 4 network cards
  - 2 for Public interface (with IPMP on it)
  - 2 for Private interface (for cluster: heartbeat & nodes information exchange)
- 1 disk array with 1 spare disk

### Partitioning

While you install Solaris, you should make a slice called /globaldevices containing at least 512MB. This slice should be in UFS (ZFS does not work as global device at the moment).

If you didn't create this slice during Solaris installation, you can:

- Use the format command to create a new slice
- Use newfs command to format filesystem as UFS
- Mount this filesystem with global option in /globaldevices
- Add it in /etc/vfstab, for example:

```bash
/dev/did/dsk/d6s3 /dev/did/rdsk/d6s3 /global/.devices/node@2 ufs 2 no global
```

Note: Since Sun Cluster 3.2 update 2, you don't need /globaldevices anymore and use ZFS as default root system.

### Hostname Configuration

Change the hostname to match the cluster nomenclature you want:
[Changing Solaris hostname]({{< ref "docs/Solaris/Misc/changing-hostname-on-solaris.md" >}})

**Do not forget to apply the same /etc/hosts file to all cluster nodes!!! And when you make changes, change it on every node!**

### Patches

**Use Sun Update Manager if you have a graphical interface to update all the available packages. If you don't have graphical interfaces, please install all available patches to avoid installation problems.**

### IPMP Configuration

You need to configure at least 2 interfaces for your public network. Follow this documentation:
[IPMP Configuration]({{< ref "docs/Solaris/Network/ipmp_configuration.md">}})

You don't have to do it for your private network because it will be automatically done by the cluster during installation.

### Activate all network cards

With your 4 network cards, you should activate all your cards to be easily recognized during installation. First, run ifconfig -a to check if all your cards are plumbed. If not, enable them:

```bash
touch /etc/hostname.e1000g2
touch /etc/hostname.e1000g3
ifconfig e1000g2 plumb
ifconfig e1000g3 plumb
```

### Remove RPC and Webconsole binding

If you have installed the latest Solaris version, you may encounter Node integration problems due to RPC binding. This is a new SUN security feature. As we need to allow communication between nodes, we need to disable binding on RPC protocol (and could do it for the webconsole as well). **You should do this operation on each node.**

- Ensure that the local_only property of rpcbind is set to false:

```bash
svcprop network/rpc/bind:default
```

If local_only is set to true, run these commands and refresh the service:

```bash
$ svccfg
svc:> select network/rpc/bind
svc:/network/rpc/bind> setprop config/local_only=false
svc:/network/rpc/bind> quit
svcadm refresh network/rpc/bind:default
```

Now communication between nodes works.

- Ensure that the tcp_listen property of webconsole is set to true:

```bash
svcprop /system/webconsole:console
```

If tcp_listen is not true, run these commands and restart service:

```bash
$ svccfg
svc:> select system/webconsole
svc:/system/webconsole> setprop options/tcp_listen=true
svc:/system/webconsole> quit
/usr/sbin/smcwebserver restart
```

It is needed for Sun Cluster Manager communication. To verify if the port is listening on \*.6789 you can execute:

```bash
netstat -a
```

If you want a faster solution to do those 2 things, use these commands:

```bash
svccfg -s network/rpc/bind setprop config/local_only=false
svcadm refresh network/rpc/bind:default
svccfg -s system/webconsole setprop options/tcp_listen=true
/usr/sbin/smcwebserver restart
```

### Profile

Configure the root profile (~/.profile) or for all users (/etc/profile) by adding these lines:

```bash
PATH=$PATH:/usr/cluster/bin/
```

Now refresh your configuration:

```bash
source ~/.profile
```

or

```bash
source /etc/profile
```

### Ending

Restart all your nodes when everything is finished.

## Installation

First of all, download the [Sun Cluster package](https://www.sun.com) (normally in zip) and uncompress it on all nodes. You should have a "Solaris_x86" folder.

Now launch the installer on all the nodes:

```bash
cd /Solaris_x86
./installer
```

We'll need to install Sun Cluster Core and Quorum (if we want to add more than 2 nodes now or later).

## Configuration

### Wizard configuration

Before launching installation, you should know there are 2 ways to configure all the nodes:

- One by one
- All in one shot

If you want to do all in one shot, you should exchange all root ssh public keys between all nodes.

```bash
scinstall
```

```
  *** Main Menu ***

    Please select from one of the following (*) options:

      * 1) Create a new cluster or add a cluster node
        2) Configure a cluster to be JumpStarted from this install server
        3) Manage a dual-partition upgrade
        4) Upgrade this cluster node
        5) Print release information for this cluster node

      * ?) Help with menu options
      * q) Quit

    Option:
```

Answer: 1

```
  *** New Cluster and Cluster Node Menu ***

    Please select from any one of the following options:

        1) Create a new cluster
        2) Create just the first node of a new cluster on this machine
        3) Add this machine as a node in an existing cluster

        ?) Help with menu options
        q) Return to the Main Menu

    Option:
```

Answer: 1

```
  *** Create a New Cluster ***


    This option creates and configures a new cluster.

    You must use the Java Enterprise System (JES) installer to install the
    Sun Cluster framework software on each machine in the new cluster
    before you select this option.

    If the "remote configuration" option is unselected from the JES
    installer when you install the Sun Cluster framework on any of the new
    nodes, then you must configure either the remote shell (see rsh(1)) or
    the secure shell (see ssh(1)) before you select this option. If rsh or
    ssh is used, you must enable root access to all of the new member
    nodes from this node.

    Press Control-d at any time to return to the Main Menu.


    Do you want to continue (yes/no)
```

Answer: yes

```
  >>> Typical or Custom Mode <<<

    This tool supports two modes of operation, Typical mode and Custom.
    For most clusters, you can use Typical mode. However, you might need
    to select the Custom mode option if not all of the Typical defaults
    can be applied to your cluster.

    For more information about the differences between Typical and Custom
    modes, select the Help option from the menu.

    Please select from one of the following options:

        1) Typical
        2) Custom

        ?) Help
        q) Return to the Main Menu

    Option [1]:
```

Answer: 2

```
  >>> Cluster Name <<<

    Each cluster has a name assigned to it. The name can be made up of any
    characters other than whitespace. Each cluster name should be unique
    within the namespace of your enterprise.

    What is the name of the cluster you want to establish ?
```

Answer: sun-cluster

```
  >>> Cluster Nodes <<<

    This Sun Cluster release supports a total of up to 16 nodes.

    Please list the names of the other nodes planned for the initial
    cluster configuration. List one node name per line. When finished,
    type Control-D:

    Node name:  sun-node1
    Node name:  sun-node2
    Node name (Control-D to finish):  ^D
```

Enter the 2 node names and finish with Ctrl+D.

```
    This is the complete list of nodes:

        sun-node1
        sun-node2

    Is it correct (yes/no) [yes]?
```

Answer: yes

```
    Attempting to contact "sun-node2" ... done

    Searching for a remote configuration method ... done

    The secure shell (see ssh(1)) will be used for remote execution.


Press Enter to continue:
```

Answer: yes

```
  >>> Authenticating Requests to Add Nodes <<<

    Once the first node establishes itself as a single node cluster, other
    nodes attempting to add themselves to the cluster configuration must
    be found on the list of nodes you just provided. You can modify this
    list by using claccess(1CL) or other tools once the cluster has been
    established.

    By default, nodes are not securely authenticated as they attempt to
    add themselves to the cluster configuration. This is generally
    considered adequate, since nodes which are not physically connected to
    the private cluster interconnect will never be able to actually join
    the cluster. However, DES authentication is available. If DES
    authentication is selected, you must configure all necessary
    encryption keys before any node will be allowed to join the cluster
    (see keyserv(1M), publickey(4)).

    Do you need to use DES authentication (yes/no) [no]?
```

Answer: no

```
  >>> Network Address for the Cluster Transport <<<

    The cluster transport uses a default network address of 172.16.0.0. If
    this IP address is already in use elsewhere within your enterprise,
    specify another address from the range of recommended private
    addresses (see RFC 1918 for details).

    The default netmask is 255.255.248.0. You can select another netmask,
    as long as it minimally masks all bits that are given in the network
    address.

    The default private netmask and network address result in an IP
    address range that supports a cluster with a maximum of 64 nodes and
    10 private networks.

    Is it okay to accept the default network address (yes/no) [yes]?
```

Answer: yes

```
    Is it okay to accept the default netmask (yes/no) [yes]?
```

Answer: yes

```
  >>> Minimum Number of Private Networks <<<

    Each cluster is typically configured with at least two private
    networks. Configuring a cluster with just one private interconnect
    provides less availability and will require the cluster to spend more
    time in automatic recovery if that private interconnect fails.

    Should this cluster use at least two private networks (yes/no) [yes]?
```

Answer: yes

```
  >>> Point-to-Point Cables <<<

    The two nodes of a two-node cluster may use a directly-connected
    interconnect. That is, no cluster switches are configured. However,
    when there are greater than two nodes, this interactive form of
    scinstall assumes that there will be exactly one switch for each
    private network.

    Does this two-node cluster use switches (yes/no) [yes]?
```

Answer: no

```
  >>> Cluster Transport Adapters and Cables <<<

    You must configure the cluster transport adapters for each node in the
    cluster. These are the adapters which attach to the private cluster
    interconnect.

    Select the first cluster transport adapter for "sun-node1":

        1) e1000g1
        2) e1000g2
        3) e1000g3
        4) Other

    Option:
```

Answer: 3

```
    Adapter "e1000g3" is an Ethernet adapter.

    Searching for any unexpected network traffic on "e1000g3" ... done
    Verification completed. No traffic was detected over a 10 second
    sample period.

    The "dlpi" transport type will be set for this cluster.

    Name of adapter on "sun-node2" to which "e1000g3" is connected?  e1000g3

    Select the second cluster transport adapter for "sun-node1":

        1) e1000g1
        2) e1000g2
        3) e1000g3
        4) Other

    Option:
```

Answer: 2

```
    Adapter "e1000g2" is an Ethernet adapter.

    Searching for any unexpected network traffic on "e1000g2" ... done
    Verification completed. No traffic was detected over a 10 second
    sample period.

    Name of adapter on "sun-node2" to which "e1000g2" is connected?
```

Answer: e1000g2

```
  >>> Quorum Configuration <<<

    Every two-node cluster requires at least one quorum device. By
    default, scinstall will select and configure a shared SCSI quorum disk
    device for you.

    This screen allows you to disable the automatic selection and
    configuration of a quorum device.

    The only time that you must disable this feature is when ANY of the
    shared storage in your cluster is not qualified for use as a Sun
    Cluster quorum device. If your storage was purchased with your
    cluster, it is qualified. Otherwise, check with your storage vendor to
    determine whether your storage device is supported as Sun Cluster
    quorum device.

    If you disable automatic quorum device selection now, or if you intend
    to use a quorum device that is not a shared SCSI disk, you must
    instead use clsetup(1M) to manually configure quorum once both nodes
    have joined the cluster for the first time.

    Do you want to disable automatic quorum device selection (yes/no) [no]?
```

Answer: yes

```
  >>> Global Devices File System <<<

    Each node in the cluster must have a local file system mounted on
    /global/.devices/node@<nodeID> before it can successfully participate
    as a cluster member. Since the "nodeID" is not assigned until
    scinstall is run, scinstall will set this up for you.

    You must supply the name of either an already-mounted file system or
    raw disk partition which scinstall can use to create the global
    devices file system. This file system or partition should be at least
    512 MB in size.

    If an already-mounted file system is used, the file system must be
    empty. If a raw disk partition is used, a new file system will be
    created for you.

    The default is to use /globaldevices.

    Is it okay to use this default (yes/no) [yes]?
```

Answer: yes

```
    Testing for "/globaldevices" on "sun-node1" ... done

    For node "sun-node2",
       Is it okay to use this default (yes/no) [yes]?
```

Answer: yes

```
    Is it okay to create the new cluster (yes/no) [yes]?
```

Answer: yes

```
  >>> Automatic Reboot <<<

    Once scinstall has successfully initialized the Sun Cluster software
    for this machine, the machine must be rebooted. After the reboot, this
    machine will be established as the first node in the new cluster.

    Do you want scinstall to reboot for you (yes/no) [yes]?
```

Answer: yes

```
    During the cluster creation process, sccheck is run on each of the new
    cluster nodes. If sccheck detects problems, you can either interrupt
    the process or check the log files after the cluster has been
    established.

    Interrupt cluster creation for sccheck errors (yes/no) [no]?
```

Answer: no

```
    The Sun Cluster software is installed on "sun-node2".

    Started sccheck on "sun-node1".
    Started sccheck on "sun-node2".

    sccheck completed with no errors or warnings for "sun-node1".
    sccheck completed with no errors or warnings for "sun-node2".


    Configuring "sun-node2" ... done
    Rebooting "sun-node2" ...

    Waiting for "sun-node2" to become a cluster member ...
```

### Manual configuration

Here is an example of setting up the first node and allowing another node:

```bash
scinstall -i \
-C PA-TLH-CLU-UAT-1 \
-F \
-T node=PA-TLH-SRV-UAT-1,node=PA-TLH-SRV-UAT-2,authtype=sys \
-w netaddr=10.255.255.0,netmask=255.255.255.0,maxnodes=16,maxprivatenets=2,numvirtualclusters=1 \
-A trtype=dlpi,name=nge2 -A trtype=dlpi,name=nge3 \
-B type=switch,name=PA-TLH-SWI-IN-1 -B type=switch,name=PA-TLH-SWI-IN-2 \
-m endpoint=:nge2,endpoint=PA-TLH-SWI-IN-1 \
-m endpoint=:nge3,endpoint=PA-TLH-SWI-IN-2 \
-P task=quorum,state=INIT
```

### Quorum

If you've made the installation with a quorum, you'll need to set it up with the webremote or with these commands. First, you need to list all LUNs with DID format:

```bash
didadm -l
```

```
1        LD-TLH-SRV-UAT-1:/dev/rdsk/c0t0d0 /dev/did/rdsk/d1
2        LD-TLH-SRV-UAT-1:/dev/rdsk/c1t0d0 /dev/did/rdsk/d2
5        LD-TLH-SRV-UAT-1:/dev/rdsk/c2t201500A0B856312Cd31 /dev/did/rdsk/d5
5        LD-TLH-SRV-UAT-1:/dev/rdsk/c2t201400A0B856312Cd31 /dev/did/rdsk/d5
5        LD-TLH-SRV-UAT-1:/dev/rdsk/c3t202400A0B856312Cd31 /dev/did/rdsk/d5
5        LD-TLH-SRV-UAT-1:/dev/rdsk/c3t202500A0B856312Cd31 /dev/did/rdsk/d5
6        LD-TLH-SRV-UAT-1:/dev/rdsk/c5t600A0B800056312C000009CB4AAA2A14d0 /dev/did/rdsk/d6
7        LD-TLH-SRV-UAT-1:/dev/rdsk/c5t600A0B800056381A00000E0A4AAA2A14d0 /dev/did/rdsk/d7
```

Choose the LUN you wish to use for your quorum:

```bash
/usr/cluster/bin/clquorum add /dev/did/rdsk/d6
```

Then, activate it:

```bash
/usr/cluster/bin/clquorum enable /dev/did/rdsk/d6
```

To finish, you need to reset it:

```bash
/usr/cluster/bin/clquorum reset
```

Now you're able to configure your cluster.

### Network

#### Cluster connections

To check cluster interconnect, please use this command:

```bash
clinterconnect status
```

To enable a network card interconnection:

```bash
clinterconnect enable hostname:card
```

example:

```bash
clinterconnect enable localhost:e1000g0
```

#### Check network interconnect interfaces

To check if all interfaces are running, configure IPs on each of the private (cluster) IPs. Then broadcast a ping:

```bash
ping -s 10.255.255.255
```

You can change this private IP address range with your own (defaults are 172.16.0.255)

#### Check traffic

Use the snoop command to see traffic, for example:

```bash
snoop -d <interface> <ip>
```

example:

```bash
snoop -d nge0 192.168.76.2
```

#### Get Fiber Channel WWN

To get Fiber Channel identifiers, run this command:

```bash
fcinfo hba-port
```

```
HBA Port WWN: 2100001b32892934
        OS Device Name: /dev/cfg/c2
        Manufacturer: QLogic Corp.
        Model: 375-3356-02
        Firmware Version: 4.04.01
        FCode/BIOS Version:  BIOS: 1.24; fcode: 1.24; EFI: 1.8;
        Serial Number: 0402R00-0906696990
        Driver Name: qlc
        Driver Version: 20081115-2.29
        Type: N-port
        State: online
        Supported Speeds: 1Gb 2Gb 4Gb
        Current Speed: 4Gb
        Node WWN: 2000001b32892934
HBA Port WWN: 2101001b32a92934
        OS Device Name: /dev/cfg/c3
        Manufacturer: QLogic Corp.
        Model: 375-3356-02
        Firmware Version: 4.04.01
        FCode/BIOS Version:  BIOS: 1.24; fcode: 1.24; EFI: 1.8;
        Serial Number: 0402R00-0906696990
        Driver Name: qlc
        Driver Version: 20081115-2.29
        Type: N-port
        State: online
        Supported Speeds: 1Gb 2Gb 4Gb
        Current Speed: 4Gb
        Node WWN: 2001001b32a92934
```

## Manage

### Get cluster state

To get cluster state, simply run the scstat command:

```bash
scstat
```

```
------------------------------------------------------------------

-- Cluster Nodes --

                    Node name           Status
                    ---------           ------
  Cluster node:     PA-TLH-SRV-PRD-1    Online
  Cluster node:     PA-TLH-SRV-PRD-2    Online
  Cluster node:     PA-TLH-SRV-PRD-3    Online
  Cluster node:     PA-TLH-SRV-PRD-6    Online
  Cluster node:     PA-TLH-SRV-PRD-4    Online
  Cluster node:     PA-TLH-SRV-PRD-5    Online

------------------------------------------------------------------

-- Cluster Transport Paths --

                    Endpoint               Endpoint               Status
                    --------               --------               ------
  Transport path:   PA-TLH-SRV-PRD-1:nge3  PA-TLH-SRV-PRD-6:nge3  Path online
  Transport path:   PA-TLH-SRV-PRD-1:nge2  PA-TLH-SRV-PRD-6:nge2  Path online
  Transport path:   PA-TLH-SRV-PRD-1:nge3  PA-TLH-SRV-PRD-2:nge3  Path online
  Transport path:   PA-TLH-SRV-PRD-1:nge3  PA-TLH-SRV-PRD-5:nge3  Path online
  Transport path:   PA-TLH-SRV-PRD-1:nge2  PA-TLH-SRV-PRD-2:nge2  Path online
  Transport path:   PA-TLH-SRV-PRD-1:nge2  PA-TLH-SRV-PRD-5:nge2  Path online
  Transport path:   PA-TLH-SRV-PRD-1:nge3  PA-TLH-SRV-PRD-4:nge3  Path online
  Transport path:   PA-TLH-SRV-PRD-1:nge2  PA-TLH-SRV-PRD-4:nge2  Path online
  Transport path:   PA-TLH-SRV-PRD-1:nge3  PA-TLH-SRV-PRD-3:nge3  Path online
  Transport path:   PA-TLH-SRV-PRD-1:nge2  PA-TLH-SRV-PRD-3:nge2  Path online
  Transport path:   PA-TLH-SRV-PRD-2:nge2  PA-TLH-SRV-PRD-5:nge2  Path online
  Transport path:   PA-TLH-SRV-PRD-2:nge3  PA-TLH-SRV-PRD-3:nge3  Path online
  Transport path:   PA-TLH-SRV-PRD-2:nge2  PA-TLH-SRV-PRD-4:nge2  Path online
  Transport path:   PA-TLH-SRV-PRD-2:nge3  PA-TLH-SRV-PRD-6:nge3  Path online
  Transport path:   PA-TLH-SRV-PRD-2:nge2  PA-TLH-SRV-PRD-3:nge2  Path online
  Transport path:   PA-TLH-SRV-PRD-2:nge3  PA-TLH-SRV-PRD-5:nge3  Path online
  Transport path:   PA-TLH-SRV-PRD-2:nge2  PA-TLH-SRV-PRD-6:nge2  Path online
  Transport path:   PA-TLH-SRV-PRD-2:nge3  PA-TLH-SRV-PRD-4:nge3  Path online
  Transport path:   PA-TLH-SRV-PRD-3:nge2  PA-TLH-SRV-PRD-5:nge2  Path online
  Transport path:   PA-TLH-SRV-PRD-3:nge3  PA-TLH-SRV-PRD-6:nge3  Path online
  Transport path:   PA-TLH-SRV-PRD-3:nge2  PA-TLH-SRV-PRD-4:nge2  Path online
  Transport path:   PA-TLH-SRV-PRD-3:nge3  PA-TLH-SRV-PRD-5:nge3  Path online
  Transport path:   PA-TLH-SRV-PRD-3:nge3  PA-TLH-SRV-PRD-4:nge3  Path online
  Transport path:   PA-TLH-SRV-PRD-3:nge2  PA-TLH-SRV-PRD-6:nge2  Path online
  Transport path:   PA-TLH-SRV-PRD-6:nge3  PA-TLH-SRV-PRD-5:nge3  Path online
  Transport path:   PA-TLH-SRV-PRD-6:nge2  PA-TLH-SRV-PRD-5:nge2  Path online
  Transport path:   PA-TLH-SRV-PRD-6:nge3  PA-TLH-SRV-PRD-4:nge3  Path online
  Transport path:   PA-TLH-SRV-PRD-6:nge2  PA-TLH-SRV-PRD-4:nge2  Path online
  Transport path:   PA-TLH-SRV-PRD-4:nge2  PA-TLH-SRV-PRD-5:nge2  Path online
  Transport path:   PA-TLH-SRV-PRD-4:nge3  PA-TLH-SRV-PRD-5:nge3  Path online

------------------------------------------------------------------

-- Quorum Summary from latest node reconfiguration --

  Quorum votes possible:      11
  Quorum votes needed:        6
  Quorum votes present:       11


-- Quorum Votes by Node (current status) --

                    Node Name           Present Possible Status
                    ---------           ------- -------- ------
  Node votes:       PA-TLH-SRV-PRD-1    1        1       Online
  Node votes:       PA-TLH-SRV-PRD-2    1        1       Online
  Node votes:       PA-TLH-SRV-PRD-3    1        1       Online
  Node votes:       PA-TLH-SRV-PRD-6    1        1       Online
  Node votes:       PA-TLH-SRV-PRD-4    1        1       Online
  Node votes:       PA-TLH-SRV-PRD-5    1        1       Online


-- Quorum Votes by Device (current status) --

                    Device Name         Present Possible Status
                    -----------         ------- -------- ------
  Device votes:     /dev/did/rdsk/d28s2 5        5       Online

------------------------------------------------------------------

-- Device Group Servers --

                         Device Group        Primary             Secondary
                         ------------        -------             ---------


-- Device Group Status --

                              Device Group        Status
                              ------------        ------


-- Multi-owner Device Groups --

                              Device Group        Online Status
                              ------------        -------------

------------------------------------------------------------------
------------------------------------------------------------------

-- IPMP Groups --

              Node Name           Group   Status         Adapter   Status
              ---------           -----   ------         -------   ------
------------------------------------------------------------------
```

### Registering Resources

You can look at the available resources:

```bash
$ clrt list
SUNW.LogicalHostname:2
SUNW.SharedAddress:2
```

Here we need to use more resources like HA Storage (HAStoragePlus) and GDS:

```bash
clrt register SUNW.HAStoragePlus
clrt register SUNW.gds
```

Now we can verify:

```bash
$ clrt list
SUNW.LogicalHostname:2
SUNW.SharedAddress:2
SUNW.HAStoragePlus:6
SUNW.gds:6
```

### Creating Resource Group

An RG (Resource Group) can contain, for example, a VIP (Virtual IP or Logical Host):

```bash
clrg create sun-rg
```

You can also specify an RG on a specific node:

```bash
clrg create -n sun-node1 sun-rg
```

### Creating Logical Host (VIP) Resource

All your requested VIPs should be in the /etc/hosts file **on each node**, e.g.:

```bash
#
# Internet host table
#
::1             localhost
127.0.0.1       localhost
192.168.0.72    sun-node1
192.168.0.77    sun-node2
192.168.0.79    my_app1_vip
192.168.0.80    my_app2_vip
```

Now, activate it:

```bash
clrslh create -g sun-rg -h my_app1_vip my_app1_vip
```

- sun-rg: Resource group (created before)
- my_app1_vip: name of the VIP in the hosts files
- my_app1_vip: name of the VIP resource in cluster

To specify on only one node:

```bash
clrslh create -g sun-rg -h lh -N ipmp0sun-node1 my_app1_vip
```

### Creating FileSystem Resource

Once your LUNs have been created, be sure you can see all available DIDs on all nodes:

```bash
didadm -l
```

and compare to the 'format' command. Everything should look similar. If not, please run these commands on all nodes:

```bash
didadm -C
didadm -r
```

This will clear all deleted LUNs and add all newly created LUNs in cluster DID configuration.

Now create a zpool for each of your services. Once done, use them as filesystem resource:

```bash
clrs create -g sun-rg -t SUNW.HAStoragePlus -p zpools=my_app1 my_app1-fs
```

- sun-rg: name of the resource group
- my_app1: zpool name
- my_app1-fs: filesystem cluster resource name

### Creating a GDS Resource

A GDS is used to use custom scripts for starting, stopping or probing (status) an application. To integrate a GDS in an RG:

```bash
clrs create -g sun-rg -t SUNW.gds -p Start_command="/bin/myscript.sh start" -p Stop_command="/bin/myscript.sh stop" -p Probe_command="/bin/myscript.sh status" -p resource_dependencies=my_app-fs -p Network_aware=false my_app1-script
```

This will create a GDS with your Zpool as dependency. This means it should be up before the start of the GDS.

Note: You don't need to put the VIP as resource dependency because Sun cluster does it for you by default.

### Modify / view resource properties

You may need to change some properties or get information about them. To show them:

```bash
clrs show -v my_resource
```

And to set a property:

```bash
clrs set -p my_property=value my_resource
```

ex:

```bash
clrs set -p START_TIMEOUT=60 ressource_gds
clrs set -p Probe_command="/mnt/test/bin/service_cluster.pl status my_rg" ressource_gds
```

### Activating Resource Group

To activate the RG:

```bash
clrg manage sun-rg
```

Now if you want to use it (this will activate all the resources in the RG):

```bash
clrg online sun-rg
```

You can also specify a node by adding -n:

```bash
clrg online -n node1 sun-rg
```

## Maintenance

### Boot in non cluster mode

#### Reboot with command line

If you need to enter non-cluster mode, use this command:

```bash
reboot -- -x
```

#### Boot from grub

Simply edit this line by adding -x at the end during server boot:

```
kernel /platform/i86pc/multiboot -x
```

### Remove node from cluster

Simply run this command:

```bash
clnode remove
```

## FAQ

### Can't integrate cluster

#### Solution 1

During installation, if you get this kind of problem:

```
Waiting for "sun-node2" to become a cluster member ...
```

Please follow this step:
[Remove RPC and Webconsole binding](#remove-rpc-and-webconsole-binding)

#### Solution 2

[Remove node configuration](#remove-node-from-cluster) and retry.

### The cluster is in installation mode

If at the end of the installation you encounter this kind of problem (a message like "The cluster is in installation mode" or "Le cluster est en mode installation") this means you need to configure something before configuring your RG or RS.

If you have the WebUI (http://127.0.0.1:6789 for example), you might be able to resolve your problem with it. But in this case, if you may have installed the Quorum, [you need to configure it as well](#quorum).

### How to change Private Interconnect IP for cluster?

The cluster install wanted to use a .0.0 as the private interconnect, and when installed, one private interconnect ended up on 172.16.0 and one ended up on 172.16.1, causing one private interconnect to fault. I found an article that indicated you could edit the cluster configuration by first booting each machine in non-cluster mode (boot-x, I actually did a reboot and then a stop A on the reboot and then a boot -x).

Edit the file /etc/cluster/ccr/infrastructure and then incorporate your changes using:

```bash
/usr/cluster/lib/sc/ccradm -o -i /etc/cluster/ccr/infrastructure
```

After I modified the file to change both private interconnects to be on the 172.16.0 subnet, the second private interconnect came online. Once the second private interconnect came up, I was able to run scsetup, select an additional quorum drive and then set the cluster out of install mode.

### Some commands cannot be executed on a cluster in Install mode

This is generally the case in a 2-node cluster when Quorum is not already set. As described in the man page:

```
   Specify the installation-mode setting for the cluster. You can
   specify either enabled or disabled for the installmode property.

   While the installmode property is enabled, nodes do not attempt to
   reset their quorum configurations at boot time. Also, while in this
   mode, many administrative functions are blocked. When you first
   install a cluster, the installmode property is enabled.

   After all nodes have joined the cluster for the first time, and
   shared quorum devices have been added to the configuration, you must
   explicitly disable the installmode property. When you disable the
   installmode property, the quorum vote counts are set to default
   values. If quorum is automatically configured during cluster
   creation, the installmode property is disabled as well after quorum
   has been configured.
```

However, if you don't want to add a quorum or would like to use it now, simply run this command:

```bash
cluster set -p installmode=disabled
```

### Disk path offline

DID number 3 corresponds to and is reserved for the disk array management and may be seen by the cluster. As it cannot be written (because disk arrays show it in read-only) by the cluster, it shows errors. However, these are not actual errors and you can carefully use your cluster.

#### Method 1

To recover your DID as cleanly as possible, run this command on all the cluster nodes:

```bash
devfsadm
```

Then **if it's the same on all nodes and only if it's like that**, you can safely run this command:

```bash
scgdevs
```

If you get errors like this, please use Method 2:

```bash
$ scgdevs
Configuring DID devices
/usr/cluster/bin/scdidadm: Could not open "/dev/rdsk/c0t0d0s2" to verfiy device ID - Device busy.
/usr/cluster/bin/scdidadm: Could not stat "../../devices/scsi_vhci/disk@g600a0b80005634b400005a334accd4d9:c,raw" - No such file or directory.
Warning: Path node loaded - "../../devices/scsi_vhci/disk@g600a0b80005634b400005a334accd4d9:c,raw".
Configuring the /dev/global directory (global devices)
obtaining access to all attached disks
```

#### Method 2

This second method is manual. You can see a format WWN **command ending with 31**, ex:

```
3       PA-TLH-SRV-PRD-1:/dev/rdsk/c4t600A0B800048A9B6000008304AC37A31d0 /dev/did/rdsk/d3
```

If you really want to disable these kinds of messages, connect to all nodes integrated in the cluster and run this command:

```bash
$ scdidadm -C
scdidadm:  Unable to remove driver instance "3" - No such device or address.
Updating shared devices on node 1
Updating shared devices on node 2
Updating shared devices on node 3
Updating shared devices on node 4
Updating shared devices on node 5
Updating shared devices on node 6
```

Now we can verify everything is ok:

```bash
$ scdidadm -l
1        PA-TLH-SRV-PRD-1:/dev/rdsk/c0t0d0 /dev/did/rdsk/d1
2        PA-TLH-SRV-PRD-1:/dev/rdsk/c1t0d0 /dev/did/rdsk/d2
14       PA-TLH-SRV-PRD-1:/dev/rdsk/c4t600A0B800048A9B6000008304AC37AA3d0 /dev/did/rdsk/d14
15       PA-TLH-SRV-PRD-1:/dev/rdsk/c4t600A0B80005634B40000594B4AC37A8Fd0 /dev/did/rdsk/d15
16       PA-TLH-SRV-PRD-1:/dev/rdsk/c4t600A0B800048A9B60000082E4AC37A6Ad0 /dev/did/rdsk/d16
17       PA-TLH-SRV-PRD-1:/dev/rdsk/c4t600A0B80005634B4000059494AC37A5Dd0 /dev/did/rdsk/d17
18       PA-TLH-SRV-PRD-1:/dev/rdsk/c4t600A0B800048A9B60000082C4AC37A3Dd0 /dev/did/rdsk/d18
19       PA-TLH-SRV-PRD-1:/dev/rdsk/c4t600A0B80005634B4000059474AC37A2Bd0 /dev/did/rdsk/d19
20       PA-TLH-SRV-PRD-1:/dev/rdsk/c4t600A0B800048A9B60000082A4AC37A07d0 /dev/did/rdsk/d20
21       PA-TLH-SRV-PRD-1:/dev/rdsk/c4t600A0B80005634B4000059454AC379F8d0 /dev/did/rdsk/d21
22       PA-TLH-SRV-PRD-1:/dev/rdsk/c4t600A0B800048A9B6000008284AC379B4d0 /dev/did/rdsk/d22
23       PA-TLH-SRV-PRD-1:/dev/rdsk/c4t600A0B80005634B4000059434AC3799Ad0 /dev/did/rdsk/d23
24       PA-TLH-SRV-PRD-1:/dev/rdsk/c4t600A0B80005634B4000058BF4AC1DF47d0 /dev/did/rdsk/d24
25       PA-TLH-SRV-PRD-1:/dev/rdsk/c4t600A0B80005634B4000058BD4AC1DF32d0 /dev/did/rdsk/d25
26       PA-TLH-SRV-PRD-1:/dev/rdsk/c4t600A0B80005634B4000058BB4AC1DF20d0 /dev/did/rdsk/d26
27       PA-TLH-SRV-PRD-1:/dev/rdsk/c4t600A0B80005634B4000058B94AC1DF0Dd0 /dev/did/rdsk/d27
28       PA-TLH-SRV-PRD-1:/dev/rdsk/c4t600A0B80005634B4000007104A9D6B81d0 /dev/did/rdsk/d28
```

DID 3 is not present anymore. If you want to update everything:

```bash
$ scdidadm -r
Warning: DID instance "3" has been detected to support SCSI2 Reserve/Release protocol only. Adding path "PA-TLH-SRV-PRD-3:/dev/rdsk/c3t202300A0B85634B4d31" creates more than 2 paths to this device and can lead to unexpected node panics.
DID subpath "/dev/rdsk/c3t202300A0B85634B4d31s2" created for instance "3".
Warning: DID instance "3" has been detected to support SCSI2 Reserve/Release protocol only. Adding path "PA-TLH-SRV-PRD-3:/dev/rdsk/c2t201200A0B85634B4d31" creates more than 2 paths to this device and can lead to unexpected node panics.
DID subpath "/dev/rdsk/c2t201200A0B85634B4d31s2" created for instance "3".
Warning: DID instance "3" has been detected to support SCSI2 Reserve/Release protocol only. Adding path "PA-TLH-SRV-PRD-3:/dev/rdsk/c3t202200A0B85634B4d31" creates more than 2 paths to this device and can lead to unexpected node panics.
DID subpath "/dev/rdsk/c3t202200A0B85634B4d31s2" created for instance "3".
Warning: DID instance "3" has been detected to support SCSI2 Reserve/Release protocol only. Adding path "PA-TLH-SRV-PRD-3:/dev/rdsk/c2t201300A0B85634B4d31" creates more than 2 paths to this device and can lead to unexpected node panics.
DID subpath "/dev/rdsk/c2t201300A0B85634B4d31s2" created for instance "3".
```

### Force uninstall

This is not recommended, but if you can't uninstall and want to force it, here is the procedure:

- Stop all cluster nodes (`scshutdown -y -g 0`) and start them again [in non-cluster mode](#boot-in-non-cluster-mode)

```bash
ok boot -x
```

- Remove the Sun Cluster packages

```bash
pkgrm SUNWscu SUNWscr SUNWscdev SUNWscvm SUNWscsam SUNWscman SUNWscsal SUNWmdm
```

- Remove the configurations

```bash
rm -r /var/cluster /usr/cluster /etc/cluster
rm /etc/inet/ntp.conf
rm -r /dev/did
rm -r /devices/pseudo/did*
rm /etc/path_to_inst (make sure you have backup copy of this file)
```

ATTENTION: If you create a new path_to_inst at boottime with 'boot -ra' you should be on the physical boot device. It may not be possible to write a path_to_inst on a boot mirror (SVM or VxVM).

- Edit configuration files

  - edit /etc/vfstab to remove did and global entries
  - edit /etc/nsswitch.conf to remove cluster references

- Reboot the node with -a option (necessary to write a new path_to_inst file)

```bash
reboot -- -rav
```

reply "y" to "do you want to rebuild path_to_inst?"

- In case of reinstalling, then...

```bash
mkdir /globaldevices; rmdir /global
```

- Uncomment /globaldevices entry from /etc/vfstab
- newfs /dev/rdsk/c?t?d?s? (wherever /globaldevices was mounted)
- mount /globaldevices
- scinstall

### How to Change Sun Cluster Node Names

Make a copy of /etc/cluster/ccr/infrastructure:

```bash
cp /etc/cluster/ccr/infrastructure /etc/cluster/ccr/infrastructure.old
```

Edit /etc/cluster/ccr/infrastructure and change node names as you want. For example, change srv01 to server01 and srv02 to server02.

If necessary, change the Solaris node name:

```bash
echo server01 > /etc/nodename
```

Regenerate the checksum for the infrastructure file:

```bash
/usr/cluster/lib/sc/ccradm -i /etc/cluster/ccr/infrastructure -o
```

Shut down Sun Cluster and boot both nodes:

```bash
cluster shutdown -g 0 -y ok boot
```

### Can't switch an RG from one node to another

I had a problem switching an RG on a Solaris 10u7 with Sun Cluster 3.2u2 (installed patches: 126107-33, 137104-02, 142293-01, 141445-09). The ZFS volume wouldn't mount on another node. In the `/var/adm/messages` file, I saw this message when trying to mount an RG:

```
Dec 16 15:34:30 LD-TLH-SRV-PRD-3 zfs: [ID 427000 kern.warning] WARNING: pool 'ulprod-ld_mysql' could not be loaded as it was last accessed by another system (host: LD-TLH-SRV-PRD-2 hostid: 0x27812152). See: http://www.sun.com/msg/ZFS-8000-EY
```

In fact, it's a bug that can be bypassed by putting the RG offline:

```bash
clrg offline <RG_name>
```

Then manually mount and unmount the zpool:

```bash
zpool import <zpool_name>
zpool export <zpool_name>
```

Now put the RG online:

```bash
clrg online -n <node_name> <rg_name>
```

If the problem still occurs, look in the log files and if you see something like this:

```
Aug 17 22:23:28 minipardus SC[,SUNW.HAStoragePlus:8,clstorage,zfspool,hastorageplus_prenet_start]: [ID 148650 daemon.notice] Started searching for devices in '/dev/dsk' to find the importable pools.
Aug 17 22:23:35 minipardus SC[,SUNW.HAStoragePlus:8,clstorage,zfspool,hastorageplus_prenet_start]: [ID 547433 daemon.notice] Completed searching the devices in '/dev/dsk' to find the importable pools.
Aug 17 22:23:35 minipardus SC[,SUNW.HAStoragePlus:8,clstorage,zfspool,hastorageplus_prenet_start]: [ID 471757 daemon.error] cannot import pool 'qnap': '/var/cluster/run/HAStoragePlus/zfs' is not a valid directory
Aug 17 22:23:35 minipardus SC[,SUNW.HAStoragePlus:8,clstorage,zfspool,hastorageplus_prenet_start]: [ID 117328 daemon.error] The pool 'qnap' failed to import and populate cachefile.
Aug 17 22:23:35 minipardus SC[,SUNW.HAStoragePlus:8,clstorage,zfspool,hastorageplus_prenet_start]: [ID 292307 daemon.error] Failed to import:qnap
```

If that's the case, it's apparently fixed in Sun Cluster 3.2u3.

To avoid installing this update, create this folder '/var/cluster/run/HAStoragePlus/zfs':

```bash
mkdir -p /var/cluster/run/HAStoragePlus/zfs
```

Check if file "/etc/cluster/eventlog/eventlog.conf" contains the line "EC_zfs - - - /usr/cluster/lib/sc/events/zpool_cachefile_plugin.so".

If it's not the case, the content should look like:

```

# Class         Subclass        Vendor  Publisher       Plugin location                         Plugin parameters

EC_Cluster      -               -       -               /usr/cluster/lib/sc/events/default_plugin.so
EC_Cluster      -               -       gds             /usr/cluster/lib/sc/events/gds_plugin.so
EC_Cluster      -               -       -               /usr/cluster/lib/sc/events/commandlog_plugin.so
EC_zfs          -               -       -               /usr/cluster/lib/sc/events/zpool_cachefile_plugin.so
```

Now mount the RG where you want, it should work.

### Cluster is unavailable when a node crashes on a 2-node cluster

Two types of problems can arise from cluster partitions: **split brain and amnesia**. Split brain occurs when the cluster interconnect between Solaris hosts is lost and the cluster becomes partitioned into subclusters, and each subcluster believes that it is the only partition. A subcluster that is not aware of the other subclusters could cause a conflict in shared resources, such as duplicate network addresses and data corruption.

Amnesia occurs if all the nodes leave the cluster in staggered groups. An example is a two-node cluster with nodes A and B. If node A goes down, the configuration data in the CCR is updated on node B only, and not node A. If node B goes down at a later time, and if node A is rebooted, node A will be running with old contents of the CCR. This state is called amnesia and might lead to running a cluster with stale configuration information.

You can avoid split brain and amnesia by giving each node one vote and mandating a majority of votes for an operational cluster. A partition with the majority of votes has a quorum and is enabled to operate. This majority vote mechanism works well if more than two nodes are in the cluster. In a two-node cluster, a majority is two. If such a cluster becomes partitioned, an external vote enables a partition to gain quorum. This external vote is provided by a quorum device. A quorum device can be any disk that is shared between the two nodes.

#### Recovering from amnesia

Scenario: Two node cluster (nodes A and B) with one Quorum Device, nodeA has gone bad, and amnesia protection is preventing nodeB from booting up.

Amnesia occurs if all the nodes leave the cluster in staggered groups. An example is a two-node cluster with nodes A and B. If node A goes down, the configuration data in the CCR is updated on node B only, and not node A. If node B goes down at a later time, and if node A is rebooted, node A will be running with old contents of the CCR. This state is called amnesia and might lead to running a cluster with stale configuration information.

**Warning: this is a dangerous operation**

- Boot nodeB in non-cluster mode:

```bash
reboot -- -x
```

- Edit nodeB's file /etc/cluster/ccr/global/infrastructure as follows:
  - Change the value of "cluster.properties.installmode" from "disabled" to "enabled"
  - Change the number of votes for nodeA from "1" to "0", in the property line "cluster.nodes.<NodeA's id>.properties.quorum_vote".
  - Delete all lines with "cluster.quorum_devices" to remove knowledge of the quorum device.

```
...
cluster.properties.installmode  enabled
...
cluster.nodes.1.properties.quorum_vote  1
...
```

- On the first node (the master one, the first to boot) run:

```bash
/usr/cluster/lib/sc/ccradm -i /etc/cluster/ccr/infrastructure -o
```

or (depending on version)

```bash
/usr/cluster/lib/sc/ccradm recover -o /etc/cluster/ccr/global/infrastructure
```

- Reboot nodeB in cluster mode:

```bash
reboot
```

If you have more than 2 nodes, use the same command but without "-o":

```bash
/usr/cluster/lib/sc/ccradm recover /etc/cluster/ccr/global/infrastructure
```

## References

[Installation of Sun Cluster (old)]({{< ref "docs/Servers/HighAvailability/SunCluster/installation_of_sun_cluster.md">}})
[https://en.wikipedia.org/wiki/Solaris_Cluster](https://en.wikipedia.org/wiki/Solaris_Cluster)  
[https://opensolaris.org/os/community/ha-clusters/translations/french/relnote_fr/](https://opensolaris.org/os/community/ha-clusters/translations/french/relnote_fr/)  
[Resources Properties](https://docs.sun.com/app/docs/doc/819-2974/6n57pdk2o?a=view#indexterm-535)  
[https://docs.sun.com/app/docs/doc/819-0177/cbbbgfij?l=ja&a=view](https://docs.sun.com/app/docs/doc/819-0177/cbbbgfij?l=ja&a=view)  
[https://www.vigilanttechnologycorp.com/genasys/weblogRender.jsp?LogName=Sun%20Cluster](https://www.vigilanttechnologycorp.com/genasys/weblogRender.jsp?LogName=Sun%20Cluster)  
[https://docs.sun.com/app/docs/doc/820-2558/gdrna?l=fr&a=view](https://docs.sun.com/app/docs/doc/820-2558/gdrna?l=fr&a=view)  
[https://wikis.sun.com/display/SunCluster/%28English%29+Sun+Cluster+3.2+1-09+Release+Notes#%28English%29SunCluster3.21-09ReleaseNotes-optgdfsinfo](https://wikis.sun.com/display/SunCluster/%28English%29+Sun+Cluster+3.2+1-09+Release+Notes#%28English%29SunCluster3.21-09ReleaseNotes-optgdfsinfo) (Gold Mine)  
[Deploying highly available zones with Solaris Cluster 3.2](/pdf/deploying_highly_available_zones_with_solaris_cluster_3.2.pdf)
