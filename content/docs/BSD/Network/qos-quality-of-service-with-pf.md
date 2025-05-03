---
weight: 999
url: "/Faire_de_la_QOS_(Quality_Of_Service)_avec_PF/"
title: "QoS (Quality Of Service) with PF"
description: "How to implement Quality of Service (QoS) using Packet Filter (PF) in OpenBSD to prioritize different types of network traffic."
categories: ["Network", "OpenBSD", "Firewall"]
date: "2007-10-07T10:14:00+02:00"
lastmod: "2007-10-07T10:14:00+02:00"
tags: ["pf", "qos", "openbsd", "firewall", "bandwidth", "hfsc"]
toc: true
---

## Introduction

Hierarchical Fair Service Curve (HFSC) alias QOS:
Quality of Service (QoS) is an attempt to give priority to a packet type or data connection on a per session basis. Hierarchical Fair Service Curve takes QoS to the next level over CBQ by focusing on guaranteed real-time, adaptive best-effort, and hierarchical link-sharing service.

Though this may sound difficult, it is really easy to use once you understand the basics.

What HSFC means without technical jargon is, you have the ability to setup rules to govern how data leaves the system. For example...

- You may choose to have ack packets labeled with the highest priority to guarantee those packets go out first. Ack packets are the way you tell the remote system you have received the last payload and to continue to send the next. This will make sure you data transfers go as fast as they can even on a saturated connection.

- What if you are an avid gamer and other users on your network are slowing your connection down or causing you to loose your connection. You choose to give priority to your gaming traffic over normal web traffic. This way you can play games without slowing down and keep your latency low while other users on the network browse the web and download files.

- What if you are running a web server and you find the majority of your data is text based and is less than 10KB per page, but you do have a few larger data files around 5MB. You decide you want to serve out data quickly in the beginning of the connection and slow down after a few seconds. You can setup HFSC to serve out the first few seconds of a connection at full speed, lets say 100KB/sec and then slow the connection down after 5 seconds to 25KB/sec. This allows you to serve out your page at full speed and still allow people to download the 5MB files at slow speed, saving band with for other new web clients.

Quality of Service gives you the tools you need to shape traffic.

## Set of commands

Let's take a look at the basic set of commands in HFSC and why you would use them in the real world:

### bandwidth

The percentage of the total connection speed this queue is allowed to borrow from the total queue or other queues. This variable can be not equal zero(0) and the bandwidth values for all of the queues can not exceed 100% of the available connection. This value is not to be confused with "bandwidth" as a value to describe the amount of data that can be transferred to or from a server, but a as directive specifying the amount of total connection speed this queue can borrow.

### priority

The level specifies the order in which a service is to occur relative to other queues. The higher the number or value, the higher the priority. This directive is a simple way of saying which packets are first out of the gate compared to others. Let's say you have gaming data and bulk web data. You want gaming data to be first since it is interactive and bulk web traffic can wait. Set the gaming data queue at least 1 priority level higher than the bulk web traffic queue.

### qlimit

The amount of "slots" available to a queue to save outgoing packets when the amount of available bandwidth has been exceeded. This value is 50 by default. When the total amount of bandwidth has been reached on the outgoing interface or higher queues are taking up all of the bandwidth then no more data can be sent. The qlimit will put the packets the queue can not send out into slots in memory in the order that they arrive. When bandwidth is available the qlimit slots will be emptied in the order they arrived; first in, first out. If the qlimit reaches the maximum value of qlimit, the packets will be dropped. Look at qlimit slots as "emergency use only," but as a better alternative to dropping the packets out right. Also, do not think that setting the qlimit really high will solve the problem of bandwidth starvation and packet drops. What you want to do is setup a queue with the proper bandwidth boundaries so that packets only go into the qlimit slots for a short time, if ever.

### realtime

The amount of bandwidth that is guaranteed to the queue no matter what any other queue needs. Realtime can be set from 0% to 80% of total connection bandwidth. Let's say you want to make sure that your web server gets 25KB/sec of bandwidth no matter what. Setting the realtime value will give the webserver queue the bandwidth it needs even if other queues want to share its bandwidth.

