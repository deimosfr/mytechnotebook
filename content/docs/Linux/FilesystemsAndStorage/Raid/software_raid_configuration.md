---
weight: 999
url: "/Configuration_d'un_Raid_logiciel/"
title: "Software RAID Configuration"
description: "Learn how to set up, monitor, and optimize software RAID configurations on Linux systems"
categories: ["Linux", "Storage", "System Administration"]
date: "2014-08-08T08:29:00+02:00"
lastmod: "2014-08-08T08:29:00+02:00"
tags: ["raid", "mdadm", "storage", "performance", "linux"]
toc: true
---

## Introduction

Not everyone can afford a RAID 5 card with proper disks. That's why a small software RAID 5 can be a good solution, especially for home use!

## Creating a RAID

### RAID 1

To create a RAID 1, it's simple. You just need 2 disks with 2 partitions of the same size, then run this command:

```bash
mdadm --create --assume-clean --verbose /dev/md0 --level=1 --raid-devices=2 /dev/sdb1 /dev/sdc1
```

- create: creating a raid
- assume-clean: allows having a directly usable raid, without a complete synchronization. This requires being 100% sure that both disks/partitions are blank.
- level: the type of raid (here RAID 1)
- raid-devices: the number of disks used

## Monitoring

To see if everything is working properly, here are several solutions:

- The /proc/mdstat file:

```bash
$ cat /proc/mdstat
Personalities : [raid6] [raid5] [raid4]
md0 : active raid5 sdb1[0] sde1[3] sdd1[2] sdc1[1]
      2930279808 blocks level 5, 64k chunk, algorithm 2 [4/4] [UUUU]

unused devices: <none>
```

- The mdadm command which will allow us to have an exact view of the raid status:

```bash
$ mdadm --detail /dev/md0
/dev/md0:
        Version : 00.90
  Creation Time : Sun Apr 12 13:36:39 2009
     Raid Level : raid5
     Array Size : 2930279808 (2794.53 GiB 3000.61 GB)
  Used Dev Size : 976759936 (931.51 GiB 1000.20 GB)
   Raid Devices : 4
  Total Devices : 4
Preferred Minor : 0
    Persistence : Superblock is persistent

    Update Time : Fri May 22 00:04:17 2009
          State : clean
 Active Devices : 4
Working Devices : 4
 Failed Devices : 0
  Spare Devices : 0

         Layout : left-symmetric
     Chunk Size : 64K

           UUID : ab6f6d7f:29cf4645:9eee7aa6:0e7eea1b
         Events : 0.34

    Number   Major   Minor   RaidDevice State
       0       8       17        0      active sync   /dev/sdb1
       1       8       33        1      active sync   /dev/sdc1
       2       8       49        2      active sync   /dev/sdd1
       3       8       65        3      active sync   /dev/sde1
```

- And finally the best option, to be alerted in case of problems:

```bash
mdadm --monitor --mail=xxx@mycompany.com --delay=1800 /dev/md0
```

I won't go into details, the options speak for themselves.

## Problem Cases

Relative to what we've seen above, here's what happens in case of problems:

```bash
> cat /proc/mdstat
Personalities : [raid6] [raid5] [raid4]
md0 : active raid5 sdc1[1] sde1[3] sdd1[2]
      2930279808 blocks level 5, 64k chunk, algorithm 2 [4/3] [_UUU]

unused devices: <none>
```

And finally, a little mdadm:

```bash
> mdadm --detail /dev/md0
/dev/md0:
        Version : 00.90
  Creation Time : Sun Apr 12 13:36:39 2009
     Raid Level : raid5
     Array Size : 2930279808 (2794.53 GiB 3000.61 GB)
  Used Dev Size : 976759936 (931.51 GiB 1000.20 GB)
   Raid Devices : 4
  Total Devices : 3
Preferred Minor : 0
    Persistence : Superblock is persistent

    Update Time : Sat Dec 26 21:28:19 2009
          State : clean, degraded
 Active Devices : 3
Working Devices : 3
 Failed Devices : 0
  Spare Devices : 0

         Layout : left-symmetric
     Chunk Size : 64K

           UUID : ab6f6d7f:29cf4645:9eee7aa6:0e7eea1b
         Events : 0.3727284

    Number   Major   Minor   RaidDevice State
       0       0        0        0      removed
       1       8       33        1      active sync   /dev/sdc1
       2       8       49        2      active sync   /dev/sdd1
       3       8       65        3      active sync   /dev/sde1
```

## Repairing Your RAID 5

Replace the problematic disk, then add it to your raid:

```bash
> mdadm /dev/md0 -a /dev/sdb1
mdadm: added /dev/sdb1
```

Now, you can monitor the restoration via this command:

