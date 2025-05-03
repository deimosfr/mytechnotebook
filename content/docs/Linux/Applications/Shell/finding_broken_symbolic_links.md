---
weight: 999
url: "/Trouver_des_liens_symboliques_cassés/"
title: "Finding Broken Symbolic Links"
description: "How to find and clean up broken symbolic links on Unix systems"
categories: ["Linux"]
date: "2009-09-20T16:00:00+02:00"
lastmod: "2009-09-20T16:00:00+02:00"
tags: ["Linux", "Servers", "Maintenance"]
toc: true
---

## Introduction

After a few years on heavily used servers, especially those with many users, you may end up with broken symbolic links. Here's a solution to find them all at once.

## Usage

To search for all broken links:

```bash
find . -type l ! -exec test -e {} \; -print
```

And to delete them in the same operation:

```bash
find . -type l ! -exec test -e {} \; -print0 | xargs -0 rm
```
