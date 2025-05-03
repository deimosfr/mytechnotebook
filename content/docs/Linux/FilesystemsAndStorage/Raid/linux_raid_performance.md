---
weight: 999
url: "/Linux_RAID_performances/"
title: "Linux RAID Performance"
description: "Guide to optimizing Linux RAID performance, including chunk size calculation, stride settings, and performance considerations for different RAID levels."
categories: ["Linux"]
date: "2012-08-31T20:45:00+02:00"
lastmod: "2012-08-31T20:45:00+02:00"
tags: ["RAID", "Performance", "Linux", "Storage"]
toc: true
---

I won't discuss the different RAID types, but rather refer you to Wikipedia for that[^1]. For using software RAID under Linux, [I recommend this documentation](Configuration_d'un_Raid_logiciel.html)[^2]. We'll focus on performance since that's our topic here. RAID 0 is the most performant of all RAID types, but it obviously has data security issues if a disk fails.

The MTBF (Mean Time Between Failure) is also important for RAID systems. This is an estimate of how long the RAID will function properly before a disk is detected as failed.

## Chunk Size

The "Chunk size" (or stripe size or element size for some vendors) is the amount of data (in KiB segments) written or read on each device before moving to another segment. The Round Robin algorithm is used for movement. The chunk size must be an integer multiple of the block size. The larger the chunk size, the faster the write speed for very large data, but inversely slower for small data. If the average size of IO requests is smaller than the chunk size, the request will be placed on a single disk in the RAID, canceling all the advantages of RAID. Reducing the chunk size will split large files into smaller pieces distributed across multiple disks, improving performance. However, the positioning time of the chunks will be reduced. Some hardware does not allow writing until a stripe is complete, eliminating this positioning latency effect.

A good rule for defining chunk size is to divide the size of IO operations by the number of disks in the RAID (minus parity disks if using RAID5 or 6).

{{< table "table-hover table-striped" >}}
| **Notes** |
|-----------|
| Quick reminder: |
| - RAID 0: No parity |
| - RAID 1: No parity |
| - RAID 5: 1 parity disk |
| - RAID 6: 2 parity disks |
| - RAID 10: No parity disks |
{{< /table >}}

If you have no idea about your IO, choose a value between 32KB and 128KB, taking a multiple of 2KB (or 4KB if you have larger block sizes). The chunk size (stripe size) is an important factor in your RAID's performance. If the stripe is too wide, the raid can have a "hot spot" which will be the disk receiving the most IO and will reduce your RAID's performance. It's obvious that the best performance is when data is spread across all disks. The correct formula is therefore:

Chunk size = average request IO size (avgrq-sz) / number of disks

To get the average request size, I invite you to check the Systat documentation[^3] where we talk about [Iostat]({{< ref "docs/Linux/Misc/sysstat_essential_tools_for_analyzing_performance_issues.md#iostat">}}) and [Sar]({{< ref "docs/Linux/Misc/sysstat_essential_tools_for_analyzing_performance_issues.md#sar">}}).

- To see the chunk size on a RAID (here md0):

```bash
> cat /sys/block/md0/md/chunk_size
131072
```

So here it's 128KB.

Here's another way to see it:

```bash
> cat /proc/mdstat
Personalities : [raid10]
md0 : active raid10 sdc2[3] sda2[1] sdb2[0] sdd2[2]
      1949426688 blocks super 1.0 128K chunks 2 near-copies [4/4] [UUUU]
unused devices: <none>
```

Or even:

```bash {linenos=table,hl_lines=[16]}
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

- It's possible to define the chunk size when creating the RAID with the -c or --chunk argument. Let's also see how to calculate it optimally. First, let's use iostat to get the avgrq-sz value:

```bash {linenos=table,hl_lines=[4,5,6,7,8]}
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

Now let's do the calculation to get the chunk size in KiB:

```bash
> echo "47.55*512/1024" | bc -l
23.77500000000000000000
```

