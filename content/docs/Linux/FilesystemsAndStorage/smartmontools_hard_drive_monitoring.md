---
weight: 999
url: "/Smartmontool_\\:_Surveillance_des_disques_dur/"
title: "Smartmontools: Hard Drive Monitoring"
description: "How to install and configure Smartmontools to monitor hard drive health and performance on various Linux and BSD systems."
categories: ["Linux", "Debian", "BSD", "System Administration", "Monitoring"]
date: "2011-01-14T20:24:00+02:00"
lastmod: "2011-01-14T20:24:00+02:00"
tags: ["smartmontools", "monitoring", "hard drive", "FreeBSD", "Debian", "diagnostics", "maintenance", "system health"]
toc: true
---

## Introduction

Smartmontools is a tool for analyzing hard drives and their most critical physical characteristics. It consists of two parts: smartd daemon, which checks parameters every 30 minutes and writes the results to `/var/log/syslog`, and the smartctl command which requires root privileges and is used to display all the information.

## Activation / Installation of smartmontools

### Debian

Installation requires root privileges. The package name varies depending on your Debian version. The example below is for Sarge.

```bash
> aptitude install smartmontools
Lecture des listes de paquets... Fait
Construction de l'arbre des dependances... Fait
Les NOUVEAUX paquets suivants seront installes :
 smartmontools
0 mis a jour, 1 nouvellement installes, 0 a enlever et 60 non mis a jour.
Il est necessaire de prendre 222ko dans les archives.
Apres depaquetage, 508ko d'espace disque supplementaires seront utilises.
Reception de : 1 http://ftp.fr.debian.org unstable/main smartmontools 5.32-3 [222kB]
222ko receptionnes en 0s (272ko/s)
Selection du paquet smartmontools precedemment deselectionne.
(Lecture de la base de données... 67466 fichiers et repertoires deja installes.)
Depaquetage de smartmontools (a partir de .../smartmontools_5.32-3_i386.deb) ...
Parametrage de smartmontools (5.32-3) ...
Not starting S.M.A.R.T. daemon smartd, disabled via /etc/default/smartmontools
```

As you can see, the daemon has not been started immediately. You need to edit `/etc/default/smartmontools` and uncomment the lines `start_smartd=yes` and `smartd_opts="--interval=1800"`:

```bash
# Defaults for smartmontools initscript (/etc/init.d/smartmontools)
# This is a POSIX shell fragment

# list of devices you want to explicitly enable S.M.A.R.T. for
# not needed if the device is monitored by smartd
# enable_smart="/dev/hda /dev/hdb"

# uncomment to start smartd on system startup
start_smartd=yes

# uncomment to pass additional options to smartd on startup
smartd_opts="--interval=1800"
```

Once the changes are validated, start the daemon:

```bash
/etc/init.d/smartmontools start
Enabling S.M.A.R.T. for: /dev/hda /dev/hdb.
Starting S.M.A.R.T. daemon: smartd.
23:21 root@revolution /# smartctl -a /dev/hda
smartctl version 5.32 Copyright (C) 2002-4 Bruce Allen
```

The smartd daemon will now regularly check your disk information and record it in your logs:

```bash
cat /var/log/syslog | grep smartd
Mar 17 10:48:34 slut smartd[990]: Configuration file /etc/smartd.conf was parsed, found DEVICESCAN, scanning devices
```

And there you go, it's ready.

### FreeBSD

To install smartmontools:

```bash
pkg_add -r smartmontool
```

Then start it like this:

```bash
/usr/local/etc/rc.d/smartd start
```

Edit the `/usr/local/etc/smartd.conf` configuration and add this line (adapting to your email):

```bash
DEVICESCAN -a -m my@mail.com
```

Next, if we want smartd to start at every boot, add this line:

```bash
smartd_enable="YES"
```

## Fine Tuning

### Debian

To fine-tune the smartmontools configuration, edit the `/etc/smartd.conf` file and look for the DEVICESCAN line to add your own settings, as in this example:

```bash
DEVICESCAN -H -l error -l selftest -t -f -m admin@webank.fr -M exec /usr/bin/mail -s (S/../.././02|L/../../6/03)
```

The DEVICESCAN directive indicates that you want to apply this configuration to all hard disks detected as SMART compatible on the system. It can be replaced by the name of a device `/dev/hdx` or `/dev/sdx`.

```bash
/dev/hda -H -l error -l selftest -t -f -m admin@webank.fr -M exec /usr/bin/mail -s (S/../.././02|L/../../6/03)
/dev/hdc -H -l error -l selftest -t -f -m admin@webank.fr -M exec /usr/bin/mail -s (S/../.././02|L/../../6/03)
```