### upperlimit

The amount of bandwidth the queue can _never_ exceed. For example, say you want to setup a new mail server and you want to make sure that the server never takes up more than 50% of your available bandwidth. Or let's say you have a p2p user you need to limit. Using the upperlimit value will keep them from abusing the connection.

### linkshare

This value has the exact same use as "bandwidth" above. If you decide to use both "bandwidth" and "linkshare" in the same rule pf(OpenBSD 4.1) will just drop "linkshare". For this reason we are not going to use it.

### default

The default queue. As data connections or rules that are specifically put into a queue will be put into this queue rule. This directive must be in one rule. You can _not_ have two(2) default directives in any two(2) rules.

## Setup HFSC

Now, let's take a look at a common HFSC queue setup. The following group of rules splits data into 6 subsets and gives each one of them specific data tasks and limits. You do not have to follow this example exactly, especially since you have the definitions above. Let me explain what each line does and why it is used then you can decide for yourself.

```bash
#Comcast Upload = 768Kb/s (queue at 97%)
 altq on $ExtIf bandwidth 744Kb hfsc queue { ack, dns, ssh, bulk, bittor, spamd }
   queue ack        bandwidth 80% priority 7 qlimit 500 hfsc (realtime 50%)
   queue dns        bandwidth  7% priority 6 qlimit 500 hfsc (realtime  5%)
   queue ssh        bandwidth 10% priority 5 qlimit 500 hfsc (realtime 20%) {ssh_bulk, ssh_login}
    queue ssh_login bandwidth 90% priority 5 qlimit 500 hfsc
    queue ssh_bulk  bandwidth 10% priority 4 qlimit 500 hfsc
   queue bulk       bandwidth  1% priority 4 qlimit 500 hfsc (realtime 5% default)
   queue bittor     bandwidth  1% priority 3 qlimit 500 hfsc (upperlimit 99%)
   queue spamd      bandwidth  1% priority 2 qlimit 500 hfsc (upperlimit 1%)
```

The first line is simply a comment. It reminds one that comcast's total upload bandwidth is 768Kb/s (kilobits per second). You never want to use exactly the total upload speed, but a few kilobytes less. On comcast 97% works very well. Why? Because you want to use your queue as the limiting factor in the connection. When you send out data and you saturate your link the router you connect to will decide what packets go first and that is what we want HSFC to do. You can _not_ trust your upstream router to forward packets correctly. It may be that the router maintainer does not care, the router does not have the ability, or the router it too old. So, we limit the upload speed to just under the total available. "Doesn't that waste some bandwidth then?" Yes, in this example we are not using 3KB/s, but remember we are making sure the upstream routers sends out the packets in the order we want, not what they decide. This makes all the difference with ACK packets especially and will actually increase the available bandwidth on a saturated connection

The second line is the parent queue for the external interface ($ExtIf), it shows we are using "hfsc queue" and lists out all six(6) of the child queues (ack, dns, ssh, bulk, bittor, spamd). This is where we specify the bandwidth limit at 97% of the total 768Kb = 744Kb.

The next set of lines specify the 6 child queues and also two sub-child queues in the ssh rule. All of these rules use the external interface and are limited by the parent queue's bandwidth limitations.

Note: REMEMBER: Do not set your upload bandwidth too high otherwise the queue in pf will be useless. A safe rule is to set the maximum bandwidth at around 97% of the total upload speed available to you. Setting your max speed lower is preferable to setting it too high.

```bash
queue ack bandwidth 80% priority 7 qlimit 500 hfsc (realtime 50%)
```

This is the ack queue. it can borrow as little as 80% of the total bandwidth, it is the highest priority at 7, and has a very high queue limit of 500 slots. The realtime of 50% means this queue is guaranteed at least 50% of the total bandwidth no matter what any other rules want.

