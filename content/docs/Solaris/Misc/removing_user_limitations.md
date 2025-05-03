---
weight: 999
url: "/Enlever_les_limitation_utilisateurs/"
title: "Removing User Limitations"
description: "Guide on how to remove JVM memory limitations for users on Solaris systems"
categories: ["Linux", "Solaris"]
date: "2008-10-21T17:23:00+02:00"
lastmod: "2008-10-21T17:23:00+02:00"
tags: ["Servers", "Operating Systems", "Memory Management", "Java", "Solaris"]
toc: true
---

## Introduction

At work, I experienced an issue with JVM memory where Xmx and Xms were set to 4.5GB, and part of the memory was going to swap until a Full GC triggered, at which point the reserved memory (RSS) that had gone to swap returned to RAM.

The warning messages I received were:

```
pages faults
```

I searched quite a bit before finding a solution.

## Removing User Limitations

First, for the JVM, you need to add an option at launch:

```
-XX+UseISM
```

This activates the use of Intimate Shared Memory.

However, users are limited under Solaris in terms of allocatable memory, so to increase the limit:

```bash
projadd -U qa user.qa
projmod -sK "project.max-shm-memory=(privileged,64G,deny)" user.pmavro
```

Then log back in as user pmavro and verify that everything worked correctly:

```bash
$ prctl -n project.max-shm-memory -i project user.pmavro
project: 100: user.pmavro
NAME    PRIVILEGE       VALUE    FLAG   ACTION                       RECIPIENT
project.max-shm-memory
       privileged      64,0GB      -   deny                                 -
       system          16,0EB    max   deny                                 -
```

Once this is done, the RES memory in the top command equals the maximum right from startup, and there are no more page fault issues. :-)
