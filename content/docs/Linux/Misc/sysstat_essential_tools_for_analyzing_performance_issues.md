---
weight: 999
url: "/Sysstat_\\:_Des_outils_indispensable_pour_analyser_des_problèmes_de_performances/"
title: "Sysstat: Essential Tools for Analyzing Performance Issues"
description: "Learn how to use Sysstat tools including iostat and sar to monitor and analyze system performance issues across Linux and Solaris systems."
categories:
  - "Linux"
  - "Monitoring"
  - "Debian"
date: "2013-08-24T12:57:00+02:00"
lastmod: "2013-08-24T12:57:00+02:00"
tags:
  - "Performance"
  - "System Monitoring"
  - "Linux"
  - "Solaris"
  - "iostat"
  - "sar"
  - "Network"
toc: true
---

![Sysstat](/images/sysstat_logo.avif)

## Introduction

Sysstat is a package containing the sar and iostat binaries. The latter is used to monitor only disk I/O, while sar is used to monitor almost everything.

## Installation

### Debian

On Debian, you need to install sysstat:

```bash
aptitude install sysstat
```

### Red Hat

On Red Hat, you need to install sysstat:

```bash
yum install sysstat
```

### Solaris

On Solaris, you'll need to use [Sun Freeware packages](/Pkg-get_:_Mise_en_place_d'un_système_de_repository_pour_Solaris/) to find this tool:

```bash
pkg-get install CSWsysstat
```

## iostat

iostat allows you to measure disk I/O. If you want to test the performance of disks mounted on your machines, I recommend using [screen]({{< ref "docs/Linux/Applications/Shell/screen_most_used_commands.md" >}}) to see real-time results.

### Linux

On Linux, here's how to use it on the sda disk, for example:

```bash {linenos=table,hl_lines=["4-8"]}
> iostat -x sda 1 5

avg-cpu:  %user   %nice %system %iowait  %steal   %idle
           0,21    0,00    0,29    0,05    0,00   99,45

Device:         rrqm/s   wrqm/s     r/s     w/s   rsec/s   wsec/s avgrq-sz avgqu-sz   await  svctm  %util
sda               0,71     1,25    1,23    0,76    79,29    15,22    47,55     0,01    2,84   0,73   0,14

avg-cpu:  %user   %nice %system %iowait  %steal   %idle
           0,00    0,00    0,00    0,00    0,00  100,00

Device:         rrqm/s   wrqm/s     r/s     w/s   rsec/s   wsec/s avgrq-sz avgqu-sz   await  svctm  %util
sda               0,00     0,00    1,00    0,00    16,00     0,00    16,00     0,00    1,00   1,00   0,10

avg-cpu:  %user   %nice %system %iowait  %steal   %idle
           0,00    0,00    0,00    0,00    0,00  100,00
```

- -x: extended statistics mode
- 1: This means that every second, iostat will analyze the performance of all disks
- 5: For 5 seconds

The first lines of the iostat command provide an average of I/O since the machine was booted. Here are some explanations of the columns:

- r/s: read operations per second
- w/s: write operations per second
- await: wait time (r/s + w/s)

To know if there are sequential reads/writes on the disk:

- rrqm/s: number of merged read requests per second
- wrqm/s: number of merged write requests per second

The more sequential the reads, the faster (on non-SSD disks), the more scattered the sectors, the longer the wait.

For testing, here is some useful information[^1]:

- Advise to drop cache for whole file

```bash
dd if=ifile iflag=nocache count=0
```

- Ensure drop cache for whole file

```bash
dd of=ofile oflag=nocache conv=notrunc,fdatasync count=0
```

- Drop cache for part of file

```bash
dd if=ifile iflag=nocache skip=10 count=10 of=/dev/null
```

- Stream data just using readahead cache

```bash
dd if=ifile of=ofile iflag=nocache oflag=nocache
```

### Solaris

On Solaris, the commands are a bit different:

```bash
iostat -xcnCXTdz 1
```

To stress the disk, we will use the dd command, which is a low-level command:

```bash
dd if=/dev/zero of=/export/home/dd.img bs=10485760 count=100
```

