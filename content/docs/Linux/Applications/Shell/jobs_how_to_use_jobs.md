---
weight: 999
url: "/Jobs_\\:_Utilisation_des_jobs/"
title: "Jobs: How to Use Jobs"
description: "Learn how to use Linux jobs functionality to run multiple tasks in parallel without needing multiple shells."
categories: ["Linux"]
date: "2013-05-06T20:02:00+02:00"
lastmod: "2013-05-06T20:02:00+02:00"
tags: ["Linux", "Shell", "Terminal", "Command Line", "Background Jobs"]
toc: true
---

## Introduction

Jobs allow you to have multiple tasks running in parallel. The advantage is that you don't need to open multiple shells to launch multiple applications.

## Usage

The `jobs` command allows you to see exactly what is running in the background:

```bash
jobs
```

```
[1]  + running    tail -f /var/log/syslog
```

Here we can see that there is a tail -f running.

If you want to launch a command so that it becomes a job (running in the background), start it like this:

```bash
tail -f /var/log/syslog &
```

- If you forgot the '&' symbol at the end of your command, no worries, there is a way to fix it. **Press 'Ctrl+Z' (^Z) to pause the current command, then type 'bg' which means background.**

- You can check the status of your command with the `jobs` command. Then, if you want to bring back the command you just put in the background, simply use the `fg` command which means foreground.

- If you want to exit your shell, you will lose all your ongoing jobs. To prevent this from happening, you need to launch a nohup like this:

```bash
nohup my_command &
```

Again, if you're forgetful and forgot the nohup, there's a solution. After launching your command, you need to:

- Press Ctrl+Z (^Z)
- Type: `bg`
- Then type: `disown`

This will do the equivalent of the previous command.

- To kill a job, for example the first one:

```bash
kill %1
```

To see what's happening with background processes:

```bash
lsof -p$!
```
