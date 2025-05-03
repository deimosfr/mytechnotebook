---
weight: 999
url: "/Tester_le_load_sur_ses_serveurs_et_applications_web/"
title: "Testing Load on Servers and Web Applications"
description: "A guide to different tools for load testing web servers and applications, comparing curl-loader, httperf, Siege, and Tsung."
categories: ["Linux", "Apache", "Development"]
date: "2008-09-10T15:57:00+02:00"
lastmod: "2008-09-10T15:57:00+02:00"
tags: ["Servers", "Network", "Development", "httperf", "http_load", "Testing"]
toc: true
---

## Introduction

A good way to see how your Web applications and server will behave under high load is by testing them with a simulated load. We tested several free software tools that do such testing to see which work best for what kinds of sites.

If you leave out the load-testing packages that are no longer maintained, non-free, or fail the installation process in some obscure way, you are left with five candidates: curl-loader, httperf, Siege, Tsung, and Apache JMeter.

## curl-loader

The purpose of [curl-loader](https://curl-loader.sourceforge.net/) is to "deliver a powerful and flexible open source testing solution as a real alternative to Spirent Avalanche and IXIA IxLoad." It relies on the mature and flexible [cURL](https://curl.haxx.se/) library to manage requests, authentication, and sessions.

Building the application is straightforward: download, untar, and make the code inside its directory. I had to add an #include <limits.h> to the file ip_secondary.c to make it build, probably due to some recent changes in the glibc headers. You also need to install the OpenSSL libraries and headers to compile and run curl-loader.

Getting started takes a bit more effort. curl-loader's configuration interface is split into two parts. The first part is a configuration file that contains the parameters for a specific scenario. Its simple format consists of newline-delimited variable assignments of the form VAR=VALUE. You can find a bunch of self-explanatory examples in the directory conf-examples in the curl-loader source tree. Because curl-loader can use multiple IP addresses to realistically simulate requests from different clients, you need to adjust the values for INTERFACE, CLIENTS_NUM_MAX, NETMASK, IP_ADDR_MIN, and IP_ADDR_MAX to suit your network environment. The number of maximum parallel clients must of course be supported by the IP address range you specify.

The second part is the command-line interface of the curl-loader binary. It offers one essential parameter, -f, to which you're expected to pass the location of the scenario file you wish to use. The other arguments let you fine-tune the test run. For example, by default curl-loader will issue its requests from a single thread. While this helps conserve system resources and increase performance, it's advisable to use the -t option to add one thread per additional CPU core you wish to utilize.

A screen-oriented result display shows a synopsis of the test results, updated at regular intervals, and more detailed results are available in the log files curl-loader generates. The file with the .log extension contains information about any errors, while the .ctx file shows per-client statistics and the .txt file shows statistics over time.

If you're looking for a utility similar to curl-loader but written in Python, check out [Pylot](https://www.pylot.org/). It comes with a GUI and uses XML as its configuration format.

## httperf

[httperf](https://www.hpl.hp.com/research/linux/httperf/) is a command-line single-thread load tester developed by Hewlett-Packard Labs.

You set the bulk of httperf's configuration via command-line parameters. Configuration files play an auxiliary role if you wish to specify a session scenario.

A sample invocation for a total of 5,000 connections, each one of which should try to issue 50 requests, looks like this:

```bash
httperf --server=localhost --uri=/ --num-conns=5000 --num-calls=50
```

The first line of output will show arguments that have been assigned their defaults because you haven't specified them:

```bash
httperf --client=0/1 --server=localhost --port=80 --uri=/ \ --send-buffer=4096 --recv-buffer=16384 \ --num-conns=5000 --num-calls=50
```

Unlike curl-loader, httperf doesn't keep you updated about the status of the test run, nor does it write detailed log files. It only shows a summary of the test results at the end of the test run. There's a debugging switch that helps you see what it's currently doing, but it has to be enabled at compile time.

I like the fact that httperf lets you specify all parameters on the command line. This lets you prototype your load test quickly and then put the finished invocation into a shell script. Running different tests sequentially or in parallel is also a no-brainer.

httperf could be a little smarter in its interpretation of its command-line arguments, though. For example, the separation of the target URI into server and path parts seems unnecessary, plus the latter one is specified by the --uri option, which is a misnomer, as valid URIs may contain the server name as well.

Try [http_load](https://www.acme.com/software/http_load/) if you want a simpler tool in the style of httperf.

## Siege

[Siege](https://www.joedog.org/JoeDog/Siege) is similar to httperf in that it can be configured almost fully with command-line arguments. But Siege uses multiple concurrent threads to send its requests, has fewer low-level options than httperf, and is built to work with a list of URLs. It's also easier to use than httperf because its options are named in a more straightforward way.

A Siege run with default parameters can be as simple as

```bash
siege localhost
```

However, this doesn't access the full power of Siege, which is to test a list of URLs in a largely unpredictable way like a real user would. To gather URLs, Siege's author offers the auxiliary program Sproxy. After installing, run it like this:

```bash
sproxy -v -o urls.txt
```

Sproxy will keep the terminal open and list all recorded URLs, plus write them to the file urls.txt for Siege.

Configure your browser to use localhost:9001 as an HTTP proxy. Then it's time to start browsing your site, thereby letting Sproxy record information about the URLs for Siege.

When you have gathered some URLs you would like to test, start the next test run like this:

```bash
siege -v --internet --file=urls.txt
```

The --internet argument instructs Siege to hit the URLs in a random fashion, like a horde of surfers would.

## Tsung

[Tsung](https://tsung.erlang-projects.org/) works in a way similar to using Siege with an URL file, but it offers more elaborate features, such as random user agents, session simulation, and dynamic data feedback. It also performs better by using Erlang's Green Threads.

This comes at a price, though: Tsung doesn't offer the ad-hoc command-line invocation we know from Siege, curl-loader, and httperf. You must either manually create a scenario file in ~/.tsung/tsung.xml or use the recording mode of Tsung, which works like Siege's Sproxy:

* Start the Tsung recording proxy with tsung recorder, then visit the target URLs with localhost:8090 as proxy server.
* Open the newly created session record in ~/.tsung and [edit the details of your scenario](https://tsung.erlang-projects.org/user_manual.html#htoc28).
* Save it as ~/.tsung/tsung.xml.
* Use tsung start to start the test.

Step two, the crafting of the configuration file, is the most difficult part. If you want to use the advanced features of Tsung you will definitely need to get acquainted with its format.

As an additional bonus, Tsung can also put PostgreSQL and Jabber servers under load.

## Conclusion

Each of these tools has its advantages and disadvantages. All are documented well, offer a painless installation, and run reliably.

Here's a possible decision path to find the right tool for your particular job: start out simple with Siege and see if you can get away with its simple performance and feature set. After that, try httperf, which has a slightly expanded feature set and runs faster. If you need to set up more complex scenarios, move on to curl-loader and Tsung, which have the largest feature sets and best performance, but especially Tsung takes time to get used to.

Apache JMeter is the only GUI-based application in the crowd. Its feature set is pretty impressive, offering some unique things like content pre- and postprocessors. You should give it a try if you prefer GUI apps.