Here is a small shell script that does everything for you (`bench_disk.sh`):

```bash
#!/bin/sh
# Made by Pierre Mavro

echo "What size of file would you like to test (in Mo)? (ex. 10240 for 10Go) :"
read size
echo "Choose your requiered device :"
df | awk '{ print $1 }'
read device
echo ""
echo "Please enter to confirm : a test_array_file file of $size will be created in $device"
read ok
echo ""
echo "Starting disk bench (Ctrl+C to stop)..."
dd if=/dev/zero of=$device/test_array_file bs=1024k count=$size &
iostat -nmCxz 1
```

## sar

sar is a tool that will allow us to monitor many things. When installing sysstat, sar will set itself up in crontab to regularly execute probes placed in `/var/log/sa`.

You can change the default crontab at any time:

```bash
# Run system activity accounting tool every 10 minutes
*/10 * * * * root /usr/lib64/sa/sa1 -S DISK 1 1
# 0 * * * * root /usr/lib64/sa/sa1 -S DISK 600 6 &
# Generate a daily summary of process accounting at 23:53
53 23 * * * root /usr/lib64/sa/sa2 -A
```

Uncomment all lines or adjust according to your needs. To read these files afterward, use the sar command like this:

```bash
sar -d -f /var/log/sa/saXX
```

- XX: day of the month

One important thing before using sar, create an alias in your bashrc or the preferences file for your favorite shell so that the hours are displayed correctly:

```bash
alias sar='LANG=C sar'
```

### Disks

To monitor disks, use the -d option:

```bash {linenos=table,hl_lines=["15-19"]}
> sar -d 1 2

13:57:03          DEV       tps  rd_sec/s  wr_sec/s  avgrq-sz  avgqu-sz     await     svctm     %util
13:57:04      dev8-16      0,00      0,00      0,00      0,00      0,00      0,00      0,00      0,00
13:57:04       dev8-0      0,00      0,00      0,00      0,00      0,00      0,00      0,00      0,00
13:57:04     dev253-0      0,00      0,00      0,00      0,00      0,00      0,00      0,00      0,00
13:57:04     dev253-1      0,00      0,00      0,00      0,00      0,00      0,00      0,00      0,00

13:57:04          DEV       tps  rd_sec/s  wr_sec/s  avgrq-sz  avgqu-sz     await     svctm     %util
13:57:05      dev8-16      0,00      0,00      0,00      0,00      0,00      0,00      0,00      0,00
13:57:05       dev8-0      0,00      0,00      0,00      0,00      0,00      0,00      0,00      0,00
13:57:05     dev253-0      0,00      0,00      0,00      0,00      0,00      0,00      0,00      0,00
13:57:05     dev253-1      0,00      0,00      0,00      0,00      0,00      0,00      0,00      0,00

Moyenne:         DEV       tps  rd_sec/s  wr_sec/s  avgrq-sz  avgqu-sz     await     svctm     %util
Moyenne:     dev8-16      0,00      0,00      0,00      0,00      0,00      0,00      0,00      0,00
Moyenne:      dev8-0      0,00      0,00      0,00      0,00      0,00      0,00      0,00      0,00
Moyenne:    dev253-0      0,00      0,00      0,00      0,00      0,00      0,00      0,00      0,00
Moyenne:    dev253-1      0,00      0,00      0,00      0,00      0,00      0,00      0,00      0,00
```

As you can see, the last lines correspond to the averages since the machine was booted.

If the disk devices don't make much sense to you, you can use the -p option:

