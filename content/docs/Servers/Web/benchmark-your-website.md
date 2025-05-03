---
weight: 999
url: "/Benchmarker_son_site_web/"
title: "Benchmark Your Website"
description: "How to benchmark your website to test its performance and understand how it behaves under load"
categories: ["Web", "Server", "Performance"]
date: "2013-01-16T07:47:00+02:00"
lastmod: "2013-01-16T07:47:00+02:00"
tags: ["Apache", "Nginx", "Varnish", "Benchmark", "Performance"]
toc: true
---

![Apache](/images/apache_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 2.3 |
| **Operating System** | Ubuntu 12.10 |
| **Website** | [Apache Website](https://www.apache.org) |
| **Last Update** | 16/01/2013 |
| **Others** | Server: Debian 7 |
{{< /table >}}

## Introduction

It's useful to know how many connections your web server can handle. It's important to see how your server behaves under load. That's why it's necessary to benchmark it. We'll see here how to benchmark it, and then we'll look at some differences when using cache servers.

For my tests, I started with an Nginx server, then added a Varnish server in front.

## Installation

For benchmarking, there's the ab command (Apache Benchmark). To install it:

```bash
aptitude install apache2-utils
```

## Running benchmarks

Here's how to use the ab command:

```bash
ab -c <occurrences> -t <time> <website>
```

- occurrences: defines the number of parallel requests
- time: the time (in seconds) that the tests should run
- website: the website to benchmark (use a complete address, like an index.php or index.html)

The tests below are conducted between a server with a 1Gb/s bandwidth and a client with 100Mb/s through the Internet. However, it's important to understand that the source of your tests is extremely important.
Indeed, a local network or loopback will be much more telling in terms of benchmarks than an Internet network.

{{< alert context="info" text="I recommend doing your benchmarks on loopback if you have the possibility" />}}

### Wordpress

Here's the command I used to benchmark my blog:

```bash {linenos=table,hl_lines=[17]}
> ab -c 5 -t 30 http://blog.deimos.fr/index.php
This is ApacheBench, Version 2.3 <$Revision: 655654 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking blog.deimos.fr (be patient)
Completed 5000 requests
Completed 10000 requests
Finished 11072 requests


Server Software:        nginx
Server Hostname:        blog.deimos.fr
Server Port:            80

Document Path:          /index.php
Document Length:        0 bytes

Concurrency Level:      5
Time taken for tests:   30.003 seconds
Complete requests:      11072
Failed requests:        0
Write errors:           0
Non-2xx responses:      11072
Total transferred:      2911936 bytes
HTML transferred:       0 bytes
Requests per second:    369.03 [#/sec] (mean)
Time per request:       13.549 [ms] (mean)
Time per request:       2.710 [ms] (mean, across all concurrent requests)
Transfer rate:          94.78 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        4    6   3.5      6     109
Processing:     4    7  12.6      7     388
Waiting:        4    7  12.6      7     388
Total:          9   14  13.1     13     395

Percentage of the requests served within a certain time (ms)
  50%     13
  66%     13
  75%     14
  80%     14
  90%     14
  95%     14
  98%     15
  99%     16
 100%    395 (longest request)
```

Here, the server is capable of handling 369 requests per second. That's pretty good, but **the CPU on the server side was close to 80%**.

Now, let's add [a cache server]({{< ref "docs/Servers/Web/Caches">}}) like [Varnish]({{< ref "docs/Servers/Web/Caches/varnish_a_website_accelerator.md">}}), then run the benchmarks again:

```bash {linenos=table,hl_lines=[26]}
> ab -c 5 -t 30 http://blog.deimos.fr/index.php
This is ApacheBench, Version 2.3 <$Revision: 655654 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking blog.deimos.fr (be patient)
Completed 5000 requests
Completed 10000 requests
Finished 12541 requests


Server Software:        nginx
Server Hostname:        blog.deimos.fr
Server Port:            80

Document Path:          /index.php
Document Length:        0 bytes

Concurrency Level:      5
Time taken for tests:   30.000 seconds
Complete requests:      12541
Failed requests:        0
Write errors:           0
Non-2xx responses:      12541
Total transferred:      4548468 bytes
HTML transferred:       0 bytes
Requests per second:    418.03 [#/sec] (mean)
Time per request:       11.961 [ms] (mean)
Time per request:       2.392 [ms] (mean, across all concurrent requests)
Transfer rate:          148.06 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        3    5   8.1      5     390
Processing:     4    7   8.6      6     390
Waiting:        4    6   8.6      6     390
Total:          8   12  11.9     11     396

Percentage of the requests served within a certain time (ms)
  50%     11
  66%     12
  75%     12
  80%     12
  90%     13
  95%     13
  98%     14
  99%     14
 100%    396 (longest request)
```

We get a decent gain, the big difference is that **the CPU is below 20% with Varnish**!

### Mediawiki

I also benchmarked the wiki. Again, the CPU was overloaded without a cache server and the information speaks for itself:

```bash {linenos=table,hl_lines=[27]}
> ab -c 5 -t 30 http://wiki.deimos.fr/index.php
This is ApacheBench, Version 2.3 <$Revision: 655654 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking wiki.deimos.fr (be patient)
Finished 2442 requests


Server Software:        nginx
Server Hostname:        wiki.deimos.fr
Server Port:            80

Document Path:          /index.php
Document Length:        0 bytes

Concurrency Level:      5
Time taken for tests:   30.000 seconds
Complete requests:      2442
Failed requests:        0
Write errors:           0
Non-2xx responses:      2442
Total transferred:      1015872 bytes
HTML transferred:       0 bytes
Requests per second:    81.40 [#/sec] (mean)
Time per request:       61.425 [ms] (mean)
Time per request:       12.285 [ms] (mean, across all concurrent requests)
Transfer rate:          33.07 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        3    6   9.2      5     352
Processing:    19   56  32.0     54     410
Waiting:       19   56  32.0     54     409
Total:         26   61  34.8     59     761

Percentage of the requests served within a certain time (ms)
  50%     59
  66%     61
  75%     62
  80%     63
  90%     65
  95%     69
  98%     94
  99%    166
 100%    761 (longest request)
```

And with Varnish:

```bash {linenos=table,hl_lines=[26]}
> ab -c 5 -t 30 http://wiki.deimos.fr/index.php
This is ApacheBench, Version 2.3 <$Revision: 655654 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking wiki.deimos.fr (be patient)
Completed 5000 requests
Completed 10000 requests
Finished 10691 requests


Server Software:        nginx
Server Hostname:        wiki.deimos.fr
Server Port:            80

Document Path:          /index.php
Document Length:        0 bytes

Concurrency Level:      5
Time taken for tests:   30.002 seconds
Complete requests:      10691
Failed requests:        0
Write errors:           0
Non-2xx responses:      10691
Total transferred:      5149681 bytes
HTML transferred:       0 bytes
Requests per second:    356.34 [#/sec] (mean)
Time per request:       14.032 [ms] (mean)
Time per request:       2.806 [ms] (mean, across all concurrent requests)
Transfer rate:          167.62 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        4    6  11.9      6     391
Processing:     5    8  10.9      7     391
Waiting:        5    7  10.9      7     391
Total:         10   14  16.2     13     398

Percentage of the requests served within a certain time (ms)
  50%     13
  66%     13
  75%     14
  80%     14
  90%     14
  95%     15
  98%     16
  99%     16
 100%    398 (longest request)
```
