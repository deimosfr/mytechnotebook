---
weight: 999
url: "/find-and-du-locate-large-files-and-directories/"
title: "Locate Large Files and Directories"
description: "How to use find and du commands to identify large files and directories consuming disk space"
categories: ["Linux", "System Administration", "Command Line"]
date: "2008-08-31T10:59:00+02:00"
lastmod: "2008-08-31T10:59:00+02:00"
tags: ["Find", "Du", "Disk Usage", "Files", "Directories"]
toc: true
---

## Finding Large Directories

To find large directories, use the du command and sort the output.

For example, to output the 10 largest directories in /var, sorted in ascending size order, use the following command:

```bash
du -ko /var
```

To avoid crossing file system boundaries, that is, to see the directory usage in / but not in the other mounted files systems (/var, /opt, and so on), add the d option to the du command:

```bash
du -kod /var
```

Another super useful tool is ncdu. An ncurses disk usage viewer that provides a more user-friendly interface for navigating through directories and their sizes.

```
ncdu /
```

## Finding Large Files

To find large files, use the find command and sort the output.

- Example 1: To find all plain files (not block, character, symbolic links, and so on) in a file system larger than 200,000 512-byte blocks (approximately 100 Mbytes) and sort on field 7 (file size) while numerically ignoring leading blanks, do this:

```bash
find / -size +200000 -type f -ls
```

- Example 2: To find all plain files (not block, character, symbolic links, and so on) in a /var file system larger than 1,000 512-byte blocks (approximately 500 Kbytes) and sort on field 7 (file size) while numerically ignoring leading blanks, do this:

```bash
find /var -size +1000 -type f -ls
```
