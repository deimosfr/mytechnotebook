---
weight: 999
url: "/Bwm-ng_\\:_Mesurer_la_consommation_de_bande_passante_en_temps_réel/"
title: "Bwm-ng: Measure Bandwidth Consumption in Real Time"
description: "How to use bwm-ng to monitor network bandwidth usage in real-time on Linux systems."
categories: ["Linux", "Network", "Monitoring"]
date: "2009-11-28T16:14:00+02:00"
lastmod: "2009-11-28T16:14:00+02:00"
tags: ["bwm-ng", "bandwidth", "monitoring", "network"]
toc: true
---

## Introduction

Sometimes, you need to measure the current total bandwidth. bwm-ng is your friend :-)

## Installation

On Debian, it's easy:

```bash
apt-get install bwm-ng
```

## Usage

Simply run the command to get it working directly:

```bash
$ bwm-ng
  bwm-ng v0.6 (probing every 0.500s), press 'h' for help
  input: /proc/net/dev type: rate
  |         iface                   Rx                   Tx                Total
  ==============================================================================
               lo:           0.00 KB/s            0.00 KB/s            0.00 KB/s
             eth0:        2275.89 KB/s           57.56 KB/s         2333.45 KB/s
  ------------------------------------------------------------------------------
            total:        2275.89 KB/s           57.56 KB/s         2333.45 KB/s
```
