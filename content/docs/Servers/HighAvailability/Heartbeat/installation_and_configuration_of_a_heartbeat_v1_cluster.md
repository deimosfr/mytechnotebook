---
weight: 999
url: "/Installation_et_configuration_d'un_cluster_Heartbeat_V1/"
title: "Installation and Configuration of a Heartbeat V1 Cluster"
description: "This guide explains how to set up and configure a Heartbeat V1 cluster for high availability services on Linux."
categories: ["Linux", "CentOS", "Ubuntu"]
date: "2009-02-18T18:31:00+02:00"
lastmod: "2009-02-18T18:31:00+02:00"
tags: ["Heartbeat", "High Availability", "Cluster", "Service Failover"]
toc: true
---

## Introduction

The main idea to ensure service availability is to have multiple machines (at least two) running simultaneously. These machines form what we call a cluster, and each machine is a node in the cluster. Each machine will check if the others are still responding by taking their pulse. If a machine stops working, the others will ensure the service.

Once the cluster is configured, it is accessed through a single unique IP address, which is the cluster's IP address; the cluster itself consists of multiple nodes.

To implement this kind of technique, we will use the HeartBeat application, which will monitor the machines and apply a series of scripts defined by the user if necessary (i.e. if a machine fails).

## Preliminary Information

The test configuration consists of two modest machines; twin1 and twin2.

The cluster will be accessible through the address 192.168.252.205.

### Twin 1

- Master / Primary
- Pentium 3 - 500 MHz - 128MB RAM
- Fresh minimal installation of Debian
- Kernel upgraded to 2.6.18-4-686 (instead of 386)
- IP: 192.168.252.200
- Hostname: twin1
- 13.5GB disk split into 5 partitions (/ of 4.5GB, /srv/prod of 3GB, /srv/intern of 3GB, /srv/other of 2GB and swap of 1GB)

### Twin 2

- Slave / Secondary
- Pentium 3 - 700 MHz - 128MB RAM
- Fresh minimal installation of Debian Stable
- Kernel upgraded to 2.6.18-4-686 (instead of 386)
- IP: 192.168.252.201
- Hostname: twin2
- 10GB disk split into 4 partitions (/ of 3GB, /srv/prod of 3GB, /srv/intern of 3GB and swap of 1GB)

## Machine Status

On both machines, here's what has been done.

- Installed Debian 4.0 CD and manually configured the network.
- Updated some basic packages via the Internet.

```bash
sudo apt-get update
sudo apt-get dist-upgrade
```

- Upgraded the kernel to take advantage of the Pentium 3 (not necessary). After the upgrade, a reboot to use the new kernel.

```bash
sudo apt-get install linux-image-2.6.18-686
sudo reboot
```

## Installing HeartBeat

HeartBeat is very simple to install on Debian. The package can be found in the standard repositories. All you need to do is run the following command on both machines (Twin 1 and Twin 2):

```bash
sudo apt-get install heartbeat-2
```

Once the service is installed, you will get an error saying that the ha.cf configuration file is not present. This is perfectly normal.

Now let's look at the basic configuration of HeartBeat.

## Basic Configuration

HeartBeat's configuration is based on 3 fundamental files. During the basic configuration, I will present a minimalist configuration with the goal of:

- providing a minimalist base that you can adapt to your needs over time.

The three configuration files are strictly identical on both machines; in your case, they should be copied to each machine in the cluster. In our example, this means Twin 1 and Twin 2.

### ha.cf

```bash
sudo vi /etc/ha.d/ha.cf
```

Here is the ha.cf file I use on Twin 1 and Twin 2:

```bash
bcast         eth0

debugfile     /var/log/ha-debug
logfile       /var/log/ha-log
logfacility   local0

keepalive     2
deadtime      10
warntime      6
initdead      60

udpport       694
node          twin1
node          twin2

auto_failback off

apiauth         mgmtd   uid=root
crm             on
respawn         root    /usr/lib/heartbeat/mgmtd -v
```

Let's examine this configuration file in detail:

- bcast: indicates the network interface through which we will take the pulse.
- debugfile: indicates the debugging file to use.
- logfile: indicates the activity log to use.
- logfacility: indicates that we are also using the syslog facility.
- keepalive: indicates the delay between two pulse beats. This delay is given by default in seconds; to use milliseconds, add the ms suffix.
- deadtime: indicates the time needed before considering a node as dead. This time is given by default in seconds; to use milliseconds, add the ms suffix.
- warntime: indicates the delay before sending a warning for late pulses. This delay is given by default in seconds; to use milliseconds, add the ms suffix.
- initdead: indicates a specific deadtime for configurations where the network takes some time to start. initdead is normally at least twice the deadtime. This delay is given by default in seconds; to use milliseconds, add the ms suffix.
- udpport: indicates the port to use for taking the pulse.
- node: specifies the names of the machines that are part of the cluster. This name must be identical to that returned by the uname -n command.
- auto_failback: indicates the behavior to adopt if the master node returns to the cluster. If set to on, when the master node returns to the cluster, everything is transferred to it. If set to off, services continue to run on the slave even when the master returns to the cluster.
- apiauth: indicates that the root uid has the right to administer remotely via the GUI (**this also implies adding root to the haresources group**)
- crm: allows remote connections
- respawn: the binary to launch to open the port for remote control

