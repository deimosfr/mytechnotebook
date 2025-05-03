---
weight: 999
url: "/Connaitre_le_page_size_de_sa_machine/"
title: "How to Check the Page Size of Your Machine"
description: "Methods to determine memory page size in different operating systems for more efficient memory allocation"
categories: ["System Administration", "Development"]
date: "2007-05-09T13:11:00+02:00"
lastmod: "2007-05-09T13:11:00+02:00"
tags: ["memory", "page size", "system programming", "C", "Windows", "Unix"]
toc: true
---

## Determining the Page Size

Most operating systems allow programs to determine what the page size is so that they can allocate memory more efficiently.

### UNIX and POSIX-based Operating Systems

UNIX and POSIX-based systems use the C function sysconf().

Edit a test.c file and paste it:

```c
#include <stdio.h>     // printf(3)
#include <unistd.h>    // sysconf(3)

int
main(void)
{
        printf("The page size for this system is %ld bytes\n", sysconf(_SC_PAGESIZE)); //_SC_PAGE_SIZE is OK too.
        return 0;
}
```

Then compile it with gcc:

```bash
gcc -o test test.c
```

Now launch it:

```bash
# ./test
The page size for this system is 4096 bytes
```

### Win32-based Operating Systems (Windows 9x, NT, ReactOS)

Win32-based operating system use the C function GetSystemInfo() function in kernel32.dll

```c
#include <windows.h>

SYSTEM_INFO si;

GetSystemInfo(&si);
printf("The page size for this system is %u bytes\n", si.dwPageSize);
```
