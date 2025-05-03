---
weight: 999
url: "/Connaitre_le_temps_d'exécution_d'une_ou_plusieurs_commandes/"
title: "Measuring Execution Time of One or Multiple Commands"
description: "How to measure the execution time of commands in Unix-like systems"
categories: ["Linux", "Shell", "System Administration"]
date: "2009-11-19T06:45:00+02:00"
lastmod: "2009-11-19T06:45:00+02:00"
tags: ["time", "performance", "bash", "unix", "commands"]
toc: true
---

## Introduction

You may need to know exactly the time it takes to execute some commands. The `time` command is perfect for this purpose.

This command is a bash builtin.

## Examples

The last semicolon `;` is important. For example:

```bash
time { rm -rf /folder/bar && mkdir -p /folder/bar ; echo "done" ; }
```

This will give you the time taken to execute the entire block of commands as a single unit.