```bash
> sar -d -p 1 2

14:13:28          DEV       tps  rd_sec/s  wr_sec/s  avgrq-sz  avgqu-sz     await     svctm     %util
14:13:29          sdb      0,00      0,00      0,00      0,00      0,00      0,00      0,00      0,00
14:13:29          sda      0,00      0,00      0,00      0,00      0,00      0,00      0,00      0,00
14:13:29    VolGroup-lv_root      0,00      0,00      0,00      0,00      0,00      0,00      0,00      0,00
14:13:29    VolGroup-lv_swap      0,00      0,00      0,00      0,00      0,00      0,00      0,00      0,00

14:13:29          DEV       tps  rd_sec/s  wr_sec/s  avgrq-sz  avgqu-sz     await     svctm     %util
14:13:30          sdb      0,00      0,00      0,00      0,00      0,00      0,00      0,00      0,00
14:13:30          sda      0,00      0,00      0,00      0,00      0,00      0,00      0,00      0,00
14:13:30    VolGroup-lv_root      0,00      0,00      0,00      0,00      0,00      0,00      0,00      0,00
14:13:30    VolGroup-lv_swap      0,00      0,00      0,00      0,00      0,00      0,00      0,00      0,00

Moyenne:         DEV       tps  rd_sec/s  wr_sec/s  avgrq-sz  avgqu-sz     await     svctm     %util
Moyenne:         sdb      0,00      0,00      0,00      0,00      0,00      0,00      0,00      0,00
Moyenne:         sda      0,00      0,00      0,00      0,00      0,00      0,00      0,00      0,00
Moyenne:   VolGroup-lv_root      0,00      0,00      0,00      0,00      0,00      0,00      0,00      0,00
Moyenne:   VolGroup-lv_swap      0,00      0,00      0,00      0,00      0,00      0,00      0,00      0,00
```

### CPU

To analyze the CPU:

```bash
> sar -u 1 3

14:20:56        CPU     %user     %nice   %system   %iowait    %steal     %idle
14:20:57        all      0,00      0,00      0,00      0,00      0,00    100,00
14:20:58        all      0,00      0,00      0,00      0,00      0,00    100,00
14:20:59        all      0,00      0,00      0,99      0,00      0,00     99,01
Moyenne:       all      0,00      0,00      0,33      0,00      0,00     99,67
```

### Memory

To monitor memory:

```bash
> sar -r 1 30
Linux 3.2.0-3-amd64 (deb-pmavro) 	13/09/2012 	_x86_64_	(2 CPU)

18:17:01    kbmemfree kbmemused  %memused kbbuffers  kbcached  kbcommit   %commit  kbactive   kbinact
18:17:02       275472   3617496     92,92    220604   1635456   4294436     55,08   2053992   1312976
18:17:03       275612   3617356     92,92    220604   1635140   4294136     55,08   2054064   1312708
18:17:04       268048   3624920     93,11    220616   1642852   4302016     55,18   2054192   1320488
18:17:05       276356   3616612     92,90    220616   1634612   4293660     55,07   2054348   1312128
```

You can also use vmstat to monitor memory:

```bash
> vmstat -n 1 30
procs -----------memory---------- ---swap-- -----io---- -system-- ----cpu----
 r  b   swpd   free   buff  cache   si   so    bi    bo   in   cs us sy id wa
 0  0  52600 272952 220800 1636336    0    1   572   587  471  222 19  3 75  3
 3  0  52600 272828 220804 1635904    0    0     0   140 2969 6041 16  3 79  2
 0  0  52600 275108 220804 1635492    0    0     0     0 3016 6002 18  3 80  0
 0  0  52600 274984 220804 1635184    0    0     0   220 2327 4608 14  3 83  0
 0  0  52600 277836 220804 1635044    0    0     0     0 2868 5820 24  4 72  0
```

- free, buff, and cache: the amount of memory in KiB that is idle
- si and so: correspond to swap usage
- swpd: the size in KiB of swap used

To monitor the rate of change:

```bash
> sar -R 1 30
Linux 3.2.0-3-amd64 (deb-pmavro) 	13/09/2012 	_x86_64_	(2 CPU)

18:18:42      frmpg/s   bufpg/s   campg/s
18:18:43     -2203,00      0,00   2202,00
18:18:44       713,00      0,00      8,00
18:18:45      -279,00      0,00    404,00
18:18:46       186,00      0,00    -13,00
18:18:47       -93,00      0,00     26,00
18:18:48      -155,00      2,00    -18,00
18:18:49        62,00      0,00    -44,00
```

### Swap

To analyze swap:

