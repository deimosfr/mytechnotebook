---
weight: 999
url: "/Installation_et_configuration_d\\'un_cluster_Heartbeat_V2/"
title: "Installation and Configuration of a Heartbeat V2 Cluster"
description: "Guide for installing, configuring and monitoring services with Heartbeat 2 cluster software for high availability systems."
categories: ["Linux", "Security", "Monitoring"]
date: "2009-02-18T18:37:00+02:00"
lastmod: "2009-02-18T18:37:00+02:00"
tags: ["Heartbeat", "Cluster", "High Availability", "Linux HA"]
toc: true
---

## Introduction

[Heartbeat](https://linux-ha.org) is one of the most widely used cluster solutions because it's very flexible, powerful and stable. We'll see how to proceed with installation, configuration and service monitoring for Heartbeat 2. Compared to version 1, Heartbeat 2 can manage more than 2 nodes.

Unlike Heartbeat version 1, this solution cannot be set up in 30 minutes. You'll need to spend a few hours on it...

## Installation

### Via official repositories

Let's go, as usual, it's relatively simple:

```bash
apt-get install heartbeat-2
```

### Via external packages

For those who don't want to use Debian packages because they're not complete enough, go to [https://download.opensuse.org/repositories/server:/ha-clustering/](https://download.opensuse.org/repositories/server:/ha-clustering/) and download the packages for your distribution (I took the Debian 64-bit ones, for example). Create a folder where you'll put all the packages and navigate to it, then:

```bash
dpkg -i *
apt-get -f install
```

This will install all packages and dependencies.

If you want to use the graphical interface, you'll also need to install these packages:

```bash
apt-get install python-central python python-glade2 python-xml python-gtk2 python-gtkmvc
```

Perform this on both nodes. Here they are called:

- deb-node1
- deb-node2

## hosts

To simplify the configuration and HA transactions, properly configure your `/etc/hosts` file on all nodes:

```
192.168.0.155   deb-node1
192.168.0.150   deb-node2
```

A DNS server can also do the job!

## NTP

Make sure all servers have the same time. It is therefore advisable to synchronize them on the same NTP server (e.g., ntpdate 0.pool.ntp.org)

## sysctl.conf

We're going to modify the management of "core" files. Heartbeat recommends modifying the default configuration by adding the following line in `/etc/sysctl.conf`:

```
kernel/core_uses_pid=1
```

Then apply this configuration:

```bash
/etc/init.d/procps.sh restart
```

## Configuration

First, let's get the example configuration files:

```bash
cd /usr/share/doc/heartbeat-2
gzip -d ha.cf.gz haresources.gz
cp authkeys ha.cf haresources /etc/ha.d/
```

Then we'll start editing the configuration files. Go to `/etc/ha.d/`

### ha.cf

An example of the ha.cf configuration file can be found in `/usr/share/doc/heartbeat/ha.cf.gz`.

Here is the content of the `/etc/ha.d/ha.cf` file once modified:

```
logfacility     local0 # For syslog integration
auto_failback   off # Services don't automatically fail back when a node comes back up
rtprio          5 # Assigns priority to the heartbeat service
deadping        20 # Node problem if it doesn't respond after 20 sec

autojoin        other # Allows other nodes to connect to the cluster
node            deb-node1 # My first node
node            deb-node2 # My second node

crm             on # We allow connections via GUI
apiauth	mgmtd	uid=root
respawn	root	/usr/lib/heartbeat/mgmtd -v

respawn hacluster /usr/lib/heartbeat/ipfail # Launches a command when the node boots
apiauth ipfail gid=haclient uid=root # We authorize root in the haclient group
```

### authkeys

This file is used for information exchange between the different nodes. I can choose sha1 as the encryption algorithm, followed by my passphrase:

```bash
#
#       Authentication file.  Must be mode 600
#
#
#       Must have exactly one auth directive at the front.
#       auth    send authentication using this method-id
#
#       Then, list the method and key that go with that method-id
#
#       Available methods: crc sha1, md5.  Crc doesn't need/want a key.
#
#       You normally only have one authentication method-id listed in this file
#
#       Put more than one to make a smooth transition when changing auth
#       methods and/or keys.
#
#
#       sha1 is believed to be the "best", md5 next best.
#
#       crc adds no security, except from packet corruption.
#               Use only on physically secure networks.
#
auth 1
1 sha1 This is my key !
```

Don't forget to set the correct permissions on this file:

```bash
chmod 0600 authkeys
```

### Replications and tests

Now that our files are ready, we'll upload them to the other nodes:

```bash
scp ha.cf autkeys deb-node2:/etc/ha.d/
```

Now we just need to restart the nodes:

```bash
/etc/init.d/heartbeat restart
ssh root@deb-node2 /etc/init.d/heartbeat restart
```

Now we'll check if everything works:

```bash
crm_mon -i 5
```

The number 5 corresponds to the number of seconds the monitor will try to check the cluster status. If it tries to connect continuously and you have no conclusive results, check your logs:

```
$ tail /var/log/syslog
Jun 25 13:35:40 deb-node1 cib: [2739]: info: #========= Input message message start ==========#
Jun 25 13:35:40 deb-node1 cib: [2739]: info: MSG: Dumping message with 6 fields
Jun 25 13:35:40 deb-node1 cib: [2739]: info: MSG[0] : [t=cib]
Jun 25 13:35:40 deb-node1 cib: [2739]: info: MSG[1] : [cib_op=cib_query]
Jun 25 13:35:40 deb-node1 cib: [2739]: info: MSG[2] : [cib_callid=3]
Jun 25 13:35:40 deb-node1 cib: [2739]: info: MSG[3] : [cib_callopt=256]
Jun 25 13:35:40 deb-node1 cib: [2739]: info: MSG[4] : [cib_clientid=bc414c54-4630-456c-ad8a-6ab01b7e2152]
Jun 25 13:35:40 deb-node1 cib: [2739]: info: MSG[5] : [cib_clientname=2775]
```

Here we encounter **a problem that appears if heartbeat has started and kept an old configuration in memory**. So stop all nodes and delete the contents of the `/var/lib/heartbeat/crm` folder:

```bash
/etc/init.d/heartbeat stop
rm /var/lib/heartbeat/crm/*
/etc/init.d/heartbeat start
```

And now we'll look at the state of our cluster again:

```bash
$ crm_mon -i5

Refresh in 5s...
============
Last updated: Tue Jun 26 06:58:38 2007
Current DC: deb-node2 (18ef92e2-cff0-469b-ab77-a2ef0e45ebd7)
2 Nodes configured.
0 Resources configured.
```

## Configuration of cluster services

The file that will interest us is `/var/lib/heartbeat/crm/cib.xml` (alias CIB).

### crm_config

```
   <configuration>
     <crm_config>
       <cluster_property_set id="cib-bootstrap-options">
         <attributes>
           <nvpair name="default_resource_stickiness" id="options-default_resource_stickiness" value="0"/>
           <nvpair id="options-transition_idle_timeout" name="transition_idle_timeout" value="61s"/>
           <nvpair id="symmetric_cluster" name="symmetric_cluster" value="true"/>
           <nvpair id="no_quorum_policy" name="no_quorum_policy" value="stop"/>
           <nvpair id="default_resource_failure_stickiness" name="default_resource_failure_stickiness" value="0"/>
           <nvpair id="stonith_enabled" name="stonith_enabled" value="false"/>
           <nvpair id="stop_orphan_resources" name="stop_orphan_resources" value="true"/>
           <nvpair id="stop_orphan_actions" name="stop_orphan_actions" value="true"/>
           <nvpair id="remove_after_stop" name="remove_after_stop" value="true"/>
           <nvpair id="is_managed_default" name="is_managed_default" value="true"/>
           <nvpair id="short_resource_names" name="short_resource_names" value="true"/>
         </attributes>
       </cluster_property_set>
     </crm_config>
```

- **default_resource_stickiness**: Do you prefer to keep your service on the current node or move it to one that has more available resources? This option is equivalent to "auto_failback on" except that resources can move to nodes other than the one on which they were activated.

  - _0_: If the value is 0, resources will be automatically assigned.
  - _> 0_: The resource will prefer to return to its original node, but it can also move if this node is not available. A higher value will reinforce the resource to stay where it currently is.
  - _< 0_: The resource will prefer to move elsewhere than where it currently is. A higher value will cause this resource to move.
  - _INFINITY_: The resource will always return to its initial place unless forced (node off, stand by...). This option is equivalent to "auto_failback off" except that resources can move to nodes other than the one on which they were activated.
  - _-INFINITY_: Resources will automatically go to another node.

- **options-transition_idle_timeout (60s by default)**: If no action has been detected during this time, the transition is considered to have failed. If each initialized operation has a higher timeout, then it will be taken into account.

- **symmetric_cluster (true by default)**: Specifies that resources can be launched on any node. Otherwise, you will need to create specific "constraints" for each resource.

- **no_quorum_policy (stop by default)**:

  - _ignore_: Pretends that we have a quorum
  - _freeze_: Does not start any resources not present in our partition. Resources in our partition can be moved to another node with the partition (fencing disabled).
  - _stop_: Stops all resources activated in our partition (fencing disabled)

- **default_resource_failure_stickiness**: Force a resource to migrate after a failure.

- **stonith_enabled (false by default)**: Failed nodes will be fenced.

- **stop_orphan_resources (true by default)**: If a resource is found with no definition.
  - _true_: Stops the action
  - _false_: Ignores the action

This affects more CRM's behavior when the resource is deleted by an admin without stopping it first.

- **stop_orphan_actions (true by default)**: If a recursive action is found and there is no definition:
  - _true_: Stops the action
  - _false_: Ignores the action

This affects more CRM's behavior when the interval for a recurring action is modified.

- **remove_after_stop**: This removes the resource from the "status" section of the CIB.

- **is_managed_default (true by default)**: Unless the resource definition says otherwise:

  - _true_: The resource will be started, stopped, monitored and moved if needed.
  - _false_: The resource will not be started if stopped, stopped if started, and not if it has scheduled recurring actions.

- **short_resource_names (false by default, true recommended)**: This option is for compatibility with versions prior to 2.0.2

### nodes

```
     <nodes>
       <node id="cdd9b426-ba04-4bcd-98a8-76f1d5a8ecb0" uname="deb-node3" type="normal"/>
       <node id="18ef92e2-cff0-469b-ab77-a2ef0e45ebd7" uname="deb-node2" type="normal"/>
       <node id="08f85a1b-291f-46a7-b2d8-cab46788e23d" uname="deb-node1" type="normal"/>
     </nodes>
```

Here we will define our nodes.

- id: It is automatically generated when heartbeat is started and changes depending on several criteria.
- uname: Node names.
- type: normal, member or ping.

If you don't know what to put because the ids haven't been generated yet, put this:

```
      <node/>
```

### resources

```
 26      <resources>
 27        <group id="group_1">
 28          <primitive class="ocf" id="IPaddr_192_168_0_90" provider="heartbeat" type="IPaddr">
 29            <operations>
 30              <op id="IPaddr_192_168_0_90_mon" interval="5s" name="monitor" timeout="5s"/>
 31            </operations>
 32            <instance_attributes id="IPaddr_192_168_0_90_inst_attr">
 33              <attributes>
 34                <nvpair id="IPaddr_192_168_0_90_attr_0" name="ip" value="192.168.0.90"/>
 35              </attributes>
 36            </instance_attributes>
 37            <instance_attributes id="IPaddr_192_168_0_90">
 38              <attributes>
 39                <nvpair id="IPaddr_192_168_0_90-apache2" name="apache2" value="started"/>
 40              </attributes>
 41            </instance_attributes>
 42          </primitive>
 43          <primitive class="lsb" provider="heartbeat" type="apache2" id="apache2_2">
 44            <operations>
 45              <op id="apache2_2_mon" interval="120s" name="monitor" timeout="60s"/>
 46            </operations>
 47            <instance_attributes id="apache2_2">
 48              <attributes>
 49                <nvpair id="apache2_2-deb-node2" name="deb-node2" value="started"/>
 50                <nvpair id="apache2_2-group_1" name="group_1" value="stopped"/>
 51                <nvpair name="target_role" id="apache2_2-target_role" value="started"/>
 52              </attributes>
 53            </instance_attributes>
 54          </primitive>
 55          <instance_attributes id="group_1">
 56            <attributes>
 57              <nvpair id="group_1-deb-node2" name="deb-node2" value="started"/>
 58              <nvpair id="group_1-apache2_2" name="apache2_2" value="stopped"/>
 59            </attributes>
 60          </instance_attributes>
 61        </group>
```

- L26: Here we have created a group. I strongly recommend this to avoid making mistakes. The group is called "group_1".
- L27: Then we insert "IPaddr_192_168_0_90" as a name for the definition of an IP address. The type being "IPaddr".
- L29: We insert a monitoring operation to check every 5 sec if everything is going well with a timeout of 5 sec.
- L34: At the attributes level, we set the virtual IP address we want to use.
- L39: We define an instance that will allow Apache2 (not yet configured) and the IP address to work.
- L43: Now we create the primitive for Apache2
- L45: We will monitor Apache2 every 120s with a timeout of 60s
- L49: Apache2 must start on node2 by default
- L50: We indicate that group
- L51: The state of Apache2 must be started by default

The simplest way to avoid headaches is to use the Heartbeat GUI. It will do the configuration for you :-)

### constraints

```
 98      <constraints>
 99        <rsc_location id="rsc_location_group_1" rsc="group_1">
100          <rule id="prefered_location_group_1" score="100">
101            <expression attribute="#uname" id="prefered_location_group_1_expr" operation="eq" value="deb-node1"/>
102          </rule>
103        </rsc_location>
104        <rsc_location id="cli-prefer-apache2_2" rsc="apache2_2">
105          <rule id="cli-prefer-rule-apache2_2" score="INFINITY">
106            <expression id="cli-prefer-expr-apache2_2" attribute="#uname" operation="eq" value="deb-node2" type="string"/>
107          </rule>
108        </rsc_location>
109      </constraints>
110    </configuration>
111  </cib>
```

The constraints are used to indicate what should start before what. It's up to you to decide according to your needs.

## Management of Cluster Services

To start, I strongly encourage you to [consult this page](https://www.linux-ha.org/v2/AdminTools/crm_resource).

To list services, here is the command:

```bash
crm_resource -L
```

### Preparing a cluster service

Before continuing, make sure the configuration files of the two servers are identical, and that the services are stopped (**here apache2**):

```bash
/etc/init.d/apache2 stop
```

Then make sure that the services managed by Heartbeat are no longer automatically started when Linux boots:

```bash
update-rc.d -f apache2 remove
```

You can now type (on both servers):

```bash
/etc/init.d/heartbeat start
```

After a few moments, the services should normally have started on the first machine, while the other is waiting.

### Switching a service to another node

To switch, for example, our Apache2 service (**apache2_2**) to the second node (**deb-node2**):

```bash
crm_resource -M -r apache2_2 -H deb-node2
```

If you want to switch a service from a node to itself, you'll get the following error:

```
Error performing operation: apache2_2 is already active on deb-node2
```

### Starting a service

Starting a cluster service is quite simple. Here we want to start Apache2 (**apache2_2**)

```bash
crm_resource -r apache2_2 -p target_role -v started
```

If you want to start Apache2 on the second node, you must first reallocate it before starting it ([documentation here](Installation_et_configuration_avancée_d'un_cluster_Heartbeat.html#Basculer_un_service_sur_un_autre_noeud)). Ex:

```bash
crm_resource -M -r apache2_2 -H deb-node2
crm_resource -r apache2_2 -p target_role -v started
```

This will switch apache2_2 which was running (**so currently stopped**) to the second node, but it won't start it. You'll need to run the second line to start it.

### Adding an additional node

So here's the deal, there are several solutions. I'll give you the one that seems best to me. In the `/etc/ha.d/ha.cf` file, check this (otherwise you'll have to reload your heartbeat configuration, which will cause a small outage):

```
autojoin        other # Allows other nodes to connect to the cluster
```

This will authorize nodes that are not entered in the ha.cf file to connect.

For this, a minimum of security is required, which is why you still need to copy the ha.cf and authkeys files to deb-node3 (my new node).
**Note**: First edit the ha.cf file to add the new node. This will allow you to have this node _hardcoded_ at the next startup:

```
autojoin        other # Allows other nodes to connect to the cluster
node            deb-node1 # My first node
node            deb-node2 # My second node
node            deb-node3 # My 3rd node
```

Now we send all this to the new node and don't forget to set the correct permissions:

```bash
cd /etc/ha.d/
scp ha.cf authkeys deb-node3:/etc/ha.d/
ssh deb-node3 chown hacluster:haclient /var/lib/heartbeat/crm/*
```

And now the 3rd node joins the cluster:

```bash
$ /etc/init.d/heartbeat restart
$ crm_mon -i5
Node: deb-node2 (18ef92e2-cff0-469b-ab77-a2ef0e45ebd7): online
Node: deb-node1 (08f85a1b-291f-46a7-b2d8-cab46788e23d): online
Node: deb-node3 (cdd9b426-ba04-4bcd-98a8-76f1d5a8ecb0): online
```

And that's it, it's integrated without interruptions :-)

If you create a new machine by "rsync" from a machine in the cluster, you must delete the following files on the new machine:

```bash
rm /var/lib/heartbeat/hb_uuid
```

## FAQ

[Clusterlab FAQ (Excellent site)](https://www.clusterlabs.org/mw/FAQ)

## Resources

- [Documentation on setting up Heartbeat 2](/pdf/heartbeat_2.pdf)
- [Documentation on Setting up a Web Server with Apache, LVS and Heartbeat 2](/pdf/setup_web_server_cluster.pdf)
- [Documentation on Heartbeat2 Xen cluster with drbd8 and OCFS2](/pdf/heartbeat2_xen_cluster_with_drbd8_and_ocfs2.pdf)
- [Enable high availability for composite applications](/pdf/l-haccmdb-pdf.pdf)
- [https://en.wikipedia.org/wiki/High_availability](https://en.wikipedia.org/wiki/High_availability)
