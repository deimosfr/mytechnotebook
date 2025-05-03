---
weight: 999
url: "/Différences_entre_du_et_df/"
title: "Differences between du and df"
description: "Understanding the differences between du and df commands in Linux systems and how to troubleshoot space discrepancies between these two utilities."
categories: ["Linux", "Backup"]
date: "2008-03-19T14:23:00+02:00"
lastmod: "2008-03-19T14:23:00+02:00"
tags: ["Servers", "Linux", "Disk Management", "Troubleshooting"]
toc: true
---

## Issue

You will certainly ask yourself one day: "What is the difference between the du command and the df command?"

According to the man page for du:

```bash
du - estimate file space usage
```

And df displays exactly the space taken on your hard drive. I encountered a case where on a 30 GB partition, I had this:

```bash
du -sh | grep /home : 13 G used
df -h /home : 26 G used
```

The df command is not wrong.

## Solutions

### Check the number and size of blocks on your partition

Indeed, when formatting your partition, you can choose the block size. By default it's 8KB. This means that if you have many files of just a few bytes, they will still take up 8KB each.

The only solution is to backup your data, reformat with a much smaller block size, and then restore the backup.

### Check processes that are writing to this partition

It can happen that processes that were writing to the partition have crashed. In this case, the size remains in virtual memory and that's when problems occur. So analyze your partition:

```bash
lsof /home
```

and check the running processes:

```bash
ps -aux
```