```bash
> sar -W 1 3

14:22:10     pswpin/s pswpout/s
14:22:11         0,00      0,00
14:22:12         0,00      0,00
14:22:13         0,00      0,00
Moyenne:        0,00      0,00
```

### Network

For network analysis:

```bash
> sar -n DEV 1 2

14:23:47        IFACE   rxpck/s   txpck/s    rxkB/s    txkB/s   rxcmp/s   txcmp/s  rxmcst/s
14:23:48           lo      0,00      0,00      0,00      0,00      0,00      0,00      0,00
14:23:48         eth0      1,00      1,00      0,06      0,17      0,00      0,00      0,00

14:23:48        IFACE   rxpck/s   txpck/s    rxkB/s    txkB/s   rxcmp/s   txcmp/s  rxmcst/s
14:23:49           lo      0,00      0,00      0,00      0,00      0,00      0,00      0,00
14:23:49         eth0      2,00      1,00      0,12      0,38      0,00      0,00      1,00

Moyenne:       IFACE   rxpck/s   txpck/s    rxkB/s    txkB/s   rxcmp/s   txcmp/s  rxmcst/s
Moyenne:          lo      0,00      0,00      0,00      0,00      0,00      0,00      0,00
Moyenne:        eth0      1,50      1,00      0,09      0,28      0,00      0,00      0,50
```

### IO Operations

To monitor all I/O operations:

```bash
> sar -B 1 30
Linux 3.2.0-3-amd64 (deb-pmavro) 	13/09/2012 	_x86_64_	(2 CPU)

18:19:33     pgpgin/s pgpgout/s   fault/s  majflt/s  pgfree/s pgscank/s pgscand/s pgsteal/s    %vmeff
18:19:34         0,00      0,00    106,00      0,00   3105,00      0,00      0,00      0,00      0,00
18:19:35         0,00      0,00     58,00      0,00   1567,00      0,00      0,00      0,00      0,00
18:19:36         0,00      0,00     82,00      0,00   1039,00      0,00      0,00      0,00      0,00
18:19:37         0,00      0,00    205,00      0,00   1530,00      0,00      0,00      0,00      0,00
18:19:38         0,00     44,00    131,00      0,00   1192,00      0,00      0,00      0,00      0,00
```

### Processes

It is possible to get a lot of information about a specific process using the pidstat command:

```bash
> pidstat -p 2365 1 50
Linux 3.2.0-4-amd64 (ZG020194) 	05/07/2013 	_x86_64_	(2 CPU)

02:12:57 PM       PID    %usr %system  %guest    %CPU   CPU  Command
02:12:58 PM      2365    1.00    0.00    0.00    1.00     0  awesome
02:12:59 PM      2365    0.00    0.00    0.00    0.00     0  awesome
02:13:00 PM      2365    0.00    1.00    0.00    1.00     0  awesome
```

It is also possible to monitor I/O (-d) or even the top 5 processes in page fault:

```bash
> pidstat -T CHILD -r 2 5
Linux 3.2.0-4-amd64 (ZG020194) 	05/07/2013 	_x86_64_	(2 CPU)

02:16:02 PM       PID minflt-nr majflt-nr  Command
02:16:04 PM      2252         1         0  VBoxService
02:16:04 PM      2365        50         0  awesome
02:16:04 PM      4938         7         0  firefox
02:16:04 PM      5051       171         0  pidstat
```

## FAQ

### sar: can't open /var/adm/sa/saXX: No such file or directory

You want to use the "sar" command on Solaris to perform monitoring or performance analysis on your server, but when you execute the command, you get an error similar to the following:

```
sar: can't open /var/adm/sa/saXX: No such file or directory
```

The answer is in the manpage for "sadc". You need to execute the following command, and you should be able to execute the command without issue after running this:

```bash
su sys -c "/usr/lib/sa/sadc /var/adm/sa/sa`date +%d`"
```

## Resources
- http://www.cyberciti.biz/open-source/command-line-hacks/linux-monitor-process-using-pidstat/

[^1]: http://comments.gmane.org/gmane.comp.gnu.coreutils.general/904