The highest priority queue is for ack (acknowledge) packets. Ack packets are the method your system tells the remote systems that you have received the payload they sent and to send the next one. By prioritizing these packets you can keep your transfer rates high even on a highly saturated link. For example, if you are downloading a file and you receive a chunk of data the remote system will not send you the next chunk of data until you send them an OK. The OK is the ack packet. When you send the ack packet the remote system knows you got the packet and it has checked out, thus it will send the next one. If on the other hand you delay ack packets, the transfer rate will diminish quickly because the remote system won't send anything new until you respond.

```bash
queue dns bandwidth 7% priority 6 qlimit 500 hfsc (realtime 5%)
```

This is the dns queue. it can borrow as little as 7% of the total bandwidth, it is the second highest priority at 6, and has a very high queue limit of 500 slots. The realtime of 5% means this queue is guaranteed at least 5% of the total bandwidth no matter what any other rules want.

This queue is simple to make sure dns packets get out on time. Though this is not really necessary your web browsing users will be thankful. When you go to a site or enter a URL the clients need the ip of the server. This rule simply allows dns queries to go out before other traffic.

```bash
queue ssh bandwidth 10% priority 5 qlimit 500 hfsc (realtime 20%) {ssh_login, ssh_bulk}
   queue ssh_login bandwidth 90% priority 5 qlimit 500 hfsc
   queue ssh_bulk bandwidth 10% priority 4 qlimit 500 hfsc
```

This is the ssh parent and child queues. The parent queue can borrow as little as 10% of the total bandwidth, it is at priority at 5, and has a very high queue limit of 500 slots. The realtime of 20% means this queue is guaranteed at least 20% of the total bandwidth.

The two(2) child queues are for ssh's interactive logins (ssh_login) and bulk transfer data like scp/sftp (ssh_bulk). These two queues are under the parent queue and both divide the parent's bandwidth of 10% of the total bandwidth. In this example we want to make sure interactive ssh like the console has at least 90% of the bandwidth. The rest of the bandwidth at 10% is used for bulk transfers like scp and sftp transfers. Both child queues do have the ability to share bandwidth from each other. The priorities of the ssh child queues are independent of all of the other queues. We could have picked any other priorities as long as ssh_login was higher than ssh_bulk.

```bash
queue bulk bandwidth 1% priority 4 qlimit 500 hfsc (realtime 5% default)
```

This is the bulk queue. The bulk queue can borrow as little as 1% of the total bandwidth, it is at priority at 4, and has a very high queue limit of 500 slots. The realtime of 5% means this queue is guaranteed at least 5% of the total bandwidth no matter what any other rules want.

This queue is where all of the general traffic will go. If one does not specify a queue for a rule, that traffic will go here. Notice the directive "default" after the realtime tag.

One also has the option of changing the realtime speed over time. In the following example the bulk queue has been changed. This time we will transfer 37kb/s for 5000 milliseconds and then drop the speeds down to 10kb/sec. This might be useful to keep short bursts fast, but slow down big downloads.

```bash
queue bulk bandwidth 1% priority 4 qlimit 500 hfsc (realtime 37kb 5000 10kb default)
```

```bash
queue bittor bandwidth 1% priority 3 qlimit 500 hfsc (upperlimit 99%)
```

This is bittor queue. The bittor queue can borrow as little as 1% of the total bandwidth, it is at priority at 3, and has a very high queue limit of 500 slots. Notice this rule does not have a real time directive. This is because we have decided that bittor traffic is expendable and we want to make sure this queue gives up all bandwidth to higher priority queues that need it. The upperlimit directive makes sure this rule will never borrow more than 99% of the total bandwidth from any other queue.

This rule is here to show that one can use peer sharing tools and still have control of their network. You will notice that remote clients using peer 2 peer sharing tools connecting will hammer your connection. This rule will allow the data to transfer at up to 99% of your full speed, but if another queue needs the bandwidth, the bittor queue will be pruned almost instantly to 1%. Imagine if you are getting the latest OpenBSD distro through a torrent and then you want to browse the web. Normally, you would experience a slow connection because you are fighting for bandwidth. With this rule your browsing traffic gets the bandwidth it needs instantly. The bittor queue on the other hand gets reduced and starts using the qlimit slots until you are done using the bandwidth browsing. Best of both worlds.

