---
weight: 999
url: "/Strace_et_Ltrace_\\:_tracez_les_appels_systèmes_et_librairies/"
title: "Strace and Ltrace: Trace System and Library Calls"
description: "How to use strace and ltrace tools to monitor system calls and library calls for debugging and troubleshooting on Linux systems"
categories: ["Debian", "Linux"]
date: "2012-01-20T10:32:00+02:00"
lastmod: "2012-01-20T10:32:00+02:00"
tags:
  [
    "debugging",
    "system calls",
    "strace",
    "ltrace",
    "troubleshooting",
    "Linux",
    "monitoring",
  ]
toc: true
---

## Introduction

strace is a debugging tool on Linux used to monitor system calls made by a program and all the signals it receives, similar to the "truss" tool on other Unix systems. It's made possible through a Linux kernel feature called ptrace.

The most common use is to launch a program using strace, which displays a list of system calls made by the program. This is useful when a program continually crashes or doesn't behave as expected. For example, using strace can reveal that the program is trying to access a file that doesn't exist or can't be read.

Another use is to use the -p option to attach it to a running program. This is useful when a program stops responding, and can reveal, for example, that the process is blocked waiting to make a network connection.

Since strace only details system calls, it can't be used as a code debugger like Gdb. However, it remains simpler to use than a code debugger and is an extremely useful tool for system administrators.

In this documentation, I won't discuss ltrace much because its usage is quite similar to strace.

## Installation

### Debian

To install on Debian:

```bash
aptitude install strace ltrace
```

### Red Hat

```bash
yum install strace ltrace
```

## Usage

For example, if we want to debug an issue with an Apache server:

```bash
> strace -f /etc/init.d/httpd restart
execve("/etc/init.d/httpd", ["/etc/init.d/httpd", "restart"], [/* 32 vars */]) = 0
brk(0)                                  = 0x1c61000
mmap(NULL, 4096, PROT_READ|PROT_WRITE, MAP_PRIVATE|MAP_ANONYMOUS, -1, 0) = 0x7ff218642000
access("/etc/ld.so.preload", R_OK)      = -1 ENOENT (No such file or directory)
open("/etc/ld.so.cache", O_RDONLY)      = 3
fstat(3, {st_mode=S_IFREG|0644, st_size=31346, ...}) = 0
mmap(NULL, 31346, PROT_READ, MAP_PRIVATE, 3, 0) = 0x7ff21863a000
close(3)                                = 0
open("/lib64/libtinfo.so.5", O_RDONLY)  = 3
read(3, "\177ELF\2\1\1\0\0\0\0\0\0\0\0\0\3\0>\0\1\0\0\0@\310@\331=\0\0\0"..., 832) = 832
...
```

The -f option of strace traces child processes as they are created by currently traced processes following the fork system call.

All you need to do is analyze the lines to see the issue. This can be tedious depending on the number of lines, but generally the information about your problem is here.

### Redirecting Output to a File

If we want to redirect all of strace's output (initially on error output) to a file using the -o option:

```bash
strace -o apache -f /etc/init.d/httpd restart
```

### Working with Standard Output

As you now know, strace works on error output, so if you want to work on it with grep or other commands on-the-fly (without redirecting to a file), you'll need to use redirection:

```bash
strace -f /etc/init.d/httpd restart 2>&1
```

### Working with Specific Kernel Calls Only

If you want to get only open and access type calls, for example:

```bash
strace -e open,access ls
```

Here are some examples of system calls you can try:

```bash
strace -e trace=set
strace -e trace=open
strace -e trace=read
strace -e trace=file
strace -e trace=process
strace -e trace=network
strace -e trace=signal
strace -e trace=ipc
strace -e trace=desc //descriptors
strace -e read=set
```

### Increasing the Number of Characters to Display

You can increase the display size using the -s option followed by the desired size (5000 for example):

```bash
strace -o apache -f -s 5000 /etc/init.d/httpd restart
```

### Attaching to an Existing PID

If we want to trace a process that's already running, it's possible. To do this, simply use the -p argument:

```bash
strace -f -s 5000 -p <PID>
```

### Getting Statistics

If you want to get statistics, we'll use the -c option:

```bash {linenos=table,hl_lines=[1]}
> strace -c uname
Linux
% time     seconds  usecs/call     calls    errors syscall
------ ----------- ----------- --------- --------- ----------------
  -nan    0.000000           0         1           read
  -nan    0.000000           0         1           write
  -nan    0.000000           0         3           open
  -nan    0.000000           0         5           close
  -nan    0.000000           0         4           fstat
  -nan    0.000000           0        10           mmap
  -nan    0.000000           0         3           mprotect
  -nan    0.000000           0         2           munmap
  -nan    0.000000           0         3           brk
  -nan    0.000000           0         1         1 access
  -nan    0.000000           0         1           execve
  -nan    0.000000           0         1           uname
  -nan    0.000000           0         1           arch_prctl
------ ----------- ----------- --------- --------- ----------------
100.00    0.000000                    36         1 total
```

