---
weight: 999
url: "/Installation_of_Sun_Cluster_(old)/"
title: "Installation of Sun Cluster"
description: "Guide for installing and configuring Sun Cluster software on Solaris systems."
categories: ["Solaris", "Cluster"]
date: "2010-07-26T12:01:00+02:00"
lastmod: "2010-07-26T12:01:00+02:00"
tags: ["Cluster", "High Availability", "Sun", "Solaris"]
toc: true
---

## Introduction

This is not very complicated. You just have to follow this documentation

## Installation

For the installation, download the Availability Suite on [Sun Website](https://www.sun.com).

### First Node

After downloading it, then unzip and launch the installer:

```bash
cd ./java_es_05Q4_cluster/Solaris_x86
./installer

Read license and accept it:


   Welcome to the Sun Java(TM) Enterprise System; serious software made
   simple..

   Before you begin, please refer to the Release Notes and Installation Guides
   for Sun Java Enterprise System and Sun Cluster software at http://docs.sun.
   com.


   Copyright 2005 Sun Microsystems, Inc. All rights reserved.
   Use is subject to license terms.

   <Press ENTER to Continue>



   Before you install this product, you must read and accept the entire
   Software License Agreement under which this product is licensed for your use
   .

   <Press ENTER to display the Software License Agreement>
```

Press "Enter"

```bash
   If you have read and accept all the terms of the entire Software License
   Agreement, answer 'yes', and the installation will continue.
   If you do not accept all the terms of the Software License Agreement, answer
   'no', and the installation program will end without installing the product.
   Have you read, and do you accept, all of the terms of the preceding Software
   License Agreement [No] {"<" goes back, "!" exits}? yes
```

Say yes.

```bash
Note that English will always be installed

1.French

2.Spanish

3.Korean

4.Traditional Chinese

5.Simplified Chinese

6.German

7.Japanese

8.English only


   Please enter a comma separated list of languages you would like supported
   with this installation [8] {"<" goes back, "!" exits}
```

Don't choose anything! Stay in english mode:

```bash
Installation Type
-----------------

   Do you want to install the full set of Sun Java(TM) Enterprise System
   Products and Services? (yes/no) [Yes] {"<" goes back, "!" exits} yes
```

Say yes for the full install:

```bash
Shared Component Upgrades Required
-----------------------------------

The shared components listed below are currently installed. They will be
upgraded for compatibility with the products you chose to install.

Component    Package
--------------------
JavaActivationFramework       SUNWjaf
        8.0.0.0 (installed)
        8.1 (required)
JavaMail       SUNWjmail
        8.0.0.0 (installed)
        8.1 (required)

   Enter 1 to upgrade these shared components and 2 to cancel  [1] {"<" goes
   back, "!" exits}: 1
```

Choose upgrade:

Enter 1 to upgrade

```bash
Checking System Status

        Available disk space...    : OK

        Memory installed...        : OK

        Operating system patches...: OK

        Operating system resources...: OK


System ready for installation



   Enter 1 to continue [1] {"<" goes back, "!" exits} 1
```

Enter 1 to continue.

```bash
Screen for selecting Type of Configuration

1. Configure Now - Selectively override defaults or express through

2. Configure Later - Manually configure following installation


   Select Type of Configuration [1] {"<" goes back, "!" exits} 2
```

Continue but configure **later**:

2. Configure Later - Manually configure following installation

```bash
Ready to Install
----------------
The following components will be installed.

Product: Java Enterprise System
Location: /var/sadm/prod/entsys
Space Required: 110.15 MB
-------------------------------
        Sun Cluster 3.1 8/05
           Sun Cluster Core
        Sun Cluster Agents for Sun Java(TM) System
           HA/Scalable Sun Java System Web Server
           HA Sun Java System Application Server
           HA Sun Java System Message Queue
           HA Sun Java System Calendar Server
           HA Sun Java System Administration Server
           HA Sun Java System Directory Server

1. Install
2. Start Over
3. Exit Installation

 What would you like to do [1] {"<" goes back, "!" exits}? 1
```

Choose 1 to install:

1. Install

```bash
Java Enterprise System Sun Cluster
|-1%--------------25%-----------------50%-----------------75%--------------100%|


Installation Complete


Software installation has completed successfully. You can view the installation
summary and log by using the choices below. Summary and log files are available
in /var/sadm/install/logs/.



Your next step is to perform the postinstallation configuration and
verification tasks documented in the Postinstallation Configuration and Startup
Chapter of the Sun Java(TM) Enterprise System Installation Guide. See: http:
//docs.sun.com/doc/819-2328.

   Enter 1 to view installation summary and Enter 2 to view installation logs
   [1] {"!" exits}
```

Now you can do Ctrl+C to exit installation.

Then you have to enter the new PATH:

```bash
PATH=$PATH:/usr/cluster/bin
MANPATH=$MANPATH:/usr/cluster/man:/usr/share/man
```

### Second Node

```bash
    Your responses indicate the following options to scinstall:

      scinstall -ik \
           -C cluster3 \
           -N host08 \
           -A trtype=dlpi,name=bge3 -A trtype=dlpi,name=ce0 \
           -m endpoint=:bge3,endpoint=switch1 \
           -m endpoint=:ce0,endpoint=switch2

    Are these the options you want to use (yes/no) [yes]?

    Do you want to continue with the install (yes/no) [yes]?
```

At the end of installation of the first node:

```bash
sconf -p | grep install
```

For more verbosity:

```bash
sconf -pvv | grep install
```

## Utilities

- scinstall: cluster installation
- scsetup: modify configuration
- scstat: cluster stats < 3.2
- scconf: view / modify configuration < 3.2
- sconf -p: cluster configuration informations
- scrgadm: view / modify RG
- clinfo -n: node number

With the 3.2 version, new commands arrived. You can also use both of command type.

To administrate the cluster, you have 2 web interfaces:

- [https://node:6789](https://node:6789)
- [https://node:3000](https://node:3000)

## Configuration

**All your nodes and hostnames must be entered in `/etc/hosts`.** It's necessary for the configuration

- Informations: ALOM, LOM or RSC are like DRAC cards on Dell.

### First Node

After downloading Sun Cluster, mount or burn the CD. Then launch the installation:

```bash
scinstall
```

```bash
  *** Main Menu ***

    Please select from one of the following (*) options:

      * 1) Install a cluster or cluster node
        2) Configure a cluster to be JumpStarted from this install server
        3) Add support for new data services to this cluster node
        4) Upgrade this cluster node
        5) Print release information for this cluster node

      * ?) Help with menu options
      * q) Quit

    Option:
```

Choose 1:

1. Install a cluster or cluster node

```bash
  *** Install Menu ***

    Please select from any one of the following options:

        1) Install all nodes of a new cluster
        2) Install just this machine as the first node of a new cluster
        3) Add this machine as a node in an existing cluster

       ?) Help with menu options
        q) Return to the Main Menu

    Option:
```

Choose 2:

2. Install just this machine as the first node of a new cluster

```bash
  *** Installing just the First Node of a New Cluster ***


    This option is used to establish a new cluster using this machine as
    the first node in that cluster.

    Once the cluster framework software is installed, you will be asked
    for the name of the cluster. Then, you will have the opportunity to
    run sccheck(1M) to test this machine for basic Sun Cluster
    pre-configuration requirements.

    After sccheck(1M) passes, you will be asked for the names of the
    other nodes which will initially be joining that cluster. Unless this
    is a single-node cluster, you will be also be asked to provide
    certain cluster transport configuration information.

    Press Control-d at any time to return to the Main Menu.


    Do you want to continue (yes/no) [yes]?
```

Answer "yes"!

```bash
  >>> Type of Installation <<<

    There are two options for proceeding with cluster installation. For
    most clusters, a Typical installation is recommended. However, you
    might need to select the Custom option if not all of the Typical
    defaults can be applied to your cluster.

    For more information about the differences between the Typical and
    Custom installation methods, select the Help option from the menu.

    Please select from one of the following options:

        1) Typical
        2) Custom

       ?) Help
        q) Return to the Main Menu

    Option [1]:
```

Choose 2:

2. Custom

```bash
  >>> Software Patch Installation <<<

    If there are any Solaris or Sun Cluster patches that need to be added
    as part of this Sun Cluster installation, scinstall can add them for
    you. All patches that need to be added must first be downloaded into
    a common patch directory. Patches can be downloaded into the patch
    directory either as individual patches or as patches grouped together
    into one or more tar, jar, or zip files.

    If a patch list file is provided in the patch directory, only those
    patches listed in the patch list file are installed. Otherwise, all
    patches found in the directory will be installed. Refer to the
    patchadd(1M) man page for more information regarding patch list files.

    Do you want scinstall to install patches for you (yes/no) [yes]? no
```

Choose no.

```bash
  >>> Cluster Name <<<

    Each cluster has a name assigned to it. The name can be made up of
    any characters other than whitespace. Each cluster name should be
    unique within the namespace of your enterprise.

    What is the name of the cluster you want to establish?
```

Now choose the name of the cluster (eg. ClusterSolaris)

```bash
  >>> Check <<<

    This step allows you to run sccheck(1M) to verify that certain basic
    hardware and software pre-configuration requirements have been met.
    If sccheck(1M) detects potential problems with configuring this
    machine as a cluster node, a report of failed checks is prepared and
    available for display on the screen. Data gathering and report
    generation can take several minutes, depending on system
    configuration.

    Do you want to run sccheck (yes/no) [yes]? yes
```

Choose yes

```bash
sccheck: Requesting explorer data and node report from sola1.
sccheck: sola1: Explorer finished.
sccheck: sola1: Starting single-node checks.
sccheck: sola1: Single-node checks finished.




Press Enter to continue:
```

Now press "Enter"

```bash
  >>> Cluster Nodes <<<

    This Sun Cluster release supports a total of up to 16 nodes.

    Please list the names of the other nodes planned for the initial
    cluster configuration. List one node name per line. When finished,
    type Control-D:

    Node name (Control-D to finish):   Sola2

Unknown host.

    Node name (Control-D to finish):  ^D


    This is the complete list of nodes:

        sola1

    This is a single-node cluster.
    Is that correct (yes/no) [yes]?  yes
```

Enter the name of all the nodes. This node (first one is automatically inserted). Then press Ctrl+d and answer yes.

```bash
  >>> Authenticating Requests to Add Nodes <<<

    Once the first node establishes itself as a single node cluster,
    other nodes attempting to add themselves to the cluster configuration
    must be found on the list of nodes you just provided. You can modify
    this list using scconf(1M) or other tools once the cluster has been
    established.

    By default, nodes are not securely authenticated as they attempt to
    add themselves to the cluster configuration. This is generally
    considered adequate, since nodes which are not physically connected
    to the private cluster interconnect will never be able to actually
    join the cluster. However, DES authentication is available. If DES
    authentication is selected, you must configure all necessary
    encryption keys before any node will be allowed to join the cluster
    (see keyserv(1M), publickey(4)).

    Do you need to use DES authentication (yes/no) [no]?  no
```

If you really want to have encrypted DES authentication say yes. But I choose no, because it's on the private interfaces.

```bash
  >>> Network Address for the Cluster Transport <<<

    The private cluster transport uses a default network address of
    172.16.0.0. But, if this network address is already in use elsewhere
    within your enterprise, you may need to select another address from
    the range of recommended private addresses (see RFC 1918 for details).

    If you do select another network address, bear in mind that the Sun
    Cluster software requires that the rightmost two octets always be
    zero.

    The default netmask is 255.255.0.0. You can select another netmask,
    as long as it minimally masks all bits given in the network address.

    Is it okay to accept the default network address (yes/no) [yes]?  yes

    Is it okay to accept the default netmask (yes/no) [yes]?  yes
```

Say yes to have the default network address.

```bash
  >>> Point-to-Point Cables <<<

    The two nodes of a two-node cluster may use a directly-connected
    interconnect. That is, no cluster transport junctions are configured.
    However, when there are greater than two nodes, this interactive form
    of scinstall assumes that there will be exactly two cluster transport
    junctions.

    Does this two-node cluster use transport junctions (yes/no) [yes]?  yes
```

Say yes.

```bash
  >>> Cluster Transport Junctions <<<

    All cluster transport adapters in this cluster must be cabled to a
    transport junction, or "switch". And, each adapter on a given node
    must be cabled to a different junction. Interactive scinstall
    requires that you identify two switches for use in the cluster and
    the two transport adapters on each node to which they are cabled.

    What is the name of the first junction in the cluster [switch1]?  switch1

    What is the name of the second junction in the cluster [switch2]?  switch2
```

Give here the switches names.

```bash
  >>> Cluster Transport Adapters and Cables <<<

    You must configure at least two cluster transport adapters for each
    node in the cluster. These are the adapters which attach to the
    private cluster interconnect.

    What is the name of the first cluster transport adapter?  vmxnet0
```

Type the name of your interface. May be you have to choose between some interfaces if they are detected (like on vmware)

```bash
Adapter "vmxnet0" is already in use as a public network adapter.

    What is the name of the first cluster transport adapter?  vmxnet1

    Will this be a dedicated cluster transport adapter (yes/no) [yes]?  yes

    All transport adapters support the "dlpi" transport type. Ethernet
    and Infiniband adapters are supported only with the "dlpi" transport;
    however, other adapter types may support other types of transport.
    For more information on which transports are supported with which
    adapters, please refer to the scconf_transp_adap family of man pages
    (scconf_transp_adap_hme(1M), ...).

    Is "vmxnet1" an Ethernet adapter (yes/no) [no]?  yes

    Is "vmxnet1" an Infiniband adapter (yes/no) [no]?  no

    The "dlpi" transport type will be set for this cluster.

    Name of the junction to which "vmxnet1" is connected [switch1]?

    Each adapter is cabled to a particular port on a transport junction.
    And, each port is assigned a name. You can explicitly assign a name
    to each port. Or, for Ethernet and Infiniband switches, you can
    choose to allow scinstall to assign a default name for you. The
    default port name assignment sets the name to the node number of the
    node hosting the transport adapter at the other end of the cable.

    For more information regarding port naming requirements, refer to the
    scconf_transp_jct family of man pages (e.g.,
    scconf_transp_jct_dolphinswitch(1M)).

    Use the default port name for the "vmxnet1" connection (yes/no) [yes]?

    What is the name of the second cluster transport adapter?  vmxnet2

    Will this be a dedicated cluster transport adapter (yes/no) [yes]?

    Name of the junction to which "vmxnet2" is connected [switch2]?

    Use the default port name for the "vmxnet2" connection (yes/no) [yes]?
```

Follow step by step and answer for the network card and switch.

```bash
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

    Is it okay to use this default (yes/no) [yes]?  yes
```

Choose yes if you want to have GFS activated on your cluster.

```bash
/globaldevices is not a directory or file system mount point.
Cannot use "/globaldevices".

    Do you want to select another file system (yes/no) [yes]?  yes

    Do you want to use an already existing file system (yes/no) [yes]?  yes

    What is the name of the file system?  global
```

### Second Node

To see DID local disks:

```bash
scdidadm -l
```

And for all shared and local disks:

```bash
scdidadm -L
```

For the Quorum device (or global device), you have to choose a duplicate device to define it:

```bash
sconf -a -q globaldev=d4
```

or

```bash
scsetup
```

Note: The global device **must** be writable (hard drive is the best quorum device)

To verify disk group imported in a node:

```bash
vxdg list
```

Then to resynchronize

## Others stuffs

### Resynchronize group

This command is for resynchronizing group:

```bash
sconf -c -D name=nfsdg,sync
```

### DPM

The disk path monitoring is a tool for monitor disks. To show all monitored disks:

```bash
scdpm -p all:all
```

To monitor a partition for a node:

```bash
scdpm -m node:d4
```

### Maintenance

If you want to upgrade or do maintenance on a node, you can place it in maintenance mode. All its quorum were deleted. When it will reboot normally, then it will took back its quorum.

```bash
scconf -c -q node=node1 maintstate
```

During a boot, on the OBP:

- ok boot --> node integration in the cluster
- ok boot -x --> don't enter the node in the cluster (this won't load cluster kernel modules)

To shutdown all node of the cluster:

```bash
scshutdown
```

### SVM

SVM are Solaris Volume Manager, it's like LVM on Linux. If the raid is sofware, each hard drive get a partition for MetaDB. MetaDB is the SVM informations replicated on each drives.

If the MetaDB is corrupted on some drives, and the addition of all your hard drive are:

- Total disks < 50% = **PANIC (reboot)**
- Total disks >= 50% = continue to run but if there is a reboot, it could not start. It must have more than 50%

To bypass this limitation of 50%, type this command:

```bash
echo "set md:root_mirror_flag=1" >> /etc/system
```

Mediators are additional votes in crash case to increase the MetaDB percentage. When a mediator is activated, it ask to other nodes to block their own mediator. So it becomes the **golden mediator**.
