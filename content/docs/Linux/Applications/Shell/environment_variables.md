---
weight: 999
url: "/Env_\\:_variables_d'environnements/"
title: "Environment Variables"
description: "A guide on how to manage environment variables in Linux systems."
categories: ["Linux"]
date: "2007-02-21T15:03:00+02:00"
lastmod: "2007-02-21T15:03:00+02:00"
tags: ["Servers", "Linux", "Mac OS X", "Windows"]
toc: true
---

For environment variables, when you don't use them often, it's not always easy to remember the commands.

To display your PATH:

```bash
echo $PATH
```

To display all environment variables:

```bash
env
```

To add something to your PATH:

```bash
PATH=$PATH:/path/to/add
```

To add a new environment variable:

```bash
MYVAR=/toto
export $MYVAR
```

Check with the `env` command and it works! :-)
