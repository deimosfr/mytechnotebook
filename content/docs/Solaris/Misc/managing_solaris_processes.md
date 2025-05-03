---
weight: 999
url: "/Manager_les_processes_Solaris/"
title: "Managing Solaris Processes"
description: "A comprehensive guide to managing processes in Solaris operating system, including monitoring with prstat, killing processes, and dealing with zombie processes."
categories: ["Unix", "Solaris", "System Administration"]
date: "2010-04-12T09:52:00+02:00"
lastmod: "2010-04-12T09:52:00+02:00"
tags: ["prstat", "process management", "Solaris", "zombie processes", "system administration"]
toc: true
---

## Introduction

A process is any program that is running on the system. All processes are assigned a unique process identification (PID) number, which is used by the kernel to track and manage the process. The PID numbers are used by the root and regular users to identify and control their processes.

## prstat

The prstat command examines and displays information about active processes on the system.

This command enables you to view information by specific processes, user identification (UID) numbers, central processing unit (CPU) IDs, or processor sets. By default, the prstat command displays information about all processes sorted by CPU usage. To use the prstat command, perform the command:

```bash
# prstat
   PID USERNAME  SIZE   RSS STATE  PRI NICE      TIME  CPU PROCESS/NLWP       
  1641 root     4864K 4520K cpu0    59    0   0:00:00 0.5% prstat/1
  1635 root     1504K 1168K sleep   59    0   0:00:00 0.3% ksh/1
     9 root     6096K 4072K sleep   59    0   0:00:29 0.1% svc.configd/11
   566 root       82M   30M sleep   29   10   0:00:36 0.1% java/14
  1633 root     2232K 1520K sleep   59    0   0:00:00 0.1% in.rlogind/1
   531 root     8200K 2928K sleep   59    0   0:00:12 0.1% dtgreet/1
   474 root       21M 7168K sleep   59    0   0:00:11 0.1% Xsun/1
   236 root     4768K 2184K sleep   59    0   0:00:03 0.0% inetd/4
    86 root     3504K 1848K sleep   59    0   0:00:01 0.0% nscd/24
     7 root     5544K 1744K sleep   59    0   0:00:06 0.0% svc.startd/12
   154 root     2280K  824K sleep   59    0   0:00:01 0.0% in.routed/1
   509 root     6888K 2592K sleep   59    0   0:00:02 0.0% httpd/1
   240 root     5888K 1256K sleep   59    0   0:00:01 0.0% sendmail/1
   145 root     2944K  816K sleep   59    0   0:00:01 0.0% httpd/1
   347 daemon   2608K  776K sleep   59    0   0:00:00 0.0% nfsmapid/3
   206 root     1288K  600K sleep   59    0   0:00:00 0.0% utmpd/1
   344 daemon   2272K 1248K sleep   60  -20   0:00:00 0.0% nfsd/2
   241 smmsp    5792K  960K sleep   59    0   0:00:00 0.0% sendmail/1
   107 root     2584K  784K sleep   59    0   0:00:00 0.0% syseventd/14
   123 root     3064K  880K sleep   59    0   0:00:00 0.0% picld/4
   146 lp       2976K  448K sleep   59    0   0:00:00 0.0% httpd/1
Total: 53 processes, 171 lwps, load averages: 0.02, 0.04, 0.07
```

To quit the prstat command, type q.

The table shows the column headings and their meanings in a prstat report.
Column Headings for the prstat Report.

{{< table "table-hover table-striped" >}}
| Default Column Heading | Description |
|-----------------------|-------------|
| PID | The PID number of the process. |
| USERNAME | The login name or UID of the owner of the process. |
| SIZE | The total virtual memory size of the process. |
| RSS | The resident set size of the process in kilobytes, megabytes, or gigabytes. |
| STATE | The state of the process:<br>* cpu - The process is running on the CPU.<br>* sleep - The process is waiting for an event to complete.<br>* run - The process is in the run queue.<br>* zombie - The process terminated, and the parent is not waiting.<br>* stop - The process is stopped. |
| PRI | The priority of the process. |
| NICE | The value used in priority computation. |
| TIME | The cumulative execution time for the process. |
| CPU | The percentage of recent CPU time used by the process. |
| PROCESS/NLWP | The name of the process/the number of lightweight processes (LWPs) in the process. |
{{< /table >}}

*Note:* The kernel and many applications are now multithreaded. A thread is a logical sequence of program instructions written to accomplish a particular task. Each application thread is independently scheduled to run on an LWP, which functions as a virtual CPU. LWPs in turn, are attached to kernel threads, which are scheduled to run on actual CPUs.

*Note:* Use the priocntl(1) command to assign processes to a priority class and to manage process priorities. The nice(1) command is only supported for backward compatibility to previous Solaris OS releases. The priocntl command provides more flexibility in managing processes.

The table shows the options for the prstat command.

{{< table "table-hover table-striped" >}}
| Option | Description |
|--------|-------------|
| -a | Displays separate reports about processes and users at the same time. |
| -c | Continuously prints new reports below previous reports. |
| -n nproc | Restricts the number of output lines. |
| -p pidlist | Reports only on processes that have a PID in the given list. |
| -s key | Sorts output lines by key in descending order. The five possible keys include: cpu, time, size, rss, and pri. You can use only one key at a time. |
| -S key | Sorts output lines by key in ascending order. |
| -t | Reports total usage summary for each user. |
| -u euidlist | Reports only processes that have an effective user ID (EUID) in the given list. |
| -U uidlist | Reports only processes that have a real UID in the given list. |
{{< /table >}}

