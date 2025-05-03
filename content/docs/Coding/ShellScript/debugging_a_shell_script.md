---
weight: 999
url: "/Debugger_un_script_shell/"
title: "Debugging a Shell Script"
description: "How to debug shell scripts using built-in options like -v and -x to trace execution and understand script behavior."
categories: ["Linux"]
date: "2008-03-05T09:49:00+02:00"
lastmod: "2008-03-05T09:49:00+02:00"
tags: ["Linux", "Bash", "Scripting", "Debugging", "Shell"]
toc: true
---

Shell scripts are often criticized for not having an integrated debugger. But this is false!

When programming in bash, there are command line options to see what is being read and then executed in a script.

## Example

Let's say we have the script MyScript.sh:

```bash
#!/bin/sh
touch unFichier
if [ -f ./unFichier ]; then
  rm ./unFichier
fi
```

If we execute it this way:

```bash
/bin/bash -v -x ./MyScript.sh
```

We'll get output like this:

```bash
#!/bin/sh
touch unFichier
+ touch unFichier
if [ -f ./unFichier ]; then
  rm ./unFichier
fi
+ '['-f ./unFichier ']'
+ rm ./unFichier
```

The normal lines are the lines and blocks that are read, while those with a + in front are the ones being executed.
