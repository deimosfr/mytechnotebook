---
title: 'Useful print technics'
slug: useful-print-technics/
description: 'Useful print technics'
categories: ['3D Print']
tags: ['3D Print']
date: '2025-12-11T23:53:13+01:00'
---

## Introduction

This page is a collection of useful print technics.

## Bridges without supports

If you have a bridge, and you want to print it without supports, you should have:

- A good quality filament
- A good quality printer
- An enclosure
- A good air flow (Prusa MK4S/CoreOne or similar)

Then you can play with some parameters on Prusa Slicer to get a good result:

- Speed -> Speed for print moves -> Bridges: 25 mm/s
- Advanced -> Flow -> Bridge flow ratio: 1.7

In the end you can have something like this:

![3d print no bridge support](../../static/images/3dprint_bridge_no_support.avif)