## Kill Frozen Process

You use the kill or pkill commands to terminate one or more processes.

The format for the kill command is:

```bash
kill -signal PID
```

To show all of the available signals used with the kill command:

```bash
kill -l
```

The format for the pkill command is:

```bash
pkill -signal Process
```

Before you can terminate a process, you must know its name or PID. Use either the ps or pgrep command to locate the PID for the process.

The following examples uses the pgrep command to locate the PID for the mail processes.

```bash
# pgrep -l mail
  241 sendmail
  240 sendmail
# pkill sendmail
```

The following examples use the ps and pkill commands to locate and terminate the sendmail process.

```bash
# ps -e | grep sendmail
   241 ?           0:00 sendmail
   240 ?           0:02 sendmail
# kill 241
```

To terminate more than one process at the same time, use the following syntax:

```bash
 kill -signal PID PID PID PID
 pkill signal process process
```

You use the kill command without a signal on the command line to send the default Signal 15 to the process. This signal usually causes the process to terminate.

The table shows some signals and names.

{{< table "table-hover table-striped" >}}
| Signal Number | Signal Name | Event | Default Action |
|--------------|------------|-------|---------------|
| 1 | SIGHUP | Hangup | Exit |
| 2 | SIGINT | Interrupt | Exit |
| 9 | SIGKILL | Kill | Exit |
| 15 | SIGTERM | Terminate | Exit |
{{< /table >}}

* 1, SIGHUP - A hangup signal to cause a telephone line or terminal connection to be dropped. For certain daemons, such as inetd and in.named, a hangup signal will cause the daemon to reread its configuration file.
* 2, SIGINT - An interrupt signal from your keyboard--usually from a Control-C key combination.
* 9, SIGKILL - A signal to kill a process. A process cannot ignore this signal.
* 15, SIGTERM - A signal to terminate a process in an orderly manner. Some processes ignore this signal.

A complete list of signals that the kill command can send can be found by executing the command kill -l, or by referring to the man page for signal:

```bash
man -s3head signal
```

Some processes can be written to ignore Signal 15. Processes that do not respond to a Signal 15 can be terminated by force by using Signal 9 with the kill or pkill commands. You use the following syntax:

```bash
kill -9 PID
pkill -9 process
```

Caution: Use the kill -9 or pkill -9 command as a last resort to terminate a process. Using the -9 signal on a process that controls a database application or a program that updates files can be disastrous. The process is terminated instantly with no opportunity to perform an orderly shutdown.

When a workstation is not responding to your keyboard or mouse input, the CDE might be frozen. In such cases, you may be able to remotely access your workstation by using the rlogin command or by using the telnet command from another system.
Killing the Process for a Frozen Login

After you are connected remotely to your system, you can invoke the pkill command to terminate the corrupted session on your workstation.

In the following examples, the rlogin command is used to log in to sys42, from which you can issue a pkill or a kill command.

```bash
# rlogin sys-02
Password: 
Last login: Sun Oct 24 13:44:51 from sys-01
Sun Microsystems Inc.   SunOS 5.10      s10_68  Sep. 20, 2004
# pkill -9 Xsun
```

or

```bash
# ps -e | grep Xsun
   442 ?           0:01 Xsun
# kill -9 442
```

## Kill Zombie Processes

Sometimes we encounter processes called zombies. Nothing to do with a George A. Romero movie, it's a much more prosaic phenomenon. When a process ends, almost all associated resources are released, except for the corresponding entry in the OS process table. The reason is simple: the parent process must be able to retrieve the return code of its child process, so we cannot abruptly erase everything.

Typically, the entry is removed from the process table when the parent retrieves this return code, which is called reaping (the Reaper being our Grim Reaper). If the parent process, for one reason or another, does not read this code, the entry remains in the process table. Let's see what these zombies look like, and how to get rid of them.

In itself, it's generally not very troublesome, since apart from the entry in the process table, there's nothing left: no memory consumed, no CPU used. However, I see two disadvantages:

* the process table has a limited size (30000 entries by default), and each zombie occupies a spot, so if many are generated, this can become problematic on a busy system
* the zombie is really, really ugly (and it doesn't smell good either)

Fortunately, it is possible to substitute for the negligent parent and finally allow these zombies to find eternal rest. Let's first see how a zombie appears on the system:

```bash
$ ps -edf 
```

Or:

```bash
$ ps -ecl 
```

To remove them, we'll use the preap command:

```bash
$ preap 7450
7450: exited with status 0
$ preap 7542
7542: exited with status 0
$ preap 7544
7544: exited with status 0
$ preap 7546
7546: exited with status 0
$ ps -edf | grep defunct
```

And there you have it, much more effective than the hero of Brain Dead! Note, if you want to manually create a zombie for your tests, you can use the following command, provided by the excellent c0t0d0s0.org:

```bash
nohup perl -e "if (fork()>0) {while (1) {sleep 100*100;};};"
```
