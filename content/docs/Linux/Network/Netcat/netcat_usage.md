---
weight: 999
url: "/Netcat\\:\_utilisation/"
title: "Netcat: Usage"
description: "Learn how to use Netcat, the Swiss Army knife utility for network operations including port scanning, file transfers, chat, proxying and more."
categories: ["Linux", "Debian", "Ubuntu", "Network"]
date: "2009-11-27T21:29:00+02:00"
lastmod: "2009-11-27T21:29:00+02:00"
tags: ["Network", "Servers", "Development", "Linux", "Ubuntu", "Debian"]
toc: true
---

## Introduction

Netcat is often referred to as a "Swiss Army knife" utility, and for a good reason. Just like the multi-function usefulness of the venerable Swiss Army pocket knife, netcat's functionality is as helpful. Some of its features include port scanning, transferring files, port listening and it can be used a backdoor.

In 2006 netcat was ranked #4 in "[1] Top 100 Network Security Tools" survey, so it's definitely a tool to know.

## Installation

If you're on Debian or Debian based system such as Ubuntu do the following:

```bash
sudo aptitude install netcat
```

## Usage

Let's start with a few very simple examples and build up on those.

If you remember, I said that netcat was a Swiss Army knife. What would a Swiss Army knife be if it also wasn't a regular knife, right? That's why netcat can be used as a replacement of telnet:

```bash
nc www.google.com 80
```

It's actually much more handy than the regular telnet because you can terminate the connection at any time with ctrl+c, and it handles binary data as regular data (no escape codes, nothing).

You may add "-v" parameter for more verboseness, and two -v's (-vv) to get statistics of how many bytes were transmitted during the connection.

Netcat can also be used as a server itself. If you start it as following, it will listen on port 12345 (on all interfaces):

```bash
nc -l -p 12345
```

If you now connect to port 12345 on that host, everything you type will be sent to the other party, which leads us to using netcat as a chat server. Start the server on one computer:

```bash
# On a computer A with IP 10.10.10.10
nc -l -p 12345
```

And connect to it from another:

```bash
# On computer B
nc 10.10.10.10 12345
```

Now both parties can chat!

Talking of which, the chat can be turned to make two processes talk to each other, thus making nc do I/O over network! For example, you can send the whole directory from one computer to another by piping tar to nc on the first computer, and redirecting output to another tar process on the second.

Suppose you want to send files in /data from computer A with IP 192.168.1.10 to computer B (with any IP). It's as simple as this:

```bash
# On computer A with IP 192.168.1.10
tar -cf - /data
```

```bash
# On computer B
nc 192.168.1.10 6666
```

Don't forget to combine the pipeline with pipe viewer from previous article in this series to get statistics on how fast the transfer is going!

A single file can be sent even easier:

```bash
# On computer A with IP 192.168.1.10
cat file
```

```bash
# On computer B
nc 192.168.1.10 6666 > file
```

You may even copy and restore the whole disk with nc:

```bash
# On computer A with IP 192.168.1.10
cat /dev/hdb
```

```bash
# On computer B
nc 192.168.1.10 6666 > /dev/hdb
```

Note: It turns out that "-l" can't be used together with "-p" on a Mac! The solution is to replace "-l -p 6666" with just "-l 6666". Like this:

```bash
nc -l 6666
```

nc now listens on port 6666 on a Mac computer

An uncommon use of netcat is port scanning. Netcat is not the best tool for this job, but it does it ok (the best tool is nmap):

```bash
$ nc -v -n -z -w 1 192.168.1.2 1-1000
(UNKNOWN) [192.168.1.2] 445 (microsoft-ds) open
(UNKNOWN) [192.168.1.2] 139 (netbios-ssn) open
(UNKNOWN) [192.168.1.2] 111 (sunrpc) open
(UNKNOWN) [192.168.1.2] 80 (www) open
(UNKNOWN) [192.168.1.2] 25 (smtp) : Connection timed out
(UNKNOWN) [192.168.1.2] 22 (ssh) open
```

The "-n" parameter here prevents DNS lookup, "-z" makes nc not to receive any data from the server, and "-w 1" makes the connection timeout after 1 second of inactivity.

Another uncommon behavior is using netcat as a proxy. Both ports and hosts can be redirected. Look at this example:

```bash
nc -l -p 12345
```

This starts a nc server on port 12345 and all the connections get redirected to google.com:80. If you now connect to that computer on port 12345 and do a request, you will find that no data gets sent back. That's correct, because we did not set up a bidirectional pipe. If you add another pipe, you can get the data back on another port:

```bash
nc -l -p 12345
```

After you have sent the request on port 12345, connect on port 12346 to get the data.

Probably the most powerful netcat's feature is making any process a server:

```bash
nc -l -p 12345 -e /bin/bash
```

The "-e" option spawns the executable with it's input and output redirected via network socket. If you now connect to the host on port 12345, you may use bash:

```bash
$ nc localhost 12345
ls -las
total 4288
   4 drwxr-xr-x 15 pkrumins users    4096 2009-02-17 07:47 .
   4 drwxr-xr-x  4 pkrumins users    4096 2009-01-18 21:22 ..
   8 -rw-------  1 pkrumins users    8192 2009-02-16 19:30 .bash_history
   4 -rw-r--r--  1 pkrumins users     220 2009-01-18 21:04 .bash_logout
   ...
```

The consequences are that nc is a popular hacker tool as it is so easy to create a backdoor on any computer. On a Linux computer you may spawn /bin/bash and on a Windows computer cmd.exe to have total control over it.

That's everything I can think of. Do you know any other netcat uses that I did not include?

## Advanced usage

### netcat as a portscanner

```bash
nc -v -n -z -w 1 127.0.0.1 22-1000
```

### Create a single-use TCP (or UDP) proxy

Redirect the local port 2000 to the remote port 3000. The same but UDP:

```bash
nc -u -l -p 2000 -c "nc -u example.org 3000"
```

It may be used to "convert" TCP client to UDP server (or viceversa):

```bash
nc -l -p 2000 -c "nc -u example.org 3000"
```

## References

http://www.catonmat.net/blog/unix-utilities-netcat/
