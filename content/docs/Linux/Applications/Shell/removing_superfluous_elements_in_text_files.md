---
weight: 999
url: "/Suppression_des_éléments_superflus_dans_un_fichier_texte/"
title: "Removing Superfluous Elements in Text Files"
description: "How to clean text files by removing unwanted characters like ^M line endings and empty lines using various command-line tools."
categories: ["Linux"]
date: "2009-09-20T15:07:00+02:00"
lastmod: "2009-09-20T15:07:00+02:00"
tags: ["Text Processing", "Command Line", "sed", "perl"]
toc: true
---

## Introduction

The title of this documentation is admittedly ambiguous, but it's difficult to be more specific. When you're developing or writing text, unwanted characters may appear in your files, such as tabulations on empty lines, "^M" characters at the end of each line, or similar artifacts depending on the editor you've chosen.

Here's how to save precious bytes :-p and most importantly, make your documents "clean".

## Removing ^M Characters at the End of Lines

Have you used Windows Wordpad? Too bad, why not use a real OS? :-p Use this command on your file to clean it:

```bash
perl -pi -e 's:^V^M::g' my_dirty_file > my_clean_file
```

## Removing Empty Lines

```bash
sed '/./,$!d' my_dirty_file > my_clean_file
```

or

```bash
tr -d "\n" < file1 > file2
```
