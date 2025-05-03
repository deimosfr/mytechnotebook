---
weight: 999
url: "/SNMP_\\:_Le_protocole_de_gestion_réseaux/"
title: "SNMP: The Network Management Protocol"
description: "A comprehensive guide to understanding and implementing SNMP (Simple Network Management Protocol) on Linux systems including configuration for both v1 and v3 versions, MIB exploration, and usage examples."
categories: ["Network", "Servers", "Red Hat", "Debian"]
date: "2011-12-25T20:37:00+02:00"
lastmod: "2011-12-25T20:37:00+02:00"
tags: ["SNMP", "Network Management", "MIB", "Linux", "Monitoring", "Red Hat", "Debian"]
toc: true
---

## Introduction

Network management systems are based on three main elements: a supervisor, nodes, and agents. In SNMP terminology, the term "manager" is more commonly used than "supervisor". The supervisor is the console that allows the network administrator to execute management requests. Agents are entities located at each interface, connecting the managed equipment (node) to the network and allowing information to be retrieved on different objects.

Switches, hubs, routers, workstations, and servers (physical or virtual) are examples of equipment containing manageable objects. These manageable objects can be hardware information, configuration parameters, performance statistics, and other objects that are directly related to the current behavior of the equipment in question. These objects are classified in a tree-like database called MIB (Management Information Base). SNMP enables communication between the supervisor and agents to collect the desired objects in the MIB.

The network management architecture proposed by the SNMP protocol is therefore based on three main elements:

* The managed devices are network elements (bridges, switches, hubs, routers or servers) containing "managed objects" which can be hardware information, configuration elements or statistical information
* The agents, which are network management applications residing in a device, are responsible for transmitting the local management data of the device in SNMP format
* Network management systems (NMS), which are the consoles through which administrators can perform administration tasks

## Versions

There are 3 versions of the SNMP protocol:

* V1: The first version uses communities to have access to the protocol
* V2: This version suffers from incompatible implementations (no standards, each manufacturer does as they wish)
* V3: Uses USM (User Security Model) to improve security on:
  * Hashed user authentication
  * Encryption of data in transit

## Management of Information Bases

SNMP contains hierarchical information in a database for each device. The data is encapsulated as objects called OIDs represented by:

* A table that can contain multiple values
* A scalar for a single value
* 2 types of integers for:
  * Counters: non-negative integer, increases to max, then values are reset to zero
  * Gauges: negative or non-negative integer, remains at the max value

## Installation

### Client

#### Debian

On Debian, we'll need the snmp package:

```bash
aptitude install snmp
```

#### Red Hat

On Red Hat, we'll need to install the net-snmp-utils package:

```bash
yum install net-snmp-utils
```

### Server

#### Debian

On Debian, we'll need the snmpd package:

```bash
aptitude install snmpd
```

#### Red Hat

On Red Hat, we'll need to install the net-snmp package:

```bash
yum install net-snmp
```

## Configuration

### Server v1

For a server in v1 configuration, you can configure access using the snmpconf command or the configuration file:

```bash
# read only(ro) or write community(rw) | shared secret | source | oid |
# oid: ".1" = everythings
rocommunity       public
rwcommunity       password      192.168.0.0/24     .1
```

The first line allows anybody to access the rocommunity. The second allows write access to the 192.168.0.0/24 range.

If you use the snmpconf command, create a configuration file. But first, you'll need to move the current one because the command might cause problems otherwise.

### Server v3

Version 3 of the SNMP protocol is different from version 1 in its operation. Before starting, we will stop the SNMP server and it's very important that it is stopped for the following steps:

```bash
/etc/init.d/snmpd stop
```

We'll then install the development package to have a very useful tool on Red Hat:

```bash
yum install net-snmp-devel
```

We'll need to create a user to whom we'll assign rights (the password must be >= 8 characters):

```bash
> net-snmp-config --create-snmpv3-user -ro -A auth_passphrase -a sha -X private_passphrase -x AES username
adding the following line to /var/lib/net-snmp/snmpd.conf:
   createUser username SHA "auth_passphrase" AES private_passphrase
adding the following line to /etc/snmp/snmpd.conf:
  rouser username
```

Here we see that the net-snmp-config tool modified the configuration file by adding the username user. It also registered other information in the `/var/lib/net-snmp/snmpd.conf` file.

Now let's create rights for this user. For that we'll need to create:

