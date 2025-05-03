---
weight: 999
url: "/Rediriger_l\\'output_d\\'un_service_vers_un_fichier/"
title: "Redirect a Service Output to a File"
description: "How to redirect a program's output to a file using GDB."
categories: ["Linux"]
date: "2011-05-01T20:08:00+02:00"
lastmod: "2011-05-01T20:08:00+02:00"
tags: ["Linux", "GDB", "Debugging", "Services"]
toc: true
---

## Introduction

It can be useful to redirect a program's output to a file. Here's how to do it.

## Usage

```bash
yes 'Y'|gdb -ex 'p close(1)' -ex 'p creat("/tmp/output.txt",0600)' -ex 'q' -p pid
```

This command uses the GDB debugger to attach to a running process and reassign the file handle to a file.

The two commands executed in GDB are:

```
p close(1) which closes STDOUT
```

and

```
p creat("/tmp/filename",0600)
```

which creates a file and opens it for output to which the process will be assigned.

Sequentially, this command closes the STDOUT file handle, creates a new output file, and captures the output to this file.
