---
weight: 999
url: "/Planifier_des_tâches/"
title: "Task Scheduling"
description: "Guide on how to schedule tasks on Linux systems using at and cron utilities"
categories: ["Linux"]
date: "2009-12-11T20:54:00+02:00"
lastmod: "2009-12-11T20:54:00+02:00"
tags: ["at", "cron", "task scheduling", "system administration", "automation"]
toc: true
---

## Introduction

Use the at command to automatically execute a job only once at a specified time.

## at

The format for the at command is:

```bash
at -m -q queuename time date
at -r job
at -l
```

The table shows the options you can use to instruct the cron process on how to execute an at job.

{{< table "table-hover table-striped" >}}
| Option | Description |
|--------|-------------|
| -m | Sends mail to the user after the job has finished |
| -r job | Removes a scheduled at job from the queue |
| -q queuename | Specifies a specific queue |
| time | Specifies a time for the command to execute |
| -l | Reports all jobs scheduled for the invoking user |
| date | Specifies an optional date for the command to execute, which is either a month name followed by a day number or a day of the week |
{{< /table >}}

For example, to create an at job to run at 9:00 p.m. to locate and verify the file type of core files from the `/export/home` directory, perform the command:

```bash
# at 9:00 pm
at> find /export/home -name core -exec file {} \; >> /var/tmp/corelog
at> 
