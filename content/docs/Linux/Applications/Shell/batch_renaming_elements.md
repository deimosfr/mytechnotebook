---
weight: 999
url: "/Renommer_en_masse_des_éléments/"
title: "Batch Renaming Elements"
description: "Tips and techniques for batch renaming files and elements in a system"
categories: ["Linux"]
date: "2009-09-20T15:03:00+02:00"
lastmod: "2009-09-20T15:03:00+02:00"
tags: ["command line", "batch processing", "file management"]
toc: true
---

## Introduction

Here are some quick tips (albeit a bit crude) to make bulk modifications to elements.

## Rename All Uppercase Elements to Lowercase

```bash
find . -type f
```

## Make Elements Sequential

```bash
for i in `seq -f %03g 5 50 111`; do echo $i ; done
```

The seq command will give:

```
foo01
foo04
foo07
foo10
```
