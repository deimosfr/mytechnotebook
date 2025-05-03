---
weight: 999
url: "/SystemTap_\\:_Profilez_et_utilisez_rapidement_des_fonctionnalités_du_kernel/"
title: "SystemTap: Profile and Quickly Use Kernel Features"
description: "Learn how to use SystemTap to analyze and diagnose performance issues in Linux systems without kernel recompilation or rebooting."
categories: ["Linux", "RHEL", "Red Hat"]
date: "2011-12-31T14:15:00+02:00"
lastmod: "2011-12-31T14:15:00+02:00"
tags: ["Kernel", "Debugging", "Performance", "Monitoring", "Tracing"]
toc: true
---

![SystemTap](/images/systemtaplogo.avif)

## Introduction

SystemTap is a software tool that simplifies information gathering on a Linux system. It allows you to analyze and diagnose performance or functionality issues. It eliminates the need for kernel recompilation, rebooting, and other steps typically required for low-level data collection.

SystemTap provides a command-line interface similar to awk and a C-like scripting language that allows you to write tools directly for a live kernel. Beyond tracing/probing, it's useful for complex tasks requiring real-time analysis and programmatic responses to events.

SystemTap can profile all system calls that use kprobes. Unlike OProfile, SystemTap captures 100% of events.

One major advantage of SystemTap is that it was designed to be run in production environments. This means you can have:

* A development machine: Where you compile and test SystemTap kernel modules with all the tools and libraries necessary for proper development.
* Production machines: With minimal packages installed, just the desired compiled module retrieved from the development machine.

## Installation

### Red Hat

#### Development

Make sure you have the debuginfo repository. If not, create this file with the following content (adapt as needed):

```bash
[rhel-debuginfo]
name=Red Hat Enterprise Linux $releasever - $basearch - Debug
baseurl=ftp://ftp.redhat.com/pub/redhat/linux/enterprise/$releasever/en/os/$basearch/Debuginfo/
enabled=1
```

Then install these packages:

```bash
yum install systemtap kernel-debuginfo kernel-devel gcc
```

#### Production

As explained in the introduction, the advantage of SystemTap is that on a production machine, there's only one package to install:

```bash
yum install systemtap-runtime
```

When using a module, check compatibility between the development and production environments with this command:

```bash
modinfo modulename | grep vermagic
```

If you want to see all modules with their associated kernel version at once:

```bash {linenos=table,hl_lines=[1]}
> for i in `lsmod | awk '{ print $1 }' | egrep -v '^Module$'` ; do modinfo $i | grep vermagic | xargs echo "$i :" ; done
fuse : vermagic: 2.6.32-220.el6.x86_64 SMP mod_unload modversions
autofs4 : vermagic: 2.6.32-220.el6.x86_64 SMP mod_unload modversions
sunrpc : vermagic: 2.6.32-220.el6.x86_64 SMP mod_unload modversions
nf_conntrack_ipv4 : vermagic: 2.6.32-220.el6.x86_64 SMP mod_unload modversions
nf_defrag_ipv4 : vermagic: 2.6.32-220.el6.x86_64 SMP mod_unload modversions
...
```

## Creating a SystemTap Script

You can find many examples on the SystemTap website: http://sourceware.org/systemtap/examples/

The scripts must use dot notation and support wildcards. Here's what a hello world script would look like in SystemTap:

```c
#! /usr/bin/env stap
probe begin { println("hello world") exit () }
```

### Functions

Many functions already exist and can be used from `/usr/share/systemtap/tapset/*`.
To call a function, use the `probe` keyword. Here's an example:

```c
#! /usr/bin/env stap
probe kernel.function("foo")
probe kernel.function("*").return
```

Here's some information you can find in these scripts to better understand what's happening:

* `ioscheduler.elv_next_request`: Detects when a request (disk read/write type) is retrieved from the queue
* `ioscheduler.elv_next_request.return`: When a request is returned (like a return function in any language)
* `process.exec`: The process is going to execute a new program
* `process.release`: When the desired process will be released from memory (non-Zombie state)
* `netdev.receive`: See the arrival of any network data on all cards
* `tcp.sendmsg`: When the kernel sends a TCP frame
* `vm.pagefault`: See page faults (when memory is physically allocated/when data is taken from swap...)

If you want to get a complete list of all functions available from the kernel:

```bash
stap -p2 -e 'probe kernel.function("*") {}' | sort -u
```

## Execution

### Stap

As you may have understood, the `stap` command is used to run scripts. This command is to be used on a development machine and is divided into 5 levels/steps:

1. Parsing the script (checking the SystemTap script syntax)
2. Comparing and verifying symbols/functions with those in kernel-debuginfo
3. Converting SystemTap code to C
4. Creating a Kernel module
5. Loading and launching the module (as well as unloading it, all requiring root privileges)

It's possible to not create a script and directly run the `stap` command in a shell. You can also define an execution level for this command (for example 2 with the -p option). If you don't specify one, all 5 levels will be executed and the module will remain loaded until it's closed (via Ctrl+C).

Let's look at an example command line that will place a tracer on the kernel's `sys_open()` function and display all calls with their arguments:

```bash
stap -e 'probe syscall.open {printf("%s: %s\n", execname(), argstr)}'
```

*Note: you can use the -k option to keep a trace of your compilation in /tmp to check if there was a problem during a step*

## FAQ

### SystemTap ERROR: Build-id mismatch

If you observe this type of behavior when compiling a SystemTap module, it's likely that the kernel-dev version is different from the version of the kernel currently running on your machine. To verify this, compare the Build-id of the running kernel:

```bash
> eu-readelf -n /boot/vmlinuz-`uname -r` | grep "Build ID"
    Build ID: efbb1bd2e40f890370b8f2fc536c991a2d4abda7
```

With that of the kernel-dev module:

```bash
> eu-readelf -n /usr/lib/debug/lib/modules/my_kernel_version/vmlinux | grep "Build ID"
    Build ID: adcf5270333d375aa5a034523b006373e8f54e48
```

Here we can see there's a version problem, which will result in this type of message:

```
Pass 5: starting execution
ERROR: Build-id mismatch: "kernel" vs. "vmlinux" byte 0 (0xad vs 0xef) address 0xffffffff814f7380 rc 0
Warning: /usr/bin/staprun exited with status: 1
Pass 5: run completed in 10usr/10sys/232real ms.
Pass 5: execution failed.Try again with an additional '--vp 00001' option.
```

To solve your problem, just take the correct version of each [required packages](#installation).
