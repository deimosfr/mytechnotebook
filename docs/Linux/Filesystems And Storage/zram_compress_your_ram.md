---
title: "ZRam: Compress Your RAM"
slug: zram-compress-your-ram/
description: "Learn how to use ZRam to create compressed block devices in RAM instead of swap for better system performance."
categories: ["Linux", "Filesystems And Storage"]
tags: ["zram", "compression", "performance", "debian"]
---

# ZRam: Compress Your RAM

ZRam is a Linux kernel module that creates a compressed block device in RAM. Think of it as a way to "squeeze" more data into your physical memory. Instead of swapping out to a slow disk when RAM gets full, Linux can swap into this compressed ZRam device, which is significantly faster.

This is especially great for systems with limited RAM, as it trades a little bit of CPU for a lot more usable memory space.

## Installation

Debian makes this incredibly easy with the `zram-tools` package. It handles all the setup scripts for us.

We will install the package:

```bash
sudo apt-get install zram-tools
```

## Configuration

By default, the package sets up a basic configuration, but let's tweak it to get the best performance. We want to use the `zstd` compression algorithm (which offers a great balance of speed and compression ratio) and dedicate a healthy chunk of RAM to it.

We will edit the configuration file `/etc/default/zramswap`:

=== "/etc/default/zramswap"

    ```bash
    # Compression algorithm selection
    # speed: lz4 > zstd > lzo
    # compression: zstd > lzo > lz4
    # This is not inclusive of all available algorithms.
    # See /sys/block/zram0/comp_algorithm (when zram module is loaded) to see
    # what your kernel provides.
    ALGO=zstd

    # Specifies the amount of RAM that should be used for zram
    # based on a percentage the total amount of available memory
    PERCENT=60

    # Specifies a static amount of RAM that should be used for zram
    # This is in MiB
    #SIZE=256

    # Specifies the priority for the swap devices
    # See swapon(8) for more details
    #PRIORITY=100
    ```

Once we've saved our changes, we need to restart the service to apply them:

```bash
sudo systemctl restart zramswap
```

## Verification

Let's check if our ZRam device is up and running. The `zramctl` tool gives us a nice summary.
You should see output similar to this, showing the algorithm (`zstd`) and the disk size:

```bash
$ zramctl
NAME       ALGORITHM DISKSIZE  DATA COMPR TOTAL STREAMS MOUNTPOINT
/dev/zram0 lz4           1.2G 28.8M   18M 18.6M       2 [SWAP]
```

We can also verify that it's being heavily utilized as a high-priority swap device:

```bash
$ swapon --show
NAME       TYPE      SIZE  USED PRIO
/dev/zram0 partition 1.2G 29.3M  100
```

## Resources
* https://wiki.archlinux.org/title/Zram
* https://wiki.debian.org/ZRam