We then divide this value by the number of disks (let's say 2) and round it to the nearest multiple of 2:

Chunk Size(KB) = 23.775/2 = 11.88 ≈ 8

Here the chunk size to use is 8, since it's the multiple of 2 closest to 11.88.

{{< alert context="warning" text="Remember that it's not recommended to go below 32K!" />}}

To create a RAID 0 while defining the chunk size:

```bash
mdadm -C /dev/md0 -l 0 -n 2 --chunk-size=32 /dev/sd[ab]1
```

## Stride

Stride is a parameter passed during RAID construction that optimizes how the filesystem places its data blocks on the disks before moving to the next ones. With extXn, you can optimize using the -E option which corresponds to the number of filesystem blocks in a chunk. To calculate the stride:

Stride = chunk size / block size

For a RAID 0 with a chunk size of 64KiB (64 KiB / 4KiB = 16) for example:

```bash
mkfs.ext4 -b 4096 -E stride=16 /dev/mapper/vg1-lv0
```

Some disk controllers make a physical abstraction of block groups, making it impossible for the kernel to know them. Here's an example to see the stride size:

```bash {linenos=table,hl_lines=[3]}
> dumpe2fs /dev/mapper/vg1-lv0 | grep -i stride
dumpe2fs 1.42 (29-Nov-2011)
RAID stride:              16
```

Here, the size is 16 KiB.

To calculate the stride, there is also a website: http://busybox.net/~aldot/mkfs_stride.html

## The Round Robin

Non-parity RAIDs allow data segmentation across multiple disks to increase performance using the Round Robin algorithm. The segment size is defined when creating the RAID and refers to the chunk size.  
The size of a RAID is defined by the smallest disk when creating the RAID. The size can vary in the future if all disks are replaced with larger capacity disks. Disk resynchronization will occur and the filesystem can be extended.

So for tuning Round Robin, you need to correctly tune the chunk size and stride for optimal algorithm usage! That's all :-)

## Parity RAIDs

One of the major performance constraints of RAID 5 and 6 is the parity calculation. For data to be written, parity calculations must be performed on the raid first. Only then can the parity and data be written.

{{< alert context="warning" text="Avoid RAID 5 and 6 if writing data represents <b>more than 20% of activity</b>" />}}

Each data update requires 4 IO operations:

1. Data to be updated is first read from the disks
2. Update of the new data (but parity is not yet correct)
3. Reading blocks from the same stripe and calculating parity
4. Final writing of new data to disks and parity

In RAID 5, it's recommended to use stripe cache:

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

## RAID 1

The RAID driver writes to the bitmap when changes have been detected since the last synchronization. A major drawback of RAID 1 is during a power outage, since it needs to be completely rebuilt. With the 'write-intent' bitmap, only the parts that have changed will need to be synchronized, greatly reducing rebuild time.

If a disk fails and is removed from the RAID, md stops clearing bits in the bitmap. If the same disk is reintroduced into the RAID, md will only need to resynchronize the difference. When creating the RAID, if the '--write-intent' bitmap option is combined with '--write-behind', write requests to devices with the '--write-mostly' option will not wait for the requests to complete before writing to disk. The '--write-behind' option can be used for RAID1 with slow connections.

New mdraid arrays support write intent bitmaps. This helps the system identify problematic parts of an array; thus, in case of an incorrect shutdown, only the problematic parts will need to be resynchronized, not the entire disk. This drastically reduces the time required for resynchronization. Newly created arrays will automatically have a write intent bitmap added when possible. For example, arrays used as swap and very small arrays (such as /boot arrays) will not benefit from getting write intent bitmaps. It is possible to add write intent bitmaps to previously existing arrays once the device update is complete via the mdadm --grow command. However, write intent bitmaps incur performance impact (about 3-5% on a bitmap size of 65536, but can increase to 10% or more on smaller bitmaps, such as 8192). This means that if a write intent bitmap is added to an array, it is better to keep the size relatively large. The recommended size is 65536.[^6]

To see if a RAID is persistent:

```bash {linenos=table,hl_lines=[9]}
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

To add write intent bitmap (internal):

```bash
mdadm /dev/md0 --grow --bitmap=internal
```

To add write intent bitmap (external):

```bash
mdadm /dev/md0 --grow --bitmap=/mnt/mon_fichier
```

And to remove it:

```bash
mdadm /dev/md0 --grow --bitmap=none
```

To define the slow and fastest disk:

```bash
mdadm -C /dev/md0 -l1 -n2 -b /tmp/md0 --write-behind=256 /dev/sdal --write-mostly /dev/sdbl
```

[^1]: http://en.wikipedia.org/wiki/RAID
[^2]: [Configuration of a Software RAID](Configuration_d'un_Raid_logiciel.html)
[^3]: [Sysstat: Essential tools for analyzing performance problems]({{< ref "docs/Linux/Misc/sysstat_essential_tools_for_analyzing_performance_issues.md" >}})
[^4]: http://kernel.org/doc/Documentation/md.txt
[^5]: http://kernel.org/doc/Documentation/md.txt
[^6]: http://www.linux-raid.com/