Adding this line to the configuration file allows sending an email to admin@domain.com using your system's mail command. The -t option indicates that we want to be informed in case the "Pre-Fail" or "Old-age" attribute shows errors, if the health test (option -H) fails, or if the error and selftest logs evolve (-l). You can choose from a range of options to best adjust according to your needs. For example, you can deliberately ignore an attribute using the -I option. Adding the -I 194 option indicates that we want to receive an email in case of failure but ignoring attribute number 194 (temperature). The -s option allows you to define the periodicity of the tests to be performed (version >5.30 required). In this example, we perform a short test (S/) every day at 2 a.m., and a long test every Saturday at 3 a.m. It's also possible to modify the email that will be sent by smartd in case of failure by creating a script that will be called instead of /bin/mail.

### FreeBSD

To receive daily emails indicating the state of your disks, add this to the `/etc/periodic.conf` file:

```bash
daily_status_smart_devices="/dev/ad4 /dev/ad6 /dev/ad8 /dev/ad10 /dev/ad12"
```

Obviously, use your own devices.

## Diagnostics and Troubleshooting

Since smartd writes to `/var/log/syslog`, it's easy to search with a grep command as in the following example:

```bash
> grep smartd /var/log/syslog
Mar 17 10:48:34 slut smartd[990]: Configuration file /etc/smartd.conf was parsed, found DEVICESCAN, scanning devices
Mar 17 10:48:34 slut smartd[990]: Device: /dev/hda, opened
Mar 17 10:48:34 slut smartd[990]: Device: /dev/hda, found in smartd database.
Mar 17 10:48:35 slut smartd[990]: Device: /dev/hda, is SMART capable. Adding to "monitor" list.
Mar 17 10:48:35 slut smartd[990]: Device: /dev/hdb, opened
Mar 17 10:48:35 slut smartd[990]: Device: /dev/hdb, not ATA, no IDENTIFY DEVICE Structure
Mar 17 10:48:35 slut smartd[990]: Monitoring 1 ATA and 0 SCSI devices
Mar 17 10:48:35 slut smartd: Lancement smartd succeeded
Mar 17 10:48:35 slut smartd[2421]: smartd has fork()ed into background mode. New PID=2421.
Mar 17 13:48:35 slut smartd[2421]: Device: /dev/hda, SMART Prefailure Attribute: 8 Seek_Time_Performance changed from 246 to 247
Mar 17 15:48:35 slut smartd[2421]: Device: /dev/hda, SMART Prefailure Attribute: 8 Seek_Time_Performance changed from 247 to 246
Mar 17 17:18:35 slut smartd[2421]: Device: /dev/hda, SMART Prefailure Attribute: 8 Seek_Time_Performance changed from 246 to 247
```

How to interpret these lines? The drive shows a constant value that varies between 246 and 247. If the value suddenly changes from 247 to 500, this is abnormal behavior.

Using the smartctl command requires root privileges. Let's look at the different attributes of the command.

```bash
smarctl -h

smartctl version 5.33 [i386-redhat-linux-gnu] Copyright (C) 2002-4 Bruce Allen
Home page is http://smartmontools.sourceforge.net/[1]
Usage: smartctl [options] device
h, --help, --usage
Display this help and exit
i, --info
Show identity information for device
a, --all
Show all SMART information for device
```

```bash
smartctl -i /dev/hda

=== START OF INFORMATION SECTION ===
Device Model:     Maxtor 6E040L0
Serial Number:    E1KTPXFE
Firmware Version: NAR61590
User Capacity:    41,110,142,976 bytes
Device is:        In smartctl database [for details use: -P show]
ATA Version is:   7
ATA Standard is:  ATA/ATAPI-7 T13 1532D revision 0
Local Time is:    Thu Mar 17 22:21:52 2005 CET
SMART support is: Available - device has SMART capability.
SMART support is: Enabled
```