Personally, I prefer the off value for auto_failback so that I can manually return to normal when the production load is less important. Moreover, if the master node suffers from a serious problem, we could have a start/stop loop that would result in a continuous relay of services. This can become problematic.

### haresources

Now, we will define the master node, the cluster's IP address, and the services that need to be assured. Initially (for basic configuration), we will only set up the cluster's IP address.

To configure this aspect of HeartBeat, we edit the haresources file using the following command:

```bash
sudo vi /etc/ha.d/haresources
```

**CAUTION: the resource configuration file must be STRICTLY identical on each node.**

In our example, the resource configuration file looks like this:

```bash
twin1 IPaddr::192.168.252.205
```

This means we define the master node as twin1 and the cluster IP address as 192.168.252.205. We'll stop here for the basic configuration, but I invite you to read the advanced configuration section for more information.

### authkeys

The authkeys file allows cluster nodes to identify each other. This file is edited via the command:

```bash
sudo vi /etc/ha.d/authkeys
```

**CAUTION: the authentication key configuration file must be STRICTLY identical on each node.**

In our example, the authkeys file looks like this:

```bash
auth 2
1 md5 "cluster twin test"
2 crc
```

The auth keyword indicates which identification system we will use. If the link is not physically secure, it is necessary to use md5 encryption with a more intelligent key string than this one.

In our case, there are only the two twins on a small local network, so we use the crc system.

Once the file is edited, don't forget to protect it by entering the command:

```bash
sudo chmod 600 /etc/ha.d/authkeys
```

### Starting Up and First Tests

Before making our first tests, I suggest installing an SSH server on each of the servers.

Check if SSH connections to each node via their own IP addresses (192.168.252.200 and 192.168.252.201 in our example) work. If so, we can do a small test.

Start the HeartBeat service on the master node; namely Twin 1 via the following command:

```bash
sudo /etc/init.d/heartbeat start
```

Then, using the same command, activate HeartBeat on Twin 2.

You can now try an SSH connection to the cluster address; namely: 192.168.252.205 in our example. Once identified, you notice that the command prompt tells you that you are on twin1.

Now, firmly remove the power cord from Twin 1.

Your SSH console will report an error (normal... ;-) ).

Restart an SSH connection to the cluster address; namely: 192.168.252.205 in our example. The service responds! Yet, Twin 1 is stopped. Once identified, you notice that the command prompt indicates:

```bash
deimos@twin2:~$
```

Everything works. Now, let's imagine that Twin 1 comes back into the cluster (plug Twin 1's power back in).

The cluster still runs on Twin 2. To force the switch to Twin 1 (since we are in auto_failback off mode), we enter the following command on Twin 2:

```bash
sudo /etc/init.d/heartbeat restart
```

And Twin 1 automatically takes over.

## Advanced Configuration

The configuration line in the /etc/ha.d/haresources file as defined above only indicates that the master node is Twin 1 and the cluster IP address is 192.168.252.205. According to this configuration, HeartBeat does the minimum, no script or service is started or stopped on the different nodes composing the cluster.

Let's see how we can configure the haresources file with more finesse for specific needs.

### Assigning Multiple IP Addresses to the Cluster

We wrote a line like this during the basic configuration:

```bash
twin1 IPaddr::192.168.252.205
```

This means that we assign the IP address 192.168.252.205 to the cluster. However, we may need more than one IP address to access the cluster. That is, we have multiple entry points to the cluster and each entry point (IP address) will allow access to the cluster.

To do this, we simply configure the nodes with the following /etc/ha.d/haresources file:

```bash
twin1 IPaddr::192.168.252.205 IPaddr::192.168.252.206
```

And you can assign as many additional IP addresses as desired.

### Specifying Actions to Take

Once the basic configuration is established, it's the master node (if all goes well) that provides services to the outside. If all your nodes are configured exactly the same and run the same services permanently, this doesn't pose any problems.

If we want to be more precise in the node configuration, we also indicate the actions to take in the cluster when a node fails. This aspect is primordial if certain services must access exclusive resources (i.e. only one node at a time can use a particular resource). This is the case with mirrored disks via the network (see RAID-1 over IP) where only one node can access the mirrored disk in read-write mode.

To address this kind of problem, we can specify actions to take by completing the /etc/ha.d/haresources file. To specify these actions, we use scripts in the style of the Init System V (or additional HeartBeat scripts or even homemade scripts) to indicate what the cluster should do in case of a switch from one node to another.

Let's take a simple example (but not an interesting one ;-)): the ssh service. The Init script for this service is located in /etc/init.d/ (along with many others; apache, inetd, rsync,...). We configure the /etc/ha.d/haresources file (on each node, don't forget it!) as follows:

```bash
twin1 IPaddr::192.168.252.205 ssh
```

This line gives HeartBeat the following information:

