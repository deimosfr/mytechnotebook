---
weight: 999
url: "/DTrace_\\:_détection_de_problèmes_en_temps_réel/"
title: "DTrace: Real-time Problem Detection"
description: "This article covers DTrace, a powerful tracing system designed for real-time problem detection at kernel or application level, with practical examples and scripts."
categories: ["Linux", "Solaris"]
date: "2008-12-26T18:24:00+02:00"
lastmod: "2008-12-26T18:24:00+02:00"
tags:
  [
    "Troubleshooting",
    "Performance",
    "Monitoring",
    "System Administration",
    "DTrace",
    "Solaris",
  ]
toc: true
---

## Introduction

[DTrace](https://fr.wikipedia.org/wiki/DTrace) is a tracing system designed by Sun Microsystems for real-time problem detection at both kernel and application levels. It has been available since November 2003 and was integrated as part of Solaris 10 in January 2005. DTrace is the first component of the OpenSolaris project whose code was released under the Common Development and Distribution License (CDDL).

DTrace is designed to provide information that allows users to tune applications and the operating system itself. It's built to be used in production environments. The impact of the probes is minimal when tracing is active, and there's no performance impact for inactive probes. This is important because a system includes tens of thousands of probes, many of which can be active.

Tracing programs (often called scripts) are written using the D programming language (not to be confused with D). D is a subset of the C language with additional functions and predefined variables specific to tracing operations. A program written in D structurally resembles a program written in AWK.

Since it can be time-consuming to create your own scripts each time, I've included here all the ones I've used.

## Print Utilization Statistics per Process

Brendan Gregg developed [prustat](/others/prustat.zip) to display the top processes sorted by CPU, Memory, Disk or Network utilization:

```bash
$ prustat -c -t 10 5

 PID   %CPU   %Mem  %Disk   %Net  COMM
7176   0.88   0.70   0.00   0.00  dtrace
7141   0.00   0.43   0.00   0.00  sshd
7144   0.11   0.24   0.00   0.00  sshd
   3   0.34   0.00   0.00   0.00  fsflush
7153   0.03   0.19   0.00   0.00  bash
  99   0.00   0.22   0.00   0.00  nscd
7146   0.00   0.19   0.00   0.00  bash
  52   0.00   0.17   0.00   0.00  vxconfigd
7175   0.07   0.09   0.00   0.00  sh
  98   0.00   0.16   0.00   0.00  kcfd
```

This script is super useful for getting a high-level understanding of what is happening on a Solaris server. Golden!

## File System Flush Activity

On Solaris systems, the pagedaemon is responsible for scanning the page cache and adjusting the MMU reference bit of each dirty page it finds. When the fsflush daemon runs, it scans the page cache looking for pages with the MMU reference bit set, and schedules these pages to be written to disk. The [fsflush.d D script](/others/fsflush.d.zip) provides a detailed breakdown of pages scanned, and the number of nanoseconds that were required to scan "SCANNED" pages:

```bash
$ fsflush.d
  SCANNED   EXAMINED     LOCKED   MODIFIED   COALESCE   RELEASES   TIME(ns)
     4254       4255          1          1          0          0    2695024
     4254       4255          1          0          0          0    1921518
     4254       4255          6          0          0          0    1989044
     4254       4255          1          0          0          0    2401266
     4254       4255          4          1          0          0    2562138
     4254       4255         89          4          0          0    2425988
     4254       3744         80         25          0          0    2394895
     4254       4255         28          8          0          0    1776222
     4254       4255        216          8          0          0    2350826
     4254       4255        108          7          0          0    2356146
```

Now you might be wondering why "SCANNED" is less than "EXAMINED?" This is due to a bug in fsflush, and a bug report was filed to address this anomaly. Tight!

## Seek Sizes

Prior to Solaris 10, determining if an application accessed data in a sequential or random pattern required reviewing mounds of truss(1m) and vxtrace(1m) data. With the introduction of DTrace and Brendan Gregg's [seeksize.d D script](/others/seeksize.d.zip), this question is trivial to answer:

```bash
$ seeksize.d
Sampling... Hit Ctrl-C to end.
^C

     PID  CMD
    7312  dd if=/dev/dsk/c1t1d0s2 of=/dev/null bs=1048576

           value  ------------- Distribution ------------- count
              -1 |                                         0
               0 |@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 1762
               1 |                                         0

       0  sched

           value  ------------- Distribution ------------- count
        -1048576 |                                         0
         -524288 |@@@@                                     1
         -262144 |                                         0
         -131072 |                                         0
          -65536 |                                         0
          -32768 |                                         0
          -16384 |                                         0
           -8192 |                                         0
           -4096 |                                         0
           -2048 |                                         0
           -1024 |                                         0
            -512 |                                         0
            -256 |                                         0
            -128 |                                         0
             -64 |                                         0
             -32 |                                         0
             -16 |                                         0
              -8 |                                         0
              -4 |                                         0
              -2 |                                         0
              -1 |                                         0
               0 |                                         0
               1 |                                         0
               2 |                                         0
               4 |                                         0
               8 |                                         0
              16 |                                         0
              32 |                                         0
              64 |                                         0
             128 |@@@@                                     1
             256 |@@@@                                     1
             512 |@@@@                                     1
            1024 |@@@@                                     1
            2048 |                                         0
            4096 |                                         0
            8192 |                                         0
           16384 |@@@@                                     1
           32768 |@@@@                                     1
           65536 |@@@@@@@@                                 2
          131072 |                                         0
          262144 |                                         0
          524288 |@@@@                                     1
         1048576 |                                         0
```

This script measures the seek distance between consecutive reads and writes, and provides a histogram with the seek distances. For applications that are using sequential access patterns (e.g., dd in this case), the distribution will be small. For applications accessing data in a random nature (e.g., sched in this example), you will see a wide distribution. Shibby!

## Print Overall Paging Activity

Prior to the introduction of DTrace, it was difficult to extract data on which files and disk devices were active at a specific point in time. With the introduction of fspaging.d, you can get a detailed view of which files are being accessed:

```bash
$ fspaging.d
Event      Device                                                    Path RW     Size   Offset
get-page                                           /lib/sparcv9/libc.so.1        8192
get-page                                  /usr/lib/sparcv9/libdtrace.so.1        8192
get-page                                           /lib/sparcv9/libc.so.1        8192
get-page                                                   /lib/libc.so.1        8192
put-page                  /etc/svc/volatile/system-system-log:default.log        8192
put-page                              /etc/svc/volatile/svc_nonpersist.db        8192
put-page                              /etc/svc/volatile/svc_nonpersist.db        8192
put-page                                /etc/svc/volatile/init-next.state        8192
put-page                        /etc/svc/volatile/network-ssh:default.log        8192
put-page                       /etc/svc/volatile/network-pfil:default.log        8192
```

This is a super useful script! Niiiiiiiiiiice!

## Getting System Wide errno Information

When system calls have problems executing, they usually return a value to indicate success or failure, and set the global "ERRNO" variable to a value indicating what went wrong. To get a system-wide view of which system calls are erroring out, we can use Brendan Gregg's [errinfo D script](/others/errinfo.zip):

```bash
$ errinfo -c
Sampling... Hit Ctrl-C to end.
^C
            EXEC          SYSCALL  ERR  COUNT  DESC
          ttymon             read   11      1  Resource temporarily unavailable
           utmpd            ioctl   25      2  Inappropriate ioctl for device
            init            ioctl   25      4  Inappropriate ioctl for device
            nscd         lwp_kill    3     13  No such process
             fmd         lwp_park   62     48  timer expired
            nscd         lwp_park   62     48  timer expired
      svc.startd         lwp_park   62     48  timer expired
           vxesd           accept    4     49  interrupted system call
     svc.configd         lwp_park   62     49  timer expired
           inetd         lwp_park   62     49  timer expired
      svc.startd           portfs   62    490  timer expired
```

This will display the process, system call, and errno number and description from /usr/src/sys/errno.h! Jeah!

## I/O per Process

Several Solaris utilities provide a summary of the time spent waiting for I/O ([which is a meaningless metric](https://www.opensolaris.org/jive/thread.jspa?threadID=719&tstart=0)), but fail to provide facilities to easily correlate I/O activity with a process. With the introduction of [psio.pl](https://brendangregg.com/psio.html), you can see exactly which processes are responsible for generating I/O:

```bash
$ psio.pl
     UID   PID  PPID %I/O    STIME TTY      TIME CMD
    root  7312  7309 70.6 16:00:59 pts/2   02:36 dd if=/dev/dsk/c1t1d0s2 of=/dev/null bs=1048576
    root     0     0  0.0 10:24:18 ?       00:02 sched
    root     1     0  0.0 10:24:18 ?       00:03 /sbin/init
    root     2     0  0.0 10:24:18 ?       00:00 pageout
    root     3     0  0.0 10:24:18 ?       00:51 fsflush
    root     7     1  0.0 10:24:20 ?       00:06 /lib/svc/bin/svc.startd
    root     9     1  0.0 10:24:21 ?       00:14 /lib/svc/bin/svc.configd
```

Once you find I/O intensive processes, you can use [fspaging](/others/fspaging.d.zip), [iosnoop](/others/iosnoop.zip), and [rwsnoop](/others/rwsnoop.zip) to get additional information:

```bash
$ iosnoop -n
MAJ MIN   UID   PID D    BLOCK   SIZE       COMM PATHNAME
136   8     0   990 R   341632   8192     dtrace /lib/sparcv9/ld.so.1
136   8     0   990 R   341568   8192     dtrace /lib/sparcv9/ld.so.1
136   8     0   990 R 14218976   8192     dtrace /lib/sparcv9/libc.so.1
 [ ... ]
```

```bash
$ iosnoop -e
DEVICE    UID   PID D    BLOCK   SIZE       COMM PATHNAME
dad1        0   404 R   481712   8192      vxsvc /lib/librt.so.1
dad1        0     3 W   516320   3072    fsflush /var/adm/utmpx
dad1        0     3 W 18035712   8192    fsflush /var/adm/wtmpx
 [ ... ]
```

```bash
$ rwsnoop
 UID    PID CMD          D   BYTES FILE
 100    902 sshd         R      42 /devices/pseudo/clone@0:ptm
 100    902 sshd         W      80
 100    902 sshd         R      65 /devices/pseudo/clone@0:ptm
 100    902 sshd         W     112
 100    902 sshd         R      47 /devices/pseudo/clone@0:ptm
 100    902 sshd         W      96
   0    404 vxsvc        R    1024 /etc/inet/protocols
  [ ... ]
```

Smooooooooooth!

## I/O Sizes Per Process

As a Solaris administrator, we are often asked to identify application I/O sizes. This information can be acquired for a single process with truss(1m), or system wide with Brendan Gregg's [bitesize.d D script](https://brendangregg.com/DTrace/bitesize.d):

```bash
$ bitesize.d
Sampling... Hit Ctrl-C to end.

    3  fsflush

           value  ------------- Distribution ------------- count
             512 |                                         0
            1024 |@                                        1
            2048 |                                         0
            4096 |@@                                       2
            8192 |@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@    39
           16384 |                                         0

    7312  dd if=/dev/dsk/c1t1d0s2 of=/dev/null bs=1048576

           value  ------------- Distribution ------------- count
              16 |                                         0
              32 |                                         2
              64 |                                         0
             128 |                                         0
             256 |                                         0
             512 |                                         2
            1024 |                                         0
            2048 |                                         0
            4096 |                                         0
            8192 |                                         0
           16384 |                                         0
           32768 |                                         0
           65536 |                                         0
          131072 |@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 76947
          262144 |                                         0
```

If only Dorothy could see this!

## TCP Top

Snoop(1m) and ethereal are amazing utilities, and provide a slew of options to filter data. When you don't have time to wade through snoop data or download and install ethereal, you can use [tcptop](/others/tcptop.zip) to get an overview of TCP activity on a system:

```bash
$ tcptop 5

2005 Jul 19 14:09:06,  load: 0.01,  TCPin:   2679 Kb,  TCPout:     12 Kb

UID    PID LADDR           LPORT RADDR           RPORT      SIZE NAME
  0   7138 192.168.1.3     44084 192.18.108.40      21       544 ftp
  0    352 192.168.1.3        22 192.168.1.8     49805      1308 sshd
100   7134 192.168.1.3     44077 192.168.1.1        22      1618 ssh
  0   7138 192.168.1.3     44089 24.98.83.96     51731   2877524 ftp
```

Now this is some serious bling!

## Who's Paging and DTrace Enhanced vmstat

With Solaris 9, the "-p" option was added to vmstat to break paging activity up into "executable," "anonymous" and "filesystem" page types:

```bash
$ vmstat -p 5
     memory           page          executable      anonymous      filesystem
   swap  free  re  mf  fr  de  sr  epi  epo  epf  api  apo  apf  fpi  fpo  fpf
 1738152 832320 5   9   0   0   0    0    0    0    0    0    0    1    0    0
 1683280 818800 0   2   0   0   0    0    0    0    0    0    0    0    0    0
 1683280 818800 0   0   0   0   0    0    0    0    0    0    0    0    0    0
```

This was super useful information, but unfortunately doesn't provide the executable responsible for the paging activity. With the introduction of [whospaging.d](/others/whospaging.d.zip), you can get paging activity per process:

```bash
$ whospaging.d

Who's waiting for pagein (milliseconds):
Who's on cpu (milliseconds):
  svc.configd                                                      0
  sendmail                                                         0
  svc.startd                                                       0
  sshd                                                             0
  nscd                                                             1
  dtrace                                                           3
  fsflush                                                         14
  dd                                                            1581
  sched                                                         3284
```

Once we get the process name that is responsible for the paging activity, we can use [dvmstat](/others/dvmstat.zip) to break down the types of pages the application is paging (similar to vmstat -p, but per process!):

```bash
$ dvmstat -p 0
    re   maj    mf   fr  epi  epo  api  apo  fpi  fpo     sy
     0     0     0 13280    0    0    0    0    0 13280      0
     0     0     0 13504    0    0    0    0    0 13504      0
     0     0     0 13472    0    0    0    0    0 13472      0
     0     0     0 13472    0    0    0    0    0 13472      0
     0     0     0 13248    0    0    0    0    0 13248      0
     0     0     0 13376    0    0    0    0    0 13376      0
     0     0     0 13664    0    0    0    0    0 13664      0
```

Once we have an idea of which pages are being paged in or out, we can use iosnoop, rwsnoop and fspaging.d to find out which files or devices the application is writing to! Since these rockin' scripts go hand in hand, I am placing them together. Shizam!

And without further ado, number 1 goes to ... (_drum roll_)

## I/O Top

After careful thought, I decided to make [iotop](/others/iotop.zip) and [rwtop](/others/rwtop.zip) #1 on my top ten list. I have long dreamed of a utility that could tell me which applications were actively generating I/O to a given file, device or file system. With the introduction of iotop and rwtop, my wish came true:

```bash
$ iotop 5

2005 Jul 19 13:33:15,  load: 0.24,  disk_r:  95389 Kb,  disk_w:      0 Kb

  UID    PID   PPID CMD              DEVICE  MAJ MIN D            BYTES
    0     99      1 nscd             dad1    136   8 R            16384
    0   7037   7033 find             dad1    136   8 R          2266112
    0   7036   7033 dd               sd7      32  58 R         15794176
    0   7036   7033 dd               sd6      32  50 R         15826944
    0   7036   7033 dd               sd5      32  42 R         15826944
    0   7036   7033 dd               vxio21000 100 21000 R         47448064
```

```bash
$ rwtop 5
2005 Jul 24 10:47:26,  load: 0.18,  app_r:      9 Kb,  app_w:      8 Kb

  UID    PID   PPID CMD              D            BYTES
  100    922    920 bash             R                3
  100    922    920 bash             W               15
  100    902    899 sshd             R             1223
  100    926    922 ls               R             1267
  100    902    899 sshd             W             1344
  100    926    922 ls               W             2742
  100    920    917 sshd             R             2946
  100    920    917 sshd             W             4819
    0    404      1 vxsvc            R             5120
```

## References

http://brendangregg.com/  
[DTrace User Guide](https://docs.sun.com/app/docs/doc/817-6223)  
[Observing I/O Behavior with the DTraceToolkit](https://prefetch.net/articles/observeiodtk.html)  
[DTrace Toolkit](https://brendangregg.com/dtrace.html)  
[DTrace Topics](https://www.solarisinternals.com/wiki/index.php/DTrace_Topics)

### Dtrace GUI

http://www.netbeans.org/kb/docs/ide/NetBeans_DTrace_GUI_Plugin_0_4.html
