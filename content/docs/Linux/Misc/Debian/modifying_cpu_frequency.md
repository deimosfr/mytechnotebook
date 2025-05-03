---
weight: 999
url: "/Modifier_la_fréquence_de_son_processeur/"
title: "Modifying CPU Frequency"
description: "How to modify CPU frequency to save power and reduce fan noise on your computer"
categories: ["Linux", "Ubuntu"]
date: "2009-08-02T15:07:00+02:00"
lastmod: "2009-08-02T15:07:00+02:00"
tags: ["CPU", "Power Management", "Laptop Mode", "Ubuntu"]
toc: true
---

## Introduction

So, you have an irritatingly loud CPU fan which is making you consider whether or not launching your laptop through the nearest window is a good idea. And you want to save the planet too! Well, before you do that, why not give CPU frequency scaling a go. Look at laptop mode!

## Installation

You'll need to install the package:

```bash
apt-get install laptop-mode-tools
```

If you're on ubuntu, you can look at **ubuntu-laptop-mode package**.

## Configuration

Modify this line to 1 to always activate laptop mode (`/etc/laptop-mode/laptop-mode.conf`):

```bash
ENABLE_LAPTOP_MODE_ON_AC=1
```

Then you also can modify this line to autoadapt the CPU with your needs (`/etc/laptop-mode/conf.d/cpufreq.conf`):

```bash
CONTROL_CPU_FREQUENCY=1
```

After restart the service:

```bash
/etc/init.d/laptop-mode restart
```

## Advanced Configuration

### Loading Module

OK, first of all, roll up your sleeves and insert the p4_clockmod module:

```bash
sudo modprobe p4_clockmod
```

This shouldn't return any output.

Now, add the line (`/etc/modules`):

```bash
p4_clockmod
```

to ensure the CPU clock scaling module starts with the system.

### CPU frequency scaling monitor

Now, to add the CPU frequency scaling monitor applet to the panel, right click over an empty area in the panel, select 'add to panel', and select the CPU frequency applet. Hopefully it'll pop up showing the CPU frequency now. Reboot your laptop if it doesn't seem to be working immediately.

### Set frequency

Finally, if you (like I did) get miffed off with the laptop 'lagging' when needing a quick boost of power, you can manually set the frequency you want it to run at. Sometimes 250Mhz just isn't enough!

```bash
sudo cpufreq-selector -f 1000000
```

Would set the CPU frequency to 1GHz - Easy eh?

There you have it, CPU frequency scaling in 5 minutes and a good deal cheaper than lobbing your computer through a window!
