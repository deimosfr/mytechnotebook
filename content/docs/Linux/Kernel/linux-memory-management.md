---
weight: 999
url: "/La_gestion_de_la_mémoire_sous_Linux/"
title: "Linux Memory Management"
description: "An in-depth guide to memory management in Linux, including page types, dirty/clean page reclamation, OOM handling, memory leak detection, and swap configuration."
categories: ["Linux", "System Administration", "Performance"]
date: "2013-05-06T14:41:00+02:00"
lastmod: "2013-05-06T14:41:00+02:00"
tags: ["Linux", "Memory Management", "Kernel", "Swap", "OOM"]
toc: true
---

{{< table "table-hover table-striped" >}}
| | |
|------|------|
| **Software version** | Kernel 2.6.32+ |
| **Operating System** | Red Hat 6.3<br>Debian 7 |
| **Website** | [Kernel Website](https://www.kernel.org) |
| **Last Update** | 06/05/2013 |
{{< /table >}}

## Pages

When looking at `/proc/meminfo`, 'Inactive Clean' pages correspond to free pages. If the kernel needs to allocate pages to a process, it can take these pages from the free page list or from inactive clean. Pages being used by processes are referenced as active pages. In the case of shared memory, a page can have multiple processes referencing it.

- As long as the page has at least one process referencing it, it will remain in the active list.
- When all processes have released their reference to the page, it becomes inactive. If an active page has been modified by the process that referenced it, it will become an inactive dirty page. Dirty pages contain data that can be written to disk.
- If the page has not been modified since the last read from disk, it will be an inactive clean page. It is then available for allocation.
- Free pages are pages that have not yet been allocated to a process

To find dirty and clean memory:

```bash
awk 'BEGIN {}
/Shared_Clean/{ CLEAN += $2; }
/Shared_Dirty/{ DIRTY += $2; }
END {
print "Shared_Clean: " CLEAN;
print "Shared_Dirty: " DIRTY;
}' /proc/1/smaps
```

## Reclaiming dirty pages

Memory is not unlimited, so the kernel cannot keep dirty pages in RAM forever. Pages that are dirty but no longer used are flushed to disk. The kernel pdflush thread handles this task, and there are at least 2 threads minimum to handle this operation. The use of memory cache creates a strong need for control over how pages are reclaimed. To see the number of pdflush threads currently present:

```bash
> sysctl vm.nr_pdflush_threads
vm.nr_pdflush_threads = 0
```

When flushing dirty pages to disk, the goal is to avoid IO bursts that could saturate the disk. The pdflush daemon writes data to disk in a constant and gentle manner. By default, the 2 pdflush threads will do their actions, but others can be launched (up to 8) if the 2 are saturated to ensure parallel writes on multiple disks. If the ratio of dirty pages to available RAM pages reaches a certain percentage, processes will block while pdflush synchronously writes data. There are several options to improve pdflush:

```bash
# Percentage (total memory) of dirty pages at which pdflush should start writing
vm.dirty_background_ratio=<value>
# Percentage (total memory) of dirty pages at which the process itself should start writing dirty data
vm.dirty_ratio=<value>
# Interval at which pdflush will wake up (100ths/sec) (Observation time). Set to 0 to disable
vm.dirty_writeback_centisecs=<value>
# Defines when data is old enough (100ths/sec) to be intercepted by pdflush (wait time).
vm.dirty_expire_centisecs=<value>
```

If during a ps command, you can see kswap and pdflush, and these 2 elements are in state D (iowait), it's probably caused by the kernel.

To commit all dirty pages and buffers:

```bash
sync
echo s > /proc/sysrq-trigger
```

## Reclaiming clean pages

There are several ways to empty caches. To write all clean pages in the page cache to disk:

```bash
echo 1 > /proc/sys/vm/drop_caches
```

{{< alert context="warning" text="Be careful not to do this during production hours due to the IO it causes" />}}

It is also possible to flush dentries and inodes:

```bash
echo 2 > /proc/sys/vm/drop_caches
```

To flush all clean pages, dentries and inodes:

```bash
echo 3 > /proc/sys/vm/drop_caches
```

## OOM

OOM (Out Of Memory) can happen. There is a process called oomkiller for this. When there is no more swap, no more RAM, it will kill processes. It will trigger if:

- You have no more memory space (including RAM)
- There are no more available pages in the [ZONE_NORMAL or ZONE_HIGHMEM]({{< ref "docs/Linux/Kernel/memory_addressing_and_allocation.md#uma" >}})[^1]
- There is no more available memory in the page mapping table

It is also possible to [add swap on the fly]({{< ref "docs/Linux/FilesystemsAndStorage/swap_creating_dynamic_swap.md" >}})[^2] to avoid a crash due to an OOM.

To see the OOM-Kill immunity level on a process, check the process score (here PID 1):

```bash
> cat /proc/1/oom_score
0
```

To manually request oom-kill to launch:

```bash
echo f > /proc/sysrq-trigger
```

However, it will not kill processes if there is enough memory space. You can check the logs (messages or syslog) to see its result.

It's possible to protect daemons from oom-kill this way:

```bash
echo n > /proc/<PID>/oom_adj
```

- n: corresponds to the score that will be multiplied by 2

{{< alert context="info" text="Note that oom_adj is deprecated on recent kernels. You should use /proc/<PID>/oom_score_adj instead" />}}

{{< alert context="warning" text="Just because you set a process's OOM score doesn't mean its children will inherit it! Be careful about this!" />}}

Finally, it's possible to disable oom-kill:

```bash
vm.panic_on_oom=1
```

A small dedication to budding developers: this is not a solution to fix memory leaks!

## Detecting memory leaks

There are 2 types of memory leaks:

- Virtual: when a process makes requests that are not in the virtual address space (vsize)
- Real: when a process fails to free memory (RSS)

- Use [sar]({{< ref "docs/Linux/Misc/sysstat_essential_tools_for_analyzing_performance_issues.md" >}}) to see system-side exchanges:

```bash
sar -R 1 120
```

- Use the ps command combined with watch:

```bash
watch -n 1 'ps axo pid,comm,rss,vsize | grep sshd'
```

- Use Valgrind for C programs:

```bash
valgrind --tool=memcheck cat /prox/$$/maps
```

To better understand what's happening with ps or top:

- VIRT stands for the virtual size of a process, which is the sum of memory it is actually using, memory it has mapped into itself (for instance the video card's RAM for the X server), files on disk that have been mapped into it (most notably shared libraries), and memory shared with other processes. VIRT represents how much memory the program is able to access at the present moment.

- RES stands for the resident size, which is an accurate representation of how much actual physical memory a process is consuming. (This also corresponds directly to the %MEM column) This will virtually always be less than the VIRT size, since most programs depend on the C or other library.

- SHR indicates how much of the VIRT size is actually sharable memory or libraries. In the case of libraries, it does not necessarily mean that the entire library is resident. For example, if a program only uses a few functions in a library, the whole library is mapped and will be counted in VIRT and SHR, but only the parts of the library file containing the functions being used will actually be loaded in and be counted under RES.

## Swap

It is sometimes necessary to remove all pages or segments of a process from the main memory. In this case, the process will be said to be swapped, and all data belonging to it will be stored in mass memory. This can happen for processes that have been dormant for a long time, while the operating system needs to allocate memory to active processes. Code pages or segments (program) will never be swapped, but simply reassigned, as they can be found in the file corresponding to the program (the executable file). For this reason, the operating system prohibits write access to an executable file that is in use; symmetrically, it is impossible to launch the execution of a file as long as it is held open for write access by another process.[^3]

It's important to allocate enough swap to your systems, even if you have plenty of RAM. If you have multiple disks, don't hesitate to create multiple partitions and give them equal priority at the mount level in fstab:

```bash
[...]
/dev/mapper/vg0-swap none            swap    sw,pri=3              0       0
/dev/mapper/vg1-swap none            swap    sw,pri=3              0       0
/dev/mapper/vg2-swap none            swap    sw,pri=3              0       0
[...]
```

Thanks to this, the kswapd daemon will do round robin, just like a RAID 0 would to increase performance.

How to know what size of swap to allocate to a system? This is a rather complex question that has already generated a lot of debate. Here is a formula that works pretty well:

{{< table "table-hover table-striped" >}}
| RAM | SWAP |
|-----|------|
| Between 1 and 2 GB | 1.5 x the size of RAM |
| Between 2 and 8 GB | equal to the size of RAM |
| More than 8 GB | 0.75 x the size of RAM |
{{< /table >}}

Searching for inactive pages can take CPU. On systems with a lot of RAM, searching for and unmapping inactive pages consume more disk and CPU than writing anonymous pages to disk. It is therefore possible to configure how the kernel will swap. This ranges from 0 to 100. The higher the swappiness value, the more the system is forced to swap, which reduces I/O as shown in the table below[^4]:

{{< table "table-hover table-striped" >}}
| vm.swappiness value | Total I/O | Average Swap |
|---------------------|-----------|--------------|
| 0 | 273.57 MB/s | 0 MB |
| 20 | 273.75 MB/s | 0 MB |
| 40 | 273.52 MB/s | 0 MB |
| 60 | 229.01 MB/s | 23068 MB |
| 80 | 195.63 MB/s | 25587 MB |
| 100 | 184.30 MB/s | 26006 MB |
{{< /table >}}

Here's how to modify swappiness:

```bash
vm.swappiness=60
```

You should know that the kernel likes to swap anonymous pages when the % of memory mapped in the page tables + vm.swappiness >= 100

There are also other interesting parameters to reduce wait time:

```bash
# Number of pages the kernel will read for a page fault. This helps reduce disk head movements. Default is 2 to the power of page-cluster
vm.page-cluster=<value>
# Controls for how long a process is protected from paging when there is a memory dump (in seconds)
vm.swap_token_timeout=<seconds>
```

If there is little memory left, the kernel will start by killing processes in user space, starting with those that make the worst use of memory (memory access relative to the allocation that processes make).

## References

[^1]: [Memory addressing and allocation#UMA](/L'adressage_mémoire_et_son_allocation/#UMA)
[^2]: [SWAP: Creating dynamic swap]({{< ref "docs/Linux/FilesystemsAndStorage/swap_creating_dynamic_swap.md" >}})[^2]
[^3]: [https://fr.wikipedia.org/wiki/Mémoire_virtuelle#Swapping](https://fr.wikipedia.org/wiki/Mémoire_virtuelle#Swapping)
[^4]: [https://www.linuxvox.com/2009/10/what-is-the-linux-kernel-parameter-vm-swappiness/](https://www.linuxvox.com/2009/10/what-is-the-linux-kernel-parameter-vm-swappiness/)

[Memory management and tuning options in Red Hat Enterprise Linux](/pdf/memory_management_and_tuning_options_in_red_hat_enterprise_linux.pdf)