```bash
queue spamd bandwidth 1% priority 2 qlimit 500 hfsc (upperlimit 1%)
```

This is spamd queue. The spamd queue can borrow as little as 1% of the total bandwidth, it is at lowest priority of 2, and has a very high queue limit of 500 slots. Notice this rule does not have a real time directive. This is because we have decided that spamd traffic is expendable and we want to make sure this queue gives up all bandwidth if higher priority queues need it. The upperlimit directive makes sure this rule will never borrow more than 1% of the total bandwidth from any other queue.

This rule is used for spammers and is linked to the spamd daemon to annoy spammers. Since the traffic on this queue has very low traffic requirements we have decided to set the upper and lower bounds at 1% of the total bandwidth. Even with 100 spammers connected less than 1KB/sec is more than enough to annoy them. Even if you had more spammers connected the queue would never use more than 1% of the bandwidth. Any extra packets would go into qlimit and if that fills would then packets would be dropped. No problem since the data is expendable.

Now that we have taken a detailed look at the queue rules and directives, we now need to look at a way to apply those queues to our pf rules.

Here we have two(2) examples of rules you can use queuing on. Notice the queue names we used above like ack, bulk, ssh_login, and ssh_bulk at the end of the rules. Also, notice the order that we have put the two queues in on each rule. The first queue name in "bulk, ack" is for general data and the second "ack" is for special short length packets.

```bash
pass out on $ExtIf inet proto tcp from ($ExtIf) to any flags S/SA modulate state queue (bulk, ack)
pass out on $ExtIf inet proto tcp from ($ExtIf) to any port ssh flags S/SA modulate state queue (ssh_bulk, ssh_login)
```

The first rule is passing out bulk traffic on the external interface and prioritizing ack packets. The second rule is passing out data on port 22(ssh) and prioritizing the interactive ssh traffic. This traffic is originating on our internal network or on the firewall itself.

If we decided to have a rule with only one queue directive it would look like this:

```bash
pass out on $ExtIf inet proto tcp from ($ExtIf) to any flags S/SA modulate state queue (bulk)
```

You can also queue data on the return trip on an external stateful connection. Remember you can _not_ queue data coming into the box, only going out. Let's say you have a web server and clients from the outside connect to you and you want their data responses to be queued. The following works perfectly.

```bash
pass in on $ExtIf inet proto tcp from any to ($ExtIf) port www flags S/SA modulate state queue (bulk, ack)
```

So, now you have read all about queuing and you have applied the queue tags to your rules. Now you need to verify that what you setup works actually does what you thought it should do. You should first install "pftop" from the OpenBSD package collection. It is a very easy install without any dependencies. You can also install the latest version from source without issue.

The following is an example output from "table #8" in pftop. To get to this table start pftop and press #8 on the keyboard.

```bash
pfTop: Up Queue 1-9/9, View: queue, Cache: 10000

  QUEUE          BW SCH  PRIO     PKTS    BYTES   DROP_P   DROP_B QLEN BORROW SUSPEN     P/S     B/S
root_rl0       744K hfsc    0        0        0        0        0    0                     0       0
 ack           595K hfsc    7        0        0        0        0    0                     0       0
 dns          52080 hfsc    6        0        0        0        0    0                     0       0
 ssh          74400 hfsc    5        0        0        0        0    0                     0       0
  ssh_login   66960 hfsc    5       83    13538        0        0    0                   0.2      26
  ssh_bulk     7440 hfsc    4       11     3042        0        0    0                     0       0
 bulk          7440 hfsc    4      406    44540        0        0    0                    80     403
 bittor        7440 hfsc    3        0        0        0        0    0                     0       0
 spamd         7440 hfsc    2    24424  1412491        0        0  140                    15     923
```

The output above is similar to what you are looking for. You need to test each type of queue you setup to make sure you see the packets being added to the correct queue. For example, you could ssh to another machine going out the external interface and as you do so you should see packets being add to the "ssh_login" queue.
