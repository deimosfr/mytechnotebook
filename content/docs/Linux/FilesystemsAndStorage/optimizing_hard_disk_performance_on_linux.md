---
weight: 999
url: "/Optimiser_les_performances_des_disques_dur_sur_Linux/"
title: "Optimizing Hard Disk Performance on Linux"
description: "This guide explains how to optimize disk performance on Linux, covering both mechanical hard drives and SSDs, with various techniques like partition alignment, scheduler optimization, and read-ahead settings."
categories: ["Red Hat", "Debian", "Storage"]
date: "2014-01-02T14:47:00+02:00"
lastmod: "2014-01-02T14:47:00+02:00"
tags:
  [
    "kernel",
    "disk",
    "performance",
    "ssd",
    "schedulers",
    "trim",
    "lvm",
    "noatime",
    "benchmark",
    "io",
  ]
toc: true
---

![Linux](/images/poweredbylinux.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | Kernel 2.6.32+ |
| **Operating System** | Red Hat 6.3<br />Debian 7 |
| **Website** | [Kernel Website](https://www.kernel.org) |
| **Last Update** | 02/01/2014 |
{{< /table >}}

## Introduction

Physical hard drives are currently the slowest components in our machines, whether they are mechanical hard drives with platters or even SSDs! But there are ways to optimize their performance according to specific needs. In this article, we'll look at several aspects that should help you understand why bottlenecks can occur, how to avoid them, and solutions for benchmarking.

## What causes slowness?

There are several factors that can cause slow disk I/O. If you have mechanical disks, they will be slower than SSDs and have additional constraints:

- The rotation speed of the disks
- The reading speed per second will be better on the outer part of the disks (the part furthest from the center)
- Data that is not aligned on the disk
- Small partitions present on the end of the disk
- The speed of the bus on which the disks are connected
- The seek time, corresponding to the time it takes for the read head to move

## Partition alignment

Alignment consists of matching the logical blocks of partitions with the physical blocks to limit read/write operations and thus not hinder performance.

Current SSDs work internally on blocks of 1 or 2 MiB, which is 1,048,576 or 2,097,152 bytes respectively. Considering that a sector stores 512 bytes, it will take 2,048 sectors to store 1,048,576 bytes.
While traditionally operating systems started the first partition at the 63rd sector, the latest versions take into account the constraints of SSDs.
Thus [Parted](./parted_:_résoudre_les_problèmes_de_partionnnement_sur_gros_filesystems.html#partitionnement) can automatically align the beginning of partitions on multiples of 2,048 sectors.

To ensure proper alignment of partitions, enter the following command with administrative privileges and verify that the number of sectors at the beginning of each of your partitions is a multiple of 2,048. Here's the command for an MSDOS partition table:

```bash
fdisk -lu /dev/sdX
```

If you have a GPT partition table:

```bash
parted -l /dev/sdX
```

## Differences between electromechanical disks and SSDs

To better understand the difference between the inner part of the disk (closest to the center) and the outer part of the disk (farthest from the center), let me show you a test between a USB drive (linear, equivalent to an SSD) and a hard drive (non-linear). For this we will use a benchmark tool called bonnie++. So we install it:

```bash
aptitude install bonnie++
```

And we'll launch a capture on both disks with the zcav utility which allows us to test the throughput in raw mode:

```bash
for i in sdc sdd ; do
    zcav -c1 /dev/$i >> ~/$i.zcav
done
```

The -c option allows you to specify the number of times to read the entire disk.

Then we'll generate a graph of the data with Gnuplot:

```bash
> gnuplot
set term png crop
set output '~/zcav_steppings.avif'
set xlabel 'Blocks'
set xrange [0:]
set ylabel 'Disk throughput (MiB/s)'
set yrange [0:]
set border 3
set xtics nomirror
set ytics nomirror
set title 'Zoned Constant Angular Velocity (ZCAV) steppings'
set key below box
plot '/home/pmavro/sdc.zcav' u 1:2 t '160 GiB Hard Drive' with lines, \
'/home/pmavro/sdd.zcav' u 1:2 t '64 GiB USB key' with lines
```

As a reminder, I wrote an [article on Gnuplot](./gnuplot_:_grapher_des_données_facilement.html). Here's the result:

![Zcav steppings.png](/images/zcav_steppings.avif)

We can see very clearly that the hard disk performs very well at the beginning and suffers on the inside. The speed is almost twice as high on the outside as on the inside, and this is explained by the oscillating arm that reads more data over the same period of time on the outside.

## Different data buses

There are different types of buses:

- PCI
- PCI-X
- PCIe
- AGP
- ...

**Today PCI-X is the fastest.** If you have RAID cards, you can check the clock speed to ensure it's running at its maximum to deliver as much throughput as possible. You can find all the necessary information for these buses on Wikipedia. It's important to check:

- Bus size: 32/64 bits
- Clock speed

Here's a small approximate summary (varies with technological advances):

{{< table "table-hover table-striped" >}}
| Object | Latency | Throughput |
|--------|---------|------------|
| 10krpm disk | 3ms | 50MB/s |
| Swap access | 8ms | 50MB/s |
| SSD disk | 0.5ms | 100MB/s |
| Gigabit Ethernet | 1ms | 133MB/s |
| PCI interface | 0.1us | 133MB/s |
| malloc/mmap | 0.1us | |
| fork | 0.1ms | |
| gettimeofday | 1us | |
| context switch | 3us | |
| RAM | 80ns | 8GB/s |
| PCI-Express 16x interface | 10ns | 8GB/s |
| L2 Cache | 5ns | |
| L1 Cache | 1ns | |
| L1 Cache | 0.3ns | 40GB/s |
{{< /table >}}

There are also SCSI type buses. These are a bit special, but **you need to be careful not to mix different clock speeds on the same bus**, bus sizes, passive/active terminations...
You can use the sginfo command to retrieve all the SCSI parameters of your devices:

```bash
aptitude install sg-utils
sginfo -a /dev/sda
```

## Caches and transfer rates

Recent disk controllers have built-in caches to speed up read and write access. By default, many manufacturers disable write caching to avoid any data corruption. However, it's possible to configure this cache to greatly accelerate access. Additionally, when these controllers are equipped with a battery, the cards are capable of keeping data for a few hours to a few days. Once the machine is turned on, the card will take care of writing the data to the disk(s).

To calculate the transfer rate of a disk in bytes/second:

```
rate = (sectors per track * rpm * 512) / 60
```

For disks using ZCAV, you need to replace sectors per track with the average bytes per track:

```
speed = ((average sectors per track * 512 * rpm) /60) / 1000000
```

## I/O requests and caches

![Blockio cache.gif](/images/blockio_cache.avif)

High-level I/O requests such as read/write operations made by the Linux Virtual Filesystem layer must be transformed into block device requests. The kernel then proceeds to queue each block device. Each physical block makes its own queuing request. Queued requests are "Request Descriptors". They describe the data structures that the kernel needs to handle I/O requests. A "request descriptor" can point to an I/O transfer that will in turn point to several disk blocks.

When an I/O request on a device is issued, a special request structure is queued in the "Request Queue" for the device in question. The request structure contains pointers designating the sectors on the disk or the "Buffer Cache (Page Cache)". If the request is:

- to read data, the transfer will be from disk to memory.
- to write data, the transfer will be from memory to disk.

I/O request scheduling is a joint effort. A high-level driver places an I/O request in the request queue. This request is sent to the scheduler which will use an algorithm to process it. To avoid any bottleneck effect, the request will not be processed immediately, but put in blocked or connected mode. Once a certain number of requests is reached, the queue will be disconnected and a low-level driver will handle I/O request transfers to move blocks (disk) and pages (memory).

The unit used for I/O transfer is a page. Each page transferred from the disk corresponds to a page in memory. You can find out the size of cache pages and buffer cache like this:

```bash
> grep -ie '^cache' -ie '^buffer' /proc/meminfo
Buffers:          233184 kB
Cached:          2035636 kB
```

- Buffers: used for storing filesystem metadata
- Cached: used for caching data files

In User mode, programs don't have direct access to the contents of the buffers. Buffers are managed by the kernel in the "kernel space". The kernel must copy data from these buffers into the user mode space of the process that requested the file/inode represented by cache blocks or memory pages.

## Sequential read accesses

{{< alert context="info" text="Using read-ahead technology only makes sense for applications that read data sequentially! There's no benefit for random accesses." />}}

When you make disk accesses, the kernel tries to read data sequentially from the disk. Read-ahead allows reading more blocks than requested to anticipate the demand and store this data in cache. Because when a block is read, it's more than very common to need to read the next block, which is why it can be interesting to tune read-ahead. The advantages of this method are that:

- The kernel is able to respond more quickly to the demand
- The disk controller load is reduced
- Response times are greatly improved

{{< alert context="info" text="The algorithm is designed to stop itself if it detects too many random accesses so as not to impair performance. So don't be afraid to test this feature." />}}

The read-ahead algorithm is managed by 2 values:

- The current window: it controls the amount of data that the kernel will have to process when it makes I/O accesses.
- The ahead window

![Disks read ahead.jpg](/images/disks_read_ahead.avif)

When an application requests page access for reading in the buffer cache (which are part of the current window), the I/Os are done on the ahead window! However, when the application has finished reading on the current window, the ahead window becomes the new current window and a new ahead is created.
**If access to a page is in the current window, the size of the new ahead window will then be increased by 2 pages. If the read-ahead throughput is low, the ahead window size will be gradually reduced.**

To know the size of read-aheads in sectors (1 sector = 512 Bytes):

```bash
> blockdev --getra /dev/sda
256
```

Or in kilobyte:

```bash
> cat /sys/block/sda/queue/read_ahead_kb
128
```

If you want to benchmark to see the best performance you can achieve with your disks:

```bash {linenos=table,hl_lines=["20-21"]}
> DEV="sda" ; for V in 4 8 16 32 64 128 256 512 1024 2048 4096 8192; do echo $V; echo $V > /sys/block/$DEV/queue/read_ahead_kb && hdparm -t /dev/$DEV | grep "Timing"; done
4
 Timing buffered disk reads: 120 MB in  3.03 seconds =  39.58 MB/sec
8
 Timing buffered disk reads: 194 MB in  3.02 seconds =  64.15 MB/sec
16
 Timing buffered disk reads: 268 MB in  3.02 seconds =  88.73 MB/sec
32
 Timing buffered disk reads: 268 MB in  3.00 seconds =  89.25 MB/sec
64
 Timing buffered disk reads: 272 MB in  3.01 seconds =  90.38 MB/sec
128
 Timing buffered disk reads: 272 MB in  3.01 seconds =  90.46 MB/sec
256
 Timing buffered disk reads: 272 MB in  3.01 seconds =  90.24 MB/sec
512
 Timing buffered disk reads: 272 MB in  3.02 seconds =  90.18 MB/sec
1024
 Timing buffered disk reads: 270 MB in  3.00 seconds =  89.93 MB/sec
2048
 Timing buffered disk reads: 272 MB in  3.00 seconds =  90.58 MB/sec
4096
 Timing buffered disk reads: 272 MB in  3.01 seconds =  90.33 MB/sec
8192
 Timing buffered disk reads: 270 MB in  3.00 seconds =  89.99 MB/sec
```

However, you should take this information with a pinch of salt because you would need to test it with the application you want to run on this disk to get a truly satisfactory result. So override the value in /sys to change it. Insert it in `/etc/rc.local` to make it persistent.

{{< alert context="info" text="The initial read-ahead window is equal to half of the configured one. The configured one corresponds to the maximum size of the read-ahead window!" />}}

You can get a report like this:

```bash
> blockdev --report /dev/sda
RO    RA   SSZ   BSZ   StartSec            Size   Device
rw   256   512  1024          0    250000000000   /dev/sda
```

## Schedulers

When the kernel receives multiple I/O requests simultaneously, it must manage them to avoid conflicts. The best solution (from a performance perspective) for disk access is sequential data addressed in logical blocks. In addition, I/O requests are prioritized based on their size. The smaller they are, the higher they are placed in the queue, since the disk will be able to deliver this type of data much more quickly than for large ones.

To avoid bottlenecks, the kernel ensures that all processes get I/Os. It's the scheduler's role to ensure that I/Os at the bottom of the queue are processed and not always postponed.
When adding an entry to the queue, the kernel first tries to expand the current queue and insert the new request into it. If this is not possible, the new request will be assigned to another queue that uses an "elevator" algorithm.

To determine which I/O scheduler (elevator algorithm) is in use:

```bash
> grep CONFIG_DEFAULT_IOSCHED /boot/config-`uname -r`
CONFIG_DEFAULT_IOSCHED="cfq"
```

Here are the schedulers you may find:

- deadline: less efficiency, but less response time
- anticipatory: longer wait times, but better efficiency
- noop: the simplest, designed to save CPU
- cfq: tries to be as homogeneous as possible in all aspects

For more official information: http://www.kernel.org/doc/Documentation/block/

To find out the scheduler currently being used:

```
> cat /sys/block/sda/queue/scheduler
noop deadline [cfq]
```

So it's the algorithm in brackets that's being used. To change the scheduler:

```
> echo noop > /sys/block/sda/queue/scheduler
[noop] deadline cfq
```

Don't forget to put this line in /etc/rc.local if you want it to persist.

{{< alert context="warning" text="Don't do the following on a production machine or you risk having severe slowdowns for a few seconds" />}}

To write all data in the cache to disk, clear the caches, and ensure the use of the newly chosen algorithm:

```bash
sync
sysctl -w vm.drop_caches=3
```

### cfq

This is the default scheduler on Linux, which stands for: Completely Fair Queuing. This scheduler maintains 64 request queues, the IOs are addressed via the Round Robin algorithm to these queues. The addressed requests are used to minimize the movements of the read heads and thus gain speed.

Possible options are:

- quantum: the total number of requests placed on the dispatch queue per cycle
- queued: the maximum number of requests allowed per queue

To tune the CFQ elevator a bit, we need this package installed to have the ionice command:

```bash
aptitude install util-linux
```

ionice allows you to change the read/write priority on a process. Here's an example:

```bash
ionice -p1000 -c2 -n7
```

- -p1: activates the request on PID 1000
- -c2: allows specifying the desired class:
  - 0: none
  - 1: real-time
  - 2: best-effort
  - 3: idle
- -n7: allows specifying the priority on the chosen command/pid between 0 (most important) and 7 (least important)

### deadline

Each scheduler request is assigned an expiration date. When this time has passed, the scheduler moves this request to the disk. To avoid too much solicitation for movements, the deadline scheduler will also handle other requests to a new location on the disk.

It's possible to tune certain parameters such as:

- read_expire: number of milliseconds before each I/O read request expires
- write_expire: number of milliseconds before each I/O write request expires
- fifo_batch: the number of requests to move from the scheduler list to the 'block device' queue
- writes_starved: allows setting the preference on how many times the scheduler must do reads before doing writes. Once the number of reads is reached, the data will be moved to the 'block device' queue and writes will be processed.
- front_merge: a merger (addition) of requests at the bottom of the queue is the normal way requests are processed to be inserted into the queue. After a writes_starved, requests attempt to be added to the beginning of the queue. To disable this feature, set it to 0.

Here are some examples of optimizations I found for [DRBD which uses the deadline scheduler](https://www.drbd.org/users-guide/s-latency-tuning.html):

- Disable front merges:

```bash
echo 0 > /sys/block/<device>/queue/iosched/front_merges
```

- Reduce read I/O deadline to 150 milliseconds (the default is 500ms):

```bash
echo 150 > /sys/block/<device>/queue/iosched/read_expire
```

- Reduce write I/O deadline to 1500 milliseconds (the default is 3000ms):

```bash
echo 1500 > /sys/block/<device>/queue/iosched/write_expire
```

### anticipatory

In many situations, an application that reads blocks, waits, and resumes will read the blocks that follow the blocks just read. But if the desired data is not in the blocks following the last read, there will be additional latency. To avoid this kind of inconvenience, the anticipatory scheduler will respond to this need by trying to find the blocks that will be requested and put them in cache. The performance gain can then be greatly improved.

Read and write access requests are processed in batches. Each batch corresponds in fact to a grouped response time.

Here are the options:

- read_expire: number of milliseconds before each I/O read request expires
- write_expire: number of milliseconds before each I/O write request expires
- antic_expire: how long to wait for another request before reading the next one

### noop

The noop option allows for not using an intelligent algorithm. It serves requests as they come in. It's notably used for [host machines in virtualization](./kvm_:_mise_en_place_de_kvm.html#disks). Or disks that incorporate [TCQ technology](https://en.wikipedia.org/wiki/Tagged_Command_Queuing) to prevent two algorithms from overlapping and causing performance loss instead of gain.

## Optimizations for SSDs

You now understand the importance of options and the differences between disks as explained above. For SSDs, there's some tuning to do if you want to have the best performance while optimizing their lifespan.

### Alignment

One of the first things to do is to create properly aligned partitions. Here's an example of creating [aligned partitions](#partition-alignment):

```bash
datas_device=/dev/sdb
parted -s -a optimal $datas_device mklabel gpt
parted -s -a optimal $datas_device mkpart primary ext4 0% 100%
parted -s $datas_device set 1 lvm on
```

- line 1: we create a gpt type label for large partitions (greater than 2Tb)
- line 2: we create a partition that takes the entire disk
- line 3: we indicate that this partition will be of LVM type

### TRIM

The TRIM function is disabled by default. You'll also need a kernel at least equal to 2.6.33. In order to use TRIM, you'll need to use one of the filesystems designed for SSDs that support this technology:

- Btrfs
- Ext4
- XFS
- JFS

In your fstab, you'll then need to add the 'discard' option to enable TRIM:

```bash {linenos=table,hl_lines=[8,13]}
# /etc/fstab: static file system information.
#
# Use 'blkid' to print the universally unique identifier for a
# device; this may be used with UUID= as a more robust way to name devices
# that works even if disks are added and removed. See fstab(5).
#
# <file system> <mount point>   <type>  <options>       <dump>  <pass>
/dev/mapper/vg-root /               ext4    noatime,nodiratime,discard,errors=remount-ro 0       1
# /boot was on /dev/sda2 during installation
UUID=f41d22fd-a348-42aa-b1a3-4997d19555c8 /boot           ext2    defaults,noatime,nodiratime        0       2
# /boot/efi was on /dev/sda1 during installation
UUID=3104-A1D4  /boot/efi       vfat    defaults        0       1
/dev/mapper/vg-home /home           ext4    noatime,nodiratime,discard         0       2
/dev/mapper/vg-swap none            swap    sw              0       0
```

#### On LVM

It's also possible to enable TRIM on LVM (`/etc/lvm/lvm.conf`):

```bash {linenos=table,hl_lines=[12]}
[...]
    # Issue discards to a logical volumes's underlying physical volume(s) when
    # the logical volume is no longer using the physical volumes' space (e.g.
    # lvremove, lvreduce, etc).  Discards inform the storage that a region is
    # no longer in use.  Storage that supports discards advertise the protocol
    # specific way discards should be issued by the kernel (TRIM, UNMAP, or
    # WRITE SAME with UNMAP bit set).  Not all storage will support or benefit
    # from discards but SSDs and thinly provisioned LUNs generally do.  If set
    # to 1, discards will only be issued if both the storage and kernel provide
    # support.
    # 1 enables; 0 disables.
    issue_discards = 1
[...]
```

### noatime

It's possible to disable access times on files. By default, each time you access a file, the access date is recorded on it. If there are many concurrent accesses on a partition, it ends up being felt enormously. That's why you can disable it if this feature is not useful to you. In your fstab, add the noatime option:

```bash
/dev/mapper/vg-home /home           ext4    defaults,noatime        0       2
```

It's also possible to use the same functionality for folders:

```bash
/dev/mapper/vg-home /home           ext4    defaults,noatime,nodiratime        0       2
```

### Scheduler

Simply use [deadline](#deadline). CFQ is not optimal (although it has been revised for SSDs), we don't want to work unnecessarily. Add the elevator option:

```bash {linenos=table,hl_lines=[9]}
# If you change this file, run 'update-grub' afterwards to update
# /boot/grub/grub.cfg.
# For full documentation of the options in this file, see:
#   info -f grub -n 'Simple configuration'

GRUB_DEFAULT=0
GRUB_TIMEOUT=5
GRUB_DISTRIBUTOR=`lsb_release -i -s 2> /dev/null || echo Debian`
GRUB_CMDLINE_LINUX_DEFAULT="quiet elevator=deadline"
GRUB_CMDLINE_LINUX=""

# Uncomment to enable BadRAM filtering, modify to suit your needs
# This works with Linux (no patch required) and with any kernel that obtains
# the memory map information from GRUB (GNU Mach, kernel of FreeBSD ...)
#GRUB_BADRAM="0x01234567,0xfefefefe,0x89abcdef,0xefefefef"

# Uncomment to disable graphical terminal (grub-pc only)
#GRUB_TERMINAL=console

# The resolution used on graphical terminal
# note that you can use only modes which your graphic card supports via VBE
# you can see them in real GRUB with the command `vbeinfo'
#GRUB_GFXMODE=640x480

# Uncomment if you don't want GRUB to pass "root=UUID=xxx" parameter to Linux
#GRUB_DISABLE_LINUX_UUID=true

# Uncomment to disable generation of recovery mode menu entries
#GRUB_DISABLE_RECOVERY="true"

# Uncomment to get a beep at grub start
#GRUB_INIT_TUNE="480 440 1"
```

#### SSD Detection

It's possible thanks to [UDEV](./Udev_:_Utilisation_d'un_socket_pour_parler_avec_les_devices_kernel.html) to automatically define the scheduler to use depending on the type of disk (platter or SSD):

```bash
# set deadline scheduler for non-rotating disks
ACTION=="add|change", KERNEL=="sd[a-z]", ATTR{queue/rotational}=="0", ATTR{queue/scheduler}="deadline"

# set cfq scheduler for rotating disks
ACTION=="add|change", KERNEL=="sd[a-z]", ATTR{queue/rotational}=="1", ATTR{queue/scheduler}="cfq"
```

### Limiting writes

We'll also limit the use of disk writes by putting tmpfs where temporary files are often written. Insert this into your fstab:

```bash
[...]
tmpfs	/tmp		tmpfs	defaults,noatime,mode=1777	0	0
tmpfs	/var/lock	tmpfs	defaults,noatime,mode=1777	0	0
tmpfs	/var/run	tmpfs	defaults,noatime,mode=1777	0	0
```

## References

http://www.mjmwired.net/kernel/Documentation/block/queue-sysfs.txt
http://www.ocztechnologyforum.com/forum/showthread.php?85495-Tuning-under-linux-read_ahead_kb-wont-hold-its-value
http://static.usenix.org/event/usenix07/tech/full_papers/riska/riska_html/main.html