```bash
smartctl -a /dev/hda

=== START OF READ SMART DATA SECTION ===
SMART overall-health self-assessment test result: PASSED

General SMART Values:
Offline data collection status:  (0x82) Offline data collection activity
                                     was completed without error.
                                     Auto Offline Data Collection: Enabled.
Self-test execution status:      (   0) The previous self-test routine completed
                                     without error or no self-test has ever
                                     been run.
Total time to complete Offline
data collection:                 (1021) seconds.
Offline data collection
capabilities:                    (0x5b) SMART execute Offline immediate.
                                     Auto Offline data collection on/off support.
                                     Suspend Offline collection upon new
                                     command.
                                     Offline surface scan supported.
                                     Self-test supported.
                                     No Conveyance Self-test supported.
                                     Selective Self-test supported.
SMART capabilities:            (0x0003) Saves SMART data before entering
                                     power-saving mode.
                                     Supports SMART auto save timer.
Error logging capability:        (0x01) Error logging supported.
                                     No General Purpose Logging support.
Short self-test routine
recommended polling time:        (   2) minutes.
Extended self-test routine
recommended polling time:        (  17) minutes.

SMART Attributes Data Structure revision number: 16
Vendor Specific SMART Attributes with Thresholds:
ID# ATTRIBUTE_NAME          FLAG     VALUE WORST THRESH TYPE      UPDATED  WHEN_FAILED RAW_VALUE
3 Spin_Up_Time            0x0027   252   252   063    Pre-fail  Always       -       2463
4 Start_Stop_Count        0x0032   253   253   000    Old_age   Always       -       18
5 Reallocated_Sector_Ct   0x0033   253   253   063    Pre-fail  Always       -       0
6 Read_Channel_Margin     0x0001   253   253   100    Pre-fail  Offline      -       0
7 Seek_Error_Rate         0x000a   253   252   000    Old_age   Always       -       0
8 Seek_Time_Performance   0x0027   247   238   187    Pre-fail  Always       -       46214
9 Power_On_Minutes        0x0032   241   241   000    Old_age   Always       -       950h+09m
10 Spin_Retry_Count        0x002b   252   252   157    Pre-fail  Always       -       0
11 Calibration_Retry_Count 0x002b   253   252   223    Pre-fail  Always       -       0
12 Power_Cycle_Count       0x0032   253   253   000    Old_age   Always       -       22
192 Power-Off_Retract_Count 0x0032   253   253   000    Old_age   Always       -       13
193 Load_Cycle_Count        0x0032   253   253   000    Old_age   Always       -       72
194 Temperature_Celsius     0x0032   253   253   000    Old_age   Always       -       31
195 Hardware_ECC_Recovered  0x000a   253   252   000    Old_age   Always       -       25095
196 Reallocated_Event_Count 0x0008   253   253   000    Old_age   Offline      -       0
197 Current_Pending_Sector  0x0008   253   253   000    Old_age   Offline      -       0
198 Offline_Uncorrectable   0x0008   253   253   000    Old_age   Offline      -       0
199 UDMA_CRC_Error_Count    0x0008   199   199   000    Old_age   Offline      -       0
200 Multi_Zone_Error_Rate   0x000a   253   252   000    Old_age   Always       -       0
201 Soft_Read_Error_Rate    0x000a   251   138   000    Old_age   Always       -       1746
202 TA_Increase_Count       0x000a   253   252   000    Old_age   Always       -       0
203 Run_Out_Cancel          0x000b   253   252   180    Pre-fail  Always       -       137
204 Shock_Count_Write_Opern 0x000a   253   252   000    Old_age   Always       -       0
205 Shock_Rate_Write_Opern  0x000a   253   252   000    Old_age   Always       -       0
207 Spin_High_Current       0x002a   252   252   000    Old_age   Always       -       0
208 Spin_Buzz               0x002a   252   252   000    Old_age   Always       -       0
209 Offline_Seek_Performnce 0x0024   187   183   000    Old_age   Offline      -       0
99 Unknown_Attribute       0x0004   253   253   000    Old_age   Offline      -       0
100 Unknown_Attribute       0x0004   253   253   000    Old_age   Offline      -       0
101 Unknown_Attribute       0x0004   253   253   000    Old_age   Offline      -       0

SMART Error Log Version: 1
No Errors Logged

SMART Self-test log structure revision number 1
No self-tests have been logged.  [To run self-tests, use: smartctl -t]

SMART Selective self-test log data structure revision number 1
SPAN  MIN_LBA  MAX_LBA  CURRENT_TEST_STATUS
  1        0        0  Not_testing
  2        0        0  Not_testing
  3        0        0  Not_testing
  4        0        0  Not_testing
  5        0        0  Not_testing
Selective self-test flags (0x0):
After scanning selected spans, do NOT read-scan remainder of disk.
If Selective self-test is pending on power-up, resume after 0 minute delay.
```

Now we need to interpret the information such as disk uptime, temperature, and most importantly for us, errors. For this we mainly observe the last two columns: WHEN_FAILED and RAW_VALUE, and the section just below: SMART Error Log Version: 1 No Errors Logged.

An example:

```
5 Reallocated_Sector_Ct 0x0033 016 016 063 Pre-fail Always FAILING_NOW 598
```

Here we see that sector reallocation has failed. You should therefore monitor this part. If the number indicated quickly increases to higher figures, take the necessary measures: back up your data and possibly contact support.

## Conclusion

Smartmontools is simple to use and very comprehensive. Note, however, that such a tool does not replace the most important thing: regular backup of your data.

## FAQ

### Problems during updates

Sometimes during a package update, things may go wrong and you may not know why. The problem is actually quite simple. Just stop the service:

```bash
/etc/init.d/smartmontool stop
```

then restart the update.

### The service won't start

This problem can occur when SMART is simply not enabled. To enable it, just type this command:

```bash
smartctl -s on /dev/sda
```

Then try to start smartmontools:

```bash
/etc/init.d/smartmontools start
```

## Resources
- [Checking Hard Disk Sanity With Smartmontools](/pdf/checking_hard_disk_sanity_with_smartmontools.pdf)
- http://www.davidandrzejewski.com/2009/03/15/freebsd-monitor-your-disks-health-with-smartmontools/