* A group: we'll define a group to integrate users
* A view: this view will be used to specify what the defined group is authorized to see (in relation to an SNMP tree)
* Access: we map access and authentication methods to the group and the chosen view

```bash
...
rouser username
#       groupName      securityModel securityName
group  mygroup   usm        username
#       name           incl/excl     subtree         mask(optional)
view   myview   included   .1
#       group          context sec.model sec.level prefix read   write  notif
access mygroup any auth exact myview none none
```

Pay attention to the order of insertion of the lines, they are important for the configuration to work.

Then start the SNMP service:

```bash
/etc/init.d/snmpd start
```

## MIBs

MIBs are translated in this form:

![Snmp mib](/images/snmp_mib.avif)

The definition of a MIB is therefore in the following form:

* Prefix: IP-MIB::ipForwarding.0 (the 0 is mandatory for scalar values, otherwise it doesn't work)
* Numeric ID: .1.3.6.1.2.1.4.1.0
* The full name of the object: .iso.org.dod.internet.mgmt.mib-2.ip.ipForwarding.0

The last number is an index corresponding to the value of the OID (indexes working like arrays in Perl):

![Snmp index oid](/images/snmp_index_oid.avif)

### Reading a MIB

For reading, let's take the example of Linux MIBs:

```bash {linenos=table,hl_lines=[5,6,13]}
...
--
-- the IP general group
-- some objects that affect all of IPv4
--
 
ip       OBJECT IDENTIFIER ::= { mib-2 4 }
ipForwarding OBJECT-TYPE    SYNTAX     INTEGER {
                    forwarding(1),    -- acting as a router
                    notForwarding(2)  -- NOT acting as a router
               }
    MAX-ACCESS read-write
    STATUS     current
    DESCRIPTION           "The indication of whether this entity is acting as an IPv4
            router in respect to the forwarding of datagrams received
            by, but not addressed to, this entity.  IPv4 routers forward
            datagrams.  IPv4 hosts do not (except those source-routed
            via the host).
 
            When this object is written, the entity should save the
            change to non-volatile storage and restore the object from
            non-volatile storage upon re-initialization of the system.
            Note: a stronger requirement is not used because this object
            was previously defined."
    ::= { ip 1 }
...
```

* ip OBJECT IDENTIFIER ::= { **mib-2 4** }: Corresponds to the SNMP prefix
* **ipForwarding** OBJECT-TYPE: importing dependencies, such as OBJECT-TYPE
* DESCRIPTION: we have a description

#### snmpget

##### SNMP v1

To retrieve the value of this object, we'll use the snmpget command which is used to retrieve a single value:

```bash
snmpget -v1 -c public localhost IP-MIB::ipForwarding.0
```

* -v: We specify the protocol version here (1)
* -c: the community to use (check the server configuration to know it)

*Note: If you don't get anything, it's probably because you have a permission problem on the server side.*

##### SNMP v3

We've just seen for version 1, now for version 3:

```bash
snmpget -v3 localhost IP-MIB::ipForwarding.0 -l authPriv -u username -A auth_passphrase -a sha -X private_passphrase -x AES
```

* -l: This is the type of security we want for SNMPv3
  * auth: password for hashed and therefore encrypted authentication
  * priv: password for data encryption
  * authPriv: allows using both types of encryption (authentication + data)
  * authNoPriv: Having authentication without data encryption.
* -u: the username to use for authentication
* -A: the passphrase for user authentication
* -a: the hash algorithm to use for authentication
* -X: the passphrase shared with the server
* -x: the encryption algorithm to use for the shared secret

*Note: If you don't get anything, it's probably because you have a permission problem on the server side.*

If you often have multiple requests to make to the same host, you can create a file that can only contain a single host in `~/.snmp/snmp.conf` or `/etc/snmp/snmp.conf`:

```bash
defversion             3
defsecurityname        username
defsecuritylevel       authPriv
defhtype               SHA
defauthpassphrase      auth_passphrase
defprivtype            AES
defprivpassphrase      private_passphrase
```

Now, you no longer need to pass all your arguments, simply the server with the MIB.

#### snmpwalk

Snmpwalk will retrieve all values. You'll need to use grep to find the value you want:

```bash
> snmpwalk -v1 -c public localhost
SNMPv2-MIB::sysDescr.0 = STRING: Linux localhost.localdomain 2.6.32 #1 SMP Wed Nov 9 08:03:13 EST 2011 x86_64
SNMPv2-MIB::sysObjectID.0 = OID: NET-SNMP-MIB::netSnmpAgentOIDs.10
DISMAN-EVENT-MIB::sysUpTimeInstance = Timeticks: (121891) 0:20:18.91
SNMPv2-MIB::sysContact.0 = STRING: Root <root@localhost> (configure /etc/snmp/snmp.local.conf)
SNMPv2-MIB::sysName.0 = STRING: localhost.localdomain
SNMPv2-MIB::sysLocation.0 = STRING: Unknown (edit /etc/snmp/snmpd.conf)
...
```

#### snmpnetstat

This tool will retrieve OIDs and display them like the netstat command:

```bash
snmpnetstat -v1 -c public -Cs localhost
```

### Finding MIB Objects

The snmptranslate command will help us find MIB objects installed locally on the machine (`/usr/share/snmp/mibs/*`). For example, to perform a tree search:

```bash
> snmptranslate -TB '.*forward.*'
SNMP-COMMUNITY-MIB::snmpProxyTrapForwardGroup
SNMP-COMMUNITY-MIB::snmpProxyTrapForwardCompliance
IP-MIB::ipv6InterfaceForwarding
IP-MIB::ipv6IpForwarding
IP-MIB::ipForwarding
IP-FORWARD-MIB::ipForward
IP-FORWARD-MIB::ipForwardTable
IP-FORWARD-MIB::ipForwardEntry
IP-FORWARD-MIB::ipForwardMetric5
IP-FORWARD-MIB::ipForwardMetric4
IP-FORWARD-MIB::ipForwardMetric3
IP-FORWARD-MIB::ipForwardMetric2
IP-FORWARD-MIB::ipForwardMetric1
IP-FORWARD-MIB::ipForwardNextHopAS
IP-FORWARD-MIB::ipForwardInfo
IP-FORWARD-MIB::ipForwardAge
IP-FORWARD-MIB::ipForwardProto
IP-FORWARD-MIB::ipForwardType
IP-FORWARD-MIB::ipForwardIfIndex
IP-FORWARD-MIB::ipForwardNextHop
IP-FORWARD-MIB::ipForwardPolicy
IP-FORWARD-MIB::ipForwardMask
IP-FORWARD-MIB::ipForwardDest
IP-FORWARD-MIB::ipForwardNumber
IP-FORWARD-MIB::ipForwardConformance
IP-FORWARD-MIB::ipForwardCompliances
IP-FORWARD-MIB::ipForwardOldCompliance
IP-FORWARD-MIB::ipForwardCompliance
IP-FORWARD-MIB::ipForwardReadOnlyCompliance
IP-FORWARD-MIB::ipForwardFullCompliance
IP-FORWARD-MIB::ipForwardGroups
IP-FORWARD-MIB::ipForwardMultiPathGroup
IP-FORWARD-MIB::ipForwardCidrRouteGroup
IP-FORWARD-MIB::inetForwardCidrRouteGroup
```

If we want a numeric version:

```bash
> snmptranslate -On IP-FORWARD-MIB::ipForward
.1.3.6.1.2.1.4.24
```

For a complete tree version:

```bash
> snmptranslate -Tp -Of .1.3.6.1.2.1.4.24
+--ipForward(24)
   +-- -R-- Gauge     ipForwardNumber(1)
   |
   +--ipForwardTable(2)
   |  |
   |  +--ipForwardEntry(1)
   |     |  Index: ipForwardDest, ipForwardProto, ipForwardPolicy, ipForwardNextHop
   |     |
   |     +-- -R-- IpAddr    ipForwardDest(1)
   |     +-- CR-- IpAddr    ipForwardMask(2)
   |     +-- -R-- Integer32 ipForwardPolicy(3)
   |     |        Range: 0..2147483647
   |     +-- -R-- IpAddr    ipForwardNextHop(4)
   |     +-- CR-- Integer32 ipForwardIfIndex(5)
   |     +-- CR-- EnumVal   ipForwardType(6)
   |     |        Values: other(1), invalid(2), local(3), remote(4)
   |     +-- -R-- EnumVal   ipForwardProto(7)
   |     |        Values: other(1), local(2), netmgmt(3), icmp(4), egp(5), ggp(6), hello(7), rip(8), is-is(9), es-is(10), ciscoIgrp(11), bbnSpfIgp(12), ospf(13), bgp(14), idpr(15)
   |     +-- -R-- Integer32 ipForwardAge(8)
   |     +-- CR-- ObjID     ipForwardInfo(9)
...
```
