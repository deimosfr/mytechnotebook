---
weight: 999
url: "/Belier_\\:_script_your_SSH_connection/"
title: "Belier: Script Your SSH Connection"
description: "How to use Belier to simplify complex SSH connections through multiple intermediate servers"
categories: ["Linux", "SSH", "Networking"]
date: "2013-05-06T09:22:00+02:00"
lastmod: "2013-05-06T09:22:00+02:00"
tags: ["ssh", "tunneling", "debian", "expect", "automation"]
toc: true
---

![Belier](/images/belier_logo.avif)

{{< table "table-hover table-striped" >}}
| | |
|-|-|
| **Software version** | 1.2 |
| **Operating System** | Debian 7 |
| **Website** | [Belier Website](https://www.ohmytux.com/belier/) |
| **Last Update** | 06/05/2013 |
{{< /table >}}

## Introduction

Belier[^1] allows opening a shell or executing a command on a remote computer through an SSH connection. The main feature of Belier is its ability to cross several intermediate computers before completing the job.

* Belier reaches the final computer through intermediate machines.
* You can execute commands with any account available on the remote computer.
* It is possible to switch accounts on intermediate computers before accessing the final computer.
* You can open a data tunnel through every host you cross to the final host.
* Belier generates one script for each final computer to reach.

Belier aims to give a single system administrator a tool to work independently, without requiring modification of the computers they need to cross, just using the current configurations of every machine they have to work with. So it's not a revolutionary tool, but it helps to be productive quickly.

For instance, here is what is possible to do with it:

![Schema belier](/images/schema_belier.avif)[^2]

## Installation

To install it, as usual it's easy on Debian:

```bash
aptitude install belier
```

## Configuration

The configuration is really simple. Type line by line in a 'connection' file, connection information for:

1. The server name or IP
2. The username
3. The password

It should look like this:

```bash
server1_bound username password
server2_bound username password
server3_bound username password
final_server username password
```

Then generate an automatic shell connection script with Belier:

```bash
bel -e ~/connection
```

You should now have a script containing all the commands to quickly connect to your desired server:

```bash
#!/usr/bin/expect -f
set timeout 10

spawn ssh -o NoHostAuthenticationForLocalhost=yes -o StrictHostKeyChecking=no  server1_bound
expect -re  "(%|#|\$) $"
send -- "su - username\r"
expect ":"
send -- "password\r"
expect -re  "(%|#|\$) $"
send -- "ssh -o NoHostAuthenticationForLocalhost=yes -o StrictHostKeyChecking=no  server2_bound\r"
expect -re  "(%|#|\$) $"
send -- "su - username\r"
expect ":"
send -- "password\r"
expect -re  "(%|#|\$) $"
send -- "ssh -o NoHostAuthenticationForLocalhost=yes -o StrictHostKeyChecking=no  server3_bound\r"
expect -re  "(%|#|\$) $"
send -- "su - username\r"
expect ":"
send -- "password\r"
expect -re  "(%|#|\$) $"
send -- "ssh -o NoHostAuthenticationForLocalhost=yes -o StrictHostKeyChecking=no  final_server\r"
expect -re  "(%|#|\$) $"
send -- "su - username\r"
expect ":"
send -- "password\r"
expect -re  "(%|#|\$) $"
interact +++ return
```

You can now launch the command with final_server.sh :-)

[^1]: http://www.ohmytux.com/belier/
[^2]: http://carlchenet.wordpress.com/?p=22
