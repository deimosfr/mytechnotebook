---
title: "Manipulating Images via Command Line"
slug: manipulating-images-via-command-line/
description: "How to manipulate images using command line tools with ImageMagick"
categories: 
  - Linux
date: "2009-09-19T21:17:00+02:00"
lastmod: "2009-09-19T21:17:00+02:00"
tags:
  - ImageMagick
  - CLI
  - Images
---

## Introduction

It can be very useful to manipulate images through the command line. For example, to automate tasks.

## Installation

For the following manipulations, we'll need ImageMagick:

```bash
apt-get install imagemagick
```

## Usage

### Getting Information About an Image

```bash
identify -format "%wx%h" /path/to/image.jpg
```
