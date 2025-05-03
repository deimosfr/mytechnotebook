---
weight: 999
url: "/GS_\\:_Assembler_plusieurs_PDF_pour_n'en_fait_qu'un/"
title: "GS: Merge Multiple PDFs Into One"
description: "Learn how to use GhostScript (GS) to merge multiple PDF files into a single document."
categories: ["Linux", "Tools", "Documents"]
date: "2009-11-28T16:10:00+02:00"
lastmod: "2009-11-28T16:10:00+02:00"
tags: ["PDF", "GhostScript", "document management"]
toc: true
---

## Introduction

You may need to merge multiple PDF files to get only one. Here is the solution.

## Usage

```bash
gs -q -sPAPERSIZE=letter -dNOPAUSE -dBATCH -sDEVICE=pdfwrite -sOutputFile=out.pdf `ls *.pdf`
```

This command merges all PDF files in the current directory into one PDF file (the out.pdf file).