```bash
> cat /proc/mdstat
Personalities : [raid6] [raid5] [raid4]
md0 : active raid5 sdb1[4] sdc1[1] sde1[3] sdd1[2]
      2930279808 blocks level 5, 64k chunk, algorithm 2 [4/3] [_UUU]
      [>....................] recovery = 2.3% (23379336/976759936) finish=225.0min speed=70592K/sec

unused devices: <none>
```

We can also view the reconstruction like this:

```bash
> mdadm --detail /dev/md0
/dev/md0:
        Version : 00.90
  Creation Time : Sun Apr 12 13:36:39 2009
     Raid Level : raid5
     Array Size : 2930279808 (2794.53 GiB 3000.61 GB)
  Used Dev Size : 976759936 (931.51 GiB 1000.20 GB)
   Raid Devices : 4
  Total Devices : 4
Preferred Minor : 0
    Persistence : Superblock is persistent

    Update Time : Mon Dec 28 19:38:57 2009
          State : clean, degraded, recovering
 Active Devices : 3
Working Devices : 4
 Failed Devices : 0
  Spare Devices : 1

         Layout : left-symmetric
     Chunk Size : 64K

 Rebuild Status : 2% complete

           UUID : ab6f6d7f:29cf4645:9eee7aa6:0e7eea1b
         Events : 0.3772106

    Number   Major   Minor   RaidDevice State
       4       8       16        0      spare rebuilding   /dev/sdb1
       1       8       33        1      active sync   /dev/sdc1
       2       8       49        2      active sync   /dev/sdd1
       3       8       65        3      active sync   /dev/sde1
```

## Increasing RAID Performance

I won't discuss the different RAID types, but will rather leave Wikipedia for that[^1]. For using software RAID under Linux, **I recommend this documentation**[^2]. We'll focus more on performance since that's the subject here. RAID 0 is the most performant of all raids, but it obviously has its data security problems when a disk is lost.

The MTBF (Mean Time Between Failure) is also important on RAIDs. It's an estimate of the good functioning of the RAID before a disk is detected as failing.

### Chunk Size

The "Chunk size" (or stripe size or element size for some vendors) is the number (in segment size (KiB)) of data written or read for each device before moving to another segment. The algorithm used is Round Robin. The chunk size must be an integer, multiple of the block size. The larger the chunk size, the faster the write speed on very large capacity data, but conversely slower on small data. If the average size of IO requests is smaller than the size of a chunk, then the request will be placed on a single disk of the RAID, canceling all the advantages of RAID. Reducing the chunk size will break large files into smaller pieces that will be distributed across multiple disks, which will improve performance. However, the positioning time of chunks will be reduced. Some hardware doesn't allow writing until a stripe is complete, canceling this positioning latency effect.

A good rule to define the chunk size is to divide roughly the size of IO operations by the number of disks on the RAID (remove parity disks if RAID 5 or 6).

{{% alert context="info" %}}
Quick reminder:

- RAID 0: No parity
- RAID 1: No parity
- RAID 5: 1 parity disk
- RAID 6: 2 parity disks
- RAID 10: No parity disks
  {{% /alert %}}

If you have no idea about your IOs, take a value between 32KB and 128KB, taking a multiple of 2KB (or 4KB if you have larger block sizes). The chunk size (stripe size) is an important factor on the performance of your RAID. If the stripe is too wide, the raid may have a "hot spot" which will be the disk that receives the most IO and will reduce the performance of your RAID. It's obvious that the best performance is when data is spread across all disks. The good formula is therefore:

> Chunk size = average request IO size (avgrq-sz) / number of disks

To get the average request size, I invite you to check the Systat documentation[^3] where we talk about [Iostat]({{< ref "docs/Linux/Misc/sysstat_essential_tools_for_analyzing_performance_issues.md#iostat">}}) and [Sar]({{< ref "docs/Linux/Misc/sysstat_essential_tools_for_analyzing_performance_issues.md#sar">}}).

- To see the chunk size on a RAID (here md0):

```bash
> cat /sys/block/md0/md/chunk_size
131072
```

It's therefore 128KB here.

Here's another way to see it:

```bash {linenos=table,hl_lines=[4]}
> cat /proc/mdstat
Personalities : [raid10]
md0 : active raid10 sdc2[3] sda2[1] sdb2[0] sdd2[2]
      1949426688 blocks super 1.0 128K chunks 2 near-copies [4/4] [UUUU]
unused devices: <none>
```

Or even:

```bash {linenos=table,hl_lines=[20]}
> mdadm --detail /dev/md0
/dev/md0:
        Version : 1.0
  Creation Time : Sat May 12 09:35:34 2012
     Raid Level : raid10
     Array Size : 1949426688 (1859.12 GiB 1996.21 GB)
  Used Dev Size : 974713344 (929.56 GiB 998.11 GB)
   Raid Devices : 4
  Total Devices : 4
    Persistence : Superblock is persistent

    Update Time : Thu Aug 30 12:53:20 2012
          State : clean
 Active Devices : 4
Working Devices : 4
 Failed Devices : 0
  Spare Devices : 0

         Layout : near=2
     Chunk Size : 128K
           Name : N7700:2
           UUID : 1a83e7dc:daa7d822:15a1de4d:e4f6fd19
         Events : 64

    Number   Major   Minor   RaidDevice State
       0       8       18        0      active sync   /dev/sdb2
       1       8        2        1      active sync   /dev/sda2
       2       8       50        2      active sync   /dev/sdd2
       3       8       34        3      active sync   /dev/sdc2
```

- It's possible to define the chunk size when creating the RAID with the argument -c or --chunk. Let's also see how to calculate it best. First, let's use iostat to get the avgrq-sz value:

```bash {linenos=table,hl_lines=["3-7"]}
> iostat -x sda 1 5

avg-cpu: %user  %nice%system%iowait %steal  %idle
           0,21    0,00    0,29    0,05    0,00   99,45

Device:         rrqm/s   wrqm/s     r/s     w/s   rsec/s   wsec/s avgrq-sz avgqu-sz   await  svctm %util
sda               0,71     1,25    1,23    0,76    79,29    15,22    47,55     0,01    2,84   0,73   0,14

avg-cpu: %user  %nice%system%iowait %steal  %idle
           0,00    0,00    0,00    0,00    0,00  100,00

Device:         rrqm/s   wrqm/s     r/s     w/s   rsec/s   wsec/s avgrq-sz avgqu-sz   await  svctm %util
sda               0,00     0,00    1,00    0,00    16,00     0,00    16,00     0,00    1,00   1,00   0,10

avg-cpu: %user  %nice%system%iowait %steal  %idle
           0,00    0,00    0,00    0,00    0,00  100,00
```

Let's then do the calculation to get the chunk size in KiB:

```bash
> echo "47.55*512/1024" | bc -l
23.77500000000000000000
```