- The master node is twin1.
- The cluster IP address is 192.168.252.205.
- The actions to take in case of a switch are: ssh.

Let's see in detail what it means by "The actions to take are...".

The scripts used by HeartBeat must respond to 3 commands: start, stop, and status. These commands are called in the following cases:

- start: called on the node to which the cluster has just switched.
- stop: called on the node that has just failed.

So, let's imagine that Twin 1 is responding to requests addressed to the cluster. If Twin 1 falls, HeartBeat executes ssh stop on it (if possible obviously) then HeartBeat switches the IP to Twin 2 and executes ssh start on the latter.

In the ssh case, it doesn't have much interest but for mounting file systems or just sending emails, it's essential.

The search order for scripts is as follows:

- `/etc/ha.d/resource.d/`: which is installed by default and contains a number of interesting scripts.
- `/etc/init.d/`: which is the default directory for the System V Init system.

You can add custom scripts in `/etc/ha.d/resource.d/` but remember that they must be able to respond to the start, stop, and status commands.

Don't hesitate to look at the few scripts in `/etc/ha.d/resource.d/` to draw inspiration from them and know how they work (syntactically speaking).

## Practical Considerations

### Regarding the Network Connection

The solution described above uses the network bandwidth intensively. Every two seconds (or more frequently depending on the configuration), each node sends a UDP packet to each of the other nodes. In addition to that, each node will respond to each of the packets sent. I'll let you imagine the number of packets transiting on the network.

The ideal is to have machines with multiple network cards. For example, we can use two network cards: one network card for taking pulses and one network card for services.

If your cluster consists of only two nodes, then it's better to use a crossover cable between the two machines to avoid overloading the production network.

In any case, I advise you to use a network dedicated to pulse taking and service switching. Otherwise, you will be facing an overload of your production network and performances will slightly drop.

### Regarding an Additional Serial Connection

In addition to the UDP pulse, you can perform a direct pulse via a null modem serial cable that connects the two machines. This pulse is an effective complement to the network in case a network interface would fail.

To use this additional serial connection, you must introduce the following two lines in the `/etc/ha.d/ha.cf` file (just after the bcast eth0 line):

```bash
baud     19200
serial   /dev/ttyS0
```

I won't go further regarding serial connections because I haven't had the opportunity to experiment with them.
During a configuration change...

During a configuration change of the cluster (the 3 HeartBeat configuration files), you must force the heartbeat processes to reread the configuration files. Always make sure to have identical configuration files on the different nodes of your servers, otherwise, clustering would not work.

To force the rereading of the configuration files, you enter the following command on each node of the cluster, finishing with the node that has control (normally, the master node unless you are in a catastrophic situation). Prepare the consoles and commands in advance so that the time during which the nodes have different configurations is as short as possible!:

```bash
sudo /etc/init.d/heartbeat reload
```

The safest solution is to stop all heartbeats, apply the modifications, and restart them; but sometimes, production environments do not allow this kind of thing.

### Multiple Clusters on the Same Network

Recently, I added an additional cluster on a network where there was already a cluster. This posed some problems (nothing serious, just gigantic logs).

In fact, the nodes of the second cluster were receiving broadcast heartbeat packets and didn't know how to interpret them. A warning message in the heartbeat log gave something like this:

```bash
Aug  8 06:28:54 nodeprod1 heartbeat[6743]: ERROR: process_status_message: bad node [nodearch2] in message
Aug  8 06:28:54 nodeprod1 heartbeat[6743]: info: MSG: Dumping message with 9 fields
[...]
```

In fact, the nodes were sending messages in broadcast to each other without knowing what to do with them. To solve this problem, we had to remove the broadcast and use either unicast or multicast.

To do this, you just need to modify the `/etc/ha.d/ha.cf` file and comment out the bcast <interface> line. Then, depending on the number of nodes in your cluster (if you only have two nodes, unicast is sufficient) you modify the configuration.

For a unicast configuration (still following the example above) on Twin1, add the line:

```bash
ucast eth1 192.168.252.201
```

And for Twin2, add the line:

```bash
ucast eth1 192.168.252.200
```

For a multicast configuration (easier to maintain identical files or if there are more than two nodes), add the following line (on both nodes):

```bash
mcast eth1 225.0.0.1 694 1 0
mcast <interface> <multicast group> <udp port> <ttl> 0
```

The multicast group is an address you choose between 224.0.0.0 - 239.255.255.255 through which your nodes will communicate. Be careful not to have other systems using this multicast group. The ttl indicates the number of network hops the packet can make (from 1 to 255). Generally, nodes are on the same network and a value of 1 is perfect.

The zero that ends the line is a bug in heartbeat ;-).

## Resources

For additional high availability and load balancing resources, please refer to these related guides:

- [HAProxy: Load Balance Your Traffic]({{< ref "docs/Servers/HighAvailability/HA-Proxy/haproxy-load-balance-your-traffic.md" >}})
- [Installation and Configuration of a Heartbeat V2 Cluster]({{< ref "docs/Servers/HighAvailability/Heartbeat/installation_and_configuration_of_a_heartbeat_v2_cluster.md" >}})
