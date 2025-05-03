---
weight: 999
url: "/Faire_freezer_une_machine_par_overload_CPU/"
title: "Freeze a Machine by CPU Overload"
description: "How to freeze a machine by causing CPU overload, useful for testing system stability and resource limits."
categories: ["Linux", "System Administration", "Security"]
date: "2009-09-20T15:41:00+02:00"
lastmod: "2009-09-20T15:41:00+02:00"
tags: ["cpu", "performance", "bash", "system", "fork bomb"]
toc: true
---

## Introduction

Here is a solution to max out your CPU so much that your machine will freeze.

## Usage

Simply run this in a terminal:

```bash
:(){ :|:& };:
```
