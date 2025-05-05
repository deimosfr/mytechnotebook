+++
weight = 999
title = 'HyperHDR with RPI 3: Swapped colors'
description = 'Get your true colors with Raspberry Pi 3 and HyperHDR'
categories = ['Hyperhdr', 'Raspberry Pi']
tags = ['Hyperhdr', 'Raspberry Pi']
toc = true
date = '2025-05-04T15:08:37+02:00'
+++

I really love the [HyperHDR](https://github.com/awawa-dev/HyperHDR) project. I use it to control my LED strips on my TV. I've even [designed a 3D case](https://www.printables.com/model/1072533-ambilight-case-for-hyperhdrhyperion) to put my RPI 3 and all the requirements in one place.

But the problem is that I have a RPI 3 and the colors are swapped. It looks to be a kernel issue and the only solution I've found is to revert to a working kernel:

```bash
sudo rpi-update 5fc4f643d2e9c5aa972828705a902d184527ae3f
```

## Resources

- https://github.com/awawa-dev/HyperHDR/discussions/848
