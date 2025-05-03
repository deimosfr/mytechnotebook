---
weight: 999
url: "/Memcached_\\:_Mise_en_place_d'un_serveur_de_cache_pour_Apache/"
title: "Memcached: Setting up a Cache Server for Apache"
description: "How to set up Memcached as a caching server for Apache to accelerate client requests and improve performance."
categories: ["Linux", "Apache"]
date: "2009-11-27T20:20:00+02:00"
lastmod: "2009-11-27T20:20:00+02:00"
tags: ["memcached", "cache", "apache", "performance", "php"]
toc: true
---

## Introduction

Memcached is a caching server that helps accelerate client requests.

## Setup

Here's documentation for setting up a memcached cache server:

[Installing memcached And The PHP5 memcache Module](/pdf/installing_memcached_and_the_php5_memcache_module.pdf)

## Other

To get statistics on memcached and calculate hits and ratios:

```bash
echo -en "stats\r\n" "quit\r\n" | nc localhost 11211 | tr -s [:cntrl:] " "| cut -f42,48 -d" " | sed "s/\([0-9]*\)\s\([0-9]*\)/ \2\/\1*100/" | bc -l
```
