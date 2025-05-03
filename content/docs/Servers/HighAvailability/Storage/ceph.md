---
weight: 999
url: "/Ceph_\\:_performance,_reliability_and_scalability_storage_solution/"
title: "Ceph: Performance, Reliability and Scalability Storage Solution"
description: "Learn how to implement Ceph, an open-source distributed storage system that provides object, block, and file storage in a single platform for improved reliability and scalability."
categories: ["Storage", "Linux"]
date: "2014-06-02T15:06:00+02:00"
lastmod: "2014-06-02T15:06:00+02:00"
tags:
  [
    "ceph",
    "storage",
    "distributed",
    "cluster",
    "object-storage",
    "block-storage",
  ]
toc: true
---

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 0.72.2 |
| **Operating System** | Debian 7 |
| **Website** | [Ceph Website](https://ceph.com/) |
| **Last Update** | 02/06/2014 |
{{< /table >}}

![Ceph](/images/ceph_logo.avif)

## Introduction

Ceph is an open-source, massively scalable, software-defined storage system which provides object, block and file system storage in a single platform. It runs on commodity hardware-saving you costs, giving you flexibility and because it's in the Linux kernel, it's easy to consume.

Ceph is able to manage:

- **Object Storage**: Ceph provides seamless access to objects using native language bindings or radosgw, a REST interface that's compatible with applications written for S3 and Swift.
- **Block Storage**: Ceph's RADOS Block Device (RBD) provides access to block device images that are striped and replicated across the entire storage cluster.
- **File System**: Ceph provides a POSIX-compliant network file system that aims for high performance, large data storage, and maximum compatibility with legacy applications (not yet stable)

Whether you want to provide Ceph Object Storage and/or Ceph Block Device services to Cloud Platforms, deploy a Ceph Filesystem or use Ceph for another purpose, all Ceph Storage Cluster deployments begin with setting up each Ceph Node, your network and the Ceph Storage Cluster. **A Ceph Storage Cluster requires at least one Ceph Monitor and at least two Ceph OSD Daemons. The Ceph Metadata Server is essential when running Ceph Filesystem clients.**

- **OSDs**: A Ceph OSD Daemon (OSD) stores data, handles data replication, recovery, backfilling, rebalancing, and provides some monitoring information to Ceph Monitors by checking other Ceph OSD Daemons for a heartbeat. A Ceph Storage Cluster requires at least two Ceph OSD Daemons to achieve an _active + clean_ state when the cluster makes two copies of your data (Ceph makes 2 copies by default, but you can adjust it).
- **Monitors**: A Ceph Monitor maintains maps of the cluster state, including the monitor map, the OSD map, the Placement Group (PG) map, and the CRUSH map. Ceph maintains a history (called an "epoch") of each state change in the Ceph Monitors, Ceph OSD Daemons, and PGs.
- **MDSs**: A Ceph Metadata Server (MDS) stores metadata on behalf of the Ceph Filesystem (i.e., Ceph Block Devices and Ceph Object Storage do not use MDS). Ceph Metadata Servers make it feasible for POSIX file system users to execute basic commands like _ls, find, etc_. without placing an enormous burden on the Ceph Storage Cluster.

Ceph stores a client's data as objects within storage pools. Using the CRUSH algorithm, Ceph calculates which placement group should contain the object, and further calculates which Ceph OSD Daemon should store the placement group. The CRUSH algorithm enables the Ceph Storage Cluster to scale, rebalance, and recover dynamically.

### Testing case

If you want to test with [Vagrant](./vagrant_:_quickly_deploy_virtual_machines.html) and VirtualBox, I've made a Vagrantfile for it running on Debian Wheezy:

```ruby {linenos=table}
# -*- mode: ruby -*-
# vi: set ft=ruby :
ENV['LANG'] = 'C'

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

# Insert all your Vms with configs
boxes = [
    { :name => :mon1, :role => 'mon'},
    { :name => :mon2, :role => 'mon'},
    { :name => :mon3, :role => 'mon'},
    { :name => :osd1, :role => 'osd', :ip => '192.168.33.31'},
    { :name => :osd2, :role => 'osd', :ip => '192.168.33.32'},
    { :name => :osd3, :role => 'osd', :ip => '192.168.33.33'},
]

$install = <<INSTALL
wget -q -O- 'https://ceph.com/git/?p=ceph.git;a=blob_plain;f=keys/release.asc' | sudo apt-key add -
echo deb http://ceph.com/debian/ $(lsb_release -sc) main | sudo tee /etc/apt/sources.list.d/ceph.list
aptitude update
aptitude -y install ceph ceph-deploy openntpd
INSTALL

Vagrant::Config.run do |config|
  # Default box OS
  vm_default = proc do |boxcnf|
    boxcnf.vm.box       = "deimosfr/debian-wheezy"
  end

  # For each VM, add a public and private card. Then install Ceph
  boxes.each do |opts|
    vm_default.call(config)
    config.vm.define opts[:name] do |config|
        config.vm.network   :bridged, :bridge => "eth0"
        config.vm.host_name = "%s.vm" % opts[:name].to_s
        config.vm.provision "shell", inline: $install
        # Create 8G disk file and add private interface for OSD VMs
        if opts[:role] == 'osd'
            config.vm.network   :hostonly, opts[:ip]
            file_to_disk = 'osd-disk_' + opts[:name].to_s + '.vdi'
            config.vm.customize ['createhd', '--filename', file_to_disk, '--size', 8 * 1024]
            config.vm.customize ['storageattach', :id, '--storagectl', 'SATA', '--port', 1, '--device', 0, '--type', 'hdd', '--medium', file_to_disk]
        end
    end
  end
end
```

This will spawn VMs with correct hardware to run. It will also install Ceph as well. After booting those instances, you will have all the Ceph servers like that:

![Ceph diagram](/images/ceph_diagram.avif)

## Installation

### First node

To get the latest version, we're going to use the official repositories:

```bash
wget -q -O- 'https://ceph.com/git/?p=ceph.git;a=blob_plain;f=keys/release.asc' | sudo apt-key add -
echo deb http://ceph.com/debian/ $(lsb_release -sc) main | sudo tee /etc/apt/sources.list.d/ceph.list
```

Then we're going to install Ceph and ceph-deploy which help us to install in a faster way all the components:

```bash
aptitude update
aptitude install ceph ceph-deploy openntpd
```

Openntpd is not mandatory, but you need all your machine to be clock synchronized!

{{< alert context="warning" text="You absolutely need well named servers (hostname) and dns names available. A dns server is mandatory." />}}

### Other nodes

To install ceph from the first (admin) node for any kind of nodes, here is a simple solution that avoid to enter the apt key etc...:

```bash
ceph-deploy install <node_name>
```

{{< alert context="warning" text="For production usage, you should choose an LTS version like emperor" />}}

```bash
ceph-deploy install --release emperor <node_name>
```

{{< alert context="info" text="You need to exchange SSH keys to remotely be able to connect to the target machines." />}}

## Deploy

First of all, your Ceph configuration will be generated in the current directory, so it's suggested to create a dedicated folder for it:

```bash
mkdir ceph-config
cd ceph-config
```

{{< alert context="warning" text="Be sure of your network configuration for Monitor nodes as it's a nightmare to change later!!!" />}}

### Cluster

The first machine on which you'll want to start, you'll need to create the Ceph cluster. I'll do it here on the monitor server named mon1:

```bash
> ceph-deploy new mon1
[ceph_deploy.cli][INFO  ] Invoked (1.2.7): /usr/bin/ceph-deploy new mon1
[ceph_deploy.new][DEBUG ] Creating new cluster named ceph
[ceph_deploy.new][DEBUG ] Resolving host mon1
[ceph_deploy.new][DEBUG ] Monitor mon1 at 192.168.33.30
[ceph_deploy.new][DEBUG ] Monitor initial members are ['mon1']
[ceph_deploy.new][DEBUG ] Monitor addrs are ['192.168.33.30']
[ceph_deploy.new][DEBUG ] Creating a random mon key...
[ceph_deploy.new][DEBUG ] Writing initial config to ceph.conf...
[ceph_deploy.new][DEBUG ] Writing monitor keyring to ceph.mon.keyring...
```

### Monitor

A Ceph Monitor maintains maps of the cluster state, including the monitor map, the OSD map, the Placement Group (PG) map, and the CRUSH map. Ceph maintains a history (called an "epoch") of each state change in the Ceph Monitors, Ceph OSD Daemons, and PGs.

#### Add the first monitor

To create **the first** (only the first today because ceph-deploy got problems) monitor node (mon1 here):

```bash
> ceph-deploy --overwrite-conf mon create mon1
[ceph_deploy.cli][INFO  ] Invoked (1.2.7): /usr/bin/ceph-deploy --overwrite-conf mon create mon1
[ceph_deploy.mon][DEBUG ] Deploying mon, cluster ceph hosts mon1
[ceph_deploy.mon][DEBUG ] detecting platform for host mon1 ...
[ceph_deploy.sudo_pushy][DEBUG ] will use a local connection without sudo
[ceph_deploy.mon][INFO  ] distro info: Debian 7.2 wheezy
[mon1][DEBUG ] determining if provided host has same hostname in remote
[mon1][DEBUG ] deploying mon to mon1
[mon1][DEBUG ] remote hostname: mon1
[mon1][INFO  ] write cluster configuration to /etc/ceph/{cluster}.conf
[mon1][INFO  ] creating path: /var/lib/ceph/mon/ceph-mon1
[mon1][DEBUG ] checking for done path: /var/lib/ceph/mon/ceph-mon1/done
[mon1][DEBUG ] done path does not exist: /var/lib/ceph/mon/ceph-mon1/done
[mon1][INFO  ] creating keyring file: /var/lib/ceph/tmp/ceph-mon1.mon.keyring
[mon1][INFO  ] create the monitor keyring file
[mon1][INFO  ] Running command: ceph-mon --cluster ceph --mkfs -i mon1 --keyring /var/lib/ceph/tmp/ceph-mon1.mon.keyring
[mon1][INFO  ] ceph-mon: mon.noname-a 192.168.33.30:6789/0 is local, renaming to mon.mon1
[mon1][INFO  ] ceph-mon: set fsid to df462719-9477-4f22-afdf-8f237a576cad
[mon1][INFO  ] ceph-mon: created monfs at /var/lib/ceph/mon/ceph-mon1 for mon.mon1
[mon1][INFO  ] unlinking keyring file /var/lib/ceph/tmp/ceph-mon1.mon.keyring
[mon1][INFO  ] create a done file to avoid re-doing the mon deployment
[mon1][INFO  ] create the init path if it does not exist
[mon1][INFO  ] locating `service` executable...
[mon1][INFO  ] found `service` executable: /usr/sbin/service
[mon1][INFO  ] Running command: /usr/sbin/service ceph -c /etc/ceph/ceph.conf start mon.mon1
[mon1][DEBUG ] === mon.mon1 ===
[mon1][DEBUG ] Starting Ceph mon.mon1 on mon1...
[mon1][DEBUG ] Starting ceph-create-keys on mon1...
[mon1][INFO  ] Running command: ceph --admin-daemon /var/run/ceph/ceph-mon.mon1.asok mon_status
[mon1][DEBUG ] ********************************************************************************
[mon1][DEBUG ] status for monitor: mon.mon1
[mon1][DEBUG ] {
[mon1][DEBUG ]   "election_epoch": 2,
[mon1][DEBUG ]   "extra_probe_peers": [],
[mon1][DEBUG ]   "monmap": {
[mon1][DEBUG ]     "created": "0.000000",
[mon1][DEBUG ]     "epoch": 1,
[mon1][DEBUG ]     "fsid": "df462719-9477-4f22-afdf-8f237a576cad",
[mon1][DEBUG ]     "modified": "0.000000",
[mon1][DEBUG ]     "mons": [
[mon1][DEBUG ]       {
[mon1][DEBUG ]         "addr": "192.168.33.30:6789/0",
[mon1][DEBUG ]         "name": "mon1",
[mon1][DEBUG ]         "rank": 0
[mon1][DEBUG ]       }
[mon1][DEBUG ]     ]
[mon1][DEBUG ]   },
[mon1][DEBUG ]   "name": "mon1",
[mon1][DEBUG ]   "outside_quorum": [],
[mon1][DEBUG ]   "quorum": [
[mon1][DEBUG ]     0
[mon1][DEBUG ]   ],
[mon1][DEBUG ]   "rank": 0,
[mon1][DEBUG ]   "state": "leader",
[mon1][DEBUG ]   "sync_provider": []
[mon1][DEBUG ] }
[mon1][DEBUG ] ********************************************************************************
[mon1][INFO  ] monitor: mon.mon1 is running
[mon1][INFO  ] Running command: ceph --admin-daemon /var/run/ceph/ceph-mon.mon1.asok mon_status
```

#### Add a monitor

To add 2 others monitors nodes (mon2 and mon3) in the cluster, you'll need to edit the configuration on a monitor and admin node. You'll have to set the mon_host, mon_initial_members and public_network configuration in:

```bash {linenos=table,hl_lines=[3,4,7]}
[global]
fsid = 0a91b62e-cd43-4558-9ebd-4719f830cf8b
mon_initial_members = mon1,mon2,mon3
mon_host = 192.168.33.30,192.168.33.31,192.168.33.32
auth_supported = cephx
osd_journal_size = 1024
filestore_xattr_use_omap = true
public_network = 192.168.0.0/24
```

Then update the current cluster configuration:

```bash
ceph-deploy --overwrite-conf config push mon2
ceph-deploy --overwrite-conf config push mon3
```

or

```bash
ceph-deploy --overwrite-conf config push mon2 mon3
```

You now need to update the new configuration on all your monitor nodes:

```bash
ceph-deploy --overwrite-conf config push mon1 mon2 mon3
```

{{< alert context="warning" text="For production usage, you need at least 3 nodes" />}}

#### Remove a monitor

If you need to remove a monitor for maintenance:

```bash
ceph-deploy mon destroy mon1
```

### Admin node

To get the first admin node, you'll need to gather keys on a monitor node. To make it simple, all ceph-deploy should be done from that machine:

```bash
> ceph-deploy gatherkeys mon1
[ceph_deploy.cli][INFO  ] Invoked (1.2.7): /usr/bin/ceph-deploy gatherkeys osd1
[ceph_deploy.gatherkeys][DEBUG ] Checking osd1 for /etc/ceph/ceph.client.admin.keyring
[ceph_deploy.sudo_pushy][DEBUG ] will use a local connection without sudo
[ceph_deploy.gatherkeys][DEBUG ] Got ceph.client.admin.keyring key from osd1.
[ceph_deploy.gatherkeys][DEBUG ] Have ceph.mon.keyring
[ceph_deploy.gatherkeys][DEBUG ] Checking osd1 for /var/lib/ceph/bootstrap-osd/ceph.keyring
[ceph_deploy.sudo_pushy][DEBUG ] will use a local connection without sudo
[ceph_deploy.gatherkeys][DEBUG ] Got ceph.bootstrap-osd.keyring key from osd1.
[ceph_deploy.gatherkeys][DEBUG ] Checking osd1 for /var/lib/ceph/bootstrap-mds/ceph.keyring
[ceph_deploy.sudo_pushy][DEBUG ] will use a local connection without sudo
[ceph_deploy.gatherkeys][DEBUG ] Got ceph.bootstrap-mds.keyring key from osd1.
```

Then you need to exchange SSH keys to remotely be able to connect to the target machines.

### OSD

Ceph OSD Daemon (OSD) stores data, handles data replication, recovery, backfilling, rebalancing, and provides some monitoring information to Ceph Monitors by checking other Ceph OSD Daemons for a heartbeat. A Ceph Storage Cluster requires at least two Ceph OSD Daemons to achieve an active + clean state when the cluster makes two copies of your data (Ceph makes 2 copies by default, but you can adjust it).

#### Add an OSD

To deploy Ceph OSD, we'll first start to erase the remote disk and create a gpt table on the dedicated disk 'sdb':

```bash
> ceph-deploy disk zap osd1:sdb
[ceph_deploy.cli][INFO  ] Invoked (1.3.2): /usr/bin/ceph-deploy disk zap osd1:sdb
[ceph_deploy.osd][DEBUG ] zapping /dev/sdb on osd1
[osd1][DEBUG ] connected to host: osd1
[osd1][DEBUG ] detect platform information from remote host
[osd1][DEBUG ] detect machine type
[ceph_deploy.osd][INFO  ] Distro info: debian 7.2 wheezy
[osd1][DEBUG ] zeroing last few blocks of device
[osd1][INFO  ] Running command: sgdisk --zap-all --clear --mbrtogpt -- /dev/sdb
[osd1][DEBUG ] Warning: The kernel is still using the old partition table.
[osd1][DEBUG ] The new table will be used at the next reboot.
[osd1][DEBUG ] GPT data structures destroyed! You may now partition the disk using fdisk or
[osd1][DEBUG ] other utilities.
[osd1][DEBUG ] Warning: The kernel is still using the old partition table.
[osd1][DEBUG ] The new table will be used at the next reboot.
[osd1][DEBUG ] The operation has completed successfully.
```

It will create a journalized partition and a data one. Then we could create partitions on the on 'osd1' server and prepare + activate this OSD:

```bash
> ceph-deploy --overwrite-conf osd create osd1:sdb
[ceph_deploy.cli][INFO  ] Invoked (1.3.2): /usr/bin/ceph-deploy --overwrite-conf mon create osd1:sdb
[ceph_deploy.mon][DEBUG ] Deploying mon, cluster ceph hosts osd1:sdb
[ceph_deploy.mon][DEBUG ] detecting platform for host osd1 ...
ssh: Could not resolve hostname sdb: Name or service not known
[ceph_deploy.mon][ERROR ] connecting to host: sdb resulted in errors: HostNotFound sdb
[ceph_deploy][ERROR ] GenericError: Failed to create 1 monitors
```

You can see there's an error but it works:

```bash {linenos=table,hl_lines=[5]}
> ceph -s
    cluster 0a91b62e-cd43-4558-9ebd-4719f830cf8b
     health HEALTH_WARN 192 pgs degraded; 192 pgs stuck unclean
     monmap e1: 1 mons at {mon1=192.168.33.30:6789/0}, election epoch 2, quorum 0 mon1
     osdmap e5: 1 osds: 1 up, 1 in
      pgmap v8: 192 pgs, 3 pools, 0 bytes data, 0 objects
            34912 kB used, 7122 MB / 7156 MB avail
                 192 active+degraded
```

#### Get OSD status

To know the OSD status:

```bash
> ceph osd tree
# id	weight	type name	up/down	reweight
-1	0.02998	root default
-2	0.009995		host osd1
3	0.009995			osd.3	up	1
-3	0.009995		host osd3
1	0.009995			osd.1	up	1
-4	0.009995		host osd2
2	0.009995			osd.2	up	1
```

#### Remove an OSD

To remove an OSD, it's unfortunately not yet integrated in ceph-deploy. So first, look at the current status:

```bash {linenos=table,hl_lines=[5]}
> ceph osd tree
# id	weight	type name	up/down	reweight
-1	0.03998	root default
-2	0.01999		host osd1
0	0.009995			osd.0	down	0
3	0.009995			osd.3	up	1
-3	0.009995		host osd3
1	0.009995			osd.1	up	1
-4	0.009995		host osd2
2	0.009995			osd.2	up	1
```

Here I want to remove the osd.0:

```bash
osd='osd.0'
```

If the OSD wasn't down, I should put it down with this command:

```bash
> ceph osd out $osd
osd.0 is already out.
```

Then I remove it from the crushmap:

```bash
> ceph osd crush remove $osd
removed item id 0 name 'osd.0' from crush map
```

Delete the authentication part (Paxos):

```bash
ceph auth del $osd
```

Then remove the OSD:

```bash
ceph osd rm $osd
```

And now, it's definitively out:

```bash
> ceph osd tree
# id	weight	type name	up/down	reweight
-1	0.02998	root default
-2	0.009995		host osd1
3	0.009995			osd.3	up	1
-3	0.009995		host osd3
1	0.009995			osd.1	up	1
-4	0.009995		host osd2
2	0.009995			osd.2	up	1
```

### RBD

To make Block devices, you need to have a correct OSD configuration done with a created pool. You don't have anything else to have :-)

## Configuration

### OSD

#### OSD Configuration

##### Global OSD configuration

The Ceph Client retrieves the latest cluster map and the CRUSH algorithm calculates how to map the object to a placement group, and then calculates how to assign the placement group to a Ceph OSD Daemon dynamically. By default Ceph have 2 replicas and you can change it by 3 in adding those line to the Ceph configuration:

```bash
[global]
    osd pool default size = 3
    osd pool default min size = 1
```

- osd pool default size: the number of replicas
- osd pool default min size: set the minimum available replicas before putting OSD down

Configure the placement group (Total PGs = (number of OSD \* 100) / replicas numbers):

```bash
[global]
    osd pool default pg num = 100
    osd pool default pgp num = 100
```

Then you can push your new configuration:

```bash
ceph-deploy --overwrite-conf config push mon1 mon2 mon3
```

##### Network configuration

For the OSD, you've got 2 network interfaces (private and public). So to configure it properly on your admin machine by updating your configuration file as follow:

```bash
[osd]
cluster network = 192.168.33.0/24
public network = 192.168.0.0/24
```

But if you want to add specific configuration:

```bash
[osd.0]
public addr = 192.168.0.1:6801
cluster addr = 192.168.33.31

[osd.1]
public addr = 192.168.0.2:6802
cluster addr = 192.168.33.32

[osd.2]
public addr = 192.168.0.3:6803
cluster addr = 192.168.33.33
```

Then you can push your new configuration:

```bash
ceph-deploy --overwrite-conf config push mon1 mon2 mon3
```

#### Create an OSD pool

To create an OSD pool:

```bash
ceph osd pool create <pool_name> <pg_num> <pgp_num>
```

#### List OSD pools

You can list OSD:

```bash
> ceph osd lspools
0 data,1 metadata,2 rbd,3 mypool,
```

```bash
> ceph osd dump
epoch 62
fsid 0314c737-68d2-4d14-a247-53dfe7ec2a01
created 2014-01-02 16:33:14.572233
modified 2014-01-03 15:40:14.370846
flags

pool 0 'data' rep size 2 min_size 1 crush_ruleset 0 object_hash rjenkins pg_num 64 pgp_num 64 last_change 1 owner 0 crash_replay_interval 45
pool 1 'metadata' rep size 2 min_size 1 crush_ruleset 1 object_hash rjenkins pg_num 64 pgp_num 64 last_change 1 owner 0
pool 2 'rbd' rep size 2 min_size 1 crush_ruleset 2 object_hash rjenkins pg_num 64 pgp_num 64 last_change 1 owner 0
pool 3 'mypool' rep size 2 min_size 1 crush_ruleset 0 object_hash rjenkins pg_num 100 pgp_num 100 last_change 14 owner 0

max_osd 3
osd.0 up   in  weight 1 up_from 61 up_thru 61 down_at 59 last_clean_interval [49,60) 192.168.32.15:6800/5577 192.168.33.31:6801/5577 192.168.33.31:6802/5577 192.168.32.15:6801/5577 exists,up d1b21569-bcf2-45a6-87c3-597cd267bdba
osd.1 up   in  weight 1 up_from 53 up_thru 61 down_at 51 last_clean_interval [33,50) 192.168.32.72:6800/4439 192.168.33.32:6800/4439 192.168.33.32:6801/4439 192.168.32.72:6802/4439 exists,up 935e8e8f-3a0b-445b-b730-5de61ea34556
osd.2 up   in  weight 1 up_from 57 up_thru 61 down_at 55 last_clean_interval [40,54) 192.168.32.90:6800/4444 192.168.33.33:6800/4444 192.168.33.33:6801/4444 192.168.32.90:6802/4444 exists,up 4e661a3c-2f15-4037-9fcb-d48a49d0b228
```

### Crush map

The Crushmap contain a list of OSDs, a list of 'buckets' for aggregating the devices into physical locations. You can edit it to manage it manually.

To show the complete map of your Ceph:

```bash
ceph osd tree
```

To get crushmap and to edit it:

```bash
ceph osd getcrushmap -o mymap
crushtool -d mymap -o mymap.txt
```

To set a new crushmap after editing:

```bash
crushtool -c mymap.txt -o mynewmap
ceph osd setcrushmap -i mynewmap
```

To get quorum status for the monitors:

```bash
ceph quorum_status --format json-pretty
```

## Usage

### Check health

You can check cluster health:

```bash
> ceph -s
  cluster 0314c737-68d2-4d14-a247-53dfe7ec2a01
   health HEALTH_OK
   monmap e3: 3 mons at {mon1=192.168.33.31:6789/0,mon2=192.168.33.32:6789/0,mon3=192.168.33.33:6789/0}, election epoch 8, quorum 0,1,2 mon1,mon2,mon3
   osdmap e13: 3 osds: 3 up, 3 in
    pgmap v21: 192 pgs: 192 active+clean; 0 bytes data, 71248 KB used, 14244 MB / 14313 MB avail
   mdsmap e1: 0/0/1 up
```

### Change configuration on the fly

To avoid service restart on a simple modification, you can interact directly with Ceph to change some values. First of all, you can get all current values of your Ceph cluster:

```bash
> ceph --admin-daemon /var/run/ceph/ceph-mon.mon1-pub.asok config show
{ "name": "mon.mon1-pub",
  "cluster": "ceph",
  "none": "0\/5",
  "lockdep": "0\/1",
  "context": "0\/1",
[...]
  "mon_clock_drift_allowed": "0.05",
  "mon_clock_drift_warn_backoff": "5",
[...]
```

Now if I want to change one of those values :

```bash
> ceph tell osd.* injectargs '--mon_clock_drift_allowed 1'
osd.0: mon_clock_drift_allowed = '1'
osd.1: mon_clock_drift_allowed = '1'
osd.2: mon_clock_drift_allowed = '1'
```

You can change '\*' by the name of an OSD if you want to apply this configuration to a specific node.

Then do not forget to add it to your Ceph configuration and push it !

### Use object storage

#### Add object

When you have an Ceph OSD pool ready, you can add a file :

```bash
rados put <object_name> <file_to_upload> --pool=<pool_name>
```

#### List objects in a pool

You can list your pool content :

```bash
> rados -p <pool_name> ls
filename
```

You can see where it has been stored :

```bash
> ceph osd map <pool_name> <filename>
osdmap e62 pool '<pool_name>' (3) object 'filename' -> pg 3.af0f2847 (3.47) -> up [1,0] acting [1,0]
```

To locate the file on the hard drive, look at this folder (/var/lib/ceph/osd/ceph-1/current). Then look at the previous result (3.47) and the filename af0f2847. So the file will be placed here :

```bash
> ls /var/lib/ceph/osd/ceph-1/current/3.47_head
ceph\ulog__head_AF0F2847__3
```

#### Remove objects

To finish, remove it :

```bash
rados rm <filename> --pool=<pool_name>
```

### Use blocks device storage

This part becomes very interesting if you start using block devices storage. From an admin node, launch client install :

```bash
ceph-deploy install <client>
ceph-deploy admin <client>
```

On the client, load the module and add it to launch at boot :

```bash
modprobe rbd
echo "rbd" >> /etc/modules
```

#### Create a block device

Now to create a block device (you can do it on the client node if you want has it gets the admin key last pushed by ceph-deploy) :

```bash
rbd create <name> --size <size_in_megabytes> [--pool <pool_name>]
```

If you don't specify the pool name option, it will automatically be created in the 'rbd' pool.

#### List available block devices

This is so simple :

```bash
> rbd ls
bd_name
```

And you can show mapped :

```bash
> rbd showmapped
id pool image snap device
0  rbd  bd_name   -    /dev/rbd0
```

#### Get block device informations

You may need to grab informations on the device block to know where it is physically or simply the size :

```bash
> rbd info <name>
rbd image 'bd_name':
	size 4096 MB in 1024 objects
	order 22 (4096 kB objects)
	block_name_prefix: rb.0.139d.74b0dc51
	format: 1
```

#### Map and mount a block device

First your need to map it to make it appears in your device list :

```bash
rbd map <name> [--pool <pool_name>]
```

You can now format it :

```bash
mkfs.ext4 /dev/rbd/pool_name/name
```

And now mount it :

```bash
mount /dev/rbd/pool_name/name /mnt
```

#### Remove a block device

To remove, once again it's simple :

```bash
rbd rm <name> [--pool <pool_name>]
```

#### Umount and unmap a block device

Umount and umap is as easy as you think :

```bash
umount /mnt
rbd unmap /dev/rbd/pool_name/name
```

#### Advanced usage

I won't list all features, but you can look at the man to :

- clone
- export/import
- snapshot
- bench write

## FAQ

### Reset a node

You sometime needs to reset a node. It's generally needed when you're doing tests. From the admin node run this :

```bash
ceph-deploy purge <node_name>
ceph-deploy purgedata <node_name>
```

Then reinstall it with :

```bash
ceph-deploy install <node_name>
```

### Can't add a new monitor node

If you can't add a new monitor mon (here mon2)[5] :

```bash
> ceph-deploy --overwrite-conf mon create mon2
[...]
[mon2][DEBUG ] === mon.mon2 ===
[mon2][DEBUG ] Starting Ceph mon.mon2 on mon2...
[mon2][DEBUG ] failed: 'ulimit -n 32768;  /usr/bin/ceph-mon -i mon2 --pid-file /var/run/ceph/mon.mon2.pid -c /etc/ceph/ceph.conf '
[mon2][DEBUG ] Starting ceph-create-keys on mon2...
[mon2][WARNIN] No data was received after 7 seconds, disconnecting...
[mon2][INFO  ] Running command: ceph --cluster=ceph --admin-daemon /var/run/ceph/ceph-mon.mon2.asok mon_status
[mon2][ERROR ] admin_socket: exception getting command descriptions: [Errno 2] No such file or directory
[mon2][WARNIN] monitor: mon.mon2, might not be running yet
[mon2][INFO  ] Running command: ceph --cluster=ceph --admin-daemon /var/run/ceph/ceph-mon.mon2.asok mon_status
[mon2][ERROR ] admin_socket: exception getting command descriptions: [Errno 2] No such file or directory
[mon2][WARNIN] mon2 is not defined in `mon initial members`
[mon2][WARNIN] monitor mon2 does not exist in monmap
[mon2][WARNIN] neither `public_addr` nor `public_network` keys are defined for monitors
[mon2][WARNIN] monitors may not be able to form quorum
```

You have to add public network and monitor to the list in configuration file. Look here to see how to add correctly a new mon.

### health HEALTH_WARN clock skew detected

If you get this kind of problem :

```bash
> ceph -s
  cluster 0314c737-68d2-4d14-a247-53dfe7ec2a01
   health HEALTH_WARN clock skew detected on mon.mon1
   monmap e3: 3 mons at {mon1=192.168.33.31:6789/0,mon2=192.168.33.32:6789/0,mon3=192.168.33.33:6789/0}, election epoch 8, quorum 0,1,2 mon1,mon2,mon3
   osdmap e13: 3 osds: 3 up, 3 in
    pgmap v21: 192 pgs: 192 active+clean; 0 bytes data, 71248 KB used, 14244 MB / 14313 MB avail
   mdsmap e1: 0/0/1 up
```

You need to install an NTP server. The millisecond is important. Here is a workaround (`ceph.conf`):

```ini
[global]
    mon_clock_drift_allowed = 1
```

This set the possible delta to avoid the warning message.

## References

[^1]: http://www.inktank.com/what-is-ceph/
[^2]: http://ceph.com/docs/master/rados/configuration/pool-pg-config-ref/
[^3]: http://ceph.com/docs/master/rados/operations/crush-map/
[^4]: http://ceph.com/docs/master/man/8/crushtool/
[^5]: http://tracker.ceph.com/issues/5195
