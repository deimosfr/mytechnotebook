---
weight: 999
url: "/Installation_et_Configuration_de_Pacemaker/"
title: "Installation and Configuration of Pacemaker"
description: "A guide on how to install and configure Pacemaker with Corosync on Debian 6"
categories: ["Linux", "Debian"]
date: "2010-12-21T20:41:00+02:00"
lastmod: "2010-12-21T20:41:00+02:00"
tags: ["Cluster", "Pacemaker", "Corosync", "High Availability"]
toc: true
---

## Introduction

Pacemaker is the logical evolution of Heartbeat, which has merged with other open source software to achieve perfection. Pacemaker is a resource management software. It must be coupled with Corosync which will manage the exchange of information between nodes.

The installation and configuration described here is done on Debian 6.

## Installation

For installation, nothing could be simpler on Debian:

```bash
aptitude install pacemaker corosync
```

## Configuration

### Corosync

For Corosync configuration, we need to generate encryption keys for intra-cluster information exchange (between nodes):

```bash
> corosync-keygen
Corosync Cluster Engine Authentication key generator.
Gathering 1024 bits for key from /dev/random.
Press keys on your keyboard to generate entropy.
Press keys on your keyboard to generate entropy (bits = 136).
Press keys on your keyboard to generate entropy (bits = 200).
Press keys on your keyboard to generate entropy (bits = 264).
Press keys on your keyboard to generate entropy (bits = 328).
Press keys on your keyboard to generate entropy (bits = 392).
Press keys on your keyboard to generate entropy (bits = 456).
Press keys on your keyboard to generate entropy (bits = 520).
Press keys on your keyboard to generate entropy (bits = 584).
Press keys on your keyboard to generate entropy (bits = 648).
Press keys on your keyboard to generate entropy (bits = 712).
Press keys on your keyboard to generate entropy (bits = 776).
Press keys on your keyboard to generate entropy (bits = 840).
Press keys on your keyboard to generate entropy (bits = 904).
Press keys on your keyboard to generate entropy (bits = 968).
Writing corosync key to /etc/corosync/authkey.
```

Then send this configuration to the other nodes using scp for example:

```bash
scp -r /etc/corosync root@nodex:/etc/
```

Then edit the Corosync configuration file on each node, and adapt the network part for each of them:

(`/etc/corosync/corosync.conf`)

```bash
# Please read the openais.conf.5 manual page

totem {
    version: 2

    # How long before declaring a token lost (ms)
    token: 3000

    # How many token retransmits before forming a new configuration
    token_retransmits_before_loss_const: 10

    # How long to wait for join messages in the membership protocol (ms)
    join: 60

    # How long to wait for consensus to be achieved before starting a new round of membership configuration (ms)
    consensus: 3600

    # Turn off the virtual synchrony filter
    vsftype: none

    # Number of messages that may be sent by one processor on receipt of the token
    max_messages: 20

    # Limit generated nodeids to 31-bits (positive signed integers)
    clear_node_high_bit: yes

    # Disable encryption
    secauth: off

    # How many threads to use for encryption/decryption
    threads: 0

    # Optionally assign a fixed node id (integer)
    # nodeid: 1234

    # This specifies the mode of redundant ring, which may be none, active, or passive.
    rrp_mode: none

    interface {
        # The following values need to be set based on your environment
        ringnumber: 0
        bindnetaddr: 192.168.20.4
        mcastaddr: 226.94.1.1
        mcastport: 5405
    }
}

amf {
    mode: disabled
}

service {
    # Load the Pacemaker Cluster Resource Manager
    ver:       0
    name:      pacemaker
}

aisexec {
        user:   root
        group:  root
}

logging {
        fileline: off
        to_stderr: yes
        to_logfile: no
        to_syslog: yes
    syslog_facility: daemon
        debug: off
        timestamp: on
        logger_subsys {
                subsys: AMF
                debug: off
                tags: enter
```

Finally, edit the `/etc/default/corosync` file to tell it to start at machine boot:

(`/etc/default/corosync`)

```bash
 # start corosync at boot [yes|no]
 START=yes
```

Now let's activate our new configuration on all nodes:

```bash
/etc/init.d/corosync restart
```

Now we can check the status of the cluster:

```bash
> crm_mon --one-shot -V
crm_mon[427]: 2010/12/20_21:51:54 ERROR: unpack_resources: Resource start-up disabled since no STONITH resources have been defined
crm_mon[427]: 2010/12/20_21:51:54 ERROR: unpack_resources: Either configure some or disable STONITH with the stonith-enabled option
crm_mon[427]: 2010/12/20_21:51:54 ERROR: unpack_resources: NOTE: Clusters with shared data need STONITH to ensure data integrity
============
Last updated: Mon Dec 20 21:51:54 2010
Stack: openais
Current DC: zazu - partition WITHOUT quorum
Version: 1.0.9-74392a28b7f31d7ddc86689598bd23114f58978b
2 Nodes configured, 2 expected votes
0 Resources configured.
============

Online: [ zazu shenzi ]
```

Now we will configure a few things:

- No quorum (because it's a 2-node cluster)
- No Stonith (Fencing)
- For stickiness, this is to avoid auto-failback of resources

```bash
> crm
crm(live)# configure
crm(live)configure# property no-quorum-policy=ignore
crm(live)configure# property stonith-enabled=false
crm(live)configure# property default-resource-stickiness=1000
crm(live)configure# commit
crm(live)configure# bye
```