We must then divide this value by the number of disks (let's say 2) and round it to the nearest multiple of 2:

> Chunk Size(KB) = 23.775/2 = 11.88 ≈ 8

Here the chunk size to set is 8, since it's the multiple of 2 that is closest to 11.88.

{{< alert context="warning" text="Remember that it's not recommended to go below 32K!" />}}

To create a raid 0 by defining the chunk size:

```bash
mdadm -C /dev/md0 -l 0 -n 2 --chunk-size=32 /dev/sd[ab]1
```

### Stride

The Stride is a parameter that we pass during the construction of a RAID that optimizes the way the filesystem will place its data blocks on the disks before moving to the next ones. With extXn, we can optimize by using the -E option which corresponds to the number of filesystem blocks in a chunk. To calculate the stride:

> Stride = chunk size / block size

For a raid 0 having a chunk size of 64KiB (64 KiB / 4KiB = 16) for example:

```bash
mkfs.ext4 -b 4096 -E stride=16 /dev/mapper/vg1-lv0
```

Some disk controllers do a physical abstraction of block groups making it impossible for the kernel to know them. Here's an example to see the size of a stride:

```bash {linenos=table,hl_lines=[3]}
> dumpe2fs /dev/mapper/vg1-lv0 | grep -i stride
dumpe2fs 1.42 (29-Nov-2011)
RAID stride:              16
```

Here, the size is 16 KiB.

To calculate the stride, there's also a website: http://busybox.net/~aldot/mkfs_stride.html

### Round Robin

RAIDs without parity allow data segmentation across multiple disks to increase performance using the Round Robin algorithm. The segment size is defined at the creation of the RAID and refers to the chunk size.  
The size of a RAID is defined by the smallest disk at the creation of the RAID. The size can vary in the future if all disks are replaced by larger capacity disks. A resynchronization of the disks will take place and the filesystem can be extended.

So for Round Robin tuning, you need to properly tune the chunk size and stride so that the usage of the algorithm is optimal! That's all :-)

### Parity RAIDs

One of the big performance constraints of RAID 5 and 6 is parity calculation. For data to be written, parity calculation must be performed on the raid beforehand. and only then can parity and data be written.

{{< alert context="warning" text="Avoid RAID 5 and 6 if writing your data represents **more than 20% of the activity**" />}}

Each data update requires 4 IO operations:

1. The data to be updated is first read from the disks
2. Update of the new data (but the parity is not yet correct)
3. Reading of blocks of the same stripe and parity calculation
4. Final writing of new data to disks and parity

In RAID 5, it's recommended to use stripe caching:

```bash
echo 256 > /sys/block/md0/md/stripe_cache_size
```

For more information on RAID optimizations: http://kernel.org/doc/Documentation/md.txt[^4][^5]. For the optimization part, look at the following parameters:

- chunk_size
- component_size
- new_dev
- safe_mode_delay
- sync*speed*{min,max}
- sync_action
- stripe_cache_size

### RAID 1

The RAID driver writes to the bitmap when changes have been detected since the last synchronization. A major drawback of RAID 1 is during a power cut, since it has to be entirely rebuilt. With the 'write-intent' bitmap, only the parts that have changed will have to be synchronized, which greatly reduces the reconstruction time.

If a disk fails and is removed from the RAID, md stops erasing bits in the bitmap. If this same disk is reintroduced into the RAID, md will only have to resynchronize the difference. When creating the RAID, if the '--write-intent' bitmap option is combined with '--write-behind', write requests to devices with the '--write-mostly' option will not wait for the requests to be complete before writing to the disk. The '--write-behind' option can be used for RAID1 with slow connections.

The new mdraid matrices support the use of write intent bitmaps. This helps the system identify problematic parts of a matrix; thus, in case of an incorrect stop, the problematic parts will have to be resynchronized, not the entire disk. This drastically reduces the time required for resynchronization. Newly created matrices will automatically have a write intent bitmap added when possible. For example, matrices used as swap and very small matrices (such as /boot matrices) will not benefit from obtaining write intent bitmaps. It's possible to add write intent bitmap to previously existing matrices once the update on the device is completed via the mdadm --grow command. However, write intent bitmaps don't incur a performance impact (about 3-5% on a bitmap size of 65536, but can increase up to 10% or more on smaller bitmaps, such as 8192). This means that if write intent bitmap is added to a matrix, it's better to keep the size relatively large. The recommended size is 65536.[^6]

To see if a RAID is persistent:

```bash {linenos=table,hl_lines=[8]}
> mdadm --detail /dev/md0
/dev/md0:
        Version : 1.0
  Creation Time : Sat May 12 09:35:34 2012
     Raid Level : raid10
     Array Size : 1949426688 (1859.12 GiB 1996.21 GB)
  Used Dev Size : 974713344 (929.56 GiB 998.11 GB)
    Persistence : Superblock is persistent
    Update Time : Thu Aug 30 16:43:17 2012
          State : clean
 Active Devices : 4
Working Devices : 4
 Failed Devices : 0
  Spare Devices : 0

         Layout : near=2
     Chunk Size : 128K

           Name : N7700:2
           UUID : 1a83e7dc:daa7d822:15a1de4d:e4f6fd19
         Events : 64

    Number   Major   Minor   RaidDevice State
       0       8       18        0      active sync   /dev/sdb2
       1       8        2        1      active sync   /dev/sda2
       2       8       50        2      active sync   /dev/sdd2
       3       8       34        3      active sync   /dev/sdc2
```

To add the write intent bitmap (internal):

```bash
mdadm /dev/md0 --grow --bitmap=internal
```

To add the write intent bitmap (external):

```bash
mdadm /dev/md0 --grow --bitmap=/mnt/my_file
```

And to remove it:

```bash
mdadm /dev/md0 --grow --bitmap=none
```

To define the slow disk and the fastest one:

```bash
mdadm -C /dev/md0 -l1 -n2 -b /tmp/md0 --write-behind=256 /dev/sdal --write-mostly /dev/sdbl
```

## FAQ

### I have an md127 appearing and my md0 is broken

First, you need to repair the RAID with mdadm. Then, you need to add the current configuration to mdadm.conf, so that at boot time, it doesn't try to guess a wrong configuration. Simply run this command when your RAID is working properly:

```bash
mdadm --detail --scan --verbose >> /etc/mdadm/mdadm.conf
```

[^7]

## References

[^1]: http://fr.wikipedia.org/wiki/RAID_%28informatique%29
[^2]: **Software RAID Configuration**
[^3]: [Sysstat: Essential tools for analyzing performance problems]({{< ref "docs/Linux/Misc/sysstat_essential_tools_for_analyzing_performance_issues.md" >}})
[^4]: http://kernel.org/doc/Documentation/md.txt
[^5]: http://makarevitch.org/rant/raid/
[^6]: https://access.redhat.com/knowledge/docs/fr-FR/Red_Hat_Enterprise_Linux/6/html/Migration_Planning_Guide/chap-Migration_Guide-File_Systems.html
[^7]: http://www.linuxpedia.fr/doku.php/expert/mdadm
[^8]: http://tldp.org/HOWTO/Software-RAID-HOWTO-6.html
[^9]: [How to resize raid partition](/pdf/how_to_resize_raid_partitions.pdf)
