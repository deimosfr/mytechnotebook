---
weight: 999
url: "/Sosreport\\:\_générez\_et\_analysez\_des\_rapports\_de\_machine/"
title: "SoSReport: Generate and analyze machine reports"
description: "Learn how to use sosreport to generate comprehensive system reports and analyze them with sxconsole."
categories: ["Red Hat", "CentOS", "Diagnostics", "System Administration"]
date: "2013-03-15T15:22:00+02:00"
lastmod: "2013-03-15T15:22:00+02:00"
tags:
  [
    "SosReport",
    "Red Hat",
    "CentOS",
    "System Monitoring",
    "Diagnostics",
    "Troubleshooting",
    "System Reports",
  ]
toc: true
---

![SOS Report](/images/red_hat_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 2.2 |
| **Operating System** | Red Hat 6.4 |
| **Website** | [Sosreport Website](https://github.com/sosreport/sosreport) |
| **Last Update** | 15/03/2013 |
| **Others** | CentOS 6.4 |
{{< /table >}}

## Introduction

If you've ever contacted Red Hat support, you probably know that they systematically request a sosreport. It's a very practical tool for the support team as it collects command results, log files, configuration files, and other information often essential for diagnosing a problem. Personally, I've installed it on CentOS and it works very well.

Recently, sosreport has also been ported to Debian/Ubuntu, making it a very interesting tool when you have a heterogeneous Linux environment. It's possible to get a quick summary of the provided information via sxconsole. The advantage of sosreport is that it collects essential information without consuming a lot of resources. Therefore, it can be executed in production environments without risking any slowdowns (depending on the production environment, of course).

## Installation

Installing it is simple, as it's part of the packages available in the base repository:

```bash
yum install sos
```

If you want to install sxconsole to get a summary:

```bash
yum install sx
```

## Utilization

### sosreport

Using it is simple. On the machine where the problem is occurring, run:

```bash {linenos=table,hl_lines=[17,18,26]}
> sosreport --report

sosreport (version 2.2)

This utility will collect some detailed  information about the
hardware and setup of your CentOS system.
The information is collected and an archive is  packaged under
/tmp, which you can send to a support representative.
CentOS will use this information for diagnostic purposes ONLY
and it will be considered confidential information.

This process may take a while to complete.
No changes will be made to your system.

Press ENTER to continue or CTRL-C to quit.

Please enter your first name (if you have more than one) and your last name [localhost]:
Please enter the case number that you are generating this report for [None]: 123

Launching plugins. Please wait...

  Completed [45/45] ...
Creating compressed archive...

Your sosreport has been generated and saved in:
  /tmp/sosreport-localhost.123-20130315154819-d582.tar.xz

The md5sum is: 4136a18ba7a5e7e2151203b71bc3d582

Please send this file to your support representative.
```

The report option will generate an HTML file that can be used more easily in some cases. For the instructions to enter:

- localhost: ideally, enter the machine name, it's easier to find, especially if you generate reports for multiple machines
- Case number: this is the ticket number. If it's for your personal needs, you can enter anything
- Finally, in /tmp you'll find an archive with all the collected information

You can then easily extract them:

```bash
tar -xaf /tmp/sosreport-localhost.123-20130315154819-d582.tar.xz
```

### sxconsole

sxconsole will generate a summary from a sysreport archive. Here's the usage syntax:

```bash
sxconsole <ticket_number> -E -d -r <sosreport_archive> -R ~/tmp/
```

- ticket_number: this is the ticket number that was inserted when creating the sosreport. You can find this number in the name of the sosreport archive
- -E: enables all modules (cluster, storage...)
- -d: debug mode
- -r: the sosreport of a machine. If you have multiple machines, you can execute multiple -r with one machine at a time
- -R: this is a **temporary storage directory that must be created beforehand**

Here's an example using the syntax described above:

```bash
sxconsole 123 -E -d -r /tmp/sosreport-localhost.123-20130315154819-d582.tar.xz -R ~/tmp/
```

## References

https://github.com/sosreport/sosreport
