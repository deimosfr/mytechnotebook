---
weight: 999
url: "/Date_\\:_utilisation_avancée_de_la_commande_date/"
title: "Date: Advanced Usage of the Date Command"
description: "Guide to advanced date calculations and formatting using the Linux date command, including examples for calculating past dates and creating ISO 8601 timestamps."
categories: ["Linux"]
date: "2009-11-27T20:18:00+02:00"
lastmod: "2009-11-27T20:18:00+02:00"
tags: ["Linux", "Command Line", "Time Management"]
toc: true
---

## Introduction

The date command can be very useful for advanced date calculations.

## Usage

* Calculates the date 2 weeks ago from Saturday in the specified format:

```bash
date -d '2 weeks ago Saturday' +%Y-%m-%d
```

* Unix alias for date command that lets you create timestamps in ISO 8601 format:

```bash
alias timestamp='date "+%Y%m%dT%H%M%S"'
```
