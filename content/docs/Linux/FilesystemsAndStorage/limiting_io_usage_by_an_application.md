---
weight: 999
url: "/O_par_une_application/"
title: "Limiting I/O usage by an application"
description: "Learn how to identify and limit I/O usage by applications to improve system performance and responsiveness."
categories: ["Linux"]
date: "2009-11-19T07:07:00+02:00"
lastmod: "2009-11-19T07:07:00+02:00"
tags: ["Ionice", "Linux", "Performance", "System Administration"]
toc: true
---

## Introduction

Want to know why your load average is so high? Run this command to see what processes are on the run queue. Runnable processes have a status of "R", and commands waiting on I/O have a status of "D".

Once found, you may need to reduce its I/O requests, so we'll use ionice.

## Get I/O apps

To get the biggest I/O consuming applications:

```bash
ps -eo stat,pid,user,command 
```

On some older versions of Linux may require -emo instead of -eo.

And on Solaris:

```bash
ps -aefL -o s -o user -o comm 
```

## Ionice

ionice limits process I/O, to keep it from swamping the system (Linux)

This command is somewhat similar to 'nice', but constrains I/O usage rather than CPU usage. In particular, the '-c3' flag tells the OS to only allow the process to do I/O when nothing else is pending. This dramatically increases the responsiveness of the rest of the system if the process is doing heavy I/O.

There's also a '-p' flag, to set the priority of an already-running process.

```bash
ionice -c3 find /
```