### Detecting Network Problems

If, for example, you only want to work on a network layer, here's the solution:

```bash
strace -e poll,select,connect,recvfrom,sendto nc www.deimos.fr 80
```

## Example

Here's an example of this command with some explanations to help you get started with a 'ls' command:

```c
# Reading an l from the keyboard
read(10, "l"..., 1)                     = 1
# Writing the l to the screen (the one it just read)
write(10, "l"..., 1)                    = 1
# Reading the s
read(10, "s"..., 1)                     = 1
# writing the s
write(10, "\10ls"..., 3)                = 3
# reading the enter key (in C)
read(10, "
"..., 1)                    = 1
write(10, "
"..., 2)                 = 2
alarm(0)                                = 0
ioctl(10, SNDCTL_TMR_STOP or TCSETSW, {B38400 opost isig icanon echo ...}) = 0
time(NULL)                              = 1229629587
pipe([3, 4])                            = 0
gettimeofday({1229629587, 864550}, {0, 0}) = 0
# clone ----> a new process is created, in fact fork() executes the clone system call, the new pid is 4024
clone(child_stack=0, flags=CLONE_CHILD_CLEARTID|CLONE_CHILD_SETTID|SIGCHLD, child_tidptr=0xb7ddd998) = 4024
close(4)                                = 0
read(3, ""..., 1)                       = 0
close(3)                                = 0
rt_sigprocmask(SIG_BLOCK, [CHLD], [CHLD], 8) = 0
rt_sigsuspend([])                       = ? ERESTARTNOHAND (To be restarted)
--- SIGCHLD (Child exited) @ 0 (0) ---
rt_sigprocmask(SIG_BLOCK, ~[RTMIN RT_1], [CHLD], 8) = 0
rt_sigprocmask(SIG_SETMASK, [CHLD], ~[KILL STOP RTMIN RT_1], 8) = 0
# The end of the child process (4024)
wait4(-1, [{WIFEXITED(s) && WEXITSTATUS(s) == 0}], WNOHANG|WSTOPPED, {ru_utime={0, 0}, ru_stime={0, 0}, ...}) = 4024
gettimeofday({1229629587, 867012}, {0, 0}) = 0
ioctl(10, SNDCTL_TMR_TIMEBASE or TCGETS, {B38400 opost isig icanon echo ...}) = 0
ioctl(10, TIOCGPGRP, [4024])            = 0
ioctl(10, TIOCSPGRP, [3982])            = 0
ioctl(10, TIOCGWINSZ, {ws_row=38, ws_col=127, ws_xpixel=1270, ws_ypixel=758}) = 0
wait4(-1, 0xbfe3f48c, WNOHANG|WSTOPPED, 0xbfe3f434) = -1 ECHILD (No child processes)
# It wonders what time it is :-)
time(NULL)                              = 1229629587
ioctl(10, TIOCSPGRP, [3982])            = 0
fstat64(0, {st_mode=S_IFCHR|0620, st_rdev=makedev(136, 3), ...}) = 0
fcntl64(0, F_GETFL)                     = 0x2 (flags O_RDWR)
# It wonders under which UID it is running
getuid32()                              = 1000
# It rewrites the prompt
write(1, "\33]0;phil@philpep.ath.cx ~\7"..., 26) = 26
rt_sigprocmask(SIG_BLOCK, [CHLD], [CHLD], 8) = 0
# It asks again what time it is
time(NULL)                              = 1229629587
rt_sigaction(SIGINT, {0x80a8fd0, [], SA_INTERRUPT}, NULL, 8) = 0
write(10, "\33[1m\33[3m%\33[23m\33[1m\33[0m           "..., 149) = 149
time(NULL)                              = 1229629587
# It opens the file /etc/localtime
stat64("/etc/localtime", {st_mode=S_IFREG|0644, st_size=2945, ...}) = 0
ioctl(10, FIONREAD, [0])                = 0
ioctl(10, TIOCSPGRP, [3982])            = 0
ioctl(10, SNDCTL_TMR_STOP or TCSETSW, {B38400 opost isig -icanon -echo ...}) = 0
write(10, "
\33[0m\33[23m\33[24m\33[J\33[01;30m[\33[01;3"..., 105) = 105
write(10, "\33[K\33[81C  \33[01;30m18/12/08 20:46:"..., 46) = 46
# It waits for a new input
read(10,
Process 3982 detached
```

## Resources

- http://blog.philpep.org/post/Que-font-vos-processus---La-commande-strace
- http://www.hokstad.com/5-simple-ways-to-troubleshoot-using-strace.html
