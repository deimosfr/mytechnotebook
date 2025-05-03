---
weight: 999
url: "/Informations_sur_la_mémoire_vive/"
title: "RAM Information"
description: "Understanding memory management in Linux systems, including disk cache usage and how to correctly evaluate free memory."
categories: ["Linux"]
date: "2006-10-17T09:14:00+02:00"
lastmod: "2006-10-17T09:14:00+02:00"
tags: ["Linux", "Memory Management", "System Administration"]
toc: true
---

## Memory Management Overview

### Evidence of Memory Usage

When a system has been running for some time, traditional tools like 'top' often report a surprisingly small amount of free memory. For example, after about 3 hours, the machine on which I'm writing this shows less than 60 MB of free memory, although I have 512 MB on this system. But where has all this memory gone?

Most of it is used by the disk cache, which currently occupies around 290 MB. The *top* command displays this amount under the "cached" column. Memory used for the cache is essentially free, in the sense that it can be quickly reclaimed if a running program (or one that has just been launched) needs it.

### Why?

The reason Linux uses so much memory for disk cache is that unused RAM is simply wasted. Keeping data in the cache means that if something requests previously used data again, there's a good chance it will still be present in the cache.  
Retrieving information from the cache is about 1,000 times faster than reading it from the hard drive.
If the information is not in the cache, the hard drive will need to be read anyway, but in this case, there is no time lost.

### Estimating Free Memory

To get a better estimate of the amount of memory actually free for applications, run the command:

```bash
# Estimating actual free memory:
free -m
```

The -m option means megabytes and the output should look something like this:

```bash
             total       used       free     shared    buffers     cached
Mem:           503        451         52          0         14        293
-/+ buffers/cache:        143        360
Swap:         1027          0       1027
```

The *-/+ buffers/cache* line shows the amounts of free and used memory as seen by applications.  
In general, as long as swap is little used, memory usage has no impact on performance.

Note that I have *512 MB* in my machine but the free command only shows *503* available. This is mainly because the kernel cannot be *swapped* and therefore the memory it occupies will never be available.  
There may also be some regions reserved for/by hardware for various purposes, depending on the system architecture.
