---
weight: 999
url: "/L'adressage_mémoire_et_son_allocation/"
title: "Memory Addressing and Allocation"
description: "A comprehensive guide to memory addressing and allocation in Linux systems, covering virtual addressing spaces, memory allocation techniques, NUMA, TLB optimization, and more."
categories:
  - "Red Hat"
  - "Debian"
  - "Linux"
date: "2012-09-11T16:53:00+02:00"
lastmod: "2012-09-11T16:53:00+02:00"
tags:
  - "Memory"
  - "Kernel"
  - "NUMA"
  - "Performance"
  - "System Administration"
toc: true
---

![Linux](/images/poweredbylinux.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | Kernel 2.6.32+ |
| **Operating System** | Red Hat 6.3<br />Debian 7 |
| **Website** | [Kernel Website](https://www.kernel.org) |
| **Last Update** | 11/09/2012 |
{{< /table >}}

## Memory Addressing

For better efficiency, a computer's memory is divided into blocks (chunks) called pages. Page size may vary depending on the processor architecture (32 or 64 bits). RAM is divided into page frames. One page frame can contain one page. When a process wants to access a memory address, a translation from the page to the page frame must be performed. If this information is not already in memory, the kernel must perform a search to manually load this page into the page frame.

When a program needs to use memory, it uses a 'linear address'. On 32-bit systems, only 4GB of RAM can be addressed. It's possible to bypass this limit with a kernel option called [PAE](https://fr.wikipedia.org/wiki/Extension_d%27adresse_physique)[^1] allowing up to 64GB of RAM. If you need more, you must switch to a 64-bit system which will allow you to use up to 1TB.

Each process has its own page table. Each PTE (Page Table Entry) contains information about the page frame assigned to a process.

## Virtual Address Space

A process's memory in Linux is divided into several sectors:

![Virtual address space and physical address space relationship](/images/virtual_address_space_and_physical_address_space_relationship.avif)[^2]

- text: the code of the executing process, also known as 'text area'
- data: data used by the program. Initialized data will be at the beginning, followed by uninitialized data
- arguments: arguments passed to the program
- environment: environment variables available to the program
- heap: used for dynamic memory allocation (also known as brk)
- stack: used for passing arguments between procedures and dynamic memory

Some processes don't directly manage their memory addressing, leaving two possible solutions:

- The heap and stack can grow toward each other
- No more space available, the process fails to free allocated memory (Memory leaks)

The virtual memory (VMA) of processes can be viewed this way (here PID 1):

```bash
> cat /proc/1/statm
5361 386 306 35 0 104 0
```

Here's another more detailed view (PID 1):

```bash
> cat /proc/1/status
Name:	init
State:	S (sleeping)
Tgid:	1
Pid:	1
PPid:	0
TracerPid:	0
Uid:	0	0	0	0
Gid:	0	0	0	0
Utrace:	0
FDSize:	64
Groups:
VmPeak:	   21452 kB
VmSize:	   21444 kB
VmLck:	       0 kB
VmHWM:	    1544 kB
VmRSS:	    1544 kB
VmData:	     328 kB
VmStk:	      88 kB
VmExe:	     140 kB
VmLib:	    2384 kB
VmPTE:	      60 kB
VmSwap:	       0 kB
Threads:	1
SigQ:	1/11005
SigPnd:	0000000000000000
ShdPnd:	0000000000000000
SigBlk:	0000000000000000
SigIgn:	0000000000001000
SigCgt:	00000001a0016623
CapInh:	0000000000000000
CapPrm:	ffffffffffffffff
CapEff:	fffffffffffffeff
CapBnd:	ffffffffffffffff
Cpus_allowed:	1
Cpus_allowed_list:	0
Mems_allowed:	00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000001
Mems_allowed_list:	0
voluntary_ctxt_switches:	1332
nonvoluntary_ctxt_switches:	42
```

You can also get an even more detailed view of the VMA with pmap (PID 1):

```bash
> pmap 1
1:   /sbin/init
00007f319151d000     44K r-x--  /lib64/libnss_ldap.so.2
00007f3191528000   2044K -----  /lib64/libnss_ldap.so.2
00007f3191727000      4K rw---  /lib64/libnss_ldap.so.2
00007f3191728000     48K r-x--  /lib64/libnss_files-2.12.so
00007f3191734000   2048K -----  /lib64/libnss_files-2.12.so
00007f3191934000      4K r----  /lib64/libnss_files-2.12.so
00007f3191935000      4K rw---  /lib64/libnss_files-2.12.so
00007f3191936000   1572K r-x--  /lib64/libc-2.12.so
00007f3191abf000   2044K -----  /lib64/libc-2.12.so
00007f3191cbe000     16K r----  /lib64/libc-2.12.so
00007f3191cc2000      4K rw---  /lib64/libc-2.12.so
00007f3191cc3000     20K rw---    [ anon ]
00007f3191cc8000     88K r-x--  /lib64/libgcc_s-4.4.6-20120305.so.1.#prelink#.31mBid (deleted)
00007f3191cde000   2044K -----  /lib64/libgcc_s-4.4.6-20120305.so.1.#prelink#.31mBid (deleted)
00007f3191edd000      4K rw---  /lib64/libgcc_s-4.4.6-20120305.so.1.#prelink#.31mBid (deleted)
00007f3191ede000     28K r-x--  /lib64/librt-2.12.so
00007f3191ee5000   2044K -----  /lib64/librt-2.12.so
00007f31920e4000      4K r----  /lib64/librt-2.12.so
00007f31920e5000      4K rw---  /lib64/librt-2.12.so
00007f31920e6000     92K r-x--  /lib64/libpthread-2.12.so
00007f31920fd000   2048K -----  /lib64/libpthread-2.12.so
00007f31922fd000      4K r----  /lib64/libpthread-2.12.so
00007f31922fe000      4K rw---  /lib64/libpthread-2.12.so
00007f31922ff000     16K rw---    [ anon ]
00007f3192303000    252K r-x--  /lib64/libdbus-1.so.3.4.0
00007f3192342000   2048K -----  /lib64/libdbus-1.so.3.4.0
00007f3192542000      4K r----  /lib64/libdbus-1.so.3.4.0
00007f3192543000      4K rw---  /lib64/libdbus-1.so.3.4.0
00007f3192544000     36K r-x--  /lib64/libnih-dbus.so.1.0.0
00007f319254d000   2044K -----  /lib64/libnih-dbus.so.1.0.0
00007f319274c000      4K r----  /lib64/libnih-dbus.so.1.0.0
00007f319274d000      4K rw---  /lib64/libnih-dbus.so.1.0.0
00007f319274e000     96K r-x--  /lib64/libnih.so.1.0.0
00007f3192766000   2044K -----  /lib64/libnih.so.1.0.0
00007f3192965000      4K r----  /lib64/libnih.so.1.0.0
00007f3192966000      4K rw---  /lib64/libnih.so.1.0.0
00007f3192967000    128K r-x--  /lib64/ld-2.12.so
00007f3192b79000     20K rw---    [ anon ]
00007f3192b85000      4K rw---    [ anon ]
00007f3192b86000      4K r----  /lib64/ld-2.12.so
00007f3192b87000      4K rw---  /lib64/ld-2.12.so
00007f3192b88000      4K rw---    [ anon ]
00007f3192b89000    140K r-x--  /sbin/init
00007f3192dab000      8K r----  /sbin/init
00007f3192dad000      4K rw---  /sbin/init
00007f31946b1000    260K rw---    [ anon ]
00007fffa18c1000     84K rw---    [ stack ]
00007fffa18f1000      4K r-x--    [ anon ]
ffffffffff600000      4K r-x--    [ anon ]
 total            21444K
```

The working sets of a process correspond to a group of pages of a process in memory. A process's working set continuously changes throughout the program's life because memory space allocation varies all the time. For memory page allocation, the kernel constantly ensures which pages are not used in the working set. If pages contain modified data, these pages will be written to disk. If there is no data in the pages, they will be reallocated to other processes with more urgent memory needs.

Memory thrashing occurs when the system spends more time moving pages in and out of a process's working set than simply working with them.

## Ulimits

There are limitations on a system. I've already written [an article about this](./ulimit_:_utiliser_les_limites_systèmes.html)[^3]. It's possible to limit memory usage with ulimits. For example, you can define the maximum address space limit.

## Physical Address Space

Most processor architectures support different page sizes. For example:

- IA-32: 4KiB, 2MiB and 4MiB.
- IA-64: 4KiB, 8KiB, 64KiB, 256KiB, 1MiB, 4MiB, 16MiB and 256MiB.

The number of TLB entries is fixed but can be expanded by changing the page size.

## Virtual Address Mapping

When the kernel needs to access a particular memory space in a page frame, it refers to a virtual address. This is sent to the MMU (Memory Management Unit) on the processor, referencing a process's page table. This virtual address will point to a PTE in the page table. The MMU uses information transmitted by the PTE to locate the physical memory page that the virtual one points to. Each PTE contains a bit to indicate whether the page is currently in memory or has been swapped to disk.

![TLB](/images/tlb.avif)[^4]

A page table can be likened to a page directory. Paging is done by breaking down the 32 bits of linear addresses that are used as references to memory positions in several places also known as 'page branching structure'. The last 12 bits reference the memory offset in which the memory page is located. The remaining bits are used to specify the page tables. On a 32-bit system, the 20 bits will require a large page table. The linear address will then be divided into 4 segments:

- PGD: The Page Global Directory
- PMD: The Page Middle Directory
- The page table
- The offset

The PGD and PMD can be viewed as page tables that point to other page tables.

Converting linear addresses to physical ones can take some time, which is why processors have small caches also known as TLB (Translation Lookaside Buffer) that store physical addresses recently associated with virtual ones. The size of a TLB cache is the product of the number of TLBs and a processor's page size. For example, for a processor with 64 TLBs and 4KiB pages, the TLB cache will be 256KiB (64\*4).

## UMA

![UMA](/images/uma.avif)[^5]

On a 32-bit system, the kernel maps all memory up to 896MiB on the 4GiB linear address space. This allows the kernel to have direct memory access below 896MiB by looking at the linear addressing present in the kernel page tables. The kernel directly maps all memory up to 869KiB except for certain reserved regions:

1. 0KiB -> 4KiB: region reserved for BIOS (ZONE_DMA)
2. 4KiB -> 640KiB: region mapped in kernel page tables (ZONE_DMA)
3. 640KiB -> 1MiB: region reserved for BIOS and IO-type devices (ZONE_DMA)
4. 1MiB -> end of kernel data structure: for the kernel and its data structures (ZONE_DMA)
5. end of kernel data structure -> 869MiB: region mapped in kernel page tables (ZONE_NORMAL)
6. 896MiB -> 1GiB: for the kernel to map its linear addresses reserved in ZONE_HIGHMEM (PAE)[^1]

![Uma zones](/images/uma_zones.avif)[^6]

On a 64-bit system, however, it's much simpler as you can see!

## Memory Allocation

COW (Copy on Write) is a form of demand paging. When a process forks a child, the child process inherits the parent's memory addressing. Instead of wasting CPU cycles copying the parent's address spaces to the child, COW ensures that the parent and child share the same address space. COW is called this way because as soon as the child or parent tries to write to a page, the kernel will create a copy of this page and assign it to the address space of the process trying to write. This technology saves a lot of time.

When a process refers to a virtual address, several things are done before approving access to the requested memory space:

1. Verification that the memory address is valid
2. Every reference a process makes for a page in memory does not necessarily give immediate access to a page frame. This check is also done
3. Each PTE of a process in a page table that contains the bit flag specifying whether the page is present in memory or not is checked
4. Access to non-resident pages in memory will generate page faults
5. These page faults may be due to programming errors (as in many cases) and will represent memory accesses that have not yet been allocated by a process or already swapped to disk

It's important to know that:

- As soon as a process requests a page that is not present in memory, it will receive a 'page fault' error.
- When virtual memory needs to allocate a new page frame for a process, a 'minor page fault' will occur (with the help of the MMU).
- When the kernel needs to block a process while it is reading from disk, a 'major page fault' will occur

You can see the current page faults (here the PID is 1):

```bash
> ps -o minflt,majflt 1
MINFLT MAJFLT
  1297      7
```

## Different Types of RAM

Memory cache consists of static random access (SRAM). SRAM is the fastest memory of all RAMs. The advantage of SRAM (static) over dynamic (DRAM) is that it has shorter cycle times and doesn't need to be refreshed after being read. However, SRAM remains very expensive.

There are also:

1. SDRAM (Synchronous Dynamic) uses the processor clock to synchronize IO signals. By coordinating memory accesses, response time is reduced. Basically, a 100Mhz SDRAM equals 10ns access time.
2. DDR (Double Data Rate) is a variant of SDRAM and allows reading as soon as it receives rising/falling clock signals.
3. RDRAM (Rambus Dynamic) uses narrow buses to go very fast and increase throughput.
4. SRAM that we just saw

The RAMs described above range from slowest to fastest.

## NUMA

![NUMA](/images/numa.avif)[^5]

NUMA technology increases MMU performance. A 64-bit processor is mandatory to use this technology. You can check if it's present at your kernel level:

```bash
> grep -i numa /boot/config-`uname -r`
CONFIG_NUMA=y
CONFIG_AMD_NUMA=y
CONFIG_X86_64_ACPI_NUMA=y
CONFIG_NUMA_EMU=y
CONFIG_USE_PERCPU_NUMA_NODE_ID=y
CONFIG_ACPI_NUMA=y
```

Depending on the manufacturers, NUMA technology can be different. This is for example the case between AMD and Intel. For instance, Intel has an MCH controller hub where all memory accesses are routed. This simplifies cache snooping management but can potentially cause bottlenecks for memory accesses. Latency also varies depending on frequency and usage. However, AMD puts different memory ports directly at its CPU level, which surpasses any other SMP technology. With this solution, no bottlenecks, but it complicates cache snooping management which must be managed by all CPUs.

It's possible to have more information about NUMA management on a PID (here 1):

```bash
> cat /proc/1/numa_maps
00400000 default file=/sbin/init mapped=8 N0=8
00608000 default file=/sbin/init anon=1 dirty=1 active=0 N0=1
00609000 default file=/sbin/init anon=1 dirty=1 active=0 N0=1
025c7000 default heap anon=5 dirty=5 active=0 N0=5
7f86cc3c6000 default file=/lib/x86_64-linux-gnu/libdl-2.13.so mapped=2 mapmax=61 N0=2
7f86cc3c8000 default file=/lib/x86_64-linux-gnu/libdl-2.13.so
7f86cc5c8000 default file=/lib/x86_64-linux-gnu/libdl-2.13.so anon=1 dirty=1 active=0 N0=1
7f86cc5c9000 default file=/lib/x86_64-linux-gnu/libdl-2.13.so anon=1 dirty=1 active=0 N0=1
7f86cc5ca000 default file=/lib/x86_64-linux-gnu/libc-2.13.so mapped=121 mapmax=102 N0=121
7f86cc747000 default file=/lib/x86_64-linux-gnu/libc-2.13.so
7f86cc947000 default file=/lib/x86_64-linux-gnu/libc-2.13.so anon=4 dirty=4 active=0 N0=4
7f86cc94b000 default file=/lib/x86_64-linux-gnu/libc-2.13.so anon=1 dirty=1 active=0 N0=1
7f86cc94c000 default anon=4 dirty=4 active=0 N0=4
7f86cc951000 default file=/lib/x86_64-linux-gnu/libselinux.so.1 mapped=14 mapmax=55 N0=14
7f86cc96f000 default file=/lib/x86_64-linux-gnu/libselinux.so.1
7f86ccb6e000 default file=/lib/x86_64-linux-gnu/libselinux.so.1 anon=1 dirty=1 active=0 N0=1
7f86ccb6f000 default file=/lib/x86_64-linux-gnu/libselinux.so.1 anon=1 dirty=1 active=0 N0=1
7f86ccb70000 default anon=1 dirty=1 active=0 N0=1
7f86ccb71000 default file=/lib/x86_64-linux-gnu/libsepol.so.1 mapped=6 N0=6
7f86ccbb0000 default file=/lib/x86_64-linux-gnu/libsepol.so.1
7f86ccdaf000 default file=/lib/x86_64-linux-gnu/libsepol.so.1 anon=1 dirty=1 active=0 N0=1
7f86ccdb0000 default file=/lib/x86_64-linux-gnu/libsepol.so.1 anon=1 dirty=1 active=0 N0=1
7f86ccdb1000 default file=/lib/x86_64-linux-gnu/ld-2.13.so mapped=27 mapmax=101 N0=27
7f86ccfb2000 default anon=4 dirty=4 active=0 N0=4
7f86ccfce000 default anon=2 dirty=2 active=0 N0=2
7f86ccfd0000 default file=/lib/x86_64-linux-gnu/ld-2.13.so anon=1 dirty=1 active=0 N0=1
7f86ccfd1000 default file=/lib/x86_64-linux-gnu/ld-2.13.so anon=1 dirty=1 active=0 N0=1
7f86ccfd2000 default anon=1 dirty=1 active=0 N0=1
7fff85c91000 default stack anon=3 dirty=3 active=1 N0=3
7fff85d22000 default
```

If you want to have finer control of NUMA and decide on processor assignments, [you should look at cpuset]({{< ref "docs/Linux/Kernel/process_latency_and_kernel_timing.md#cpuset" >}})[^7]. It's widely used for applications that need low latency.

### numactl

You can also use the numactl command to force certain CPUs to use specific memory to gain performance. To install it on Red Hat:

```bash
yum install numactl
```

On Debian:

```bash
aptitude install numactl
```

To retrieve information from a machine:

```bash
> numactl --hardware
available: 4 nodes (0-3)
node 0 size: 8058 MB
node 0 free: 7656 MB
node 1 size: 8080 MB
node 1 free: 7930 MB
node 2 size: 8080 MB
node 2 free: 8051 MB
node 3 size: 8080 MB
node 3 free: 8062 MB
node distances:
node   0   1   2   3
 0:  10  20  20  20
 1:  20  10  20  20
 2:  20  20  10  20
 3:  20  20  20  10
```

and:

```bash
> numactl --show
policy: default
preferred node: current
physcpubind: 0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23
cpubind: 0 1 2 3
nodebind: 0 1 2 3
membind: 0 1 2 3
```

#### Assigning a PID to Specific Processors

You can bind a processor to a CPU:

```bash
numactl --physcpubind=0,1,2,3 <PID>
```

To allocate memory on the same NUMA node as the processor:

```bash
numactl --physcpubind=0 --localalloc <PID>
```

### Deactivation

If you wish to disable NUMA, you need to activate "Node Interleaving" in the BIOS. Otherwise at the grub level add (numa=off):

```bash {linenos=table,hl_lines=[6]}
# grub.conf generated by anaconda
#
# Note that you do not have to rerun grub after making changes to this file
# NOTICE:  You have a /boot partition.  This means that
#          all kernel and initrd paths are relative to /boot/, eg.
#          root (hd0,0)
#          kernel /vmlinuz-version ro root=/dev/mapper/vgos-root
#          initrd /initrd-[generic-]version.img
#boot=/dev/sda
default=0
timeout=5
splashimage=(hd0,0)/grub/splash.xpm.gz
hiddenmenu
title Red Hat Enterprise Linux Server (2.6.32-279.2.1.el6.x86_64)
	root (hd0,0)
	kernel /vmlinuz-2.6.32-279.2.1.el6.x86_64 ro root=/dev/mapper/vgos-root rd_NO_LUKS  KEYBOARDTYPE=pc KEYTABLE=fr LANG=en_US.UTF-8 rd_LVM_LV=vgos/root rd_NO_MD rd_LVM_LV=vgos/swap SYSFONT=latarcyrheb-sun16 crashkernel=128M biosdevname=0 rd_NO_DM numa=off
	initrd /initramfs-2.6.32-279.2.1.el6.x86_64.img
title Red Hat Enterprise Linux (2.6.32-279.el6.x86_64)
```

## Improving TLB Performance

The kernel allocates and empties its memory using the "Buddy System" algorithm. The purpose of this algorithm is to avoid external memory fragmentation. These fragmentations occur when there are multiple allocations and deallocations of different sizes of page frames. The memory then becomes fragmented into small blocks of free pages interspersed with blocks of allocated memory. When the kernel receives a request for allocating blocks of a page frame of size N, it will first look for available blocks that can contain this size. If none are available, it will try to find N/2 available blocks.

The "Buddy System" algorithm tries to reorder in the most contiguous way possible. You can view available memory like this:

```bash
> cat /proc/buddyinfo
Node 0, zone      DMA      5      2      2      2      2      2      2      1      2      2      2
Node 0, zone    DMA32   4326  15892   6613   1219    554    188    109      7      0      1      1
Node 0, zone   Normal   2688      0      0      0      0      0      0      0      0      0      1
```

The kernel also supports 'large pages' through the 'hugepages' mechanism (also known as bigpages, largepages, or hugetlbfs). At each context switch encountered, the kernel empties TLB entries for the exiting process.

To determine the size of hugepages:

```bash
> grep -i huge /proc/meminfo
AnonHugePages:         0 kB
HugePages_Total:       0
HugePages_Free:        0
HugePages_Rsvd:        0
HugePages_Surp:        0
Hugepagesize:       2048 kB
```

To choose the size of hugepages, there are 2 solutions:

- With sysctl:

```bash
vm.nr_hugepages=<value>
```

- With grub by adding this parameter:

```bash
[...]
kernel /vmlinuz-2.6.32-279.2.1.el6.....hugepages=<value>
[...]
```

{{< alert context="warning" text="Requesting allocation beyond what the machine can provide will result in a kernel panic" />}}

For applications to use these spaces, they must use system calls like mmap, shmat, and shmget. In the case of mmap, large pages must be available using the hugetlbfs filesystem:

```bash
mkdir /mnt/largepage
mount -t hugetlbfs none /mnt/largepage
```

## References

[^1]: http://fr.wikipedia.org/wiki/Extension_d%27adresse_physique
[^2]: http://en.wikipedia.org/wiki/Virtual_address_space
[^3]: Ulimit : Utiliser les limites systèmes
[^4]: http://www.liafa.jussieu.fr/~carton/Enseignement/Architecture/Cours/Virtual/
[^5]: http://frankdenneman.nl/2010/12/node-interleaving-enable-or-disable/
[^6]: http://www.myexception.cn/linux-unix/515530.html
[^7]: Latence des process et kernel timing#cpuset
