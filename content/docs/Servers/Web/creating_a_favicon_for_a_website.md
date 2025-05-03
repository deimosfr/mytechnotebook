---
weight: 999
url: "/Créer_une_favicon_pour_un_site_web/"
title: "Creating a Favicon for a Website"
description: "How to create and convert images into favicons for websites"
categories: ["Linux"]
date: "2009-11-30T20:31:00+02:00"
lastmod: "2009-11-30T20:31:00+02:00"
tags: ["Servers", "Linux", "Web Development", "Image Conversion"]
toc: true
---

## Introduction

A favicon is the small icon you see at the top of your browser, next to the URL. If you have a 16x16 size image, you can use the following method to convert it to the right format.

## Converting an Image

```bash
convert -colors 256 -resize 16x16 face.jpg face.ppm && ppmtowinicon -output favicon.ico face.ppm
```
