---
weight: 999
url: "/Connaitre_le_nombre_de_cores_CPU_actifs_sur_Solaris/"
title: "How to Check the Number of Active CPU Cores on Solaris"
description: "How to verify the number of active CPU cores on Solaris systems, which is useful for licensing and resource allocation"
categories: ["Solaris", "System Administration"]
date: "2012-02-06T10:10:00+02:00"
lastmod: "2012-02-06T10:10:00+02:00"
tags: ["Solaris", "CPU", "cores", "hardware", "psrinfo"]
toc: true
---

## Introduction

On certain machines (especially HP), it's very practical for licensing issues to be able to limit the number of cores. Afterwards, it's important to verify that the state is what you expect.

## Usage

To get the list of active cores:

```bash {linenos=table,hl_lines=[2,4]}
> psrinfo -pv
The physical processor has 1 virtual processor (0)  x86 (chipid 0x0 GenuineIntel family 6 model 44 step 2 clock 2800 MHz)
        Intel(r) Xeon(r) CPU           X5660  @ 2.80GHz
The physical processor has 1 virtual processor (1)  x86 (chipid 0x1 GenuineIntel family 6 model 44 step 2 clock 2800 MHz)
        Intel(r) Xeon(r) CPU           X5660  @ 2.80GHz
```

Here we can see that I only have one virtual CPU per processor (so 2 total).

## Resources
- http://serverfault.com/questions/85478/sun-solaris-find-out-number-of-processors-and-cores
