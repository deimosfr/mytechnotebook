---
weight: 999
url: "/Shell_\\:_renommer_en_masse_avec_compteur/"
title: "Shell: Batch Renaming with Counter"
description: "How to rename multiple files in batch with a counter using shell scripting."
categories: ["Linux"]
date: "2006-08-05T13:09:00+02:00"
lastmod: "2006-08-05T13:09:00+02:00"
tags: ["bash", "shell", "file management", "scripting"]
toc: true
---

Here's a script that allows you to rename JPG files while adding a counter:

```bash
export j=0 # export is only useful if you're working in interactive mode (not in a script)
for i in *.JPG ; do
    mv $i `echo $i | sed s/^/$j\ -\ /`
    j=$((j+1))
done
```

Here's a faster method for renaming files:

Instead of typing:

```bash
mv my_file.txt my_file.that_i_want_to_backup
```

You can simply do:

```bash
mv my_file.{txt,that_i_want_to_backup}
```
