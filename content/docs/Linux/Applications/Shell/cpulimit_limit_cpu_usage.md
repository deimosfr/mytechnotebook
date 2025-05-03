---
weight: 999
url: "/Cpulimit\\:_limit_CPU_usage/"
title: "Cpulimit: Limit CPU Usage"
description: "Learn how to limit CPU usage for processes in Linux using Cpulimit tool. This guide covers installation, usage examples and best practices."
categories: ["Linux", "Debian"]
date: "2013-05-06T13:13:00+02:00"
lastmod: "2013-05-06T13:13:00+02:00"
tags:
  - "CPU"
  - "Process Management"
  - "System Administration"
  - "Performance"
  - "Linux"
toc: true
---

|                  |                                                           |
| ---------------- | --------------------------------------------------------- |
| Software version | 1.7                                                       |
| Operating System | Debian 7                                                  |
| Website          | [Cpulimit Website](https://github.com/opsengine/cpulimit) |
| Last Update      | 06/05/2013                                                |

## Introduction

Cpulimit[^1] is a tool which limits the CPU usage of a process (expressed in percentage, not in CPU time). It is useful to control batch jobs, when you don't want them to eat too many CPU cycles. The goal is prevent a process from running for more than a specified time ratio. It does not change the nice value or other scheduling priority settings, but the real CPU usage. Also, it is able to adapt itself to the overall system load, dynamically and quickly.

The control of the used cpu amount is done sending SIGSTOP and SIGCONT POSIX signals to processes.
All the children processes and threads of the specified process will share the same percent of CPU.

## Installation

```bash
aptitude install cpulimit
```

## Usage

It's really easy to use it. If you want for example limit a PID to 40% of CPU:

```bash
cpulimit -p 1234 -l 40
```

This will limit 40% of one core on the PID 1234. You can also do it directly on a binary:

```bash
cpulimit -l 40 app
```

Other case could be found in the man. But the idea is clear...limiting the CPU of one app.

## References

[^1]: https://github.com/opsengine/cpulimit
