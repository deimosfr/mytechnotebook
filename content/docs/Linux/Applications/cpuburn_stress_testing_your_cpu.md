---
weight: 999
url: "/Cpuburn_\\:_stresser_son_CPU/"
title: "CPUBurn: Stress Testing Your CPU"
description: "A guide on how to stress test your CPU in Linux and monitor system statistics during testing."
categories: ["Linux", "Ubuntu", "Servers"]
date: "2013-08-20T13:14:00+02:00"
lastmod: "2013-08-20T13:14:00+02:00"
tags:
  ["Linux", "CPU", "System Monitoring", "Performance Testing", "Stress Testing"]
toc: true
---

There are plenty of system stress testing applications for Windows, but what about Linux? Here is a simple way to stress test your CPU in Linux and monitor various system statistics while doing it. This can be used for testing an overclocked system or just to burn in a new CPU.

## Installing CPUBurn

First install cpuburn:

```bash
aptitude install cpuburn
```

You can run CPU Burn-In with:

```bash
burnP6
```

## Using Stress

You can also use stress:

```bash
aptitude install stress
```

Then launch it:

```bash
stress --cpu 2 --io 1 --vm 1 --vm-bytes 128M --timeout 10s --verbose
```

## Monitoring Tools

Now for a few diagnostic tools. First install lm-sensors. This program can read various sensor chips in your system and report their outputs. This gives you access to things like CPU temperature, core voltage, and fan speeds. In Ubuntu you can install lm-sensors with:

```bash
sudo apt-get install lm-sensors
```

Once you have done this you have to configure lm-sensors so it knows about the sensors in your system. To do this run:

```bash
sudo /usr/sbin/sensors-detect
```

In most cases you can choose the default answer at each of the programs prompts. When you get to the end, the program will tell you that you need to add some lines to several system files. In my case I was asked to add the following:

```bash
To make the sensors modules behave correctly, add these lines to
/etc/modules:

#----cut here----
# I2C adapter drivers
# modprobe unknown adapter NVIDIA i2c adapter 0 at 5:00.0
# modprobe unknown adapter NVIDIA i2c adapter 1 at 5:00.0
# modprobe unknown adapter NVIDIA i2c adapter 2 at 5:00.0
# modprobe unknown adapter NVIDIA i2c adapter 0 at 4:00.0
# modprobe unknown adapter NVIDIA i2c adapter 1 at 4:00.0
# modprobe unknown adapter NVIDIA i2c adapter 2 at 4:00.0
i2c-nforce2
# Chip drivers
eeprom
k8temp
w83627hf
#----cut here----
```

At this point either restart so that the new kernel modules get loaded or load them by hand with:

```bash
sudo /sbin/modprobe <module>
```

If you load the modules by hand you do not need to restart. Now you can run lm-sensors with:

```bash
sensors
```

and you will get an output like this:

```bash
k8temp-pci-00c3
Adapter: PCI adapter
Core0 Temp:
             +42°C
Core1 Temp:
             +37°C

w83627thf-isa-0290
Adapter: ISA adapter
VCore:     +1.13 V  (min =  +0.70 V, max =  +1.87 V)
+12V:     +12.77 V  (min = +14.53 V, max =  +6.81 V)
+3.3V:     +3.20 V  (min =  +0.91 V, max =  +4.02 V)
+5V:       +5.04 V  (min =  +3.33 V, max =  +2.91 V)
-12V:     -12.03 V  (min =  -7.34 V, max =  -1.75 V)
V5SB:      +5.08 V  (min =  +0.30 V, max =  +1.67 V)
VBat:      +3.02 V  (min =  +0.86 V, max =  +2.10 V)
fan1:        0 RPM  (min = 11065 RPM, div = 2)
CPU Fan:  4821 RPM  (min = 84375 RPM, div = 8)
fan3:     6887 RPM  (min = 6887 RPM, div = 2)
M/B Temp:    +41°C  (high =  +127°C, hyst =   +32°C)   sensor = thermistor
CPU Temp:  +41.5°C  (high =   +80°C, hyst =   +75°C)   sensor = thermistor
temp3:     -98.5°C  (high =   +80°C, hyst =   +75°C)   sensor = diode
vid:      +0.000 V  (VRM Version 2.4)
alarms:   Chassis intrusion detection
beep_enable:    Sound alarm enabled
```

You can further modify the output of lm-sensors by changing the `/etc/sensors.conf` file, but I will not get into that here. lm-sensors only displays its output once each time you call it, so it only gives you a static image of system statistics. A simple python script will enable you to view the sensor outputs continuously.

```python
#!/usr/bin/python

import os, time

while 1:
        time.sleep(3)
        os.system('clear')
        os.system('sensors')
```

Note: you can do the same with this in your terminal:

```bash
watch sensors
```

This will update the sensor output every three seconds. If you have a CPU that does frequency scaling like an AMD Athlon 64, you can also output the current CPU frequency by adding these lines at the end of the while loop:

```python
fid = open('/sys/devices/system/cpu/cpu0/cpufreq/scaling_cur_freq', 'rb')
print 'CPU0 Frequency:', float(fid.read()) / 1000000, 'GHz'
fid.close()
```

If you have more than one CPU core just copy that piece of code multiple times replacing "cpu0" with "cpuN". Save the script as something like sensors.py and run it with:

```bash
python sensors.py
```

Just hit Ctrl-c to stop the script. Now its time to actually stress the system. I open several terminals, one for each CPU core I have in my system, and run one instance of CPU Burn-In in each of them. This should hopefully get the CPUs running at 100 percent. Next open another terminal and run the python sensor script so you can keep track of temperatures and core frequencies and whatever else lm-sensors is outputting.

If you are using Gnome go to System-->Administration-->"System Monitor" on the top menu bar. This will open System Monitor which is much like the Windows Task Manager and will graphically display various system information. In particular for stress testing, the Resource tab in System Monitor is helpful because it shows the CPU load.

I usually run the burn in program for a good length of time, like 10 or 12 hours while I am sleeping. Hopefully, when CPU Burn-In finishes it reports a message that no errors were found. If not, reduce your overclock or get some better CPU cooling or maybe your CPU is just defective. Whatever the case, that should be everything you need to stress test the CPU in Linux and view some diagnostics while you do so.
