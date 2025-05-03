---
weight: 999
url: "/fifo-and-cat-share-session-with-multiple-users/"
title: "Fifo and Cat: Share a Session with Multiple Users"
description: "How to use named pipes (FIFOs) and cat to visually share a terminal session with other users"
categories: ["Linux", "System Administration", "Command Line"]
date: "2007-03-18T09:29:00+02:00"
lastmod: "2007-03-18T09:29:00+02:00"
tags: ["Fifo", "Cat", "Script", "Terminal", "Sharing"]
toc: true
---

## Introduction

This technique can be used to show a remote colleague or client what you're working on on the server.

It's just for visual sharing and not interactive like with screen ([see this documentation]({{< ref "docs/Linux/Applications/Shell/screen_most_used_commands.md" >}})).

## Creating a fifo file

```bash
mkfifo /tmp/sortieScript
ls -l /tmp/sortieScript
prw-r--r-- 1 yannick yannick 0 Jul 6 02:59 /tmp/sortieScript
```

```
mkfifo - Create named pipes (FIFOs) with the given NAMEs.
A FIFO special file (a named pipe) is similar to a pipe, except that it is accessed as part of the file system.
[...] the FIFO special file has no contents on the file system
```

## Reading the file by the remote user

```bash
cat /tmp/sortieScript
```

Warning: as long as the file is not "catted", it cannot be used by the following command...

```
When a process tries to write to a FIFO that is not opened for read on the other side, the process is sent a SIGPIPE signal.
```

## Output script to this file

```bash
script -f /tmp/sortieScript
```

```
Script started, file is /tmp/sortieScript
```

From now on, everything that is typed is visible to the person "catting" the sortieScript file, including interactive sessions like vi...

Stop logging to the file with CTRL-D

```
-f Flush output after each write. This is nice for telecooperation [...]
```

## Amazing: the demo

![demo](/images/scriptfmkfifo.avif)
