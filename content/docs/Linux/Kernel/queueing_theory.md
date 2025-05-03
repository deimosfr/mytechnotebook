---
weight: 999
url: "/Théorie_des_files_d\\'attentes/"
title: "Queueing Theory"
description: "An overview of queueing theory, Little's Law, wait times, completion rates, and predicting system limits"
categories: ["Linux", "Network", "Development"]
date: "2013-07-26T13:59:00+02:00"
lastmod: "2013-07-26T13:59:00+02:00"
tags:
  [
    "queueing theory",
    "performance",
    "system optimization",
    "throughput",
    "network",
    "linux",
    "waiting times",
  ]
toc: true
---

## Introduction

[Queueing theory](https://fr.wikipedia.org/wiki/Th%C3%A9orie_des_files_d%27attente) is a mathematical theory in the field of probability, which studies optimal solutions for managing waiting lines, or queues. A queue is necessary and will create itself if not anticipated, in all cases where supply is less than demand, even temporarily. It can apply to different situations: management of aircraft taking off or landing, waiting of customers and administrators at counters, or storage of computer programs before processing. This field of research, born in 1917 from the work of Danish engineer Erlang on the management of Copenhagen telephone networks between 1909 and 1920, studies in particular arrival systems in a queue, the different priorities of each new arrival, as well as statistical modeling of execution times. It is thanks to the contributions of mathematicians Khintchine, Palm, Kendall, Pollaczek and Kolmogorov that the theory really developed.

## Little's Law

Little mathematically demonstrated what Erlang had said years earlier. Here is the law:

> L = A x W

- L: The queue length that corresponds to an average of requests waiting in the system
- A: This is the rate at which requests enter the system (x/second)
- W: Average time to satisfy a request

## Queue Size

Requests are stored in memory. Therefore L can be a tunable read/write cache or a read-only cache for measurements. There are algorithms to manage queue priorities:

- Two queues (1 for urgent requests, 1 for normal requests). Requests are processed according to their urgency.
- A single queue that will process urgent requests first. Normal requests will only be processed if there are no more urgent ones.

## Queue Size and Waiting Time

L tends to vary directly with W. Simple example:

> **100** requests = **100** requests/second _ **1s**  
> **200** requests = **100** requests/second _ **2s**  
> **800** requests = **100** requests/second \* **8s**

We can therefore predict the waiting time performance by restricting the queue size.
Similarly, we can restrict the waiting time to optimize the queue size (in memory).

## Waiting Time

The waiting time includes:

> (W = Q + S)

- Queue time: waiting time for a resource to become available
- Service time: time for a resource to execute a request

The goal is to reduce both queue time and service time. For example, a waiting time corresponds to a web page loading. Its execution time can take a certain amount of time which can be multiplied according to the number of parallel requests.

If we look from a more system-oriented perspective:

> W = Q + (T(system) + T(user))

- System time: time for a kernel task
- User time: time for a user task

You can get more information in the man page of the time command or in top. The goal is therefore:

- To reduce system time (but will block operations for users)
- Spend as much time as necessary on user tasks (but no more)

The service time corresponds to the time a process takes to execute a task inserted in the queue until the end of its execution. This type of service can be calculated as utilization / throughput.
The average service time is defined by the amount of occupancy of a resource per request: S = Occupancy time / completions per second

## User Time Algorithms

There are asymptotic complexities that help understand the necessary user time to execute x requests:

- Intractable: O(2^n)
- Polynomial: O(n^2)
- Linear: O(50n)
- Constant: O(1)
- Logarithmic: O(500 log n)

Note: O describes the order of growth rate at a time of execution or memory usage of an algorithm when its inputs grow.

When changing algorithms, performance can be improved by reducing waiting time and increasing throughput (number of actions).

## Profiling with the time Command

The time command allows you to know the execution time of a command:

```bash
> export TIME="\n%e %S %U"
> time tar -czf deimos.tgz /home/deimos
real    0m0.017s
user    0m0.001s
sys     0m0.001s
```

- real: this is the actual time the application takes to execute
- user: this is the CPU time needed to execute its instructions
- sys: this is the time that the CPU takes for kernel calls or waiting for I/O

To get the queue time: Q = W - (Tsys + Tuser)

## Completion Rate

The bandwidth (theoretical) corresponds to a data size **that can** be transferred in a certain amount of time.
The bandwidth (real) corresponds to a data size **that is** transferred in a certain amount of time.
The throughput corresponds to a size of **usable data that is** transferred in a certain amount of time.
The overhead is the difference between the real bandwidth and the throughput. It's the surplus of real data transferred.
For data transferred at a certain bandwidth, the throughput will always be less than the bandwidth.

The completion rate is the rate at which work is performed at a certain throughput or bandwidth:

> x = work done / observation time

Bandwidth is generally predictable, but overhead can be reduced to speed up processing time. For example, on a 100Mb/s network, we can measure 80Mb/s of data at the application layer (OSI). Because in fact, we have 20Mb/s of overhead (20 + 80 = 100).
If we reduce the overhead to 10, we can increase the throughput to 90Mb/s.

> B (bandwidth) = X (transfers) + O (overhead)

However, it should be known that reducing overhead can inflict undesirable behavior. For example, if we decide to reduce overhead by choosing UDP instead of TCP for an application (since UDP has less overhead than TCP). If there is packet loss on the network, the application may request retransmission of the request. As a result, the application might generate more overhead than with TCP sessions.
Depending on the applications, data does not necessarily need to be retransmitted (such as online music sites). After a certain number of packets not transmitted, the client can give an error message for bad data transmission.

## Arrival Rate and Completion Rate

The completion rate is the rate at which work is performed and is related to throughput, as well as bandwidth (depending on the context).

> A (number of arrivals) = C (number of completions) / T (observation time)

- T: the observation time corresponds to the time allotted for the execution of events in a given time (e.g., in seconds)
- A: The number of arrivals observed during the observation period (e.g., packets/s)

The number of arrivals (A) that were made during an observation time (T) are referenced as completed (C).
The occupancy time is the observation time \* the utilization time.

## Finding the Best Observation Time

It is important to use these kinds of tools during non-zero load times. Making an observation during too short a time will not allow for satisfactory results. It is therefore advised to measure the observation time in the following way:

- Collect data as frequently as possible (ideally 1 second)
- Compare the measurement values according to Little's Law
- Is the comparison good?
  - yes: the observation time is valid and L = C x W
  - no: redo the capture with a longer time interval

We will now calculate the queue length and compare it to the collected data.
Here is an example of capture:

```bash {linenos=table,hl_lines=[3,10,17,24,31,38]}
> dd if=/dev/zero of=/tmp/bigfile bs=1M count=1024 oflag=direct & sar -d -p 1 5 > sar_result ; pkill dd && rm -f /tmp/bigfile ; cat sar_result
23:45:05          DEV       tps  rd_sec/s  wr_sec/s  avgrq-sz  avgqu-sz     await     svctm     %util
23:45:06          sda    227,00   9872,00  16288,00    115,24      6,85     28,58      4,39     99,60
23:45:06     dev254-0      1,00      0,00      8,00      8,00      2,23   2836,00    284,00     28,40
23:45:06     dev254-1      0,00      0,00      0,00      0,00      0,00      0,00      0,00      0,00
23:45:06     dev254-2    206,00  10056,00    424,00     50,87      7,17     34,00      4,82     99,20
23:45:06     dev254-3   2541,00     40,00  20288,00      8,00    100,64     33,19      0,12     30,00

23:45:06          DEV       tps  rd_sec/s  wr_sec/s  avgrq-sz  avgqu-sz     await     svctm     %util
23:45:07          sda    114,85   2811,88  12110,89    129,93      9,98     89,90      8,62     99,01
23:45:07     dev254-0      0,00      0,00      0,00      0,00      0,00      0,00      0,00      0,00
23:45:07     dev254-1      0,99      0,00      7,92      8,00      0,07     68,00     68,00      6,73
23:45:07     dev254-2     84,16   2471,29      0,00     29,36      3,86     47,25     11,76     99,01
23:45:07     dev254-3     21,78    253,47    166,34     19,27     16,84   1514,73     37,09     80,79

23:45:07          DEV       tps  rd_sec/s  wr_sec/s  avgrq-sz  avgqu-sz     await     svctm     %util
23:45:08          sda    207,00   4672,00   1304,00     28,87      6,09     29,72      4,83    100,00
23:45:08     dev254-0      0,00      0,00      0,00      0,00      0,00      0,00      0,00      0,00
23:45:08     dev254-1    162,00      0,00   1296,00      8,00     59,80    369,14      4,02     65,20
23:45:08     dev254-2    185,00   4216,00      0,00     22,79      5,07     27,70      5,38     99,60
23:45:08     dev254-3      5,00    288,00      0,00     57,60      0,18     36,00     36,00     18,00

23:45:08          DEV       tps  rd_sec/s  wr_sec/s  avgrq-sz  avgqu-sz     await     svctm     %util
23:45:09          sda    173,00   3208,00      0,00     18,54      3,54     19,93      5,78    100,00
23:45:09     dev254-0      0,00      0,00      0,00      0,00      0,00      0,00      0,00      0,00
23:45:09     dev254-1      0,00      0,00      0,00      0,00      0,00      0,00      0,00      0,00
23:45:09     dev254-2    171,00   2744,00      0,00     16,05      3,30     18,76      5,85    100,00
23:45:09     dev254-3      3,00    520,00      0,00    173,33      0,16     54,67     54,67     16,40

23:45:09          DEV       tps  rd_sec/s  wr_sec/s  avgrq-sz  avgqu-sz     await     svctm     %util
23:45:10          sda    155,56   3797,98   1470,71     33,87      3,29     21,71      6,36     98,99
23:45:10     dev254-0      1,01      8,08      0,00      8,00      0,06     60,00     60,00      6,06
23:45:10     dev254-1     84,85      0,00    678,79      8,00      5,16     60,86      2,81     23,84
23:45:10     dev254-2    242,42   3498,99    791,92     17,70      6,40     26,77      4,02     97,37
23:45:10     dev254-3      1,01    242,42      0,00    240,00      0,08     76,00     76,00      7,68

23:45:10          DEV       tps  rd_sec/s  wr_sec/s  avgrq-sz  avgqu-sz     await     svctm     %util
23:45:11          sda     74,75   1252,53  13292,93    194,59      0,95     12,81      5,35     40,00
23:45:11     dev254-0     43,43      8,08    339,39      8,00      0,40      9,12      3,44     14,95
23:45:11     dev254-1     90,91      0,00    727,27      8,00      3,03     33,29      0,93      8,48
23:45:11     dev254-2     27,27   1204,04      0,00     44,15      0,28     10,81      8,30     22,63
23:45:11     dev254-3   1530,30      8,08  12234,34      8,00     25,25     16,36      0,16     25,05
...
```

Based on the columns that sar gives us, here is the formula:

> (rd_sec/s + wr_sec/s) \* await / 1000 = requests/s

So if we take the information above, we can easily calculate like this (on sda here):

```bash {linenos=table,hl_lines=[1]}
> grep sda sar_result | awk '{ printf("%s%s%s%s", "("$4" + ", $5") * ", $8" / 1000) = ", (($4+$5)*$8)/1000"\n") }'
(9872,00 + 16288,00) * 28,58 / 1000) = 732.48
(2811,88 + 12110,89) * 89,90 / 1000) = 1327.97
(4672,00 + 1304,00) * 29,72 / 1000) = 173.304
(3208,00 + 0,00) * 19,93 / 1000) = 60.952
(3797,98 + 1470,71) * 21,71 / 1000) = 110.607
(1252,53 + 13292,93) * 12,81 / 1000) = 174.528
(1837,62 + 1900,99) * 13,89 / 1000) = 48.581
(2407,92 + 22922,77) * 42,17 / 1000) = 1063.82
(290,91 + 371,72) * 4,00 / 1000) = 2.644
(1096,00 + 96,00) * 7,64 / 1000) = 8.344
(118,10 + 3126,40) * 3,74 / 1000) = 9.732
```

It is interesting to let these kinds of tools run for a while to profile an application on a server, for example, and realize during a high load period the number of requests available.

## Predicting System Limits

It's important to know that a saturated resource is the bottleneck for the entire system! You must therefore ensure that there is no funnel effect on your resources.

Reducing the number of accesses or improving bandwidth helps solve the problem, whether it's for CPU, RAM, Network, or anything else.

It is therefore important to know both the arrival rate and the necessary bandwidth to properly configure your hardware.
