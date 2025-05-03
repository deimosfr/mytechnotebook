---
weight: 999
url: "/aspirer-un-site-web/"
title: "Mirroring a Website"
description: "How to mirror an entire website locally using wget command line utility"
categories: ["Web", "Command Line"]
date: "2009-09-20T15:51:00+02:00"
lastmod: "2009-09-20T15:51:00+02:00"
tags: ["wget", "Website", "Mirror", "Download"]
toc: true
---

## Introduction

There are software applications that allow you to mirror websites, but why use them when a simple command can do the same thing?

## Usage

To mirror mozilla.org for example, use this command:

```bash
wget --random-wait -r -p -e robots=off -U mozilla http://www.mozilla.org
```

* -p: Include all files, images, etc.
* -e robots=off: Bypass the robot.txt file
* -U: Specify mozilla as the browser that will do the mirroring
* --random-wait: Allows wget to randomize download times in seconds to avoid blacklists

Other useful parameters:

* --limit-rate=20k: Limits download speed to 20k
* -b: wget continues to run even if you log out (like nohup)
* -o: $HOME/wget_log.txt log file
